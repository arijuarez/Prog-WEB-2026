.PHONY: test-hurl

test-hurl:
	@echo "Eliminando volúmenes y contenedores previos..."
	@docker compose down -v
	@echo "Generando código de base de datos con sqlc..."
	@sqlc generate
	@echo "Construyendo y levantando contenedores..."
	@docker compose up -d --build --wait
	@echo "Ejecutando pruebas con Hurl..."
	hurl --test --jobs 1 Handlers/hurl_test/*.hurl || (docker compose down -v && exit 1)
	@echo "Pruebas completadas exitosamente."
	@echo "Limpiando entorno Docker..."
	@docker compose down -v