package tui

import (
	"bytes"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/sLg00/nba-now-tui/cmd/helpers"
	"github.com/sLg00/nba-now-tui/cmd/nba/nbaAPI"
)

// mockNewsClient implements a simple mock for the NewsClient
type mockNewsClient struct{}

func (m *mockNewsClient) FetchLatestNews() ([]nbaAPI.NewsArticle, error) {
	return []nbaAPI.NewsArticle{
		{
			Title:     "Test Headline 1",
			URL:       "https://www.nba.com/news/article1",
			Timestamp: "2025-01-01",
		},
		{
			Title:     "Test Headline 2",
			URL:       "https://www.nba.com/news/article2",
			Timestamp: "2025-01-02",
		},
	}, nil
}

// TestNewsViewInitialization tests that the news view initializes correctly
func TestNewsViewInitialization(t *testing.T) {
	ts := helpers.SetupTest()
	defer ts.CleanUpTest()

	// Initialize the news view
	model, cmd, err := NewNewsView(tea.WindowSizeMsg{Width: 80, Height: 30})
	if err != nil {
		t.Fatalf("Failed to initialize news view: %v", err)
	}

	if model == nil {
		t.Fatal("Expected non-nil model")
	}

	if cmd == nil {
		t.Fatal("Expected non-nil cmd")
	}
}

// TestNewsViewRender tests that the news view renders
func TestNewsViewRender(t *testing.T) {
	ts := helpers.SetupTest()
	defer ts.CleanUpTest()

	// Initialize the news view
	model, _, err := NewNewsView(tea.WindowSizeMsg{Width: 80, Height: 30})
	if err != nil {
		t.Fatalf("Failed to initialize news view: %v", err)
	}

	// Create a test model
	testModel := teatest.NewTestModel(t, model)

	// Wait for "Loading" to appear
	teatest.WaitFor(t, testModel.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Recent News"))
	})
}

// TestNewsView_FetchNewsMsg tests handling of newsFetchedMsg
func TestNewsView_FetchNewsMsg(t *testing.T) {
	ts := helpers.SetupTest()
	defer ts.CleanUpTest()

	// Create a test NewsView
	model := &NewsModel{
		list:     list.New([]list.Item{}, list.NewDefaultDelegate(), 80, 24),
		quitting: false,
		width:    80,
		height:   30,
	}

	// Create a test newsFetchedMsg
	msg := newsFetchedMsg{
		articles: []nbaAPI.NewsArticle{
			{
				Title:     "Test Article",
				URL:       "https://example.com",
				Timestamp: "2025-01-01",
			},
		},
		err: nil,
	}

	// Update the model with the message
	updatedModel, _ := model.Update(msg)
	if updatedModel == nil {
		t.Fatal("Expected non-nil model")
	}

	newsModel, ok := updatedModel.(*NewsModel)
	if !ok {
		t.Fatal("Expected NewsModel type")
	}
	// Check that the list now has items
	if len(newsModel.list.Items()) != 1 {
		t.Errorf("Expected list to have 1 item, got: %d", len(newsModel.list.Items()))
	}

	// Check the first item's title
	item, ok := newsModel.list.Items()[0].(NewsItem)
	if !ok {
		t.Fatal("Expected NewsItem type for list item")
	}

	title := item.Title()
	if title != "Test Article" {
		t.Errorf("Expected item title 'Test Article', got: %s", title)
	}
}
