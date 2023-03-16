{ pkgs ? import <nixpkgs> {} }:
  pkgs.mkShell {
    buildInputs = [
      pkgs.python3
      pkgs.poetry
      pkgs.lttng-ust_2_12
      pkgs.curl
      pkgs.krb5
      pkgs.zlib
    ];

  shellHook = ''
    export LD_LIBRARY_PATH=${pkgs.lttng-ust_2_12.out}/lib:${pkgs.curl.out}/lib:${pkgs.stdenv.cc.cc.lib}/lib:${pkgs.krb5}/lib:${pkgs.zlib}/lib:$LD_LIBRARY_PATH
  '';
  }
