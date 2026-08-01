pub mod daemon;
pub mod store;

pub mod rpc {
    // OUT_DIR is populated by build.rs, so generated code stays out of source control.
    tonic::include_proto!("daemon");
}

pub const ADDRESS: &str = "127.0.0.1:50051";
pub const DEFAULT_CONFIG_DIR: &str = ".ssh-tunnel-manager";

pub fn config_dir() -> std::path::PathBuf {
    // Keep the legacy environment variable spelling for compatibility with Go releases.
    if let Ok(path) = std::env::var("config-dir") {
        return expand_home(&path);
    }
    dirs::home_dir()
        .unwrap_or_else(|| std::path::PathBuf::from("."))
        .join(DEFAULT_CONFIG_DIR)
}

pub fn expand_home(path: &str) -> std::path::PathBuf {
    // Expand only the current user's home; forms such as ~other remain literal paths.
    if path == "~" {
        return dirs::home_dir().unwrap_or_else(|| std::path::PathBuf::from("."));
    }
    if let Some(rest) = path.strip_prefix("~/") {
        return dirs::home_dir()
            .unwrap_or_else(|| std::path::PathBuf::from("."))
            .join(rest);
    }
    std::path::PathBuf::from(path)
}
