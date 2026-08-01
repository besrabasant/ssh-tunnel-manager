use crate::rpc::daemon_service_server::DaemonService;
use crate::rpc::*;
use crate::store::{Entry, Store, atomic_json};
use anyhow::{Context, Result, bail};
use chrono::Utc;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::fs;
use std::net::TcpListener;
use std::sync::Arc;
use tokio::process::{Child, Command};
use tokio::sync::Mutex;
use tokio::time::{Duration, sleep};
use tonic::{Request, Response, Status};

#[derive(Clone)]
pub struct Service {
    store: Store,
    tunnels: Arc<Mutex<HashMap<i32, RunningTunnel>>>,
}

struct RunningTunnel {
    entry: Entry,
    port: i32,
    child: Child,
}

#[derive(Clone, Debug, Deserialize, Serialize)]
struct SavedTunnel {
    config_name: String,
    local_port: i32,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    local_addr: String,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    remote_addr: String,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    server: String,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    user: String,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    saved_at: String,
}

impl Service {
    pub fn new(store: Store) -> Self {
        Self {
            store,
            tunnels: Arc::new(Mutex::new(HashMap::new())),
        }
    }

    pub async fn restore(&self) -> Result<()> {
        let path = self.store.active_path();
        let saved: Vec<SavedTunnel> = match fs::read(&path) {
            Ok(data) => serde_json::from_slice(&data).context("parse active_tunnels.json")?,
            Err(err) if err.kind() == std::io::ErrorKind::NotFound => return Ok(()),
            Err(err) => return Err(err.into()),
        };
        for tunnel in saved {
            if let Err(error) = self.start(&tunnel.config_name, tunnel.local_port).await {
                eprintln!("failed to restore tunnel {}: {error:#}", tunnel.config_name);
            }
        }
        Ok(())
    }

    pub async fn shutdown(&self) -> Result<()> {
        self.persist().await?;
        let mut tunnels = self.tunnels.lock().await;
        for tunnel in tunnels.values_mut() {
            let _ = tunnel.child.kill().await;
        }
        tunnels.clear();
        Ok(())
    }

    async fn persist(&self) -> Result<()> {
        let tunnels = self.tunnels.lock().await;
        if tunnels.is_empty() {
            match fs::remove_file(self.store.active_path()) {
                Ok(()) => {}
                Err(error) if error.kind() == std::io::ErrorKind::NotFound => {}
                Err(error) => return Err(error.into()),
            }
            return Ok(());
        }
        let saved_at = Utc::now().to_rfc3339();
        let saved = tunnels
            .values()
            .map(|tunnel| SavedTunnel {
                config_name: tunnel.entry.name.clone(),
                local_port: tunnel.port,
                local_addr: format!("127.0.0.1:{}", tunnel.port),
                remote_addr: format!("{}:{}", tunnel.entry.remote_host, tunnel.entry.remote_port),
                server: tunnel.entry.server.clone(),
                user: tunnel.entry.user.clone(),
                saved_at: saved_at.clone(),
            })
            .collect::<Vec<_>>();
        atomic_json(&self.store.active_path(), &saved)
    }

    async fn start(&self, name: &str, requested_port: i32) -> Result<Vec<String>> {
        self.reap().await;
        let entry = self.store.get(name)?;
        let port = if requested_port == -1 {
            entry.local_port
        } else {
            requested_port
        };
        let port = if port == 0 {
            let listener = TcpListener::bind("127.0.0.1:0")?;
            listener.local_addr()?.port() as i32
        } else {
            port
        };
        if !(1..=65535).contains(&port) {
            bail!("invalid local port {port}");
        }
        if self.tunnels.lock().await.contains_key(&port) {
            bail!("Cannot start tunnel as connection is already open on port {port}");
        }
        TcpListener::bind(("127.0.0.1", port as u16))
            .with_context(|| format!("local port {port} is unavailable"))?;

        let server = if entry.server.contains(':') {
            entry.server.clone()
        } else {
            format!("{}:22", entry.server)
        };
        let (host, ssh_port) = split_server(&server)?;
        let forward = format!(
            "127.0.0.1:{port}:{}:{}",
            entry.remote_host, entry.remote_port
        );
        let destination = format!("{}@{host}", entry.user);
        let mut command = Command::new("ssh");
        command
            .kill_on_drop(true)
            .args([
                "-N",
                "-T",
                "-L",
                &forward,
                "-i",
                &entry.key_file,
                "-p",
                &ssh_port,
            ])
            .args([
                "-o",
                "BatchMode=yes",
                "-o",
                "ExitOnForwardFailure=yes",
                "-o",
                "ConnectTimeout=10",
                "-o",
                "StrictHostKeyChecking=no",
                "-o",
                "UserKnownHostsFile=/dev/null",
            ])
            .arg(destination)
            .stdin(std::process::Stdio::null())
            .stdout(std::process::Stdio::null())
            .stderr(std::process::Stdio::piped());
        let mut child = command.spawn().context("couldn't execute OpenSSH client")?;
        sleep(Duration::from_millis(250)).await;
        if let Some(status) = child.try_wait()? {
            let stderr = child.stderr.take();
            let detail = if let Some(mut stderr) = stderr {
                use tokio::io::AsyncReadExt;
                let mut text = String::new();
                stderr.read_to_string(&mut text).await.ok();
                text.trim().to_string()
            } else {
                String::new()
            };
            bail!("SSH exited with {status}: {detail}");
        }
        self.tunnels.lock().await.insert(
            port,
            RunningTunnel {
                entry: entry.clone(),
                port,
                child,
            },
        );
        self.persist().await?;
        Ok(vec![
            "Starting tunnel setup...".into(),
            format!("Connected to {server}"),
            "Local listener set up, ready to accept connections.".into(),
            format!(
                "Tunneling \"127.0.0.1:{port}\" <==> \"{}:{}\" through \"{server}\"",
                entry.remote_host, entry.remote_port
            ),
        ])
    }

    async fn stop(&self, name: &str, port: i32) -> Result<String> {
        let target = {
            let tunnels = self.tunnels.lock().await;
            if !name.is_empty() {
                tunnels
                    .iter()
                    .find(|(_, tunnel)| tunnel.entry.name == name)
                    .map(|(port, _)| *port)
            } else if tunnels.contains_key(&port) {
                Some(port)
            } else {
                None
            }
        };
        let Some(target) = target else {
            return Ok(if !name.is_empty() {
                format!("Did not find any connection for configuration {name}")
            } else {
                format!("Did not find any connection on port {port}")
            });
        };
        let mut tunnel = self
            .tunnels
            .lock()
            .await
            .remove(&target)
            .expect("target exists");
        tunnel.child.kill().await.context("stop SSH process")?;
        self.persist().await?;
        Ok(format!(
            "Closing existing connection on port {target} for {}",
            tunnel.entry.name
        ))
    }

    async fn active(&self) -> Vec<ActiveTunnel> {
        self.reap().await;
        self.tunnels
            .lock()
            .await
            .values()
            .map(|tunnel| ActiveTunnel {
                name: tunnel.entry.name.clone(),
                local_port: tunnel.port,
                remote_addr: format!("{}:{}", tunnel.entry.remote_host, tunnel.entry.remote_port),
                local_addr: format!("127.0.0.1:{}", tunnel.port),
                server: tunnel.entry.server.clone(),
            })
            .collect()
    }

    async fn reap(&self) {
        self.tunnels
            .lock()
            .await
            .retain(|_, tunnel| matches!(tunnel.child.try_wait(), Ok(None)));
    }

    async fn active_port_for(&self, name: &str) -> Option<i32> {
        self.tunnels
            .lock()
            .await
            .values()
            .find(|t| t.entry.name == name)
            .map(|t| t.port)
    }
}

fn split_server(server: &str) -> Result<(String, String)> {
    if let Some(value) = server.strip_prefix('[') {
        let (host, port) = value.rsplit_once("]:").context("bad SSH server address")?;
        return Ok((host.into(), port.into()));
    }
    let (host, port) = server.rsplit_once(':').context("bad SSH server address")?;
    Ok((host.into(), port.into()))
}

fn ok() -> i32 {
    ResponseStatus::Success as i32
}
fn error() -> i32 {
    ResponseStatus::Error as i32
}
fn internal(error: anyhow::Error) -> Status {
    Status::internal(format!("{error:#}"))
}

#[tonic::async_trait]
impl DaemonService for Service {
    async fn list_configurations(
        &self,
        request: Request<ListConfigurationsRequest>,
    ) -> Result<Response<ListConfigurationsResponse>, Status> {
        let pattern = request.into_inner().search_pattern.to_lowercase();
        let mut configs = self.store.list().map_err(internal)?;
        if !pattern.is_empty() {
            configs.retain(|entry| fuzzy_match(&pattern, &entry.name.to_lowercase()));
        }
        let result = if configs.is_empty() {
            if pattern.is_empty() {
                "\nNo configurations found\n".into()
            } else {
                format!("\nNo configurations found with search pattern \"{pattern}\" \n")
            }
        } else {
            configs
                .iter()
                .map(format_config)
                .collect::<Vec<_>>()
                .join("\n")
        };
        Ok(Response::new(ListConfigurationsResponse {
            result,
            configs: configs.into_iter().map(Into::into).collect(),
        }))
    }

    async fn list_configurations_json(
        &self,
        _: Request<ListConfigurationsJsonRequest>,
    ) -> Result<Response<ListConfigurationsJsonResponse>, Status> {
        Ok(Response::new(ListConfigurationsJsonResponse {
            configs: self
                .store
                .list()
                .map_err(internal)?
                .into_iter()
                .map(Into::into)
                .collect(),
        }))
    }

    async fn add_configuration(
        &self,
        request: Request<AddOrUpdateConfigurationRequest>,
    ) -> Result<Response<AddOrUpdateConfigurationResponse>, Status> {
        let request = request.into_inner();
        let Some(data) = request.data else {
            return Ok(Response::new(mutation_legacy(
                error(),
                "missing config data",
                None,
            )));
        };
        let entry: Entry = data.clone().into();
        match self.store.add(&entry) {
            Ok(()) => Ok(Response::new(mutation_legacy(
                ok(),
                &format!("Successfully added new configuration {}", entry.name),
                Some(data),
            ))),
            Err(err) => Ok(Response::new(mutation_legacy(
                error(),
                &err.to_string(),
                None,
            ))),
        }
    }

    async fn add_configuration_json(
        &self,
        request: Request<AddOrUpdateConfigurationRequest>,
    ) -> Result<Response<MutationResponse>, Status> {
        let response = self.add_configuration(request).await?.into_inner();
        Ok(Response::new(MutationResponse {
            status: response.status,
            message: response.message,
            data: response.data,
        }))
    }

    async fn update_configuration(
        &self,
        request: Request<AddOrUpdateConfigurationRequest>,
    ) -> Result<Response<AddOrUpdateConfigurationResponse>, Status> {
        let request = request.into_inner();
        if let Some(port) = self.active_port_for(&request.name).await {
            return Ok(Response::new(mutation_legacy(
                error(),
                &format!("Cannot update configuration as connection is open on port {port}"),
                None,
            )));
        }
        let Some(data) = request.data else {
            return Ok(Response::new(mutation_legacy(
                error(),
                "missing config data",
                None,
            )));
        };
        let entry: Entry = data.clone().into();
        match self.store.update(&request.name, &entry) {
            Ok(()) => Ok(Response::new(mutation_legacy(
                ok(),
                &format!("Successfully updated configuration {}", request.name),
                Some(data),
            ))),
            Err(err) => Ok(Response::new(mutation_legacy(
                error(),
                &err.to_string(),
                None,
            ))),
        }
    }

    async fn update_configuration_json(
        &self,
        request: Request<AddOrUpdateConfigurationRequest>,
    ) -> Result<Response<MutationResponse>, Status> {
        let response = self.update_configuration(request).await?.into_inner();
        Ok(Response::new(MutationResponse {
            status: response.status,
            message: response.message,
            data: response.data,
        }))
    }

    async fn fetch_configuration(
        &self,
        request: Request<FetchConfigurationRequest>,
    ) -> Result<Response<FetchConfigurationResponse>, Status> {
        let name = request.into_inner().name;
        let response = match self.store.get(&name) {
            Ok(entry) => FetchConfigurationResponse {
                status: ok(),
                message: String::new(),
                data: Some(entry.into()),
            },
            Err(_) => FetchConfigurationResponse {
                status: error(),
                message: format!("No configurations found with name {name}"),
                data: None,
            },
        };
        Ok(Response::new(response))
    }

    async fn delete_configuration(
        &self,
        request: Request<DeleteConfigurationRequest>,
    ) -> Result<Response<DeleteConfigurationResponse>, Status> {
        let name = request.into_inner().name;
        if let Some(port) = self.active_port_for(&name).await {
            return Err(Status::failed_precondition(format!(
                "cannot delete configuration {name:?} as a connection is open on port {port}"
            )));
        }
        let response = match self.store.remove(&name) {
            Ok(()) => DeleteConfigurationResponse {
                result: format!("Successfully deleted configuration {name}"),
                status: ok(),
                message: format!("Successfully deleted configuration {name}"),
            },
            Err(err) => DeleteConfigurationResponse {
                result: err.to_string(),
                status: error(),
                message: err.to_string(),
            },
        };
        Ok(Response::new(response))
    }

    async fn delete_configuration_json(
        &self,
        request: Request<DeleteConfigurationRequest>,
    ) -> Result<Response<MutationResponse>, Status> {
        let response = self.delete_configuration(request).await?.into_inner();
        Ok(Response::new(MutationResponse {
            status: response.status,
            message: response.message,
            data: None,
        }))
    }

    async fn start_tunnel(
        &self,
        request: Request<StartTunnelRequest>,
    ) -> Result<Response<StartTunnelResponse>, Status> {
        let request = request.into_inner();
        let response = match self.start(&request.config_name, request.local_port).await {
            Ok(events) => StartTunnelResponse {
                result: events.join("\n") + "\n",
                status: ok(),
                message: String::new(),
                events,
            },
            Err(err) => {
                let message = format!("Failed to start tunneling: {err:#}");
                StartTunnelResponse {
                    result: message.clone() + "\n",
                    status: error(),
                    message,
                    events: vec![],
                }
            }
        };
        Ok(Response::new(response))
    }

    async fn kill_tunnel(
        &self,
        request: Request<KillTunnelRequest>,
    ) -> Result<Response<KillTunnelResponse>, Status> {
        let request = request.into_inner();
        let result = self
            .stop(&request.config_name, request.local_port)
            .await
            .map_err(internal)?;
        Ok(Response::new(KillTunnelResponse {
            result: result.clone(),
            status: ok(),
            message: result,
        }))
    }

    async fn list_active_tunnels(
        &self,
        _: Request<ListActiveTunnelsRequest>,
    ) -> Result<Response<ListActiveTunnelsResponse>, Status> {
        let tunnels = self.active().await;
        let result = tunnels
            .iter()
            .map(format_active)
            .collect::<Vec<_>>()
            .join("\n");
        Ok(Response::new(ListActiveTunnelsResponse { result, tunnels }))
    }

    async fn list_active_tunnels_json(
        &self,
        _: Request<ListActiveTunnelsJsonRequest>,
    ) -> Result<Response<ListActiveTunnelsJsonResponse>, Status> {
        Ok(Response::new(ListActiveTunnelsJsonResponse {
            tunnels: self.active().await,
        }))
    }
}

fn mutation_legacy(
    status: i32,
    message: &str,
    data: Option<TunnelConfig>,
) -> AddOrUpdateConfigurationResponse {
    AddOrUpdateConfigurationResponse {
        result: message.into(),
        status,
        message: message.into(),
        data,
    }
}

fn format_config(entry: &Entry) -> String {
    let port = if entry.local_port == 0 {
        "auto".into()
    } else {
        entry.local_port.to_string()
    };
    let local = if entry.local_port == 0 {
        "auto".into()
    } else {
        format!("127.0.0.1:{}", entry.local_port)
    };
    let description = if entry.description.trim().is_empty() {
        String::new()
    } else {
        format!(" ({})", entry.description)
    };
    format!(
        "\n:{port}\n- Connection:                  {}{description}\n- Remote Address:              {}:{}\n- Local Address:               {local}\n- SSH server:                  {}\n- User:                        {}\n- Private key:                 {}\n",
        entry.name, entry.remote_host, entry.remote_port, entry.server, entry.user, entry.key_file
    )
}

fn format_active(tunnel: &ActiveTunnel) -> String {
    format!(
        "\n:{}\n- Connection:                  {}\n- Remote Address:              {}\n- Local Address:               {}\n- SSH server:                  {}\n",
        tunnel.local_port, tunnel.name, tunnel.remote_addr, tunnel.local_addr, tunnel.server
    )
}

fn fuzzy_match(pattern: &str, value: &str) -> bool {
    let mut chars = pattern.chars();
    let mut current = chars.next();
    for ch in value.chars() {
        if current == Some(ch) {
            current = chars.next();
        }
        if current.is_none() {
            return true;
        }
    }
    current.is_none()
}
