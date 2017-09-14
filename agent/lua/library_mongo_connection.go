package lua

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/telemetryapp/go-lua"
	"github.com/telemetryapp/goluago/util"
)

var mongoConnectionFunctions = map[string]func(s *mgo.Session) lua.Function{
	"db": func(s *mgo.Session) lua.Function {
		return func(l *lua.State) int {
			pushMongoDatabase(l, s, lua.CheckString(l, 1))

			return 1
		}
	},

	"live_servers": func(s *mgo.Session) lua.Function {
		return func(l *lua.State) int {
			pushArray(l)

			for index, server := range s.LiveServers() {
				util.DeepPush(l, server)
				l.RawSetInt(-2, index+1)
			}

			return 1
		}
	},

	"database_names": func(s *mgo.Session) lua.Function {
		return func(l *lua.State) int {
			pushArray(l)

			dbNames, err := s.DatabaseNames()

			if err != nil {
				lua.Errorf(l, "%s", err.Error())
			}

			for index, dbName := range dbNames {
				util.DeepPush(l, dbName)
				l.RawSetInt(-2, index+1)
			}

			return 1
		}
	},

	"close": func(s *mgo.Session) lua.Function {
		return func(l *lua.State) int {
			s.Close()

			return 0
		}
	},
}

func pushGoConnection(l *lua.State, connectionString string) {
	ci, err := parseMongoURI(connectionString)
	if err != nil {
		lua.Errorf(l, "%s", err.Error())
	}

	//s, err := mgo.Dial(connectionString)
	s, err := mgo.DialWithInfo(ci)

	if err != nil {
		lua.Errorf(l, "%s", err.Error())
	}

	l.NewTable()

	for name, fn := range mongoConnectionFunctions {
		l.PushGoFunction(fn(s))
		l.SetField(-2, name)
	}

	l.CreateTable(0, 1)
	l.PushGoFunction(func(l *lua.State) int {
		s.Close()

		return 0
	})
	l.SetField(-2, "__gc")
	l.SetMetaTable(-2)
}

func parseMongoURI(rawURI string) (*mgo.DialInfo, error) {
	uri, err := url.Parse(rawURI)
	if err != nil {
		return nil, err
	}

	info := mgo.DialInfo{
		Addrs:    strings.Split(uri.Host, ","),
		Database: strings.TrimPrefix(uri.Path, "/"),
		Timeout:  10 * time.Second,
		FailFast: true,
	}

	if uri.User != nil {
		info.Username = uri.User.Username()
		info.Password, _ = uri.User.Password()
	}

	uriSsl := false
	var sslKeyPath, sslCertPath, sslCAPath string
	uriSslSkipVerify := false

	query := uri.Query()
	for key, values := range query {
		var value string
		if len(values) > 0 {
			value = values[0]
		}

		switch key {
		case "authSource":
			info.Source = value
		case "authMechanism":
			info.Mechanism = value
		case "gssapiServiceName":
			info.Service = value
		case "replicaSet":
			info.ReplicaSetName = value
		case "maxPoolSize":
			poolLimit, err := strconv.Atoi(value)
			if err != nil {
				return nil, errors.New("bad value for maxPoolSize: " + value)
			}
			info.PoolLimit = poolLimit
		case "ssl":
			ssl, err := strconv.ParseBool(value)
			if err != nil {
				return nil, errors.New("bad value for ssl: " + value)
			}
			if ssl {
				uriSsl = true
			}
		case "sslSkipVerify":
			sslSkipVerify, err := strconv.ParseBool(value)
			if err != nil {
				return nil, errors.New("bad value for sslSkipVerify: " + value)
			}
			if sslSkipVerify {
				uriSslSkipVerify = true
			}
		case "sslKey":
			sslKeyPath = value
		case "sslCert":
			sslCertPath = value
		case "sslCA":
			sslCAPath = value
		case "connect":
			if value == "direct" {
				info.Direct = true
				break
			}
			if value == "replicaSet" {
				break
			}
			fallthrough
		default:
			return nil, errors.New("unsupported connection URL option: " + key + "=" + value)
		}
	}

	// deal with TLS
	if uriSsl {
		tlsConfig := tls.Config{}

		if uriSslSkipVerify {
			tlsConfig.InsecureSkipVerify = true
		}

		if sslCAPath != "" {
			caBytes, err := ioutil.ReadFile(sslCAPath)
			if err != nil {
				return nil, errors.New("could not read CA data from '" + sslCAPath + "'")
			}
			caPool := x509.NewCertPool()
			ok := caPool.AppendCertsFromPEM(caBytes)
			if !ok {
				// TODO: log something?
			}
			tlsConfig.RootCAs = caPool
		}

		if sslKeyPath != "" && sslCertPath != "" {
			cert, err := tls.LoadX509KeyPair(sslCertPath, sslKeyPath)
			if err != nil {
				return nil, errors.New("could not load cert and/or key from '" + sslCertPath + " / '" + sslKeyPath + "'")
			}

			tlsConfig.Certificates = []tls.Certificate{cert}
		}

		info.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			dialer := &net.Dialer{
				Timeout:   time.Duration(10 * time.Second),
				Deadline:  time.Now().Add(30 * time.Second),
				KeepAlive: time.Duration(60 * time.Second),
			}
			return tls.DialWithDialer(dialer, "tcp", addr.String(), &tlsConfig)
			//return tls.Dial("tcp", addr.String(), &tlsConfig)
		}
	}

	return &info, nil
}
