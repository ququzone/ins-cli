use clap::Parser;
use foundry_cli::opts::{EthereumOpts, TransactionOpts};

#[derive(Debug, Parser)]
pub struct RegisterArgs {
    #[clap(flatten)]
    tx: TransactionOpts,

    #[clap(flatten)]
    eth: EthereumOpts,

    sig: Option<String>,

    args: Vec<String>,
}

impl RegisterArgs {
    pub async fn run(self) -> eyre::Result<()> {
        Ok(())
    }
}
