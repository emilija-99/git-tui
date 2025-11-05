package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type fileItem struct{ status, path string }

func (f fileItem) Title() string       { return f.path }
func (f fileItem) Description() string { return f.status }
func (f fileItem) FilterValue() string { return f.path }

type Model struct {
	dir       string
	keys      KeyMap
	list      list.Model
	diff      string
	showDiff  bool
	commitBox textarea.Model
	spin      spinner.Model
	loading   bool
	err       error
}

func New(dir string) Model {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 40, 20)
	l.Title = "Changes"
	ta := textarea.New()
	ta.Placeholder = "Commit message. Enter to commit. Esc to cancel."
	sp := spinner.New()
	return Model{
		dir: dir, keys: DefaultKeyMap(), list: l,
		commitBox: ta, spin: sp,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.loadStatus(), spinner.Tick)
}

func (m Model) loadStatus() tea.Cmd {
	m.loading = true
	return func() tea.Msg {
		ctx, cancel := TimeoutCtx(5 * time.Second)
		defer cancel()
		lines, err := Status(ctx, m.dir)
		fmt.Printf("%s", lines)
		if err != nil {
			return errMsg{err}
		}
		items := make([]list.Item, 0, len(lines))
		for _, ln := range lines {
			fmt.Printf("ln: ", ln)
			if strings.TrimSpace(ln) == "" {
				continue
			}
			// ovaj mu vraca vrednosti, --porcelain ima status
			// ovde hocu da ispisem po grupama, staged.. da napravim detaljam opis i onda na tab diff radi
			st := strings.TrimSpace(ln[:3])
			path := strings.TrimSpace(ln[4:])
			items = append(items, fileItem{status: st, path: path})
		}
		return statusMsg{items}
	}
}

type statusMsg struct{ items []list.Item }
type diffMsg struct{ text string }
type errMsg struct{ e error }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statusMsg:
		m.loading = false
		m.list.SetItems(msg.items)
		return m, nil
	case diffMsg:
		m.loading = false
		m.diff = msg.text
		return m, nil
	case errMsg:
		m.loading = false
		m.err = msg.e
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Refresh):
			return m, m.loadStatus()
		case key.Matches(msg, m.keys.ToggleDiff):
			m.showDiff = !m.showDiff
			return m, m.loadDiff(false)
		case key.Matches(msg, m.keys.Stage):
			return m, m.stageSelected()
		case key.Matches(msg, m.keys.Unstage):
			return m, m.unstageSelected()
		case key.Matches(msg, m.keys.Commit):
			if m.commitBox.Focused() {
				msgTxt := strings.TrimSpace(m.commitBox.Value())
				if msgTxt != "" {
					return m, m.commit(msgTxt)
				}
			} else {
				m.commitBox.Focus()
			}
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) currentPath() (string, bool) {
	it, ok := m.list.SelectedItem().(fileItem)
	if !ok {
		return "", false
	}
	return it.path, true
}

func (m Model) loadDiff(staged bool) tea.Cmd {
	path, ok := m.currentPath()
	if !ok {
		return nil
	}
	m.loading = true
	return func() tea.Msg {
		ctx, cancel := TimeoutCtx(5 * time.Second)
		defer cancel()
		d, err := Diff(ctx, m.dir, path, staged)
		if err != nil {
			return errMsg{err}
		}
		if strings.TrimSpace(d) == "" {
			d = "(no diff)"
		}
		return diffMsg{d}
	}
}

func (m Model) stageSelected() tea.Cmd {
	path, ok := m.currentPath()
	if !ok {
		return nil
	}
	return func() tea.Msg {
		ctx, c := TimeoutCtx(5 * time.Second)
		defer c()
		if err := Add(ctx, m.dir, path); err != nil {
			return errMsg{err}
		}
		return statusMsg{nil} // trigger reload
	}
}

func (m Model) unstageSelected() tea.Cmd {
	path, ok := m.currentPath()
	if !ok {
		return nil
	}
	return func() tea.Msg {
		ctx, c := TimeoutCtx(5 * time.Second)
		defer c()
		if err := Unstage(ctx, m.dir, path); err != nil {
			return errMsg{err}
		}
		return statusMsg{nil}
	}
}

func (m Model) commit(msg string) tea.Cmd {
	return func() tea.Msg {
		ctx, c := TimeoutCtx(10 * time.Second)
		defer c()
		if err := Commit(ctx, m.dir, msg); err != nil {
			return errMsg{err}
		}
		return statusMsg{nil}
	}
}

var (
	title  = lipgloss.NewStyle().Bold(true)
	footer = lipgloss.NewStyle().Faint(true)
)

func (m Model) View() string {
	if m.commitBox.Focused() {
		return title.Render("Commit") + "\n\n" + m.commitBox.View()
	}
	left := m.list.View()
	right := m.diff
	if !m.showDiff {
		right = "tab: toggle diff"
	}
	body := lipgloss.JoinHorizontal(lipgloss.Top, left, "\n", right)
	hint := footer.Render("j/k: move  space: stage  u: unstage  c: commit  P: push  f: pull  r: refresh  q: quit")
	if m.err != nil {
		hint += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("error: "+m.err.Error())
	}
	return body + "\n\n" + hint
}
