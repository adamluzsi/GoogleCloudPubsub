package publisher

import "cloud.google.com/go/pubsub"

func chunkSlice(slices [][]byte) [][][]byte {
	chunks := make([][][]byte, 0)
	chunk := make([][]byte, 0, pubsub.MaxPublishBatchSize)

	for _, slice := range slices {
		if len(chunk) != pubsub.MaxPublishBatchSize {
			chunk = append(chunk, slice)
		} else {
			chunks = append(chunks, chunk)
			chunk = make([][]byte, 0, pubsub.MaxPublishBatchSize)
		}
	}

	if len(chunk) != 0 {
		chunks = append(chunks, chunk)
	}

	return chunks
}
