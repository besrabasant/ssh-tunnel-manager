//! Thin tonic adapters that translate protobuf messages to daemon operations.

use super::Service;
use crate::rpc::daemon_service_server::DaemonService;
use crate::rpc::*;
use crate::store::Entry;
use anyhow::Result;
use tonic::{Request, Response, Status};

fn ok() -> i32 {
    // Prost stores enum fields as i32 values in generated response structures.
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
        // Match characters in order to preserve the former fuzzy-search behavior.
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
