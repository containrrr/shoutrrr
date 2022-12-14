use shoutrrr::spec::ServiceSpec;
use std::{fs};
use eyre::WrapErr;
use genco::{prelude::*};


pub struct Generator {
    args: super::Args
}

pub fn generate(args: super::Args, spec: &ServiceSpec) -> eyre::Result<()> {
    Generator::new(args).generate(spec)
}

impl Generator {
    pub fn new(args: super::Args) -> Self {
        Self { args }
    }

    pub fn generate(&self, spec: &ServiceSpec) -> eyre::Result<()> {

        let out_file = self.args.out_dir.join(format!("{}_config.gen.js", spec.scheme));

        if !self.args.out_dir.exists() {
            fs::create_dir_all(self.args.out_dir.as_path()).wrap_err("Failed to create output directory")?;
        }

        println!("Using output file {out_file:?}");

        // let shoutrrr = &js::import("@containrrr/shoutrrr", "shoutrrr").into_default();
        let service_config = &js::import("@containrrr/shoutrrr", "ServiceConfig");
        let name = capitalized(&spec.scheme);

        let tokens = quote! {
            export default class $(&name)Config extends $service_config {
                
                const name = $(quoted(&name))
                
            }
        };

        for (i, line) in tokens.to_file_vec()?.iter().enumerate() {
            println!("{i} {line}");
        }

        

        Ok(())
    }
}

fn capitalized(s: &str) -> String {
    s[0..1].to_uppercase() + &s[1..]
}