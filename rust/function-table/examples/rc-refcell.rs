use std::{cell::RefCell, collections::HashMap, rc::Rc};

// rust-jp で教えてもらったぱたーん。RCべんり。
// ただ、実行時のオーバーヘッドがかかるだけで良い事ないと思うよって話だった。
fn main() {
    let map = Rc::new(RefCell::new(HashMap::<_, Box<dyn Fn() -> i32>>::new()));
    let map1 = Rc::clone(&map);
    map.borrow_mut().insert("a", Box::new(|| 10));
    map.borrow_mut().insert("b", Box::new(|| 20));
    map.borrow_mut()
        .insert("c", Box::new(move || map1.borrow().get("a").unwrap()() + 1));
    dbg!(map.borrow().get("c").unwrap()());
}
