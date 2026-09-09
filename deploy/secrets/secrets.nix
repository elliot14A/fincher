let
  keys = import ../keys.nix;
in
{
  "fincher.env.age".publicKeys = [ keys.gcloud keys.host_fincher_prod ];
}
