#!/bin/bash

# Dieses Skript installiert und aktiviert die DS3231 RTC auf einem Raspberry Pi

# Check for root privileges
if [ "$EUID" -ne 0 ]; then
  echo "Bitte das Skript als root oder mit sudo ausführen."
  exit 1
fi

echo "1. Überprüfen, ob I2C aktiviert ist..."
if ! grep -q "^dtparam=i2c_arm=on" /boot/config.txt; then
  echo "Aktivieren von I2C..."
  echo "dtparam=i2c_arm=on" >> /boot/config.txt
fi

# Installiere I2C Tools
echo "2. I2C Tools und RTC Treiber installieren..."
apt-get update
apt-get install -y i2c-tools

# Überprüfe, ob DS3231 erkannt wird
echo "3. Prüfen, ob die DS3231 RTC erkannt wird..."
i2cdetect -y 1

echo "4. RTC Modul dem Kernel hinzufügen..."
if ! grep -q "^dtoverlay=i2c-rtc,ds3231" /boot/config.txt; then
  echo "Hinzufügen des RTC Moduls in /boot/config.txt..."
  echo "dtoverlay=i2c-rtc,ds3231" >> /boot/config.txt
fi

# Deaktiviere den Fake-HW-Clock Dienst (falls installiert)
echo "5. Deaktivieren des Fake-HW-Clock-Dienstes..."
apt-get -y remove fake-hwclock
systemctl disable fake-hwclock
systemctl stop fake-hwclock

# Synchronisiere die Uhrzeit von der RTC
echo "6. Synchronisiere die Uhrzeit von der RTC..."
hwclock -r

# Stelle sicher, dass die Systemuhrzeit bei jedem Boot synchronisiert wird
echo "7. Überprüfe den Service zum Synchronisieren der RTC mit der Systemuhr..."

# Füge einen RTC-Synchronisationsbefehl zu /etc/rc.local hinzu, falls nicht vorhanden
if ! grep -q "hwclock -s" /etc/rc.local; then
  sed -i -e '$i \hwclock -s\n' /etc/rc.local
fi

echo "8. Neustart des Raspberry Pi für Änderungen erforderlich. Möchtest du jetzt neu starten? (y/n)"
read -r RESTART

if [ "$RESTART" == "y" ]; then
  echo "Neustart..."
  reboot
else
  echo "Bitte starte das System manuell neu, um die Änderungen zu übernehmen."
fi
