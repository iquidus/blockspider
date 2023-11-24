package crawler

import (
	"time"

	"github.com/ethereum/go-ethereum/log"
)

type logObject struct {
	blockNo uint64
	blocks  int
	txns    int
}

func (l *logObject) add(o *logObject) {
	l.blockNo = o.blockNo
	l.blocks++
	l.txns += o.txns
}

func (l *logObject) clear() {
	l.blockNo = 0
	l.blocks = 0
	l.txns = 0
}

func startLogger(c chan *logObject, logger log.Logger) {
	// Start logging goroutine
	go func(ch chan *logObject) {
		start := time.Now()
		stats := &logObject{
			0,
			0,
			0,
		}
	logLoop:
		for {
			lo, more := <-ch
			if more {
				stats.add(lo)
				if stats.blocks >= 1000 || time.Now().After(start.Add(time.Minute)) {
					logger.Info("Imported new chain segment",
						"blocks", stats.blocks,
						"head", stats.blockNo,
						"transactions", stats.txns,
						"took", time.Since(start))
					stats.clear()
					start = time.Now()
				}
			} else {
				if stats.blocks > 0 {
					logger.Info("Imported new chain segment",
						"blocks", stats.blocks,
						"head", stats.blockNo,
						"transactions", stats.txns,
						"took", time.Since(start))
				}
				break logLoop
			}
		}
	}(c)
}
