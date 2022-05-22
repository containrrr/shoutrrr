package util

import (
	"strings"

	"github.com/containrrr/shoutrrr/pkg/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Partition Message", func() {
	limits := types.MessageLimit{
		ChunkSize:      2000,
		TotalChunkSize: 6000,
		ChunkCount:     10,
	}
	When("given a message that exceeds the max length", func() {
		When("not splitting by lines", func() {
			It("should return a payload with chunked messages", func() {

				items, _ := testPartitionMessage(42, limits, 100)

				Expect(len(items[0].Text)).To(Equal(1994))
				Expect(len(items[1].Text)).To(Equal(1999))
				Expect(len(items[2].Text)).To(Equal(205))
			})
			It("omit characters above total max", func() {
				items, _ := testPartitionMessage(62, limits, 100)

				Expect(len(items[0].Text)).To(Equal(1994))
				Expect(len(items[1].Text)).To(Equal(1999))
				Expect(len(items[2].Text)).To(Equal(1999))
				Expect(len(items[3].Text)).To(Equal(5))
			})
			It("should handle messages with a size modulus of chunksize", func() {
				items, _ := testPartitionMessage(20, limits, 100)
				// Last word fits in the chunk size
				Expect(len(items[0].Text)).To(Equal(2000))

				items, _ = testPartitionMessage(40, limits, 100)
				// Now the last word of the first chunk will be concatenated with
				// the first word of the second chunk, and so it does not fit in the chunk anymore
				Expect(len(items[0].Text)).To(Equal(1994))
				Expect(len(items[1].Text)).To(Equal(1999))
				Expect(len(items[2].Text)).To(Equal(5))
			})
			When("the message is empty", func() {
				It("should return no items", func() {
					items, _ := testPartitionMessage(0, limits, 100)
					Expect(items).To(BeEmpty())
				})
			})
			When("given an input without whitespace", func() {
				It("should not crash, regardless of length", func() {
					testString := ""
					for inputLen := 1; inputLen < 8000; inputLen++ {
						testString += "z"
						items, omitted := PartitionMessage(testString, limits, 7)
						included := 0
						for ii, item := range items {
							expectedSize := limits.ChunkSize
							// If this is the last chunk, and we haven't reached the total maximum
							if inputLen < limits.TotalChunkSize && ii == len(items)-1 {
								// ...and the chunk wouldn't be empty
								chunkOffset := inputLen % limits.ChunkSize
								if chunkOffset > 0 {
									// expect the chunk to contain the "rest" of the runes
									expectedSize = chunkOffset
								}
								// the last chunk should never be empty, so treat it as one of the full ones
							}
							included += len(item.Text)
							Expect(item.Text).To(HaveLen(expectedSize))
						}
						Expect(omitted + included).To(Equal(inputLen))
					}
				})
			})
		})
		When("splitting by lines", func() {
			It("should return a payload with chunked messages", func() {
				items, omitted := testMessageItemsFromLines(18, limits, 2)

				Expect(len(items[0].Text)).To(Equal(200))
				Expect(len(items[8].Text)).To(Equal(200))

				Expect(omitted).To(Equal(0))
			})
			It("omit characters above total max", func() {
				items, omitted := testMessageItemsFromLines(19, limits, 2)

				Expect(len(items[0].Text)).To(Equal(200))
				Expect(len(items[8].Text)).To(Equal(200))

				Expect(omitted).To(Equal(100))
			})
			It("should trim characters above chunk size", func() {
				hundreds := 42
				repeat := 21
				items, omitted := testMessageItemsFromLines(hundreds, limits, repeat)

				Expect(len(items[0].Text)).To(Equal(limits.ChunkSize))
				Expect(len(items[1].Text)).To(Equal(limits.ChunkSize))

				// Trimmed characters do not count towards the total omitted count
				Expect(omitted).To(Equal(0))
			})

			It("omit characters above total chunk size", func() {
				hundreds := 100
				repeat := 20
				items, omitted := testMessageItemsFromLines(hundreds, limits, repeat)

				Expect(len(items[0].Text)).To(Equal(limits.ChunkSize))
				Expect(len(items[1].Text)).To(Equal(limits.ChunkSize))
				Expect(len(items[2].Text)).To(Equal(limits.ChunkSize))

				maxRunes := hundreds * 100
				expectedOmitted := maxRunes - limits.TotalChunkSize

				Expect(omitted).To(Equal(expectedOmitted))
			})

		})

	})
})

const hundredChars = "this string is exactly (to the letter) a hundred characters long which will make the send func error"

func testMessageItemsFromLines(hundreds int, limits types.MessageLimit, repeat int) (items []types.MessageItem, omitted int) {

	builder := strings.Builder{}

	ri := 0
	for i := 0; i < hundreds; i++ {
		builder.WriteString(hundredChars)
		ri++
		if ri == repeat {
			builder.WriteRune('\n')
			ri = 0
		}
	}

	items, omitted = MessageItemsFromLines(builder.String(), limits)

	maxChunkSize := Min(limits.ChunkSize, repeat*100)

	expectedChunkCount := Min(limits.TotalChunkSize/maxChunkSize, Min(hundreds/repeat, limits.ChunkCount-1))
	Expect(len(items)).To(Equal(expectedChunkCount), "Chunk count")

	return
}

func testPartitionMessage(hundreds int, limits types.MessageLimit, distance int) (items []types.MessageItem, omitted int) {
	builder := strings.Builder{}

	for i := 0; i < hundreds; i++ {
		builder.WriteString(hundredChars)
	}

	items, omitted = PartitionMessage(builder.String(), limits, distance)

	contentSize := Min(hundreds*100, limits.TotalChunkSize)
	// expectedChunkCount := CeilDiv(contentSize, limits.ChunkSize-1)
	expectedOmitted := Max(0, (hundreds*100)-contentSize)

	ExpectWithOffset(1, omitted).To(Equal(expectedOmitted))
	// ExpectWithOffset(1, len(items)).To(Equal(expectedChunkCount))

	return
}
