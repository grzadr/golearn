package parsing

import (
	"encoding/json"
	"fmt"
	"testing"
)

// TestEntry represents the JSON structure we'll generate for testing
type TestEntry struct {
	Value   float64  `json:"value"` // Changed from Value to value to match JSON
	Aliases []string `json:"aliases"`
}

// generateTestData creates a map of TestEntry with specified size
func generateTestData(size int) []byte {
	entries := make(map[string]TestEntry, size)
	for i := 1; i < size+1; i++ {
		entries[fmt.Sprintf("unit%d", i)] = TestEntry{
			Value: float64(i),
			Aliases: []string{
				fmt.Sprintf("alias1_%d", i),
				fmt.Sprintf("alias2_%d", i),
				fmt.Sprintf("alias3_%d", i),
			},
		}
	}

	data, err := json.Marshal(entries)
	if err != nil {
		panic(err)
	}
	return data
}

// BenchmarkInitialization benchmarks the initialization of both data structures
func BenchmarkInitialization(b *testing.B) {
	sizes := []int{2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}

	for _, size := range sizes {
		testData := generateTestData(size)

		b.Run(fmt.Sprintf("UnitRecords/size_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := newUnitRecords(testData)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkLookup benchmarks the lookup operations for both data structures
func BenchmarkLookup(b *testing.B) {
	sizes := []int{2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}

	for _, size := range sizes {
		testData := generateTestData(size)

		// Create instances once before the benchmark
		mapRecords, err := newUnitRecords(testData)
		if err != nil {
			b.Fatal(err)
		}

		// Test best case (first element)
		b.Run(fmt.Sprintf("UnitRecords/size_%d/best", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = mapRecords.findUnit("alias1_0")
			}
		})

		// Test worst case (last element)
		lastAlias := fmt.Sprintf("alias1_%d", size-1)
		b.Run(fmt.Sprintf("UnitRecords/size_%d/worst", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = mapRecords.findUnit(lastAlias)
			}
		})
	}
}
