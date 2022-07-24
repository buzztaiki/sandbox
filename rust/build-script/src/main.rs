include!(concat!(env!("OUT_DIR"), "/hello.rs"));

fn main() {
    println!("{}", message());

    if cfg!(moo) {
        println!("moooooo");
    }
    if cfg!(cow = "moo") {
        println!("moooooo!!!");
    }
    if cfg!(libinput_1_19) {
        println!("has libinput 1.19");
    }
}
