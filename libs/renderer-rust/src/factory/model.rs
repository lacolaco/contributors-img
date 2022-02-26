use std::fmt;

use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug)]
pub struct Contributor {
    pub avatar_url: String,
    pub contributions: i32,
    pub html_url: String,
    pub id: i32,
    pub login: String,
}

impl fmt::Display for Contributor {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "Contributor: {}", self)
    }
}
