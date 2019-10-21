package loki_conn

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/hellodudu/grafana_loki/internal/util"
)

type LokiConn struct {
	sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
	opts      *Options
	waitGroup util.WaitGroupWrapper
}

// loki push request struct
type Entry struct {
	TS   string `json:"ts"`
	Line string `json:"line"`
}

type Stream struct {
	Labels  string   `json:"labels"`
	Entries []*Entry `json:"entries"`
}

type PushRequest struct {
	Streams []*Stream `json:"streams"`
}

func New(opts *Options) (*LokiConn, error) {
	l := &LokiConn{
		opts: opts,
	}

	l.ctx, l.cancel = context.WithCancel(context.Background())

	return l, nil
}

// Main starts an instance of loki_conn and returns an
// error if there was a problem starting up.
func (l *LokiConn) Main() error {

	exitCh := make(chan error)
	var once sync.Once
	exitFunc := func(err error) {
		once.Do(func() {
			if err != nil {
				log.Fatal("LokiConn Main() error:", err)
			}
			exitCh <- err
		})
	}

	l.waitGroup.Wrap(func() {
		exitFunc(l.HTTPRequest())
	})

	err := <-exitCh
	return err
}

func (l *LokiConn) Exit() {
	l.cancel()
	//l.waitGroup.Wait()
}

func (l *LokiConn) HTTPRequest() error {

	for {
		select {
		case <-l.ctx.Done():
			return nil
		default:
		}

		t := time.Now()

		req := &PushRequest{
			Streams: make([]*Stream, 0),
		}

		entry := &Entry{
			TS:   time.Now().Format(time.RFC3339),
			Line: "[info] heartbeat",
		}

		entries := make([]*Entry, 0)
		entries = append(entries, entry)

		labels := "{loki_conn=\"connection\"}"
		req.Streams = append(req.Streams, &Stream{Labels: labels, Entries: entries})
		reqJSON, err := json.Marshal(req)
		if err != nil {
			log.Println("marshal json error:", err)
			d := time.Since(t)
			time.Sleep(l.opts.Interval - d)
			continue
		}

		request, err := http.NewRequest("POST", l.opts.URL, bytes.NewBuffer(reqJSON))
		request.Header.Set("X-Custom-Header", "myvalue")
		request.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			log.Println("http request with error:", err)
			d := time.Since(t)
			time.Sleep(l.opts.Interval - d)
			continue
		}

		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Status:", resp.Status, ", Body:", string(body))

		d := time.Since(t)
		time.Sleep(l.opts.Interval - d)
	}
}
