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

	// _, err := rdb.FTCreate(ctx,
	// 	"vector_idx",
	// 	&redis.FTCreateOptions{
	// 		OnHash: true,
	// 		Prefix: []any{"doc:"},
	// 	},
	// 	&redis.FieldSchema{
	// 		FieldName: "content",
	// 		FieldType: redis.SearchFieldTypeText,
	// 	},
	// 	&redis.FieldSchema{
	// 		FieldName: "genre",
	// 		FieldType: redis.SearchFieldTypeTag,
	// 	},
	// 	&redis.FieldSchema{
	// 		FieldName: "embedding",
	// 		FieldType: redis.SearchFieldTypeVector,
	// 		VectorArgs: &redis.FTVectorArgs{
	// 			HNSWOptions: &redis.FTHNSWOptions{
	// 				Dim:            384,
	// 				DistanceMetric: "L2",
	// 				Type:           "FLOAT32",
	// 			},
	// 		},
	// 	},
	// ).Result()

	// if err != nil {
	// 	panic(err)
	// }

	return &RedisClient{
		client: rdb,
	}
}

func (rdb *RedisClient) Set(key string, embedding []float64) {
	ctx := context.Background()
	// Convert embedding to byte array
	byteEmbedding := floatsToBytes(embedding)

	// Store in Redis
	rdb.client.HSet(ctx, key, "doc_embedding", byteEmbedding)
	fmt.Println("Embedding stored successfully!")
}

func floatsToBytes(fs []float64) []byte {
	buf := make([]byte, len(fs)*8)
	for i, f := range fs {
		binary.LittleEndian.PutUint64(buf[i*4:], math.Float64bits(f))
	}
	return buf
}

func (rdb *RedisClient) SearchVector(ctx context.Context, queryVector []float64) {
	// Convert query vector to binary
	queryBytes := floatsToBytes(queryVector)

	// Execute Redis search query
	searchCmd := rdb.client.Do(ctx, "FT.SEARCH", "my_index", "(*)=>[KNN 3 @doc_embedding $query_vector]",
		"PARAMS", "2", "query_vector", queryBytes, "DIALECT", "2")

	// Process results
	results, err := searchCmd.Result()
	if err != nil {
		fmt.Println("Error running search:", err)
		return
	}
	fmt.Println("Search results:", results)
}
