use std::collections::{HashMap, HashSet};
use std::io;
use std::os::unix::prelude::*;
use std::path::Path;

use evdev::uinput;
use input::event::keyboard::{KeyState, KeyboardEventTrait};
use input::event::pointer::ButtonState;
use input::event::{DeviceEvent, EventTrait, KeyboardEvent, PointerEvent};
use input::{Device, DeviceConfigResult, Event, Libinput, LibinputInterface};

struct Interface {
    iface: libinput::Interface,
    devnames: HashSet<String>,
}

type ConfigFn = Box<dyn Fn(&mut Device) -> DeviceConfigResult>;

impl Interface {
    fn new(config: &HashMap<String, ConfigFn>) -> Self {
        Self {
            iface: libinput::Interface,
            devnames: config.keys().map(|x| x.to_string()).collect(),
        }
    }
}

impl LibinputInterface for Interface {
    fn open_restricted(&mut self, path: &Path, flags: i32) -> Result<RawFd, i32> {
        let fd = self.iface.open_restricted(path, flags)?;
        let device = evdev::Device::open(path).map_err(|e| e.raw_os_error().unwrap())?;

        if device.name().and_then(|x| self.devnames.get(x)).is_some() {
            libinput::grab(fd)?
        }
        Ok(fd)
    }

    fn close_restricted(&mut self, fd: RawFd) {
        self.iface.close_restricted(fd)
    }
}

fn main() {
    let mut config = HashMap::<String, ConfigFn>::new();
    config.insert(
        "Kinesis Advantage2 Keyboard".to_string(),
        Box::new(|_| Ok(())),
    );
    config.insert(
        "Kensington Expert Wireless TB Mouse".to_string(),
        Box::new(|device| device.config_middle_emulation_set_enabled(true)),
    );

    let mut input = Libinput::new_with_udev(Interface::new(&config));
    input.udev_assign_seat("seat0").unwrap();

    let mut vdevices = HashMap::new();
    loop {
        input.dispatch().unwrap();
        for event in &mut input {
            match event {
                Event::Device(DeviceEvent::Added(ev)) => {
                    let mut dev = ev.device();
                    println!("added: {}", dev.name());

                    if let Some(config_fn) = config.get(dev.name()) {
                        match new_vertial_device(&dev) {
                            Ok(vdev) => {
                                vdevices.insert(dev.sysname().to_string(), vdev);
                                config_fn(&mut dev).unwrap();
                                println!("vdevice created");
                            }
                            Err(e) => println!("failed to create vdevice {:?}", e),
                        }
                    }
                }
                Event::Keyboard(ev) => {
                    if let Some(ref mut vdev) = vdevices.get_mut(ev.device().sysname()) {
                        vdev.emit(convert_to_evdev_keyboard_event(&ev).as_ref()).unwrap();
                    }
                },
                Event::Pointer(ev) => {
                    if let Some(ref mut vdev) = vdevices.get_mut(ev.device().sysname()) {
                        match convert_to_evdev_pointer_event(&ev) {
                            Some(x) => vdev.emit(x.as_ref()).unwrap(),
                            None => println!("event {:?} not supported", ev),
                        }

                        ;
                    }
                }
                _ => {
                    println!("Got event: {:?}", event);
                }
            }
        }
    }
}

fn new_vertial_device(device: &Device) -> io::Result<uinput::VirtualDevice> {
    let mut keys = evdev::AttributeSet::<evdev::Key>::new();
    for code in (0..libc::KEY_CNT).map(|x| x as u16) {
        keys.insert(evdev::Key::new(code));
    }

    let mut axes = evdev::AttributeSet::<evdev::RelativeAxisType>::new();
    for code in (0..libc::REL_CNT).map(|x| x as u16) {
        axes.insert(evdev::RelativeAxisType(code));
    }

    uinput::VirtualDeviceBuilder::new()?
        .name(format!("virtual-{}-{}", device.sysname(), device.name()).as_str())
        .with_keys(&keys)?
        .with_relative_axes(&axes)?
        .build()
}

fn convert_to_evdev_keyboard_event(ev: &KeyboardEvent) -> Vec<evdev::InputEvent> {
    vec![
        evdev::InputEvent::new(
            evdev::EventType::KEY,
            ev.key() as u16,
            match ev.key_state() {
                KeyState::Pressed => 1,
                KeyState::Released => 0,
            },
        )
    ]
}

fn convert_to_evdev_pointer_event(pev: &PointerEvent) -> Option<Vec<evdev::InputEvent>> {
    match pev {
        PointerEvent::Motion(ev) => {
            // dx, dy は accelarate された値。
            // comositor と組合せると二重に accel される事になってしまうが。どうしようかな。
            Some(vec![
                evdev::InputEvent::new(
                    evdev::EventType::RELATIVE,
                    evdev::RelativeAxisType::REL_X.0,
                    ev.dx() as i32,
                ),
                evdev::InputEvent::new(
                    evdev::EventType::RELATIVE,
                    evdev::RelativeAxisType::REL_Y.0,
                    ev.dy() as i32,
                ),
            ])
        },
        // タッチパッドの場合これになるんだけど、libinput event を evdev event に再構築するの無理くさいかもしれないぞ
        // 無理ではないかもしれないけど、多分すごく大変
        // 再構築するなら、libinput のソースみちゃった方がはやいかもな。
        PointerEvent::MotionAbsolute(_) => None,
        PointerEvent::Button(ev) => Some(
            vec![
                evdev::InputEvent::new(
                    evdev::EventType::KEY,
                    ev.button() as u16,
                    match ev.button_state() {
                        ButtonState::Pressed => 1,
                        ButtonState::Released => 0,
                    },
                )
            ]
        ),
        PointerEvent::Axis(_) => None,
    }
}
