.PHOOY: goose-status goose-up goose-down reset fix 

goose-status:
	goose postgres "host=localhost port=5432 user=baloo password=junglebook dbname=lenslocked sslmode=disable" status

goose-up:
	goose postgres "host=localhost port=5432 user=baloo password=junglebook dbname=lenslocked sslmode=disable" up

goose-down:
	goose postgres "host=localhost port=5432 user=baloo password=junglebook dbname=lenslocked sslmode=disable" down

goose-reset:
	goose postgres "host=localhost port=5432 user=baloo password=junglebook dbname=lenslocked sslmode=disable" reset


fix:
	goose fix