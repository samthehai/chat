#############
# VARIABLES #
#############
SQL_MIGRATE_DB_CONFIG = config/database/dbconfig.yml
SQL_MIGRATE_DB_LABEL = development

DATASOURCE_HOST = 0.0.0.0
DATASOURCE_USER = chat
DATASOURCE_PASS = chat
DATASOURCE_PORT = 5432
DATASOURCE_DATABASE = chat
SEED_DATA_PATH = config/database/testdata/seed.sql

#############
# COMMANDS  #
#############

download:
	@echo Download go.mod dependencies
	@go mod download

install-tools: download
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

migrate-new:
ifdef name
	@sql-migrate new -config=${SQL_MIGRATE_DB_CONFIG} -env=${SQL_MIGRATE_DB_LABEL} ${name}
else
	@echo "please specify name=<value>"
endif

clean-postgres:
	@PGPASSWORD=${DATASOURCE_PASS} dropdb \
		-h ${DATASOURCE_HOST} \
		-p ${DATASOURCE_PORT} \
		-U ${DATASOURCE_USER} \
		-e \
		${DATASOURCE_DATABASE} && \
		sleep 1 && \
		PGPASSWORD=${DATASOURCE_PASS} \
		createdb \
		-h ${DATASOURCE_HOST} \
		-p ${DATASOURCE_PORT} \
		-U ${DATASOURCE_USER} \
		-e \
		${DATASOURCE_DATABASE} \

migrate-up:
	@sql-migrate up -config=${SQL_MIGRATE_DB_CONFIG} -env=${SQL_MIGRATE_DB_LABEL} ${DATASOURCE_DATABASE}

migrate-down:
	@sql-migrate down -config=${SQL_MIGRATE_DB_CONFIG} -env=${SQL_MIGRATE_DB_LABEL}

seed: clean-postgres migrate-up
	@PGPASSWORD=${DATASOURCE_PASS} \
	  psql \
		-h ${DATASOURCE_HOST} \
		-p ${DATASOURCE_PORT} \
		-U ${DATASOURCE_USER} \
		-w ${DATASOURCE_PASS} \
		-e \
    -f ${SEED_DATA_PATH}

lint:
	@go mod tidy
	@golangci-lint run

lintfix:
	@golangci-lint run --fix

gqlgen:
	@echo gqlgen generating...
	rm -rf internal/interfaces/graph/generated/generated.go
	@go run github.com/99designs/gqlgen

serve:
	@go run cmd/main.go

start-redis:
	redis-server /usr/local/etc/redis.conf
