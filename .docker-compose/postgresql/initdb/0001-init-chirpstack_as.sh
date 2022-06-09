#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    create role app with login password 'app';
    create database app with owner app;
    create role app_test with login password 'app_test';
    create database app_test with owner app_test;
EOSQL
