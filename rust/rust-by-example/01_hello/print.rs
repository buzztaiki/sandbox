fn main() {
    // NOTE: print系のマクロは std::fmt で定義されているそう。
    // NOTE: see https://doc.rust-lang.org/std/fmt/
    println!("{} days", 31);
    println!("{0}, this is {1}. {1}, this is {0}", "Alice", "Bob");
    println!("{subject} {verb} {object}",
             object="the lazy dog",
             subject="the quick brown fox",
             verb="jumps over");

    println!("{} of {:b} people know binary", 1, 2); // 1 of 10 ...
    // lpad
    println!("{number:>width$}", number=1, width=6);
    // rpad
    println!("{number:<width$}", number=1, width=6);
    // lpad by 0
    println!("{number:>0width$}", number=1, width=6);
    println!("{number:>03}", number=1);
 
    println!("My name is {0}, {1} {0}", "Bond", "James");

    // NOTE: std::fmt::Display trait を実装する必要がある
    // struct Structure(i32);
    // println!("This struct `{}` won't print", Structure(3));

    let pi = 3.141592;
    println!("Pi is roughly {:.3}", pi);
    println!("Pi is roughly {:.width$}", pi, width=3);
}
