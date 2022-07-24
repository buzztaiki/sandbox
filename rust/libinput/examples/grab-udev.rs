use std::os::unix::prelude::*;
use std::path::Path;

use input::event::{DeviceEvent, EventTrait};
use input::{Event, Libinput, LibinputInterface};

struct Interface {
    iface: libinput::Interface,
}

impl LibinputInterface for Interface {
    fn open_restricted(&mut self, path: &Path, flags: i32) -> Result<RawFd, i32> {
        let fd = self.iface.open_restricted(dbg!(path), flags)?;
        let device = evdev::Device::open(path).map_err(|e| e.raw_os_error().unwrap())?;
        if let Some("Kinesis Advantage2 Keyboard") = device.name() {
            libinput::grab(fd)?;
        }
        Ok(fd)
    }
    
    fn close_restricted(&mut self, fd: RawFd) {
        self.iface.close_restricted(fd)
    }
}

fn main() {
    let mut input = Libinput::new_with_udev(Interface{ iface: libinput::Interface });
    input.udev_assign_seat("seat0").unwrap();
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
