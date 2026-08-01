use anyhow::{Context, Result, bail};
use clap::{CommandFactory, Parser, Subcommand};
use clap_complete::{Shell, generate};
use sshtm::rpc::daemon_service_client::DaemonServiceClient;
use sshtm::rpc::*;
use std::io::{self, Write};
use tonic::transport::{Channel, Endpoint};

#[derive(Parser)]
#[command(name = "sshtm", version, about = "Manage SSH tunnels with ease")]
struct Cli {
    // No subcommand intentionally launches the interactive menu.
    #[command(subcommand)]
    command: Option<Commands>,
}

#[derive(Subcommand)]
enum Commands {
    #[command(alias = "l", alias = "ls")]
    List {
        search_pattern: Option<String>,
    },
    #[command(alias = "a")]
    Add,
    #[command(alias = "e")]
    Edit {
        name: String,
    },
    #[command(alias = "d", alias = "del")]
    Delete {
        name: String,
    },
    #[command(alias = "t")]
    Tunnel {
        name: String,
        local_port: Option<i32>,
    },
    Active,
    #[command(alias = "k", alias = "terminate")]
    Kill {
        identifier: String,
    },
    Completion {
        shell: Shell,
    },
    /// Print the version number of SSH Tunnel Manager.
    Version,
}

async fn client() -> Result<DaemonServiceClient<Channel>> {
    // Bound connection and request timeouts keep scripts from hanging on a dead daemon.
    let endpoint = Endpoint::from_static("http://127.0.0.1:50051")
        .connect_timeout(std::time::Duration::from_secs(10))
        .timeout(std::time::Duration::from_secs(20));
    Ok(DaemonServiceClient::connect(endpoint)
        .await
        .context("did not connect to sshtmd")?)
}

#[tokio::main]
async fn main() {
    // Render the complete anyhow context chain and expose failures through the exit code.
    if let Err(error) = run().await {
        eprintln!("{error:#}");
        std::process::exit(1);
    }
}

async fn run() -> Result<()> {
    let cli = Cli::parse();
    match cli.command {
        Some(Commands::Completion { shell }) => {
            // Generate completions from the same Clap command tree used for parsing.
            generate(shell, &mut Cli::command(), "sshtm", &mut io::stdout());
        }
        Some(Commands::Version) => {
            println!("SSH Tunnel Manager Version {}", env!("CARGO_PKG_VERSION"))
        }
        Some(Commands::List { search_pattern }) => {
            let response = client()
                .await?
                .list_configurations(ListConfigurationsRequest {
                    search_pattern: search_pattern.unwrap_or_default(),
                })
                .await?
                .into_inner();
            if response.configs.is_empty() {
                print!("{}", response.result);
            } else {
                for config in &response.configs {
                    print_config(config);
                }
            }
        }
        Some(Commands::Active) => {
            let response = client()
                .await?
                .list_active_tunnels(ListActiveTunnelsRequest {})
                .await?
                .into_inner();
            if response.tunnels.is_empty() {
                print!("{}", response.result);
            } else {
                for tunnel in &response.tunnels {
                    print_active(tunnel);
                }
            }
        }
        Some(Commands::Add) => {
            let data = prompt_config(None)?;
            let name = data.name.clone();
            let response = client()
                .await?
                .add_configuration(AddOrUpdateConfigurationRequest {
                    name,
                    data: Some(data),
                })
                .await?
                .into_inner();
            print_mutation(response.status, &response.message);
        }
        Some(Commands::Edit { name }) => {
            let mut daemon = client().await?;
            // Fetch first so pressing Enter in a prompt retains the existing field value.
            let fetched = daemon
                .fetch_configuration(FetchConfigurationRequest { name: name.clone() })
                .await?
                .into_inner();
            if fetched.status == ResponseStatus::Error as i32 {
                bail!("{}", fetched.message);
            }
            let data = prompt_config(fetched.data)?;
            let response = daemon
                .update_configuration(AddOrUpdateConfigurationRequest {
                    name,
                    data: Some(data),
                })
                .await?
                .into_inner();
            print_mutation(response.status, &response.message);
        }
        Some(Commands::Delete { name }) => {
            let response = client()
                .await?
                .delete_configuration(DeleteConfigurationRequest { name })
                .await?
                .into_inner();
            print_mutation(response.status, &response.message);
        }
        Some(Commands::Tunnel { name, local_port }) => {
            let response = client()
                .await?
                .start_tunnel(StartTunnelRequest {
                    config_name: name,
                    local_port: local_port.unwrap_or(-1),
                })
                .await?
                .into_inner();
            for event in &response.events {
                // Structured events preserve daemon setup progress without parsing result text.
                println!("{event}");
            }
            if !response.message.is_empty() {
                print_mutation(response.status, &response.message);
            }
            if response.events.is_empty() && response.message.is_empty() {
                print!("{}", response.result);
            }
        }
        Some(Commands::Kill { identifier }) => {
            // Match the legacy CLI: an integer identifies a port; all other text is a name.
            let (config_name, local_port) = match identifier.parse::<i32>() {
                Ok(port) => (String::new(), port),
                Err(_) => (identifier, 0),
            };
            let response = client()
                .await?
                .kill_tunnel(KillTunnelRequest {
                    config_name,
                    local_port,
                })
                .await?
                .into_inner();
            print_mutation(response.status, &response.message);
        }
        None => interactive().await?,
    }
    Ok(())
}

async fn interactive() -> Result<()> {
    // Keep a dependency-light line-oriented UI for users who invoke sshtm without arguments.
    loop {
        println!("\nSSH Tunnel Manager");
        println!("1) List configurations  2) Active tunnels  3) Start  4) Kill  5) Quit");
        match prompt("Select", "")?.as_str() {
            "1" => {
                let response = client()
                    .await?
                    .list_configurations(ListConfigurationsRequest::default())
                    .await?
                    .into_inner();
                if response.configs.is_empty() {
                    print!("{}", response.result);
                }
                for config in &response.configs {
                    print_config(config);
                }
            }
            "2" => {
                let response = client()
                    .await?
                    .list_active_tunnels(ListActiveTunnelsRequest {})
                    .await?
                    .into_inner();
                if response.tunnels.is_empty() {
                    print!("{}", response.result);
                }
                for tunnel in &response.tunnels {
                    print_active(tunnel);
                }
            }
            "3" => {
                let name = prompt("Configuration name", "")?;
                let response = client()
                    .await?
                    .start_tunnel(StartTunnelRequest {
                        config_name: name,
                        local_port: -1,
                    })
                    .await?
                    .into_inner();
                for event in &response.events {
                    println!("{event}");
                }
                if !response.message.is_empty() {
                    print_mutation(response.status, &response.message);
                }
            }
            "4" => {
                let identifier = prompt("Configuration name or local port", "")?;
                let (config_name, local_port) = identifier
                    .parse::<i32>()
                    .map(|p| (String::new(), p))
                    .unwrap_or((identifier, 0));
                let response = client()
                    .await?
                    .kill_tunnel(KillTunnelRequest {
                        config_name,
                        local_port,
                    })
                    .await?
                    .into_inner();
                print_mutation(response.status, &response.message);
            }
            "5" | "q" | "quit" => return Ok(()),
            _ => eprintln!("Unknown selection"),
        }
    }
}

fn prompt_config(existing: Option<TunnelConfig>) -> Result<TunnelConfig> {
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

fn prompt(label: &str, default: &str) -> Result<String> {
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

fn print_mutation(status: i32, message: &str) {
    // Keep failures on stderr so command output remains usable in shell pipelines.
    if status == ResponseStatus::Error as i32 {
        eprintln!("{message}");
    } else {
        println!("{message}");
    }
}

fn print_config(config: &TunnelConfig) {
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

fn print_active(tunnel: &ActiveTunnel) {
    println!(
        ":{}\n- Connection:                  {}\n- Remote Address:              {}\n- Local Address:               {}\n- SSH server:                  {}",
        tunnel.local_port, tunnel.name, tunnel.remote_addr, tunnel.local_addr, tunnel.server
    );
}
