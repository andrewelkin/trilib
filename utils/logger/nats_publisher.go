package logger

import "github.com/nats-io/nats.go"

type NatsLogger struct {
	subj     string
	natsConn *nats.Conn
}

func NewNatsLogger(subject string, natsConn *nats.Conn) *NatsLogger {
	nl := &NatsLogger{
		subj:     subject,
		natsConn: natsConn,
	}
	return nl
}

func (nl *NatsLogger) Write(p []byte) (n int, err error) {
	nl.natsConn.Publish(nl.subj, p)
	return len(p), nil
}
