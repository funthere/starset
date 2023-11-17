# starset

Simple E-commerce Microservice written in Go.

Stack: Go, PostgreSQL, Jwt, REST API


## Getting Started

Make sure to set the .env file (see: .env.example).

Clone the project and run this command to start/stop the services:

```bash
make up
```
```bash
make down
```

### Documentation
[https://documenter.getpostman.com/view/718196/2s9YXh5i4e](https://documenter.getpostman.com/view/718196/2s9YXh5i4e)

After register and login the user, put the token in the request header for product and order endpoint (see the example on the docs)


## Updates
- Nov 15: Change DB to PostgreSQL [here](https://github.com/funthere/starset/pull/6)
- Nov 15: Dockerize the services [here](https://github.com/funthere/starset/pull/7)
- Nov 16: Separate DB for each services [here](https://github.com/funthere/starset/commit/3c4586a2ec50ecf2a7ebffc57319a0e4c44ba5ec)
