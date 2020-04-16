package iam

import (
	"fmt"

	awssdk "github.com/aws/aws-sdk-go/aws"
	awsarn "github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"

	"github.com/redradrat/cloud-objects/aws"
)

func createPolicyAttachment(svc iamiface.IAMAPI, attachType AttachmentType, policyArn, targetArn awsarn.ARN) error {
	var err error
	switch attachType {
	case RoleAttachmentType:
		_, err = svc.AttachRolePolicy(&iam.AttachRolePolicyInput{
			PolicyArn: awssdk.String(policyArn.String()),
			RoleName:  awssdk.String(FriendlyNamefromARN(targetArn)),
		})
	case UserAttachmentType:
		_, err = svc.AttachUserPolicy(&iam.AttachUserPolicyInput{
			PolicyArn: awssdk.String(policyArn.String()),
			UserName:  awssdk.String(FriendlyNamefromARN(targetArn)),
		})
	default:
		return aws.NewInstanceError(ErrAttachmentTypeUnknown, fmt.Sprintf("unknown attachment type '%s", attachType))
	}
	if err != nil {
		return err
	}

	return nil
}

func deletePolicyAttachment(svc iamiface.IAMAPI, attachType AttachmentType, policyArn, targetArn awsarn.ARN) error {
	var err error

	switch attachType {
	case RoleAttachmentType:
		_, err = svc.DetachRolePolicy(&iam.DetachRolePolicyInput{
			PolicyArn: awssdk.String(policyArn.String()),
			RoleName:  awssdk.String(FriendlyNamefromARN(targetArn)),
		})
	case UserAttachmentType:
		_, err = svc.DetachUserPolicy(&iam.DetachUserPolicyInput{
			PolicyArn: awssdk.String(policyArn.String()),
			UserName:  awssdk.String(FriendlyNamefromARN(targetArn)),
		})
	default:
		return aws.NewInstanceError(ErrAttachmentTypeUnknown, fmt.Sprintf("unknown attachment type '%s", attachType))
	}
	if err != nil {
		if err.(awserr.Error).Code() != iam.ErrCodeNoSuchEntityException {
			return err
		}
	}

	return nil
}

const (
	RoleAttachmentType       AttachmentType = "role"
	UserAttachmentType       AttachmentType = "user"
	ErrAttachmentTypeUnknown aws.ErrorCode  = "AttachmentType is unknown"
)

type AttachmentType string

type PolicyAttachmentInstance struct {
	PolicyRef awsarn.ARN
	Type      AttachmentType
	TargetRef awsarn.ARN
	created   bool
}

func NewPolicyAttachmentInstance(policyRef awsarn.ARN, attType AttachmentType, ref awsarn.ARN) *PolicyAttachmentInstance {
	return &PolicyAttachmentInstance{PolicyRef: policyRef, TargetRef: ref, Type: attType}
}

//  Create attaches the referenced policy on referenced target type
func (pa *PolicyAttachmentInstance) Create(svc iamiface.IAMAPI) error {
	if err := createPolicyAttachment(svc, pa.Type, pa.PolicyRef, pa.TargetRef); err != nil {
		return err
	}
	pa.created = true
	return nil
}

// Update for PolicyAttachmentInstance doesn't do anything
func (pa *PolicyAttachmentInstance) Update(svc iamiface.IAMAPI) error {
	// PolicyAttachment not updateable
	return nil
}

// Delete removes the referenced Policy from referenced target type
func (pa *PolicyAttachmentInstance) Delete(svc iamiface.IAMAPI) error {
	if !pa.IsCreated() {
		return aws.NewInstanceNotYetCreatedError("PolicyAttachment not yet created")
	}

	if err := deletePolicyAttachment(svc, pa.Type, pa.PolicyRef, pa.TargetRef); err != nil {
		return err
	}

	return nil
}

// Returns the target ARN when created
func (pa *PolicyAttachmentInstance) ARN() awsarn.ARN {
	return pa.TargetRef
}

func (pa *PolicyAttachmentInstance) IsCreated() bool {
	return pa.created
}
