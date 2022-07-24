trait Handler {
    type Input;
    type Output;

    fn handle(&self, input: Self::Input) -> Self::Output;

    fn chain<Sink>(self, sink: Sink) -> ChainHandler<Self, Sink>
    where
        Sink: Handler,
        Self: Sized,
    {
        ChainHandler::new(self, sink)
    }
}

struct Incr;
impl Handler for Incr {
    type Input = u32;
    type Output = u32;
    fn handle(&self, input: Self::Input) -> Self::Output {
        input + 1
    }
}

struct ToString;
impl Handler for ToString {
    type Input = u32;
    type Output = String;

    fn handle(&self, input: Self::Input) -> Self::Output {
        input.to_string()
    }
}

struct Crub;
impl Handler for Crub {
    type Input = String;
    type Output = String;

    fn handle(&self, input: Self::Input) -> Self::Output {
        format!("🦀🦀{input}🦀🦀")
    }
}

struct ChainHandler<A, B> {
    a: A,
    b: B,
}

impl<A, B> ChainHandler<A, B>
where
    A: Handler,
    B: Handler,
{
    fn new(a: A, b: B) -> Self {
        Self { a, b }
    }
}

impl<A, B, M> Handler for ChainHandler<A, B>
where
    A: Handler<Output = M>,
    B: Handler<Input = M>,
{
    type Input = A::Input;
    type Output = B::Output;

    fn handle(&self, input: Self::Input) -> Self::Output {
        self.b.handle(self.a.handle(input))
    }
}

fn main() {
    let h = Incr.chain(ToString).chain(Crub);
    println!("{}", h.handle(10));
}
