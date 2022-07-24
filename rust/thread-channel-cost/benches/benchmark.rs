use std::sync::mpsc::channel;
use std::thread::spawn;

use bencher::{benchmark_group, benchmark_main, Bencher};
use threadpool::ThreadPool;

fn bench_collect(bench: &mut Bencher) {
    let xs = text();
    bench.iter(|| {
        let _xs = xs
            .chars()
            .map(|x| x.to_ascii_lowercase())
            .collect::<Vec<_>>();
    });
}

fn bench_spawn(bench: &mut Bencher) {
    let xs = text();
    bench.iter(|| {
        for x in xs.chars() {
            spawn(move || {
                let _x = x.to_ascii_lowercase();
            });
        }
    })
}

fn bench_channel(bench: &mut Bencher) {
    let (tx, rx) = channel::<char>();
    spawn(move || {
        for x in rx {
            let _x = x.to_ascii_lowercase();
        }
    });

    let xs = text();
    bench.iter(|| {
        for x in xs.chars() {
            tx.send(x).unwrap();
        }
    });
}

fn bench_thread_pool(bench: &mut Bencher) {
    let pool = ThreadPool::new(5);
    let xs = text();
    bench.iter(|| {
        for x in xs.chars() {
            pool.execute(move || {
                let _x = x.to_ascii_lowercase();
            });
        }
    })
}

fn bench_tokio(bench: &mut Bencher) {
    tokio::runtime::Builder::new_multi_thread()
        .enable_all()
        .build()
        .unwrap()
        .block_on(async {
            let xs = text();
            bench.iter(|| {
                for x in xs.chars() {
                    tokio::task::spawn(async move {
                        let _x = x.to_ascii_lowercase();
                    });
                }
            })
        });
}

fn bench_tokio_blocking(bench: &mut Bencher) {
    tokio::runtime::Builder::new_multi_thread()
        .enable_all()
        .build()
        .unwrap()
        .block_on(async {
            let xs = text();
            bench.iter(|| {
                for x in xs.chars() {
                    tokio::task::spawn_blocking(move || {
                        let _x = x.to_ascii_lowercase();
                    });
                }
            })
        });
}

fn text() -> String {
    vec!['a'; 1024].iter().collect()
}

benchmark_group!(
    benches,
    bench_collect,
    bench_spawn,
    bench_channel,
    bench_thread_pool,
    bench_tokio,
    bench_tokio_blocking,
);
benchmark_main!(benches);
