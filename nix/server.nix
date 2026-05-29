{ yt-dlp-web-ui-frontend, buildGo124Module, lib, makeWrapper, yt-dlp, ... }:
let
  fs = lib.fileset;
  common = import ./common.nix { inherit lib; };
in
buildGo124Module {
  pname = "yt-dlp-web-ui";
  inherit (common) version;
  src = fs.toSource rec {
    root = ../.;
    fileset = fs.difference root (fs.unions [
      ### LIST OF FILES TO IGNORE ###
      # frontend (this is included by the frontend.nix drv instead)
      ../frontend
      # documentation
      ../examples
      # docker
      ../Dockerfile
      ../docker-compose.yml
      # nix
      ./devShell.nix
      ../.envrc
      ./tests
      # make
      ../Makefile # this derivation does not use the project Makefile
      # repo commons
      ../.github
      ../README.md
      ../LICENSE
      ../.gitignore
      ../.vscode
    ]);
  };

  # https://github.com/golang/go/issues/44507
  preBuild = ''
    cp -r ${yt-dlp-web-ui-frontend} frontend
  '';

  nativeBuildInputs = [ makeWrapper ];

  postInstall = ''
    wrapProgram $out/bin/yt-dlp-web-ui \
      --prefix PATH : ${lib.makeBinPath [ yt-dlp ]}
  '';

  vendorHash = "sha256-namJ99iUb70HtcfUtUAGsx3dmmC5fhuOBHGGPqYGcaA=";

  meta = common.meta // {
    mainProgram = "yt-dlp-web-ui";
  };
}
