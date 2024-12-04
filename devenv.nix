{ pkgs, lib, config, inputs, ... }:

{
  packages = with pkgs; [
    rsync
    just
  ];

  languages.go.enable = true;

  processes.rsyncd.exec = ''
    rsync --daemon --no-detach --verbose --config=rsyncd.conf --address=127.0.0.1 --port=1873
  '';
}
