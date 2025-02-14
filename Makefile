.PHONY: memory postgres

memory:
	@echo "Запуск с in-memory хранилищем..."
	@export STORAGE_TYPE=redis; \
	docker-compose up -d redis; \
	docker-compose up -d url_shortener

postgres:
	@echo "Запуск с PostgreSQL..."
	@export STORAGE_TYPE=postgres; \
	docker-compose up -d postgres; \
	docker-compose up -d url_shortener