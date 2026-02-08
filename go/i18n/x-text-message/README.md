# x/text/messages の使い方

- `golang.org/x/text/message.Printer` を使ったコードを書く。フォーマット指定子は `fmt` パッケージと同じに書く。
- `gotext update` を実行する。詳細な引数は `Makefile` を参照。
- 実行すると `locales` ディレクトリに `out.gotext.json` が生まれてくる。
- これを翻訳者に翻訳してもらって、`messages.gotext.json` として保存する。
- 再度 `gotext update` すると、翻訳結果が埋め込まれた go ファイルが生成される。
- これをロードする事で、翻訳されたメッセージが使えるようになる。
