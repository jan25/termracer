package data

const defaultContent = `
Ctrl+s		New race
Ctrl+j/k	Scroll
Ctrl+c 		Quit`

const raceModeContent = `
Ctrl+e End race
Ctrl+c Quit`

// ControlsData wraps content for controls view
type ControlsData struct {
	Content string
}

// NewControlsData creates new instance of ControlsData
func NewControlsData() *ControlsData {
	return &ControlsData{
		Content: defaultContent,
	}
}

// DefaultControls sets content for default controls view
func (c *ControlsData) DefaultControls() {
	c.Content = defaultContent
}

// StartRace sets content for controls for race in progresse
func (c *ControlsData) StartRace() {
	c.Content = raceModeContent
}
