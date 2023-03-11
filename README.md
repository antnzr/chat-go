## Chat App

POST /signup
POST /login
GET /me
GET /dialogs

* ws events
send_message
receive_message

Tools:
golang,
postgresql, pgx
gin
websockets

## Generate rsa key for jwt, and paste them to .env
```sh
ssh-keygen -t rsa -b 4096 -m PEM -f access_jwtRS256.key
openssl rsa -in access_jwtRS256.key -pubout -outform PEM -out access_jwtRS256.key.pub
base64 -w 0 access_jwtRS256.key
base64 -w 0 access_jwtRS256.key.pub \
```

```sh
ssh-keygen -t rsa -b 4096 -m PEM -f refresh_jwtRS256.key
openssl rsa -in refresh_jwtRS256.key -pubout -outform PEM -out refresh_jwtRS256.key.pub
base64 -w 0 refresh_jwtRS256.key
base64 -w 0 refresh_jwtRS256.key.pub \
```
