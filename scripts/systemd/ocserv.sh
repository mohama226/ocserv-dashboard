#!/bin/bash
# ======================================================================
# Script: systemd_ocserv.sh
# ======================================================================

source /etc/os-release
source ./scripts/lib.sh

OCSERV_PORT="${OCSERV_PORT:-443}"
OC_NET="${OC_NET:-172.16.24.0/24}"
OCSERV_DNS="${OCSERV_DNS:-1.1.1.1}"
ETH="${ETH:-}"

auto_detect_interface() {
  if [[ -z "${ETH}" ]]; then
    ETH="$(ip -o -4 route show to default 2>/dev/null | awk '{print $5}' | head -n1 || true)"
    [[ -n "${ETH}" ]] || die "Could not detect external interface. Set ETH manually."
    log "Auto-detected external interface: ${ETH}"
  fi
}
auto_detect_interface

# ===============================
# Install Ocserv (OS-based)
# ===============================
install_ocserv(){

case "$ID" in

    ubuntu|debian)

        apt update
        apt install -y ocserv

        ;;

    almalinux|rocky|rhel|centos)

        dnf install -y epel-release
        dnf install -y ocserv || {
            echo "ocserv package not found, compiling from source"
            exit 1
        }

        ;;

    *)

        die "Unsupported OS: $ID"

        ;;

esac

}

install_ocserv

log "Installing dependencies..."
sudo apt-get install -y gnutls-bin openssl iptables iptables-persistent || true
sudo dnf install -y gnutls openssl iptables iptables-services || true

# ==============================================================
# 2. Generate Ocserv Certificates
# ==============================================================
generate_ocserv_certs() {
  log "Generating SSL certificates..."

  sudo mkdir -p /etc/ocserv/certs
  sudo touch /etc/ocserv/ocpasswd

  SERVER_CERT="cert.pem"
  SERVER_KEY="key.pem"

  SSL_CN="${SSL_CN:-End-way-Cisco-VPN}"
  SSL_ORG="${SSL_ORG:-End-way}"
  SSL_EXPIRE="${SSL_EXPIRE:-3650}"

  sudo certtool --generate-privkey --outfile ca-key.pem

  cat <<_EOF_ | sudo tee ca.tmpl >/dev/null
cn = "${SSL_CN}"
organization = "${SSL_ORG}"
serial = 1
expiration_days = ${SSL_EXPIRE}
ca
signing_key
cert_signing_key
crl_signing_key
_EOF_

  sudo certtool --generate-self-signed \
    --load-privkey ca-key.pem \
    --template ca.tmpl \
    --outfile ca-cert.pem

  sudo certtool --generate-privkey --outfile "${SERVER_KEY}"

  cat <<_EOF_ | sudo tee server.tmpl >/dev/null
cn = "${SSL_CN}"
organization = "${SSL_ORG}"
serial = 2
expiration_days = ${SSL_EXPIRE}
signing_key
encryption_key
tls_www_server
_EOF_

  sudo certtool --generate-certificate \
    --load-privkey "${SERVER_KEY}" \
    --load-ca-certificate ca-cert.pem \
    --load-ca-privkey ca-key.pem \
    --template server.tmpl \
    --outfile "${SERVER_CERT}"

  sudo cp "${SERVER_CERT}" /etc/ocserv/certs/cert.pem
  sudo cp "${SERVER_KEY}" /etc/ocserv/certs/cert.key
}

if [[ ! -f /etc/ocserv/certs/cert.pem ]]; then
  generate_ocserv_certs
fi

# ==============================================================
# 3. Ocserv Main Configuration
# ==============================================================
OCSERV_CONF="/etc/ocserv/ocserv.conf"
MANAGED_HEADER="# Managed by ocserv-dashboard install.sh"

write_ocserv_conf_systemd() {
  log "Writing Ocserv configuration..."
  sudo tee "$OCSERV_CONF" >/dev/null <<EOT
# ===============================================
# Managed by ocserv-dashboard install.sh
# ===============================================

auth = "certificate"
enable-auth = "plain[passwd=/etc/ocserv/ocpasswd]"
ca-cert = /etc/ocserv/ssl/ca-cert.pem
crl = /etc/ocserv/ssl/crl.pem
cert-user-oid = 2.5.4.3

run-as-user = root
run-as-group = root

socket-file = /var/run/ocserv-socket
isolate-workers = true
max-clients = 1024

server-cert = /etc/ocserv/certs/cert.pem
server-key  = /etc/ocserv/certs/cert.key

dns = ${OCSERV_DNS}
ipv4-network = ${OC_NET}

tcp-port = ${OCSERV_PORT}
udp-port = ${OCSERV_PORT}

pre-login-banner="${OCSERV_PRE_LOGIN_BANNER}"
EOT

OCSERV_BANNER=$(echo "$OCSERV_BANNER" | awk '{printf "%s\\n", $0}' | sed 's/\\n$//')
printf 'banner = "%s"\n' "$OCSERV_BANNER" | sudo tee -a "$OCSERV_CONF" > /dev/null
}

if [[ ! -f "$OCSERV_CONF" ]]; then
    write_ocserv_conf_systemd
elif ! head -n 5 "$OCSERV_CONF" | grep -q "$MANAGED_HEADER"; then
    write_ocserv_conf_systemd
fi

sudo mkdir -p /etc/ocserv/defaults /etc/ocserv/groups /etc/ocserv/users
sudo touch /etc/ocserv/defaults/group.conf

# ==============================================================
# 4. Enable Kernel Forwarding
# ==============================================================
log "Enabling IP forwarding..."

sudo sysctl -w net.ipv4.ip_forward=1
echo "net.ipv4.ip_forward = 1" | sudo tee /etc/sysctl.d/99-ocserv.conf >/dev/null
sudo sysctl --system

# ==============================================================
# 5. Firewall Rules / NAT
# ==============================================================
log "Configuring firewall..."

sudo iptables -A INPUT -p tcp --dport "${OCSERV_PORT}" -j ACCEPT || true
sudo iptables -A INPUT -p udp --dport "${OCSERV_PORT}" -j ACCEPT || true

sudo iptables -t nat -A POSTROUTING -s "${OC_NET}" -o "${ETH}" -j MASQUERADE || true

sudo iptables -A FORWARD -s "${OC_NET}" -o "${ETH}" -j ACCEPT || true
sudo iptables -A FORWARD -d "${OC_NET}" -m state --state ESTABLISHED,RELATED -j ACCEPT || true

sudo sh -c "iptables-save > /etc/iptables/rules.v4"

# ==============================================================
# 6. Start & Enable Ocserv Service
# ==============================================================
info "Starting Ocserv..."

sudo systemctl daemon-reload
sudo systemctl enable ocserv.service
sudo systemctl restart ocserv.service

OCSERV_VERSION=$(ocserv --version | head -n 1)

info "ocserv ${OCSERV_VERSION} installed successfully!"
ok "Ocserv VPN deployment completed successfully!"
