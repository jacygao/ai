package redis

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password docs
		DB:       0,  // use default DB
		Protocol: 2,
	})

	ctx := context.Background()
	rdb.FTDropIndexWithArgs(ctx,
		"vector_idx",
		&redis.FTDropIndexOptions{
			DeleteDocs: true,
		},
	)
	_, err := rdb.FTCreate(ctx,
		"vector_idx",
		&redis.FTCreateOptions{
			OnHash: true,
			Prefix: []any{"docs:"},
		},
		&redis.FieldSchema{
			FieldName: "content",
			FieldType: redis.SearchFieldTypeText,
		},
		&redis.FieldSchema{
			FieldName: "embedding",
			FieldType: redis.SearchFieldTypeVector,
			VectorArgs: &redis.FTVectorArgs{
				HNSWOptions: &redis.FTHNSWOptions{
					Dim:            384,
					DistanceMetric: "COSINE",
					Type:           "FLOAT32",
				},
			},
		},
	).Result()

	if err != nil {
		panic(err)
	}

	return &RedisClient{
		client: rdb,
	}
}

func (rdb *RedisClient) Set(key string, content string, embedding []float32) {
	ctx := context.Background()
	// Convert embedding to byte array
	byteEmbedding := floatsToBytes(embedding)
	fmt.Println(content)

	// Store in Redis
	_, err := rdb.client.HSet(
		ctx,
		fmt.Sprintf("docs:%v", key),
		map[string]any{
			"content":   content,
			"embedding": byteEmbedding,
		}).Result()

	if err != nil {
		fmt.Println("Error storing embedding:", err)
		return
	}

	fmt.Println("Embedding stored successfully!")
}

func floatsToBytes(fs []float32) []byte {
	buf := make([]byte, len(fs)*4)

	for i, f := range fs {
		u := math.Float32bits(f)
		binary.NativeEndian.PutUint32(buf[i*4:], u)
	}

	return buf
}

func (rdb *RedisClient) SearchVector(ctx context.Context, queryVector []float32) []string {
	// Convert query vector to binary
	queryBytes := floatsToBytes(queryVector)
	// fmt.Println(queryBytes)
	// Execute Redis search query
	results, err := rdb.client.FTSearchWithArgs(
		ctx,
		"vector_idx",
		"*=>[KNN 5 @embedding $vec AS vector_distance]",
		&redis.FTSearchOptions{
			Return: []redis.FTSearchReturn{
				{FieldName: "vector_distance"},
				{FieldName: "content"},
			},
			DialectVersion: 2,
			Params: map[string]any{
				"vec": queryBytes,
			},
		}).Result()

	if err != nil {
		fmt.Println("Error running search:", err)
		return nil
	}

	found := []string{}

	for _, doc := range results.Docs {
		found = append(found, doc.Fields["content"])
	}
	return found
}
