{ ... }:

{
  virtualisation.oci-containers.backend = "podman";
  virtualisation.podman.enable = true;

  virtualisation.oci-containers.containers.mcp-clickhouse = {
    image = "ghcr.io/clickhouse/mcp-clickhouse:latest";
    environment = {
      CLICKHOUSE_HOST = "127.0.0.1";
      CLICKHOUSE_PORT = "8123";
      CLICKHOUSE_USER = "default";
      CLICKHOUSE_DATABASE = "fincher";
      CLICKHOUSE_SECURE = "false";
      CLICKHOUSE_MCP_SERVER_TRANSPORT = "http";
      CLICKHOUSE_MCP_BIND_HOST = "127.0.0.1";
      CLICKHOUSE_MCP_BIND_PORT = "8000";
      CLICKHOUSE_MCP_ALLOWED_HOSTS = "127.0.0.1:8000,localhost:8000,127.0.0.1,localhost,*";
      CLICKHOUSE_MCP_AUTH_DISABLED = "true";
    };
    extraOptions = [ "--network=host" ];
  };
}
