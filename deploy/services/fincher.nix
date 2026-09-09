{ config, pkgs, inputs, ... }:

let
  fincherPkg = inputs.self.packages.${pkgs.system}.default;
in
{
  users.users.fincher = {
    isSystemUser = true;
    group = "fincher";
  };
  users.groups.fincher = {};

  age.secrets."fincher.env" = {
    file = ../secrets/fincher.env.age;
    owner = "fincher";
    group = "fincher";
    mode = "0400";
  };

  systemd.services.fincher = {
    description = "Fincher Autonomous Media Delivery Operations Platform";
    wantedBy = [ "multi-user.target" ];
    after = [ "network-online.target" "clickhouse.service" "agenix.service" ];
    wants = [ "network-online.target" "clickhouse.service" ];

    serviceConfig = {
      User = "fincher";
      Group = "fincher";
      ExecStart = "${fincherPkg}/bin/fincher";
      Restart = "always";
      RestartSec = 5;
      StateDirectory = "fincher";
      WorkingDirectory = "/var/lib/fincher";
      EnvironmentFile = config.age.secrets."fincher.env".path;
      LimitNOFILE = 65536;
    };

    environment = {
      FINCHER_ENV = "production";
      FINCHER_PORT = "8080";
      FINCHER_TURSO_URL = "/var/lib/fincher/fincher.db";
      FINCHER_MCP_URL = "http://127.0.0.1:8000/mcp";
      FINCHER_CLICKHOUSE_DSN = "127.0.0.1:9000";
      FINCHER_GEMINI_MODEL = "gemini-2.5-flash";
      FINCHER_GEMINI_OPTIONS = "location=asia-south1";
    };
  };
}
