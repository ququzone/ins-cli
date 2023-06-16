use clap::{Parser, Subcommand};

use super::RegisterArgs;

#[derive(Debug, Parser)]
#[clap(name = "ins-cli", version = "0.1.0")]
pub struct Cmd {
    #[clap(subcommand)]
    pub sub: Subcommands,
}

#[derive(Debug, Subcommand)]
#[clap(
    about = "INS operations from your command line.",
    after_help = "Find more information can refer code: https://github.com/ququzone/ins-cli",
    next_display_order = None
)]
pub enum Subcommands {
    #[clap(name = "register")]
    #[clap(visible_aliases = &["reg"])]
    #[clap(about = "Register INS name")]
    Register(RegisterArgs),
}
