pgvector:
	@docker run -d \
    	--name pgvector \
		-e POSTGRES_USER=admin \
        -e POSTGRES_PASSWORD=admin \
        -e POSTGRES_DB=vectordb \
        -v pgvector_data:/var/lib/postgresql/data \
        -p 5432:5432 \
        pgvector/pgvector:pg17
psql:
	@psql -h localhost -U admin -d vectordb