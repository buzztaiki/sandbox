{ pkgs ? import <nixpkgs> {} }:
  pkgs.mkShell {
    nativeBuildInputs = [
      pkgs.llvmPackages_10.libllvm
      pkgs.llvmPackages_10.clang
      pkgs.python3
      pkgs.poetry
    ];
}
