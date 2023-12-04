# Forticlient 7 VPN Only - Podman

**Note: Due to limitation of Forticlient implementation and Docker hardcoded bind-mount behaviour, this image should only be run using Podman.**

**Note: This image will detect container internal environment, and will refuse to execute if NOT podman.**

**Note: Forticlient requires a TTY to work, so getting update from stdout may miss data. Be careful when you automate something.**

**Note: You should disable outside network manager for newly created virtual interface to prevent unexpected issue.**

If your network environment is unstable, please restart vpn container and check log, since Linux version is so crappy, single packet fault may lead to disconnect.

## Usage

Copy and Modify `.env.example` to `.env` then set the following environment variable before you start: 

|          Env Var         |  Spec  |
|:------------------------:|:------:|
|`FORTIVPN_PASSWD`| VPN Password |
|`ALLOW_INSECURE`| Auto-Answer Allow Insecure CA|
|`FORTIVPN_SRV`| Server Address (host:port) |
|`FORTIVPN_USR`| VPN Username |
|`GOST_USERNAME`| Gost Username |
|`GOST_PASSWD`| Gost Password |
|`FORTIVPN_TOTP_SECRET`| FortiVPN TOTP Secret (Manually extract from FortiToken) |

> To extract secret, check: https://jonstoler.me/blog/extracting-fortitoken-mobile-totp-secret (Make sure your time is synced)
> For legacy Android device (Android 7.1/8.1), the uuid decryption key is: `android_id (mapped to each app) + device hardware serial number`
> If none could be detected by app, the uuid decryption key is: `"23571113171923" repeat, until 16 chars`

When starting your container at the very beginning, add `--dns=none --env-file .env --device=/dev/net/tun --cap-add=NET_ADMIN` to prevent further issue.

No external volume is required. This container only works for AMD64, since FortiNet does NOT offer other architecture.

SOCKS5 proxy is exposed on 10800 port, with authentication and TLS encryption. Please use Gost v2 to connect and transfer it to plain text SOCKS5 on your local host to ensure your safety.

After it works, just go access network via socks5 10800, `-p 10800:10800`.

Note: Currently, it just meet my personal needs, you should modify `gosrc/answerBot/main.go` if you need certificate authentication.

2FA authentication may not be supported. Possibly, you could read a file contains current TOTP and modify `gosrc/answerBot/main.go` to read it and then auto-input TOTP code to satisfy your needs.

So finally your command should be like: `sudo podman run -d --env-file .env --device=/dev/net/tun --cap-add=NET_ADMIN --security-opt "seccomp=unconfined" -i -t -p 10800:10800 --dns=none ghcr.io/kmahyyg/fortivpn:7`

After that, download GOST v2 and run `gost -L socks5://127.0.0.1:1080 -F socks5+tls://${GOST_USERNAME}:${GOST_PASSWD}@<REMOTE MACHINE IP>:10800`, then set browser socks5 proxy to localhost:1080 or use smart traffic division, enjoy.

For your reference, a single container with no workload will cost 64M RAM and 0.1% CPU.

## Credit

- FortiClient made possible by Fortinet, URL: https://www.fortinet.com/support/product-downloads#vpn
- S6-Overlay from https://github.com/just-containers/s6-overlay
- GOSU from https://github.com/tianon/gosu 
- GOST v2 from https://github.com/ginuerzh/gost , v3 from https://github.com/go-gost/gost

## License

### FortiClient

Copyright © 2023 Fortinet, Inc. All rights reserved. Fortinet®, FortiGate®, FortiCare® and FortiGuard®, and certain other marks are registered trademarks of Fortinet, Inc., and other Fortinet names herein may also be registered and/or common law trademarks of Fortinet. All other product or company names may be trademarks of their respective owners. Performance and other metrics contained herein were attained in internal lab tests under ideal conditions, and actual performance and other results may vary. Network variables, different network environments and other conditions may affect performance results. Nothing herein represents any binding commitment by Fortinet, and Fortinet disclaims all warranties, whether express or implied, except to the extent Fortinet enters a binding written contract, signed by Fortinet’s General Counsel, with a purchase that expressly warrants that the identified product will perform according to certain expressly-identified performance metrics and, in such event, only the specific performance metrics expressly identified in such binding written contract shall be binding on Fortinet. For absolute clarity, any such warranty will be limited to performance in the same ideal conditions as in Fortinet’s internal lab tests. Fortinet disclaims in full any covenants, representations, and guarantees pursuant hereto, whether express or implied. Fortinet reserves the right to change, modify, transfer, or otherwise revise this publication without notice, and the most current version of the publication shall be applicable.

### Podman Resources

 forticlient-7-podman
 Copyright (C) 2023  Patmeow Limited
 
 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU Affero General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.
 
 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU Affero General Public License for more details.
 
 You should have received a copy of the GNU Affero General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.

