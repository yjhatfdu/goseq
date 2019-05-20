package goseq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

//var LogLevel = map[string]int{
//	"TRACE":       10,
//	"Trace":       10,
//	"DEBUG":       20,
//	"Debug":       20,
//	"INFO":        30,
//	"Information": 30,
//	"WARN":        40,
//	"Warning":     40,
//	"ERROR":       50,
//	"Error":       50,
//	"FATAL":       60,
//	"Fatal":       60,
//}

var LogLevelNumber = map[int]string{
	10: "TRACE",
	20: "DEBUG",
	30: "INFO",
	40: "WARN",
	50: "ERROR",
	60: "FATAL",
}

const (
	TRACE = 10
	DEBUG = 20
	INFO  = 30
	WARN  = 40
	ERROR = 50
	FATAL = 60
)

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
	minLevel    int
	chunked     bool
}

func NewLogger(apiEndpoint string, apiKey string, bufferSize int, stream bool) *Logger {
	l := Logger{
		apiEndpoint: apiEndpoint,
		apiKey:      apiKey,
		ch:          make(chan message, bufferSize),
		chunked:     !stream,
	}
	go l.connect()
	return &l
}

func NewLoggerLocal(bufferSize int) *Logger {
	l := Logger{
		ch:        make(chan message, bufferSize),
		stdErrOut: true,
	}
	go l.connectLocal()
	return &l
}

func (l *Logger)connectLocal()  {
	for {
		select {
		case msg := <-l.ch:
			_, m, err := marshalMsg(msg)
			if err != nil {
				continue
			}

			if l.stdErrOut {
				_, _ = os.Stdout.WriteString(fmt.Sprintf("%v\t%v\t%v\n", time.Now().Format("2006-01-02T15:04:00"),
					getLogLevelName(msg.level), m))
			}
		case <-l.closeCh:
			return
		}
	}
}

func (l *Logger) connect() {
	pr, pw := io.Pipe()
	buf := bytes.NewBuffer(make([]byte, 0))
	locker := sync.Mutex{}
	if l.chunked {
		go func() {
			c := &http.Client{}
			t := time.NewTicker(time.Second)
			for range t.C {
				locker.Lock()
				r, err := http.NewRequest("POST", l.apiEndpoint+"/api/events/raw", buf)
				r.Header.Set("X-Seq-ApiKey", l.apiKey)
				r.Header.Set("content-type", "application/vnd.serilog.clef")
				_, err = c.Do(r)
				buf.Reset()
				if err != nil {
					fmt.Println(err)
				}
				locker.Unlock()
			}
		}()
	} else {
		go func() {
			c := &http.Client{}
			for {
				r, err := http.NewRequest("POST", l.apiEndpoint+"/api/events/raw", pr)
				r.Header.Set("X-Seq-ApiKey", l.apiKey)
				r.Header.Set("content-type", "application/vnd.serilog.clef")
				_, err = c.Do(r)
				if err != nil {
					fmt.Println(err)
				}
				time.Sleep(time.Second)
				pr, pw = io.Pipe()
			}
		}()
	}

	for {
		select {
		case msg := <-l.ch:
			d, m, err := marshalMsg(msg)
			if err != nil {
				continue
			}
			if l.chunked {
				locker.Lock()
				buf.Write(d)
				locker.Unlock()
			} else {
				_, _ = pw.Write(d)
			}
			if l.stdErrOut {
				_, _ = os.Stdout.WriteString(fmt.Sprintf("%v\t%v\t%v\n", time.Now().Format(time.RFC3339),
					getLogLevelName(msg.level), m))
			}
		case <-l.closeCh:
			_ = pw.Close()
			return
		}
	}
}

func marshalMsg(msg message) ([]byte, string, error) {
	params, err := ExtractParams(msg.messageTemplate, msg.params...)
	if err != nil {
		return nil, "", err
	}
	if len(msg.params) == 0 {
		params["@m"] = msg.messageTemplate
	} else {
		params["@m"] = RenderMsgTemplate(msg.messageTemplate, params)
	}
	params["@l"] = LogLevelNumber[msg.level]
	params["@t"] = time.Now().Format(time.RFC3339)
	d, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	d = append(d, '\n')
	return d, params["@m"], nil
}

func (l *Logger) SetOutputLevel(level int) {
	l.minLevel = level
}

func (l *Logger) EnableConsole() {
	l.stdErrOut = true
}
func (l *Logger) DisableConsole() {
	l.stdErrOut = false
}
func (l *Logger) Log(level int, messageTemplate string, params ...interface{}) {
	if level < l.minLevel {
		return
	}
	msg := message{
		messageTemplate: messageTemplate,
		level:           level,
		params:          params,
	}
	select {
	case l.ch <- msg:
	default:
		return
	}
}
func (l *Logger) Trace(messageTemplate string, params ...interface{}) {
	l.Log(TRACE, messageTemplate, params...)
}

func (l *Logger) Debug(messageTemplate string, params ...interface{}) {
	l.Log(DEBUG, messageTemplate, params...)
}

func (l *Logger) Info(messageTemplate string, params ...interface{}) {
	l.Log(INFO, messageTemplate, params...)
}

func (l *Logger) Warn(messageTemplate string, params ...interface{}) {
	l.Log(WARN, messageTemplate, params...)
}

func (l *Logger) Error(messageTemplate string, params ...interface{}) {
	l.Log(ERROR, messageTemplate, params...)
}

func (l *Logger) Fatal(messageTemplate string, params ...interface{}) {
	l.Log(FATAL, messageTemplate, params...)
}
