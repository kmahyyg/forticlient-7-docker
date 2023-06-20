#!/bin/bash
set -xe
sudo tar -vxzf fortiauto_2023-06-20.tar.gz --owner=root --group=root -C /
sudo install -m 644 -o root -g root /usr/local/fortivpn_autobot/fortivpn-autobot.service /usr/lib/systemd/system/fortivpn-autobot.service
