## The App

| RESOURCE | HTTP METHOD | ROUTE                | DESCRIPTION              |
| -------- | ----------- | -------------------- | ------------------------ |
| auth     | POST        | /api/v1/auth/signup  | Register a new user      |
| auth     | POST        | /api/v1/auth/login   | Login user               |
| auth     | GET         | /api/v1/auth/refresh | Request a new token pair |
| auth     | GET         | /api/v1/auth/logout  | Logout user              |
| users    | GET         | /api/v1/users/me     | Get me                   |
| users    | GET         | /api/v1/users/:id    | Find user by id          |
| users    | GET         | /api/v1/users        | Find users               |
| users    | PATCH       | /api/v1/users        | Update user              |
| users    | DELETE      | /api/v1/users        | Delete user              |

#### TOOLS:

- golang
- gin
- postgres, pgx
- zap

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
