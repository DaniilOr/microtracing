#!/bin/sh

openssl req -newkey rsa:2048 -nodes -x509 -keyout key_auth.pem -out certificate_auth.pem \
 -subj "/C=RU/ST=Moscow/L=Moscow/O=Development/OU=Dev/CN=auth" \
 -addext "subjectAltName = DNS:auth,IP:0.0.0.0"


openssl req -newkey rsa:2048 -nodes -x509 -keyout key_trans.pem -out certificate_trans.pem \
 -subj "/C=RU/ST=Moscow/L=Moscow/O=Development/OU=Dev/CN=transactions" \
 -addext "subjectAltName = DNS:transactions,IP:0.0.0.0"
# 1. Generate CA's private key and self-signed certificate
#openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-key.pem -out ca-cert.pem -subj "/C=RU/ST=Moscow/L=Moscow/O=Dev/OU=Dev/CN=netology.local" -addext "subjectAltName=DNS:netology.local,IP:0.0.0.0"

#echo "CA's self-signed certificate"
#openssl x509 -in ca-cert.pem -noout -text

# 2. Generate web server's private key and certificate signing request (CSR)
#openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj  "/C=RU/ST=Moscow/L=Moscow/O=Dev/OU=Dev/CN=netology.local" -addext "subjectAltName=DNS:netology.local,IP:0.0.0.0"

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
#openssl x509 -req -in server-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem

#echo "Server's signed certificate"
#openssl x509 -in server-cert.pem -noout -text
