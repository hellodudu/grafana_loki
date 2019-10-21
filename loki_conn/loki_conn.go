import (
	"log"
	"sync"
)

type LokiConn struct {
	sync.RWMutex
	opts      *Options
	waitGroup util.WaitGroupWrapper
}

func New(opts *Options) (*LokiConn, error) {

	l := &LokiConn{
		opts: opts,
	}

	return l, nil
}

// Main starts an instance of loki_conn and returns an
// error if there was a problem starting up.
func (l *LokiConn) Main() error {
	ctx := &Context{l}

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
		exitFunc(LokiConn.Post())
	})

	err := <-exitCh
	return err
}

func (l *LokiConn) Exit() {

	l.waitGroup.Wait()
}
