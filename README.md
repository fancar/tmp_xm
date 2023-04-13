# XM Golang Exercise - v22.0.0

## Setup and run
- config file: config.toml.
- config gile example: docker exec -it xm_app_1 xm configfile
- to run execute: docker compose up -d OR docker-compose up -d

## Swagger link
-  http://localhost:8085/api

## default creds
- login with admin/admin to get your jwt token
- then put the token in form on top in swagger

- login link:
	 http://localhost:8085/api#!/CompanyService/CompanyService_Login
	 body:
		 {
		  "password": "admin",
		  "user": "admin"
		 }


