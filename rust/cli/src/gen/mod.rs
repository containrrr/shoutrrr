use std::{io::Cursor, path::PathBuf};

use clap::{arg, ValueEnum};
use eyre::bail;
use shoutrrr::spec::ServiceSpec;
use crate::clapex;

mod go;
mod markdown;
mod ecmascript;

#[derive(clap::Parser, Debug, Clone)]
#[command(
    about = "Display a parsed service spec"
)]
pub struct Args {
    /// Target to generate for
    #[arg()]
    target: GenTarget,

    /// Spec file to read
    #[arg(id = "SPEC-PATH", value_parser = clapex::ExistingFileValueParser{})]
    spec_path: PathBuf,

    /// Output root directory
    #[arg(id = "OUT-DIR", short = 'o')]
    out_dir: PathBuf,
}

#[derive(ValueEnum, Clone, Debug)]
#[value(rename_all = "lowercase")]
pub enum GenTarget {
    #[value(alias("javascript"), alias("js"))]
    ECMAScript,
    #[value(alias = "golang")]
    Go,
    Rust,
    #[value(alias = "md")]
    Markdown,
}



impl Args {
    pub fn run(&self) -> eyre::Result<()> {

        let buf = std::fs::read(self.spec_path.as_path())?;
        let mut spec: ServiceSpec = serde_yaml::from_reader(Cursor::new(buf))?;

        spec.init_props()?;

        match &self.target {
            GenTarget::ECMAScript => ecmascript::generate(self.clone(), &spec),
            GenTarget::Markdown => markdown::generate(self.clone(), &spec),
            GenTarget::Go => go::generate(self.clone(), &spec),
            target => bail!("target {target:?} is not implemented")
        }?;


        Ok(())
    }
}

