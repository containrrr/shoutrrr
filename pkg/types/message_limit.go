package types

type MessageLimit struct {
	ChunkSize 	   int
	TotalChunkSize int

	// Maximum number of chunks (including the last chunk for meta data)
	ChunkCount int
}