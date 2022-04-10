#!/usr/bin/env bash

set -e

function generate_docs() {
  SERVICE=$1
  DOCSPATH=./docs/services/$SERVICE
  echo -en "Creating docs for \e[96m$SERVICE\e[0m... "
  mkdir -p "$DOCSPATH"
  go run ./shoutrrr docs -f markdown "$SERVICE" > "$DOCSPATH"/config.md
  if [ $? ]; then
    echo -e "Done!"
  fi
}

if [[ -n "$1" ]]; then
  generate_docs "$1"
  exit 0
fi

for S in ./pkg/services/*; do
  SERVICE=$(basename "$S")
  if [[ "$SERVICE" == "standard" ]] || [[ "$SERVICE" == "xmpp" ]]  || [[ -f "$S" ]]; then
    continue
  fi
  generate_docs "$SERVICE"
done
