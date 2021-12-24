package transaction

import (
	"math"
	"math/rand"
	"sort"
	"sync"
)

// transaction defines business transaction
type transaction struct {
	void     bool
	amount   float64
	serial   uint64
	skip     bool
	discount float64
} // struct size: 26 bytes; struct align: 8 bytes; struct aligned size: 40 bytes; struct ptr scan size: 0 bytes; - 🌺 gopium @1pkg

// aggregate defines compressed set of transactions
type aggregate struct {
	total float64  `gopium:"filter_pads,false_sharing_cpu_l1,struct_annotate_comment,add_tag_group_force"`
	_     [56]byte `gopium:"filter_pads,false_sharing_cpu_l1,struct_annotate_comment,add_tag_group_force"`
} // struct size: 64 bytes; struct align: 8 bytes; struct aligned size: 64 bytes; struct ptr scan size: 0 bytes; - 🌺 gopium @1pkg

// generate creates n pseudo random transactions
func generate(number uint) []transaction {
	// generate n pseudo random transactions
	transactions := make([]transaction, 0, number)
	for i := 0; i < int(number); i++ {
		transactions = append(transactions, transaction{
			void:     i%10 == 0,
			amount:   math.Abs(rand.Float64()),
			serial:   uint64(i) + 1,
			skip:     i%25 == 0,
			discount: rand.Float64(),
		})
	}
	// and shuffle them
	for i := range transactions {
		j := rand.Intn(i + 1)
		transactions[i], transactions[j] = transactions[j], transactions[i]
	}
	return transactions
}

// normalize preprocess list of transactions before compressing
func normalize(transactions []transaction) []transaction {
	// filter and normalize transactions
	normalized := make([]transaction, 0, len(transactions))
	for _, trx := range transactions {
		if trx.skip || trx.serial == 0 {
			continue
		}
		if trx.void {
			trx.amount = -trx.amount
		}
		trx.discount = math.Abs(trx.discount)
		normalized = append(normalized, trx)
	}
	// sort transactions by serial
	sort.Slice(normalized, func(i int, j int) bool {
		return normalized[i].serial < normalized[j].serial
	})
	return normalized
}

// compress builds single aggregate from provided normalized transactions list
func compress(transactions []transaction) aggregate {
	var amount, discont aggregate
	var wg sync.WaitGroup
	wg.Add(2)
	// run amount calculation in separate goroutine
	go func() {
		for _, tr := range transactions {
			amount.total += tr.amount
		}
		wg.Done()
	}()
	// run discounts calculation in separate goroutine
	go func() {
		for _, tr := range transactions {
			discont.total += tr.discount
		}
		wg.Done()
	}()
	wg.Wait()
	// apply discount logic to final aggregate
	if discont.total > amount.total/2 {
		discont.total = amount.total / 2
	}
	result := amount
	result.total -= discont.total
	return result
}
