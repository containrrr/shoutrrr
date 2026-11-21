package util

import (
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/containrrr/shoutrrr/pkg/types"
)

func TestPartitionMessageBounds(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	minInt := -maxInt - 1
	cases := []struct {
		name     string
		input    string
		limits   types.MessageLimit
		distance int
		want     []string
		omitted  int
	}{
		{"search before input", "abcdef", types.MessageLimit{ChunkSize: 2, TotalChunkSize: 6, ChunkCount: 4}, 4, []string{"ab", "cd", "ef"}, 0},
		{"search before current chunk", "a bcdefg", types.MessageLimit{ChunkSize: 3, TotalChunkSize: 8, ChunkCount: 4}, 100, []string{"a", "bcd", "efg"}, 0},
		{"largest search distance", "abcdef", types.MessageLimit{ChunkSize: 2, TotalChunkSize: 6, ChunkCount: 4}, maxInt, []string{"ab", "cd", "ef"}, 0},
		{"negative search distance", "ab cde", types.MessageLimit{ChunkSize: 2, TotalChunkSize: 6, ChunkCount: 4}, minInt, []string{"ab", " c", "de"}, 0},
		{"boundary whitespace", "ab cde", types.MessageLimit{ChunkSize: 2, TotalChunkSize: 6, ChunkCount: 4}, 1, []string{"ab", "cd", "e"}, 0},
		{"unicode and total limit", "界🙂éßxy", types.MessageLimit{ChunkSize: 2, TotalChunkSize: 5, ChunkCount: 4}, maxInt, []string{"界🙂", "éß", "x"}, 1},
		{"metadata slot", "abcdef", types.MessageLimit{ChunkSize: 2, TotalChunkSize: 6, ChunkCount: 2}, 1, []string{"ab"}, 4},
		{"unbounded positive limits", "abcdef", types.MessageLimit{ChunkSize: maxInt, TotalChunkSize: maxInt, ChunkCount: maxInt}, maxInt, []string{"abcdef"}, 0},
		{"zero chunk size", "abc", types.MessageLimit{ChunkSize: 0, TotalChunkSize: 3, ChunkCount: 4}, 0, nil, 3},
		{"negative chunk size", "abc", types.MessageLimit{ChunkSize: -1, TotalChunkSize: 3, ChunkCount: 4}, 0, nil, 3},
		{"smallest chunk size", "abc", types.MessageLimit{ChunkSize: minInt, TotalChunkSize: 3, ChunkCount: 4}, 0, nil, 3},
		{"zero total size", "abc", types.MessageLimit{ChunkSize: 2, TotalChunkSize: 0, ChunkCount: 4}, 0, nil, 3},
		{"negative total size", "abc", types.MessageLimit{ChunkSize: 2, TotalChunkSize: -1, ChunkCount: 4}, 0, nil, 3},
		{"smallest total size", "abc", types.MessageLimit{ChunkSize: 2, TotalChunkSize: minInt, ChunkCount: 4}, 0, nil, 3},
		{"metadata only", "abc", types.MessageLimit{ChunkSize: 2, TotalChunkSize: 3, ChunkCount: 1}, 0, nil, 3},
		{"zero chunk count", "abc", types.MessageLimit{ChunkSize: 2, TotalChunkSize: 3, ChunkCount: 0}, 0, nil, 3},
		{"negative chunk count", "abc", types.MessageLimit{ChunkSize: 2, TotalChunkSize: 3, ChunkCount: -1}, 0, nil, 3},
		{"smallest chunk count", "abc", types.MessageLimit{ChunkSize: 2, TotalChunkSize: 3, ChunkCount: minInt}, 0, nil, 3},
		{"empty input and invalid limits", "", types.MessageLimit{ChunkSize: minInt, TotalChunkSize: minInt, ChunkCount: minInt}, maxInt, nil, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			items, omitted := PartitionMessage(tc.input, tc.limits, tc.distance)
			var texts []string
			for _, item := range items {
				texts = append(texts, item.Text)
			}
			if !reflect.DeepEqual(texts, tc.want) || omitted != tc.omitted {
				t.Fatalf("got (%q, %d omitted), want (%q, %d omitted)", texts, omitted, tc.want, tc.omitted)
			}
		})
	}
}

func FuzzPartitionMessageBounds(f *testing.F) {
	maxInt := int(^uint(0) >> 1)
	minInt := -maxInt - 1
	f.Add("abcdef", 2, 6, 4, 4)
	f.Add("a bcdefg", 3, 8, 4, 100)
	f.Add("界🙂éßxy", 2, 5, 4, maxInt)
	f.Add("abc", minInt, minInt, minInt, minInt)
	f.Add("abc", 0, 3, maxInt, maxInt)
	f.Add("abc", maxInt, maxInt, maxInt, maxInt)
	f.Fuzz(func(t *testing.T, input string, chunkSize, totalSize, count, distance int) {
		limits := types.MessageLimit{ChunkSize: chunkSize, TotalChunkSize: totalSize, ChunkCount: count}
		items, omitted := PartitionMessage(input, limits, distance)
		runeCount := utf8.RuneCountInString(input)
		if omitted < 0 || omitted > runeCount {
			t.Fatalf("invalid omitted count: %d of %d", omitted, runeCount)
		}
		if chunkSize <= 0 || totalSize <= 0 || count <= 1 {
			if len(items) != 0 || omitted != runeCount {
				t.Fatalf("unusable limits must omit all input: %d items, %d omitted", len(items), omitted)
			}
			return
		}
		if len(items) > count-1 {
			t.Fatalf("metadata chunk not reserved: %d items for limit %d", len(items), count)
		}
		var texts []string
		returned := 0
		for _, item := range items {
			size := utf8.RuneCountInString(item.Text)
			if size > chunkSize {
				t.Fatalf("chunk has %d runes, limit %d", size, chunkSize)
			}
			returned += size
			texts = append(texts, item.Text)
		}
		consumed := runeCount - omitted
		if consumed > totalSize || returned > consumed || consumed-returned > len(items) {
			t.Fatalf("invalid rune accounting: returned=%d consumed=%d total limit=%d", returned, consumed, totalSize)
		}
		if !strings.ContainsAny(input, " \n") && (returned != consumed || strings.Join(texts, "") != string([]rune(input)[:consumed])) {
			t.Fatal("partitioning without whitespace lost or duplicated input")
		}
	})
}
