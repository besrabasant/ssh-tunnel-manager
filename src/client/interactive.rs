//! Interactive menu shown when `sshtm` is invoked without a subcommand.

use super::connection::connect;
use super::output::{print_config, print_mutation, print_start, print_tunnel};
use super::prompt::prompt;
use crate::rpc::*;
use anyhow::Result;

pub(super) async fn run() -> Result<()> {
    // Keep a dependency-light line-oriented UI for users who invoke sshtm without arguments.
    loop {
        println!("\nSSH Tunnel Manager");
        println!("1) List configurations  2) Active tunnels  3) Start  4) Kill  5) Quit");
        match prompt("Select", "")?.as_str() {
            "1" => {
                let response = connect()
                    .await?
                    .list_configurations(ListConfigurationsRequest::default())
                    .await?
                    .into_inner();
                if response.configs.is_empty() {
                    print!("{}", response.result);
                }
                response.configs.iter().for_each(print_config);
            }
            "2" => {
                let response = connect()
                    .await?
                    .list_active_tunnels(ListActiveTunnelsRequest {})
                    .await?
                    .into_inner();
                if response.tunnels.is_empty() {
                    print!("{}", response.result);
                }
                response.tunnels.iter().for_each(print_tunnel);
            }
            "3" => {
                let response = connect()
                    .await?
                    .start_tunnel(StartTunnelRequest {
                        config_name: prompt("Configuration name", "")?,
                        local_port: -1,
                    })
                    .await?
                    .into_inner();
                print_start(&response);
            }
            "4" => {
                let identifier = prompt("Configuration name or local port", "")?;
                let (config_name, local_port) = identifier
                    .parse::<i32>()
                    .map(|port| (String::new(), port))
                    .unwrap_or((identifier, 0));
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
            "5" | "q" | "quit" => return Ok(()),
            _ => eprintln!("Unknown selection"),
        }
    }
}
