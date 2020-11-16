package cloudobject

import awssdk "github.com/aws/aws-sdk-go/aws"

const (
	AWSProvider = "aws"
)

// CloudObject is the interface for interacting with any cloud provider resource abstracted in this repository
type CloudObject interface {
	Create(CloudObjectSpec) (Secrets, error)
	Read() error
	Update(CloudObjectSpec) (Secrets, error)
	Delete(bool) error
	Status() Status
	ID() ID
	Exists() (bool, error)
}

type Status interface {
	String() string
	ProviderID() ProviderID
}

// CloudObjectSpec should be an interface that Object Specs should implement
type CloudObjectSpec interface {
	Valid() (bool, error)
}

type ProviderType string
type ProviderID struct {
	Type ProviderType
	Value string
}

func (id ProviderID) String() string {
	return id.Value
}

func (id ProviderID) StringPtr() *string {
	return awssdk.String(id.Value)
}

type ID string

func (id ID) String() string {
	return string(id)
}

func (id ID) StringPtr() *string {
	return awssdk.String(string(id))
}

type Store interface {
	Persist(object *CloudObject) error
	Retrieve(id ID) *CloudObject
}

type Secrets interface {
	Map() map[string]string
}
