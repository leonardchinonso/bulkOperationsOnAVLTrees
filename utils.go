package main

func SplitByteArray(b *[]byte) (*[]byte, *[]byte) {
	n := len(*b)

	if n < 4 {
		return nil, nil
	}

	mark := n / 4
	if n%4 != 0 {
		mark += 1
	}

	first := (*b)[:mark]
	second := (*b)[mark+1:]

	return &first, &second
}
