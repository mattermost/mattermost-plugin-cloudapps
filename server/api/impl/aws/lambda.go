// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package aws

import (
	"encoding/json"
	glog "log"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/lambda"

	"github.com/mattermost/mattermost-plugin-apps/apps"
)

// InvokeLambda runs a lambda function with specified name and returns a payload
func (c *Client) InvokeLambda(appID apps.AppID, appVersion apps.AppVersion, functionName, invocationType string, request interface{}) ([]byte, error) {
	payload, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "Error marshaling request payload")
	}

	name, err := getFunctionName(appID, appVersion, functionName)
	if err != nil {
		return nil, errors.Wrap(err, "can't get function name")
	}

	input := &lambda.InvokeInput{
		FunctionName:   aws.String(name),
		InvocationType: aws.String(invocationType),
		LogType:        aws.String("tail"),
		Payload:        payload,
	}
	glog.Printf("input.Validate(): %#+v\n", input.Validate())

	result, err := c.Service().lambda.Invoke(&lambda.InvokeInput{
		FunctionName:   aws.String(name),
		InvocationType: aws.String(invocationType),
		Payload:        payload,
	})
	glog.Printf("payload: %#+v\n", string(payload))
	glog.Printf("err: %#+v\n", err)

	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if ok {
			glog.Printf("awsErr.Message(): %#+v\n", awsErr.Message())
			glog.Printf("awsErr.Error(): %#+v\n", awsErr.Error())
		}

		return nil, errors.Wrapf(err, "Error calling function %s", name)
	}
	glog.Printf("result: %#+v\n", result)
	glog.Printf("string(result.Payload): %#+v\n", string(result.Payload))

	resp := struct {
		Body string `json:"body"`
	}{}

	err = json.Unmarshal(result.Payload, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "Error marshaling request payload")
	}
	glog.Printf("resp: %#+v\n", resp)

	return []byte(resp.Body), nil
}
