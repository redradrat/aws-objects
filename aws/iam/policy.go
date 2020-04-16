package iam

import (
	"encoding/json"
	"fmt"

	awssdk "github.com/aws/aws-sdk-go/aws"
	awsarn "github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"

	"github.com/redradrat/cloud-objects/aws"
)

func createPolicy(svc iamiface.IAMAPI, polName, polDesc string, pd PolicyDocument) (*iam.CreatePolicyOutput, error) {
	b, err := json.Marshal(&pd)
	if err != nil {
		return nil, err
	}

	result, err := svc.CreatePolicy(&iam.CreatePolicyInput{
		PolicyDocument: awssdk.String(string(b)),
		Description:    awssdk.String(polDesc),
		PolicyName:     awssdk.String(polName),
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func updatePolicy(svc iamiface.IAMAPI, policyArn awsarn.ARN, pd PolicyDocument) (*iam.CreatePolicyVersionOutput, error) {
	b, err := json.Marshal(&pd)
	if err != nil {
		return nil, err
	}

	result, err := svc.CreatePolicyVersion(&iam.CreatePolicyVersionInput{
		PolicyDocument: awssdk.String(string(b)),
		PolicyArn:      awssdk.String(policyArn.String()),
		SetAsDefault:   awssdk.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func deletePolicy(svc iamiface.IAMAPI, arn awsarn.ARN) (*iam.DeletePolicyOutput, error) {
	res, err := svc.DeletePolicy(&iam.DeletePolicyInput{
		PolicyArn: awssdk.String(arn.String()),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func getPolicy(svc iamiface.IAMAPI, arn awsarn.ARN) (*iam.GetPolicyOutput, error) {

	result, err := svc.GetPolicy(&iam.GetPolicyInput{
		PolicyArn: awssdk.String(arn.String()),
	})

	if err != nil {
		fmt.Println("Error", err)
		return nil, err
	}

	return result, nil
}

func getPolicyVersion(svc iamiface.IAMAPI, po iam.GetPolicyOutput) (*iam.GetPolicyVersionOutput, error) {

	result, err := svc.GetPolicyVersion(&iam.GetPolicyVersionInput{
		PolicyArn: po.Policy.Arn,
		VersionId: po.Policy.DefaultVersionId,
	})

	if err != nil {
		fmt.Println("Error", err)
		return nil, err
	}

	return result, nil

}

type PolicyInstance struct {
	Name           string
	Description    string
	PolicyDocument PolicyDocument
	arn            awsarn.ARN
}

func NewPolicyInstance(name, description string, policyDoc PolicyDocument) *PolicyInstance {
	return &PolicyInstance{Name: name, Description: description, PolicyDocument: policyDoc}
}

func NewExistingPolicyInstance(svc iamiface.IAMAPI, arn awsarn.ARN) (*PolicyInstance, error) {
	var pi *PolicyInstance
	emptyarn := awsarn.ARN{}.String()
	if arn.String() == emptyarn {
		return pi, fmt.Errorf("given ARN is empty")
	}

	out, err := getPolicy(svc, arn)
	if err != nil {
		return pi, err
	}

	pdout, err := getPolicyVersion(svc, *out)
	if err != nil {
		return pi, err
	}
	var pd PolicyDocument
	json.Unmarshal([]byte(awssdk.StringValue(pdout.PolicyVersion.Document)), &pd)
	pi = &PolicyInstance{
		Name:           awssdk.StringValue(out.Policy.PolicyName),
		Description:    awssdk.StringValue(out.Policy.Description),
		PolicyDocument: pd,
		arn:            arn,
	}

	return pi, nil
}

//  Create attaches the referenced policy on referenced target type and returns the target ARN
func (p *PolicyInstance) Create(svc iamiface.IAMAPI) error {
	var newarn awsarn.ARN
	out, err := createPolicy(svc, p.Name, p.Description, p.PolicyDocument)
	if err != nil {
		return err
	}
	newarn, err = awsarn.Parse(awssdk.StringValue(out.Policy.Arn))
	if err != nil {
		return err
	}
	p.arn = newarn
	return nil
}

// Update for PolicyInstance creates a new Policy version an sets it as active; then returns the arn
func (p *PolicyInstance) Update(svc iamiface.IAMAPI) error {
	if !p.IsCreated() {
		return aws.NewInstanceNotYetCreatedError(fmt.Sprintf("Policy '%s' not yet created", p.Name))
	}

	_, err := updatePolicy(svc, p.arn, p.PolicyDocument)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes the referenced Policy from referenced target type and returns the target ARN
func (p *PolicyInstance) Delete(svc iamiface.IAMAPI) error {
	if !p.IsCreated() {
		return aws.NewInstanceNotYetCreatedError(fmt.Sprintf("Policy '%s' not yet created", p.Name))
	}

	_, err := deletePolicy(svc, p.arn)
	if err != nil {
		return err
	}
	return nil
}

func (p *PolicyInstance) ARN() awsarn.ARN {
	return p.arn
}

func (p *PolicyInstance) IsCreated() bool {
	return p.arn.String() != awsarn.ARN{}.String()
}
