package hash

func BKDRHash(s string) uint64 {
	var seed uint64 = 13131
	var hash uint64 = 0
	len := len(s)

	for i := 0; i < len; i++ {
		hash = hash*seed + uint64(s[i])
	}

	return hash & 0x7FFFFFFFFFFFFFFF
}
