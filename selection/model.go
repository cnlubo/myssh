package selection

import (
	"bytes"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cnlubo/myssh/common"
	"github.com/mattn/go-runewidth"
	"strings"
	"text/template"
)

// Model implements the bubbletea.Model for a selection prompt.
// Model is a data container used to store TUI status information
type Model struct {
	*Selection
	// Err holds errors that may occur during the execution of
	// the selection prompt.
	Err error

	// MaxWidth limits the width of the view using the Selection's WrapMode.
	MaxWidth int

	filterInput textinput.Model
	// currently displayed choices, after filtering and pagination
	currentChoices []*Choice
	// number of available choices after filtering
	availableChoices int
	// index of current selection in currentChoices slice
	currentIdx        int
	scrollOffset      int
	width             int
	height            int
	tmpl              *template.Template
	resultTmpl        *template.Template
	requestedPageSize int
	quitting          bool

	// pageData data set rendered in real time on the current page
	pageData []*Choice
	// index global real time index
	index int
	// maxIndex global max index
	maxIndex int
	// pageIndex real time index of current page
	pageIndex int
	// pageMaxIndex current page max index
	pageMaxIndex int
}

// ensure that the Model interface is implemented.
var _ tea.Model = &Model{}

// NewModel returns a new selection prompt model for the
// provided choices.
func NewModel(selection *Selection) *Model {
	return &Model{Selection: selection}
}

// Init initializes the selection prompt model.
func (m *Model) Init() tea.Cmd {
	m.reindexChoices()

	if len(m.Choices) == 0 {
		m.Err = fmt.Errorf("no choices provided")

		return tea.Quit
	}

	if m.Template == "" {
		m.Err = fmt.Errorf("empty template")

		return tea.Quit
	}

	m.tmpl, m.Err = m.initTemplate()
	if m.Err != nil {
		return tea.Quit
	}

	m.resultTmpl, m.Err = m.initResultTemplate()
	if m.Err != nil {
		return tea.Quit
	}

	m.filterInput = m.initFilterInput()

	m.currentChoices, m.availableChoices = m.filteredAndPagedChoices()

	m.requestedPageSize = m.PageSize

	m.pageIndex = 0
	m.pageMaxIndex = m.PageSize - 1
	m.index = 0
	m.maxIndex = len(m.Choices) - 1
	if m.PageSize > len(m.Choices) || m.PageSize < 1 {
		m.PageSize = len(m.Choices)
		m.pageData = m.Choices
	} else {
		m.pageData = m.Choices[:m.PageSize]
	}
	return textinput.Blink
}

func (m *Model) initTemplate() (*template.Template, error) {
	tmpl := template.New("view")
	tmpl.Funcs(m.ExtendedTemplateFuncs)
	tmpl.Funcs(common.UtilFuncMap())
	tmpl.Funcs(template.FuncMap{
		//"IsScrollDownHintPosition": func(idx int) bool {
		//	return m.canScrollDown() && (idx == len(m.currentChoices)-1)
		//},
		//"IsScrollUpHintPosition": func(idx int) bool {
		//	return m.canScrollUp() && idx == 0 && m.scrollOffset > 0
		//},
		"Selected": func(c *Choice) string {
			if m.SelectedChoiceStyle == nil {
				return c.String
			}
			return m.SelectedChoiceStyle(m, c, 1)
		},
		"Unselected": func(c *Choice) string {
			if m.UnselectedChoiceStyle == nil {
				return c.String
			}
			return m.UnselectedChoiceStyle(m, c, 1)
		},
		"Cursor": func() string {
			if m.Cursor == "" {
				return DefaultCursor
			}
			return m.Cursor
		},
	})
	tmpl.Funcs(common.ColorFuncMap)

	return tmpl.Parse(m.Template)
}

func (m *Model) initResultTemplate() (*template.Template, error) {
	if m.ResultTemplate == "" {
		return nil, nil // nolint:nilnil
	}

	tmpl := template.New("result")
	tmpl.Funcs(m.ExtendedTemplateFuncs)
	tmpl.Funcs(common.UtilFuncMap())
	tmpl.Funcs(template.FuncMap{
		"Final": func(c *Choice) string {
			if m.FinalChoiceStyle == nil {
				return c.String
			}

			return m.FinalChoiceStyle(c)
		},
	})

	return tmpl.Parse(m.ResultTemplate)
}

func (m *Model) initFilterInput() textinput.Model {
	filterInput := textinput.NewModel()
	filterInput.Prompt = ""
	filterInput.TextStyle = m.FilterInputTextStyle
	filterInput.BackgroundStyle = m.FilterInputBackgroundStyle
	filterInput.PlaceholderStyle = m.FilterInputPlaceholderStyle
	filterInput.CursorStyle = m.FilterInputCursorStyle
	filterInput.Placeholder = m.FilterPlaceholder
	filterInput.Width = 80
	filterInput.Focus()

	return filterInput
}

// Value returns the choice that is currently selected or the final
// choice after the prompt has concluded.
func (m *Model) Value() (*Choice, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	if len(m.currentChoices) == 0 {
		return nil, fmt.Errorf("no choices")
	}

	if m.currentIdx < 0 || m.currentIdx >= len(m.currentChoices) {
		return nil, fmt.Errorf("choice index out of bounds")
	}

	return m.currentChoices[m.currentIdx], nil
}

// Update updates the model based on the received message.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.Err != nil {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case keyMatches(msg, m.KeyMap.Abort):
			m.Err = common.ErrAborted
			m.quitting = true

			return m, tea.Quit
		case keyMatches(msg, m.KeyMap.Select):
			if len(m.currentChoices) == 0 {
				return m, nil
			}

			m.quitting = true

			return m, tea.Quit
		case keyMatches(msg, m.KeyMap.ClearFilter):
			m.filterInput.Reset()
			m.currentChoices, m.availableChoices = m.filteredAndPagedChoices()
		case keyMatches(msg, m.KeyMap.Down):
			m.cursorDown()
		case keyMatches(msg, m.KeyMap.Up):
			m.cursorUp()
		case keyMatches(msg, m.KeyMap.ScrollDown):
			m.scrollDown()
		case keyMatches(msg, m.KeyMap.ScrollUp):
			m.scrollUp()
		case keyMatches(msg, m.KeyMap.Left):
			m.scrollUp()
		case keyMatches(msg, m.KeyMap.Right):
			m.scrollDown()
		default:
			return m.updateFilter(msg)
		}

		return m, nil
	case tea.WindowSizeMsg:
		m.width = zeroAwareMin(msg.Width, m.MaxWidth)

		if m.height != msg.Height {
			m.height = msg.Height
			m.forceUpdatePageSizeForHeight()
		}

		return m, nil
	case error:
		m.Err = msg

		return m, tea.Quit
	}

	var cmd tea.Cmd

	return m, cmd
}

func (m *Model) forceUpdatePageSizeForHeight() {
	if m.requestedPageSize == 0 {
		m.PageSize = len(m.Choices)
	} else {
		m.PageSize = min(len(m.Choices), m.requestedPageSize)
	}

	for m.PageSize > 1 {
		m.currentIdx = 0
		m.scrollOffset = 0
		m.currentChoices, m.availableChoices = m.filteredAndPagedChoices()

		if len(strings.Split(m.View(), "\n")) <= m.height {
			break
		}

		m.PageSize--
	}
}

func (m *Model) updateFilter(msg tea.Msg) (*Model, tea.Cmd) {
	if m.Filter == nil {
		return m, nil
	}

	previousFilter := m.filterInput.Value()

	var cmd tea.Cmd
	m.filterInput, cmd = m.filterInput.Update(msg)

	if m.filterInput.Value() != previousFilter {
		m.currentIdx = 0
		m.scrollOffset = 0
		m.currentChoices, m.availableChoices = m.filteredAndPagedChoices()
	}

	return m, cmd
}

// View renders the selection prompt.
func (m *Model) View() string {
	viewBuffer := &bytes.Buffer{}

	if m.quitting {
		view, err := m.resultView()
		if err != nil {
			m.Err = err

			return ""
		}

		return m.wrap(view)
	}

	// the cursor only needs to be displayed correctly
	cursor := common.RenderedText(m.Cursor, "red")
	var data string
	// template functions may be displayed dynamically at the head, tail and data area
	// of the list, and a dynamic index(globalDynamicIndex) needs to be added
	for i, obj := range m.pageData {
		// cursor prefix (selected lines need to be displayed,
		// non-selected lines need not be displayed)
		var cursorPrefix string
		// the rendering style of each row of data (the rendering color
		// of selected rows and non-selected rows is different)
		var dataLine string
		// consider three cases when calculating globalDynamicIndex:
		//
		// first page: pageIndex(real time page index)、index(global real time index) keep the two consistent
		//	  1. feat (Introducing new features)
		//	  2. fix (Bug fix)
		//	  3. docs (Writing docs)
		//	  4. style (Improving structure/format of the code)
		//	  5. refactor (Refactoring code)
		//	» [6] test (When adding missing tests)
		//
		// slide down to page: pageIndex fixed to maximum, index increasing with sliding
		//	  2. fix (Bug fix)
		//	  3. docs (Writing docs)
		//	  4. style (Improving structure/format of the code)
		//	  5. refactor (Refactoring code)
		//	  6. test (When adding missing tests)
		//	» [7] chore (Changing CI/CD)
		//
		// swipe up to page: pageIndex fixed to minimum, index decrease with sliding
		//	» [3] docs (Writing docs)
		//	  4. style (Improving structure/format of the code)
		//	  5. refactor (Refactoring code)
		//	  6. test (When adding missing tests)
		//	  7. chore (Changing CI/CD)
		//	  8. perf (Improving performance)
		//
		// in three cases, `m.index - m.pageIndex = n`, `n` is the distance between the global real-time
		// index and the page real-time index. when traversing the page data area, think of the traversal
		// index i as a real-time page index pageIndex, `i + n =` i corresponding global index
		globalDynamicIndex := i + (m.index - m.pageIndex)
		if i == m.pageIndex {
			// keep a space between the cursor and the selected data style
			cursorPrefix = cursor + " "
			// m: A copy of the current object and pass it to the user-defined rendering function to facilitate
			//    the user to read some state information for rendering
			//
			// obj: The single data currently traversed to the data area; pass it to the user-defined rendering
			//      function to help users know the current data that needs to be rendered
			//
			// globalDynamicIndex: The global data index corresponding to the current traverse data; pass it
			//                     to the user-defined rendering function to help users achieve rendering
			//                     actions such as adding serial numbers
			dataLine = m.SelectedChoiceStyle(m, obj, globalDynamicIndex) + "\n"
		} else {
			// the cursor is not displayed on the unselected line, and the selected line is aligned with the blank character
			cursorPrefix = common.GenSpaces(runewidth.StringWidth(m.Cursor) + 1)
			dataLine = m.UnselectedChoiceStyle(m, obj, globalDynamicIndex) + "\n"
		}
		data += cursorPrefix + dataLine
	}
	//fmt.Println(data)
	// avoid panics if Quit is sent during Init
	if m.tmpl == nil {
		return ""
	}

	err := m.tmpl.Execute(viewBuffer, map[string]interface{}{
		"Prompt":        defaultHeader + "\n" + m.Prompt,
		"IsFiltered":    m.Filter != nil,
		"FilterPrompt":  m.FilterPrompt,
		"FilterInput":   m.filterInput.View(),
		"Choices":       m.currentChoices,
		"NChoices":      len(m.currentChoices),
		"SelectedIndex": m.currentIdx,
		"PageSize":      m.PageSize,
		"IsPaged":       m.PageSize > 0 && len(m.currentChoices) > m.PageSize,
		"AllChoices":    m.Choices,
		"NAllChoices":   len(m.Choices),
		"TerminalWidth": m.width,
	})
	if err != nil {
		m.Err = err

		return "Template Error: " + err.Error()
	}

	return m.wrap(viewBuffer.String())
}

func (m *Model) resultView() (string, error) {
	viewBuffer := &bytes.Buffer{}

	if m.ResultTemplate == "" {
		return "", nil
	}

	if m.resultTmpl == nil {
		return "", fmt.Errorf("rendering confirmation without loaded template")
	}

	choice, err := m.Value()
	if err != nil {
		return "", err
	}

	err = m.resultTmpl.Execute(viewBuffer, map[string]interface{}{
		"FinalChoice":   choice,
		"Prompt":        m.Prompt,
		"AllChoices":    m.Choices,
		"NAllChoices":   len(m.Choices),
		"TerminalWidth": m.width,
	})
	if err != nil {
		return "", fmt.Errorf("execute confirmation template: %w", err)
	}

	return viewBuffer.String(), nil
}

func (m *Model) wrap(text string) string {
	if m.WrapMode == nil {
		return text
	}

	return m.WrapMode(text, m.width)
}

func (m *Model) filteredAndPagedChoices() ([]*Choice, int) {
	choices := []*Choice{}

	var available, ignored int

	for _, choice := range m.Choices {
		if m.Filter != nil && !m.Filter(m.filterInput.Value(), choice) {
			continue
		}

		available++

		if m.PageSize > 0 && len(choices) >= m.PageSize {
			break
		}

		if (m.PageSize > 0) && (ignored < m.scrollOffset) {
			ignored++

			continue
		}

		choices = append(choices, choice)
	}

	return choices, available
}

func (m *Model) canScrollDown() bool {
	if m.PageSize <= 0 || m.availableChoices <= m.PageSize {
		return false
	}

	if m.scrollOffset+m.PageSize >= len(m.Choices) {
		return false
	}

	return true
}

func (m *Model) canScrollUp() bool {
	return m.scrollOffset > 0
}

func (m *Model) cursorDown() {
	if m.currentIdx == len(m.currentChoices)-1 && m.canScrollDown() {
		m.scrollDown()
	}

	m.currentIdx = min(len(m.currentChoices)-1, m.currentIdx+1)
}

func (m *Model) cursorUp() {
	if m.currentIdx == 0 && m.canScrollUp() {
		m.scrollUp()
	}

	m.currentIdx = max(0, m.currentIdx-1)
}

func (m *Model) scrollDown() {
	if m.PageSize <= 0 || m.scrollOffset+m.PageSize >= m.availableChoices {
		return
	}

	m.currentIdx = max(0, m.currentIdx-1)
	m.scrollOffset++
	m.currentChoices, m.availableChoices = m.filteredAndPagedChoices()
}

func (m *Model) scrollUp() {
	if m.PageSize <= 0 || m.scrollOffset <= 0 {
		return
	}

	m.currentIdx = min(len(m.currentChoices)-1, m.currentIdx+1)
	m.scrollOffset--
	m.currentChoices, m.availableChoices = m.filteredAndPagedChoices()
}

func (m *Model) reindexChoices() {
	for i, choice := range m.Choices {
		choice.Index = i
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func zeroAwareMin(a int, b int) int {
	switch {
	case a == 0:
		return b
	case b == 0:
		return a
	default:
		return min(a, b)
	}
}