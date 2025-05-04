package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/jacygao/ai/llm/ollama"
	"github.com/jacygao/ai/vector/redis"
)

func BuildVectors(redisClient *redis.RedisClient, data []string) {
	counter := 1
	for _, s := range data {
		vector, err := ollama.Embed(s)
		if err != nil {
			fmt.Printf("Error storing vector %s \n", s)
		}
		redisClient.Set(strconv.Itoa(counter), s, ConvertFloat64ToFloat32(vector))
		counter++
	}
}

func ConvertFloat64ToFloat32(input []float64) []float32 {
	// Create a slice with the same length as the input
	output := make([]float32, len(input))

	// Convert each float64 value to float32
	for i, v := range input {
		output[i] = float32(v)
	}

	return output
}

func SearchVector(redisClient *redis.RedisClient, query string) []string {
	vector, err := ollama.Embed(query)
	if err != nil {
		fmt.Printf("Error searching vector %s \n", query)
	}
	return redisClient.SearchVector(context.Background(), ConvertFloat64ToFloat32(vector))
}

// Example usage
func main() {
	corpus := []string{
		"Jacy is a software engineer.",
		"Charlotte will become a lawyer in a few weeks after her official admission.",
		"Jacy married Charlotte 3 years ago.",
		"Matt and Charlotte are best friends.",
		"Matt and Charlotte grow up together and they have known each other since they were 3 years old.",
		"Christiane is Charlotte's mum.",
	}

	redisClient := redis.NewRedisClient()
	BuildVectors(redisClient, corpus)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter your question : ")
		originalQuery, _ := reader.ReadString('\n')
		if originalQuery == "exit\n" {
			fmt.Println("Exiting...")
			break
		}

		found := SearchVector(redisClient, originalQuery)

		ollama.Chat(originalQuery, found)
		fmt.Println("")
	}
}
