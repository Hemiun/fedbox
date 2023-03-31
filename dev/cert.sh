#!/bin/bash

mkdir cert
cd cert

openssl genrsa -out server.key 4096
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650 -subj "/C=PH/ST=Manila/L=Manila/CN=localhost:30000/emailAddress=mail@gmail.com"

