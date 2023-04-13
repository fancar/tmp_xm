# XM Golang Exercise - v22.0.0

## Setup and run
- config file: config.toml.
- all config options: docker exec -it xm_app_1 xm configfile
- to execute: docker compose up -d OR docker-compose up -d

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

- Create payload example:
	{
	  "Company": {
	    "name" : "company name",
	    "employeesCnt": 111,
	    "description": "here is some description",
	    "id": "65fb89c1-d145-41ca-b9b5-47ae7f79a70c",
	    "registered": true,
	    "type": "NonProfit"
	  }
	}

