// NOTE: x#[xxx(yyy)] を アトリビュートというらしい。C# (EMCA-335) 由来らしい。多分 Java のアノテーションのようなもの。
// NOTE: derive アトリビュートは、trait のデフォルト実装をその対象に適用するとかそういうものっぽい。
// NOTE: see https://doc.rust-lang.org/reference/attributes/derive.html
#[derive(Debug)]
struct Structure(i32);

#[derive(Debug)]
struct Deep(Structure);

#[derive(Debug)]
struct Person<'a> {
    name: &'a str,
    age: u8
}

fn main() {
    // Debug trait が実装されていれば {:?} で出力できる。
    println!("{:?}", Structure(3));
    println!("{:?}", Deep(Structure(3)));

    let peter = Person { name:"Peter", age:27 };
    println!("{:?}", peter);
    // {:#?} なら pretty print
    println!("{:#?}", peter);
}
    
