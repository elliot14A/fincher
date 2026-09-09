{ ... }:

{
  services.clickhouse = {
    enable = true;
    extraUsersConfig = ''
      <clickhouse>
        <users>
          <default>
            <password remove="1" />
            <password_sha256_hex>35ccf4cbcbd8c466c1132b6b211cf5f7b5aa59f2c50a370bd19a688429d7eba8</password_sha256_hex>
            <profile>default</profile>
            <quota>default</quota>
            <networks>
              <ip>::/0</ip>
            </networks>
          </default>
        </users>
      </clickhouse>
    '';
  };

  systemd.tmpfiles.rules = [
    "d /var/lib/clickhouse 0700 clickhouse clickhouse -"
  ];
}
