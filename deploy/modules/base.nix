{ pkgs, ... }:

{
  nix.settings = {
    experimental-features = [ "nix-command" "flakes" ];
    trusted-users = [ "root" "deploy" ];
    auto-optimise-store = true;
  };

  nix.gc = {
    automatic = true;
    dates = "weekly";
    options = "--delete-older-than 14d";
  };

  time.timeZone = "UTC";
  i18n.defaultLocale = "en_US.UTF-8";

  networking.firewall = {
    enable = true;
    allowedTCPPorts = [ 22 80 443 8080 ];
  };

  environment.systemPackages = with pkgs; [
    curl
    wget
    git
    htop
    jq
    sqlite
    clickhouse-cli
  ];
}
