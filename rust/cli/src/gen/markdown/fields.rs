use shoutrrr::spec::{SpecProp, PropType};

use std::fmt::Write;


pub fn write_primary<T: Write>(output: &mut T, prop: &SpecProp) -> eyre::Result<()> {
    write!(output, "*  __{}__", prop.name)?;

    if !prop.description.is_empty() {
        write!(output," - {}", prop.description)?;
    }

    if prop.required {
        writeln!(output, " (**Required**)  ")?;
    } else {
        write!(output, "  \n  Default: ")?;
        match &prop.default_value {
            None => writeln!(output, "*empty*  "),
            Some(def_val) => {
                // if prop.prop_type == PropType::Toggle {
                    /*
                    defaultValue, _ := format.ParseBool(field.DefaultValue, false)
                if defaultValue {
                    sb.WriteString("✔ ")
                } else {
                    sb.WriteString("❌ ")
                }
                     */
                // }
                writeln!(output, "`{}`  ", def_val)
            },
        }?;
    }
    Ok(())
}

pub fn write_extras<T: Write>(output: &mut T, prop: &SpecProp) -> eyre::Result<()> {

    // Skip primary alias (as it's the same as the field name)
    if let Some((_, rest)) = prop.keys.split_first() {
        if !rest.is_empty() {
            writeln!(output, "  Aliases: `{}`  ", rest.join("`, `"))?;
        }
    }

    if let PropType::Option(opt) = &prop.prop_type {
        writeln!(output, "  Possible values: `{}`  ", opt.values.join("`, `"))?;
    }

    Ok(())
}