#!/bin/bash

set -e

source /etc/os-release

GO_VERSION="1.25.1"

install_go() {

    if command -v go >/dev/null 2>&1; then
        echo "Go already installed: $(go version)"
        return
    fi

    ARCH=$(uname -m)

    case "$ARCH" in
        x86_64)
            GO_ARCH="amd64"
            ;;
        aarch64)
            GO_ARCH="arm64"
            ;;
        *)
            echo "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac

    cd /tmp

    wget -O go.tar.gz https://go.dev/dl/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz

    rm -rf /usr/local/go

    tar -C /usr/local -xzf go.tar.gz

    cat >/etc/profile.d/go.sh <<EOF
export PATH=\$PATH:/usr/local/go/bin
EOF

    export PATH=$PATH:/usr/local/go/bin

    go version
}

install_debian() {

    apt update

    apt install -y \
        curl \
        wget \
        git \
        gcc \
        g++ \
        make \
        pkg-config \
        build-essential \
        nodejs \
        npm \
        nginx \
        openssl \
        ca-certificates \
        certbot \
        postgresql-client \
        unzip \
        tar \
        jq \
        nano

}

install_rhel() {

    dnf -y update

    dnf -y install epel-release

    dnf config-manager --set-enabled crb || true

    dnf -y install \
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
        ca-certificates \
        certbot \
        unzip \
        tar \
        jq \
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
        echo "Unsupported OS: $PRETTY_NAME"
        exit 1
        ;;

esac

install_go

echo "Requirements installed successfully"
