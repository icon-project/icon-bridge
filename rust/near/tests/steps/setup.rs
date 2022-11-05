use test_helper::types::Context;
pub use std::collections::HashSet;

pub static NEW_CONTEXT: fn() -> Context = || Context::new();
