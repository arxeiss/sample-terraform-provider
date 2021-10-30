package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"flag"
	"net/http"
	"os"

	logformatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"

	"github.com/arxeiss/sample-terraform-provider/server/database"
	"github.com/arxeiss/sample-terraform-provider/server/rest"
)

const (
	defaultLogLevel = logrus.DebugLevel
	systemLogKey    = ""
)

var (
	log       *logrus.Logger
	tokenFlag *string
	addrFlag  *string
)

func init() {
	log = logrus.New()
	logQuite := flag.Bool("q", false, "turns off all logs except errors and higher, has lower priority than -v")
	logTrace := flag.Bool("v", false, "turns on tracing log level - verbose")
	tokenFlag = flag.String("token", "", "set access token to authenticate, can be used also SDC_TOKEN env variable")
	addrFlag = flag.String("addr", ":8090", "default address and port to listen on")
	flag.Parse()

	log.Formatter = &logformatter.Formatter{
		TimestampFormat: "15:04:05.000",
		FieldsOrder:     []string{systemLogKey},
	}
	log.Level = defaultLogLevel
	if logQuite != nil && *logQuite {
		log.Level = logrus.ErrorLevel
	}
	if logTrace != nil && *logTrace {
		log.Level = logrus.TraceLevel
	}
	log.Out = os.Stdout
}

func getAccessToken() string {
	if tokenFlag != nil {
		if len(*tokenFlag) > 10 {
			return *tokenFlag
		} else if len(*tokenFlag) != 0 {
			log.Warn("token provided in flag is shorter than 10 characters, ignoring")
		}
	}
	if token, has := os.LookupEnv("SDC_TOKEN"); has {
		if len(token) > 10 {
			return token
		}
		log.Warn("token provided in ENV variable is shorter than 10 characters, ignoring")
	}

	nonce := make([]byte, 20)
	if _, err := rand.Read(nonce); err != nil {
		log.Fatal("cannot generate random access token: ", err)
	}
	token := base64.RawURLEncoding.EncodeToString(nonce)
	log.Info("Access token for this session is: ", token)
	return token
}

func main() {
	log.Info("Starting server")
	db, err := database.Open("superdupercloud.db", "superdupercloud.sql", log.WithField(systemLogKey, "db"))
	if err != nil {
		log.Fatal("fail to connect to DB: ", err)
	}
	defer func() {
		_ = db.Close()
	}()

	accessToken := getAccessToken()

	addr := ":8090"
	if addrFlag != nil {
		addr = *addrFlag
	}
	s := rest.NewHTTPServer(log.WithField(systemLogKey, "http"), db, accessToken)
	// always returns error. ErrServerClosed on graceful close
	if err := s.Run(addr); !errors.Is(err, http.ErrServerClosed) {
		log.Error("Serve failed: ", err.Error())
	}
}
