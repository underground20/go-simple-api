up:
	docker compose up -d

generate-swagger:
	swagger generate spec -o ./swagger.json --scan-models