## エラーの基本
- https://doc.rust-jp.rs/book-ja/ch09-00-error-handling.html
- https://doc.rust-lang.org/std/error/trait.Error.html

## Result と Option

エラー付きの値とヌルを許容する値。Rust は基本エラーを返り値で表現する。

- Result: 結果またはエラー。
- Option: 結果はたは無し。

道具:

- `?` (`try!`): Ok なら値を返して Err だった場合に関数を抜けてくれる。何故か Option でも使える。使い勝手は for + flatmap と雰囲気似てなくもないけど、実体は早期リターン。
- `match`
- `if let`
- `match return`: 外れを引いたときに 早期リターンすると式の評価を抜けて関数に値を返せる。けど素直に `?` 使った方が良いと思う。
- `unwrap` series: `unwrap_or`, `unwrap_or_else`, `unwrap_or_defualt` 等を使わないと panic
- `Result#ok -> Option`, `Result#err -> Option`: ok, err の片方だけを取る
- `map_or`, `map_or_else`: Option, Result から値を取り出す
- `and_then`: flatmap
- `Option#ok_or -> Result`, `Option#ok_or_else -> Result`: Option から Result にする


## 複数の型のエラーを扱う方法について
ある関数から利用してる関数群が複数のエラーを返し得るときに、なるべく `?` を使いたい。毎回 match させるのはさすがにつらいし。if let return unless unwrap はしたくない。
どうするか。

今の rust のエラーは基本的には `error::Error` を実装している (必須ではないが)。
つまり `error::Error` との相互変換として色々構築していく事になる。

- `Box<dyn error::Error>` を使う:
  - 簡単。真面目に扱う場合いくつか気を付ける事があるらしく、外部ライブラリが使える状況では anyhow を使うと良いそう。
  - see: https://doc.rust-lang.org/rust-by-example/error/multiple_error_types/boxing_errors.html
- 複数のエラーを表現できる enum を定義して From, Display, error:Error を実装する:
  - わりと手間がかかる。thiserror を使うと手間がはぶける。
  - `?` で他のエラーから独自エラーに変換するために From を実装する。
  - `error:Error` として扱えるようにするために `Display` も実装する。
  - see: https://doc.rust-lang.org/std/error/trait.Error.html
- anyhow
  - `Box<dyn error::Error>` の代替。
  - https://crates.io/crates/anyhow
- thiserror
  - エラーの実装を楽にしてくれる。
  - https://crates.io/crates/thiserror

以下のあたりが参考になる
- https://qiita.com/legokichi/items/d4819f7d464c0d2ce2b8
- https://nick.groenen.me/posts/rust-error-handling/
- https://dev.to/seanchen1991/a-beginner-s-guide-to-handling-errors-in-rust-40k2

## Option に対する ? について

なんでか知らないけど Option に `?` が使える。

```rust
    /// Applies the "?" operator. A return of `Ok(t)` means that the
    /// execution should continue normally, and the result of `?` is the
    /// value `t`. A return of `Err(e)` means that execution should branch
    /// to the innermost enclosing `catch`, or return from the function.
    ///
    /// If an `Err(e)` result is returned, the value `e` will be "wrapped"
    /// in the return type of the enclosing scope (which must itself implement
    /// `Try`). Specifically, the value `X::from_error(From::from(e))`
    /// is returned, where `X` is the return type of the enclosing function.
    #[lang = "into_result"]
    #[unstable(feature = "try_trait", issue = "42327")]
    fn into_result(self) -> Result<Self::Ok, Self::Error>;

    /// Wrap an error value to construct the composite result. For example,
    /// `Result::Err(x)` and `Result::from_error(x)` are equivalent.
    #[lang = "from_error"]
    #[unstable(feature = "try_trait", issue = "42327")]
    fn from_error(v: Self::Error) -> Self;

    /// Wrap an OK value to construct the composite result. For example,
    /// `Result::Ok(x)` and `Result::from_ok(x)` are equivalent.
    #[lang = "from_ok"]
    #[unstable(feature = "try_trait", issue = "42327")]
    fn from_ok(v: Self::Ok) -> Self;
```

そして `Option` は `ops::Try` を実装している。つまり `?` を使う事で `Result` に変換する事ができる。
なのだけど、まだ unstable なので使えない。

そうすると `?` が `Option` で使える理由はなんだろう。`impl From<Result> for Option` を実装しているから Option -> Result -> Option といった変換がかかったのかと邪推してたのだけど。

参考:
- https://github.com/rust-lang/rust/issues/42327

