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
  -e PEERS=1 \
  -e PEERDNS=auto \
  -v /etc/wireguard/wg0.conf:/config/wg0.conf \
  -v /lib/modules:/lib/modules \
  --restart unless-stopped \
  ghcr.io/linuxserver/wireguard