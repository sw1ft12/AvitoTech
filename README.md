exec queries.sql in pg
docker build .
docker run -e SERVER_ADDRESS= -e POSTGRES_CONN= --network= <container_id>