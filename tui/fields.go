package tui

type field interface {
	draw() string
	onInput(tui *Tui, symbol []int)
}

type input struct {
	prompt   string
	contents string
}

func (self input) draw() string {
	return self.prompt + ": " + self.contents
}

func (self *input) onInput(tui *Tui, symbol []int) {
	if len(symbol) == 1 {
		if symbol[0] == 10 {
			tui.NextPosition()
		} else if symbol[0] == 127 {
			if len(self.contents) > 0 {
				self.contents = self.contents[:len(self.contents)-1]
			}
		} else {
			self.contents = self.contents + string(rune(symbol[0]))
		}
	}
}

func newInput(prompt string) *input {
	return &input{
		prompt:   prompt,
		contents: "",
	}
}

type button struct {
	name string
}

func (button button) draw() string {
	return button.name
}

func (button *button) onInput(tui *Tui, symbol []int) {

}

func newButton(name string) *button {
	return &button{
		name: name,
	}
}
