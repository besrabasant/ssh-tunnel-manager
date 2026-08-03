//! Line-oriented configuration and menu prompts.

use crate::rpc::TunnelConfig;
use anyhow::{Context, Result, bail};
use std::io::{self, Write};

pub(super) fn prompt_config(existing: Option<TunnelConfig>) -> Result<TunnelConfig> {
    // The same prompt sequence handles creation and editing through optional defaults.
    let old = existing.unwrap_or_default();
    let name = prompt("Name", &old.name)?;
    let description = prompt("Description", &old.description)?;
    let server = prompt("SSH server", &old.server)?;
    let user = prompt("SSH user", &old.user)?;
    let key_file = prompt("Private key", &old.key_file)?;
    let remote_host = prompt("Remote host", &old.remote_host)?;
    let remote_port = parse_port("Remote port", old.remote_port, false)?;
    let local_port = parse_port("Local port (0 for automatic)", old.local_port, true)?;
    Ok(TunnelConfig {
        name,
        description,
        server,
        user,
        key_file,
        remote_host,
        remote_port,
        local_port,
    })
}

pub(super) fn prompt(label: &str, default: &str) -> Result<String> {
    if default.is_empty() {
        print!("{label}: ");
    } else {
        print!("{label} [{default}]: ");
    }
    io::stdout().flush()?;
    let mut value = String::new();
    io::stdin().read_line(&mut value)?;
    // Empty input accepts the displayed default, which is empty for new profiles.
    let value = value.trim().to_string();
    Ok(if value.is_empty() {
        default.to_string()
    } else {
        value
    })
}

fn parse_port(label: &str, default: i32, allow_zero: bool) -> Result<i32> {
    let value = prompt(label, &default.to_string())?
        .parse::<i32>()
        .with_context(|| format!("{label} must be a number"))?;
    if value < 0 || value > 65535 || (!allow_zero && value == 0) {
        bail!(
            "{label} must be between {} and 65535",
            if allow_zero { 0 } else { 1 }
        );
    }
    Ok(value)
}
