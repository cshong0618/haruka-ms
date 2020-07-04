package audit

import (
	"context"
	"github.com/cshong0618/haruka/post/pkg/utils/errors"
	log "github.com/sirupsen/logrus"
)

func Start(ctx context.Context,
	topic string,
	msg ...interface{}) *logCloser {
	// check if msg length is multiple of 2
	if len(msg)%2 != 0 {
		// we cannot process the fields,
		//	log them in a slice
		log.WithFields(log.Fields{
			"msg": msg,
			"status": "START",
		}).Info(topic)
	} else {
		fields := log.Fields{
			"status": "START",
		}

		for i := 0; i < len(msg); i += 2 {
			key, ok := msg[i].(string)
			if !ok {
				continue
			}

			fields[key] = msg[i+1]
		}

		log.WithFields(fields).Info(topic)
	}

	lc := &logCloser{
		ctx: ctx,
		topic: topic,
	}
	return lc
}

type LogCloser interface {
	Capture(werr *errors.WrappedError) *logCloser
	End()
}

type logCloser struct {
	ctx   context.Context
	topic string
	werr  *errors.WrappedError
}

func (lc *logCloser) Capture(werr *errors.WrappedError) *logCloser {
	lc.werr = werr
	return lc
}

func (lc *logCloser) End() {
	errMsg := lc.werr.Error()

	if errMsg == "" {
		// No errors
		log.WithFields(log.Fields{
			"status": "END",
		}).Info(lc.topic)
	} else {
		// error
		log.WithFields(log.Fields{
			"status": "ERROR",
			"error": errMsg,
		}).Error(lc.topic)
	}
}
