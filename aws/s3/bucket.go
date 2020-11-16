package s3

import (
	"fmt"
	awssdk "github.com/aws/aws-sdk-go/aws"
	awsarn "github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	awss3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/redradrat/cloud-objects/aws"
	"github.com/redradrat/cloud-objects/aws/kms"
	"github.com/redradrat/cloud-objects/cloudobject"
)

const (
	BucketTopic = "bkt"
	ZeroResultsList = "list-buckets returned 0 results"
)

// Bucket represents the S3 Bucket CloudObject
type Bucket struct {
	name    string
	status  BucketStatus
	session *awss3.S3
}

type BucketStatus struct {
	awss3.Bucket
	Encrypted bool
	ARN string
}
func (status BucketStatus) String() string {
	return status.Bucket.String()
}
func (status BucketStatus) ProviderID() cloudobject.ProviderID {
	return cloudobject.ProviderID{
		Type: cloudobject.AWSProvider,
		Value: status.ARN,
	}
}

type BucketSecrets struct {}
func (secrets BucketSecrets) Map() map[string]string {
	out := make(map[string]string)
	// TODO: Fill map
	return out
}

/////////////
/// FUNCS ///
/////////////

func (b *Bucket) Create(spec cloudobject.CloudObjectSpec) (cloudobject.Secrets, error) {
	var err error
	var keyFound bool

	// Check whether encryption key already exists
	var key *kms.Key
	key, err = kmsKeySession(b)
	if err != nil {
		return nil, err
	}
	keyFound, err = key.Exists()
	if err != nil {
		return nil, err
	}

	// It's fair to assume, that we get an S3 BucketSpec here.
	assertedSpec, ok := spec.(*BucketSpec)
	if !ok {
		// If not, we're throwing an error here... ya done messed up.
		return nil, cloudobject.SpecInvalidError{Message: "got unsupported spec"}
	}

	// If the S3 Bucket already exists, we're done here... you're trying to play us for a fool!
	exists, err := b.Exists()
	if err != nil {
		return nil, err
	}
	if !exists {
		// So now we should be good to go ahead with Bucket creation
		input := assertedSpec.CreateBucketInput(b.ID().String())
		_, err = b.session.CreateBucket(&input)
		if err != nil {
			return nil, err
		}
	}

	// As there are a few post-creation settings we call our bucket config helper. This is externalized to serve for
	// Update() as well.
	err = ensureBucketConfig(assertedSpec, b)
	if err != nil {
		return nil, err
	}

	// Ensure Bucket Encryption
	if !keyFound {
		_, err := key.Create(&kms.KeySpec{
			KeyUsage: kms.EncryptDecryptKeyUsage,
			KeyType:  kms.SymmetricDefaultKeyType,
		})
		if err != nil {
			return nil, err
		}
	}
	encinput := assertedSpec.PutBucketEncryptionInput(b.ID().String(), key)
	_, err = b.session.PutBucketEncryption(&encinput)
	if err != nil {
		return nil, err
	}

	// TODO: Compile secret
	return nil, nil
}

func kmsKeySession(b *Bucket) (*kms.Key, error) {
	kmsSession, err := session.NewSession(&b.session.Config)
	if err != nil {
		return nil, err
	}
	key, err := kms.NewKey(b.name, kmsSession)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func ensureBucketConfig(assertedSpec *BucketSpec, b *Bucket) error {
	var err error

	// Ensure Bucket ACL
	aclinput := assertedSpec.PutBucketAclInput(b.ID().String())
	_, err = b.session.PutBucketAcl(&aclinput)
	if err != nil {
		return err
	}

	// Ensure Bucket Versioning
	versinput := assertedSpec.PutBucketVersioningInput(b.ID().String())
	_, err = b.session.PutBucketVersioning(&versinput)
	if err != nil {
		return err
	}

	// Ensure Bucket Transfer Acceleration
	accelinput := assertedSpec.PutBucketAccelerationInput(b.ID().String())
	_, err = b.session.PutBucketAccelerateConfiguration(&accelinput)
	if err != nil {
		return err
	}

	// Ensure Bucket Public Block
	blockinput := assertedSpec.PutPublicAccessBlockInput(b.ID().String())
	_, err = b.session.PutPublicAccessBlock(&blockinput)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bucket) Read() error {
	// Call AWS to describe our S3 Bucket
	out, err := b.session.ListBuckets(&awss3.ListBucketsInput{})
	if err != nil {
		return err
	}
	// If our output bucket list is 0, there is literally no buckets and thus also not what we're looking for
	if len(out.Buckets) == 0 {
		return cloudobject.NotExistsError{Message: fmt.Sprintf(ZeroResultsList)}
	}

	// Range over all returned bucket objects to find ours. If not found, we error here.
	var foundBucket *awss3.Bucket
	for _, bucket := range out.Buckets {
		if *bucket.Name == *b.ID().StringPtr() {
			foundBucket = bucket
		}
	}
	if foundBucket == nil {
		return cloudobject.NotExistsError{Message: fmt.Sprintf("Bucket with id '%s' not found", b.ID().String())}
	}

	// Set the status to what we get from the object. In bucket case this is just the name basically. So not really
	// worth that much. But let's be such good bois here.
	b.status.Bucket = *foundBucket

	// Construct the ARN for status
	b.status.ARN = awsarn.ARN{
		Partition: b.session.PartitionID,
		Service:   b.session.ServiceID,
		Resource:  b.ID().String(),
	}.String()

	enc, err := b.session.GetBucketEncryption(&awss3.GetBucketEncryptionInput{
		Bucket: b.ID().StringPtr(),
	})
	if err != nil {
		return err
	} else {
		b.status.Encrypted = len(enc.ServerSideEncryptionConfiguration.Rules) != 0
	}

	return nil
}

func (b *Bucket) Update(spec cloudobject.CloudObjectSpec) (cloudobject.Secrets, error) {
	// It's fair to assume, that we get an S3 BucketSpec here.
	assertedSpec, ok := spec.(*BucketSpec)
	if !ok {
		// If not, we're throwing an error here... ya done messed up.
		return nil, cloudobject.SpecInvalidError{Message: "got unsupported spec"}
	}

	// Ensure updatable config is set
	err := ensureBucketConfig(assertedSpec, b)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Delete deletes the bucket. Purge has no effect, bucket will always be purged.
func (b *Bucket) Delete(_ bool) error {
	// First, let's check whether our bucket actually exists
	exists, err := b.Exists()
	if err != nil {
		return err
	}

	// If our bucket doesn't exist, there is nothing to do for us here
	if !exists {
		return nil
	}

	// compile DeleteBucketInput
	input := awss3.DeleteBucketInput{
		Bucket: b.ID().StringPtr(),
	}

	// Now let's go for it... delete that naughty Bucket!! (kill it with fire, pwetty please)
	if _, err := b.session.DeleteBucket(&input); cloudobject.IgnoreNotExistsError(err) != nil {
		return err
	}

	return nil

}

func (b *Bucket) Status() cloudobject.Status {
	return b.status
}

func (b *Bucket) ID() cloudobject.ID {
	return cloudobject.ID(aws.CloudObjectResource(BucketTopic, b.name))
}

func (b *Bucket) Exists() (bool, error) {
	return cloudobject.Exists(b)
}

////////////
/// SPEC ///
////////////

type GrantsSpec struct {

	// User will have full access to the bucket incl ACL privileges
	FullAccess []string

	// User will have write access to the bucket
	WriteAccess []string

	// User will have read access to the bucket
	ReadAccess []string

	// User will be allowed to write the ACL
	ACLWriteAccess []string

	// User will be allowed to read the ACL
	ACLReadAccess []string
}

type BucketSpec struct {

	// The Location to create the bucket in. Defaults to us-east-1.
	Location string

	// The canned ACL to apply to the bucket. (e.g. "private", "public-read", "public-read-write", "authenticated-read")
	ACL string

	// ObjectLock enables object locking capabilities on the bucket
	ObjectLock bool

	// Enables object versioning capabilities on the bucket
	Versioning bool

	// TransferAcceleration enables transfer acceleration capabilities on the bucket
	TransferAcceleration bool

	// Acl restrictions
	BlockPublicAcls bool
	IgnorePublicAcls bool

	// Policy restrictions
	BlockPublicPolicy     bool
	RestrictPublicBuckets bool

	//// Grants is a spec to grant AWS IAM Users access to different levels
	//Grants GrantsSpec
}

func (b BucketSpec) Valid() (bool, error) {
	return true, nil
}


///////////////
/// HELPERS ///
///////////////

// NewInstance returns a new S3 Bucket object
func NewBucket(name string, session client.ConfigProvider) (*Bucket, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("given name is empty")
	}
	if len(name) > 200 {
		return nil, fmt.Errorf("given name is longer than 200 characters")
	}

	bucket := Bucket{
		name:    name,
		session: awss3.New(session),
	}

	return &bucket, nil
}

func (b BucketSpec) CreateBucketInput(id string) awss3.CreateBucketInput {
	in := awss3.CreateBucketInput{
		Bucket:                     awssdk.String(id),
		ObjectLockEnabledForBucket: awssdk.Bool(b.ObjectLock),
	}
	if b.Location != "" {
		in.CreateBucketConfiguration = &awss3.CreateBucketConfiguration{LocationConstraint: awssdk.String(b.Location)}

	}
	return in
}

func (b BucketSpec) PutBucketAclInput(id string) awss3.PutBucketAclInput {
	in := awss3.PutBucketAclInput{
		ACL:    awssdk.String(b.ACL),
		Bucket: awssdk.String(id),
	}
	return in
}

func (b BucketSpec) PutBucketVersioningInput(id string) awss3.PutBucketVersioningInput {
	status := "Suspended"
	if b.Versioning {
		status = "Enabled"
	}

	in := awss3.PutBucketVersioningInput{
		Bucket: awssdk.String(id),
		VersioningConfiguration: &awss3.VersioningConfiguration{
			Status:    awssdk.String(status),
		},
	}
	return in

}

func (b BucketSpec) PutBucketAccelerationInput(id string) awss3.PutBucketAccelerateConfigurationInput {
	status := "Suspended"
	if b.TransferAcceleration {
		status = "Enabled"
	}

	in := awss3.PutBucketAccelerateConfigurationInput{
		AccelerateConfiguration: &awss3.AccelerateConfiguration{
			Status: awssdk.String(status),
		},
		Bucket: awssdk.String(id),
	}
	return in
}

func (b BucketSpec) PutBucketEncryptionInput(id string, key *kms.Key) awss3.PutBucketEncryptionInput {
	in := awss3.PutBucketEncryptionInput{
		ServerSideEncryptionConfiguration: &awss3.ServerSideEncryptionConfiguration{
			Rules: []*awss3.ServerSideEncryptionRule{{
				ApplyServerSideEncryptionByDefault: &awss3.ServerSideEncryptionByDefault{
					KMSMasterKeyID: key.ID().StringPtr(),
					SSEAlgorithm:   awssdk.String(awss3.ServerSideEncryptionAwsKms),
				},
			}},
		},
		Bucket: awssdk.String(id),
	}
	return in
}

func (b BucketSpec) PutPublicAccessBlockInput(id string) awss3.PutPublicAccessBlockInput {
	in := awss3.PutPublicAccessBlockInput{
		PublicAccessBlockConfiguration: &awss3.PublicAccessBlockConfiguration{
			BlockPublicAcls:       awssdk.Bool(b.BlockPublicAcls),
			BlockPublicPolicy:     awssdk.Bool(b.BlockPublicPolicy),
			IgnorePublicAcls:      awssdk.Bool(b.IgnorePublicAcls),
			RestrictPublicBuckets: awssdk.Bool(b.RestrictPublicBuckets),
		},
		Bucket: awssdk.String(id),
	}
	return in
}
