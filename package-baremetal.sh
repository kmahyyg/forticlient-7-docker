#!/bin/bash
set -ex

TMPDIR=$(mktemp -d)
CURDIR=$(pwd)
CURDATE=$(date -Idate)

mkdir -p ${TMPDIR}/etc/fortivpn_autobot
mkdir -p ${TMPDIR}/usr/local/fortivpn_autobot
cp -af ./bin/go-fortivpn-daemon ${TMPDIR}/usr/local/fortivpn_autobot
cp -af ./bareMetalAssets/default ${TMPDIR}/etc/fortivpn_autobot
cp -af ./bareMetalAssets/prestart.sh ${TMPDIR}/usr/local/fortivpn_autobot
cp -af ./bareMetalAssets/fortivpn-autobot.service ${TMPDIR}/usr/local/fortivpn_autobot
cd ${TMPDIR}
tar czvf ${CURDIR}/bin/fortiauto_${CURDATE}.tar.gz .
cd ${CURDIR}
rm -rf ${TMPDIR}
