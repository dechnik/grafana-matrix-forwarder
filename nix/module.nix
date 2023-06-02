self: { config, lib, pkgs, ... }:
let
  cfg = config.services.grafana-matrix-forwarder;
in {
  options.services.grafana-matrix-forwarder = {
    enable = lib.mkEnableOption "grafana-matrix-forwarder";
    package = lib.mkOption {
      type = lib.types.package;
      default = self.packages.${pkgs.system}.default;
      description = "The package implementing grafana-matrix-forwarder";
    };
    matrixAuthFile = lib.mkOption {
      type = lib.types.path;
      description = "File path containing GMF_MATRIX_USER,GMF_MATRIX_PASSWORD,GMF_MATRIX_HOMESERVER envs";
      default = null;
    };
    serverHost = lib.mkOption {
      type = lib.types.str;
      default = "0.0.0.0";
      description = "Host address the server connects to.";
    };
    serverPort = lib.mkOption {
      type = lib.types.int;
      default = 6000;
      description = "Port to run the webserver on.";
    };
    openFirewall = lib.mkOption {
      type = lib.types.bool;
      default = false;
      description = "Whether to open port in the firewall for the server.";
    };
  };
  config = lib.mkIf cfg.enable {
    # environment.systemPackages = [ pkgs.lldap ];
    networking.firewall =
      lib.mkIf cfg.openFirewall { allowedTCPPorts = [ cfg.serverPort ]; };

    systemd.services.grafana-matrix-forwarder = {
      wantedBy = [ "multi-user.target" ];
      after = [ "network.target" ];
      serviceConfig = {
        Environment = [
          "GMF_SERVER_HOST=${cfg.serverHost}"
          "GMF_SERVER_PORT=${toString cfg.serverPort}"
        ];
        EnvironmentFile = cfg.matrixAuthFile;
        ExecStart = "${cfg.package}/bin/grafana-matrix-forwarder -env";
        Restart = "on-failure";
        Type = "simple";
        User = "gmf";
        Group = "gmf";
      };
    };

    users.users.gmf = {
      group = "gmf";
      createHome = false;
      isSystemUser = true;
    };
    users.groups.gmf = { };
  };
}
