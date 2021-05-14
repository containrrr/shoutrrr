#!/usr/bin/env bash

for S in ./pkg/services/*; do
  SERVICE=$(basename "$S")
  if [[ "$SERVICE" == "standard" ]] || [[ -f "$S" ]]; then
    continue
  fi
  DOCSPATH=./docs/services/$SERVICE
  echo -en "Creating docs for \e[96m$SERVICE\e[0m... "
  mkdir -p "$DOCSPATH"
  go run ./cli docs -f markdown "$SERVICE" > "$DOCSPATH"/config.md
  if [ $? ]; then
    echo -e "Done!"
  fi
done
