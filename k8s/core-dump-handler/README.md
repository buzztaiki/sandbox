WSL + Kind でうまくいかない。

WSl 自体がコンテナで動いてて、Kind もコンテナで動いてる。

なので、以下の結果はどこで見ても全部同じになる。それはそう。

```
❯❯ sysctl kernel/core_pattern
kernel.core_pattern = /tmp/core
```

そのときに、どう動くのかみたいな事になりそう。
あと、なんでか知らんけど WSL で coredump 出てくれない気がする。
