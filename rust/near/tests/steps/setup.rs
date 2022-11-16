pub use std::collections::HashSet;
use test_helper::types::Context;

pub static NEW_CONTEXT: fn() -> Context = || Context::new();
