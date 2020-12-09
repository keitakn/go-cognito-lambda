package main

import (
	"github.com/aws/aws-lambda-go/events"
)

type createUserPoolsCustomMessageEventParams struct {
	TriggerSource string
	UserPoolId    string
	UserName      string
	Sub           string
	CodeParameter string
}

func createUserPoolsCustomMessageEvent(p *createUserPoolsCustomMessageEventParams) *events.CognitoEventUserPoolsCustomMessage {
	cc := &events.CognitoEventUserPoolsCallerContext{
		AWSSDKVersion: "",
		ClientID:      "",
	}

	ch := &events.CognitoEventUserPoolsHeader{
		Version:       "",
		TriggerSource: p.TriggerSource,
		Region:        "",
		UserPoolID:    p.UserPoolId,
		CallerContext: *cc,
		UserName:      p.UserName,
	}

	ua := map[string]interface{}{
		"sub": p.Sub,
	}

	cm := map[string]string{
		"key": "",
	}

	req := &events.CognitoEventUserPoolsCustomMessageRequest{
		UserAttributes:    ua,
		CodeParameter:     p.CodeParameter,
		UsernameParameter: p.UserName,
		ClientMetadata:    cm,
	}

	res := &events.CognitoEventUserPoolsCustomMessageResponse{
		SMSMessage:   "SMSMessage",
		EmailMessage: "EmailMessage",
		EmailSubject: "EmailSubject",
	}

	ev := &events.CognitoEventUserPoolsCustomMessage{
		CognitoEventUserPoolsHeader: *ch,
		Request:                     *req,
		Response:                    *res,
	}

	return ev
}
