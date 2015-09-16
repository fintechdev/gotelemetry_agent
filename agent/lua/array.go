package lua

import (
	"github.com/Shopify/go-lua"
)

func pushArray(l *lua.State) {
	l.NewTable()

	l.NewTable()
	l.PushBoolean(true)
	l.SetField(-2, arrayMarkerField)
	l.SetMetaTable(-2)
}
