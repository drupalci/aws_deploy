#!/bin/bash

# Name:        provision.sh
# Author:      Nick Schuch
# Description: Provisions a Vagrant environment for Golang.

# Install Golang so we can compile and run.
cd /tmp
wget https://storage.googleapis.com/golang/go1.4.1.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.4.1.linux-amd64.tar.gz
mkdir -p /opt/golang
chmod -R 777 /opt/golang

echo 'GOPATH=/opt/golang' >> /etc/environment
echo 'PATH=$PATH:/usr/local/go/bin' >> /etc/profile
echo 'PATH=$PATH:/opt/golang/bin' >> /etc/profile
export GOPATH=/opt/golang
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:/opt/golang/bin

source /vagrant/scripts/crosscompile.bash
go-crosscompile-build-all

# Some other random packages.
export DEBIAN_FRONTEND=noninteractive
sudo -E apt-get install -y vim make git
