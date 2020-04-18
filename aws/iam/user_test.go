package iam

import (
	"fmt"
	"strings"
	"testing"

	awssdk "github.com/aws/aws-sdk-go/aws"
	awsarn "github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsiam "github.com/aws/aws-sdk-go/service/iam"
	"github.com/stretchr/testify/assert"

	"github.com/redradrat/cloud-objects/aws"
)

///////////////
// IAM MOCKS //
///////////////

const (
	ReferenceUserName         = "thisismyuser"
	ReferenceExistingUserName = "thisismyexistinguser"
	ReferenceUserId           = "AIDA1234567890EXAMPLE"
	ReferenceAccessKeyId      = "AKIA1234567890EXAMPLE"
	ReferenceAccessKeySecret  = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
)

func (m *mockIAMClient) CreateUser(input *awsiam.CreateUserInput) (*awsiam.CreateUserOutput, error) {
	if *input.UserName == ReferenceExistingUserName {
		return nil, fmt.Errorf("User already exists")
	}

	if strings.Contains(*input.UserName, "/") {
		return nil, fmt.Errorf("malformed User Name")
	}
	// Check if input values are still as we want them to be
	assert.Equal(m.t, getReferenceCreateUserInput(), input)
	out := createMockCreateUserOutput(input)
	return out, nil
}

func (m *mockIAMClient) UpdateUser(input *awsiam.UpdateUserInput) (*awsiam.UpdateUserOutput, error) {

	// Check if input values are still as we want them to be
	assert.Equal(m.t, getReferenceUpdateUserInput(), input)

	return &awsiam.UpdateUserOutput{}, nil
}

func (m *mockIAMClient) DeleteUser(input *awsiam.DeleteUserInput) (*awsiam.DeleteUserOutput, error) {
	// Check if input values are still as we want them to be
	assert.Equal(m.t, awssdk.String(ReferenceExistingUserName), input.UserName)

	return &awsiam.DeleteUserOutput{}, nil
}

func (m *mockIAMClient) GetUser(input *awsiam.GetUserInput) (*awsiam.GetUserOutput, error) {
	var ro *awsiam.GetUserOutput
	if awssdk.StringValue(input.UserName) == FriendlyNamefromARN(getReferenceUserNonExistingArn()) {
		return ro, awserr.New(awsiam.ErrCodeNoSuchEntityException, "", fmt.Errorf("entity not found"))
	}

	assert.Equal(m.t, getReferenceGetUserInput(), input)
	return createMockGetUserOutput(input), nil
}

func (m *mockIAMClient) CreateAccessKey(input *awsiam.CreateAccessKeyInput) (*awsiam.CreateAccessKeyOutput, error) {
	return createMockCreateAccessKeyOutput(input), nil
}

func (m *mockIAMClient) DeleteAccessKey(input *awsiam.DeleteAccessKeyInput) (*awsiam.DeleteAccessKeyOutput, error) {
	assert.Equal(m.t, ReferenceExistingUserName, *input.UserName)
	assert.Equal(m.t, ReferenceAccessKeyId, *input.AccessKeyId)
	return &awsiam.DeleteAccessKeyOutput{}, nil
}

func (m *mockIAMClient) CreateLoginProfile(input *awsiam.CreateLoginProfileInput) (*awsiam.CreateLoginProfileOutput, error) {
	assert.Equal(m.t, false, *input.PasswordResetRequired)
	return createMockCreateLoginProfileOutput(input), nil
}

func (m *mockIAMClient) DeleteLoginProfile(input *awsiam.DeleteLoginProfileInput) (*awsiam.DeleteLoginProfileOutput, error) {
	assert.Equal(m.t, ReferenceExistingUserName, *input.UserName)
	return &awsiam.DeleteLoginProfileOutput{}, nil
}

func createMockCreateAccessKeyOutput(input *awsiam.CreateAccessKeyInput) *awsiam.CreateAccessKeyOutput {
	return &awsiam.CreateAccessKeyOutput{
		AccessKey: &awsiam.AccessKey{
			AccessKeyId:     awssdk.String(ReferenceAccessKeyId),
			CreateDate:      awssdk.Time(getReferenceCreateTimestamp()),
			SecretAccessKey: awssdk.String(ReferenceAccessKeySecret),
			UserName:        input.UserName,
		},
	}
}

func createMockCreateLoginProfileOutput(input *awsiam.CreateLoginProfileInput) *awsiam.CreateLoginProfileOutput {
	return &awsiam.CreateLoginProfileOutput{
		LoginProfile: &awsiam.LoginProfile{
			CreateDate:            awssdk.Time(getReferenceCreateTimestamp()),
			UserName:              input.UserName,
			PasswordResetRequired: input.PasswordResetRequired,
		},
	}
}

func createMockCreateUserOutput(input *awsiam.CreateUserInput) *awsiam.CreateUserOutput {
	return &awsiam.CreateUserOutput{
		User: &awsiam.User{
			Arn:        awssdk.String(getReferenceUserNonExistingArn().String()),
			CreateDate: awssdk.Time(getReferenceCreateTimestamp()),
			Path:       input.Path,
			UserId:     awssdk.String(ReferenceUserId),
			UserName:   input.UserName,
			Tags:       input.Tags,
		},
	}
}

func createMockGetUserOutput(input *awsiam.GetUserInput) *awsiam.GetUserOutput {

	return &awsiam.GetUserOutput{
		User: &awsiam.User{
			Arn:              awssdk.String(getReferenceUserExistingArn().String()),
			CreateDate:       awssdk.Time(getReferenceCreateTimestamp()),
			PasswordLastUsed: awssdk.Time(getReferenceUpdateTimestamp()),
			UserId:           awssdk.String(ReferenceUserId),
			UserName:         input.UserName,
		},
	}
}

///////////
// TESTS //
///////////

func TestUserInstance_Create(t *testing.T) {
	// Setup Test
	mockSvc := &mockIAMClient{t: t}

	usrIns := NewUserInstance(ReferenceUserName, false, false)
	err := usrIns.Create(mockSvc)
	assert.NoError(t, err)
	assert.True(t, usrIns.IsCreated(mockSvc))

	usrIns = NewUserInstance(ReferenceExistingUserName, false, false)
	err = usrIns.Create(mockSvc)
	assert.Error(t, err)

	// AccessKey
	usrIns = NewUserInstance(ReferenceUserName, false, true)
	err = usrIns.Create(mockSvc)
	assert.NoError(t, err)
	assert.True(t, usrIns.IsCreated(mockSvc))
	assert.NotNil(t, usrIns.AccessKey())
	assert.Nil(t, usrIns.LoginProfileCredentials())
	assert.Equal(t, ReferenceAccessKeyId, usrIns.AccessKey().Id())
	assert.Equal(t, ReferenceAccessKeySecret, usrIns.AccessKey().Secret())

	// LoginProfile Credentials
	usrIns = NewUserInstance(ReferenceUserName, true, false)
	err = usrIns.Create(mockSvc)
	assert.NoError(t, err)
	assert.True(t, usrIns.IsCreated(mockSvc))
	assert.NotNil(t, usrIns.LoginProfileCredentials())
	assert.Nil(t, usrIns.AccessKey())
	assert.Equal(t, ReferenceUserName, usrIns.LoginProfileCredentials().Username())
	assert.True(t, len(usrIns.LoginProfileCredentials().Password()) == 20)
}

func TestUserInstance_Update(t *testing.T) {
	// Setup Test
	mockSvc := &mockIAMClient{t: t}

	usrIns := NewUserInstance(ReferenceUserName, false, false)

	err := usrIns.Update(mockSvc)
	assert.Error(t, err)
	assert.True(t, err.(aws.InstanceError).IsOfErrorCode(aws.ErrAWSInstanceNotYetCreated))
	assert.False(t, usrIns.IsCreated(mockSvc))

	usrIns = NewExistingUserInstance(ReferenceUserName, false, false, getReferenceUserExistingArn())
	err = usrIns.Update(mockSvc)
	assert.NoError(t, err)
}

func TestUserInstance_Delete(t *testing.T) {
	// Setup Test
	mockSvc := &mockIAMClient{t: t}

	rolIns := NewUserInstance(ReferenceUserName, false, false)

	err := rolIns.Delete(mockSvc)
	assert.Error(t, err)
	assert.True(t, err.(aws.InstanceError).IsOfErrorCode(aws.ErrAWSInstanceNotYetCreated))
	assert.False(t, rolIns.IsCreated(mockSvc))

	rolIns = NewExistingUserInstance(ReferenceUserName, false, false, getReferenceUserExistingArn())
	err = rolIns.Delete(mockSvc)
	assert.NoError(t, err)
}

func TestNewUserInstance(t *testing.T) {
	pi := NewUserInstance(ReferenceUserName, false, false)
	assert.Equal(t, getReferenceUserInstance(), pi)
}

func TestNewExistingUserInstance(t *testing.T) {
	ri := NewExistingUserInstance(ReferenceUserName, false, false, getReferenceUserExistingArn())
	riWithArn := getReferenceUserInstance()
	riWithArn.arn = getReferenceUserExistingArn()
	assert.Equal(t, riWithArn, ri)
}

func TestUserInstance_ARN(t *testing.T) {
	assert.True(t, awsarn.IsARN(getReferenceUserInstance().ARN().String()))
}

func TestUserInstance_createUser(t *testing.T) {
	// Setup Test
	mockSvc := &mockIAMClient{t: t}

	_, err := createUser(mockSvc, "test/name")
	assert.Error(t, err)
}

/////////////
// HELPERS //
/////////////

func getReferenceUserExistingArn() awsarn.ARN {
	arn, _ := awsarn.Parse(fmt.Sprintf("arn:aws:iam::123456789012:role/%s", ReferenceExistingUserName))
	return arn
}

func getReferenceUserNonExistingArn() awsarn.ARN {
	arn, _ := awsarn.Parse(fmt.Sprintf("arn:aws:iam::123456789012:role/%s-NE", ReferenceUserName))
	return arn
}

func getReferenceGetUserInput() *awsiam.GetUserInput {
	return &awsiam.GetUserInput{UserName: awssdk.String(ReferenceUserName)}
}

func getReferenceCreateUserInput() *awsiam.CreateUserInput {
	return &awsiam.CreateUserInput{
		UserName: awssdk.String(ReferenceUserName),
	}
}

func getReferenceUpdateUserInput() *awsiam.UpdateUserInput {
	return &awsiam.UpdateUserInput{
		NewUserName: awssdk.String(ReferenceUserName),
		UserName:    awssdk.String(ReferenceExistingUserName),
	}
}

func getReferenceUserInstance() *UserInstance {
	return &UserInstance{
		Name: ReferenceUserName,
	}
}
