fn exclusive_range(from: i64, to: i64) {
    for n in from..to {
        if n % 15 == 0 {
            println!("fizzbuzz");
        } else if n % 3 == 0 {
            println!("fizz");
        } else if n % 5 == 0 {
            println!("buzz");
        } else {
            println!("{}", n);
        }
    }
}

fn inclusive_range(from: i64, to: i64) {
    for n in from..=to {
        if n % 15 == 0 {
            println!("fizzbuzz");
        } else if n % 3 == 0 {
            println!("fizz");
        } else if n % 5 == 0 {
            println!("buzz");
        } else {
            println!("{}", n);
        }
    }
}

fn use_iter() {
    let names = vec!["Bob", "Frank", "Ferris"];

    for name in names.iter() {
        match name {
            // NOTE: iter の場合 & (リファレンス) にしないと怒られる
            &"Ferris" => println!("There is a rustacean among us!"),
            _ => println!("Hello {}", name),
        }
    }

    // NOTE: iter の場合コレクションを消費しないので、再度利用できる (まだよくわかってない)
    for name in names.iter() {
        println!("Hello {} again!", name);
    }
}

fn use_into_iter() {
    let names = vec!["Bob", "Frank", "Ferris"];

    for name in names.into_iter() {
        match name {
            // into_iter の場合 & (リファレンス) にすると怒られる
            "Ferris" => println!("There is a rustacean among us!"),
            _ => println!("Hello {}", name),
        }
    }

    // NOTE: into_iter の場合は消費するので、再利用できない (まだ消費が何かよくからない)
    // NOTE: 最初の into_iter で以下のコンパイルエラー:
    //     `names` moved due to this method call
    // NOTE: 次に names を使う所で以下のエラー:
    //     value borrowed here after move
    // NOTE: rustc --explain E0382 で詳細な説明が見れる。やさしい。でもむずい。
    // for name in names.iter() {
    //     println!("Hello {} again!", name);
    // }
}

fn use_iter_mat() {
    let mut names = vec!["Bob", "Frank", "Ferris"];

    // NOTE: iter_mut 使うと mutable な値が変える。pointer ですね。危険なのかな、これ。
    for name in names.iter_mut() {
        *name = match name {
            &mut "Ferris" => "There is a rustacean among us!",
            _ => "Hello",
        }
    }

    println!("names: {:?}", names);
}

fn main() {
    println!("## exclusive_range");
    exclusive_range(1, 11);
    println!("## inclusive_range");
    inclusive_range(1, 10);
    println!("## iter");
    use_iter();
    println!("## into_iter");
    use_into_iter();
    println!("## iter_mat");
    use_iter_mat();
}
