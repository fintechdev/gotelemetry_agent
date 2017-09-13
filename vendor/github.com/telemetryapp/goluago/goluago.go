package goluago

import (
	"github.com/telemetryapp/go-lua"
	"github.com/telemetryapp/goluago/pkg/crypto/hmac"
	"github.com/telemetryapp/goluago/pkg/encoding/base64"
	"github.com/telemetryapp/goluago/pkg/encoding/json"
	"github.com/telemetryapp/goluago/pkg/env"
	"github.com/telemetryapp/goluago/pkg/fmt"
	"github.com/telemetryapp/goluago/pkg/net/url"
	"github.com/telemetryapp/goluago/pkg/regexp"
	"github.com/telemetryapp/goluago/pkg/strings"
	"github.com/telemetryapp/goluago/pkg/time"
	"github.com/telemetryapp/goluago/pkg/uuid"
	"github.com/telemetryapp/goluago/util"
)

func Open(l *lua.State) {
	regexp.Open(l)
	strings.Open(l)
	json.Open(l)
	time.Open(l)
	fmt.Open(l)
	url.Open(l)
	util.Open(l)
	hmac.Open(l)
	base64.Open(l)
	env.Open(l)
	uuid.Open(l)
}
