use anyhow::Result;
use clap::Parser;
use rgrep::*;

// cargo run --quiet -- "Re[^\\s]+" "src/*.rs"
fn main() -> Result<()> {
    let config: GrepConfig = GrepConfig::parse();
    config.match_with_default_strategy()?;

    Ok(())
}