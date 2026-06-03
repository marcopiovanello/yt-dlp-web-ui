{ inputsFrom ? [ ], mkShell, yt-dlp, nodejs, go }:
mkShell {
  inherit inputsFrom;
  packages = [
    yt-dlp
    nodejs
    go
  ];
}
