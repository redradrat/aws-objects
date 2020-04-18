package iam

import (
	"fmt"

	awssdk "github.com/aws/aws-sdk-go/aws"
	awsarn "github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsiam "github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"

	"github.com/redradrat/cloud-objects/aws"
)

func createUser(svc iamiface.IAMAPI, userName string) (*awsiam.CreateUserOutput, error) {

	result, err := svc.CreateUser(&awsiam.CreateUserInput{
		UserName: awssdk.String(userName),
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func updateUser(svc iamiface.IAMAPI, userName string, arn awsarn.ARN) (*awsiam.UpdateUserOutput, error) {

	result, err := svc.UpdateUser(&awsiam.UpdateUserInput{
		NewUserName: awssdk.String(userName),
		UserName:    awssdk.String(FriendlyNamefromARN(arn)),
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func deleteUser(svc iamiface.IAMAPI, arn awsarn.ARN) (*awsiam.DeleteUserOutput, error) {

	res, err := svc.DeleteUser(&awsiam.DeleteUserInput{
		UserName: awssdk.String(FriendlyNamefromARN(arn)),
	})
	if err != nil {
		if err.(awserr.Error).Code() != awsiam.ErrCodeNoSuchEntityException {
			return nil, err
		}
	}

	return res, nil
}

func getUser(svc iamiface.IAMAPI, userArn awsarn.ARN) (*awsiam.GetUserOutput, error) {

	result, err := svc.GetUser(&awsiam.GetUserInput{
		UserName: awssdk.String(FriendlyNamefromARN(userArn)),
	})

	if err != nil {
		fmt.Println("Error", err)
		return nil, err
	}

	return result, nil
}

type UserInstance struct {
	Name string
	arn  awsarn.ARN
}

func NewUserInstance(name string) *UserInstance {
	return &UserInstance{Name: name}
}

func NewExistingUserInstance(name string, arn awsarn.ARN) *UserInstance {
	return &UserInstance{
		Name: name,
		arn:  arn,
	}
}

// Create creates an AWS User
func (u *UserInstance) Create(svc iamiface.IAMAPI) error {
	var newarn awsarn.ARN
	out, err := createUser(svc, u.Name)
	if err != nil {
		return err
	}
	newarn, err = awsarn.Parse(awssdk.StringValue(out.User.Arn))
	if err != nil {
		return err
	}
	u.arn = newarn
	return nil
}

func (u *UserInstance) Update(svc iamiface.IAMAPI) error {
	if !u.IsCreated(svc) {
		return aws.NewInstanceNotYetCreatedError(fmt.Sprintf("User '%s' not yet created", u.Name))
	}

	_, err := updateUser(svc, u.Name, u.arn)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserInstance) Delete(svc iamiface.IAMAPI) error {
	if !u.IsCreated(svc) {
		return aws.NewInstanceNotYetCreatedError(fmt.Sprintf("User '%s' not yet created", u.Name))
	}

	_, err := deleteUser(svc, u.arn)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserInstance) ARN() awsarn.ARN {
	return u.arn
}

func (u *UserInstance) IsCreated(svc iamiface.IAMAPI) bool {
	return u.arn.String() != awsarn.ARN{}.String()
}
