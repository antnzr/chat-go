## The App

| RESOURCE | HTTP METHOD | ROUTE                | DESCRIPTION                |
| -------- | ----------- | -------------------- | -------------------------- |
| users    | GET         | /api/v1/users/me     | Get me                     |
| users    | PATCH       | /api/v1/users        | Update user                |
| auth     | POST        | /api/v1/auth/signup  | Register a new user        |
| auth     | POST        | /api/v1/auth/login   | Login user                 |
| auth     | GET         | /api/v1/auth/refresh | Request the new token pair |
| auth     | GET         | /api/v1/auth/logout  | Logout user                |

#### TOOLS:

- golang
- gin
- postgresql, pgx

#### Generate rsa key for jwt, and paste them to .env or use a script in a `scripts` dir

```sh
ssh-keygen -t rsa -b 4096 -m PEM -f access_jwtRS256.key \
openssl rsa -in access_jwtRS256.key -pubout -outform PEM -out access_jwtRS256.key.pub \
base64 -w 0 access_jwtRS256.key \
base64 -w 0 access_jwtRS256.key.pub \
```

```sh
ssh-keygen -t rsa -b 4096 -m PEM -f refresh_jwtRS256.key \
openssl rsa -in refresh_jwtRS256.key -pubout -outform PEM -out refresh_jwtRS256.key.pub \
base64 -w 0 refresh_jwtRS256.key \
base64 -w 0 refresh_jwtRS256.key.pub \
```
