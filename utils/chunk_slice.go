package utils

func ChunkSlice[T any](slice []T, chunkSize int) [][]T {
	var chunks [][]T
	if len(slice) == 0 {
		return make([][]T, 0)
	}
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		// necessary check to avoid slicing beyond
		// slice capacity
		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}
