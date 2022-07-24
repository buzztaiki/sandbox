use input::event::keyboard::KeyboardEventTrait;
use input::event::pointer::Axis;
use input::event::{DeviceEvent, EventTrait, KeyboardEvent, PointerEvent};
use input::{DeviceCapability, Event, Libinput};

fn main() {
    let mut input = Libinput::new_with_udev(libinput::Interface);
    input.udev_assign_seat("seat0").unwrap();
    loop {
        input.dispatch().unwrap();
        for event in &mut input {
            let device = event.device();
            let message = match event {
                Event::Keyboard(KeyboardEvent::Key(ev)) => {
                    format!("key: {:?}, {:?}", ev.key(), ev.key_state())
                }
                Event::Pointer(PointerEvent::Motion(ev)) => {
                    format!(
                        "pointer motion: dx={:?}/{:?}, dy={:?}/{:?}",
                        ev.dx_unaccelerated(),
                        ev.dx(),
                        ev.dy_unaccelerated(),
                        ev.dy(),
                    )
                }
                Event::Pointer(PointerEvent::MotionAbsolute(ev)) => format!(
                    "pointer motion absolute: {:?}, {:?}",
                    ev.absolute_x(),
                    ev.absolute_y()
                ),
                Event::Pointer(PointerEvent::Button(ev)) => {
                    format!("pointer button: {:?}, {:?}", ev.button(), ev.button_state())
                }
                Event::Pointer(PointerEvent::Axis(ev)) => {
                    let f = |axis| {
                        if ev.has_axis(axis) {
                            format!(
                                "{:?}/{:?}",
                                ev.axis_value(axis),
                                ev.axis_value_discrete(axis)
                            )
                        } else {
                            "<none>".to_string()
                        }
                    };

                    format!(
                        "pointer axis: {:?}, h={}, v={}",
                        ev.axis_source(),
                        f(Axis::Horizontal),
                        f(Axis::Vertical),
                    )
                }
                Event::Device(DeviceEvent::Added(ev)) => {
                    let dev = ev.device();
                    let all_caps = vec![
                        DeviceCapability::Keyboard,
                        DeviceCapability::Pointer,
                        DeviceCapability::Touch,
                        DeviceCapability::TabletTool,
                        DeviceCapability::TabletPad,
                        DeviceCapability::Gesture,
                        DeviceCapability::Switch,
                    ];
                    let caps = all_caps.iter().filter(|x| dev.has_capability(**x));
                    format!(
                        "device added: name={:?}, sysname={:?}, vendor={:?}, product={:?}, capabilities={}",
                        dev.name(),
                        dev.sysname(),
                        dev.id_vendor(),
                        dev.id_product(),
                        caps.map(|x| format!("{:?}", x)).collect::<Vec<_>>().join(","),
                    )
                }
                Event::Device(DeviceEvent::Removed(ev)) => {
                    format!("device added: name={:?}", ev.device().name())
                }
                _ => format!("event: {:?}", event),
            };
            println!("{:<80}: {}", device.name(), message);
        }
    }
}
