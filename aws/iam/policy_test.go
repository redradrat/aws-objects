package iam

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	awsarn "github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsiam "github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/stretchr/testify/assert"

	"github.com/redradrat/cloud-objects/aws"
)

///////////////
// IAM MOCKS //
///////////////

type mockIAMClient struct {
	iamiface.IAMAPI
	t *testing.T
}

const (
	ReferencePolicyName        = "thisismypolicy"
	ReferencePolicyDescription = "description"
	ReferenceV1VersionId       = "v1"
	ReferenceV2VersionId       = "v2"
	ReferencePolicyId          = "ANPA1234567890EXAMPLE"
)

func (m *mockIAMClient) CreatePolicy(input *awsiam.CreatePolicyInput) (*awsiam.CreatePolicyOutput, error) {
	if strings.Contains(*input.PolicyName, "/") {
		return nil, fmt.Errorf("malformed Policy Name")
	}
	// Check if input values are still as we want them to be
	assert.Equal(m.t, getReferenceCreatePolicyInput(), input)
	return createMockCreatePolicyOutput(input), nil
}

func (m *mockIAMClient) CreatePolicyVersion(input *awsiam.CreatePolicyVersionInput) (*awsiam.CreatePolicyVersionOutput, error) {
	// Check if input values are still as we want them to be
	assert.Equal(m.t, getReferenceCreatePolicyVersionInput(), input)
	out := createMockCreatePolicyVersionOutput(input)
	return out, nil
}

func (m *mockIAMClient) DeletePolicy(input *awsiam.DeletePolicyInput) (*awsiam.DeletePolicyOutput, error) {
	// Check if input values are still as we want them to be
	assert.True(m.t, arn.IsARN(*input.PolicyArn))
	assert.Equal(m.t, awssdk.String(getReferencePolicyExistingArn().String()), input.PolicyArn)

	return &awsiam.DeletePolicyOutput{}, nil
}

func (m *mockIAMClient) GetPolicy(input *awsiam.GetPolicyInput) (*awsiam.GetPolicyOutput, error) {
	var po *awsiam.GetPolicyOutput
	if awssdk.StringValue(input.PolicyArn) == getReferencePolicyNonExistingArn().String() {
		return po, awserr.New(awsiam.ErrCodeNoSuchEntityException, "", fmt.Errorf("entity not found"))
	}

	assert.Equal(m.t, getReferenceGetPolicyInput(), input)
	return createMockGetPolicyOutput(input), nil
}

func (m *mockIAMClient) GetPolicyVersion(input *awsiam.GetPolicyVersionInput) (*awsiam.GetPolicyVersionOutput, error) {
	var pvo *awsiam.GetPolicyVersionOutput
	if awssdk.StringValue(input.PolicyArn) == getReferencePolicyNonExistingArn().String() {
		return pvo, awserr.New(awsiam.ErrCodeNoSuchEntityException, "", fmt.Errorf("entity not found"))
	}

	assert.Equal(m.t, getReferenceGetPolicyVersionInput(), input)
	return createMockGetPolicyVersionOutput(input), nil
}

func createMockCreatePolicyOutput(input *awsiam.CreatePolicyInput) *awsiam.CreatePolicyOutput {
	return &awsiam.CreatePolicyOutput{
		Policy: &awsiam.Policy{
			Arn:              awssdk.String(getReferencePolicyExistingArn().String()),
			AttachmentCount:  awssdk.Int64(0),
			CreateDate:       awssdk.Time(getReferenceCreateTimestamp()),
			DefaultVersionId: awssdk.String(ReferenceV1VersionId),
			Description:      input.Description,
			IsAttachable:     awssdk.Bool(true),
			Path:             input.Path,
			PolicyId:         awssdk.String(ReferencePolicyId),
			PolicyName:       awssdk.String(ReferencePolicyName),
			UpdateDate:       awssdk.Time(getReferenceCreateTimestamp()),
		},
	}
}

func createMockGetPolicyVersionOutput(input *awsiam.GetPolicyVersionInput) *awsiam.GetPolicyVersionOutput {
	return &awsiam.GetPolicyVersionOutput{
		PolicyVersion: &awsiam.PolicyVersion{
			CreateDate:       awssdk.Time(getReferenceUpdateTimestamp()),
			Document:         awssdk.String(string(getMarshaledReferencePolicyDocument())),
			IsDefaultVersion: awssdk.Bool(true),
			VersionId:        awssdk.String(ReferenceV2VersionId),
		},
	}
}

func createMockCreatePolicyVersionOutput(input *awsiam.CreatePolicyVersionInput) *awsiam.CreatePolicyVersionOutput {
	return &awsiam.CreatePolicyVersionOutput{
		PolicyVersion: &awsiam.PolicyVersion{
			CreateDate:       awssdk.Time(getReferenceUpdateTimestamp()),
			Document:         input.PolicyDocument,
			IsDefaultVersion: input.SetAsDefault,
			VersionId:        awssdk.String(ReferenceV2VersionId),
		},
	}
}

func createMockGetPolicyOutput(input *awsiam.GetPolicyInput) *awsiam.GetPolicyOutput {
	return &awsiam.GetPolicyOutput{
		Policy: &awsiam.Policy{
			Arn:              input.PolicyArn,
			AttachmentCount:  awssdk.Int64(1),
			CreateDate:       awssdk.Time(getReferenceCreateTimestamp()),
			DefaultVersionId: awssdk.String(ReferenceV2VersionId),
			Description:      awssdk.String(ReferencePolicyDescription),
			IsAttachable:     awssdk.Bool(true),
			PolicyName:       awssdk.String(ReferencePolicyName),
			UpdateDate:       awssdk.Time(getReferenceUpdateTimestamp()),
		},
	}
}

///////////
// TESTS //
///////////

func TestPolicyInstance_Create(t *testing.T) {
	// Setup Test
	mockSvc := &mockIAMClient{t: t}

	pd := getReferencePolicyDocument()
	polIns := NewPolicyInstance(ReferencePolicyName, ReferencePolicyDescription, pd)
	err := polIns.Create(mockSvc)
	assert.NoError(t, err)
	assert.True(t, polIns.IsCreated())
}

func TestPolicyInstance_Update(t *testing.T) {
	// Setup Test
	mockSvc := &mockIAMClient{t: t}

	pd := getReferencePolicyDocument()
	polIns := NewPolicyInstance(ReferencePolicyName, ReferencePolicyDescription, pd)

	err := polIns.Update(mockSvc)
	assert.Error(t, err)
	assert.True(t, err.(aws.InstanceError).IsOfErrorCode(aws.ErrAWSInstanceNotYetCreated))
	assert.False(t, polIns.IsCreated())

	polIns = NewExistingPolicyInstance(ReferencePolicyName, ReferencePolicyDescription, pd, getReferencePolicyExistingArn())
	err = polIns.Update(mockSvc)
	assert.NoError(t, err)
}

func TestPolicyInstance_Delete(t *testing.T) {
	// Setup Test
	mockSvc := &mockIAMClient{t: t}

	pd := getReferencePolicyDocument()
	polIns := NewPolicyInstance(ReferencePolicyName, ReferencePolicyDescription, pd)

	err := polIns.Delete(mockSvc)
	assert.Error(t, err)
	assert.True(t, err.(aws.InstanceError).IsOfErrorCode(aws.ErrAWSInstanceNotYetCreated))
	assert.False(t, polIns.IsCreated())

	polIns = NewExistingPolicyInstance(ReferencePolicyName, ReferencePolicyDescription, pd, getReferencePolicyExistingArn())
	err = polIns.Delete(mockSvc)
	assert.NoError(t, err)
}

func TestNewPolicyInstance(t *testing.T) {
	pi := NewPolicyInstance(ReferencePolicyName, ReferencePolicyDescription, getReferencePolicyDocument())
	assert.Equal(t, getReferencePolicyInstance(), pi)
}

func TestNewPolicyInstanceFromExisting(t *testing.T) {
	pi := NewExistingPolicyInstance(ReferencePolicyName, ReferencePolicyDescription, getReferencePolicyDocument(), getReferencePolicyExistingArn())
	piWithArn := getReferencePolicyInstance()
	piWithArn.arn = getReferencePolicyExistingArn()
	assert.Equal(t, piWithArn, pi)
}

func TestPolicyInstance_ARN(t *testing.T) {
	assert.True(t, awsarn.IsARN(getReferencePolicyInstance().ARN().String()))
}

func TestPolicyInstance_createPolicy(t *testing.T) {
	// Setup Test
	mockSvc := &mockIAMClient{t: t}

	_, err := createPolicy(mockSvc, "test/name", "test", getReferencePolicyDocument())
	assert.Error(t, err)
}

/////////////
// HELPERS //
/////////////

func getReferencePolicyExistingArn() awsarn.ARN {
	arn, _ := awsarn.Parse(fmt.Sprintf("arn:aws:iam::123456789012:policy/%s", ReferencePolicyName))
	return arn
}

func getReferencePolicyNonExistingArn() awsarn.ARN {
	arn, _ := awsarn.Parse(fmt.Sprintf("arn:aws:iam::123456789012:policy/%s-NE", ReferencePolicyName))
	return arn
}

func getReferenceCreateTimestamp() time.Time {
	return time.Date(2020, time.April, 14, 20, 0, 0, 0, time.UTC)
}

func getReferenceUpdateTimestamp() time.Time {
	return time.Date(2020, time.April, 14, 21, 0, 0, 0, time.UTC)
}

func getReferenceGetPolicyInput() *awsiam.GetPolicyInput {
	return &awsiam.GetPolicyInput{PolicyArn: awssdk.String(getReferencePolicyExistingArn().String())}
}

func getReferenceGetPolicyVersionInput() *awsiam.GetPolicyVersionInput {
	return &awsiam.GetPolicyVersionInput{
		PolicyArn: awssdk.String(getReferencePolicyExistingArn().String()),
		VersionId: awssdk.String(ReferenceV2VersionId),
	}
}

func getReferenceCreatePolicyInput() *awsiam.CreatePolicyInput {
	return &awsiam.CreatePolicyInput{
		Description:    awssdk.String(ReferencePolicyDescription),
		PolicyDocument: awssdk.String(string(getMarshaledReferencePolicyDocument())),
		PolicyName:     awssdk.String(ReferencePolicyName),
	}
}

func getReferenceCreatePolicyVersionInput() *awsiam.CreatePolicyVersionInput {
	return &awsiam.CreatePolicyVersionInput{
		PolicyArn:      awssdk.String(getReferencePolicyExistingArn().String()),
		PolicyDocument: awssdk.String(string(getMarshaledReferencePolicyDocument())),
		SetAsDefault:   awssdk.Bool(true),
	}
}

func getReferencePolicyInstance() *PolicyInstance {
	return &PolicyInstance{
		Name:           ReferencePolicyName,
		Description:    ReferencePolicyDescription,
		PolicyDocument: getReferencePolicyDocument(),
	}
}

func getReferencePolicyDocument() PolicyDocument {
	return PolicyDocument{
		Version: PolicyVersion20121017,
		Statement: []StatementEntry{
			{
				Effect:    "Allow",
				Action:    []string{"god:Mode"},
				Resource:  []string{"*"},
				Condition: nil,
			},
		},
	}
}

func getMarshaledReferencePolicyDocument() []byte {
	marshaledPolicy := getReferencePolicyDocument()
	b, _ := json.Marshal(&marshaledPolicy)
	return b
}
