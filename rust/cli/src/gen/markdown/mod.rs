
use shoutrrr::spec::{self, ServiceSpec, URLPart};
use std::{fs, collections::{HashMap, HashSet}};
use eyre::WrapErr;
use std::fmt::Write;

mod fields;
mod url;

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

        let out_file = self.args.out_dir.join(format!("{}_config.gen.md", spec.scheme));

        if !self.args.out_dir.exists() {
            fs::create_dir_all(self.args.out_dir.as_path()).wrap_err("Failed to create output directory")?;
        }

        eprintln!("Using output file {out_file:?}");

        let mut query_fields = Vec::new();
        let mut url_fields = HashMap::new();

        for (_, prop) in &spec.props {
            if prop.url_parts.is_empty() || prop.url_parts.contains(&URLPart::Query) {
                query_fields.push(prop);
            } else {
                for url_part in prop.url_parts.iter().copied() {
                    url_fields.insert(url_part, prop);
                }
            }
        }

        query_fields.sort_by(|a, b| a.name.cmp(&b.name));

        let url_tokens = self.gen_url_tokens(&url_fields, &spec.scheme)?;

        println!("{url_tokens}");

        let query_tokens = self.gen_query_tokens(&query_fields)?;

        println!("{query_tokens}");

        Ok(())
    }

 

    pub(crate) fn gen_url_tokens(&self, props: &HashMap<spec::URLPart, &spec::SpecProp>, scheme: &str) -> eyre::Result<String> {
        
        let mut output = "### URL Fields\n\n".to_string();
        let mut fields_printed = HashSet::new();

        let slug_map = url::get_slugs(&scheme, &props);

        dbg_pls::color!(&slug_map);

        let props: Vec<_> = URLPart::iter().iter().filter_map(|up| props.get(up)).collect();
        
        // Maybe it's better to have the fields sorted alphabetically? 
        // Leaving them in URL order for now
        //props.sort_by(|a, b| a.name.cmp(&b.name));

        for prop in props {

            if fields_printed.contains(&prop.name) {
                continue;
            }

            fields::write_primary(&mut output, prop)?;

            output.push_str("  URL part: <code class=\"service-url\">");

            url::write_url_map(&mut output, &prop.url_parts, &slug_map)?;

            output.push_str("</code>  \n");
            fields_printed.insert(&prop.name);

        }

        Ok(output)
    }

    pub(crate) fn gen_query_tokens(&self, query_fields: &Vec<&spec::SpecProp>) -> eyre::Result<String> {
        let mut output = "### Query/Param Props\n".to_string();

        writeln!(output)?;
        writeln!(output, "Props can be either supplied using the params argument, or through the URL using  \n`?key=value&key=value` etc.")?;
        writeln!(output)?;

        for field in query_fields {
            fields::write_primary(&mut output, field)?;
            fields::write_extras(&mut output, field)?;
            writeln!(output)?;
        }

        Ok(output)
    }

}