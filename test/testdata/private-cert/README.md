This folder contains a self-signed certificate for testing purposes ONLY.

The files were created like this:
```
openssl req -config openssl.conf -new -newkey rsa:2048 -days 825 -nodes -x509 \
    -keyout tls.key -out tls.crt -extensions req_ext
```

Expiry of 825 days is due to an MacOS/Go issue, see https://myupbeat.wordpress.com/2022/09/09/self-signed-certificates-not-standards-compliant/.
