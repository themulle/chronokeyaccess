#!/bin/bash
# Preinstall script for chronokeyaccess_web

# Add user and group if they don't exist
if ! id "chronokeyaccess" >/dev/null 2>&1; then
    useradd -r -m -s /usr/sbin/nologin chronokeyaccess
    mkdir -m 700 /home/chronokeyaccess/.ssh
    ssh-keygen -t rsa -b 2048 -q -N "" -f /home/chronokeyaccess/.ssh/id_rsa
    chown -R chronokeyaccess:chronokeyaccess /home/chronokeyaccess/.ssh

fi

# Create directory for configuration files if it doesn't exist
if [ ! -d "/etc/chronokeyaccess" ]; then
    mkdir -p /etc/chronokeyaccess
    chown chronokeyaccess:chronokeyaccess /etc/chronokeyaccess
fi
