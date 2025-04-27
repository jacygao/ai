package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

func floats32ToString(floats []float32) string {
	strVals := make([]string, len(floats))
	for i, val := range floats {
		// Format each float into a string
		strVals[i] = fmt.Sprintf("%f", val)
	}

	// Join them with comma + space
	joined := strings.Join(strVals, ", ")

	// pgvector requires bracketed notation for vector input, e.g. [0.1, 0.2, 0.3]
	return "[" + joined + "]"
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create the connection pool
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	// 1. Ensure pgvector extension is enabled
	_, err = dbpool.Exec(context.Background(), "CREATE EXTENSION IF NOT EXISTS vector;")
	if err != nil {
		log.Fatalf("Failed to create extension: %v\n", err)
		os.Exit(1)
	}

	// 2. Create table (if not existing)
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS documents (
        id SERIAL PRIMARY KEY,
        content TEXT,
        embedding vector(1536)
    );
    `
	_, err = dbpool.Exec(context.Background(), createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v\n", err)
	}

	// 3. Create index (if not existing)
	createIndexSQL := `
    CREATE INDEX IF NOT EXISTS documents_embedding_idx
    ON documents USING ivfflat (embedding vector_l2_ops) WITH (lists = 100);
    `
	_, err = dbpool.Exec(context.Background(), createIndexSQL)
	if err != nil {
		log.Fatalf("Failed to create index: %v\n", err)
	}

	// 4. Initialize OpenAI client
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY is not set")
	}
	openaiClient := openai.NewClient(apiKey)

	// 5. Insert sample documents
	docs := []string{
		"PostgreSQL is an advanced open-source relational database.",
		"OpenAI provides GPT-based models to generate text embeddings.",
		"pgvector allows storing embeddings in a Postgres database.",
	}

	for _, doc := range docs {
		err = insertDocument(context.Background(), dbpool, openaiClient, doc)
		if err != nil {
			log.Printf("Failed to insert document '%s': %v\n", doc, err)
		}
	}

	// 6. Query for similarity
	queryText := "How to store embeddings in Postgres?"
	similarDocs, err := searchSimilarDocuments(context.Background(), dbpool, openaiClient, queryText, 5)
	if err != nil {
		log.Fatalf("Search failed: %v\n", err)
	}

	fmt.Println("=== Most Similar Documents ===")
	for _, doc := range similarDocs {
		fmt.Printf("- %s\n", doc)
	}
}

// insertDocument generates an embedding for `content` using the OpenAI API
// and inserts it into the documents table.
func insertDocument(ctx context.Context, dbpool *pgxpool.Pool, client *openai.Client, content string) error {
	// 1) Get embedding from OpenAI
	embedResp, err := client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Model: openai.AdaEmbeddingV2, // "text-embedding-ada-002"
		Input: []string{content},
	})
	if err != nil {
		return fmt.Errorf("CreateEmbeddings API call failed: %w", err)
	}

	// 2) Convert embedding to bracketed string for pgvector
	embedding := embedResp.Data[0].Embedding // []float32
	embeddingStr := floats32ToString(embedding)

	// 3) Insert into PostgreSQL
	insertSQL := `
        INSERT INTO documents (content, embedding)
        VALUES ($1, $2::vector)
    `
	_, err = dbpool.Exec(ctx, insertSQL, content, embeddingStr)
	if err != nil {
		return fmt.Errorf("failed to insert document: %w", err)
	}

	return nil
}

// searchSimilarDocuments takes a user query, gets the embedding, and returns
// the top-k similar documents based on vector similarity.
func searchSimilarDocuments(ctx context.Context, pool *pgxpool.Pool, client *openai.Client, query string, k int) ([]string, error) {
	// 1) Get the embedding for the user’s query via OpenAI
	embedResp, err := client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Model: openai.AdaEmbeddingV2, // "text-embedding-ada-002"
		Input: []string{query},
	})
	if err != nil {
		return nil, fmt.Errorf("CreateEmbeddings API call failed: %w", err)
	}

	// 2) Convert the OpenAI embedding to the bracketed string format for pgvector
	queryEmbedding := embedResp.Data[0].Embedding // []float32
	queryEmbeddingStr := floats32ToString(queryEmbedding)
	// e.g. "[0.123456, 0.789012, ...]"

	// 3) Build the SELECT statement that orders by vector similarity
	selectSQL := fmt.Sprintf(`
        SELECT content
        FROM documents
        ORDER BY embedding <-> '%s'::vector
        LIMIT %d;
    `, queryEmbeddingStr, k)

	// 4) Run the query
	rows, err := pool.Query(ctx, selectSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to query documents: %w", err)
	}
	defer rows.Close()

	// 5) Read the matching documents
	var contents []string
	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		contents = append(contents, content)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return contents, nil
}
