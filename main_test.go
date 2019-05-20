package goseq

import (
	"testing"
	"time"
)

import (
	"math/rand"
)

func TestConnect(t *testing.T) {
	Connect("http://127.0.0.1:9999/api/seq/seqtest/stream/test2", "", false)
	for {
		Trace("timestamp: {timestamp}, random: {random}", time.Now().UnixNano(), rand.Int())
		Debug("timestamp: {timestamp}, random: {random}", time.Now().UnixNano(), rand.Int())
		Info("requestId: {requestID}, took {reqTime}, timestamp: {timestamp}, random: {random}", time.Now().UnixNano(), rand.Int())
		Warn("timestamp: {timestamp}, random: {random}", time.Now().UnixNano(), rand.Int())
		Error("timestamp: {timestamp}, random: {random}", time.Now().UnixNano(), rand.Int())
		Fatal("timestamp: {timestamp}, random: {random}", time.Now().UnixNano(), rand.Int())
		time.Sleep(100 * time.Millisecond)
	}
}
func TestConnect2(t *testing.T) {
	Connect("http://127.0.0.1:9999/api/seq/seqtest/stream/stream2", "", true)
	for {
		Trace("timestamp: {timestamp}, random: {random}", time.Now().UnixNano(), rand.Int())
		Debug("timestamp: {timestamp}, random: {random}", time.Now().UnixNano(), rand.Int())
		Info("timestamp: {timestamp}, random: {random}", time.Now().UnixNano(), rand.Int())
		Warn("timestamp: {timestamp}, random: {random}", time.Now().UnixNano(), rand.Int())
		Error("timestamp: {timestamp}, random: {random}", time.Now().UnixNano(), rand.Int())
		Fatal("timestamp: {timestamp}, random: {random}", time.Now().UnixNano(), rand.Int())
		time.Sleep(100 * time.Millisecond)
	}
}
func TestConnect3(t *testing.T) {
	Connect("http://127.0.0.1:5341", "", true)
	for {
		Trace("timestamp: {timestamp}, random: {random}", time.Now().UnixNano(), rand.Int())
		Debug("timestamp: {timestamp}, random: {random}", time.Now().UnixNano(), rand.Int())
		Info("timestamp: {timestamp}, random: {random}", time.Now().UnixNano(), rand.Int())
		Warn("timestamp: {timestamp}, random: {random}", time.Now().UnixNano(), rand.Int())
		Error("timestamp: {timestamp}, random: {random}", time.Now().UnixNano(), rand.Int())
		Fatal("timestamp: {timestamp}, random: {random}", time.Now().UnixNano(), rand.Int())
		time.Sleep(100 * time.Millisecond)
	}
}
