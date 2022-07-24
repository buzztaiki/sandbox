use std::fmt;

struct List(Vec<i32>);

impl fmt::Display for List {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        let vec = &self.0;
        // NOTE: ? を使うか try! を使わないと以下の警告がでる
        // = note: `#[warn(unused_must_use)]` on by default
        // = note: this `Result` may be an `Err` variant, which should be handled
        //
        // NOTE: try! マクロを使うとエラーが起きたときにそのエラーを返してくれるらしい。まだエラー習ってないけど、エラーも返値なのかしら。
        // https://doc.rust-jp.rs/rust-by-example-ja/std/result/question_mark.html を見ると Result 型が Either 型っぽく値とエラーを返せる様子。
        // つまりそこからエラーを取り出すわけだ。
        try!(write!(f, "["));

        for (count, v) in vec.iter().enumerate() {
            // ? は try? のシンタックスシュガー
            if count != 0 {
                write!(f, ", ")?;
            }
            write!(f, "{}", v)?;
        }

        write!(f, "]")
    }
}

fn main() {
    println!("{}", List(vec![1, 2, 3]));
}
