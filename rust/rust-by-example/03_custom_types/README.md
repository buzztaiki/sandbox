enum は scala や kotolin の sealed class と雰囲気似てる。
LinkedList の例とかまさにそう。

```
enum List {
    Cons(u32, Box<List>),
    Nil,
}
```


