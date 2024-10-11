#!/bin/bash

# Check if the user provided a code as an argument
if [ -z "$1" ]; then
  echo "Usage: $0 <code>"
  exit 1
fi

# Store the provided code in a variable
CODE=$1

# Run the codechecker_rpi command with the provided code
RESULT=$(chronokeyaccess_cli -t "$CODE")

# Check if the output of codechecker_rpi is "ok"
if [ "$RESULT" == "ok" ]; then
  echo "Code verified, opening door..."
  # Run the dooropener_rpi command
  chronokeyaccess_dooropener
else
  echo "Code verification failed: $RESULT"
fi