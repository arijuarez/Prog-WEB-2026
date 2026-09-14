test: 
	docker compose down -v
	@echo "Volumnes viejos borrados..."
	sqlc generate
	docker compose build
	@echo "Imagen construida..."
	@echo "Corriendo tests"
	docker compose run --rm api go test -v
	@echo "Testeos realizados exitosamente"
	docker compose down -v
	