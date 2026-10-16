package s3

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awss3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
)

func TestGetBlobByKeyReportsStorageFailures(t *testing.T) {
	for _, tc := range []struct {
		name        string
		unavailable bool
		code        string
	}{
		{name: "invalid credentials", code: "InvalidAccessKeyId"},
		{name: "unavailable endpoint", unavailable: true, code: "RequestError"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/xml")
				w.WriteHeader(http.StatusForbidden)
				_, _ = io.WriteString(w, "<Error><Code>InvalidAccessKeyId</Code><Message>Invalid test identity</Message></Error>")
			}))
			t.Cleanup(server.Close)
			if tc.unavailable {
				server.Close()
			}
			useTestClient(t, server.URL)
			q := &query.GetBlobByKey{Key: "file.txt"}
			err := getBlobByKey(context.Background(), q)
			if err == nil {
				t.Fatal("storage failure must not be reported as success")
			}
			cause, ok := errors.Cause(err).(awserr.Error)
			if !ok || cause.Code() != tc.code {
				t.Fatalf("expected AWS error %s, got %v", tc.code, err)
			}
			if q.Result != nil {
				t.Fatal("failed read returned a blob")
			}
		})
	}
}

func TestGetBlobByKeySuccessfulRead(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = io.WriteString(w, "saved content")
	}))
	t.Cleanup(server.Close)
	useTestClient(t, server.URL)
	q := &query.GetBlobByKey{Key: "file.txt"}
	if err := getBlobByKey(context.Background(), q); err != nil {
		t.Fatal(err)
	}
	if q.Result == nil || string(q.Result.Content) != "saved content" || q.Result.Size != 13 || q.Result.ContentType != "text/plain" {
		t.Fatalf("unexpected blob: %+v", q.Result)
	}
}

func TestWrapPreservesAWSError(t *testing.T) {
	for _, underlying := range []error{nil, io.ErrUnexpectedEOF} {
		original := awserr.New("AccessDenied", "denied", underlying)
		wrapped := wrap(original, "storage operation failed")
		if wrapped == nil || errors.Cause(wrapped) != original {
			t.Fatalf("AWS error was lost: %v", wrapped)
		}
	}
	if wrap(nil, "successful operation") != nil {
		t.Fatal("nil error must remain nil")
	}
}

func useTestClient(t *testing.T, endpoint string) {
	t.Helper()
	previous := DefaultClient
	previousBucket := env.Config.BlobStorage.S3.BucketName
	env.Config.BlobStorage.S3.BucketName = "test-bucket"
	t.Cleanup(func() {
		DefaultClient = previous
		env.Config.BlobStorage.S3.BucketName = previousBucket
	})
	DefaultClient = awss3.New(session.Must(session.NewSession(&aws.Config{
		Endpoint:         aws.String(endpoint),
		Region:           aws.String("us-east-1"),
		Credentials:      credentials.NewStaticCredentials("test", "test", ""),
		S3ForcePathStyle: aws.Bool(true),
		MaxRetries:       aws.Int(0),
		HTTPClient:       &http.Client{Timeout: 2 * time.Second},
	})))
}
