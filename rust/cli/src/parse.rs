use std::{io::Cursor, path::PathBuf, fs};

use clap::{arg};
use dbg_pls::color;
use shoutrrr::spec::ServiceSpec;

#[derive(clap::Parser, Debug, Clone)]
#[command(
    about = "Display a parsed service spec"
)]
pub struct Args {
    /// Spec file to read
    #[arg(id = "spec")]
    spec_path: PathBuf,
}

impl Args {
    pub fn run(&self) -> eyre::Result<()> {
        let buf = fs::read(self.spec_path.as_path())?;
        let mut spec: ServiceSpec = serde_yaml::from_reader(Cursor::new(buf))?;

        spec.init_props()?;

        println!("{}", color(&spec));
        Ok(())
    }
}