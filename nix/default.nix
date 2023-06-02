{ lib, buildGoModule, ... }:
buildGoModule rec {
  name = "grafana-matrix-forwarder";
  version = "0.7.0";
  proxyVendor = true;

  vendorSha256 = "sha256-5P+TGUR6gSorDnsJ3xMxmu9spbH/ZeXcT563kgHFG8A=";

  src = ../src;
}
