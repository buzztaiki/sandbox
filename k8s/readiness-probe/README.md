# ReadinessProbe の検証

https://kubernetes.io/ja/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/

- 一定期間で `/tmp/healthy` を作ったり削除したりする簡単な go のプログラムを書いた。
  - `--no-recover` を付けると削除したらそのままにする
- ReadinessProbe がこのファイルを監視するようにした。

## Readiness Probe が deployment の unavailable に影響するか

`--no-recover` オプションを付けずに確認した。
- readiness probe が失敗の状態 (not readiness) になると unavailable が一つ増える。
- readiness probe が成功に戻ると available になる。


## ずっと unavailable だった場合に 新しい pod は生まれてくるか

`--no-recover` オプションを付けて一度死んだら返ってこないようにして検証した。
- pod は新しく生まれてこない。
- ようするに、普通に liveness probe も指定しましょうということ。
