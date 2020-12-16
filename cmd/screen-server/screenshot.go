package main

import (
	"io"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/function61/gokit/sync/syncutil"
)

// thanks https://github.com/BurntSushi/xgbutil/blob/master/_examples/screenshot/main.go
func (s *Screen) Screenshot(output io.Writer) error {
	X, err := s.getXConn()
	if err != nil {
		return err
	}

	// root window automatically includes all child windows. background can end up as transparent,
	// but it doesn't matter
	screenshot, err := xgraphics.NewDrawable(X, xproto.Drawable(X.RootWin()))
	if err != nil {
		return err
	}

	return screenshot.WritePng(output)
}

func (s *Screen) getXConn() (*xgbutil.XUtil, error) {
	defer syncutil.LockAndUnlock(&s.xUtilConnMu)()

	var err error
	if s.xUtilConn == nil {
		s.xUtilConn, err = xgbutil.NewConnDisplay(s.XScreenNumberWithColon())
	}

	return s.xUtilConn, err
}
