fn main() {
    // Because of the annotation, the compiler knows that `elem` has type u8.
    let elem = 5u8;

    // Create an empty vector (a growable array).
    let mut vec = Vec::new();
    // At this point the compiler doesn't know the exact type of `vec`, it
    // just knows that it's a vector of something (`Vec<_>`).

    // NOTE: mut を付けないと以下のエラーになる。変数の代入可否だけじゃないって事なんだな。
    // error[E0596]: cannot borrow `vec` as mutable, as it is not declared as mutable
    //   --> inference.rs:11:5
    //    |
    // 6  |     let vec = Vec::new();
    //    |         --- help: consider changing this to be mutable: `mut vec`
    // ...
    // 11 |     vec.push(elem);
    //    |     ^^^ cannot borrow as mutable
    


    // Insert `elem` in the vector.
    vec.push(elem);
    // Aha! Now the compiler knows that `vec` is a vector of `u8`s (`Vec<u8>`)
    // TODO ^ Try commenting out the `vec.push(elem)` line


    println!("{:?}", vec);

    // NOTE: push しないと方が決まらないよっていうエラーになる
    // error[E0282]: type annotations needed for `std::vec::Vec<T>`
    //  --> inference.rs:6:19
    //   |
    // 6 |     let mut vec = Vec::new();
    //   |         -------   ^^^^^^^^ cannot infer type for type parameter `T`
    //   |         |
    //   |         consider giving `vec` the explicit type `std::vec::Vec<T>`, where the type parameter `T` is specified

    // vec.push(5i8);
    // NOTE 違う型をいれると
    // error[E0308]: mismatched types
    //   --> inference.rs:24:14
    //    |
    // 24 |     vec.push(5i8);
    //    |              ^^^ expected `u8`, found `i8`
    //    |
}
