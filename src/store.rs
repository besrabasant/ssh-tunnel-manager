use crate::rpc::TunnelConfig;
use anyhow::{Context, Result, bail};
use serde::{Deserialize, Serialize};
use std::fs;
use std::path::{Path, PathBuf};

#[derive(Clone, Debug, Default, Deserialize, Eq, PartialEq, Serialize)]
#[serde(rename_all = "PascalCase")]
pub struct Entry {
    pub name: String,
    pub description: String,
    pub server: String,
    pub user: String,
    pub key_file: String,
    pub remote_host: String,
    pub remote_port: i32,
    pub local_port: i32,
}

impl Entry {
    pub fn validate(&self) -> Result<()> {
        let mut missing = Vec::new();
        for (name, value) in [
            ("Name", self.name.as_str()),
            ("Server", self.server.as_str()),
            ("User", self.user.as_str()),
            ("KeyFile", self.key_file.as_str()),
            ("RemoteHost", self.remote_host.as_str()),
        ] {
            if value.is_empty() {
                missing.push(format!("{name} is required."));
            }
        }
        if self.remote_port == 0 {
            missing.push("RemotePort is required.".into());
        }
        if missing.is_empty() {
            Ok(())
        } else {
            bail!("Entry is not valid. {}", missing.join(" "))
        }
    }
}

impl From<TunnelConfig> for Entry {
    fn from(value: TunnelConfig) -> Self {
        Self {
            name: value.name,
            description: value.description,
            server: value.server,
            user: value.user,
            key_file: value.key_file,
            remote_host: value.remote_host,
            remote_port: value.remote_port,
            local_port: value.local_port,
        }
    }
}

impl From<Entry> for TunnelConfig {
    fn from(value: Entry) -> Self {
        Self {
            name: value.name,
            description: value.description,
            server: value.server,
            user: value.user,
            key_file: value.key_file,
            remote_host: value.remote_host,
            remote_port: value.remote_port,
            local_port: value.local_port,
        }
    }
}

#[derive(Clone, Debug)]
pub struct Store {
    dir: PathBuf,
}

impl Store {
    pub fn new(dir: impl Into<PathBuf>) -> Result<Self> {
        let dir = dir.into();
        fs::create_dir_all(&dir)
            .with_context(|| format!("couldn't create directory {}", dir.display()))?;
        Ok(Self { dir })
    }

    fn path(&self, name: &str) -> Result<PathBuf> {
        if name.is_empty()
            || name.contains('/')
            || name.contains('\\')
            || name == "."
            || name == ".."
        {
            bail!("invalid configuration name {name:?}");
        }
        Ok(self.dir.join(format!("{name}.json")))
    }

    pub fn list(&self) -> Result<Vec<Entry>> {
        let mut paths = fs::read_dir(&self.dir)
            .with_context(|| format!("couldn't read directory {}", self.dir.display()))?
            .filter_map(Result::ok)
            .map(|item| item.path())
            .filter(|path| {
                path.is_file()
                    && path.extension().and_then(|v| v.to_str()) == Some("json")
                    && path.file_name().and_then(|v| v.to_str()) != Some("active_tunnels.json")
            })
            .collect::<Vec<_>>();
        paths.sort();
        paths
            .into_iter()
            .map(|path| self.read_path(&path))
            .collect()
    }

    fn read_path(&self, path: &Path) -> Result<Entry> {
        let data =
            fs::read(path).with_context(|| format!("couldn't read file {}", path.display()))?;
        serde_json::from_slice(&data)
            .with_context(|| format!("couldn't parse JSON file {}", path.display()))
    }

    pub fn get(&self, name: &str) -> Result<Entry> {
        self.read_path(&self.path(name)?)
    }

    pub fn add(&self, entry: &Entry) -> Result<()> {
        let path = self.path(&entry.name)?;
        if path.exists() {
            bail!("A configuration already exists with name {}", entry.name);
        }
        entry.validate()?;
        atomic_json(&path, entry)
    }

    pub fn update(&self, old_name: &str, entry: &Entry) -> Result<()> {
        let old_path = self.path(old_name)?;
        if !old_path.exists() {
            bail!("configuration {old_name:?} does not exist");
        }
        entry.validate()?;
        let new_path = self.path(&entry.name)?;
        atomic_json(&new_path, entry)?;
        if old_path != new_path {
            fs::remove_file(old_path)?;
        }
        Ok(())
    }

    pub fn remove(&self, name: &str) -> Result<()> {
        fs::remove_file(self.path(name)?)
            .with_context(|| format!("Cannot delete configuration {name}"))
    }

    pub fn active_path(&self) -> PathBuf {
        self.dir.join("active_tunnels.json")
    }
}

pub fn atomic_json(path: &Path, value: &impl Serialize) -> Result<()> {
    let data = serde_json::to_vec_pretty(value)?;
    let temp = path.with_extension(format!("json.tmp.{}", std::process::id()));
    fs::write(&temp, data).with_context(|| format!("couldn't write file {}", temp.display()))?;
    fs::rename(&temp, path).with_context(|| format!("couldn't replace file {}", path.display()))?;
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    fn entry() -> Entry {
        Entry {
            name: "dev".into(),
            server: "host".into(),
            user: "me".into(),
            key_file: "/key".into(),
            remote_host: "db".into(),
            remote_port: 5432,
            local_port: 15432,
            description: "Development".into(),
        }
    }

    #[test]
    fn reads_and_writes_go_compatible_json() {
        let temp = tempfile::tempdir().unwrap();
        let store = Store::new(temp.path()).unwrap();
        store.add(&entry()).unwrap();
        let text = fs::read_to_string(temp.path().join("dev.json")).unwrap();
        assert!(text.contains("\"KeyFile\""));
        assert_eq!(store.get("dev").unwrap(), entry());
    }

    #[test]
    fn rejects_path_traversal_names() {
        let temp = tempfile::tempdir().unwrap();
        let store = Store::new(temp.path()).unwrap();
        assert!(store.get("../secret").is_err());
    }
}
