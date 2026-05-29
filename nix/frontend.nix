{ lib
, stdenv
, nodejs
, pnpm
}:
let common = import ./common.nix { inherit lib; }; in
stdenv.mkDerivation (finalAttrs: {
  pname = "yt-dlp-web-ui-frontend";

  inherit (common) version;

  src = lib.fileset.toSource {
    root = ../frontend;
    fileset = ../frontend;
  };

  buildPhase = ''
    npm run build
  '';

  installPhase = ''
    mkdir -p $out/dist
    cp -r dist/* $out/dist
  '';

  nativeBuildInputs = [
    nodejs
    pnpm.configHook
  ];

  pnpmDeps = pnpm.fetchDeps {
    inherit (finalAttrs) pname version src;
    fetcherVersion = 2;
    hash = "sha256-0UcNHvFoI3yENAq9fuZEfK6Pxo09FHI+Yzkf2BK/BUY=";
  };

  inherit (common) meta;
})
