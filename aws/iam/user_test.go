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

func createMockCreateUserOutput(input *awsiam.CreateUserInput) *awsiam.CreateUserOutput {
	return &awsiam.CreateUserOutput{
		User: &awsiam.User{
			Arn:        awssdk.String(getReferencePolicyExistingArn().String()),
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
			Arn:              awssdk.String(getReferencePolicyExistingArn().String()),
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

	rolIns := NewUserInstance(ReferenceUserName)
	err := rolIns.Create(mockSvc)
	assert.NoError(t, err)
	assert.True(t, rolIns.IsCreated(mockSvc))

	rolIns = NewUserInstance(ReferenceExistingUserName)
	err = rolIns.Create(mockSvc)
	assert.Error(t, err)
}

func TestUserInstance_Update(t *testing.T) {
	// Setup Test
	mockSvc := &mockIAMClient{t: t}

	usrIns := NewUserInstance(ReferenceUserName)

	err := usrIns.Update(mockSvc)
	assert.Error(t, err)
	assert.True(t, err.(aws.InstanceError).IsOfErrorCode(aws.ErrAWSInstanceNotYetCreated))
	assert.False(t, usrIns.IsCreated(mockSvc))

	usrIns = NewExistingUserInstance(ReferenceUserName, getReferenceUserExistingArn())
	err = usrIns.Update(mockSvc)
	assert.NoError(t, err)
}

func TestUserInstance_Delete(t *testing.T) {
	// Setup Test
	mockSvc := &mockIAMClient{t: t}

	rolIns := NewUserInstance(ReferenceUserName)

	err := rolIns.Delete(mockSvc)
	assert.Error(t, err)
	assert.True(t, err.(aws.InstanceError).IsOfErrorCode(aws.ErrAWSInstanceNotYetCreated))
	assert.False(t, rolIns.IsCreated(mockSvc))

	rolIns = NewExistingUserInstance(ReferenceUserName, getReferenceUserExistingArn())
	err = rolIns.Delete(mockSvc)
	assert.NoError(t, err)
}

func TestNewUserInstance(t *testing.T) {
	pi := NewUserInstance(ReferenceUserName)
	assert.Equal(t, getReferenceUserInstance(), pi)
}

func TestNewExistingUserInstance(t *testing.T) {
	ri := NewExistingUserInstance(ReferenceUserName, getReferenceUserExistingArn())
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
