{
  description = "A terrible web ui for yt-dlp. Designed to be self-hosted.";

  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    pre-commit-hooks-nix.url = "github:cachix/pre-commit-hooks.nix";
  };

  outputs = inputs@{ self, flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        inputs.pre-commit-hooks-nix.flakeModule
      ];
      systems = [
        "x86_64-linux"
      ];
      perSystem = { config, self', pkgs, ... }: {

        packages = {
          yt-dlp-web-ui-frontend = pkgs.callPackage ./nix/frontend.nix { };
          default = pkgs.callPackage ./nix/server.nix {
            inherit (self'.packages) yt-dlp-web-ui-frontend;
          };
        };

        checks = import ./nix/tests { inherit self pkgs; };

        pre-commit = {
          check.enable = true;
          settings = {
            hooks = {
              ${self'.formatter.pname}.enable = true;
              deadnix.enable = true;
              nil.enable = true;
              statix.enable = true;
            };
          };
        };

        devShells.default = pkgs.callPackage ./nix/devShell.nix {
          inputsFrom = [ config.pre-commit.devShell ];
        };

        formatter = pkgs.nixpkgs-fmt;
      };
      flake = {
        nixosModules.default = import ./nix/module.nix self.packages;
      };
    };
}
