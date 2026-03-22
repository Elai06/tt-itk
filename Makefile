migrate-create:
	goose -dir migrations create $(name) sql

migrate-up:
	goose -dir migrations postgres up

migrate-down:
	goose -dir migrations public down