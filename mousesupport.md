```go
func main() {
    // ... existing code ...

    p := tea.NewProgram(&wrapper,
        tea.WithAltScreen(),
        tea.WithMouseAllMotion(),    // Enable mouse motion events
        tea.WithMouseCellMotion(),   // Enable cell-based mouse events
    )

    // ... rest of main ...
}
```

2. Update your ModelWrapper to handle mouse events:

```go
// photoboard/cmd/main.go
func (m ModelWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.MouseMsg:
        switch msg.Type {
        case tea.MouseLeft:
            // Handle left click
            return m.handleMouseClick(msg)
        case tea.MouseRight:
            // Handle right click
            return m.handleRightClick(msg)
        case tea.MouseWheelUp:
            return m.handleMouseWheel(-1)
        case tea.MouseWheelDown:
            return m.handleMouseWheel(1)
        }
    }
    // ... rest of Update method ...
}

func (m ModelWrapper) handleMouseClick(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
    if m.model.ShowPreview {
        // Handle clicks in preview mode
        return m.handlePreviewClick(msg)
    }

    // Handle clicks in list mode
    if msg.Y >= 1 && msg.Y <= len(m.model.List.Items()) {
        m.model.List.Select(msg.Y - 1)
        // Simulate enter key press for double clicks
        if msg.Type == tea.MouseLeft && time.Since(m.lastClick) < 500*time.Millisecond {
            return m.model.Update(tea.KeyMsg{Type: tea.KeyEnter}, m.nav, m.prev, m.ui)
        }
        m.lastClick = time.Now()
    }
    return m, nil
}

func (m ModelWrapper) handleRightClick(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
    if !m.model.ShowPreview {
        // Start preview on right click
        if msg.Y >= 1 && msg.Y <= len(m.model.List.Items()) {
            m.model.List.Select(msg.Y - 1)
            return m.model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")}, m.nav, m.prev, m.ui)
        }
    }
    return m, nil
}

func (m ModelWrapper) handleMouseWheel(direction int) (tea.Model, tea.Cmd) {
    if m.model.ShowPreview {
        // Scroll preview
        if direction < 0 {
            m.model.Viewport.LineUp(1)
        } else {
            m.model.Viewport.LineDown(1)
        }
    } else {
        // Scroll list
        if direction < 0 {
            m.model.List.CursorUp()
        } else {
            m.model.List.CursorDown()
        }
    }
    return m, nil
}

func (m ModelWrapper) handlePreviewClick(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
    // Add preview-specific click handling here
    // For example, clicking could exit preview mode
    if msg.Type == tea.MouseLeft {
        m.model.ShowPreview = false
        return m, nil
    }
    return m, nil
}
```

3. Add necessary fields to ModelWrapper:

```go
// photoboard/cmd/main.go
type ModelWrapper struct {
    model     model.Model
    nav       model.Navigator
    prev      model.Previewer
    ui        model.UIHandler
    lastClick time.Time  // For handling double clicks
    width     int
    height    int
}
```

4. Add visual feedback for mouse hover (optional):

```go
// photoboard/internal/ui/styles.go
type DefaultStyler struct {
    // ... existing fields ...
    hoverStyle lipgloss.Style
}

func NewDefaultStyler() *DefaultStyler {
    return &DefaultStyler{
        // ... existing initialization ...
        hoverStyle: lipgloss.NewStyle().
            Background(lipgloss.Color("205")).
            Foreground(lipgloss.Color("0")),
    }
}
```

5. Update the list delegate to show hover effects:

```go
// photoboard/cmd/main.go
func main() {
    // ... existing code ...

    delegate := list.NewDefaultDelegate()
    delegate.SetSpacing(1)
    delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
        BorderForeground(lipgloss.Color("205"))

    m := model.Model{
        List:        list.New(items, delegate, 0, 0),
        CurrentPath: dir,
        Styler:      styler,
    }

    // ... rest of main ...
}
```

6. Add mouse region tracking (optional, for more complex interactions):

```go
// photoboard/internal/model/model.go
type Model struct {
    // ... existing fields ...
    mouseRegions []string
}

func (m Model) View() string {
    if m.ShowPreview {
        return m.renderPreviewWithMouseRegions()
    }
    return m.renderListWithMouseRegions()
}

func (m Model) renderListWithMouseRegions() string {
    // Add clickable regions to the list items
    var s strings.Builder
    for i, item := range m.List.Items() {
        region := fmt.Sprintf("item-%d", i)
        s.WriteString(tea.MouseArea(item.String(), region))
        s.WriteString("\n")
    }
    return m.Styler.ListStyle().Render(s.String())
}
```

This implementation provides:
- Left click to select items
- Double click to open/enter
- Right click to preview
- Mouse wheel scrolling
- Visual feedback on hover (optional)
- Clickable regions (optional)

You can customize the mouse behavior further by:
- Adding context menus on right click
- Implementing drag and drop
- Adding scroll bars
- Adding clickable buttons or icons
- Adding tooltips on hover

Let me know if you'd like to implement any of these additional features or if you need help with specific mouse interactions!
