#!/bin/bash

# Preinstall script to enable RTC on Raspberry Pi with Ubuntu

set -e

# Function to check if a module is already loaded
module_loaded() {
    lsmod | grep -q "$1"
}

# Function to check if a service exists and is enabled
service_enabled() {
    systemctl is-enabled "$1" &>/dev/null
}

# Step 1: Ensure i2c-tools are installed
if ! dpkg -l | grep -q "i2c-tools"; then
    echo "Installing i2c-tools..."
    sudo apt-get update && sudo apt-get install -y i2c-tools
else
    echo "i2c-tools already installed."
fi

# Step 2: Check if I2C is enabled in config.txt
if ! grep -q "^dtparam=i2c_arm=on" /boot/firmware/config.txt; then
    echo "Enabling I2C interface..."
    echo "dtparam=i2c_arm=on" | sudo tee -a /boot/firmware/config.txt > /dev/null
else
    echo "I2C interface already enabled."
fi

# Step 3: Load the RTC module if not already loaded
if ! module_loaded "rtc_ds1307"; then
    echo "Loading rtc-ds1307 module..."
    sudo modprobe rtc-ds1307
else
    echo "rtc-ds1307 module already loaded."
fi

# Step 4: Add the RTC device only if not already added
if [ ! -e /sys/class/i2c-adapter/i2c-1/1-0068/rtc ]; then
    echo "Adding RTC device..."
    echo "ds1307 0x68" | sudo tee /sys/class/i2c-adapter/i2c-1/new_device > /dev/null
else
    echo "RTC device already added."
fi

# Step 5: Sync system time to rtc
echo "Syncing system time to rtc..."
date
sudo hwclock -w

# Step 6: Add rtc-ds1307 to modules.conf for auto-load at boot
if ! grep -q "^rtc-ds1307" /etc/modules-load.d/modules.conf 2>/dev/null; then
    echo "Adding rtc-ds1307 to /etc/modules-load.d/modules.conf..."
    echo "rtc-ds1307" | sudo tee -a /etc/modules-load.d/modules.conf > /dev/null
else
    echo "rtc-ds1307 already added to modules.conf."
fi

# Step 7: Create /etc/rtc script if not exists
if [ ! -f /etc/rtc ]; then
    echo "Creating /etc/rtc script..."
    sudo bash -c 'cat << EOF > /etc/rtc
#!/bin/bash
echo "ds1307 0x68" | sudo tee /sys/class/i2c-adapter/i2c-1/new_device > /dev/null
sudo hwclock -s
EOF'
    sudo chmod +x /etc/rtc
else
    echo "/etc/rtc script already exists."
fi

# Step 8: Create systemd service if not exists
if [ ! -f /etc/systemd/system/rtc.service ]; then
    echo "Creating rtc.service..."
    sudo bash -c 'cat << EOF > /etc/systemd/system/rtc.service
[Unit]
Description=RTC Clock
Before=cloud-init-local.service
Requires=systemd-modules-load.service
After=systemd-modules-load.service

[Service]
Type=oneshot
ExecStart=/etc/rtc

[Install]
WantedBy=multi-user.target
EOF'
    sudo systemctl daemon-reload
else
    echo "rtc.service already exists."
fi

# Step 9: Enable rtc.service if not already enabled
if ! service_enabled "rtc.service"; then
    echo "Enabling rtc.service..."
    sudo systemctl enable rtc.service
else
    echo "rtc.service already enabled."
fi

echo "RTC setup complete. You can now reboot the system."
