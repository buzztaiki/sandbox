use std::{env, process};

use input::Libinput;

fn main() {
    let mut input = Libinput::new_from_path(libinput::Interface);
    let args = env::args().collect::<Vec<_>>();
    if args.len() != 2 {
        eprintln!("usaege: {} <device>", args[0]);
        process::exit(1);
    }
    let path = args[1].as_str();

    let _device = input.path_add_device(path).unwrap();
    loop {
        input.dispatch().unwrap();
        for event in &mut input {
            println!("Got event: {:?}", event);
        }
    }
}
