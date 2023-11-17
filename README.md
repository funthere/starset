# starset

Simple E-commerce Microservice written in Go.

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
- Nov:
  1. [Change DB to PostgreSQL](https://github.com/funthere/starset/pull/6)
  2. [Dockerize the services](https://github.com/funthere/starset/pull/7)
  3. [Separate DB for each services](https://github.com/funthere/starset/commit/3c4586a2ec50ecf2a7ebffc57319a0e4c44ba5ec)
