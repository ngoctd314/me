---
title: "Automatic connect to VPN"
date: 2022-11-08
authors: ["ngoctd"]
draft: false
tags: ["linux"]
---

The OpenVPN is an open source Virtual Private Network (VPN) project. It creates secure connections over the Internet using a custom security protocol that utilizes SSL/TLS.

In this article, i will introduce naive solution for automatic connect to VPN using openvpn and systemd management.

**Create vpn script**
```bash
# Create script connect to vpn
mkdir -p ~/scripts
cd ~/scripts
touch vpn.sh
sudo chmod +x vpn.sh
```

**Install oauthtool**
```bash
sudo apt install oathtool 
```

**~/scripts/vpn.sh**
```bash
VPN_USER="" # insert vpn user here
VPN_PASSWORD=""  # insert vpn password here
OTP_KEY="" # insert otp key here
OVPN_FILE="" # insert path to .ovpn file here. example /home/admicro-bigdata.ovpn

VPN_AUTH="$(oathtool -b --totp $OTP_KEY)$VPN_PASSWORD" 
echo $VPN_PASSWORD | sudo -S bash -c "openvpn --config $OVPN_FILE --auth-user-pass <(echo -e '$VPN_USER\n$VPN_AUTH') --daemon"
```

**Create systemd service**
```bash
cd /lib/systemd/system
sudo touch vpn.service
```

**/lib/systemd/system/vpn.service**
```txt
[Unit]
Description=Auto connect vccorp's vpn service.

[Service]
Type=forking
User=root
ExecStart= # path to vpn file example: /home/ngoctd/scripts/vpn.sh
ExecReload= # path to vpn file example: /home/ngoctd/scripts/vpn.sh
ExecStop=sudo killall openvpn

Restart=on-failure
RestartSec=10s

[Install]
WantedBy=multi-user.target
```

**start,stop,auto restart service**
```bash
sudo systemctl enable vpn # auto start when vpn was killed
sudo systemctl start vpn # start vpn
sudo systemctl status vpn # check vpn status
sudo systemctl stop vpn # stop vpn
```