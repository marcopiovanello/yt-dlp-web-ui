{ self, pkgs }: {
  testServiceStarts = pkgs.testers.runNixOSTest (_: {
    name = "service-starts";
    nodes = {
      machine = _: {
        imports = [
          self.nixosModules.default
        ];

        services.yt-dlp-web-ui = {
          enable = true;
          downloadDir = "/var/lib/yt-dlp-web-ui";
        };
      };
    };
    testScript = ''
      machine.wait_for_unit("yt-dlp-web-ui")
    '';
  });
}
