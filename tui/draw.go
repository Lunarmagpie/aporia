package tui

import (
	"fmt"
	"strings"
)

const horizontal = "─"
const vertical = "│"

const tlCorner = "┌"
const trCorner = "┐"
const blCorner = "└"
const brCorner = "┘"

const boxHeight = 6
const boxWidth = 30

const _asciiArt = `
                                        ....                                         
                                   :YPB#####BBGP5YJ7~:.                              
                                   .~7?YPG#&@@@@@@@@@&#GY?~.                         
                                           .^!JP#@@@@@@@@@@#GY!.                     
                    .                           .^?P#@@@&&&@@@@&P7:                  
                   :!~^:.                           .!P&@@&&&&&@@@&P!                
                   .77!!~~^.                           :J#@@&&&&&&@@@B?.             
                    !77!~~!~~:.                          .J#@@&&&&&&@@@#?.           
                    ^777!~~~~~^:.......    ..:::^^^^^^^^^^~^5@@&&&##&&@@@B~          
                     !777!~~~~~~~~!!~~~~~~~~!!!~~~~~~~~!!7!  !#@&&####&&@@@Y.        
                     .77!~~~~!!~~~!!!!!!!~~~~~~~~~~~!!7777^   :B@&&#BBB#&&@@B^       
                      :!~~!~?BJ~75!~!!!!!~~^^~~~~~~!77777~     :#@&#BGGB#&&@@#^      
                      ^~~~~~J##BBP!~!~!!!!~~^^~~~~!77777!       ~@@&#BGGB#&&@@&^     
                     ^~!YY!~~7?J7!~!~~~~~~!!~^:^~~!7777!.        5@&#BGPGB#&&@@#.    
                    .~J&BJY!~~~~~~~~!!!!~~~!!~^:^~~!77~.         :&@&#GPPGB#&@@@5    
                    .J@#7B@Y~!!!~!JB#&&P?~!!!!~^:^~!7^            G@&#BP5PGB&&@@@7   
                    :Y@#JBB!~!!~7B@@@B?JG?^~~!~~::^~.             J@&#BP55PB#&@@@B   
                   .~!P#GPJ~!!!~G@@@#YY@@PJ!~~!~:::               J@&#BP5Y5GB&&@@@~  
                   ^!~~~!!!~~!~~G@@@&#5BPY#~~~~^::.               5@&#BPYYYPB#&@@@?  
                   .~~~~~~~~~!!~!?YPBBPPGG?~~~~:::               .#@&#G5YJYPG#&@@@Y  
                    .^!!!~~~~~!!~~~~~!777~~!~~^:.                7@&&BP5JJY5G#&@@@Y  
 !5:                  .^~!!!!!!!!!~~~~~~~~~~^::.                .#@&#BPYJJY5G#&@@@?  
 J@B:                    .:^~~!!!!~~~~~~~~^::.                  Y@&&BG5J?JYPB#&@@@^  
 ^@@B:                       .:~~~~~~~^^^^::::..               J@@&#G5J??J5PB&&@@B   
  P@@#^                        ~~~~~~~~~~~~~~^:::.            J@@&#G5Y???Y5G#&@@@?   
  ^&@@&?                       ~!!!!!!!!!!!!!!~^^::.        :P@@&#G5Y???J5PB&&@@G    
   7@@@@P^                    .~!!!!!!!!!!!!!!!!~~^::..    7&@&#BG5J???JYPB#&@@#:    
    J@@@@&5~                  :!!!!!!!!!!!!!!!!!!!~~^:::.!G@@&#GPYJ???JYPB#&@@&^     
     ?@@@@@@G7:               ~!!!!!!!!!!!!!!!!!!!~~~~~JB@@&#BP5J????Y5GB#&@@#^      
      !&@@@@@@&P7:           .~!!!!!!!!!!!!!!~~~~~!7YG&@@&#BP5J????J5PB#&&@@G:       
       :G@@@@&@@@&BY7^.       ~~~~~~~~~~~~~~~~!7JPB&@@&#BGPYJ????J5PB#&&@@&J         
         ?#@@&&&&&&@@&#G5?!~^:~!!!~~~~!7?JY5GB#&@&&#BBGPYJJ???JY5PB#&&@@@G^          
          .Y&@@&&&####&&@@@@&&########&&&&@&&&#BBGPP5YJ?????JY5GB#&&@@@B7            
            :Y#@@&&###BBBBB############BBBGGPP55YYJJ???JJYY5PGB#&&@@@B7              
              .7G&@@&&##BBGGGGPPPPPPP5555YYYYJJJJJJYYY5PGGB##&&@@@#5~                
                 :?P&@@@&&##BBGGGPPPP555555555PPPGGGBB##&&&@@@@#5!.                  
                    .!YG&@@@@&&&################&&&&&@@@@@&#P?^.                     
                        .~7YPB&&@@@@@@@@@@@@@@@@@@@&&BG5?!:.                         
                             .:^!7JY5PPGBBBBGGPYJ?!^:.                               
`

func (self *Tui) Draw() error {
	// Reset cursor position
	if self.shouldBeRedrawn {
		fmt.Print("\033[H\033[0J")
		ascii := newAscii(_asciiArt)
		self.shouldBeRedrawn = false

		boxHeight := self.drawBox()
		ascii.draw(ascii.calculatePosition(self.TermSize.Cols, self.TermSize.Lines, boxHeight))
	}

	// Draw the message
	if self.lastDrawnMessage != self.message {
		self.lastDrawnMessage = self.message
		moveCursor(2, 0)
		fmt.Print("\033[K")
		self.drawLine(self.message)
	}
	
	// Draw the currently selected field
	thisLine := self.fields[self.position].draw()

	moveCursor(self.position + 3, 0)

	fmt.Print("\033[K")
	self.drawLine(thisLine)

	moveCursor(self.position + 3, len(thisLine) + 2)

	return nil
}

// Draw the vertical margin.
func drawMargin(height int) {
	for i := 0; i < (height-boxHeight)/2; i++ {
		fmt.Print("\n\r")
	}
}

func eraseLine(num int) {
	fmt.Print("\033[", num, "K")
}

// Draw the box. Return the vertical lines taken up.
func (self Tui) drawBox() int {
	fmt.Print(tlCorner, strings.Repeat(horizontal, boxWidth-2), trCorner, "\n\r")

	self.drawLine(self.message)

	for _, field := range self.fields {
		self.drawLine(field.draw())
	}

	fmt.Print(blCorner, strings.Repeat(horizontal, boxWidth-2), brCorner, "\n\r")

	return 3 + len(self.fields)
}

func (self Tui) drawLine(text string) {
	fmt.Print(vertical)
	fmt.Print(text)
	fmt.Print(strings.Repeat(" ", maxInt(boxWidth-2-len(text), 0)))
	fmt.Print(vertical, "\n\r")
}
