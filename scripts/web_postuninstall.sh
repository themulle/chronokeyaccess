#!/bin/bash
# Postuninstall script for chronokeyaccess_web

# Stop the service
systemctl stop chronokeyaccess_web.service

# Disable the service
systemctl disable chronokeyaccess_web.service

# Remove the service file
rm -f /etc/systemd/system/chronokeyaccess_web.service

# Reload systemd to reflect the removal
systemctl daemon-reload