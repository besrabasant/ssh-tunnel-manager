pub mod daemon;
pub mod store;

pub mod rpc {
    tonic::include_proto!("daemon");
}

pub const ADDRESS: &str = "127.0.0.1:50051";
pub const DEFAULT_CONFIG_DIR: &str = ".ssh-tunnel-manager";

pub fn config_dir() -> std::path::PathBuf {
    if let Ok(path) = std::env::var("config-dir") {
        return expand_home(&path);
    }
    dirs::home_dir()
        .unwrap_or_else(|| std::path::PathBuf::from("."))
        .join(DEFAULT_CONFIG_DIR)
}

pub fn expand_home(path: &str) -> std::path::PathBuf {
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
