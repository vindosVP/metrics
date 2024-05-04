package utils

import "testing"

func BenchmarkSha256Hash(b *testing.B) {
	data := []byte("data")
	key := "secret-key"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Sha256Hash(data, key)
	}
}
