#!/bin/sh
apt update
apt install -y wget bzip2 cmake g++-arm-linux-gnueabihf
export CC=arm-linux-gnueabihf-gcc
cd /tmp
#libusb-1.0
cd /tmp
wget https://github.com/libusb/libusb/releases/download/v1.0.26/libusb-1.0.26.tar.bz2
tar xjf libusb-1.0.26.tar.bz2
cd libusb-1.0.26
./configure --host=arm-linux --enable-udev=no
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
CC=arm-linux-gnueabihf-gcc GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 go build -o $1/twSdrPower.arm  -ldflags="-extldflags -w -X main.version=$2 -X main.commit=$3"

