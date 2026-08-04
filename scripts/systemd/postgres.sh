#!/bin/bash
# ===============================
# Script: systemd_postgres.sh
# ===============================

source /etc/os-release

set -e
trap 'die "Error on line $LINENO"' ERR

# Load helpers
source ./scripts/lib.sh

# ===============================
# Load environment variables
# ===============================
ENV_FILE=".env"

if [ ! -f "$ENV_FILE" ]; then
    die ".env file not found"
fi

export $(grep -v '^#' "$ENV_FILE" | xargs)

[ -z "$POSTGRES_DB" ] && die "POSTGRES_DB is not set"
[ -z "$POSTGRES_USER" ] && die "POSTGRES_USER is not set"
[ -z "$POSTGRES_PASSWORD" ] && die "POSTGRES_PASSWORD is not set"

ok "Environment loaded"

# ===============================
# Install PostgreSQL (OS-based)
# ===============================
ok "Installing PostgreSQL..."

case "$ID" in

    ubuntu|debian)

        apt update -y
        apt install -y postgresql postgresql-client

        ;;

    almalinux|rocky|rhel|centos)

        dnf install -y postgresql-server postgresql-contrib

        postgresql-setup --initdb

        systemctl enable postgresql
        systemctl start postgresql

        ;;

    *)

        die "Unsupported OS: $ID"

        ;;

esac

ok "PostgreSQL installed"

# ===============================
# Start service (Debian-based)
# ===============================
if [[ "$ID" == "ubuntu" || "$ID" == "debian" ]]; then
    ok "Starting PostgreSQL..."

    systemctl enable postgresql
    systemctl restart postgresql

    ok "PostgreSQL is running"
fi

# ===============================
# Create USER (safe)
# ===============================
ok "Creating user..."

sudo -u postgres psql <<EOF
DO \$\$
BEGIN
   IF NOT EXISTS (
      SELECT 1 FROM pg_roles WHERE rolname = '$POSTGRES_USER'
   ) THEN
      CREATE USER $POSTGRES_USER WITH PASSWORD '$POSTGRES_PASSWORD';
   END IF;
END
\$\$;
EOF

# ===============================
# Create DATABASE
# ===============================
ok "Creating database..."

if ! sudo -u postgres psql -tAc "SELECT 1 FROM pg_database WHERE datname='$POSTGRES_DB'" | grep -q 1; then
    sudo -u postgres psql -c "CREATE DATABASE $POSTGRES_DB OWNER $POSTGRES_USER"
fi

# ===============================
# Grant privileges
# ===============================
ok "Granting privileges..."

sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE $POSTGRES_DB TO $POSTGRES_USER;"

ok "Database and user configured"

# ===============================
# Final output
# ===============================
ok "PostgreSQL setup complete"

echo "--------------------------------------"
echo "Database : $POSTGRES_DB"
echo "User     : $POSTGRES_USER"
echo "Host     : localhost"
echo "Port     : 5432"
echo "--------------------------------------"

ok "PostgreSQL configured ✅"
