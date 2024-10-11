#!/bin/bash

# Variablen
DOMAIN="themulle.de"
EMAIL="themulle@gmail.com"  # Ersetze durch deine E-Mail-Adresse für Let's Encrypt
NGINX_CONF="/etc/nginx/sites-available/$DOMAIN"
NGINX_CONF_LINK="/etc/nginx/sites-enabled/$DOMAIN"
SERVICE_PORT="8080"

# Update des Systems und Installation von Nginx und Certbot
sudo apt update
sudo apt install -y nginx certbot python3-certbot-nginx

# Nginx Konfiguration für die Domain erstellen
sudo tee $NGINX_CONF > /dev/null <<EOL
server {
    listen 80;
    server_name $DOMAIN www.$DOMAIN;
    return 301 http://$DOMAIN\$request_uri;
}

server {
    listen 443;
    server_name $DOMAIN www.$DOMAIN;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}

EOL

# Symbolischen Link erstellen, um die Konfiguration zu aktivieren
sudo ln -s $NGINX_CONF $NGINX_CONF_LINK

# Überprüfen der Nginx-Konfiguration
sudo nginx -t

# Nginx neu laden
sudo systemctl reload nginx

# Zertifikat von Let's Encrypt holen und automatisch konfigurieren
sudo certbot --nginx -d $DOMAIN -d www.$DOMAIN --non-interactive --agree-tos --email $EMAIL

# Automatische Zertifikatserneuerung testen
sudo certbot renew --dry-run

# Nginx neustarten, um SSL zu aktivieren
sudo systemctl restart nginx

echo "Nginx ist konfiguriert und läuft als Reverse Proxy auf 127.0.0.1:$SERVICE_PORT mit SSL für $DOMAIN"
