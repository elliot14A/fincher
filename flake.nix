{
  description = "Fincher autonomous media delivery operations environment & deployment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";

    disko.url = "github:nix-community/disko";
    disko.inputs.nixpkgs.follows = "nixpkgs";

    agenix.url = "github:ryantm/agenix";
    agenix.inputs.nixpkgs.follows = "nixpkgs";

    deploy-rs.url = "github:serokell/deploy-rs";
    deploy-rs.inputs.nixpkgs.follows = "nixpkgs";

    nixos-anywhere.url = "github:nix-community/nixos-anywhere";
    nixos-anywhere.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs =
    inputs@{
      self,
      nixpkgs,
      flake-utils,
      disko,
      agenix,
      deploy-rs,
      nixos-anywhere,
      ...
    }:
    let
      supportedSystems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      system = "x86_64-linux";
    in
    flake-utils.lib.eachSystem supportedSystems (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            # Go Toolchain & Live Reload
            go
            gopls
            golangci-lint
            air

            # Data & SQLite CLIs
            sqlite
            lazysql
            clickhouse-cli

            # Frontend & Utilities
            bun
            curl
            jq
            just
            google-cloud-sdk

            # Deployment & Secrets
            deploy-rs.packages.${system}.default
            agenix.packages.${system}.default
            nixos-anywhere.packages.${system}.default
            age
            openssh
          ];

          shellHook = ''
            export CGO_ENABLED=1
          '';
        };

        packages = {
          default = pkgs.buildGoModule {
            pname = "fincher";
            version = "0.1.0";
            src = ./.;
            vendorHash = "sha256-Dd3DXgcOMo5ZlYa2P0ih/1j/9PLGVIQweHXFUATDhmA=";
            subPackages = [ "cmd/fincher" "cmd/seed" ];
            env.CGO_ENABLED = 1;
            ldflags = [ "-s" "-w" ];
          };
        };
      }
    )
    // {
      nixosConfigurations.fincher-prod = nixpkgs.lib.nixosSystem {
        inherit system;
        specialArgs = { inherit inputs; };
        modules = [
          ./deploy/hosts/fincher-prod
        ];
      };

      deploy.nodes.fincher-prod = {
        hostname = "34.93.123.20";
        sshUser = "deploy";
        profiles.system = {
          user = "root";
          path = deploy-rs.lib.${system}.activate.nixos self.nixosConfigurations.fincher-prod;
        };
      };
    };
}
