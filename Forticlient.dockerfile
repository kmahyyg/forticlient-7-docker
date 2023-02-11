FROM debian:stable
# Author Notes
LABEL ARCH="amd64"
LABEL MAINTAINER="kmahyyg <16604643+kmahyyg@users.noreply.github.com>"
LABEL PRIVILEGE_REQUEST="--device=/dev/net/tun --cap-add=NET_ADMIN --security-opt seccomp=unconfined"
LABEL ENV_REQUEST="FORTIVPN_PASSWD,ALLOW_INSECURE,FORTIVPN_SRV,FORTIVPN_USR"
LABEL DEV_DEPENDENCIES="wget procps tree gosu openssl sudo vim xxd"
LABEL VERSION="7.0.7.0246-VPN_ONLY-deb"
# Environment
ENV FORTIVPN_CLI="/opt/forticlient/vpn"
ENV DEBIAN_FRONTEND=noninteractive
ENV S6_KEEP_ENV=1
# Do not modify
WORKDIR /tmp
# Installation of Software
ADD fortirun.expect /
ADD resolv.conf /etc/resolv.conf
RUN apt update -y && \
    apt install curl gnupg2 gzip xz-utils expect ca-certificates iproute2 -y && \
    curl -L -O https://github.com/just-containers/s6-overlay/releases/download/v3.1.3.0/s6-overlay-noarch.tar.xz && \
    tar -C / -Jxpf /tmp/s6-overlay-noarch.tar.xz && rm /tmp/s6-overlay-noarch.tar.xz && \
    curl -L -O https://github.com/just-containers/s6-overlay/releases/download/v3.1.3.0/s6-overlay-x86_64.tar.xz && \
    tar -C / -Jxpf /tmp/s6-overlay-x86_64.tar.xz && rm /tmp/s6-overlay-x86_64.tar.xz && \
    curl -L -o - https://repo.fortinet.com/repo/7.0/debian/DEB-GPG-KEY | sudo apt-key add - && \
    curl -o /tmp/vpnagent.deb -L https://links.fortinet.com/forticlient/deb/vpnagent && \
    apt install -y /tmp/vpnagent.deb && \
    useradd -u 1000 -U -m fortiuser && \
    rm -rf /var/cache/apt/* /tmp/vpnagent.deb && \
    curl -L -O https://github.com/ginuerzh/gost/releases/download/v2.11.5/gost-linux-amd64-2.11.5.gz && \
    gunzip /tmp/gost-linux-amd64-2.11.5.gz && \
    mv /tmp/gost-linux-amd64-2.11.5 /usr/bin/gost && \
    chmod +x /usr/bin/gost /fortirun.expect
# Now go ahead, add service script
ADD s6-rc.d/fortivpn /etc/s6-overlay/s6-rc.d/fortivpn
ADD s6-rc.d/gost /etc/s6-overlay/s6-rc.d/gost
ADD s6-rc.d/user/contents.d/gost /etc/s6-overlay/s6-rc.d/user/contents.d/gost
# Now run
EXPOSE 10800
ENTRYPOINT ["/init"]
