package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/sirupsen/logrus"
)

var singletonSessionOpts *SessionOpts

// DefaultSessionConfig returns a singleton Session with default AWS config.
func DefaultSessionConfig() *SessionOpts {
	if singletonSessionOpts != nil {
		singletonSessionOpts = &SessionOpts{
			Opts: &session.Options{
				Config: aws.Config{
					Region: aws.String(DefaultAWSRegion),
					// TODO(gsagula): we should supply a Retryer since it provides a more robust strategy. Ideally, we should delegate this
					// kind of operations to a proxy or something that knows how to do that well.
					// https://github.com/aws/aws-sdk-go/blob/99cd35c8c7d369ba8c32c46ed306f6c88d24cfd7/aws/request/retryer.go#L20
					MaxRetries: aws.Int(DefaultAWSClientRetries),
				},
			},
		}
	}
	return singletonSessionOpts
}

// SessionOpts supplies options for creating a session.
type SessionOpts struct {
	Opts *session.Options
}

// CreateSession creates a new AWS session.
func (s *SessionOpts) CreateSession(logger log.FieldLogger) (*session.Session, error) {
	if s.Opts == nil {
		s.Opts = DefaultSessionConfig().Opts
	}

	sess, err := session.NewSessionWithOptions(*s.Opts)
	if err != nil {
		return nil, err
	}

	setLoggerHandler(sess, logger)

	return sess, nil
}

func setLoggerHandler(sess *session.Session, logger log.FieldLogger) {
	sess.Handlers.Send.PushFront(func(r *request.Request) {
		logger.Debugf("%s: %s%s\n%s", r.HTTPRequest.Method, r.HTTPRequest.URL.Host, r.HTTPRequest.URL.RawPath, r.Params)
	})
}