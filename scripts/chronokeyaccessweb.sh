#!/bin/bash
### BEGIN INIT INFO
# Provides:          web_rpi
# Required-Start:    $remote_fs $syslog
# Required-Stop:     $remote_fs $syslog
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: Starts web_rpi in /home/svr
# Description:       This file starts and stops web_rpi running in the home directory of user svr
### END INIT INFO

# Variablen
USER=svr
APP_DIR="/home/$USER"
APP="web_rpi"
PIDFILE="/var/run/$APP.pid"
LOGFILE="/var/log/$APP.log"

start() {
    echo "Starting $APP..."
    if [ -f $PIDFILE ]; then
        echo "$APP is already running!"
    else
        cd $APP_DIR
        sudo -u $USER nohup ./$APP > $LOGFILE 2>&1 &
        echo $! > $PIDFILE
        echo "$APP started!"
    fi
}

stop() {
    echo "Stopping $APP..."
    if [ ! -f $PIDFILE ]; then
        echo "$APP is not running!"
    else
        PID=$(cat $PIDFILE)
        kill $PID && rm -f $PIDFILE
        echo "$APP stopped!"
    fi
}

status() {
    if [ -f $PIDFILE ]; then
        echo "$APP is running with PID $(cat $PIDFILE)"
    else
        echo "$APP is not running"
    fi
}

case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    status)
        status
        ;;
    restart)
        stop
        start
        ;;
    *)
        echo "Usage: /etc/init.d/web_rpi {start|stop|status|restart}"
        exit 1
        ;;
esac

exit 0
