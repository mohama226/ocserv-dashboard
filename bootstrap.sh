#!/bin/bash

set -e

if [ -f /etc/os-release ]; then
    source /etc/os-release
else
    echo "Cannot detect OS"
    exit 1
fi


echo "Detected OS: $PRETTY_NAME"


case "$ID" in

ubuntu|debian)

    apt update

    apt install -y \
    git \
    curl \
    wget \
    sudo \
    openssl \
    ca-certificates

    ;;


almalinux|rocky|rhel|centos)

    dnf update -y

    dnf install -y \
    git \
    curl \
    wget \
    sudo \
    openssl \
    ca-certificates

    ;;


*)

    echo "Unsupported OS"
    exit 1

    ;;

esac


echo "Bootstrap packages installed successfully"
