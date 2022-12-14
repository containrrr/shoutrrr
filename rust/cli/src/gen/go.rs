use shoutrrr::spec::ServiceSpec;
use std::{fs};
use eyre::WrapErr;
// use genco::{fmt, prelude::*};


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

    pub fn generate(&mut self, spec: &ServiceSpec) -> eyre::Result<()> {

        let out_file = self.args.out_dir.join(format!("{}_config.gen.go", spec.scheme));

        if !self.args.out_dir.exists() {
            fs::create_dir_all(self.args.out_dir.as_path()).wrap_err("Failed to create output directory")?;
        }

        println!("Using output file {out_file:?}");

        // let shoutrrr = &go::import("@containrrr/shoutrrr", "shoutrrr");
        // let serviceConfig = &go::import("containrrr/shoutrrr", "ServiceConfig");
        // let name = capitalized(&spec.scheme);

        // let tokens = quote! {
            
        //     export default class $(&name)Config extends $serviceConfig {
                
        //         const name = $(quoted(&name))
                
        //     }
        // };

        // for (i, line) in tokens.to_file_vec()?.iter().enumerate() {
        //     println!("{i} {line}");
        // }

        

        Ok(())
    }
}