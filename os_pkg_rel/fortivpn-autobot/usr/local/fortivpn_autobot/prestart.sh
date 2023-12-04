#!/bin/bash
set -e

if [[ $(id -u) -ne 0 ]]
then
  echo "Must be run under root."
  exit 5
fi

modprobe tun
mkdir -p /etc/fortivpn_autobot
