#!/bin/bash
sudo yum install epel-release elrepo-release -y
sudo yum install yum-plugin-elrepo -y
sudo yum install kmod-wireguard wireguard-tools -y 
