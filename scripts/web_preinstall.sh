#!/bin/bash
# Preinstall script for chronokeyaccess_web

# Add user and group if they don't exist
if ! id "chronokeyaccess" >/dev/null 2>&1; then
    useradd --system --home /nonexistent --shell /bin/false chronokeyaccess
fi

# Create directory for configuration files if it doesn't exist
if [ ! -d "/etc/chronokeyaccess" ]; then
    mkdir -p /etc/chronokeyaccess
    chown chronokeyaccess:chronokeyaccess /etc/chronokeyaccess
fi
