#!/bin/sh
apt update
apt install -y wget bzip2 cmake g++-aarch64-linux-gnu
export CC=aarch64-linux-gnu-gcc
#libusb-1.0
cd /tmp
wget https://github.com/libusb/libusb/releases/download/v1.0.26/libusb-1.0.26.tar.bz2
tar xjf libusb-1.0.26.tar.bz2
cd libusb-1.0.26
./configure --host=aarch64-linux --enable-udev=no
make 
make install
#librtlsdr
cd /tmp
git clone https://gitea.osmocom.org/sdr/rtl-sdr.git
cd rtl-sdr/
mkdir build
cd build
cmake -D CMAKE_INSTALL_PREFIX=/usr/ ../
make
make install
cd /twSdrPower
go mod tidy
CC=aarch64-linux-gnu-gcc GOOS=linux GOARCH=arm64  CGO_ENABLED=1 go build -o $1/twSdrPower.arm64  -ldflags="-extldflags -w -X main.version=$2 -X main.commit=$3"

