importdb:
	docker exec -i postgres-db psql -U root -d master-golang < ./backupdb-master-golang.sql
exportdb:
	docker exec -i postgres-db pg_dump -U root -d master-golang > ./backupdb-master-golang.sql
server:
	go run .