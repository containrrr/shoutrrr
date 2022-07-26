use eyre::{eyre};
use color_eyre::{self, Result};
use std::{env, io::Cursor, collections::HashMap};
use dbg_pls::{color, DebugPls};

fn main() -> Result<()> {
    color_eyre::install()?;

    let mut args = env::args().into_iter();

    let _ = args.next();
    let service = args
        .next()
        .ok_or_else(|| eyre!("SERVICE argument missing"))?;

    let buf = std::fs::read(format!("../spec/{service}.yml"))?;
    let mut spec: ServiceSpec = serde_yaml::from_reader(Cursor::new(buf))?;

    spec.init_props()?;

    println!("{}", color(&spec));

    Ok(())
}

#[derive(serde::Deserialize, Debug, DebugPls)]
pub struct ServiceSpec {
    pub version: u32,
    pub scheme: String,
    #[serde(default)]
    pub options: SpecOptions,
    pub props: HashMap<String, SpecProp>,
}

impl ServiceSpec {
    fn init_props(&mut self) -> Result<()> {
        for (key, prop) in &mut self.props {
            prop.name = key.clone();
        }
        Ok(())
    }
}

#[derive(serde::Deserialize, Debug, DebugPls, Default)]
#[serde(rename_all = "camelCase")]
pub struct SpecOptions {
    pub reverse_path_prio: bool,
    pub custom_query_vars: bool,
}

#[derive(serde::Deserialize, Debug, DebugPls)]
#[serde(rename_all = "camelCase")]
pub struct SpecProp {
    #[serde(rename = "type")]
    pub prop_type: PropType,
    pub description: String,
    #[serde(skip)]
    pub name: String,
    pub default_value: Option<String>,
    pub template: Option<String>,
    #[serde(default)]
    pub required: bool,
    #[serde(default)]
    pub url_parts: Vec<String>,
    // pub title: bool,
    pub base: Option<i32>,
    #[serde(default)]
    pub keys: Vec<String>,
    #[serde(default)]
    pub values: Vec<String>,
    pub custom_type: Option<String>,
    #[serde(default)]
    pub credential: bool,
    pub item_separator: Option<String>,
    pub value_separator: Option<String>,
}

#[derive(serde::Deserialize, Debug, DebugPls)]
#[serde(rename_all = "camelCase")]
pub enum PropType {
    Text,
    Number,
    List,
    Toggle,
    Option,
}
