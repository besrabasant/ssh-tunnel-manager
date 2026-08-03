//! `sshtm` command-line executable.

#[tokio::main]
async fn main() {
    // Keep process-level error rendering here; command behavior lives in the library.
    if let Err(error) = sshtm::client::run().await {
        eprintln!("{error:#}");
        std::process::exit(1);
    }
}
