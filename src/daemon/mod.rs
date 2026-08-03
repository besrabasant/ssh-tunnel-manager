//! Daemon domain types and gRPC implementation.
//!
//! Tunnel process management lives in `tunnel`, while protobuf request/response
//! translation lives in `rpc_service`. Keeping those concerns separate makes the
//! lifecycle code testable without navigating all RPC handlers.

mod rpc_service;
mod tunnel;

use crate::store::{Entry, Store};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::sync::Arc;
use tokio::process::Child;
use tokio::sync::Mutex;

#[derive(Clone)]
/// Shared daemon state used by every concurrent tonic request handler.
pub struct Service {
    pub(super) store: Store,
    // The local port is the registry key because only one listener can own it.
    pub(super) tunnels: Arc<Mutex<HashMap<i32, RunningTunnel>>>,
}

/// Owns the child process so dropping or killing a registry entry closes its tunnel.
pub(super) struct RunningTunnel {
    pub(super) entry: Entry,
    pub(super) port: i32,
    pub(super) child: Child,
}

#[derive(Clone, Debug, Deserialize, Serialize)]
/// On-disk restoration record compatible with Go's active_tunnels.json schema.
pub(super) struct SavedTunnel {
    pub(super) config_name: String,
    pub(super) local_port: i32,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub(super) local_addr: String,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub(super) remote_addr: String,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub(super) server: String,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub(super) user: String,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub(super) saved_at: String,
}
