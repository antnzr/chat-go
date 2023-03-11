#!/bin/bash
# Run this from the root of the project ./scripts/create_rsa_keys_env_var_sh
# Script creates rsa keys for access and refresh tokens in /tmp dir
# and places base64 output to env_file as environment variables
# ACCESS_TOKEN_PRIVATE_KEY=LS0tLS1CRUdJTiBSU0EgUFJJVkFURS.....
# ACCESS_TOKEN_PUBLIC_KEY=LS0tLS1CRUdJTi......

set -e

temp_dir=/tmp
access_jwt_key=access_jwtRS256.key
refresh_jwt_key=refresh_jwtRS256.key
env_file=.env

# create $env_file if doesn't exist
if [ ! -f "$env_file" ]; then
  echo "create '$env_file'"
  touch $env_file
fi

################################################################
echo "generate token pair for access token jwt"
ssh-keygen -t rsa -b 4096 -m PEM -f $temp_dir/$access_jwt_key -N "" <<<y >/dev/null
openssl rsa -in $temp_dir/$access_jwt_key -pubout -outform PEM -out "$temp_dir/$access_jwt_key.pub"

accessPrivateKey=$(base64 -w 0 $temp_dir/$access_jwt_key)
accessPublicKey=$(base64 -w 0 "$temp_dir/$access_jwt_key.pub")

echo "append access token environment vars to env_file"
echo "ACCESS_TOKEN_PRIVATE_KEY=$accessPrivateKey" >>$env_file
echo "ACCESS_TOKEN_PUBLIC_KEY=$accessPublicKey" >>$env_file
################################################################
echo "generate token pair for refresh token jwt"
ssh-keygen -t rsa -b 4096 -m PEM -f $temp_dir/$refresh_jwt_key -N "" <<<y >/dev/null
openssl rsa -in $temp_dir/$refresh_jwt_key -pubout -outform PEM -out "$temp_dir/$refresh_jwt_key.pub"

refreshPrivateKey=$(base64 -w 0 $temp_dir/$refresh_jwt_key)
refreshPublicKey=$(base64 -w 0 "$temp_dir/$refresh_jwt_key.pub")

echo "append refresh token environment vars to env_file"
echo "REFRESH_TOKEN_PRIVATE_KEY=$refreshPrivateKey" >>$env_file
echo "REFRESH_TOKEN_PUBLIC_KEY=$refreshPublicKey" >>$env_file
################################################################
