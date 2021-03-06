#!/bin/bash

docker stop wireguard && docker rm wireguard

docker run -d \
  --name=wireguard \
  --cap-add=NET_ADMIN \
  --cap-add=SYS_MODULE \
  -e PUID=1000 \
  -e PGID=1000 \
  -e TZ=Europe/London \
  -e SERVERPORT=35053 \
  -e ALLOWEDIPS=0.0.0.0/0 \
  -p 35053:51820/udp \
  -e PEERDNS=auto \
  -v /etc/wireguard:/config \
  -v /lib/modules:/lib/modules \
  --sysctl net.ipv6.conf.all.disable_ipv6=0 \
  --restart unless-stopped \
  ghcr.io/linuxserver/wireguard

  docker exec -it wireguard bash