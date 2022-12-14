use clap::{Parser, Subcommand};
use color_eyre::{self, Result};


use shoutrrr_cli::{parse, gen};

#[derive(Parser)]
pub struct Args {
    #[command(subcommand)]
    command: Command
}

#[derive(Subcommand, Clone, Debug)]
enum Command {
    #[command(alias = "gen")]
    Generate(gen::Args),
    Parse(parse::Args),
}

fn main() -> Result<()> {
    color_eyre::install()?;

    let cmd = Args::parse();
    match cmd.command {
        Command::Parse(c) => c.run(),
        Command::Generate(c) => c.run(),
    }
}
