package false_sharing_cpu_l1

import "testing"

func BenchmarkGenerate(b *testing.B) {
	// run benchmark for generation
	for n := 0; n < b.N; n++ {
		generate(1000000)
	}
}

func BenchmarkNormalize(b *testing.B) {
	// generate one millon of transactions
	// and reset benchmark timer
	transactions := generate(1000000)
	b.ResetTimer()
	// run benchmark for normalization
	for n := 0; n < b.N; n++ {
		normalize(transactions)
	}
}

func BenchmarkCompress(b *testing.B) {
	// generate one millon of normalized transactions
	// and reset benchmark timer
	transactions := normalize(generate(1000000))
	b.ResetTimer()
	// run benchmark for compressing
	for n := 0; n < b.N; n++ {
		compress(transactions)
	}
}
