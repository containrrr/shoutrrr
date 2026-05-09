package types

// MessageLimit is used for declaring the payload limits for services upstream APIs
type MessageLimit struct {
	ChunkSize      uint
	TotalChunkSize uint

	// Maximum number of chunks (including the last chunk for meta data)
	ChunkCount uint
}
