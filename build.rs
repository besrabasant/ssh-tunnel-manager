fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Rebuild generated Rust bindings only when the protocol contract changes.
    println!("cargo:rerun-if-changed=rpc/daemon.proto");
    // Generate both prost messages and tonic client/server implementations.
    tonic_prost_build::configure().compile_protos(&["rpc/daemon.proto"], &["rpc"])?;
    Ok(())
}
