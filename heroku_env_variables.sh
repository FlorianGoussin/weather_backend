#!/bin/bash

HEROKU_APP_NAME="floriang-weather-api"

# Read each line in the .env file and set it as a config variable on Heroku
while IFS= read -r line; do
  if [[ $line =~ ^[^#]*= ]]; then
    heroku config:set "$line" -a "$HEROKU_APP_NAME"
  fi
done < .env