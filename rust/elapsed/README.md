# Ruby の Benchmark みたいな crate

```rust
use core::time;
use std::thread::sleep;

use elapsed::Elapsed;

fn main() {
    let mut elp = Elapsed::new();
    elp.measure("whole", |elp| {
        for _ in 1..10000 {
            elp.measure("sleep", |_| {sleep(time::Duration::from_micros(1))});
        }
    });
    println!("{}", elp);
}
```

```
cargo run
    Finished dev [unoptimized + debuginfo] target(s) in 0.00s
     Running `target/debug/elapsed`
whole: 703 ms / 1 times
sleep: 634 ms / 9999 times
```



