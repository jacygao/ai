package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/jacygao/ai/llm/ollama"
)

// BM25 parameters
const (
	k1 = 1.5  // Controls term frequency saturation
	b  = 0.75 // Controls document length normalization
)

var top = 3
var skips = make(map[string]bool)

func cleanText(text string) string {
	words := strings.Fields(text) // Split text into words
	var filteredWords []string

	for _, word := range words {
		if skips[strings.ToLower(word)] {
			continue
		}
		filteredWords = append(filteredWords, word)
	}

	return strings.Join(filteredWords, " ")
}

// BM25 Index structure
type BM25Index struct {
	Corpus      []string
	InvertedIdx map[string]map[int]int // term -> {docID -> term frequency}
	DocLengths  map[int]int            // docID -> document length
	AvgDL       float64
	N           int // Total number of documents
}

// Tokenizer function
func tokenize(text string) []string {
	return strings.Fields(strings.ToLower(text)) // Simple whitespace-based tokenizer
}

// Compute IDF (Inverse Document Frequency)
func computeIDF(N, df int) float64 {
	return math.Log((float64(N) - float64(df) + 0.5) / (float64(df) + 0.5))
}

// Build BM25 Index
func BuildBM25Index(corpus []string) BM25Index {
	invertedIdx := make(map[string]map[int]int)
	docLengths := make(map[int]int)
	totalLength := 0

	for docID, doc := range corpus {
		doc = cleanText(doc)
		tokens := tokenize(doc)
		docLengths[docID] = len(tokens)
		totalLength += len(tokens)

		for _, token := range tokens {
			if _, exists := invertedIdx[token]; !exists {
				invertedIdx[token] = make(map[int]int)
			}
			invertedIdx[token][docID]++
		}
	}

	avgDL := float64(totalLength) / float64(len(corpus))

	return BM25Index{
		Corpus:      corpus,
		InvertedIdx: invertedIdx,
		DocLengths:  docLengths,
		AvgDL:       avgDL,
		N:           len(corpus),
	}
}

// Compute BM25 Score for a document
func BM25Score(query []string, docID int, index BM25Index) float64 {
	score := 0.0
	docLength := index.DocLengths[docID]

	for _, term := range query {
		termFreq := index.InvertedIdx[term][docID]
		docFreq := len(index.InvertedIdx[term])
		idf := computeIDF(index.N, docFreq)

		numerator := float64(termFreq) * (k1 + 1)
		denominator := float64(termFreq) + k1*(1-b+b*float64(docLength)/index.AvgDL)

		score += idf * (numerator / denominator)
	}

	return score
}

type BM25Result struct {
	DocID int
	Text  string
	Score float64
}

// Example usage
func main() {
	// "a", "is", "are", "and"
	skips["a"] = true
	skips["is"] = true
	skips["are"] = true
	skips["and"] = true

	corpus := []string{
		"Jacy is a software engineer.",
		"Charlotte will become a lawyer in a few weeks after her official admission.",
		"Jacy married Charlotte 3 years ago.",
		"Matt and Charlotte are best friends.",
		"Matt and Charlotte grow up together and they have known each other since they were 3 years old.",
		"Christiane is Charlotte's mum.",
	}

	index := BuildBM25Index(corpus)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter your question : ")
		originalQuery, _ := reader.ReadString('\n')
		if originalQuery == "exit\n" {
			fmt.Println("Exiting...")
			break
		}
		originalQuery = cleanText(originalQuery)
		query := tokenize(originalQuery)

		foundDocs := []string{}
		var results []BM25Result

		for docID := range index.Corpus {
			score := BM25Score(query, docID, index)
			fmt.Printf("BM25 Score for Document %d: %.4f\n", docID, score)
			results = append(results, BM25Result{docID, corpus[docID], score})
		}

		sort.Slice(results, func(a int, b int) bool {
			return results[a].Score > results[b].Score
		})

		for i := 0; i < top; i++ {
			foundDocs = append(foundDocs, results[i].Text)
		}
		ollama.Chat(originalQuery, foundDocs)
		fmt.Println("")
	}
}
