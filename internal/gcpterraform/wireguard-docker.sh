#!/bin/bash

docker stop wireguard && docker rm wireguard

docker run -d \
  --name=wireguard \
  --cap-add=NET_ADMIN \
  --cap-add=SYS_MODULE \
  -e PUID=1000 \
  -e PGID=1000 \
  -e TZ=Europe/London \
  -e SERVERPORT=51820 \
  -e ALLOWEDIPS=0.0.0.0/0 \
  -p 51820:51820/udp \
  -e PEERS=1 \
  -e PEERDNS=auto \
  -v /etc/wireguard/wg0.conf:/etc/wireguard/wg0.conf:rw \
  -v /lib/modules:/lib/modules \
  --restart unless-stopped \
  ghcr.io/linuxserver/wireguard

  docker exec -it wireguard bash