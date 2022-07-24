use std::env;
use std::fs;
use std::path::Path;

fn hello() {
    // target/<target>/<package>/out/hello.rs にファイルを生成して出力する
    // `include!(concat!(env!("OUT_DIR"), "/hello.rs"))` ってするとインクルードできる
    let out_dir = env::var_os("OUT_DIR").unwrap();
    let dest_path = Path::new(&out_dir).join("hello.rs");
    fs::write(
        &dest_path,
        "pub fn message() -> &'static str {
            \"Hello, World!\"
        }
        ",
    )
    .unwrap();
    // build.rs が変更されたらビルド時に build.rs を再実行する
    println!("cargo:rerun-if-changed=build.rs");
}

fn libz() {
    // libz-sys を利用した C のファイルをコンパイルする
    let mut cfg = cc::Build::new();
    cfg.file("src/zuser.c");
    if let Some(include) = std::env::var_os("DEP_Z_INCLUDE") {
        cfg.include(include);
    }
    cfg.compile("zuser");
    // src/zuser.c が変更されたら build.rs を再実行する
    println!("cargo:rerun-if-changed=src/zuser.c");
}

fn cfg_flag() {
    // cfg(moo), cfg(cow = "moo") を設定する
    println!("cargo:rustc-cfg=moo");
    println!("cargo:rustc-cfg=cow=\"moo\"");
}

fn pkgconfig() {
    // システムにインストールされている libinput のバージョンが 1.19 以上なら cfg(libinput_1_19) を設定する
    // feature flag ではない事に注意 (設定不可)
    if pkg_config::Config::new()
        .atleast_version("1.19")
        .probe("libinput")
        .is_ok()
    {
        println!("cargo:rustc-cfg=libinput_1_19");
    }
}

fn main() {
    hello();
    libz();
    cfg_flag();
    pkgconfig();
}
