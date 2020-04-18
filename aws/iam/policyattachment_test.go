package iam

import (
	"fmt"
	"testing"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	awsarn "github.com/aws/aws-sdk-go/aws/arn"
	awsiam "github.com/aws/aws-sdk-go/service/iam"
	"github.com/stretchr/testify/assert"

	"github.com/redradrat/cloud-objects/aws"
)

///////////////
// IAM MOCKS //
///////////////

func (m *mockIAMClient) AttachRolePolicy(input *awsiam.AttachRolePolicyInput) (*awsiam.AttachRolePolicyOutput, error) {
	if *input.RoleName == ReferenceRoleName {
		return &awsiam.AttachRolePolicyOutput{}, fmt.Errorf("referenced Role does not exist")
	}

	assert.Equal(m.t, ReferenceExistingRoleName, *input.RoleName)
	assert.Equal(m.t, getReferencePolicyExistingArn().String(), *input.PolicyArn)
	// Check if input values are still as we want them to be
	assert.Equal(m.t, getReferenceAttachRolePolicyInput(), input)
	return createMockAttachRolePolicyOutput(input), nil
}

func (m *mockIAMClient) AttachUserPolicy(input *awsiam.AttachUserPolicyInput) (*awsiam.AttachUserPolicyOutput, error) {
	if *input.UserName == ReferenceUserName {
		return &awsiam.AttachUserPolicyOutput{}, fmt.Errorf("referenced User does not exist")
	}

	assert.Equal(m.t, ReferenceExistingUserName, *input.UserName)
	assert.Equal(m.t, getReferencePolicyExistingArn().String(), *input.PolicyArn)

	// Check if input values are still as we want them to be
	assert.Equal(m.t, getReferenceAttachUserPolicyInput(), input)
	return createMockAttachUserPolicyOutput(input), nil
}

func (m *mockIAMClient) AttachGroupPolicy(input *awsiam.AttachGroupPolicyInput) (*awsiam.AttachGroupPolicyOutput, error) {
	if *input.GroupName == ReferenceGroupName {
		return &awsiam.AttachGroupPolicyOutput{}, fmt.Errorf("referenced Group does not exist")
	}

	assert.Equal(m.t, ReferenceExistingGroupName, *input.GroupName)
	assert.Equal(m.t, getReferencePolicyExistingArn().String(), *input.PolicyArn)
	// Check if input values are still as we want them to be
	assert.Equal(m.t, getReferenceAttachGroupPolicyInput(), input)
	return createMockAttachGroupPolicyOutput(input), nil
}

func (m *mockIAMClient) ListAttachedRolePolicies(input *awsiam.ListAttachedRolePoliciesInput) (*awsiam.ListAttachedRolePoliciesOutput, error) {
	return createMockListAttachedRolePoliciesOutput(input), nil
}

func (m *mockIAMClient) ListAttachedUserPolicies(input *awsiam.ListAttachedUserPoliciesInput) (*awsiam.ListAttachedUserPoliciesOutput, error) {
	return createMockListAttachedUserPoliciesOutput(input), nil
}

func (m *mockIAMClient) ListAttachedGroupPolicies(input *awsiam.ListAttachedGroupPoliciesInput) (*awsiam.ListAttachedGroupPoliciesOutput, error) {
	return createMockListAttachedGroupPoliciesOutput(input), nil
}

func (m *mockIAMClient) DetachRolePolicy(input *awsiam.DetachRolePolicyInput) (*awsiam.DetachRolePolicyOutput, error) {
	if *input.RoleName == ReferenceRoleName {
		return &awsiam.DetachRolePolicyOutput{}, aws.NewInstanceError(aws.ErrAWSInstanceNotYetCreated, "referenced Role does not exist")
	}

	if *input.PolicyArn == getReferencePolicyNonExistingArn().String() {
		return &awsiam.DetachRolePolicyOutput{}, aws.NewInstanceError(aws.ErrAWSInstanceNotYetCreated, "referenced Policy does not exist")
	}

	// Check if input values are still as we want them to be
	assert.True(m.t, arn.IsARN(*input.PolicyArn))
	assert.Equal(m.t, awssdk.String(getReferencePolicyExistingArn().String()), input.PolicyArn)
	assert.Equal(m.t, ReferenceExistingRoleName, *input.RoleName)

	return &awsiam.DetachRolePolicyOutput{}, nil
}

func (m *mockIAMClient) DetachUserPolicy(input *awsiam.DetachUserPolicyInput) (*awsiam.DetachUserPolicyOutput, error) {
	if *input.UserName == ReferenceUserName {
		return &awsiam.DetachUserPolicyOutput{}, aws.NewInstanceError(aws.ErrAWSInstanceNotYetCreated, "referenced Role does not exist")
	}

	if *input.PolicyArn == getReferencePolicyNonExistingArn().String() {
		return &awsiam.DetachUserPolicyOutput{}, aws.NewInstanceError(aws.ErrAWSInstanceNotYetCreated, "referenced Policy does not exist")
	}

	// Check if input values are still as we want them to be
	assert.True(m.t, arn.IsARN(*input.PolicyArn))
	assert.Equal(m.t, awssdk.String(getReferencePolicyExistingArn().String()), input.PolicyArn)
	assert.Equal(m.t, ReferenceExistingUserName, *input.UserName)

	return &awsiam.DetachUserPolicyOutput{}, nil
}

func (m *mockIAMClient) DetachGroupPolicy(input *awsiam.DetachGroupPolicyInput) (*awsiam.DetachGroupPolicyOutput, error) {
	if *input.GroupName == ReferenceGroupName {
		return &awsiam.DetachGroupPolicyOutput{}, aws.NewInstanceError(aws.ErrAWSInstanceNotYetCreated, "referenced Role does not exist")
	}

	if *input.PolicyArn == getReferencePolicyNonExistingArn().String() {
		return &awsiam.DetachGroupPolicyOutput{}, aws.NewInstanceError(aws.ErrAWSInstanceNotYetCreated, "referenced Policy does not exist")
	}

	// Check if input values are still as we want them to be
	assert.True(m.t, arn.IsARN(*input.PolicyArn))
	assert.Equal(m.t, awssdk.String(getReferencePolicyExistingArn().String()), input.PolicyArn)
	assert.Equal(m.t, ReferenceExistingGroupName, *input.GroupName)

	return &awsiam.DetachGroupPolicyOutput{}, nil
}

func createMockAttachRolePolicyOutput(input *awsiam.AttachRolePolicyInput) *awsiam.AttachRolePolicyOutput {
	return &awsiam.AttachRolePolicyOutput{}
}

func createMockAttachUserPolicyOutput(input *awsiam.AttachUserPolicyInput) *awsiam.AttachUserPolicyOutput {
	return &awsiam.AttachUserPolicyOutput{}
}

func createMockAttachGroupPolicyOutput(input *awsiam.AttachGroupPolicyInput) *awsiam.AttachGroupPolicyOutput {
	return &awsiam.AttachGroupPolicyOutput{}
}

func createMockListAttachedRolePoliciesOutput(input *awsiam.ListAttachedRolePoliciesInput) *awsiam.ListAttachedRolePoliciesOutput {
	return &awsiam.ListAttachedRolePoliciesOutput{
		AttachedPolicies: []*awsiam.AttachedPolicy{
			{
				PolicyArn:  awssdk.String(getReferencePolicyExistingArn().String()),
				PolicyName: awssdk.String(ReferencePolicyName),
			},
		},
	}
}

func createMockListAttachedUserPoliciesOutput(input *awsiam.ListAttachedUserPoliciesInput) *awsiam.ListAttachedUserPoliciesOutput {
	return &awsiam.ListAttachedUserPoliciesOutput{
		AttachedPolicies: []*awsiam.AttachedPolicy{
			{
				PolicyArn:  awssdk.String(getReferencePolicyExistingArn().String()),
				PolicyName: awssdk.String(ReferencePolicyName),
			},
		},
	}
}

func createMockListAttachedGroupPoliciesOutput(input *awsiam.ListAttachedGroupPoliciesInput) *awsiam.ListAttachedGroupPoliciesOutput {
	return &awsiam.ListAttachedGroupPoliciesOutput{
		AttachedPolicies: []*awsiam.AttachedPolicy{
			{
				PolicyArn:  awssdk.String(getReferencePolicyExistingArn().String()),
				PolicyName: awssdk.String(ReferencePolicyName),
			},
		},
	}
}

///////////
// TESTS //
///////////

func TestPolicyAttachmentInstance_Create(t *testing.T) {
	// Setup Test
	mockSvc := &mockIAMClient{t: t}

	paIns := NewPolicyAttachmentInstance(getReferencePolicyExistingArn(), RoleAttachmentType, getReferenceRoleExistingArn())
	err := paIns.Create(mockSvc)
	assert.NoError(t, err)
	assert.True(t, paIns.IsCreated(mockSvc))

	paIns = NewPolicyAttachmentInstance(getReferencePolicyExistingArn(), UserAttachmentType, getReferenceUserExistingArn())
	err = paIns.Create(mockSvc)
	assert.NoError(t, err)
	assert.True(t, paIns.IsCreated(mockSvc))

	paIns = NewPolicyAttachmentInstance(getReferencePolicyExistingArn(), GroupAttachmentType, getReferenceGroupExistingArn())
	err = paIns.Create(mockSvc)
	assert.NoError(t, err)
	assert.True(t, paIns.IsCreated(mockSvc))
}

func TestPolicyAttachmentInstance_Update(t *testing.T) {
	// Setup Test
	mockSvc := &mockIAMClient{t: t}

	paIns := NewPolicyAttachmentInstance(getReferencePolicyExistingArn(), RoleAttachmentType, getReferenceRoleNonExistingArn())

	err := paIns.Update(mockSvc)
	assert.NoError(t, err)
}

func TestPolicyAttachmentInstance_Delete(t *testing.T) {
	// Setup Test
	mockSvc := &mockIAMClient{t: t}

	// Target not existing
	paIns := NewPolicyAttachmentInstance(getReferenceRoleExistingArn(), RoleAttachmentType, getReferenceRoleNonExistingArn())
	err := paIns.Delete(mockSvc)
	assert.Error(t, err)
	assert.True(t, err.(aws.InstanceError).IsOfErrorCode(aws.ErrAWSInstanceNotYetCreated))
	assert.False(t, paIns.IsCreated(mockSvc))

	paIns = NewPolicyAttachmentInstance(getReferenceRoleExistingArn(), UserAttachmentType, getReferenceUserNonExistingArn())
	err = paIns.Delete(mockSvc)
	assert.Error(t, err)
	assert.True(t, err.(aws.InstanceError).IsOfErrorCode(aws.ErrAWSInstanceNotYetCreated))
	assert.False(t, paIns.IsCreated(mockSvc))

	paIns = NewPolicyAttachmentInstance(getReferenceRoleExistingArn(), GroupAttachmentType, getReferenceGroupNonExistingArn())
	err = paIns.Delete(mockSvc)
	assert.Error(t, err)
	assert.True(t, err.(aws.InstanceError).IsOfErrorCode(aws.ErrAWSInstanceNotYetCreated))
	assert.False(t, paIns.IsCreated(mockSvc))

	// Policy non existing
	paIns = NewPolicyAttachmentInstance(getReferenceRoleNonExistingArn(), RoleAttachmentType, getReferenceRoleExistingArn())
	err = paIns.Delete(mockSvc)
	assert.Error(t, err)
	assert.True(t, err.(aws.InstanceError).IsOfErrorCode(aws.ErrAWSInstanceNotYetCreated))
	assert.False(t, paIns.IsCreated(mockSvc))

	paIns = NewPolicyAttachmentInstance(getReferenceRoleNonExistingArn(), UserAttachmentType, getReferenceUserExistingArn())
	err = paIns.Delete(mockSvc)
	assert.Error(t, err)
	assert.True(t, err.(aws.InstanceError).IsOfErrorCode(aws.ErrAWSInstanceNotYetCreated))
	assert.False(t, paIns.IsCreated(mockSvc))

	paIns = NewPolicyAttachmentInstance(getReferenceRoleNonExistingArn(), GroupAttachmentType, getReferenceGroupExistingArn())
	err = paIns.Delete(mockSvc)
	assert.Error(t, err)
	assert.True(t, err.(aws.InstanceError).IsOfErrorCode(aws.ErrAWSInstanceNotYetCreated))
	assert.False(t, paIns.IsCreated(mockSvc))

	// All good
	paIns = NewPolicyAttachmentInstance(getReferencePolicyExistingArn(), RoleAttachmentType, getReferenceRoleExistingArn())
	err = paIns.Delete(mockSvc)
	assert.NoError(t, err)

	paIns = NewPolicyAttachmentInstance(getReferencePolicyExistingArn(), UserAttachmentType, getReferenceUserExistingArn())
	err = paIns.Delete(mockSvc)
	assert.NoError(t, err)

	paIns = NewPolicyAttachmentInstance(getReferencePolicyExistingArn(), GroupAttachmentType, getReferenceGroupExistingArn())
	err = paIns.Delete(mockSvc)
	assert.NoError(t, err)
}

func TestNewPolicyAttachmentInstance(t *testing.T) {
	pi := NewPolicyAttachmentInstance(getReferencePolicyExistingArn(), RoleAttachmentType, getReferenceRoleExistingArn())
	assert.Equal(t, getReferencePolicyAttachmentInstance(), pi)
}

func TestPolicyAttachmentInstance_ARN(t *testing.T) {
	assert.True(t, awsarn.IsARN(getReferencePolicyAttachmentInstance().ARN().String()))
}

/////////////
// HELPERS //
/////////////

func getReferencePolicyAttachmentInstance() *PolicyAttachmentInstance {
	return &PolicyAttachmentInstance{
		PolicyRef: getReferencePolicyExistingArn(),
		Type:      RoleAttachmentType,
		TargetRef: getReferenceRoleExistingArn(),
	}
}

func getReferenceAttachRolePolicyInput() *awsiam.AttachRolePolicyInput {
	return &awsiam.AttachRolePolicyInput{
		PolicyArn: awssdk.String(getReferencePolicyExistingArn().String()),
		RoleName:  awssdk.String(ReferenceExistingRoleName),
	}
}

func getReferenceAttachUserPolicyInput() *awsiam.AttachUserPolicyInput {
	return &awsiam.AttachUserPolicyInput{
		PolicyArn: awssdk.String(getReferencePolicyExistingArn().String()),
		UserName:  awssdk.String(ReferenceExistingUserName),
	}
}

func getReferenceAttachGroupPolicyInput() *awsiam.AttachGroupPolicyInput {
	return &awsiam.AttachGroupPolicyInput{
		PolicyArn: awssdk.String(getReferencePolicyExistingArn().String()),
		GroupName: awssdk.String(ReferenceExistingGroupName),
	}
}
