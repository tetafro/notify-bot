package feed

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFeed(t *testing.T) {
	t.Run("fetch new item", func(t *testing.T) {
		server := httptest.NewServer(&testRSSServer{data: true})
		defer server.Close()

		st := &testStorage{
			time: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
		}

		f := NewRSSFeed(st, server.URL)
		items, err := f.Fetch()
		assert.NoError(t, err)
		assert.Len(t, items, 1)

		expected := Item{
			Title:     "Item Title",
			Image:     "https://example.com/image.png",
			Published: time.Date(2020, 1, 1, 15, 0, 0, 0, time.UTC),
		}
		assert.Equal(t, expected, items[0])
	})

	t.Run("no new items", func(t *testing.T) {
		server := httptest.NewServer(&testRSSServer{data: true})
		defer server.Close()

		st := &testStorage{
			time: time.Date(2021, 1, 1, 10, 0, 0, 0, time.UTC),
		}

		f := NewRSSFeed(st, server.URL)
		items, err := f.Fetch()
		assert.NoError(t, err)
		assert.Len(t, items, 0)
	})

	t.Run("no items", func(t *testing.T) {
		server := httptest.NewServer(&testRSSServer{data: false})
		defer server.Close()

		st := &testStorage{
			time: time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC),
		}

		f := NewRSSFeed(st, server.URL)
		items, err := f.Fetch()
		assert.NoError(t, err)
		assert.Len(t, items, 0)
	})

	t.Run("500 error url", func(t *testing.T) {
		server := httptest.NewServer(&testRSSServer{err: true})
		defer server.Close()

		f := NewRSSFeed(nil, server.URL)
		_, err := f.Fetch()
		assert.EqualError(t, err,
			"parse url: http error: 500 Internal Server Error")
	})

	t.Run("invalid url", func(t *testing.T) {
		f := NewRSSFeed(nil, "xxx://example.com")
		_, err := f.Fetch()
		assert.EqualError(t, err,
			`parse url: Get "xxx://example.com": unsupported protocol scheme "xxx"`)
	})
}

type testRSSServer struct {
	data bool
	err  bool
}

func (s *testRSSServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if s.err {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !s.data {
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="utf-8"?>` + "\n" +
			`<feed xmlns="http://www.w3.org/2005/Atom" xml:lang="en">` +
			`<title>Feed Title</title>` +
			`<id>feed_id</id>` +
			`<updated>2020-01-01T15:00:00Z</updated>` +
			`</feed>`))
		return
	}
	_, _ = w.Write([]byte(`<?xml version="1.0" encoding="utf-8"?>` + "\n" +
		`<feed xmlns="http://www.w3.org/2005/Atom" xml:lang="en">` +
		`<title>Feed Title</title>` +
		`<id>feed_id</id>` +
		`<updated>2020-01-01T15:00:00Z</updated>` +
		`<entry>` +
		`<title>Item Title</title>` +
		`<id>item_id</id>` +
		`<updated>2020-01-01T15:00:00Z</updated>` +
		`<summary type="html">&lt;img src="https://example.com/image.png"/&gt;</summary>` +
		`</entry>` +
		`</feed>`))
}

type testStorage struct {
	time time.Time
}

func (s *testStorage) GetLastUpdate(feed string) time.Time {
	return s.time
}

func (s *testStorage) SaveLastUpdate(feed string, t time.Time) error {
	return nil
}
