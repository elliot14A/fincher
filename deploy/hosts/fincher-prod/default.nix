{ inputs, modulesPath, lib, ... }:

{
  imports = [
    (modulesPath + "/profiles/qemu-guest.nix")
    inputs.disko.nixosModules.disko
    inputs.agenix.nixosModules.default
    ../../modules/disko.nix
    ../../modules/base.nix
    ../../modules/ssh.nix
    ../../services/clickhouse.nix
    ../../services/mcp-clickhouse.nix
    ../../services/caddy.nix
    ../../services/fincher.nix
  ];

  networking.hostName = "fincher-prod";
  networking.useDHCP = lib.mkDefault true;

  boot.loader.systemd-boot.enable = false;
  boot.loader.grub = {
    enable = true;
    devices = [ "/dev/sda" ];
    efiSupport = true;
    efiInstallAsRemovable = true;
    configurationLimit = 10;
  };
  boot.loader.efi.canTouchEfiVariables = false;
  boot.growPartition = true;

  boot.kernelParams = [ "console=ttyS0,115200" "earlyprintk=ttyS0,115200" ];

  services.qemuGuest.enable = true;

  system.stateVersion = "24.11";
}
