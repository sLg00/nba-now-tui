package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sLg00/nba-now-tui/cmd/nba/nbaAPI"
	"log"
	"net/url"
	"os/exec"
	"runtime"
)

type NewsItem struct {
	title string
	url   string
	date  string
}

type NewsModel struct {
	list       list.Model
	quitting   bool
	width      int
	height     int
	newsClient *nbaAPI.NewsClient
}

type newsFetchedMsg struct {
	articles []nbaAPI.NewsArticle
	err      error
}

func (i NewsItem) FilterValue() string {
	return i.title
}

func (i NewsItem) Title() string {
	return i.title
}

func (i NewsItem) Description() string {
	return i.date
}

func fetchNewsCmd(nc *nbaAPI.NewsClient) tea.Cmd {
	return func() tea.Msg {
		articles, err := nc.FetchNews()
		return newsFetchedMsg{articles: articles, err: err}
	}
}

func NewNewsView(size tea.WindowSizeMsg) (*NewsModel, tea.Cmd, error) {
	client := nbaAPI.NewClient()
	nc := nbaAPI.NewNewsClient(client.FileSystem, client.Paths)

	delegate := list.NewDefaultDelegate()
	newsModel := list.New([]list.Item{}, delegate, 0, 0)
	newsModel.Title = "Recent News"
	newsModel.SetShowHelp(false)
	newsModel.SetShowStatusBar(false)
	newsModel.SetFilteringEnabled(false)
	newsModel.SetShowPagination(true)

	m := &NewsModel{
		newsClient: nc,
		list:       newsModel,
		width:      size.Width,
		height:     size.Height,
	}

	top, right, bottom, left := DocStyle.GetMargin()
	m.list.SetSize(size.Width-left-right, size.Height-top-bottom-1)

	return m, fetchNewsCmd(nc), nil
}

func (m *NewsModel) Init() tea.Cmd {
	return fetchNewsCmd(m.newsClient)
}

func (m *NewsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		top, right, bottom, left := DocStyle.GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom-1)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keymap.Back):
			return InitMenu()
		case key.Matches(msg, Keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, Keymap.Enter):
			if item, ok := m.list.SelectedItem().(NewsItem); ok {
				log.Printf("Opening URL: %s", item.url)
				go openURL(item.url)
			}
		}
	case newsFetchedMsg:
		if msg.err != nil {
			log.Printf("Error while fetching news: %v", msg.err)
			items := []list.Item{
				NewsItem{
					title: "Error fetching news - press backspace to return",
					url:   "https://www.nba.com/news",
					date:  "Error",
				},
			}
			m.list.SetItems(items)
			return m, nil
		}

		var items []list.Item
		for _, article := range msg.articles {
			items = append(items, NewsItem{
				title: article.Title,
				url:   article.URL,
				date:  article.Timestamp,
			})
		}

		if len(items) == 0 {
			items = append(items, NewsItem{
				title: "No news found - visit NBA.com for latest news",
				url:   "https://www.nba.com/news",
				date:  "Current",
			})
		}
		// Set the items on the list
		m.list.SetItems(items)
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *NewsModel) helpView() string {
	return HelpStyle("\n" + HelpFooter() + "\n")
}

func (m *NewsModel) View() string {
	if m.quitting {
		return ""
	}

	//if m.list.Items() == nil || len(m.list.Items()) == 0 {
	//	return "Loading..."
	//}

	renderedNews := m.list.View() + "\n"
	comboView := DocStyle.Render(renderedNews + m.helpView())
	return comboView
}

// openURL opens the provided URL in the default web browser based on the operating system.
func openURL(rawURL string) {
	u, err := url.Parse(rawURL)
	if err != nil {
		log.Printf("Error parsing URL: %v", err)
		return
	}

	if u.Scheme == "" {
		u.Scheme = "http"
	}

	log.Printf("Opening URL: %s", u.String())

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", u.String())
	case "darwin":
		cmd = exec.Command("open", u.String())
	default:
		cmd = exec.Command("xdg-open", u.String())
	}

	err = cmd.Start()
	if err != nil {
		log.Printf("Error opening URL: %v", err)
	}
}
