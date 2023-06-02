{ lib, buildGoModule, ... }:
buildGoModule rec {
  name = "grafana-matrix-forwarder";
  version = "0.7.0";
  proxyVendor = true;

  vendorSha256 = "sha256-Pgx1UwVZcxFmrFeYeE3JhlmT5ZBbKdmsSRqZrtyj5dg=";

  src = ../src;
}
