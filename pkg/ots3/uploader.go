package ots3

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"io"
	"net/http"
)

type Manager struct {
	bucket       string
	sess         *session.Session
	tracer       opentracing.Tracer
	doer         contract.HttpDoer
	locationFunc func(location string) (url string)
}

type Config struct {
	accessKey    string
	accessSecret string
	region       string
	bucket       string
	sess         *session.Session
	tracer       opentracing.Tracer
	doer         contract.HttpDoer
	locationFunc func(location string) (url string)
}

type Option func(*Config)

func WithTracer(tracer opentracing.Tracer) Option {
	return func(c *Config) {
		c.tracer = tracer
	}
}

func WithHttpClient(client contract.HttpDoer) Option {
	return func(c *Config) {
		c.doer = client
	}
}

func WithLocationFunc(f func(location string) (url string)) Option {
	return func(c *Config) {
		c.locationFunc = f
	}
}

func NewManager(accessKey, accessSecret, endpoint, region, bucket string, opts ...Option) *Manager {
	c := &Config{
		doer: http.DefaultClient,
		locationFunc: func(location string) (url string) {
			return fmt.Sprintf("%s/%s/%s", endpoint, bucket, location)
		},
	}
	for _, f := range opts {
		f(c)
	}

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, accessSecret, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	sess := session.Must(session.NewSession(s3Config))

	m := &Manager{bucket, sess, c.tracer, c.doer, c.locationFunc}

	// add opentracing capabilities if opt in
	if c.tracer != nil {
		sess.Handlers.Build.PushFront(m.otHandler())
	}
	return m
}

func (m *Manager) Upload(ctx context.Context, reader io.Reader) (url string, err error) {

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(m.sess)

	result, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(m.bucket),
		Key:    aws.String(xid.New().String()),
		Body:   reader,
	})

	if err != nil {
		return "", errors.Wrap(err, "unable to upload from io reader")
	}

	return m.locationFunc(result.Location), nil
}

func (m *Manager) UploadFromUrl(ctx context.Context, url string) (newUrl string, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.Wrap(err, "cannot build request")
	}
	resp, err := m.doer.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "cannot fetch image")
	}
	body := resp.Body
	defer body.Close()
	return m.Upload(ctx, body)
}

func (m *Manager) otHandler() func(*request.Request) {
	tracer := m.tracer

	return func(r *request.Request) {
		var sp opentracing.Span

		ctx := r.Context()
		if ctx == nil || !opentracing.IsGlobalTracerRegistered() {
			sp = tracer.StartSpan(r.Operation.Name)
		} else {
			sp, ctx = opentracing.StartSpanFromContext(ctx, r.Operation.Name)
			r.SetContext(ctx)
		}
		ext.SpanKindRPCClient.Set(sp)
		ext.Component.Set(sp, "go-aws")
		ext.HTTPMethod.Set(sp, r.Operation.HTTPMethod)
		ext.HTTPUrl.Set(sp, r.HTTPRequest.URL.String())
		ext.PeerService.Set(sp, r.ClientInfo.ServiceName)

		_ = inject(tracer, sp, r.HTTPRequest.Header)

		r.Handlers.Complete.PushBack(func(req *request.Request) {
			if req.HTTPResponse != nil {
				ext.HTTPStatusCode.Set(sp, uint16(req.HTTPResponse.StatusCode))
			} else {
				ext.Error.Set(sp, true)
			}
			sp.Finish()
		})

		r.Handlers.Retry.PushBack(func(req *request.Request) {
			sp.LogFields(log.String("event", "retry"))
		})
	}
}

func inject(tracer opentracing.Tracer, span opentracing.Span, header http.Header) error {
	return tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(header))
}
