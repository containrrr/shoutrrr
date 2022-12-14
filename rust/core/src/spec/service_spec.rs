use eyre::Result;
use std::{collections::HashMap};


#[cfg_attr(feature = "dbg-pls", derive(dbg_pls::DebugPls))]
#[cfg_attr(feature = "serde", 
    derive(serde::Deserialize, serde::Serialize), 
    serde(rename_all = "camelCase", ))]
#[derive(Debug)]
pub struct ServiceSpec {
    pub version: u32,
    pub scheme: String,
    #[cfg_attr(feature = "serde", serde(default))]
    pub options: SpecOptions,
    pub props: HashMap<String, SpecProp>,
}

impl ServiceSpec {
    pub fn init_props(&mut self) -> Result<()> {
        for (key, prop) in &mut self.props {
            prop.name = key.clone();
        }
        Ok(())
    }
}

#[cfg_attr(feature = "dbg-pls", derive(dbg_pls::DebugPls))]
#[cfg_attr(feature = "serde", derive(serde::Deserialize, serde::Serialize), serde(rename_all = "camelCase"))]
#[derive(Debug, Default)]
pub struct SpecOptions {
    #[cfg_attr(feature = "serde", serde(default, 
        skip_serializing_if = "serde_util::is_false", 
        deserialize_with = "serde_util::deserialize_bool"
    ))]
    pub reverse_path_prio: bool,
    
    #[cfg_attr(feature = "serde", serde(default, 
        skip_serializing_if = "serde_util::is_false", 
        deserialize_with = "serde_util::deserialize_bool"
    ))]
    pub custom_query_vars: bool,
}

#[cfg_attr(feature = "dbg-pls", derive(dbg_pls::DebugPls))]
#[cfg_attr(feature = "serde", derive(serde::Deserialize, serde::Serialize), serde(rename_all = "camelCase"))]
#[derive(Debug, Hash, PartialEq, Eq)]
pub struct SpecProp {
    #[cfg_attr(feature = "serde", serde(flatten))]
    pub prop_type: PropType,
    pub description: String,
    #[cfg_attr(feature = "serde", serde(skip))]
    pub name: String,
    #[cfg_attr(feature = "serde", serde(skip_serializing_if = "Option::is_none", rename = "default"))]
    pub default_value: Option<String>,
    #[cfg_attr(feature = "serde", serde(skip_serializing_if = "Option::is_none"))]
    pub template: Option<String>,
    #[cfg_attr(feature = "serde", serde(default, deserialize_with = "serde_util::deserialize_bool"))]
    pub required: bool,
    #[cfg_attr(feature = "serde", serde(default, alias="urlparts"))]
    pub url_parts: Vec<URLPart>,
    // pub title: bool,
    #[cfg_attr(feature = "serde", serde(default))]
    pub keys: Vec<String>,
    #[cfg_attr(feature = "serde", serde(default, deserialize_with = "serde_util::deserialize_bool"))]
    pub credential: bool,
    #[cfg_attr(feature = "serde", serde(skip_serializing_if = "Option::is_none"))]
    pub value_separator: Option<String>,
}

#[cfg_attr(feature = "dbg-pls", derive(dbg_pls::DebugPls))]
#[cfg_attr(feature = "serde", derive(serde::Deserialize, serde::Serialize), serde(rename_all = "camelCase"))]
#[derive(Debug, Hash, PartialEq, Eq)]
pub struct NumberSpecProp {
    #[cfg_attr(feature = "serde", serde(skip_serializing_if = "Option::is_none"))]
    pub base: Option<i32>,
}

#[cfg_attr(feature = "dbg-pls", derive(dbg_pls::DebugPls))]
#[cfg_attr(feature = "serde", derive(serde::Deserialize, serde::Serialize), serde(rename_all = "camelCase"))]
#[derive(Debug, Hash, PartialEq, Eq)]
pub struct OptionSpecProp {
    #[cfg_attr(feature = "serde", serde(default))]
    pub values: Vec<String>,
}

#[cfg_attr(feature = "dbg-pls", derive(dbg_pls::DebugPls))]
#[cfg_attr(feature = "serde", derive(serde::Deserialize, serde::Serialize), serde(rename_all = "camelCase"))]
#[derive(Debug, Hash, PartialEq, Eq)]
pub struct CustomSpecProp {
    #[cfg_attr(feature = "serde",serde(skip_serializing_if = "Option::is_none"))]
    pub custom_type: Option<String>,
}

#[cfg_attr(feature = "dbg-pls", derive(dbg_pls::DebugPls))]
#[cfg_attr(feature = "serde", derive(serde::Deserialize, serde::Serialize), serde(rename_all = "camelCase"))]
#[derive(Debug, Hash, PartialEq, Eq)]
pub struct MapSpecProp {
    #[cfg_attr(feature = "serde", serde(skip_serializing_if = "Option::is_none"))]
    pub item_separator: Option<String>,
}



#[cfg_attr(feature = "dbg-pls", derive(dbg_pls::DebugPls))]
#[cfg_attr(feature = "serde", derive(serde::Deserialize, serde::Serialize), serde(rename_all = "camelCase", tag = "type"))]
#[derive(Debug, Hash, PartialEq, Eq)]
pub enum PropType {
    Text,
    Number(NumberSpecProp),
    List,
    Toggle,
    Option(OptionSpecProp),
    Color,
    Custom(CustomSpecProp),
    Map(MapSpecProp),
}

#[cfg_attr(feature = "dbg-pls", derive(dbg_pls::DebugPls))]
#[cfg_attr(feature = "serde", derive(serde::Deserialize, serde::Serialize), serde(rename_all = "camelCase"))]
#[derive(Debug, PartialEq, Eq, Hash, Clone, Copy)]
pub enum URLPart {
    Scheme,   
    Query,    
	User,     
	Password, 
	Host,     
	Port,     
	Path1,    
	Path2,    
	Path3,    
	Path4,    
	Path,
}

impl URLPart {
    pub fn iter() -> [Self; 11] {
        [
            Self::Scheme,   
            Self::Query,    
            Self::User,     
            Self::Password, 
            Self::Host,     
            Self::Port,     
            Self::Path1,    
            Self::Path2,    
            Self::Path3,    
            Self::Path4,    
            Self::Path,
        ]
    }
}


#[cfg(feature = "serde")]
mod serde_util {
    use serde::de::{Error, Unexpected};

    pub fn deserialize_bool<'de, D>(deserializer: D) -> Result<bool, D::Error> where D: serde::Deserializer<'de> {
        deserializer.deserialize_str(BoolVisitor{})
    }
    
    pub struct BoolVisitor;
    
    impl<'de> serde::de::Visitor<'de> for BoolVisitor {
        type Value = bool;
    
        fn expecting(&self, formatter: &mut std::fmt::Formatter) -> std::fmt::Result {
            formatter.write_str("true, false, yes or no")
        }
    
        fn visit_bool<E>(self, v: bool) -> Result<Self::Value, E>
            where
                E: serde::de::Error, {
            Ok(v)
        }
    
        fn visit_str<E>(self, v: &str) -> Result<Self::Value, E>
            where
                E: serde::de::Error, {
            parse_bool(v).ok_or_else(|| Error::invalid_value(Unexpected::Str(v), &self))
        }
    }
    
    pub fn is_false(v: &bool) -> bool { !v }

    pub fn parse_bool(v: &str) -> Option<bool> {
        match v.to_lowercase().as_str() {
            "yes" | "true"  | "1" => Some(true),
            "no"  | "false" | "0" => Some(false),
            _ => None
        }
    }
}
