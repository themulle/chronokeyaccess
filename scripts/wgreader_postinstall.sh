#!/bin/bash
# Postinstall script for chronokeyaccess_wgreader

# Reload systemd to pick up new service
systemctl daemon-reload

# Enable the service so it starts at boot
systemctl enable chronokeyaccess_wgreader.service

# Start the service immediately
systemctl start chronokeyaccess_wgreader.service
