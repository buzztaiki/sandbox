use std::collections::HashMap;

// エラーとしては以下の二つがおきる。わりと積んでる気がする。

// map のライフタイム
//     let f = || map.get(&1).unwrap()();
//             -- ^^^ borrowed value does not live long enough
//             |
//             value captured here
//
// }
// -
// |
// `map` dropped here while still borrowed
// borrow might be used here, when `map` is dropped and runs the destructor for type `HashMap<i32, Box<dyn Fn() -> i32>>`

// Fn に immutable borrow した後の mutable borrow
//     let f = || map.get(&1).unwrap()();
//             -- --- first borrow occurs due to use of `map` in closure
//             |
//             immutable borrow occurs here
//     {
//         map.insert(3, Box::new(f));
//         ^^^^------^^^^^^^^^^^^^^^^
//         |   |
//         |   immutable borrow later used by call
//         mutable borrow occurs here

fn main() {
    let mut map = HashMap::<i32, Box<dyn Fn() -> i32>>::new();
    map.insert(1, Box::new(|| 10));
    map.insert(2, Box::new(|| 10));
    let f = || map.get(&1).unwrap()();
    map.insert(3, Box::new(f));
}
