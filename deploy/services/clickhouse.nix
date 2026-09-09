{ ... }:

{
  services.clickhouse = {
    enable = true;
  };

  systemd.tmpfiles.rules = [
    "d /var/lib/clickhouse 0700 clickhouse clickhouse -"
  ];
}
