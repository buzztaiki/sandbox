enum Trampoline<T> {
    More(Box<dyn FnOnce() -> Trampoline<T>>),
    Done(T),
}

impl<T> Trampoline<T> {
    fn call(mut self) -> T {
        loop {
            match self {
                Self::More(f) => self = f(),
                Self::Done(x) => return x,
            }
        }
    }
}

fn even(n: u64) -> Trampoline<bool> {
    match n {
        0 => Trampoline::Done(true),
        _ => Trampoline::More(Box::new(move || odd(n - 1)))
    }
}

fn odd(n: u64) -> Trampoline<bool> {
    match n {
        0 => Trampoline::Done(false),
        _ => Trampoline::More(Box::new(move || even(n - 1)))
    }
}

fn main() {
    dbg!(even(100000000).call());
}
