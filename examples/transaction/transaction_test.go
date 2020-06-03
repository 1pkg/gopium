package transaction

import "testing"

func BenchmarkGenerate(b *testing.B) {
	// run benchmark for generation
	for n := 0; n < b.N; n++ {
		generate(1000000)
	}
}

func BenchmarkNormalize(b *testing.B) {
	// generate one million of transactions
	// and reset benchmark timer
	transactions := generate(1000000)
	b.ResetTimer()
	// run benchmark for normalization
	for n := 0; n < b.N; n++ {
		normalize(transactions)
	}
}

func BenchmarkCompress(b *testing.B) {
	// generate one million of normalized transactions
	// and reset benchmark timer
	transactions := normalize(generate(1000000))
	b.ResetTimer()
	// run benchmark for compressing
	for n := 0; n < b.N; n++ {
		compress(transactions)
	}
}
