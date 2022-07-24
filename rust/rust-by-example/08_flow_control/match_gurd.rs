fn main() {
    let pair = (2, -2);

    println!("Tell me about {:?}", pair);
    match pair {
        // NOTE: パターンマッチの後ろに if を書くとガード節になる。
        // NOTE: rustのガード節はあくまで記法の問題だけなのかな。
        (x, y) if x == y => println!("These are twins"),
        // The ^ `if condition` part is a guard
        (x, y) if x + y == 0 => println!("Antimatter, kaboom!"),
        (x, _) if x % 2 == 1 => println!("The first one is odd"),
        _ => println!("No correlation..."),
    }
}
