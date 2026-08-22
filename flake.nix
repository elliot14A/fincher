{
  description = "Fincher autonomous media delivery operations environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    { self, nixpkgs, ... }:
    let
      supportedSystems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      forEachSupportedSystem =
        f:
        nixpkgs.lib.genAttrs supportedSystems (
          system:
          f {
            pkgs = import nixpkgs { inherit system; };
          }
        );
    in
    {
      devShells = forEachSupportedSystem (
        { pkgs }:
        {
          default = pkgs.mkShell {
            packages = with pkgs; [
              # Go Toolchain
              go
              gopls
              golangci-lint

              # Data & SQLite CLIs
              sqlite
              lazysql
              clickhouse-cli

              # Frontend & Utilities
              bun
              curl
              jq
            ];

            shellHook = ''
              export CGO_ENABLED=1
            '';
          };
        }
      );
    };
}
