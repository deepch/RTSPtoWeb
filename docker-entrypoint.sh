#!/bin/bash

# Write the config string to the config file if CONFIG_STRING is provided
if [ ! -z "$CONFIG_STRING" ]; then
    echo "$CONFIG_STRING" > /config/config.json
    echo "Config written to /config/config.json"
fi

# Execute the main command
exec "$@" 