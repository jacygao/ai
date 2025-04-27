package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Example embedding (1536-dimensional vector)
	embedding := []float32{0.1, 0.2, 0.3 /* ... */, 0.9}

	// Convert embedding to byte array
	byteEmbedding := floatsToBytes(embedding)

	// Store in Redis
	rdb.HSet(ctx, "docs:1", "doc_embedding", byteEmbedding)
	fmt.Println("Embedding stored successfully!")
}

func floatsToBytes(fs []float32) []byte {
	buf := make([]byte, len(fs)*4)
	for i, f := range fs {
		binary.LittleEndian.PutUint32(buf[i*4:], math.Float32bits(f))
	}
	return buf
}
