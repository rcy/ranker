with (import <nixpkgs> {});
mkShell {
  buildInputs = [
    go
    golint
    gopls
  ];
}
