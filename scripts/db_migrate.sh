#!/bin/bash

set -e

export_current_user() {
  if [ -z "$UID" ]; then export UID="$(id -u)"; fi
  if [ -z "$GID" ]; then export GID="$(id -g)"; fi
  if [ -z "$CURRENT_UID" ]; then export CURRENT_UID=$UID:$GID; fi
}

export_db_url() {
  if [ -z "$DATABASE_URL" ]; then export $(grep -v '^#' $(pwd)/.env | xargs); fi
}

export_db_url
export_current_user

docker run --rm -it \
  --user $CURRENT_UID \
  --network=host \
  -v $(pwd)/internal/app/db:/db \
  amacneil/dbmate -d "/db/migrations" -u $DATABASE_URL $@
