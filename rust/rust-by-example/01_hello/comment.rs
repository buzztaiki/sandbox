/// ドキュメンテーションコメントは `///` らしい。markdown が使えるらしい。
/// ドキュメントを生成するには `rustdoc comment.rs` 等すればよい。
pub fn hello() {
    // 普通のコメント

    /* 
    普通のブロックコメント 
     */
    println!("hello rust!");
}

fn main() {
    hello();
}
