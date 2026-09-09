{ ... }:

let
  keys = import ../keys.nix;
in
{
  services.openssh = {
    enable = true;
    settings = {
      PasswordAuthentication = false;
      KbdInteractiveAuthentication = false;
      PermitRootLogin = "prohibit-password";
      X11Forwarding = false;
      MaxAuthTries = 3;
    };
  };

  users.users.root.openssh.authorizedKeys.keys = builtins.attrValues keys;

  users.users.deploy = {
    isNormalUser = true;
    extraGroups = [ "wheel" ];
    openssh.authorizedKeys.keys = builtins.attrValues keys;
  };

  security.sudo.wheelNeedsPassword = false;
}
