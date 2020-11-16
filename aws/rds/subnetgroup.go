package rds

import (
	"fmt"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	awsrds "github.com/aws/aws-sdk-go/service/rds"

	"github.com/redradrat/cloud-objects/aws"
	"github.com/redradrat/cloud-objects/cloudobject"
)

const (
	DBSubnetGroupTopic = "sg"
)

type SubnetGroup struct {
	name    string
	status  *SubnetGroupStatus
	session *awsrds.RDS
}

func NewSubnetGroup(name string, session client.ConfigProvider) (*SubnetGroup, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("given name is empty")
	}
	if len(name) > 200 {
		return nil, fmt.Errorf("given name is longer than 200 characters")
	}

	sg := SubnetGroup{
		name:    name,
		session: awsrds.New(session),
	}

	return &sg, nil
}

func (s *SubnetGroup) Create(spec cloudobject.CloudObjectSpec) (cloudobject.Secrets, error) {
	// It's fair to assume, that we get an RDS SubnetGroupSpec here.
	assertedSpec, ok := spec.(*SubnetGroupSpec)
	if !ok {
		// If not, we're throwing an error here... ya done messed up.
		return nil, cloudobject.SpecInvalidError{Message: "got unsupported spec"}
	}

	// If the SubnetGroup already exists, we're done here... you're trying to play us for a fool!
	exists, err := s.Exists()
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, nil
	}

	// Now let's go for it... create this SubnetGroup!
	input := assertedSpec.CreateDBSubnetGroupInput(s.ID().String())
	_, err = s.session.CreateDBSubnetGroup(&input)
	if err != nil {
		return nil, err
	}

	// re-trigger status update
	if err = s.Read(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *SubnetGroup) Read() error {
	// Call AWS to describe our DB SubnetGroup
	out, err := s.session.DescribeDBSubnetGroups(&awsrds.DescribeDBSubnetGroupsInput{
		DBSubnetGroupName: s.ID().StringPtr(),
	})
	if err != nil {
		if err.(awserr.Error).Code() == awsrds.ErrCodeDBSubnetGroupNotFoundFault {
			return cloudobject.NotExistsError{Message: fmt.Sprintf("RDS DB SubnetGroup with id '%s' not found",
				s.ID().String())}
		}
		return err
	}
	// If our output DB list is 0, we didn't find any matches -> not exists
	if len(out.DBSubnetGroups) == 0 {
		return cloudobject.NotExistsError{Message: fmt.Sprintf("RDS DB SubnetGroup with id '%s' not found",
			s.ID().String())}
	}
	if len(out.DBSubnetGroups) < 1 {
		return cloudobject.AmbiguousIdentifierError{Message: fmt.Sprintf(
			"multiple RDS DB SubnetGroups with id '%s' found", s.ID().String())}
	}
	s.status = (*SubnetGroupStatus)(out.DBSubnetGroups[0])

	return nil
}

func (s *SubnetGroup) Update(spec cloudobject.CloudObjectSpec) (cloudobject.Secrets, error) {
	// It's fair to assume, that we get an RDS SubnetGroupSpec here.
	assertedSpec, ok := spec.(*SubnetGroupSpec)
	if !ok {
		return nil, cloudobject.SpecInvalidError{Message: "got unsupported spec"}
	}

	// Let's update our status
	if err := s.Read(); err != nil {
		return nil, err
	}
	// Here we could copy the old status before we read again, and compute a delta
	_ = s.status

	// Now let's go for it... Modify the actual DB SubnetGroup
	input := assertedSpec.ModifyDBSubnetGroupInput(s.ID().String())
	_, err := s.session.ModifyDBSubnetGroup(&input)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Delete deletes the SubnetGroup. Purge has no effect, subnet group will always be purged.
func (s *SubnetGroup) Delete(_ bool) error {
	// First, let's check whether our SubnetGroup actually exists
	exists, err := s.Exists()
	if err != nil {
		return err
	}

	// If DB SubnetGroup doesn't exist, there is nothing to do for us here
	if !exists {
		return nil
	}

	input := awsrds.DeleteDBSubnetGroupInput{
		DBSubnetGroupName: s.ID().StringPtr(),
	}
	// Now let's go for it... delete that naughty SubnetGroup!! (kill it with fire, pwetty please)
	if _, err := s.session.DeleteDBSubnetGroup(&input); cloudobject.IgnoreNotExistsError(err) != nil {
		return err
	}

	return nil
}

// What does it even mean to exist? Do we all even exist?
func (s *SubnetGroup) Exists() (bool, error) {
	return cloudobject.Exists(s)
}

func (s *SubnetGroup) ID() cloudobject.ID {
	return cloudobject.ID(aws.CloudObjectResource(DBSubnetGroupTopic, s.name))
}

func (s *SubnetGroup) Status() cloudobject.Status {
	return s.status
}

////////////
/// SPEC ///
////////////

type SubnetGroupSpec struct {
	Description string

	SubnetIDs []string

	Tags map[string]string
}

func (spec *SubnetGroupSpec) Valid() (bool, error) {
	return true, nil
}

func (spec *SubnetGroupSpec) CreateDBSubnetGroupInput(id string) awsrds.CreateDBSubnetGroupInput {

	tags := awsTags(spec)

	out := awsrds.CreateDBSubnetGroupInput{
		DBSubnetGroupDescription: awssdk.String(spec.Description),
		DBSubnetGroupName:        awssdk.String(id),
		SubnetIds:                awssdk.StringSlice(spec.SubnetIDs),
		Tags:                     tags,
	}

	return out
}

func (spec *SubnetGroupSpec) ModifyDBSubnetGroupInput(id string) awsrds.ModifyDBSubnetGroupInput {
	out := awsrds.ModifyDBSubnetGroupInput{
		DBSubnetGroupDescription: awssdk.String(spec.Description),
		DBSubnetGroupName:        awssdk.String(id),
		SubnetIds:                awssdk.StringSlice(spec.SubnetIDs),
	}

	return out
}

///////////////
/// HELPERS ///
///////////////

type SubnetGroupStatus awsrds.DBSubnetGroup

func (status *SubnetGroupStatus) ProviderID() cloudobject.ProviderID {
	out := awsrds.DBSubnetGroup(*status)
	return cloudobject.ProviderID{
		Type:  cloudobject.AWSProvider,
		Value: awssdk.StringValue(out.DBSubnetGroupArn),
	}
}

func (status *SubnetGroupStatus) String() string {
	return awsrds.DBSubnetGroup(*status).String()
}

func awsTags(spec *SubnetGroupSpec) []*awsrds.Tag {
	var tags []*awsrds.Tag
	for k, v := range spec.Tags {
		tags = append(tags, &awsrds.Tag{
			Key:   awssdk.String(k),
			Value: awssdk.String(v),
		})
	}
	return tags
}
