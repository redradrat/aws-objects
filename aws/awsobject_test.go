package aws

import (
	"reflect"
	"testing"

	awsarn "github.com/aws/aws-sdk-go/aws/arn"
	"github.com/stretchr/testify/assert"
)

const TestArnString = "arn:aws:iam::123456789012:policy/somepolicyname"

func TestARNify(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    []awsarn.ARN
		wantErr bool
	}{
		{name: "SingleEntry", args: []string{TestArnString}, want: []awsarn.ARN{getReferenceARN()}, wantErr: false},
		{name: "MultipleEntries", args: []string{TestArnString, TestArnString}, want: []awsarn.ARN{getReferenceARN(), getReferenceARN()}, wantErr: false},
		{name: "NoEntries", args: []string{}, want: []awsarn.ARN{}, wantErr: false},
		{name: "NonArnSingleEntry", args: []string{"Gustl"}, want: []awsarn.ARN{}, wantErr: true},
		{name: "NonArnMultiEntry", args: []string{"Gustl", "is", TestArnString}, want: []awsarn.ARN{}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ARNify(tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ARNify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ARNify() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func getReferenceARN() awsarn.ARN {
	arn, _ := awsarn.Parse(TestArnString)
	return arn
}

func getTestInstanceError() InstanceError {
	return NewInstanceError(ErrAWSInstanceNotYetCreated, "somestring")
}

func TestInstanceError(t *testing.T) {
	assert.Equal(t, "somestring", getTestInstanceError().Error())
	assert.True(t, getTestInstanceError().IsOfErrorCode(ErrAWSInstanceNotYetCreated))
}

func TestNewInstanceError(t *testing.T) {
	assert.Equal(t, InstanceError{
		Code:    ErrAWSInstanceNotYetCreated,
		Message: "somestring",
	}, getTestInstanceError())
}

func TestNewInstanceNotYetCreatedError(t *testing.T) {
	assert.Equal(t, InstanceError{
		Code:    ErrAWSInstanceNotYetCreated,
		Message: "differentstring",
	}, NewInstanceNotYetCreatedError("differentstring"))
}
