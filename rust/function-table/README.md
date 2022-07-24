# Rust で関数テーブルみたいなのをやりたい

つまり、こういうやつ

```ruby
map = {}
map[1] = ->{ 1 }
map[2] = ->{ 10 }
map[3] = ->{ map[1].call}
puts map[3].call
```

借用チェッカーに怒られてなかなかうまくいかない。

[examples/wrap-struct.rs](examples/wrap-struct.rs) みたいにすれば、それっぽい感じになる。ちょっときびしい。
