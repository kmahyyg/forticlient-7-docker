#!/bin/bash
set -xe
sudo systemctl disable --now fortivpn-autobot
sudo rm -f /usr/lib/systemd/system/fortivpn-autobot.service
sudo rm -rf /usr/local/fortivpn_autobot
sudo rm -rf /etc/fortivpn_autobot
