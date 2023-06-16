use clap::Parser;

use ins_cli::cmd::commands::{Cmd, Subcommands};

#[tokio::main]
async fn main() -> eyre::Result<()> {
    let cmd = Cmd::parse();

    match cmd.sub {
        Subcommands::Register(cmd) => cmd.run().await?
    }

    Ok(())
}
