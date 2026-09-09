{ ... }:

{
  services.caddy = {
    enable = true;
    virtualHosts."fincher.elliot14a.work" = {
      extraConfig = ''
        reverse_proxy 127.0.0.1:8080
      '';
    };
  };
}
