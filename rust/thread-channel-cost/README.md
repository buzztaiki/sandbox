# スレッドとか channel ってどれくらいコストかかるのか知りたい

雑にベンチマークを取ってみたというだけ。

結果としてはこうなった

```console
 $  cargo bench
    Finished bench [optimized] target(s) in 0.01s
     Running unittests (target/release/deps/thread_channel_cost-60762c316631bc2e)

running 0 tests

test result: ok. 0 passed; 0 failed; 0 ignored; 0 measured; 0 filtered out; finished in 0.00s

     Running unittests (target/release/deps/benchmark-c3a1c608a31553fd)

running 6 tests
test bench_channel        ... bench:      31,002 ns/iter (+/- 38,374)
test bench_collect        ... bench:       1,765 ns/iter (+/- 116)
test bench_spawn          ... bench:  16,182,450 ns/iter (+/- 2,079,421)
test bench_thread_pool    ... bench:     265,300 ns/iter (+/- 114,097)
test bench_tokio          ... bench:     202,443 ns/iter (+/- 99,414)
test bench_tokio_blocking ... bench:   2,425,352 ns/iter (+/- 334,934)

test result: ok. 0 passed; 0 failed; 0 ignored; 6 measured
```
