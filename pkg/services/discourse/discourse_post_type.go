package discourse

import (
	"github.com/containrrr/shoutrrr/pkg/format"
	"github.com/containrrr/shoutrrr/pkg/types"
)

type postType int

type postTypeVals struct {
	Post           postType
	Topic          postType
	PrivateMessage postType

	// Enum is the EnumFormatter instance for PostTypes
	Enum types.EnumFormatter
}

// PostTypes is the enum helper for identifying the discourse post type / archetype
var PostTypes = &postTypeVals{
	Post:           0,
	Topic:          1,
	PrivateMessage: 2,

	Enum: format.CreateEnumFormatter(
		[]string{
			"regular",
			"banner",
			"private",
		}),
}

// String returns the user-facing string representation of the post type
func (pt postType) String() string {
	return PostTypes.Enum.Print(int(pt))
}

// Archetype returns the API-facing string representation of the post type
func (pt postType) Archetype() string {
	if pt == PostTypes.PrivateMessage {
		return "private_message"
	}
	return PostTypes.Enum.Print(int(pt))
}
