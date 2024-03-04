.PHOOY: run connect-db

run:
	go run main.go

connect-db:
	docker compose exec -it db psql -U baloo -d lenslocked
