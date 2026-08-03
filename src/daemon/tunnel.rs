//! OpenSSH child-process lifecycle and active-tunnel persistence.

use super::{RunningTunnel, SavedTunnel, Service};
use crate::rpc::ActiveTunnel;
use crate::store::{Store, atomic_json};
use anyhow::{Context, Result, bail};
use chrono::Utc;
use std::collections::HashMap;
use std::fs;
use std::net::TcpListener;
use std::sync::Arc;
use tokio::process::Command;
use tokio::sync::Mutex;
use tokio::time::{Duration, sleep};

impl Service {
    pub fn new(store: Store) -> Self {
        Self {
            store,
            tunnels: Arc::new(Mutex::new(HashMap::new())),
        }
    }

    pub async fn restore(&self) -> Result<()> {
        // A missing state file simply means the previous shutdown had no active tunnels.
        let path = self.store.active_path();
        let saved: Vec<SavedTunnel> = match fs::read(&path) {
            Ok(data) => serde_json::from_slice(&data).context("parse active_tunnels.json")?,
            Err(err) if err.kind() == std::io::ErrorKind::NotFound => return Ok(()),
            Err(err) => return Err(err.into()),
        };
        for tunnel in saved {
            // Restore independently so one stale profile does not prevent other tunnels.
            if let Err(error) = self.start(&tunnel.config_name, tunnel.local_port).await {
                eprintln!("failed to restore tunnel {}: {error:#}", tunnel.config_name);
            }
        }
        Ok(())
    }

    pub async fn shutdown(&self) -> Result<()> {
        // Persist before killing children; the next daemon process will recreate them.
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
            // Absence, rather than an empty array, matches the previous Go behavior.
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

    pub(super) async fn start(&self, name: &str, requested_port: i32) -> Result<Vec<String>> {
        // Remove processes that exited since the last registry operation.
        self.reap().await;
        let entry = self.store.get(name)?;
        // -1 means "use the profile default"; zero asks the OS for a free port.
        let port = if requested_port == -1 {
            entry.local_port
        } else {
            requested_port
        };
        let port = if port == 0 {
            // Bind port zero briefly to discover a currently available ephemeral port.
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

        // OpenSSH expects its port separately, while profiles store host[:port].
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
        // -N/-T create forwarding only. Batch mode prevents a daemon-side password prompt.
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
        // ExitOnForwardFailure makes early termination a reliable setup failure signal.
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
        // Persist only after OpenSSH has survived its setup window and entered the registry.
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

    pub(super) async fn stop(&self, name: &str, port: i32) -> Result<String> {
        // A textual identifier selects by profile name; a numeric identifier selects by port.
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
        // Remove first so concurrent list/start requests cannot observe a stopping tunnel.
        tunnel.child.kill().await.context("stop SSH process")?;
        self.persist().await?;
        Ok(format!(
            "Closing existing connection on port {target} for {}",
            tunnel.entry.name
        ))
    }

    pub(super) async fn active(&self) -> Vec<ActiveTunnel> {
        // Avoid reporting children that have already exited due to network or auth failure.
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
        // try_wait is non-blocking; live children return Ok(None) and stay registered.
        self.tunnels
            .lock()
            .await
            .retain(|_, tunnel| matches!(tunnel.child.try_wait(), Ok(None)));
    }

    pub(super) async fn active_port_for(&self, name: &str) -> Option<i32> {
        self.tunnels
            .lock()
            .await
            .values()
            .find(|t| t.entry.name == name)
            .map(|t| t.port)
    }
}

fn split_server(server: &str) -> Result<(String, String)> {
    // Bracketed IPv6 addresses must be split after the closing bracket, not on ':'.
    if let Some(value) = server.strip_prefix('[') {
        let (host, port) = value.rsplit_once("]:").context("bad SSH server address")?;
        return Ok((host.into(), port.into()));
    }
    let (host, port) = server.rsplit_once(':').context("bad SSH server address")?;
    Ok((host.into(), port.into()))
}
