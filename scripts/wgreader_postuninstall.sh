#!/bin/bash
# Postuninstall script for chronokeyaccess_wgreader

# Stop the service
systemctl stop chronokeyaccess_wgreader.service

# Disable the service
systemctl disable chronokeyaccess_wgreader.service

# Remove the service file
rm -f /etc/systemd/system/chronokeyaccess_wgreader.service

# Reload systemd to reflect the removal
systemctl daemon-reload