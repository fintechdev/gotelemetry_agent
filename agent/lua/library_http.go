package lua

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/telemetryapp/go-lua"
	"github.com/telemetryapp/goluago/util"
)

var httpLibrary = []lua.RegistryFunction{
	{
		"get",
		func(l *lua.State) int {
			url := lua.CheckString(l, 1)

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				lua.Errorf(l, "%s", err.Error())
			}

			argIndex := 2
			username := ""
			password := ""
			if l.IsString(argIndex) {

				username = lua.CheckString(l, argIndex)
				argIndex++

				if l.IsString(argIndex) {
					password = lua.CheckString(l, argIndex)
					argIndex++
				}

			}

			if l.IsTable(argIndex) {
				header, _ := util.PullStringTable(l, argIndex)
				for key, value := range header {
					req.Header.Set(key, value)
				}
			}

			// see if there is another table in the arguments, and extract the TLS
			// information from there
			useTLS := false
			var tlsConfig *tls.Config

			argIndex++
			if l.IsTable(argIndex) {
				useTLS = true
				var tlsSettings map[string]string
				tlsSettings, err = util.PullStringTable(l, argIndex)
				if err != nil {
					lua.Errorf(l, "Error reading TLS Settings table: %s", err.Error())
				}
				tlsConfig = tlsConfigFromLuaTable(l, tlsSettings)
			}

			var client *http.Client
			if useTLS {
				client = &http.Client{
					Transport: &http.Transport{
						TLSClientConfig: tlsConfig,
					},
				}
			} else {
				client = http.DefaultClient
			}

			if len(username) > 0 || len(password) > 0 {
				req.SetBasicAuth(username, password)
			}

			resp, err := client.Do(req)
			if err != nil {
				lua.Errorf(l, "%s", err.Error())
			}

			defer resp.Body.Close()

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				lua.Errorf(l, "%s", err.Error())
			}

			util.DeepPush(l, string(data))

			return 1
		},
	},

	{
		"post",
		func(l *lua.State) int {
			url := lua.CheckString(l, 1)
			body := lua.CheckString(l, 2)

			req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
			if err != nil {
				lua.Errorf(l, "%s", err.Error())
			}

			argIndex := 3
			username := ""
			password := ""
			if l.IsString(argIndex) {

				username = lua.CheckString(l, argIndex)
				argIndex++

				if l.IsString(argIndex) {
					password = lua.CheckString(l, argIndex)
					argIndex++
				}

			}

			if l.IsTable(argIndex) {
				header, _ := util.PullStringTable(l, argIndex)
				for key, value := range header {
					req.Header.Set(key, value)
				}
			}

			// see if there is another table in the arguments, and extract the TLS
			// information from there
			useTLS := false
			var tlsConfig *tls.Config

			argIndex++
			if l.IsTable(argIndex) {
				useTLS = true
				var tlsSettings map[string]string
				tlsSettings, err = util.PullStringTable(l, argIndex)
				if err != nil {
					lua.Errorf(l, "Error reading TLS Settings table: %s", err.Error())
				}
				tlsConfig = tlsConfigFromLuaTable(l, tlsSettings)
			}

			var client *http.Client
			if useTLS {
				client = &http.Client{
					Transport: &http.Transport{
						TLSClientConfig: tlsConfig,
					},
				}
			} else {
				client = http.DefaultClient
			}

			if len(username) > 0 || len(password) > 0 {
				req.SetBasicAuth(username, password)
			}

			resp, err := client.Do(req)
			if err != nil {
				lua.Errorf(l, "%s", err.Error())
			}

			defer resp.Body.Close()

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				lua.Errorf(l, "%s", err.Error())
			}

			util.DeepPush(l, string(data))

			return 1
		},
	},

	{
		"custom",
		func(l *lua.State) int {
			method := lua.CheckString(l, 1)
			url := lua.CheckString(l, 2)
			body := lua.OptString(l, 3, "")

			if len(method) == 0 {
				method = "POST"
			}

			var req *http.Request
			var err error
			if len(body) > 0 {
				req, err = http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
			} else {
				req, err = http.NewRequest(method, url, nil)
			}
			if err != nil {
				lua.Errorf(l, "%s", err.Error())
			}

			argIndex := 4
			username := ""
			password := ""
			if l.IsString(argIndex) {

				username = lua.CheckString(l, argIndex)
				argIndex++

				if l.IsString(argIndex) {
					password = lua.CheckString(l, argIndex)
					argIndex++
				}

			}

			if l.IsTable(argIndex) {
				header, _ := util.PullStringTable(l, argIndex)
				for key, value := range header {
					req.Header.Set(key, value)
				}
			}

			// see if there is another table in the arguments, and extract the TLS
			// information from there
			useTLS := false
			var tlsConfig *tls.Config

			argIndex++
			if l.IsTable(argIndex) {
				useTLS = true
				var tlsSettings map[string]string
				tlsSettings, err = util.PullStringTable(l, argIndex)
				if err != nil {
					lua.Errorf(l, "Error reading TLS Settings table: %s", err.Error())
				}
				tlsConfig = tlsConfigFromLuaTable(l, tlsSettings)
			}

			var client *http.Client
			if useTLS {
				client = &http.Client{
					Transport: &http.Transport{
						TLSClientConfig: tlsConfig,
					},
				}
			} else {
				client = http.DefaultClient
			}

			if len(username) > 0 || len(password) > 0 {
				req.SetBasicAuth(username, password)
			}

			resp, err := client.Do(req)
			if err != nil {
				lua.Errorf(l, "%s", err.Error())
			}

			defer resp.Body.Close()

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				lua.Errorf(l, "%s", err.Error())
			}

			util.DeepPush(l, string(data))

			return 1
		},
	},
}

func tlsConfigFromLuaTable(l *lua.State, tlsSettings map[string]string) *tls.Config {
	var tlsKeyPath, tlsCertPath string
	var tlsConfig tls.Config
	for key, value := range tlsSettings {
		switch key {
		case "InsecureSkipVerify":
			v, err := strconv.ParseBool(value)
			if err != nil {
				lua.Errorf(l, "Error in TLS Settings for '%s': %s", key, err.Error())
			}
			tlsConfig.InsecureSkipVerify = v
		case "SessionTicketsDisabled":
			v, err := strconv.ParseBool(value)
			if err != nil {
				lua.Errorf(l, "Error in TLS Settings for '%s': %s", key, err.Error())
			}
			tlsConfig.SessionTicketsDisabled = v
		case "PreferServerCipherSuites":
			v, err := strconv.ParseBool(value)
			if err != nil {
				lua.Errorf(l, "Error in TLS Settings for '%s': %s", key, err.Error())
			}
			tlsConfig.PreferServerCipherSuites = v
		case "MinVersion":
			switch value {
			case "VersionSSL30":
				tlsConfig.MinVersion = tls.VersionSSL30
			case "VersionTLS10":
				tlsConfig.MinVersion = tls.VersionTLS10
			case "VersionTLS11":
				tlsConfig.MinVersion = tls.VersionTLS11
			case "VersionTLS12":
				tlsConfig.MinVersion = tls.VersionTLS12
			default:
				lua.Errorf(l, "Error in TLS Settings for '%s': value '%s' is not an accepted value", key, value)
			}
		case "MaxVersion":
			switch value {
			case "VersionSSL30":
				tlsConfig.MaxVersion = tls.VersionSSL30
			case "VersionTLS10":
				tlsConfig.MaxVersion = tls.VersionTLS10
			case "VersionTLS11":
				tlsConfig.MaxVersion = tls.VersionTLS11
			case "VersionTLS12":
				tlsConfig.MaxVersion = tls.VersionTLS12
			default:
				lua.Errorf(l, "Error in TLS Settings for '%s': value '%s' is not an accepted value", key, value)
			}
		case "CipherSuites":
			var clientSupportedCiphers []uint16
			for _, cipher := range strings.Split(value, ",") {
				switch cipher {
				case "TLS_RSA_WITH_RC4_128_SHA":
					clientSupportedCiphers = append(clientSupportedCiphers, tls.TLS_RSA_WITH_RC4_128_SHA)
				case "TLS_RSA_WITH_3DES_EDE_CBC_SHA":
					clientSupportedCiphers = append(clientSupportedCiphers, tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA)
				case "TLS_RSA_WITH_AES_128_CBC_SHA":
					clientSupportedCiphers = append(clientSupportedCiphers, tls.TLS_RSA_WITH_AES_128_CBC_SHA)
				case "TLS_RSA_WITH_AES_256_CBC_SHA":
					clientSupportedCiphers = append(clientSupportedCiphers, tls.TLS_RSA_WITH_AES_256_CBC_SHA)
				case "TLS_RSA_WITH_AES_128_GCM_SHA256":
					clientSupportedCiphers = append(clientSupportedCiphers, tls.TLS_RSA_WITH_AES_128_GCM_SHA256)
				case "TLS_RSA_WITH_AES_256_GCM_SHA384":
					clientSupportedCiphers = append(clientSupportedCiphers, tls.TLS_RSA_WITH_AES_256_GCM_SHA384)
				case "TLS_ECDHE_ECDSA_WITH_RC4_128_SHA":
					clientSupportedCiphers = append(clientSupportedCiphers, tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA)
				case "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA":
					clientSupportedCiphers = append(clientSupportedCiphers, tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA)
				case "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA":
					clientSupportedCiphers = append(clientSupportedCiphers, tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA)
				case "TLS_ECDHE_RSA_WITH_RC4_128_SHA":
					clientSupportedCiphers = append(clientSupportedCiphers, tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA)
				case "TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA":
					clientSupportedCiphers = append(clientSupportedCiphers, tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA)
				case "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA":
					clientSupportedCiphers = append(clientSupportedCiphers, tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA)
				case "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA":
					clientSupportedCiphers = append(clientSupportedCiphers, tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA)
				case "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256":
					clientSupportedCiphers = append(clientSupportedCiphers, tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256)
				case "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256":
					clientSupportedCiphers = append(clientSupportedCiphers, tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256)
				case "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384":
					clientSupportedCiphers = append(clientSupportedCiphers, tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384)
				case "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384":
					clientSupportedCiphers = append(clientSupportedCiphers, tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384)
				default:
					lua.Errorf(l, "Error in TLS Settings for '%s': '%s' is not an accepted cipher", key, cipher)
				}
			}
		case "RootCAs":
			caBytes, err := ioutil.ReadFile(value)
			if err != nil {
				lua.Errorf(l, "Error in TLS Settings for '%s': could not read CA data from '%s': %s", key, value, err.Error())
			}
			caPool := x509.NewCertPool()
			ok := caPool.AppendCertsFromPEM(caBytes)
			if !ok {
				lua.Errorf(l, "Error in TLS Settings for '%s': could not append certs to root CA pool: %s", key, err.Error())
			}
			tlsConfig.RootCAs = caPool
		case "certFile":
			tlsCertPath = value
		case "keyFile":
			tlsKeyPath = value
		default:
			// TODO: print something
		}
	}

	// we can set these only after we parsed both paths from the table
	if tlsKeyPath != "" && tlsCertPath != "" {
		cert, err := tls.LoadX509KeyPair(tlsCertPath, tlsKeyPath)
		if err != nil {
			lua.Errorf(l, "Error in TLS Settings: could not load cert and/or key from '%s' / '%s': %s", tlsCertPath, tlsKeyPath, err.Error())
		}

		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	return &tlsConfig
}

func openHTTPLibrary(l *lua.State) {
	open := func(l *lua.State) int {
		lua.NewLibrary(l, httpLibrary)
		return 1
	}

	lua.Require(l, "telemetry/http", open, false)
	l.Pop(1)
}
