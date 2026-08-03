//! Human-readable rendering for protobuf responses.

use crate::rpc::{ActiveTunnel, ResponseStatus, StartTunnelResponse, TunnelConfig};

pub(super) fn print_mutation(status: i32, message: &str) {
    // Keep failures on stderr so command output remains usable in shell pipelines.
    if status == ResponseStatus::Error as i32 {
        eprintln!("{message}");
    } else {
        println!("{message}");
    }
}

pub(super) fn print_start(response: &StartTunnelResponse) {
    // Structured events preserve daemon setup progress without parsing legacy text.
    response.events.iter().for_each(|event| println!("{event}"));
    if !response.message.is_empty() {
        print_mutation(response.status, &response.message);
    }
    if response.events.is_empty() && response.message.is_empty() {
        print!("{}", response.result);
    }
}

pub(super) fn print_config(config: &TunnelConfig) {
    let port = if config.local_port == 0 {
        "auto".into()
    } else {
        config.local_port.to_string()
    };
    let local = if config.local_port == 0 {
        "auto".into()
    } else {
        format!("127.0.0.1:{}", config.local_port)
    };
    let description = if config.description.trim().is_empty() {
        String::new()
    } else {
        format!(" ({})", config.description)
    };
    println!(
        ":{port}\n- Connection:                  {}{description}\n- Remote Address:              {}:{}\n- Local Address:               {local}\n- SSH server:                  {}\n- User:                        {}\n- Private key:                 {}",
        config.name,
        config.remote_host,
        config.remote_port,
        config.server,
        config.user,
        config.key_file
    );
}

pub(super) fn print_tunnel(tunnel: &ActiveTunnel) {
    println!(
        ":{}\n- Connection:                  {}\n- Remote Address:              {}\n- Local Address:               {}\n- SSH server:                  {}",
        tunnel.local_port, tunnel.name, tunnel.remote_addr, tunnel.local_addr, tunnel.server
    );
}
