package iam

import (
	"encoding/json"
	"fmt"
	"strings"

	awssdk "github.com/aws/aws-sdk-go/aws"
	awsarn "github.com/aws/aws-sdk-go/aws/arn"
	awsiam "github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"

	"github.com/redradrat/cloud-objects/aws"
)

func createRole(svc iamiface.IAMAPI, rn string, roleDesc string, pd PolicyDocument) (*awsiam.CreateRoleOutput, error) {

	b, err := json.Marshal(&pd)
	if err != nil {
		return nil, err
	}

	result, err := svc.CreateRole(&awsiam.CreateRoleInput{
		AssumeRolePolicyDocument: awssdk.String(string(b)),
		Description:              awssdk.String(roleDesc),
		RoleName:                 awssdk.String(rn),
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func updateRole(svc iamiface.IAMAPI, roleArn awsarn.ARN, roleDesc string) (*awsiam.UpdateRoleOutput, error) {

	result, err := svc.UpdateRole(&awsiam.UpdateRoleInput{
		Description: awssdk.String(roleDesc),
		RoleName:    awssdk.String(FriendlyNamefromARN(roleArn)),
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func deleteRole(svc iamiface.IAMAPI, roleArn awsarn.ARN) (*awsiam.DeleteRoleOutput, error) {

	res, err := svc.DeleteRole(&awsiam.DeleteRoleInput{
		RoleName: awssdk.String(FriendlyNamefromARN(roleArn)),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func getRole(svc iamiface.IAMAPI, roleArn awsarn.ARN) (*awsiam.GetRoleOutput, error) {

	result, err := svc.GetRole(&awsiam.GetRoleInput{
		RoleName: awssdk.String(FriendlyNamefromARN(roleArn)),
	})

	if err != nil {
		fmt.Println("Error", err)
		return nil, err
	}

	return result, nil
}

func FriendlyNamefromARN(arn awsarn.ARN) string {
	splitres := strings.Split(arn.Resource, "/")
	return splitres[len(splitres)-1]
}

type RoleInstance struct {
	Name           string
	Description    string
	PolicyDocument PolicyDocument
	arn            awsarn.ARN
}

func NewRoleInstance(name string, description string, poldoc PolicyDocument) *RoleInstance {
	return &RoleInstance{Name: name, Description: description, PolicyDocument: poldoc}
}

func NewExistingRoleInstance(svc iamiface.IAMAPI, arn awsarn.ARN) (*RoleInstance, error) {
	var ri *RoleInstance
	emptyarn := awsarn.ARN{}.String()
	if arn.String() == emptyarn {
		return ri, fmt.Errorf("given ARN is empty")
	}

	out, err := getRole(svc, arn)
	if err != nil {
		return ri, err
	}

	var pd PolicyDocument
	json.Unmarshal([]byte(awssdk.StringValue(out.Role.AssumeRolePolicyDocument)), &pd)
	ri = &RoleInstance{
		Name:           awssdk.StringValue(out.Role.RoleName),
		Description:    awssdk.StringValue(out.Role.Description),
		PolicyDocument: pd,
		arn:            arn,
	}

	return ri, nil
}

// Reconcile creates or updates an AWS Role
func (r *RoleInstance) Create(svc iamiface.IAMAPI) error {
	var newarn awsarn.ARN
	out, err := createRole(svc, r.Name, r.Description, r.PolicyDocument)
	if err != nil {
		return err
	}
	newarn, err = awsarn.Parse(awssdk.StringValue(out.Role.Arn))
	if err != nil {
		return err
	}
	r.arn = newarn
	return nil
}

func (r *RoleInstance) Update(svc iamiface.IAMAPI) error {
	if !r.IsCreated() {
		return aws.NewInstanceNotYetCreatedError(fmt.Sprintf("Role '%s' not yet created", r.Name))
	}

	_, err := updateRole(svc, r.arn, r.Description)
	if err != nil {
		return err
	}
	return nil
}

func (r *RoleInstance) Delete(svc iamiface.IAMAPI) error {
	if !r.IsCreated() {
		return aws.NewInstanceNotYetCreatedError(fmt.Sprintf("Role '%s' not yet created", r.Name))
	}

	_, err := deleteRole(svc, r.arn)
	if err != nil {
		return err
	}
	return nil
}

func (r *RoleInstance) ARN() awsarn.ARN {
	return r.arn
}

func (r *RoleInstance) IsCreated() bool {
	return r.arn.String() != awsarn.ARN{}.String()
}
