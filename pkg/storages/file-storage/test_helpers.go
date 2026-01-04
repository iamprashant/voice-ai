// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package storage_files

import (
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// mockS3Client implements s3iface.S3API for testing
type mockS3Client struct {
	s3iface.S3API
	putObjectError    error
	getObjectError    error
	getObjectResponse *s3.GetObjectOutput
	presignURL        string
	presignError      error
}

func (m *mockS3Client) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return &s3.PutObjectOutput{}, m.putObjectError
}

func (m *mockS3Client) GetObjectWithContext(ctx aws.Context, input *s3.GetObjectInput, opts ...request.Option) (*s3.GetObjectOutput, error) {
	return m.getObjectResponse, m.getObjectError
}

func (m *mockS3Client) GetObjectRequest(input *s3.GetObjectInput) (*request.Request, *s3.GetObjectOutput) {
	req := &request.Request{}
	if m.presignError != nil {
		req.Error = m.presignError
	} else {
		u, _ := url.Parse(m.presignURL)
		req.HTTPRequest = &http.Request{
			URL: u,
		}
	}
	return req, &s3.GetObjectOutput{}
}
