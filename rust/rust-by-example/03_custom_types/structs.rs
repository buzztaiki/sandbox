#[derive(Debug)]
struct Person<'a> {
    // The 'a defines a lifetime
    name: &'a str,
    age: u8,
}

// A unit struct
struct Nil;

// A tuple struct
struct Pair(i32, f32);

// A struct with two fields
#[derive(Debug)]
struct Point {
    x: f32,
    y: f32,
}

// Structs can be reused as fields of another struct
#[derive(Debug)]
struct Rectangle {
    // A rectangle can be specified by where the top left and bottom right
    // corners are in space.
    top_left: Point,
    bottom_right: Point,
}

fn rect_area(rect: Rectangle) -> f32 {
    let Rectangle {
        top_left: Point { x: left, y: top},
        bottom_right: Point { x: right, y: bottom}
    } = rect;
    return (right - left) * (top - bottom);
}

fn square(point: Point, n: f32) -> Rectangle {
    return Rectangle {
        top_left: Point { y: point.y + n, ..point },
        bottom_right: Point { x: point.x + n, ..point}
    }
}

fn main() {
    // Create struct with field init shorthand
    let name = "Peter";
    let age = 27;
    let peter = Person { name, age };

    // Print debug struct
    println!("{:?}", peter);


    // Instantiate a `Point`
    let point = Point { x: 10.3, y: 0.4 };

    // Access the fields of the point
    println!("point coordinates: ({}, {})", point.x, point.y);

    // NOTE: これ便利だ。
    // Make a new point by using struct update syntax to use the fields of our other one
    let bottom_right = Point { x: 5.2, ..point };

    // `bottom_right.y` will be the same as `point.y` because we used that field from `point`
    println!("second point: ({}, {})", bottom_right.x, bottom_right.y);
    // NOTE: point をそのまま square に渡すと value moved here って言われて怒られる。
    // println!("square: {:?}", square(point, 10.0));
    // NOTE: .. で展開すると move が起きない模様。まだ rust わからん。
    println!("square: {:?}", square(Point {..point}, 10.0));

    // Destructure the point using a `let` binding
    let Point { x: top_edge, y: left_edge } = point;

    let rectangle = Rectangle {
        // struct instantiation is an expression too
        top_left: Point { x: left_edge, y: top_edge },
        bottom_right,
    };
    println!("rectangle: {:?}", rectangle);
    println!("rect_area: {:?}", rect_area(rectangle));

    // Instantiate a unit struct
    let _nil = Nil;

    // Instantiate a tuple struct
    let pair = Pair(1, 0.1);

    // Access the fields of a tuple struct
    println!("pair contains {:?} and {:?}", pair.0, pair.1);

    // Destructure a tuple struct
    let Pair(integer, decimal) = pair;

    println!("pair contains {:?} and {:?}", integer, decimal);
}
