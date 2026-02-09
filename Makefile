init: copy-env build

build:
	docker compose up -d --build

up:
	docker compose up -d

copy-env:
	cp .env .env.local