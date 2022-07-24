use std::collections::HashMap;

// println!("{}", x(map))
//                  ^^^ cyclic type of infinite size
fn main() {
    let mut map = HashMap::new();
    f(map);
    let &x = map.get(&3).unwrap();
    println!("{}", x(map))
}

fn f<F>(mut map: HashMap<i32, Box<dyn Fn(HashMap<i32, Box<F>>) -> i32>>)
where
    F: Fn(HashMap<i32, Box<F>>) -> i32,
{
    map.insert(1, Box::new(|_| 10));
    map.insert(2, Box::new(|_| 10));
    map.insert(3, Box::new(|map| map.get(&1).unwrap()(map)));
}
