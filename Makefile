.PHONY: vendor
vendor:
	go mod tidy && go mod vendor

# copy config.toml.example so we can adjust without affecting git file change
.PHONY: config
config:
	cp config.toml.example config.toml

.PHONY: server/start
server/start:
	docker-compose -f docker-compose.yml up --build