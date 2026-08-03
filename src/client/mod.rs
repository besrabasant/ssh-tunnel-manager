//! Command parsing and dispatch for the `sshtm` client.
//!
//! Connection setup, prompting, presentation, and the interactive menu are kept
//! in focused submodules so command dispatch only coordinates those pieces.

mod connection;
mod interactive;
mod output;
mod prompt;

use crate::client::connection::connect;
use crate::client::output::{print_config, print_mutation, print_start, print_tunnel};
use crate::client::prompt::prompt_config;
use crate::rpc::*;
use anyhow::{Result, bail};
use clap::{CommandFactory, Parser, Subcommand};
use clap_complete::{Shell, generate};
use std::io;

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

/// Parses command-line arguments and executes the selected client operation.
pub async fn run() -> Result<()> {
    let cli = Cli::parse();
    match cli.command {
        Some(Commands::Completion { shell }) => {
            // Generate completions from the same Clap command tree used for parsing.
            generate(shell, &mut Cli::command(), "sshtm", &mut io::stdout());
        }
        Some(Commands::Version) => {
            println!("SSH Tunnel Manager Version {}", env!("CARGO_PKG_VERSION"));
        }
        Some(Commands::List { search_pattern }) => {
            let response = connect()
                .await?
                .list_configurations(ListConfigurationsRequest {
                    search_pattern: search_pattern.unwrap_or_default(),
                })
                .await?
                .into_inner();
            if response.configs.is_empty() {
                print!("{}", response.result);
            } else {
                response.configs.iter().for_each(print_config);
            }
        }
        Some(Commands::Active) => {
            let response = connect()
                .await?
                .list_active_tunnels(ListActiveTunnelsRequest {})
                .await?
                .into_inner();
            if response.tunnels.is_empty() {
                print!("{}", response.result);
            } else {
                response.tunnels.iter().for_each(print_tunnel);
            }
        }
        Some(Commands::Add) => {
            let data = prompt_config(None)?;
            let response = connect()
                .await?
                .add_configuration(AddOrUpdateConfigurationRequest {
                    name: data.name.clone(),
                    data: Some(data),
                })
                .await?
                .into_inner();
            print_mutation(response.status, &response.message);
        }
        Some(Commands::Edit { name }) => {
            let mut daemon = connect().await?;
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
            let response = connect()
                .await?
                .delete_configuration(DeleteConfigurationRequest { name })
                .await?
                .into_inner();
            print_mutation(response.status, &response.message);
        }
        Some(Commands::Tunnel { name, local_port }) => {
            let response = connect()
                .await?
                .start_tunnel(StartTunnelRequest {
                    config_name: name,
                    local_port: local_port.unwrap_or(-1),
                })
                .await?
                .into_inner();
            print_start(&response);
        }
        Some(Commands::Kill { identifier }) => {
            // Match the legacy CLI: an integer identifies a port; all other text is a name.
            let (config_name, local_port) = match identifier.parse::<i32>() {
                Ok(port) => (String::new(), port),
                Err(_) => (identifier, 0),
            };
            let response = connect()
                .await?
                .kill_tunnel(KillTunnelRequest {
                    config_name,
                    local_port,
                })
                .await?
                .into_inner();
            print_mutation(response.status, &response.message);
        }
        None => interactive::run().await?,
    }
    Ok(())
}
