{ pkgs, ... }:

{
  languages.go.enable = true;

  packages = [
    pkgs.golangci-lint
    pkgs.poppler-utils
  ];
}
