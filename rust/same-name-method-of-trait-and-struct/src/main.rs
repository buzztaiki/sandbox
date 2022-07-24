trait Trait {
    fn name(&self) -> String;
    fn namae(&self) -> String;
    fn simei(&self) -> String;
}

struct Struct;

impl Struct {
    pub fn name(&self) -> String {
        "name".to_string()
    }
}

impl Trait for Struct {
    fn name(&self) -> String {
        // 同じ名前のメソッドがあった場合 Struct のメソッドが優先される (多分 self は Struct として修飾されてるから)
        format!("{}!!", self.name())

        // Struct 側のメソッドを呼ぶと明示したい場合は Struct::name のように関数表記を使う
        // format!("{}!!", Struct::name(self))

        // Trait 側のメソッドを呼ぶ場合も Trait::name のように型を明示して関数表記を使う (ちなみに unconditional recursion でエラーになる。偉い)
        // format!("{}!!", Trait::name(self))
    }

    fn namae(&self) -> String {
        // この場合も Struct 側が優先される (型を明示してもよい)
        format!("{}??", self.name())

        // Trait のメソッドを呼びたい場合
        // format!("{}??", Trait::name(self))
    }

    fn simei(&self) -> String {
        // namae は Struct に存在しないので Trait の方が呼ばれる
        format!("{}$$", self.namae())
    }
}

fn name<T: Trait>(x: &T) -> String {
    x.name()
}

fn main() {
    let x = Struct;
    println!("{}", name(&x));
    println!("{}", x.name());
    println!("{}", Trait::name(&x));
    println!("{}", x.namae());
    println!("{}", x.simei());
}
