{
  description = "kfilt";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    let
      kfilt = pkgs:
        pkgs.buildGo120Module rec {
          name = "kfilt";
          version = self.shortRev or "dirty";
          src = ./.;
          # this needs to be changed any time there is a change in go.mod
          # dependencies
          vendorSha256 = "sha256-c77CzpE9cPyobt87uO0QlkKD+xC/tM7wOy4orM62tnI=";
          nativeBuildInputs = [ ];
          CGO_ENABLED = 0;
          doCheck = false;
          ldflags = [
            "-s -w -X github.com/ryane/kfilt/cmd.Version=${version} -X github.com/ryane/kfilt/cmd.GitCommit=${version}"
          ];
          excludedPackages = [ "plugin/kustomize" ];
        };
      flakeForSystem = nixpkgs: system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
          kf = kfilt pkgs;
        in {
          packages = { kfilt = kf; };
          devShell = pkgs.mkShell { packages = with pkgs; [ curl ]; };
        };
    in flake-utils.lib.eachDefaultSystem
    (system: flakeForSystem nixpkgs system);
}
