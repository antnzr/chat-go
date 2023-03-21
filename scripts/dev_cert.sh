#!/bin/bash

set -e

# Need to add 127.0.0.1 home.domain to /etc/hosts
echo "remove old certificates"
sudo rm -rf ./tools/home.domain/*.pem
sudo rm -rf ./tools/*.pem

echo "generate new certificates"
docker run --rm --user $(id -u):$(id -g) -it -v $(pwd)/tools:/output ryantk/minica --domains home.domain

echo "allow read .pem files"
sudo find ./tools -name "*.pem" -exec chmod 644 {} \;

echo "add to ca-certificates"
sudo cp tools/minica.pem /usr/local/share/ca-certificates/minica.crt

echo "give permissions"
sudo chmod 644 /usr/local/share/ca-certificates/minica.crt

echo "update certificates"
sudo update-ca-certificates

### Script installs minica.pem to certificate trust store of applications using NSS
### (e.g. Firefox, Thunderbird, Chromium)
###
### Requirement: apt install libnss3-tools
###
### CA file to install (customize!)
### Retrieve Certname: openssl x509 -noout -subject -in minica.pem
###
certfile="tools/minica.pem"
certname="minica_root_ca"
###
### For cert9 (SQL)
###
for certDB in $(find ~/ -name "cert9.db"); do
    certdir=$(dirname ${certDB})
    echo "certdir --> $certdir"
    certutil -A -n "${certname}" -t "TCu,Cu,Tu" -i ${certfile} -d sql:${certdir}
done
