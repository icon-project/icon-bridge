use super::*;

#[derive(Debug, Serialize)]
pub struct HashedCollection<T: PartialEq + HASH + Eq>(pub HashSet<T>);

impl<T: PartialEq + HASH + Eq> HashedCollection<T> {
    pub fn new() -> HashedCollection<T> {
        Self(HashSet::new())
    }
    pub fn add(&mut self, element: T) {
        self.0.insert(element);
    }
}

impl<T: PartialEq + HASH + Eq> PartialEq for HashedCollection<T> {
    fn eq(&self, other: &Self) -> bool {
        self.0 == other.0
    }
}

impl<T: PartialEq + HASH + Eq> Default for HashedCollection<T> {
    fn default() -> Self {
        Self::new()
    }
}

#[derive(Debug)]
pub struct HashedValue(Value);

impl HASH for HashedValue {
    fn hash<H: HASHER>(&self, state: &mut H) {
        self.0.to_string().hash(state)
    }
}

impl PartialEq for HashedValue {
    fn eq(&self, other: &Self) -> bool {
        self.0 == other.0
    }
}

impl Eq for HashedValue {}

impl From<Value> for HashedValue {
    fn from(value: Value) -> Self {
        Self(value)
    }
}

impl FromIterator<Value> for HashedCollection<HashedValue> {
    fn from_iter<I: IntoIterator<Item = Value>>(iter: I) -> Self {
        let mut c = HashedCollection::new();
        for i in iter {
            c.add(HashedValue(i));
        }
        c
    }
}
