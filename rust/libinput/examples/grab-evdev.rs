use std::os::unix::prelude::*;
use std::path::Path;
use std::{env, process};

use input::event::{DeviceEvent, EventTrait};
use input::{Event, Libinput, LibinputInterface};

struct Interface {
    iface: libinput::Interface,
}


impl LibinputInterface for Interface {
    fn open_restricted(&mut self, path: &Path, flags: i32) -> Result<RawFd, i32> {
        let fd = self.iface.open_restricted(path, flags)?;
        libinput::grab(fd)?;
        Ok(fd)
    }

    fn close_restricted(&mut self, fd: RawFd) {
        self.iface.close_restricted(fd)
    }
}


fn main() {
    let mut input = Libinput::new_from_path(Interface{ iface: libinput::Interface });
    let args = env::args().collect::<Vec<_>>();
    if args.len() != 2 {
        eprintln!("usaege: {} <device>", args[0]);
        process::exit(1);
    }

    let path = args[1].as_str();
    input.path_add_device(path).unwrap();

    loop {
        input.dispatch().unwrap();
        for event in &mut input {
            match event {
                Event::Device(DeviceEvent::Added(ev)) => {
                    println!("added: {}", ev.device().name());
                }
                Event::Device(DeviceEvent::Removed(ev)) => {
                    println!("removed: {}", ev.device().name());
                }
                _ => {
                    println!("Got event: {:?}", event);
                }
            }
        }
    }
}
