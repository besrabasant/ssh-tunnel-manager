use anyhow::Result;
use sshtm::daemon::Service;
use sshtm::rpc::daemon_service_server::DaemonServiceServer;
use sshtm::store::Store;
use tonic::transport::Server;

#[tokio::main]
async fn main() -> Result<()> {
    // Construct one shared service instance for restoration, RPCs, and shutdown.
    let service = Service::new(Store::new(sshtm::config_dir())?);
    service.restore().await?;
    eprintln!("server listening at {}", sshtm::ADDRESS);
    Server::builder()
        .add_service(DaemonServiceServer::new(service.clone()))
        .serve_with_shutdown(sshtm::ADDRESS.parse()?, shutdown_signal())
        .await?;
    // gRPC has stopped accepting work, so it is now safe to snapshot and stop children.
    service.shutdown().await
}

#[cfg(unix)]
async fn shutdown_signal() {
    use tokio::signal::unix::{SignalKind, signal};
    let mut terminate = signal(SignalKind::terminate()).expect("install SIGTERM handler");
    // User services normally send SIGTERM, while foreground runs commonly use Ctrl-C.
    tokio::select! {
        _ = tokio::signal::ctrl_c() => {},
        _ = terminate.recv() => {},
    }
}

#[cfg(not(unix))]
async fn shutdown_signal() {
    tokio::signal::ctrl_c().await.ok();
}
