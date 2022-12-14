use std::collections::HashMap;

use shoutrrr::spec::{URLPart, SpecProp};

use std::fmt::Write;


pub fn get_slugs(scheme: &str, props: &&HashMap<URLPart, &SpecProp>) -> HashMap<URLPart, (String, String)> {
    let mut slug_map = HashMap::new();

    let has_port = props.contains_key(&URLPart::Port);
    let has_pass = props.contains_key(&URLPart::Password);
    let has_auth = props.contains_key(&URLPart::User) || has_pass;

    for part in URLPart::iter() {
        

        if let Some((slug, separator)) = match (part, props.get(&part)
                .and_then(|p| Some(p.name.as_str()))) {
            (URLPart::Scheme, _) => Some((scheme, "://")),
            (URLPart::Host, None) if has_port => Some((scheme, ":")),
            (URLPart::Host, None) => Some((scheme, "/")),
            (URLPart::Host, Some(name)) if has_port => Some((name, ":")),
            (URLPart::Host, Some(name)) => Some((name, "/")),
            (URLPart::User, Some(name)) if !has_pass => Some((name, "@")),
            (URLPart::User, Some(name)) if has_pass => Some((name, ":")),
            (URLPart::User, None) if has_auth && !has_pass => Some(("", "@")),
            (URLPart::User, None) if has_auth && has_pass => Some(("", ":")),
            (URLPart::Password, Some(name)) => Some((name, "@")),
            (URLPart::Port, Some(_)) => Some(("port", "/")),
            (_, None) => None,
            (_, Some(name)) => Some((name, "/")),
        } {
            slug_map.insert(part, (slug.to_lowercase().to_owned(), separator.to_owned()));
        }
    }

    slug_map
}

pub fn write_url_map<T: Write>(output: &mut T, highlight_parts: &Vec<URLPart>, slug_map: &HashMap<URLPart, (String, String)>) -> eyre::Result<()> {
        
    let mut prev_sep = None;
    for part in URLPart::iter() {

        if let Some(entry) = slug_map.get(&part) {
            let (slug, sep) = entry;
        
            if let Some(prev_sep) = prev_sep {
                write!(output, "{prev_sep}")?;
            }
        
            if highlight_parts.contains(&part) {
                write!(output, "<strong>{slug}</strong>")?;
            } else {
                write!(output, "{slug}")?;
            }
   
            prev_sep = Some(sep);
        }
    }
    if let Some(prev_sep) = prev_sep {
        if prev_sep == "/" {
            write!(output, "{prev_sep}")?;
        }
    }
    Ok(())
}