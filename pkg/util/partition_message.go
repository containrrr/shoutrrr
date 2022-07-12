package util

import (
	t "github.com/containrrr/shoutrrr/pkg/types"

	"strings"
)

const ellipsis = " [...]"

// PartitionMessage splits a string into chunks that is at most chunkSize runes, it will search the last distance runes
// for a whitespace to make the split appear nicer. It will keep adding chunks until it reaches maxCount chunks, or if
// the total amount of runes in the chunks reach maxTotal.
// The chunks are returned together with the number of omitted runes (that did not fit into the chunks)
func PartitionMessage(input string, limits t.MessageLimit, distance int) (items []t.MessageItem, omitted int) {
	runes := []rune(input)
	chunkOffset := 0
	maxTotal := Min(len(runes), limits.TotalChunkSize)
	maxCount := limits.ChunkCount - 1

	if len(input) == 0 {
		// If the message is empty, return an empty array
		omitted = 0
		return
	}

	for i := 0; i < maxCount; i++ {
		// If no suitable split point is found, use the chunkSize
		chunkEnd := chunkOffset + limits.ChunkSize
		// ... and start next chunk directly after this one
		nextChunkStart := chunkEnd
		if chunkEnd >= maxTotal {
			// The chunk is smaller than the limit, no need to search
			chunkEnd = maxTotal
			nextChunkStart = maxTotal
		} else {
			for r := 0; r < distance; r++ {
				rp := chunkEnd - r
				if runes[rp] == '\n' || runes[rp] == ' ' {
					// Suitable split point found
					chunkEnd = rp
					// Since the split is on a whitespace, skip it in the next chunk
					nextChunkStart = chunkEnd + 1
					break
				}
			}
		}

		items = append(items, t.MessageItem{
			Text: string(runes[chunkOffset:chunkEnd]),
		})

		chunkOffset = nextChunkStart
		if chunkOffset >= maxTotal {
			break
		}
	}

	return items, len(runes) - chunkOffset
}

// Ellipsis returns a string that is at most maxLength characters with a ellipsis appended if the input was longer
func Ellipsis(text string, maxLength int) string {
	if len(text) > maxLength {
		text = text[:maxLength-len(ellipsis)] + ellipsis
	}
	return text
}

// MessageItemsFromLines creates a set of MessageItems that is compatible with the supplied limits
func MessageItemsFromLines(plain string, limits t.MessageLimit) (batches [][]t.MessageItem) {
	maxCount := limits.ChunkCount

	lines := strings.Split(plain, "\n")
	batches = make([][]t.MessageItem, 0)
	items := make([]t.MessageItem, 0, Min(maxCount, len(lines)))

	totalLength := 0
	for _, line := range lines {

		maxLen := limits.ChunkSize

		if len(items) == maxCount || totalLength+maxLen > limits.TotalChunkSize {
			batches = append(batches, items)
			items = items[:0]
		}

		runes := []rune(line)
		if len(runes) > maxLen {
			// Trim and add ellipsis
			runes = runes[:maxLen-len(ellipsis)]
			line = string(runes) + ellipsis
		}

		if len(runes) < 1 {
			continue
		}

		items = append(items, t.MessageItem{
			Text: line,
		})

		totalLength += len(runes)
	}

	if len(items) > 0 {
		batches = append(batches, items)
	}

	return batches
}
