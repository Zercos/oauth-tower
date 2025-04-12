#!/bin/bash
export DB_PATH="../../test.db"
export SECRET_KEY="secret"
go test -v ./...