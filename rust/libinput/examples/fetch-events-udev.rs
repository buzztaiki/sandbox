use input::Libinput;

fn main() {
    let mut input = Libinput::new_with_udev(libinput::Interface);
    input.udev_assign_seat("seat0").unwrap();
    loop {
        input.dispatch().unwrap();
        for event in &mut input {
            println!("Got event: {:?}", event);
        }
    }
}
