package iam

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client/metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	session, err := session.NewSession(&awssdk.Config{
		Region: awssdk.String("eu-west-1")},
	)
	assert.NoError(t, err)
	client := Client(session)
	assert.Equal(t, client.ClientInfo, metadata.ClientInfo{
		ServiceName:   "iam",
		ServiceID:     "IAM",
		APIVersion:    "2010-05-08",
		PartitionID:   "aws",
		Endpoint:      "https://iam.amazonaws.com",
		SigningName:   "iam",
		SigningRegion: "us-east-1",
	})
}
