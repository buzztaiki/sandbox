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
