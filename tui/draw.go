package tui

import (
	"aporia/ansi"
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

const _asciiArt = `clllc'.                                            .',,,,,;,,:,............''..'',,;:ccccc'..........';:llcll:;cllllllllllllllllllll:;;;;;;;;;;;,'..........','''''''................;;;;,;;;;,.............................................,:::cccccccccc
0000x,............................................':looooooldxl'.........';clllllodxkOOOOx;........';cdk0000OxkO0000000000000000000OOxdooooooolllc;'........:dxoolllllc;............'lxxdooololc,...........................................:dOkkOOOOOOOOO
000Ol.................................... .......;looooooodxkd;........';clooooooxkOOOOOkl'......';loxkOOO0OkxO0000000000000000000000kxdoolllllloolc,.......,okkdoooooooc,...........;dkxolllllol:'.........................................'lkkkOOOOOOOOO
000x,..........................................;cloooooooodxdc'.......,:loooooodkOO00OO0d,......,:ldkO0OO00kxO00000000000000000000000Okdolooooooooool,.......:xOkxoooooooc,..........'okkkxolclcll:,.........................................:xOkOOOOOOOOO
00Oc.........................................,coolooollodxkOkc......';coooooodxkOO0OOO0kc.....':ldxOO00OO0Oxk0000000000000000000000000Okxddoooooooodxdc'.....'oOOOkdooooool;..........ckOOOkdooollc:;'.......................................,okOOOOOOOOOO
00x;.......................................':looollcloxOO00Od,.....':looooodxkOO000OOO0d;...';codk00000000kxO00000000000000000000000000OkxdooddoooooxOkl,.....:k0OOkdoooooooc,........;xOO00Oxdooooolc,......................................'cxkkOOOOOOOO
00l.......................................;loollllloxkO00O0Ol'....,coooooodxOOOO00OO00Ol'..,cooxk00O00000OxkO0000000000000000000000000000kxooodoooodkOOOo,....,dO000kxdoooooooc'......'oOOO00Okdooooool:'.....................................;dkkkOOOOOOO
0k:.....................................,coolcllooxkOOOO0O0k:'..';looooodxkOOOOO00000Oxc.':lodxO000000000kxO000000000000000000000000000000kxdoooddodkO00Oo;....ck0O00Okdoododkkl,......ck0O000Okxdooooooc;....................................'lkOOOOOOOOO
0x,...................................':lllllloodkOO0000000d;..':looooodkO00OOkxddllcc;'.,:cclxO000OO000OkkO000000000000000000000000000000OkxdoooodxkO000Od:...'colccc:;;;;;;codl,.....;xO000000Oxdooooool:'...................................:xOOOOOOOOO
0d'...................................:oollooodxkO000000O0Oo'.':loooodooolc;;,,,,,,,,,'',,,'''';:coxOO00OkO00000000000000000000000000000000OkdoooodxO00Oxoc,.....',;;;:;::;;,,,,,......,dO000000OOkdoooooooc,..................................;dOOOOOOOOO
Ol....................................:oooooodkOO00O000000Oc.':oool:;,,,;::clodxxxxxdlccloodddolc;;;:lkOkkO0O00000000000000000000000000000000kdodookOxc,'';cc;..:dkkkxxxxxxkkkkxdol;....;okO0000000kxdooooool:'................................,okOOOOOOOO
Oc....................................:oooooxkOOO000000OO0kc,;c:,,,;cldxkkkkkdl:;'..........',:lxkkxdxOOkk00O00000000000000000000000000000000OkdoodkOxloxxkOOx:.';;,''....'',;cldkkxc'..',,;lkOOO000Oxdoooooooc,...............................'okkOOOOOOO
k:....................................;loodxOOOOO00O0000O0k:'''':oxkkkkkkkxl;.... .........    ..:oxkOOkxk000000000000000000000000000000000000OxdodO000Okkkxo:'....           ...,:ddc''cdo:,,:dkO000Oxdoooooool;'.............................'lkOkOOOOOO
k;....................................;lodkO0OO00OO000000Ol..,ldkkkkkkkkxc'..............        ..:dkOOOO0000000000000000000000000000000000000kxdxO00OOOxl,...  ..               ..,c:,cxkkxl;',lxO000kdoooooool;..............................ckOkOOOOOO
x,....................................,lxkO00OOOOO000OOOo,.'lxkOkkkkkkko,. ..  ..........          .:dkO000000000000000000000000000000000000000OkxO000OOx:. .....           .     .  ..'cxkkOkxl,.'ck0OOkdooooooo;..............................:xkkOOOOOO
d'....................................;dkOOOO00000OO0Ox:.'cxkkkOkOOO0Oo'... ............         .  .ck00OO00000000000000000000000000000000000000O00O0Kk:... ....           .          .;dkkkkkkxc..,ok00kdoooool,..............................:xkOOOOOOO
d'....................................;dO0OO000O00O0Oo,.,okkkOOO00KKKo'...............              .'d0000000000000000000000000000000000000000000000KOc.......                         .;dkkkkkkko,..cxOOkdooooc...............................;xkOOOOOOO
o.....................................;dO0OO000OOO0kc..;dkOkk00KKKKKd,.  ..........                  .;kK00O000000000000000000000000000000000000000000d'.......                          .;xOOkkkkkd;..;dOOkdooo:...............................;dkOOOOOOO
o.....................................,oOO0OO00OOOx:..;dkkkO0KKXXKKk;... ..........               .  ..l00OO0000000000000000000000000000000000000000K0c. ......                           .cO0OOkkkkd;..;dOOkdoo;...............................;dkOOOOOOO
o.....................................'lkO0OOO0OOx:..,okkkOKKKXXXX0c. .. ..........            .  . ...;k0000OO0000000000000000000000000000000000000Kk;........                            'd0K0Okkkko,..:x0Odol,...............................;dkOOOOOkO
o......................................ckOO0OOOOkc..'lkkkO0KKKKKKKk,...  ..........                 .. .d0000000000000000000000000000000000000000000Kx'..........                          .:OXK0Okkkxl'.'lOOddl'...............................,dkkOOOOOO
l......................................:k000OO0Oo'..:xkkOKKKKKKKKKd'.... ..........                 .. .lO000000000000000000000000000000000000000000Ko...........                          .,kKKKKOkkkx:..:kOdoc................................,oOOOOOOOO
l......................................;dO00O00kc..'okkO0KXKKKKKKKo. ... ..........                 .. .:O00000000000000000000000000000000000000000K0l...........                           'xKKXK0kkkkl..,dOxo:................................'oOOOOOOOO
l......................................,oO0OO00x,..;dkkO0KKKKKKKXKl..... ..........                 .. .:k00O00000000000000000000000000000000000000K0c...........                        .. 'xKKKXKOkkkd,..lkxo:................................'lkOOOOOOO
l......................................'lkOOO00d. .:xkkO0KKKKXXKXKl................             ..  .  .:O0O000000000000000000000000000000000000000KO:............                      ....'dKKKXK0Okkx;..lkxo;.................................ckkOOOOOO
l......................................'ckOOO0Oo. .:xkkO0KXKKXXKKKd.... ............         .  ..  .. .cO0OO00000000000000000000000000000000000000KO:....  ......                  ..     .,kKKKKXKOkOx:..lkxo;.................................ckkOOOOOO
l.......................................:kO0O0Oo...;xkkk0KKKKXXKKKx,................         .;lodl:'. .l000000000000000000000000000000000000000000K0l....  ......           .  .;lddol;.. .:OXKXXXKOkOx:..okdl,.................................:kOOOOOOO
o.......................................;xO0O00d'..,dkkkOKKKXXKXKK0c................    .. .,d0XXKKKx;.,x0000000000000000000000000000000000000000000Kd....  ......          ...,x0XXXXK0d,.'dKXKXXXKOOkx;.'dkdl,.................................:xOOOOOOO
l.......................................;xO0O00x;..'lkkkO0KXXXKXKKXk;.................  .. .l0KKKXXKKd,cO0000000000000000000000000000000000000000000Kk,...  ......          . .oKKKKKXKK0l'c0XKXXXXKOkOo'.;xkdl'.................................;xOOOOOOO
l.......................................;d00O00Od'..,dkkk0KXXXKKXXKKx'................. .. .:OKKKKXK0l:xK000000000000000000000000000000000000000000000l. .   .....         .. .o0KKKKXKKOlcOKKKXXXX0OOkc..lOkdc'.................................;xOOOOOOO
l.......................................;d00O00OOo'..:xkkO0KXXXKXXKKKd'..  ............  .  .:kKKKX0o;o0000000000000000000000000000000000000000000000Kk;.....   ..          . .,oOKKKKOxloOKKKKKKXKOOkl,';x0xo:..................................;dOOOOOOO
c.......................................;d00O00O0Oo,..:dkkO0KXXXKKXKKKx,.. ... .........     .':clc,'cO0000000000000000000000000000000000000000000000KKk;..  ..  .           .  .';::;,;oOKKXXKKKK0Oxoodkk0Oxo:..................................;dOOOOOOO
l.......................................;dOOO0OOO00xc;:oxkOO0KXXXXXXXXKkl'...  ..........          .lO000000000000000000000000000000000000000000000000KKOo,.......   ..      ..      .ckKKKXKXXXK00OkO00000Oxo:..................................,dOkOOOOO
l.......................................;xOO0OOOOO00OOOOOOOkOKKKKXKXXXXKKkc...  .............  ...ck0000000000000000000000000000000000000000000000000000KKOd:'.........          .':oOKKKKKKKKKKK00O0000000Oxo:..................................,dOkOOOOO
l.......................................;xO0000O000000000000000KKKKKXXXXXXKOxc;'.............';cdk0K0000000000000000000000000000000000000000000000000000O0KK0Oxo:;'..........';cokOKKXXKKKKKKK0000000000000Oxo:..................................,oOOOOOOO
l.......................................;xO00000OO00000000000000000KKKKKXKKKXXK0OOkxxdddddxxk0KKKK000000000000000000000000000000000000000000000000000000000000KKK00OkxxxxxkkO0KXXXXXXKKKKKK0000000000000000Odo;..................................'oOOOOOOO
l.......................................;xOO00000OO000000000000000000000000KKKKKKKKKKKKKKKKKKK00000000000000000000000000000000000000000000000000000000000000000000000KK0KKKKKKKKKKK000000000000000000000000kdo;..................................'lOOOOOOO
l.......................................:kOO0OO000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000kdo,..................................'lkOOOOOO
l.......................................:k0O0OOO0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000Oxdo,...................................lkOOOOOO
l.......................................:k0O0OO00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000Oxdl,...................................ckOOOOOO
o.......................................ck0O0O000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000Oxdl,...................................ckOOOOOO
d'......................................:kOO0O000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000Oxdl'...................................ckOOOOOO
x,......................................:kOO00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000kdoc'...................................ckOOOOOO
k;......................................:kOO00000000000000000000000000000000000000000000000000OkxxxddxxxxxxkxxxxxxxxxxxxxxkkkkkkkkkkkkkkOOOOOOOOOOOOOOO00000000000000000000000000000000000000000000000000Okdoc....................................ckOOOOOO
k;......................................;xOO00O00000000000000000000000000000000000000000000000k:......,;:::::::::::::::::::::::::::::::ccccccccccccccclxO000000000000000000000000000000000000000000000000Oxoo:...................................'lkOOOOOO
O:......................................,oO0OOOO0000000000000000000000000000000000000000000000d'.....,:::::::::::::::::::::::::::::::::::::::::::::::::lx00O000000000000000000000000000000000000000000000Oxdo;...................................'oOOOOOOO
Oc.......................................ck0O00O000000000000000000000000000000000000000000000Ol.....;::::::::::::ccccccccccc:::::::::::::::::::::;::::::ok0OO00000000000000000000000000000000000000000000Oxol'...................................'oOOOOOOO
0l.......................................,dOO00O000000000000000000000000000000000000000000000Oc...';::ccllooddddddddddddddddooooolllcc::::::::::::::::::cxO0000000000000000000000000000000000000000000000Odo:....................................,oOOOOOOO
0l........................................:xO00O000000000000000000000000000000000000000000000O:..';loodddxxxxddddddddddddddddxxxxxxddoolc::::::::::::::;:dO0000000000000000000000000000000000000000000000kdc'....................................,oOOOOOOO
0o.........................................:xO00000000000000000000000000000000000000000000000Oc.'codxddxdxxxxxxxdddddddxxxxxxxxxddxxxxxddoolc::::::::::::oO0000000000000000000000000000000000000000000000kc'.....................................,oOOOOOOO
0d'.........................................:xOO000000000000000000000000000000000000000000000Ol,:dxdddddddxxxxxxxxxxxxxxxxxxxxxxxxxdddddddddolc::::::::::oO00000O000000000000000000000000000000000000000Ol'......................................,dOOOOOOO
Kx'..........................................,dO0000000000000000000000000000000000000000000000o:lxxdxddxxxxxxxdxxxxxxxxxxxxxxxxxddddxxxxdddxxddolc:::::::oOO0000000000000000000000000000000000000000O00kc........................................,dOOOOOOO
Kk,............................................:dk0000O00O000000000000000000000000000000000000xloxxddxxxxxxxxddxxxxxxxxxxxxxxxxxdddxxxxdxxddddxxdolc:::;:oOO0O000000000000000000000000000000000000O00ko,.........................................,dOOOOOOO
Kk;..............................................;okO0OOO000OO00000000000000000000000000000000kddxxddxxdxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxdddolc::cdOO00000000000000000000000000000000000000Oko,...........................................,dOOOOOOO
KO:................................................'cdkOOOOO0000O000000000000000000000000000000xddxdxxxxxxxxxxxxxxxxxxxxxxxxxxxxdddddddddddddddddddddoc:lkOO000000000000000000000000000000000000kdc'.............................................,dOOOOOOO
K0c...................................................,cdk0OOOOOOOO0000000000000000000000000O00kdddxxxxxxxxxxxxxxxxxxxxxxxxxxxxxddddddddddddddddddddddoloO0O00000000000000000000000000000000Oko:,................................................,dOOOOOO0
K0l......................................................,coxO00OO00000000000O00000000000000000Oxddddxxdddddddddddddddddddddddddddddddddddddddddddddxdook0000000000000000000000000000000Oxdl;...................... .............................;xOOOOOOO
K0o..........................................................;cdxOO00O000O0000000000000000000000OxodxdddddddddddddddddddddddddddddddddddddddddddxddxdodkO000000000000000000000000000Oxoc,.......................... .............................;xOOOOOO0
K0x'.............................................................,:oxkOO000O000000000000000000000Ododdxxddddddddddddddddddddddddddddddddddddddddddxdlok000000000000000000000000Okxo:,.... ......................... .............................;xOOOOOOO
0Kk,.................................................................';coxOOO0OO000000000000000000Oxddxxdddddddddddddddddddddddddddddddddddddddddddodk00000000000000000000Okxoc;'..  .....  ....................... .............................:xOOOOOOO
0KO;......................................................................,:ldxkO000OO00000OO0000O0OkdodddddddddxxxxxxxddddddddddddddddddddddddddddxO0000000000000000Okxoc;'.....................................................................:kOOOOOO0
0K0c...........................................................................';coxkOO0000OO00000000Oxoodddddddxxxxxxxdddddddddddddddddddddddddxkk000OO0000000OOkdl:;'..........................................................................:kOOOOOO0
0K0o............................................................................ ....,:ldxOOOO000OO0000OxoooddddddddddddddddddddddddddddddddddxkO000000000kkkxc;,...........   ..................................................................ckOOOOOO0
K0Kk,.....................................................................................,codxkOO00OO000OkdoooddddddddddddddddddddddddddddxkOOO000OOkkxddooo:.................................................... ...........  .................ckOOOO000
K0K0c........................................................................... ........ .,cllloddxxkkO0O00OkxdddddddddddddddddddddddxxxkOOOOkxxddoooooooooo:.................................................... ........... ..................lOOOOOOO0
KKKKd'................................................................................... .,cooooollloooodxxkOOOOOOkkxxxxddddddxxxxkkkkkkxddooooooooooooooool,.................................................... .............................'oOOOOOOO0
KKKKx;................................................................................... .'cooooooooooooooolooodddxxkkkkkxxxxxxxxxdddooooooooooooooooooooool'.................................................... .............................,oOOOOOOO0
KKKKO:................................................................................. ....:oooooooooooooooooooooooolllolllllllloooooooooooooooooooooooooool'................................................... .........  ...................,dOOOOO000
KKKK0l......................................................................................coooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooc...................................................  ......... ....................:xOOO0OO00
KKKKKd'.....................................................................................coooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooc........................... ...  ..................................................ckOOO00OO0
KKKKKx,.....................................................................................:oooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooc........................... ... ..................................................'lOOO000O00`


func (self *Tui) Draw() error {

	// Reset cursor position
	if self.shouldBeRedrawn {
		ansi.Clear()
		ascii := newAscii(_asciiArt)
		self.shouldBeRedrawn = false
		ansi.MoveCursor(0, 0)
		ascii.draw(self.TermSize)
		ansi.MoveCursor(0, 0)
		self.drawBox()
	}

	// Draw the message
	if self.lastDrawnMessage != self.message {
		self.lastDrawnMessage = self.message
		ansi.MoveCursor(2, 0)
		ansi.EraseChars(boxWidth)
		self.drawLine(self.message)
	}

	// Draw the currently selected field
	thisLine := self.fields[self.position].draw()

	ansi.MoveCursor(self.position+3, 0)
	ansi.EraseChars(boxWidth)
	self.drawLine(thisLine)

	ansi.MoveCursor(self.position+3, len(thisLine)+2)

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
func (self Tui) drawBox() {
	fmt.Print(tlCorner, strings.Repeat(horizontal, boxWidth-2), trCorner, "\n\r")

	self.drawLine(self.message)

	for _, field := range self.fields {
		self.drawLine(field.draw())
	}

	fmt.Print(blCorner, strings.Repeat(horizontal, boxWidth-2), brCorner, "\n\r")
}

func (self Tui) drawLine(text string) {
	fmt.Print(vertical)
	fmt.Print(text)
	fmt.Print(strings.Repeat(" ", maxInt(boxWidth-2-len(text), 0)))
	fmt.Print(vertical, "\n\r")
}
