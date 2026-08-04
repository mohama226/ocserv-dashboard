#!/bin/bash

set -e


source /etc/os-release


install_go(){

GO_VERSION="1.25.0"

if ! command -v go >/dev/null 2>&1; then

    cd /tmp

    wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz

    rm -rf /usr/local/go

    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz


    cat >/etc/profile.d/go.sh <<EOF
export PATH=\$PATH:/usr/local/go/bin
EOF

    source /etc/profile.d/go.sh

fi


}


install_debian(){

apt update


apt install -y \
curl \
wget \
git \
gcc \
make \
pkg-config \
build-essential \
nodejs \
npm \
nginx \
openssl \
certbot \
postgresql-client \
nano


}


install_rhel(){

dnf update -y


dnf install -y epel-release


dnf install -y \
curl \
wget \
git \
gcc \
gcc-c++ \
make \
pkgconf-pkg-config \
nodejs \
npm \
nginx \
openssl \
openssl-devel \
certbot \
postgresql \
postgresql-server \
nano



systemctl enable nginx
systemctl start nginx


}



case "$ID" in


ubuntu|debian)

    install_debian

    ;;


almalinux|rocky|rhel|centos)

    install_rhel

    ;;


*)

echo "Unsupported OS"

exit 1


;;

esac



install_go


echo "Requirements installed successfully"
