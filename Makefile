.PHOOY: run connect-db up down

run:
	go run main.go

up:
	docker compose up -d 

down:
	docker compose down

connect-db:
	docker compose exec -it db psql -U baloo -d lenslocked
