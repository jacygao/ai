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
		redisClient.Set("docs:"+strconv.Itoa(counter), vector)
		counter++
	}
}

func SearchVector(redisClient *redis.RedisClient, query string) {
	vector, err := ollama.Embed(query)
	if err != nil {
		fmt.Printf("Error searching vector %s \n", query)
	}
	redisClient.SearchVector(context.Background(), vector)
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

		SearchVector(redisClient, originalQuery)

		// sort.Slice(results, func(a int, b int) bool {
		// 	return results[a].Score > results[b].Score
		// })

		// for i := 0; i < top; i++ {
		// 	foundDocs = append(foundDocs, results[i].Text)
		// }
		// ollama.Chat(originalQuery, foundDocs)
		fmt.Println("")
	}
}
