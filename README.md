# Forticlient 7 VPN Only Docker

## Usage

Copy and Modify `.env.example` to `.env` then set the following environment variable before you start: 

|          Env Var         |  Spec  |
|:------------------------:|:------:|
|`FORTIVPN_PASSWD`| VPN Password |
|`ALLOW_INSECURE`| Auto-Answer Allow Insecure CA|
|`FORTIVPN_SRV`| Server Address (host:port) |
|`FORTIVPN_USR`| VPN Username |

When starting your container at the very beginning, add `--env-file .env --device=/dev/net/tun --cap-add=NET_ADMIN` to prevent further issue.

No external volume is required. This container only works for AMD64, since FortiNet does NOT offer other architecture.

SOCKS5 proxy is exposed on 10080 port, without any authentication. Please do NOT expose it on public internet.

After it works, just go access network via socks5 10080.

Note: Currently, it just meet my personal needs, you should modify `fortirun.expect` if you need certificate authentication.

2FA authentication may not be supported. Possibly, you could read a file contains current TOTP and modify `fortirun.expect` to read it.

## Credit

- FortiClient made possible by Fortinet, URL: https://www.fortinet.com/support/product-downloads#vpn
- S6-Overlay from https://github.com/just-containers/s6-overlay
- GOSU from https://github.com/tianon/gosu 
- GOST from https://github.com/ginuerzh/gost

## License

### FortiClient

Copyright © 2023 Fortinet, Inc. All rights reserved. Fortinet®, FortiGate®, FortiCare® and FortiGuard®, and certain other marks are registered trademarks of Fortinet, Inc., and other Fortinet names herein may also be registered and/or common law trademarks of Fortinet. All other product or company names may be trademarks of their respective owners. Performance and other metrics contained herein were attained in internal lab tests under ideal conditions, and actual performance and other results may vary. Network variables, different network environments and other conditions may affect performance results. Nothing herein represents any binding commitment by Fortinet, and Fortinet disclaims all warranties, whether express or implied, except to the extent Fortinet enters a binding written contract, signed by Fortinet’s General Counsel, with a purchase that expressly warrants that the identified product will perform according to certain expressly-identified performance metrics and, in such event, only the specific performance metrics expressly identified in such binding written contract shall be binding on Fortinet. For absolute clarity, any such warranty will be limited to performance in the same ideal conditions as in Fortinet’s internal lab tests. Fortinet disclaims in full any covenants, representations, and guarantees pursuant hereto, whether express or implied. Fortinet reserves the right to change, modify, transfer, or otherwise revise this publication without notice, and the most current version of the publication shall be applicable.

### Docker

Docker, Inc. (“Docker”) trademarks, service marks, logos and designs, as well as other works of authorship that are eligible for copyright protection (collectively termed “Marks”) are valuable assets that Docker needs to protect. Docker does not permit all uses of Docker Marks, but Docker has posted these Trademark Usage Guidelines (“Guidelines”) to assist you in properly using our Marks in specific cases that we permit. The strength of our Marks depends, in part, upon consistent and appropriate use. We ask that you properly use and credit our Marks in accordance with these Guidelines. We reserve the right to change these Guidelines at any time and solely at our discretion.

### Docker Resources

 forticlient7-docker
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

