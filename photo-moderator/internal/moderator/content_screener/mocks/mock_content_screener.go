package content_screener_mocks

import "github.com/rhargreaves/dog-walking/photo-moderator/internal/moderator/content_screener"

type MockContentScreener struct {
	ScreenImageFunc func(id string) (*content_screener.ContentScreenerResult, error)
}

func (m *MockContentScreener) ScreenImage(id string) (*content_screener.ContentScreenerResult, error) {
	return m.ScreenImageFunc(id)
}
