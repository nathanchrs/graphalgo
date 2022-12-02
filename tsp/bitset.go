package solver

type Bitset uint32

func (b Bitset) Set(idx int) Bitset {
	if idx < 0 || idx >= 32 {
		panic("Bitset: idx out of range")
	}
	return Bitset(uint32(b) | (uint32(1) << uint32(idx)))
}

func (b Bitset) Unset(idx int) Bitset {
	if idx < 0 || idx >= 32 {
		panic("Bitset: idx out of range")
	}
	return Bitset(uint32(b) & ^(uint32(1) << uint32(idx)))
}

func (b Bitset) IsSet(idx int) bool {
	if idx < 0 || idx >= 32 {
		panic("Bitset: idx out of range")
	}
	return (uint32(b) & (uint32(1) << uint32(idx))) > 0
}

func (b Bitset) IsAnySet() bool {
	return b != 0
}

func NewBitsetAllSet(n int) Bitset {
	if n < 0 || n >= 32 {
		panic("Bitset: n out of range")
	}
	return Bitset((uint32(1) << uint32(n)) - 1)
}

func CountBitArrangements(n int) int {
	if n < 0 || n >= 32 {
		panic("Bitset: n out of range")
	}
	return int(uint32(1) << uint32(n))
}
