fn main() {
    // Rustのプログラムは（ほとんどの場合）文(statement)の連続
    // statement
    // statement
    // ...

    // 変数束縛は文
    let x = 5;

    // ; で終わる式も文
    x;
    x + 1;

    // コードブロックは式
    // NOTE: コードブロックは式なのは便利だ。よい。
    let y = {
        let x2 = x * x;
        let x3 = x2 * x;
        x + x2 + x3
    };

    let z = {
        // NOTE ; で終わるからこれは () になる。文の値は () なので。その後使うところでコンパイルエラーとかになるから、大体の場合は気付けそうではある。
        2 * x;
    };

    println!("x: {:?}", x);
    println!("y: {:?}", y);
    println!("z: {:?}", z);
}
