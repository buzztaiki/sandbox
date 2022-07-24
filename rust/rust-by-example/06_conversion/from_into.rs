use std::convert::From;

#[derive(Debug)]
struct Number {
    value: i32,
}

impl From<i32> for Number {
    fn from(item: i32) -> Self {
        Number { value: item }
    }
}

fn main() {
    let num1 = Number::from(30);
    println!("My number is {:?}", num1);

    // Try removing the type declaration
    // Note: 型指定がないと何に into するか分からないからエラーになる。
    let num2: Number = 5.into();
    println!("My number is {:?}", num2);
}
