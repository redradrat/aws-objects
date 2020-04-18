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
	ReferenceExistingGroupName = "thisismyexistinggroup"
	ReferenceGroupName         = "thisismygroup"
	ReferenceGroupId           = "AGPA1234567890EXAMPLE"
)

func (m *mockIAMClient) CreateGroup(input *awsiam.CreateGroupInput) (*awsiam.CreateGroupOutput, error) {
	if *input.GroupName == ReferenceExistingGroupName {
		return nil, fmt.Errorf("Group already exists")
	}

	if strings.Contains(*input.GroupName, "/") {
		return nil, fmt.Errorf("malformed Group Name")
	}
	// Check if input values are still as we want them to be
	assert.Equal(m.t, getReferenceCreateGroupInput(), input)
	out := createMockCreateGroupOutput(input)
	return out, nil
}

func (m *mockIAMClient) UpdateGroup(input *awsiam.UpdateGroupInput) (*awsiam.UpdateGroupOutput, error) {

	// Check if input values are still as we want them to be
	assert.Equal(m.t, getReferenceUpdateGroupInput(), input)

	return &awsiam.UpdateGroupOutput{}, nil
}

func (m *mockIAMClient) DeleteGroup(input *awsiam.DeleteGroupInput) (*awsiam.DeleteGroupOutput, error) {
	// Check if input values are still as we want them to be
	assert.Equal(m.t, awssdk.String(ReferenceExistingGroupName), input.GroupName)

	return &awsiam.DeleteGroupOutput{}, nil
}

func (m *mockIAMClient) GetGroup(input *awsiam.GetGroupInput) (*awsiam.GetGroupOutput, error) {
	var ro *awsiam.GetGroupOutput
	if awssdk.StringValue(input.GroupName) == FriendlyNamefromARN(getReferenceGroupNonExistingArn()) {
		return ro, awserr.New(awsiam.ErrCodeNoSuchEntityException, "", fmt.Errorf("entity not found"))
	}

	assert.Equal(m.t, getReferenceGetGroupInput(), input)
	return createMockGetGroupOutput(input), nil
}

func createMockCreateGroupOutput(input *awsiam.CreateGroupInput) *awsiam.CreateGroupOutput {
	return &awsiam.CreateGroupOutput{
		Group: &awsiam.Group{
			Arn:        awssdk.String(getReferencePolicyExistingArn().String()),
			CreateDate: awssdk.Time(getReferenceCreateTimestamp()),
			Path:       input.Path,
			GroupId:    awssdk.String(ReferenceGroupId),
			GroupName:  input.GroupName,
		},
	}
}

func createMockGetGroupOutput(input *awsiam.GetGroupInput) *awsiam.GetGroupOutput {

	return &awsiam.GetGroupOutput{
		Group: &awsiam.Group{
			Arn:        awssdk.String(getReferencePolicyExistingArn().String()),
			CreateDate: awssdk.Time(getReferenceCreateTimestamp()),
			GroupId:    awssdk.String(ReferenceGroupId),
			GroupName:  input.GroupName,
		},
	}
}

///////////
// TESTS //
///////////

func TestGroupInstance_Create(t *testing.T) {
	// Setup Test
	mockSvc := &mockIAMClient{t: t}

	grpIns := NewGroupInstance(ReferenceGroupName)
	err := grpIns.Create(mockSvc)
	assert.NoError(t, err)
	assert.True(t, grpIns.IsCreated(mockSvc))

	grpIns = NewGroupInstance(ReferenceExistingGroupName)
	err = grpIns.Create(mockSvc)
	assert.Error(t, err)
}

func TestGroupInstance_Update(t *testing.T) {
	// Setup Test
	mockSvc := &mockIAMClient{t: t}

	grpIns := NewGroupInstance(ReferenceGroupName)

	err := grpIns.Update(mockSvc)
	assert.Error(t, err)
	assert.True(t, err.(aws.InstanceError).IsOfErrorCode(aws.ErrAWSInstanceNotYetCreated))
	assert.False(t, grpIns.IsCreated(mockSvc))

	grpIns = NewExistingGroupInstance(ReferenceGroupName, getReferenceGroupExistingArn())
	err = grpIns.Update(mockSvc)
	assert.NoError(t, err)
}

func TestGroupInstance_Delete(t *testing.T) {
	// Setup Test
	mockSvc := &mockIAMClient{t: t}

	grpIns := NewGroupInstance(ReferenceGroupName)

	err := grpIns.Delete(mockSvc)
	assert.Error(t, err)
	assert.True(t, err.(aws.InstanceError).IsOfErrorCode(aws.ErrAWSInstanceNotYetCreated))
	assert.False(t, grpIns.IsCreated(mockSvc))

	grpIns = NewExistingGroupInstance(ReferenceGroupName, getReferenceGroupExistingArn())
	err = grpIns.Delete(mockSvc)
	assert.NoError(t, err)
}

func TestNewGroupInstance(t *testing.T) {
	pi := NewGroupInstance(ReferenceGroupName)
	assert.Equal(t, getReferenceGroupInstance(), pi)
}

func TestNewExistingGroupInstance(t *testing.T) {
	ri := NewExistingGroupInstance(ReferenceGroupName, getReferenceGroupExistingArn())
	riWithArn := getReferenceGroupInstance()
	riWithArn.arn = getReferenceGroupExistingArn()
	assert.Equal(t, riWithArn, ri)
}

func TestGroupInstance_ARN(t *testing.T) {
	assert.True(t, awsarn.IsARN(getReferenceGroupInstance().ARN().String()))
}

func TestGroupInstance_createGroup(t *testing.T) {
	// Setup Test
	mockSvc := &mockIAMClient{t: t}

	_, err := createGroup(mockSvc, "test/name")
	assert.Error(t, err)
}

/////////////
// HELPERS //
/////////////

func getReferenceGroupExistingArn() awsarn.ARN {
	arn, _ := awsarn.Parse(fmt.Sprintf("arn:aws:iam::123456789012:role/%s", ReferenceExistingGroupName))
	return arn
}

func getReferenceGroupNonExistingArn() awsarn.ARN {
	arn, _ := awsarn.Parse(fmt.Sprintf("arn:aws:iam::123456789012:role/%s-NE", ReferenceGroupName))
	return arn
}

func getReferenceGetGroupInput() *awsiam.GetGroupInput {
	return &awsiam.GetGroupInput{GroupName: awssdk.String(ReferenceGroupName)}
}

func getReferenceCreateGroupInput() *awsiam.CreateGroupInput {
	return &awsiam.CreateGroupInput{
		GroupName: awssdk.String(ReferenceGroupName),
	}
}

func getReferenceUpdateGroupInput() *awsiam.UpdateGroupInput {
	return &awsiam.UpdateGroupInput{
		GroupName:    awssdk.String(ReferenceExistingGroupName),
		NewGroupName: awssdk.String(ReferenceGroupName),
	}
}

func getReferenceGroupInstance() *GroupInstance {
	return &GroupInstance{
		Name: ReferenceGroupName,
	}
}
