package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func parseVector(s string) ([]float32, error) {
	// Remove brackets if present
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")
	// Split by spaces
	parts := strings.Fields(s)
	vec := make([]float32, len(parts))
	for i, p := range parts {
		val, err := strconv.ParseFloat(strings.TrimSpace(p), 32)
		if err != nil {
			return nil, fmt.Errorf("invalid number %q: %v", p, err)
		}
		vec[i] = float32(val)
	}
	return vec, nil
}

// usage example:
// go run main.go -v1 [143 2056 2599 2184 4909 481 332 618 354 191] -v2 [143 1290 2540 2056 2599 2202 1215 373 3127 2184]
func main() {
	vec1Str := flag.String("v1", "", "First vector, space-separated in brackets (e.g. [1 2 3])")
	vec2Str := flag.String("v2", "", "Second vector, space-separated in brackets (e.g. [4 5 6])")
	flag.Parse()

	if *vec1Str == "" || *vec2Str == "" {
		fmt.Fprintln(os.Stderr, "Usage: main -v1 \"[1 2 3]\" -v2 \"[4 5 6]\"")
		os.Exit(1)
	}

	v1, err := parseVector(*vec1Str)
	if err != nil {
		log.Fatalf("Error parsing v1: %v", err)
	}
	v2, err := parseVector(*vec2Str)
	if err != nil {
		log.Fatalf("Error parsing v2: %v", err)
	}

	cos := CosineSimilarity(v1, v2)
	fmt.Printf("Cosine similarity: %f\n", cos)
}

func CosineSimilarity(vec1, vec2 []float32) float64 {
	if len(vec1) != len(vec2) {
		return 0 // Return 0 if the vectors have different lengths
	}

	var dotProduct, magnitude1, magnitude2 float64
	for i := range vec1 {
		dotProduct += float64(vec1[i] * vec2[i])
		magnitude1 += float64(vec1[i] * vec1[i])
		magnitude2 += float64(vec2[i] * vec2[i])
	}

	if magnitude1 == 0 || magnitude2 == 0 {
		return 0 // Avoid division by zero
	}

	return dotProduct / (math.Sqrt(magnitude1) * math.Sqrt(magnitude2))
}
