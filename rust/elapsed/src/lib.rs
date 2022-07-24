use std::collections::HashMap;
use std::time;

#[derive(Default, Debug)]
struct Value {
    count: usize,
    elapsed: time::Duration,
}

#[derive(Default, Debug)]
pub struct Elapsed {
    map: HashMap<String, Value>,
}

impl Elapsed {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn measure(&mut self, name: &str, f: impl FnOnce(&mut Elapsed)) {
        let start = time::Instant::now();
        let mut child = Elapsed::new();
        f(&mut child);
        let value = Value {
            count: 1,
            elapsed: start.elapsed(),
        };

        for (k, v) in child
            .map
            .into_iter()
            .chain(std::iter::once((name.to_string(), value)))
        {
            let e = self.map.entry(k).or_default();
            *e = Value {
                count: e.count + v.count,
                elapsed: e.elapsed + v.elapsed,
            };
        }
    }
}

impl std::fmt::Display for Elapsed {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        for (k, v) in self.map.iter() {
            writeln!(f, "{}: {} ms / {} times", k, v.elapsed.as_millis(), v.count)?;
        }
        Ok(())
    }
}
