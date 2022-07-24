# トレイトと構造体のメソッド名が同じ場合

https://doc.rust-jp.rs/book-ja/ch19-03-advanced-traits.html#明確化のためのフルパス記法-同じ名前のメソッドを呼ぶ

- struct 側が優先。
- trait にしかないメソッドの場合は trait の方。
- フルパス記法 (`Struct::name(self)` や `Trait::name(self)`) で書く事で明示できるし選択できる。
- `(self as dyn Trait).name()` もできるっちゃできるけど、`as` はレシーバがない時に使うくらいで良い。
