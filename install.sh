#!/bin/bash
# ==============================================================
# Script: install.sh
# Description:
#   Interactive installer for the OpenConnect VPN dashboard
# ==============================================================

# Load shared helpers
source ./scripts/lib.sh

check_base_packages(){

    REQUIRED="curl wget openssl"

    for cmd in $REQUIRED
    do
        if ! command -v $cmd >/dev/null 2>&1
        then
            echo "$cmd missing, installing..."

            if command -v apt >/dev/null 2>&1
            then
                apt update
                apt install -y $cmd

            elif command -v dnf >/dev/null 2>&1
            then
                dnf install -y $cmd
            fi
        fi
    done
}

# ===============================
# OS Detection
# ===============================
detect_os(){

    if [ ! -f /etc/os-release ]; then
        print_message error "Cannot detect operating system"
        exit 1
    fi

    source /etc/os-release

    case "$ID" in
        ubuntu|debian)
            OS_FAMILY="debian"
            PACKAGE_MANAGER="apt"
            ;;
        almalinux|rocky|rhel|centos)
            OS_FAMILY="rhel"
            PACKAGE_MANAGER="dnf"
            ;;
        *)
            print_message error "Unsupported OS: $PRETTY_NAME"
            exit 1
            ;;
    esac

    export OS_FAMILY
    export PACKAGE_MANAGER

    print_message success "Detected OS: $PRETTY_NAME"
    print_message success "Package Manager: $PACKAGE_MANAGER"
}

# ===============================
# Install Package Wrapper
# ===============================
install_package(){

    if [[ "$PACKAGE_MANAGER" == "apt" ]]; then
        apt update -y
        apt install -y "$@"

    elif [[ "$PACKAGE_MANAGER" == "dnf" ]]; then
        dnf install -y "$@"
    fi
}

# ===============================
# Default Configuration
# ===============================
ENV_FILE=".env"
HOST=$(hostname -I | awk '{print $1}')
SSL_CN="End-way-Cisco-VPN"
SSL_ORG="End-way"
SSL_EXPIRE=3650
OC_NET="172.16.24.0/24"
OCSERV_PORT=443
OCSERV_DNS="8.8.8.8"
OCSERV_BANNER="Welcome to the VPN service"
OCSERV_PRE_LOGIN_BANNER="Welcome"
LANGUAGES="en:English,it:Italiano,zh-cn:中文(简体),zh-tw:中文(繁體),ru:Русский,fa:فارسی,ar:العربية"
SECRET_KEY=$(openssl rand -hex 32)
JWT_SECRET=$(openssl rand -hex 32)
SSL_C=US
SSL_ST=CA
SSL_L=SanFrancisco
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=ocserv
POSTGRES_USER=ocserv
POSTGRES_PASSWORD=ocserv-passwd
CURRENT_RELEASE=
LATEST_RELEASE=

# ===============================
# ensure_root
# ===============================
ensure_root() {
    if ! command -v sudo >/dev/null 2>&1; then
        print_message error "❌ Error: sudo is not installed."
        exit 1
    fi
}

# ===============================
# choose_deployment
# ===============================
choose_deployment() {
    print_message info "🚀 Deployment Options:"
    print_message highlight "   [1] Docker"
    print_message highlight "   [2] Systemd Full (Ocserv + Dashboard)"
    print_message highlight "   [3] Systemd Dashboard (Standalone Setup/Upgrade)"
    print_message highlight "   [4] Update"
    print_message highlight "   [5] Uninstall"

    read -rp "Choose deployment method [1-5] (default = 1): " choice
    choice=${choice:-1}

    case "$choice" in
        1) DEPLOY_METHOD="docker" ;;
        2) DEPLOY_METHOD="systemd" ;;
        3) DEPLOY_METHOD="standalone" ;;
        4) DEPLOY_METHOD="update" ;;
        5) DEPLOY_METHOD="uninstall" ;;
        *)
            print_message warn "Invalid choice, defaulting to Docker."
            DEPLOY_METHOD="docker"
            ;;
    esac

    print_message highlight "✅ Selected deployment method: ${DEPLOY_METHOD}"
    printf "\n"
}

# ===============================
# check_docker
# ===============================
check_docker() {
    local missing=0

    if ! command -v sudo docker &> /dev/null; then
        print_message error "❌ Docker is not installed."
        missing=1
    else
        print_message success "✅ Docker is installed."
        docker_info=$(sudo docker info --format 'Server Version: {{.ServerVersion}}')
        print_message highlight "🔹 $docker_info"
    fi

    if ! sudo docker compose version &> /dev/null; then
        print_message error "❌ Docker Compose plugin missing."
        missing=1
    else
        print_message success "✅ Docker Compose plugin installed."
        compose_version=$(sudo docker compose version | head -n1)
        print_message highlight "🔹 $compose_version"
    fi

    if [[ $missing -eq 1 ]]; then
        print_message info "🔗 Installation guides:"
        print_message highlight "   Docker: https://docs.docker.com/get-docker/"
        print_message highlight "   Docker Compose: https://docs.docker.com/compose/install/"
        exit 1
    fi
}

# ===============================
# get_ip
# ===============================
get_ip() {
    print_message info "🔍 Detecting public IP ..."
    local detected_ip
    detected_ip=$(
        curl -s --max-time 5 https://api.ipify.org || \
        curl -s --max-time 5 https://ifconfig.me || \
        curl -s --max-time 5 https://checkip.amazonaws.com
    ) || true

    print_message info "Detected IP: $detected_ip"

    if [[ "$detected_ip" =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        read -rp "Use this IP? [Y/n]: " choice
        HOST=${detected_ip}
        [[ "$choice" =~ [Nn] ]] && read -rp "Enter your VPS host or IP: " HOST
    else
        read -rp "Enter your VPS host or IP: " HOST
    fi

    print_message highlight "🔧 Using host IP: ${HOST}"
    printf "\n"
}

# ===============================
# generate_secret
# ===============================
generate_secret() {
    local len=64

    if ! command -v openssl >/dev/null 2>&1; then
        install_package openssl
    fi

    openssl rand -hex 64 | head -c "$len"
}

# ===============================
# get_envs
# ===============================
# (بدون تغییر)

# ===============================
# get_site_lang
# ===============================
# (بدون تغییر)

# ===============================
# get_mirrors
# ===============================
# (بدون تغییر)

# ===============================
# set_environment
# ===============================
# (بدون تغییر)

# ===============================
# get_interface
# ===============================
# (بدون تغییر)

# ===============================
# setup_docker
# ===============================
# (بدون تغییر)

# ===============================
# setup_systemd
# ===============================
setup_systemd() {
    local full_setup="$1"

    # ... (بدون تغییر تا بخش آخر)

    if [[ "$OS_FAMILY" == "debian" ]]; then
        apt autoremove -y
        apt autoclean -y

    elif [[ "$OS_FAMILY" == "rhel" ]]; then
        dnf autoremove -y || true
        dnf clean all
    fi

    ok "Cleanup completed."
}

# ===============================
# uninstall
# ===============================
# (بدون تغییر)

# ===============================
# deploy
# ===============================
# (بدون تغییر)

# ===============================
# get_current_version
# ===============================
# (بدون تغییر)

# ===============================
# main
# ===============================
main() {
    ensure_root
    detect_os
    check_base_packages
    choose_deployment

    if [[ "$DEPLOY_METHOD" == "uninstall" ]]; then
        print_message info "⚠️ Uninstall mode selected."
        uninstall
        exit 0
    fi

    if [[ "$DEPLOY_METHOD" == "update" ]]; then
        print_message info "♻️ Update mode selected."
        ./scripts/update.sh
        exit 0
    fi

    install_package curl

    if [[ "$DEPLOY_METHOD" == "docker" ]]; then
        check_docker
    else
        ./scripts/systemd/pre_requirements.sh
    fi

    ENV_FILE=".env"
    if [[ -f "$ENV_FILE" ]]; then
        print_message info "✅ Loading environment from $ENV_FILE"
        set -o allexport
        source "$ENV_FILE"
        set +o allexport
        print_message success "✅ Environment loaded"
    else
        print_message info "⚡ No .env found. Running interactive setup..."
        get_ip
        get_envs
        get_site_lang

        read -rp "Do you want to configure development mirrors? [y/N]: " configure_mirrors
        configure_mirrors=${configure_mirrors:-N}
        if [[ "$configure_mirrors" =~ ^[Yy]$ ]]; then
            get_mirrors
        else
            print_message info "Skipping mirror configuration - using defaults"
            printf "\n"
        fi

        set_environment
    fi

    if ! get_current_version; then
        echo "Failed to get latest version"
    fi

    echo "Current Dashboard Release: $CURRENT_RELEASE"
    echo "Latest Dashboard Release: $LATEST_RELEASE"

    deploy
}

main
