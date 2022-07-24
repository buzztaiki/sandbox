#![feature(arbitrary_self_types)]

use std::ops::Deref;
use std::rc::Rc;

#[derive(Default, Debug)]
struct A {
    x: i64
}

struct B(A);
impl Deref for B {
    type Target = A;

    fn deref(&self) -> &Self::Target {
        &self.0
    }
}

impl A {
    fn with_rc(self: Rc<Self>) -> i64 {
        self.x
    }

    fn with_deref(self: &B) -> i64 {
        self.x
    }
}

fn main() {
    let a = Rc::new(A::default());
    let x = a.clone().with_rc();
    dbg!(x);

    let a = B(A::default());
    let x = a.with_deref();
    let y = a.with_deref();
    dbg!((x, y));
}
