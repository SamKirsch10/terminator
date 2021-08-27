#!/bin/bash


usage() {
    echo "Usage:"
    echo "$(basename $0) file.pfx" 
}

if ! command -v openssl &> /dev/null
then
    echo "openssl could not be found!"
    exit
fi


if [ -z $1 ]; then
    usage
    exit 1
fi

PFX_IN="$1"
PEM_OUT="$(echo $PFX_IN | sed 's/.pfx//')"

openssl pkcs12 -in $PFX_IN -out "${PEM_OUT}.pem" -nokeys
openssl pkcs12 -in $PFX_IN -out "${PEM_OUT}.withkey.pem"
openssl rsa -in "${PEM_OUT}.withkey.pem" -out "${PEM_OUT}.key"
rm "${PEM_OUT}.withkey.pem"