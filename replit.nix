{ pkgs }: {
  deps = [
    pkgs.rustc
    pkgs.cargo
    pkgs.openssl
    pkgs.pkg-config
    pkgs.postgresql
    pkgs.nodejs_20
    pkgs.nodePackages.npm
  ];
}
