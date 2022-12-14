use std::{path::PathBuf, fs, fmt::Display};

use clap::{error::{ContextValue, ContextKind, ErrorKind}, builder::{TypedValueParser}};


#[derive(Clone)]
pub(crate) struct ExistingFileValueParser;

impl ExistingFileValueParser {

    fn iter_file_siblings(path: &PathBuf) -> Vec<PathBuf> {
        if let Some(parent) = path.parent() {
            if let Ok(entries) = fs::read_dir(parent) {
                return entries.filter_map(Result::ok)
                    .map(|e| e.path())
                    .filter(|f| f != path)
                    .collect();
            }
        }
    
        return Vec::new();
    }
    
    fn possible_files(path: &PathBuf) -> Vec<String> {
        let siblings = Self::iter_file_siblings(path);
        Self::did_you_mean(&path.to_string_lossy(), siblings.iter()
            .map(|f| f.to_string_lossy()))
    }
    
    fn did_you_mean<T, I>(v: &str, possible_values: I) -> Vec<String>
    where
        T: AsRef<str>,
        I: IntoIterator<Item = T>,
    {
        let mut candidates: Vec<(f64, String)> = possible_values
            .into_iter()
            .map(|pv| (strsim::jaro_winkler(v, pv.as_ref()), pv.as_ref().to_owned()))
            .filter(|(confidence, _)| *confidence > 0.8)
            .collect();
        candidates.sort_by(|a, b| a.0.partial_cmp(&b.0).unwrap_or(std::cmp::Ordering::Equal));
        candidates.into_iter().map(|(_, pv)| pv).collect()
    }

    fn clap_error(kind: ErrorKind, cmd: &clap::Command, arg: Option<&clap::Arg>, value: Option<String>) -> clap::error::Error {
        

        let mut err = clap::Error::new(kind).with_cmd(cmd);
        
        if let Some(value) = value {
            err.insert(ContextKind::InvalidValue, ContextValue::String(value));
        } 
        if let Some(arg) = arg {
            err.insert(ContextKind::InvalidArg, ContextValue::String(arg.to_string()));
        }

        err
    }

    fn custom_error(kind: ErrorKind, cmd: &clap::Command,  message: impl Display) -> clap::Error {
        clap::Error::raw(kind, message)
            .with_cmd(cmd)
    }
}

impl TypedValueParser for ExistingFileValueParser {
    type Value = PathBuf;

    fn parse_ref(
        &self,
        cmd: &clap::Command,
        arg: Option<&clap::Arg>,
        value: &std::ffi::OsStr,
    ) -> Result<Self::Value, clap::Error> {
        TypedValueParser::parse(self, cmd, arg, value.to_owned())
    }

    fn parse(
            &self,
            cmd: &clap::Command,
            arg: Option<&clap::Arg>,
            value: std::ffi::OsString,
        ) -> Result<Self::Value, clap::Error> {
            
        if value.is_empty() {
            return Err(Self::clap_error(ErrorKind::MissingRequiredArgument, cmd, arg, None));
        }
        
        let val = Self::Value::from(&value);

        if !val.is_file() {
            // let mut err = Self::error(ErrorKind::ValueValidation, cmd, arg, Some(value.to_string_lossy().to_string()), );
            
            let val_str = val.to_string_lossy().to_string();
            let arg_str = arg.map(ToString::to_string).unwrap_or_else(|| "...".to_owned());
            if val.exists() {
                return Err(Self::custom_error(ErrorKind::ValueValidation, cmd, 
                    format!("Invalid value for {}: Path is not a file: {}", arg_str, val_str)));
                
                // err.insert (ContextKind::Custom, ContextValue::String("Path is not a file".to_owned()));
            } else {
                let mut err = Self::clap_error(ErrorKind::ValueValidation, cmd, arg, Some(val_str));
                    // format!("Invalid value for {}: Path does not exist: {}", arg_str, val_str));
                // err.insert (ContextKind::Custom, ContextValue::String(.to_owned()));
                if let Some(suggestion) = Self::possible_files(&val).pop() {
                    err.insert(ContextKind::SuggestedValue, ContextValue::String(suggestion));
                }
                return Err(err)
            }

            // return Err(err)
        }

        Ok(val)
    }
}