use std::collections::HashMap;
use std::hash::Hash;

#[derive(Default)]
struct Map<'a, K, V> {
    map: HashMap<K, Box<dyn 'a + Fn(&Self) -> V>>,
}

impl<'a, K: Eq + Hash, V> Map<'a, K, V> {
    fn insert<F: 'a + Fn(&Self) -> V>(&mut self, k: K, f: F) {
        self.map.insert(k, Box::new(f));
    }

    fn call(&self, k: K) -> Option<V> {
        self.map.get(&k).map(|f| f(self))
    }
}

fn main() {
    let mut map = Map::default();
    map.insert(1, |_| 10);
    map.insert(2, |_| 20);
    map.insert(3, |m| m.call(1).unwrap() + 1);
    println!("{:?}", map.call(3).unwrap());
}
