#!/bin/sh
set -e

echo "Generating envoy.yaml config file..."
cat /tmpl/envoy.yaml.tmpl | envsubst \$STATISTICO_RATINGS_HOST,\$STATISTICO_RATINGS_PORT > /etc/envoy/envoy.yaml

echo "Starting Envoy..."
/usr/local/bin/envoy -c /etc/envoy/envoy.yaml
