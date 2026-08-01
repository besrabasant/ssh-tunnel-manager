fn main() -> Result<(), Box<dyn std::error::Error>> {
    println!("cargo:rerun-if-changed=rpc/daemon.proto");
    tonic_prost_build::configure().compile_protos(&["rpc/daemon.proto"], &["rpc"])?;
    Ok(())
}
