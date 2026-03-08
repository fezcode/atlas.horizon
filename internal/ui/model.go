package ui

import (
	"fmt"
	"time"

	"atlas.horizon/internal/api"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Location    *api.Location
	Weather     *api.WeatherData
	Err         error
	Loading     bool
	Spinner     spinner.Model
	Width       int
	Height      int
	ShowDetails bool
	Quitting    bool
	LastRefresh time.Time
}

type (
	LocMsg   *api.Location
	WeatherMsg *api.WeatherData
	ErrMsg   error
)

func NewModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(Gold)
	
	return Model{
		Loading: true,
		Spinner: s,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.Spinner.Tick,
		m.GetInitialData(),
	)
}

func (m Model) GetInitialData() tea.Cmd {
	return func() tea.Msg {
		loc, err := api.GetLocation()
		if err != nil {
			return ErrMsg(err)
		}
		
		weather, err := api.GetWeather(loc.Lat, loc.Lon, loc.Timezone)
		if err != nil {
			return ErrMsg(err)
		}
		
		m.Location = loc // This won't work in a pure function but we return both
		return tea.Batch(
			func() tea.Msg { return LocMsg(loc) },
			func() tea.Msg { return WeatherMsg(weather) },
		)()
	}
}

// Improved Data Fetcher
func FetchData() tea.Cmd {
	return func() tea.Msg {
		loc, err := api.GetLocation()
		if err != nil {
			return ErrMsg(err)
		}
		
		weather, err := api.GetWeather(loc.Lat, loc.Lon, loc.Timezone)
		if err != nil {
			return ErrMsg(err)
		}
		
		return DataBundle{loc, weather}
	}
}

type DataBundle struct {
	Loc *api.Location
	W   *api.WeatherData
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.Quitting = true
			return m, tea.Quit
		case "r":
			m.Loading = true
			return m, tea.Batch(m.Spinner.Tick, FetchData())
		case "h":
			m.ShowDetails = !m.ShowDetails
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case DataBundle:
		m.Location = msg.Loc
		m.Weather = msg.W
		m.Loading = false
		m.LastRefresh = time.Now()
		return m, nil

	case ErrMsg:
		m.Err = msg
		m.Loading = false
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	if m.Quitting {
		return ""
	}

	if m.Err != nil {
		return AppStyle.Render(fmt.Sprintf("❌ Error: %v\n\nPress 'q' to quit.", m.Err))
	}

	if m.Loading {
		return AppStyle.Render(fmt.Sprintf("%s Fetching environmental data from orbit...", m.Spinner.View()))
	}

	// Main Layout
	title := fmt.Sprintf(" HORIZON-OS v0.1.0 | %s, %s ", m.Location.City, m.Location.Country)
	
	// Current Weather Box
	curr := m.Weather.Current
	icon := api.GetWeatherIcon(curr.WeatherCode, curr.IsDay == 1)
	condition := api.GetWeatherCondition(curr.WeatherCode)
	
	currentContent := fmt.Sprintf(
		"%s %s\n\n"+
		"Temp:     %s %.1f°C\n"+
		"Feels:    %s %.1f°C\n"+
		"Humidity: %s %.0f%%\n"+
		"Wind:     %s %.1f km/h",
		icon, ValueStyle.Render(condition),
		LabelStyle.Render(""), curr.Temperature,
		LabelStyle.Render(""), curr.ApparentTemp,
		LabelStyle.Render(""), curr.Humidity,
		LabelStyle.Render(""), curr.WindSpeed,
	)

	// Rad-Meter (UV Index)
	uv := m.Weather.Daily.UVIndex[0]
	radLevel := "SAFE"
	radStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	
	if uv > 3 && uv <= 6 {
		radLevel = "CAUTION"
		radStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00"))
	} else if uv > 6 {
		radLevel = "DANGER"
		radStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true)
	}

	radContent := fmt.Sprintf(
		"UV Index:  %s %.1f\n"+
		"Rad-Level: %s\n\n"+
		"%s",
		ValueStyle.Render(""), uv,
		radStyle.Render(radLevel),
		m.drawRadBar(uv),
	)

	// Compass / Wind Box
	windContent := m.drawWindRadar(curr.WindDirection, curr.WindSpeed)

	// Assemble boxes
	row1 := lipgloss.JoinHorizontal(lipgloss.Top,
		DrawBox("SURFACE CONDITIONS", currentContent, 30),
		lipgloss.NewStyle().MarginLeft(2).Render(DrawBox("RADIATION MONITOR", radContent, 30)),
	)

	row2 := lipgloss.NewStyle().MarginTop(1).Render(
		lipgloss.JoinHorizontal(lipgloss.Top,
			DrawBox("WIND RADAR", windContent, 30),
			lipgloss.NewStyle().MarginLeft(2).Render(m.drawForecastBox()),
		),
	)

	s := lipgloss.JoinVertical(lipgloss.Left,
		TitleStyle.Render(title),
		row1,
		row2,
		HelpStyle.Render("r: refresh • h: toggle stats • q: quit"),
	)

	return AppStyle.Render(s)
}

func (m Model) drawRadBar(uv float64) string {
	width := 24
	filled := int(uv / 11.0 * float64(width))
	if filled > width { filled = width }
	if filled < 1 { filled = 1 }

	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	
	style := lipgloss.NewStyle().Foreground(Bronze)
	if uv > 6 {
		style = style.Foreground(Red)
	}
	
	return style.Render("[" + bar + "]")
}

func (m Model) drawWindRadar(dir int, speed float64) string {
	// Simple ASCII Compass
	directions := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	idx := int((float64(dir) + 22.5) / 45.0) % 8
	
	compass := fmt.Sprintf(
		"    N\n"+
		"W   +   E\n"+
		"    S\n\n"+
		"Heading: %s (%d°)\n"+
		"Speed:   %.1f km/h",
		ValueStyle.Render(directions[idx]), dir, speed,
	)
	return compass
}

func (m Model) drawForecastBox() string {
	content := ""
	for i := 1; i < 4; i++ {
		t := m.Weather.Daily.Time[i]
		date, _ := time.Parse("2006-01-02", t)
		icon := api.GetWeatherIcon(m.Weather.Daily.WeatherCode[i], true)
		max := m.Weather.Daily.TempMax[i]
		min := m.Weather.Daily.TempMin[i]
		
		content += fmt.Sprintf("%s %s %s %.0f/%.0f°C\n", 
			LabelStyle.Render(date.Format("Mon")),
			icon,
			ValueStyle.Render(""), max, min,
		)
	}
	return DrawBox("72H PROJECTION", content, 26)
}
