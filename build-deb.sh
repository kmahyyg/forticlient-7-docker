#!/bin/bash

chmod +x ./fortivpn-autobot/DEBIAN/postinst
chmod +x ./fortivpn-autobot/usr/local/fortivpn_autobot/*
fakeroot dpkg-deb -b fortivpn-autobot
