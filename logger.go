package seq

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

var LogLevel = map[string]int{
	"TRACE":       10,
	"Trace":       10,
	"DEBUG":       20,
	"Debug":       20,
	"INFO":        30,
	"Information": 30,
	"WARN":        40,
	"Warning":     40,
	"ERROR":       50,
	"Error":       50,
	"FATAL":       60,
	"Fatal":       60,
}

var LogLevelNumber = map[int]string{
	10: "TRACE",
	20: "DEBUG",
	30: "INFO",
	40: "WARN",
	50: "ERROR",
	60: "FATAL",
}

func getLogLevelName(l int) string {
	ln := l / 10 * 10
	if ln >= 60 {
		return "FATAL"
	}
	if ln <= 10 {
		return "TRACE"
	}
	return LogLevelNumber[l]
}

type message struct {
	messageTemplate string
	level           int
	params          []interface{}
}

type Logger struct {
	ch          chan message
	apiKey      string
	apiEndpoint string
	stdErrOut   bool
	closeCh     chan int
}

func NewLogger(apiEndpoint string, apiKey string, bufferSize int) *Logger {
	l := Logger{
		apiEndpoint: apiEndpoint,
		apiKey:      apiKey,
		ch:          make(chan message, bufferSize),
	}
	l.connect()
	return &l
}

type closerBuffer struct {
	b      *bytes.Buffer
	closed bool
}

func (cb closerBuffer) Read(p []byte) (n int, err error) {
	if cb.closed {
		return 0, errors.New("EOF")
	}
	return cb.b.Read(p)
}

func (cb *closerBuffer) Close() error {
	cb.closed = true
	cb.b = nil
	return nil
}

func (l *Logger) connect() {
	buf := closerBuffer{b: bytes.NewBuffer(make([]byte, 2048))}
	go func() {
		c := &http.Client{
			Timeout: time.Second * 2,
		}
		for {
			_, err := c.Post(path.Join(l.apiEndpoint, "/api/events/raw"), "application/vnd.serilog.clef", buf)
			if err != nil {
				time.Sleep(time.Second)
				return
			}
		}
	}()
	for {
		select {
		case msg := <-l.ch:
			params, err := ExtractParams(msg.messageTemplate, msg.params)
			if err != nil {
				panic("todo")
			}
			if len(msg.params) == 0 {
				params["@m"] = msg.messageTemplate
			} else {
				params["@m"] = RenderMsgTemplate(msg.messageTemplate, params)
			}
			params["@l"] = strconv.Itoa(msg.level)
			params["@t"] = time.Now().Format(time.RFC3339)
			if l.stdErrOut {
				_, _ = os.Stdout.WriteString(fmt.Sprintf("%v\t%v\t%v\n", time.Now().Format(time.RFC3339),
					getLogLevelName(msg.level), RenderMsgTemplate(msg.messageTemplate, params)))
			}
		case <-l.closeCh:
			_ = buf.Close()
			return
		}
	}
}

func (l *Logger) EnableConsole() {
	l.stdErrOut = true
}
func (l *Logger) DisableConsole() {
	l.stdErrOut = false
}
func (l *Logger) Log(level string, messageTemplate string, params ...interface{}) {
	lv := LogLevel[level]
	if lv == 0 {
		lv = 30
	}
	msg := message{
		messageTemplate: messageTemplate,
		level:           lv,
		params:          params,
	}
	select {
	case l.ch <- msg:
	default:
		return
	}
}
func (l *Logger) Trace(messageTemplate string, params ...interface{}) {
	l.Log("TRACE", messageTemplate, params...)
}

func (l *Logger) Debug(messageTemplate string, params ...interface{}) {
	l.Log("DEBUG", messageTemplate, params...)
}

func (l *Logger) Info(messageTemplate string, params ...interface{}) {
	l.Log("INFO", messageTemplate, params...)
}

func (l *Logger) Warn(messageTemplate string, params ...interface{}) {
	l.Log("WARN", messageTemplate, params...)
}

func (l *Logger) Error(messageTemplate string, params ...interface{}) {
	l.Log("ERROR", messageTemplate, params...)
}

func (l *Logger) Fatal(messageTemplate string, params ...interface{}) {
	l.Log("FATAL", messageTemplate, params...)
}
