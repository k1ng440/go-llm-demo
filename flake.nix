{
  description = "Go LLM Demo";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        inherit (pkgs) go;
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
            golangci-lint
            delve
          ];

          shellHook = ''
            echo "=== Go LLM Demo - Development Shell ==="
            echo "Go version: $(go version)"
          '';
        };
      }
    );
}
