packages: { config, lib, pkgs, ... }:
let
  cfg = config.services.yt-dlp-web-ui;
  inherit (pkgs.stdenv.hostPlatform) system;
  pkg = packages.${system}.default;
in
{
  /*
    Some notes on the module design:
    - Usually, you don't map out all of the options like this in attrsets,
      but due to the software's nonstandard "config file overrides CLI" behavior,
      we don't want to expose a config file catchall, and as such don't use '-conf'.

    - Notably, '-driver' is missing as a configuration option. 
      This should instead be customized with idiomatic Nix, overriding 'cfg.package' with
      the desired yt-dlp package.

    - The systemd service has been sandboxed as much as possible. This restricts configuration of 
      data and logs dir. If you really need a custom data and logs dir, use BindPaths (man systemd.exec)
  */
  options.services.yt-dlp-web-ui = {
    enable = lib.mkEnableOption "yt-dlp-web-ui";
    package = lib.mkOption {
      type = lib.types.package;
      default = pkg;
      defaultText = lib.literalMD "`packages.default` from the yt-dlp-web-ui flake.";
      description = ''
        The yt-dlp-web-ui package to use.
      '';
    };

    user = lib.mkOption {
      type = lib.types.str;
      default = "yt-dlp-web-ui";
      description = lib.mdDoc ''
        User under which yt-dlp-web-ui runs.
      '';
    };

    group = lib.mkOption {
      type = lib.types.str;
      default = "yt-dlp-web-ui";
      description = lib.mdDoc ''
        Group under which yt-dlp-web-ui runs.
      '';
    };

    openFirewall = lib.mkOption {
      type = lib.types.bool;
      default = false;
      description = lib.mdDoc ''
        Whether to open the TCP port in the firewall.
      '';
    };

    host = lib.mkOption {
      default = "0.0.0.0";
      type = lib.types.str;
      description = lib.mdDoc ''
        Host where yt-dlp-web-ui will listen at.
      '';
    };

    port = lib.mkOption {
      default = 3033;
      type = lib.types.port;
      description = lib.mdDoc ''
        Port where yt-dlp-web-ui will listen at.
      '';
    };

    downloadDir = lib.mkOption {
      type = lib.types.str;
      description = lib.mdDoc ''
        The directory where yt-dlp-web-ui stores downloads.
      '';
    };

    queueSize = lib.mkOption {
      default = 2;
      type = lib.types.ints.unsigned; # >= 0
      description = lib.mdDoc ''
        Queue size (concurrent downloads).
      '';
    };

    logging = lib.mkEnableOption "logging";

    rpcAuth = lib.mkOption {
      description = lib.mdDoc ''
        RPC Authentication settings.
      '';
      default = { };
      type = lib.types.submodule {
        options = {
          enable = lib.mkEnableOption "RPC authentication";
          user = lib.mkOption {
            type = lib.types.str;
            description = lib.mdDoc ''
              Username required for auth.
            '';
          };
          passwordFile = lib.mkOption {
            type = with lib.types; nullOr str;
            default = null;
            description = lib.mdDoc ''
              Path to the file containing the password required for auth.
            '';
          };
          insecurePasswordText = lib.mkOption {
            type = with lib.types; nullOr str;
            default = null;
            description = lib.mdDoc ''
              Raw password required for auth.

              It's strongly recommended to use 'passwordFile' instead of this option.
              
              **Don't use this option unless you know what you're doing!**. 
              It writes the password to the world-readable Nix store, which is a big security risk.
              More info: https://wiki.nixos.org/wiki/Comparison_of_secret_managing_schemes 
            '';
          };
        };
      };
    };

  };
  config = lib.mkIf cfg.enable {
    assertions = [
      (lib.mkIf cfg.rpcAuth.enable {
        assertion = lib.xor (cfg.rpcAuth.passwordFile == null) (cfg.rpcAuth.insecurePasswordText == null);
        message = ''
          RPC Auth is enabled for yt-dlp-web-ui! Exactly one RPC auth password source must be set!

          Tip: You should set 'services.yt-dlp-web-ui.rpcAuth.passwordfile'!
        '';
      })
    ];

    networking.firewall.allowedTCPPorts = lib.mkIf cfg.openFirewall [ cfg.port ];

    users.users = lib.mkIf (cfg.user == "yt-dlp-web-ui") {
      yt-dlp-web-ui = {
        inherit (cfg) group;
        isSystemUser = true;
      };
    };

    users.groups = lib.mkIf (cfg.group == "yt-dlp-web-ui") { yt-dlp-web-ui = { }; };

    systemd.services.yt-dlp-web-ui = {
      description = "yt-dlp-web-ui system service";
      after = [ "network.target" ];
      path = [ cfg.package pkgs.tree ];
      wantedBy = [ "multi-user.target" ];
      serviceConfig =
        rec {
          ExecStart =
            let
              password =
                if cfg.rpcAuth.passwordFile == null
                then cfg.rpcAuth.insecurePasswordText
                else "$(cat ${cfg.rpcAuth.passwordFile})";
              args = [
                "-host ${cfg.host}"
                "-port ${builtins.toString cfg.port}"
                ''-out "${cfg.downloadDir}"''
                "-qs ${builtins.toString cfg.queueSize}"
              ] ++ (lib.optionals cfg.logging [
                "-fl"
                ''-lf "/var/log/${LogsDirectory}/yt-dlp-web-ui.log"''
              ]) ++ (lib.optionals cfg.rpcAuth.enable [
                "-auth"
                "-user ${cfg.rpcAuth.user}"
                "-pass ${password}"
              ]);
            in
            "${lib.getExe cfg.package} ${lib.concatStringsSep " " args}";
          User = cfg.user;
          Group = cfg.group;
          ProtectSystem = "strict";
          ProtectHome = "read-only";
          StateDirectory = "yt-dlp-web-ui";
          WorkingDirectory = "/var/lib/${StateDirectory}"; # equivalent to the dir above
          LogsDirectory = "yt-dlp-web-ui";
          ReadWritePaths = [
            cfg.downloadDir
          ];
          BindReadOnlyPaths = [
            builtins.storeDir
            # required for youtube DNS lookup
            "${config.environment.etc."ssl/certs/ca-certificates.crt".source}:/etc/ssl/certs/ca-certificates.crt"
          ] ++ lib.optionals (cfg.rpcAuth.enable && cfg.rpcAuth.passwordFile != null) [
            cfg.rpcAuth.passwordFile
          ];
          CapabilityBoundingSet = "";
          RestrictAddressFamilies = [ "AF_UNIX" "AF_INET" "AF_INET6" ];
          RestrictNamespaces = true;
          PrivateDevices = true;
          PrivateUsers = true;
          ProtectClock = true;
          ProtectControlGroups = true;
          ProtectKernelLogs = true;
          ProtectKernelModules = true;
          ProtectKernelTunables = true;
          SystemCallArchitectures = "native";
          SystemCallFilter = [ "@system-service" "~@privileged" ];
          RestrictRealtime = true;
          LockPersonality = true;
          MemoryDenyWriteExecute = true;
          ProtectHostname = true;
        };
    };
  };
}
