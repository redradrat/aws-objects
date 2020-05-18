package cloudobject

import awssdk "github.com/aws/aws-sdk-go/aws"

// CloudObject is the interface for interacting with any cloud provider resource abstracted in this repository
type CloudObject interface {
	Create(CloudObjectSpec) (Secrets, error)
	Read() error
	Update(CloudObjectSpec) (Secrets, error)
	Delete(purge bool) error
	Status() Status
	Id() Id
	Exists() (bool, error)
}

type Status interface {
	String() string
}

// CloudObjectSpec should be an interface that Object Specs should implement
type CloudObjectSpec interface {
	Valid() (bool, error)
}

type Id string

func (id Id) String() string {
	return string(id)
}

func (id Id) StringPtr() *string {
	return awssdk.String(string(id))
}

type Store interface {
	Persist(object *CloudObject) error
	Retrieve(id Id) *CloudObject
}

type Secrets interface {
	Map() map[string]string
}
