入出力がジェネリックな handler-chain (filter-chain) のサンプル。

こんな感じで使えるやつ。

```rust
fn main() {
    let h = Incr.chain(ToString).chain(Crub);
    println!("{}", h.handle(10));
}
```
