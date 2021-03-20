package bot

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBot(t *testing.T) {
	api := &mockAPI{}
	st := &mockStorage{}
	bot, err := NewBot(api, st, nil)
	assert.NoError(t, err)
	assert.NotNil(t, bot)
}

func TestMapToSlice(t *testing.T) {
	testCases := []struct {
		name string
		m    map[int64]struct{}
		s    []int64
	}{
		{
			name: "test-1",
			m:    map[int64]struct{}{1: {}, 2: {}, 3: {}},
			s:    []int64{1, 2, 3},
		},
		{
			name: "test-2",
			m:    map[int64]struct{}{1: {}},
			s:    []int64{1},
		},
		{
			name: "test-3",
			m:    map[int64]struct{}{},
			s:    []int64{},
		},
		{
			name: "test-4",
			m:    nil,
			s:    nil,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := mapToSlice(tt.m)
			sort.Slice(s, func(i, j int) bool {
				return s[i] < s[j]
			})
			assert.Equal(t, tt.s, s)
		})
	}
}