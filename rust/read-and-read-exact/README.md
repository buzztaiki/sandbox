# read と read_exact の違い

基本は `read` で、絶対にその長さを読める事が保証できていて、必ず読みたい場合だけ `read_exact` とかするのがいいんじゃないかな多分。

## 挙動

テストを書いて調べて感じだと
- ストリーム >= バッファ の場合はどちらの挙動も基本は同じ
- ストリーム < バッファ の場合は、`read_exact` はエラー


## 実装

`read` のドキュメントには

https://doc.rust-lang.org/std/io/trait.Read.html#tymethod.read

> It is not an error if the returned value n is smaller than the buffer size, even when the reader is not at the end of the stream yet. This may happen for example because fewer bytes are actually available right now (e. g. being close to end-of-file) or because read() was interrupted by a signal.

とある。

`read_exact` には

https://doc.rust-lang.org/std/io/trait.Read.html#method.read_exact

> Read the exact number of bytes required to fill buf.

さらに

> If this function encounters an error of the kind ErrorKind::Interrupted then the error is ignored and the operation will continue.

とある。

また、実装を見ると以下のようなコードになってる。

```
    while !buf.is_empty() {
        match this.read(buf) {
...
```
