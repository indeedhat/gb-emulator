package types

import "github.com/indeedhat/gb-emulator/internal/emu/enum"

type KeyEvent struct {
	Key  enum.KeyCode
	Down bool
}
