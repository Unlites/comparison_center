migrate_up:
	docker exec -e MIGRATE_OPERATION=up -it comparison_center_app /bin/migrate

migrate_down:
	docker exec -e MIGRATE_OPERATION=down -it comparison_center_app /bin/migrate

run:
	docker-compose up -d