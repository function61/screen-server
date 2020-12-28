package evdev

// https://github.com/torvalds/linux/blob/5c8fe583cce542aa0b84adc939ce85293de36e5e/include/uapi/linux/input-event-codes.h#L913
type Led uint16

const (
	LedNUML     Led = 0x00
	LedCAPSL    Led = 0x01
	LedSCROLLL  Led = 0x02
	LedCOMPOSE  Led = 0x03
	LedKANA     Led = 0x04
	LedSLEEP    Led = 0x05
	LedSUSPEND  Led = 0x06
	LedMUTE     Led = 0x07
	LedMISC     Led = 0x08
	LedMAIL     Led = 0x09
	LedCHARGING Led = 0x0a
)

type Rel uint16

const (
	RelX           Rel = 0x00
	RelY           Rel = 0x01
	RelZ           Rel = 0x02
	RelRX          Rel = 0x03
	RelRY          Rel = 0x04
	RelRZ          Rel = 0x05
	RelHWHEEL      Rel = 0x06
	RelDIAL        Rel = 0x07
	RelWHEEL       Rel = 0x08
	RelMISC        Rel = 0x09
	RelRESERVED    Rel = 0x0a
	RelWHEELHIRES  Rel = 0x0b
	RelHWHEELHIRES Rel = 0x0c
)

type Key uint16

const (
	KeyRESERVED   Key = 0
	KeyESC        Key = 1
	Key1          Key = 2
	Key2          Key = 3
	Key3          Key = 4
	Key4          Key = 5
	Key5          Key = 6
	Key6          Key = 7
	Key7          Key = 8
	Key8          Key = 9
	Key9          Key = 10
	Key0          Key = 11
	KeyMINUS      Key = 12
	KeyEQUAL      Key = 13
	KeyBACKSPACE  Key = 14
	KeyTAB        Key = 15
	KeyQ          Key = 16
	KeyW          Key = 17
	KeyE          Key = 18
	KeyR          Key = 19
	KeyT          Key = 20
	KeyY          Key = 21
	KeyU          Key = 22
	KeyI          Key = 23
	KeyO          Key = 24
	KeyP          Key = 25
	KeyLEFTBRACE  Key = 26
	KeyRIGHTBRACE Key = 27
	KeyENTER      Key = 28
	KeyLEFTCTRL   Key = 29
	KeyA          Key = 30
	KeyS          Key = 31
	KeyD          Key = 32
	KeyF          Key = 33
	KeyG          Key = 34
	KeyH          Key = 35
	KeyJ          Key = 36
	KeyK          Key = 37
	KeyL          Key = 38
	KeySEMICOLON  Key = 39
	KeyAPOSTROPHE Key = 40
	KeyGRAVE      Key = 41
	KeyLEFTSHIFT  Key = 42
	KeyBACKSLASH  Key = 43
	KeyZ          Key = 44
	KeyX          Key = 45
	KeyC          Key = 46
	KeyV          Key = 47
	KeyB          Key = 48
	KeyN          Key = 49
	KeyM          Key = 50
	KeyCOMMA      Key = 51
	KeyDOT        Key = 52
	KeySLASH      Key = 53
	KeyRIGHTSHIFT Key = 54
	KeyKPASTERISK Key = 55
	KeyLEFTALT    Key = 56
	KeySPACE      Key = 57
	KeyCAPSLOCK   Key = 58
	KeyF1         Key = 59
	KeyF2         Key = 60
	KeyF3         Key = 61
	KeyF4         Key = 62
	KeyF5         Key = 63
	KeyF6         Key = 64
	KeyF7         Key = 65
	KeyF8         Key = 66
	KeyF9         Key = 67
	KeyF10        Key = 68
	KeyNUMLOCK    Key = 69
	KeySCROLLLOCK Key = 70
	KeyKP7        Key = 71
	KeyKP8        Key = 72
	KeyKP9        Key = 73
	KeyKPMINUS    Key = 74
	KeyKP4        Key = 75
	KeyKP5        Key = 76
	KeyKP6        Key = 77
	KeyKPPLUS     Key = 78
	KeyKP1        Key = 79
	KeyKP2        Key = 80
	KeyKP3        Key = 81
	KeyKP0        Key = 82
	KeyKPDOT      Key = 83

	KeyZENKAKUHANKAKU   Key = 85
	Key102ND            Key = 86
	KeyF11              Key = 87
	KeyF12              Key = 88
	KeyRO               Key = 89
	KeyKATAKANA         Key = 90
	KeyHIRAGANA         Key = 91
	KeyHENKAN           Key = 92
	KeyKATAKANAHIRAGANA Key = 93
	KeyMUHENKAN         Key = 94
	KeyKPJPCOMMA        Key = 95
	KeyKPENTER          Key = 96
	KeyRIGHTCTRL        Key = 97
	KeyKPSLASH          Key = 98
	KeySYSRQ            Key = 99
	KeyRIGHTALT         Key = 100
	KeyLINEFEED         Key = 101
	KeyHOME             Key = 102
	KeyUP               Key = 103
	KeyPAGEUP           Key = 104
	KeyLEFT             Key = 105
	KeyRIGHT            Key = 106
	KeyEND              Key = 107
	KeyDOWN             Key = 108
	KeyPAGEDOWN         Key = 109
	KeyINSERT           Key = 110
	KeyDELETE           Key = 111
	KeyMACRO            Key = 112
	KeyMUTE             Key = 113
	KeyVOLUMEDOWN       Key = 114
	KeyVOLUMEUP         Key = 115
	KeyPOWER            Key = 116
	KeyKPEQUAL          Key = 117
	KeyKPPLUSMINUS      Key = 118
	KeyPAUSE            Key = 119
	KeySCALE            Key = 120

	// used by phones, remote controls, and other keypads
	KeyNUMERIC0     Key = 0x200
	KeyNUMERIC1     Key = 0x201
	KeyNUMERIC2     Key = 0x202
	KeyNUMERIC3     Key = 0x203
	KeyNUMERIC4     Key = 0x204
	KeyNUMERIC5     Key = 0x205
	KeyNUMERIC6     Key = 0x206
	KeyNUMERIC7     Key = 0x207
	KeyNUMERIC8     Key = 0x208
	KeyNUMERIC9     Key = 0x209
	KeyNUMERICSTAR  Key = 0x20a
	KeyNUMERICPOUND Key = 0x20b
	KeyNUMERICA     Key = 0x20c // Phone key A - HUT Telephony 0xb9
	KeyNUMERICB     Key = 0x20d
	KeyNUMERICC     Key = 0x20e
	KeyNUMERICD     Key = 0x20f
)

type Btn uint16

const (
	BtnMISC Btn = 0x100
	Btn0    Btn = 0x100
	Btn1    Btn = 0x101
	Btn2    Btn = 0x102
	Btn3    Btn = 0x103
	Btn4    Btn = 0x104
	Btn5    Btn = 0x105
	Btn6    Btn = 0x106
	Btn7    Btn = 0x107
	Btn8    Btn = 0x108
	Btn9    Btn = 0x109

	BtnMOUSE   Btn = 0x110
	BtnLEFT    Btn = 0x110
	BtnRIGHT   Btn = 0x111
	BtnMIDDLE  Btn = 0x112
	BtnSIDE    Btn = 0x113
	BtnEXTRA   Btn = 0x114
	BtnFORWARD Btn = 0x115
	BtnBACK    Btn = 0x116
	BtnTASK    Btn = 0x117
)

// connects the code with human readable description
// taken from https://github.com/torvalds/linux/blob/master/include/uapi/linux/input-event-codes.h
var keyCodeMap = map[uint16]string{
	0:   "RESERVED",         // KEY_RESERVED
	1:   "ESC",              // KEY_ESC
	2:   "1",                // KEY_1
	3:   "2",                // KEY_2
	4:   "3",                // KEY_3
	5:   "4",                // KEY_4
	6:   "5",                // KEY_5
	7:   "6",                // KEY_6
	8:   "7",                // KEY_7
	9:   "8",                // KEY_8
	10:  "9",                // KEY_9
	11:  "0",                // KEY_0
	12:  "MINUS",            // KEY_MINUS
	13:  "EQUAL",            // KEY_EQUAL
	14:  "BACKSPACE",        // KEY_BACKSPACE
	15:  "TAB",              // KEY_TAB
	16:  "Q",                // KEY_Q
	17:  "W",                // KEY_W
	18:  "E",                // KEY_E
	19:  "R",                // KEY_R
	20:  "T",                // KEY_T
	21:  "Y",                // KEY_Y
	22:  "U",                // KEY_U
	23:  "I",                // KEY_I
	24:  "O",                // KEY_O
	25:  "P",                // KEY_P
	26:  "LEFTBRACE",        // KEY_LEFTBRACE
	27:  "RIGHTBRACE",       // KEY_RIGHTBRACE
	28:  "ENTER",            // KEY_ENTER
	29:  "LEFTCTRL",         // KEY_LEFTCTRL
	30:  "A",                // KEY_A
	31:  "S",                // KEY_S
	32:  "D",                // KEY_D
	33:  "F",                // KEY_F
	34:  "G",                // KEY_G
	35:  "H",                // KEY_H
	36:  "J",                // KEY_J
	37:  "K",                // KEY_K
	38:  "L",                // KEY_L
	39:  "SEMICOLON",        // KEY_SEMICOLON
	40:  "APOSTROPHE",       // KEY_APOSTROPHE
	41:  "GRAVE",            // KEY_GRAVE
	42:  "LEFTSHIFT",        // KEY_LEFTSHIFT
	43:  "BACKSLASH",        // KEY_BACKSLASH
	44:  "Z",                // KEY_Z
	45:  "X",                // KEY_X
	46:  "C",                // KEY_C
	47:  "V",                // KEY_V
	48:  "B",                // KEY_B
	49:  "N",                // KEY_N
	50:  "M",                // KEY_M
	51:  "COMMA",            // KEY_COMMA
	52:  "DOT",              // KEY_DOT
	53:  "SLASH",            // KEY_SLASH
	54:  "RIGHTSHIFT",       // KEY_RIGHTSHIFT
	55:  "KPASTERISK",       // KEY_KPASTERISK
	56:  "LEFTALT",          // KEY_LEFTALT
	57:  "SPACE",            // KEY_SPACE
	58:  "CAPSLOCK",         // KEY_CAPSLOCK
	59:  "F1",               // KEY_F1
	60:  "F2",               // KEY_F2
	61:  "F3",               // KEY_F3
	62:  "F4",               // KEY_F4
	63:  "F5",               // KEY_F5
	64:  "F6",               // KEY_F6
	65:  "F7",               // KEY_F7
	66:  "F8",               // KEY_F8
	67:  "F9",               // KEY_F9
	68:  "F10",              // KEY_F10
	69:  "NUMLOCK",          // KEY_NUMLOCK
	70:  "SCROLLLOCK",       // KEY_SCROLLLOCK
	71:  "KP7",              // KEY_KP7
	72:  "KP8",              // KEY_KP8
	73:  "KP9",              // KEY_KP9
	74:  "KPMINUS",          // KEY_KPMINUS
	75:  "KP4",              // KEY_KP4
	76:  "KP5",              // KEY_KP5
	77:  "KP6",              // KEY_KP6
	78:  "KPPLUS",           // KEY_KPPLUS
	79:  "KP1",              // KEY_KP1
	80:  "KP2",              // KEY_KP2
	81:  "KP3",              // KEY_KP3
	82:  "KP0",              // KEY_KP0
	83:  "KPDOT",            // KEY_KPDOT
	85:  "ZENKAKUHANKAKU",   // KEY_ZENKAKUHANKAKU
	86:  "102ND",            // KEY_102ND
	87:  "F11",              // KEY_F11
	88:  "F12",              // KEY_F12
	89:  "RO",               // KEY_RO
	90:  "KATAKANA",         // KEY_KATAKANA
	91:  "HIRAGANA",         // KEY_HIRAGANA
	92:  "HENKAN",           // KEY_HENKAN
	93:  "KATAKANAHIRAGANA", // KEY_KATAKANAHIRAGANA
	94:  "MUHENKAN",         // KEY_MUHENKAN
	95:  "KPJPCOMMA",        // KEY_KPJPCOMMA
	96:  "KPENTER",          // KEY_KPENTER
	97:  "RIGHTCTRL",        // KEY_RIGHTCTRL
	98:  "KPSLASH",          // KEY_KPSLASH
	99:  "SYSRQ",            // KEY_SYSRQ
	100: "RIGHTALT",         // KEY_RIGHTALT
	101: "LINEFEED",         // KEY_LINEFEED
	102: "HOME",             // KEY_HOME
	103: "UP",               // KEY_UP
	104: "PAGEUP",           // KEY_PAGEUP
	105: "LEFT",             // KEY_LEFT
	106: "RIGHT",            // KEY_RIGHT
	107: "END",              // KEY_END
	108: "DOWN",             // KEY_DOWN
	109: "PAGEDOWN",         // KEY_PAGEDOWN
	110: "INSERT",           // KEY_INSERT
	111: "DELETE",           // KEY_DELETE
	112: "MACRO",            // KEY_MACRO
	113: "MUTE",             // KEY_MUTE
	114: "VOLUMEDOWN",       // KEY_VOLUMEDOWN
	115: "VOLUMEUP",         // KEY_VOLUMEUP
	116: "POWER",            // KEY_POWER
	117: "KPEQUAL",          // KEY_KPEQUAL
	118: "KPPLUSMINUS",      // KEY_KPPLUSMINUS
	119: "PAUSE",            // KEY_PAUSE
	120: "SCALE",            // KEY_SCALE
}
