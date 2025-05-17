#!/bin/sh
migrate -path /app/sql/migrations -database sqlite:///app/database/flatos.sqlite up
/app/main