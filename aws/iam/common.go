package iam

import (
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/iam"
)

const PolicyVersion20121017 PolicyVersion = "2012-10-17"

type PolicyVersion string

type PolicyDocument struct {
	Version   PolicyVersion    `json:"Version,omitempty"`
	Statement []StatementEntry `json:"Statement,omitempty"`
}

type StatementEntry struct {
	Sid       string                       `json:"Sid,omitempty"`
	Effect    string                       `json:"Effect,omitempty"`
	Principal map[string]string            `json:"Principal,omitempty"`
	Action    []string                     `json:"Action,omitempty"`
	Resource  []string                     `json:"Resource,omitempty"`
	Condition map[string]map[string]string `json:"Condition,omitempty"`
}

func Client(session client.ConfigProvider) *iam.IAM {
	// Create a IAM service client.
	return iam.New(session)
}
