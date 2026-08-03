//! gRPC channel construction for CLI and interactive operations.

use crate::rpc::daemon_service_client::DaemonServiceClient;
use anyhow::{Context, Result};
use tonic::transport::{Channel, Endpoint};

pub(super) async fn connect() -> Result<DaemonServiceClient<Channel>> {
    // Bound connection and request timeouts keep scripts from hanging on a dead daemon.
    let endpoint = Endpoint::from_static("http://127.0.0.1:50051")
        .connect_timeout(std::time::Duration::from_secs(10))
        .timeout(std::time::Duration::from_secs(20));
    Ok(DaemonServiceClient::connect(endpoint)
        .await
        .context("did not connect to sshtmd")?)
}
