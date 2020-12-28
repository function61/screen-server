package evdev

// represents numpad key's meaning for both on and off cases for num lock
//   https://en.wikipedia.org/wiki/Num_Lock
type numpadKeyPair struct {
	withNumLock    Key
	withoutNumLock Key
}

// TODO: just inline this as code?
var keysDependingOnNumLockState = map[Key]*numpadKeyPair{
	KeyKP0:   &numpadKeyPair{Key0, KeyINSERT},
	KeyKP1:   &numpadKeyPair{Key1, KeyEND},
	KeyKP2:   &numpadKeyPair{Key2, KeyDOWN},
	KeyKP3:   &numpadKeyPair{Key3, KeyPAGEDOWN},
	KeyKP4:   &numpadKeyPair{Key4, KeyLEFT},
	KeyKP5:   &numpadKeyPair{Key5, 0}, // 0 = do nothing
	KeyKP6:   &numpadKeyPair{Key6, KeyRIGHT},
	KeyKP7:   &numpadKeyPair{Key7, KeyHOME},
	KeyKP8:   &numpadKeyPair{Key8, KeyUP},
	KeyKP9:   &numpadKeyPair{Key9, KeyPAGEUP},
	KeyKPDOT: &numpadKeyPair{KeyDOT, KeyDELETE},
}

// you can pass any key code. if code is one of those affected by num lock, it returns the "normalized"
// key if num lock is on, and the alternate key if the num lock is off.
//
// TranslateNumLock(KeyKP1,  true) => Key1
// TranslateNumLock(KeyKP1, false) => KeyEND
func TranslateNumLock(code Key, numLock bool) (Key, bool) {
	pair, dependsOnNumLock := keysDependingOnNumLockState[code]
	if dependsOnNumLock {
		if numLock {
			return pair.withNumLock, true
		} else if pair.withoutNumLock != 0 {
			return pair.withoutNumLock, true
		} else {
			return 0, false
		}
	} else {
		return code, true
	}
}
