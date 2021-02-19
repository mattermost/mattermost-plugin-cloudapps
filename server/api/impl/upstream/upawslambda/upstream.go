// Copyright (c) 2020-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package upawslambda

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/service/lambda"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/api/impl/aws"
)

type Upstream struct {
	app *apps.App
	aws *aws.Client
}

func NewUpstream(app *apps.App, aws *aws.Client) *Upstream {
	return &Upstream{
		app: app,
		aws: aws,
	}
}

func (u *Upstream) OneWay(call *apps.Call) error {
	_, err := u.aws.InvokeLambda(u.app.Manifest.AppID, u.app.Manifest.Version, call.URL, lambda.InvocationTypeEvent, call)
	return err
}

func (u *Upstream) Roundtrip(call *apps.Call) (io.ReadCloser, error) {
	log.Printf("call: %#+v\n", *call)
	req := struct {
		Path       string            `json:"path"`
		HTTPMethod string            `json:"httpMethod"`
		Headers    map[string]string `json:"headers"`
		Body       interface{}       `json:"body"`
	}{
		Path:       call.URL,
		HTTPMethod: "POST",
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       call,
	}
	log.Printf("req: %#+v\n", req)

	bb, err := u.aws.InvokeLambda(u.app.Manifest.AppID, u.app.Manifest.Version, u.app.Manifest.Functions[0].Name, lambda.InvocationTypeRequestResponse, req)
	if err != nil {
		return nil, err
	}
	return ioutil.NopCloser(bytes.NewReader(bb)), err
}
