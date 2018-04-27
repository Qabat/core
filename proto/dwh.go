package sonm

import "strings"

const (
	NumNetflags = 3
)

func NewBenchmarks(benchmarks []uint64) (*Benchmarks, error) {
	b := &Benchmarks{
		Values: make([]uint64, len(benchmarks)),
	}
	copy(b.Values, benchmarks)
	if err := b.Validate(); err != nil {
		return nil, err
	}
	return b, nil
}

func (m *Benchmarks) ToArray() []uint64 {
	return m.Values
}

func UintToNetflags(flags uint64) [NumNetflags]bool {
	var fixedNetflags [3]bool
	for idx := 0; idx < NumNetflags; idx++ {
		fixedNetflags[NumNetflags-1-idx] = flags&(1<<uint64(idx)) != 0
	}

	return fixedNetflags
}

func NetflagsToUint(flags [NumNetflags]bool) uint64 {
	var netflags uint64
	for idx, flag := range flags {
		if flag {
			netflags |= 1 << uint64(NumNetflags-1-idx)
		}
	}

	return netflags
}

func (r *DealsRequest) ToLower() {
	r.ConsumerID = strings.ToLower(r.ConsumerID)
	r.SupplierID = strings.ToLower(r.SupplierID)
	r.MasterID = strings.ToLower(r.MasterID)
}

func (r *OrdersRequest) ToLower() {
	r.AuthorID = strings.ToLower(r.AuthorID)
	r.CounterpartyID = strings.ToLower(r.CounterpartyID)
}

func (o *Order) ToLower() {
	o.AuthorID = strings.ToLower(o.AuthorID)
	o.CounterpartyID = strings.ToLower(o.CounterpartyID)
}

func (d *Deal) ToLower() {
	d.ConsumerID = strings.ToLower(d.ConsumerID)
	d.SupplierID = strings.ToLower(d.SupplierID)
	d.MasterID = strings.ToLower(d.MasterID)
}

func (v *Validator) ToLower() {
	v.Id = strings.ToLower(v.Id)
}

func (c *Certificate) ToLower() {
	c.OwnerID = strings.ToLower(c.OwnerID)
	c.ValidatorID = strings.ToLower(c.ValidatorID)
}
