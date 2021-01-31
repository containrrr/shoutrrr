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

	for i := 0; i < maxCount; i++ {
		// If no suitable split point is found, use the chunkSize
		chunkEnd := chunkOffset + limits.ChunkSize
		// ... and start next chunk directly after this one
		nextChunkStart := chunkEnd
		if chunkEnd > maxTotal {
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

// MessageItemsFromLines creates a set of MessageItems that is compatible with the supplied limits
func MessageItemsFromLines(plain string, limits t.MessageLimit) (items []t.MessageItem, omitted int) {
	omitted = 0
	maxCount := limits.ChunkCount - 1

	lines := strings.Split(plain, "\n")
	items = make([]t.MessageItem, 0, Min(maxCount, len(lines)))

	totalLength := 0
	for l, line := range lines {
		if l < maxCount && totalLength < limits.TotalChunkSize {
			runes := []rune(line)
			maxLen := limits.ChunkSize
			if totalLength+maxLen > limits.TotalChunkSize {
				maxLen = limits.TotalChunkSize - totalLength
			}
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

		} else {
			omitted += len(line)
		}
	}

	return items, omitted
}
