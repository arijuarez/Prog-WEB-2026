test: 
	docker compose down -v
	@echo "Volumnes viejos borrados..."
	sqlc generate
	docker compose build
	@echo "Imagen construida..."
	docker compose up -d --remove-orphans
	@echo "Contenedor activado..."
	@echo "Corriendo tests"
	docker compose run api go test -v
	@echo "Testeos realizados exitosamente"
	docker compose down -v
	