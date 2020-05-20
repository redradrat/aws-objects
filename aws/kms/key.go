package kms

import (
	"fmt"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	awskms "github.com/aws/aws-sdk-go/service/kms"

	"github.com/redradrat/cloud-objects/aws"
	"github.com/redradrat/cloud-objects/cloudobject"
)

const (
	KMSKeyTopic = "ENCKEY"
)

type Key struct {
	name    string
	status  *KeyStatus
	session *awskms.KMS
}

func NewKey(name string, session client.ConfigProvider) (*Key, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("given name is empty")
	}
	if len(name) > 200 {
		return nil, fmt.Errorf("given name is longer than 200 characters")
	}

	key := Key{
		name:    name,
		session: awskms.New(session),
	}

	return &key, nil
}

func (k *Key) Create(spec cloudobject.CloudObjectSpec) (cloudobject.Secrets, error) {
	// It's fair to assume, that we get an KMS KeySpec here.
	assertedSpec, ok := spec.(*KeySpec)
	if !ok {
		// If not, we're throwing an error here... ya done messed up.
		return nil, cloudobject.SpecInvalidError{Message: "got unsupported spec"}
	}

	// If the KMS Key already exists, we're gonna throw an error here... you're trying to play us for a fool!
	exists, err := k.Exists()
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, nil
	}

	// Now let's go for it... create this Key!
	input := assertedSpec.CreateKeyInput()
	out, err := k.session.CreateKey(&input)
	if err != nil {
		return nil, err
	}

	// We specifically need to assign an alias to this key!
	aliasInput := assertedSpec.CreateAliasInput(k.Id().String(), *out.KeyMetadata.KeyId)
	_, err = k.session.CreateAlias(&aliasInput)
	if err != nil {
		return nil, err
	}

	// re-trigger status update
	if err = k.Read(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (k *Key) Read() error {
	// Call AWS to describe our KMS Key
	out, err := k.session.DescribeKey(&awskms.DescribeKeyInput{
		KeyId: k.Id().StringPtr(),
	})
	if err != nil {
		if err.(awserr.Error).Code() == awskms.ErrCodeNotFoundException {
			return cloudobject.NotExistsError{Message: fmt.Sprintf("KMS Key with id '%s' not found",
				k.Id().String())}
		}
		return err
	}
	k.status = (*KeyStatus)(out.KeyMetadata)

	return nil
}

func (k *Key) Update(spec cloudobject.CloudObjectSpec) (cloudobject.Secrets, error) {
	// It's fair to assume, that we get an KMS KeySpec here.
	_, ok := spec.(*KeySpec)
	if !ok {
		return nil, cloudobject.SpecInvalidError{Message: "got unsupported spec"}
	}

	// Let's update our status
	if err := k.Read(); err != nil {
		return nil, err
	}
	// Here we could copy the old status before we read again, and compute a delta
	_ = k.status

	// As the actual Key really has nothing to modify, we're just good here...
	return nil, nil
}

func (k *Key) Delete(purge bool) error {
	// First, let's check whether our KMS Key actually exists
	exists, err := k.Exists()
	if err != nil {
		return err
	}

	// If our KMS Key doesn't exist, there is nothing to do for us here, as we cannot get the KeyId anymore
	if !exists {
		return nil
	}

	// Now let's go for it...
	// First we disable the actual key
	input := awskms.DisableKeyInput{
		KeyId: k.status.KeyId,
	}
	if _, err := k.session.DisableKey(&input); err != nil {
		return err
	}

	// If we decided to purge, we trigger the Key deletion here
	if purge {
		input := awskms.ScheduleKeyDeletionInput{
			KeyId:               k.status.KeyId,
			PendingWindowInDays: awssdk.Int64(7),
		}
		if _, err := k.session.ScheduleKeyDeletion(&input); err != nil {
			return err
		}
	}

	// Secondly we delete the alias, so we're free to "create" that key again
	aliasInput := awskms.DeleteAliasInput{
		AliasName: k.Id().StringPtr(),
	}
	if _, err := k.session.DeleteAlias(&aliasInput); cloudobject.IgnoreNotExistsError(err) != nil {
		return err
	}

	return nil
}

func (k *Key) Status() cloudobject.Status {
	return k.status
}

func (k *Key) Id() cloudobject.Id {
	return cloudobject.Id(fmt.Sprintf("%s/%s", "alias", aws.CloudObjectResource(KMSKeyTopic, k.name)))
}

func (k *Key) Exists() (bool, error) {
	return cloudobject.Exists(k)
}

//////////////
/// SECRET ///
//////////////

type KeySecrets awskms.KeyMetadata

func (secrets *KeySecrets) Map() map[string]string {
	return secrets.Map()
}

//////////////
/// STATUS ///
//////////////

type KeyStatus awskms.KeyMetadata

func (status *KeyStatus) String() string {
	out := awskms.KeyMetadata(*status).String()
	return out
}

////////////
/// SPEC ///
////////////

type KeySpec struct {
	KeyUsage KeyUsage
	KeyType  KeyType
	Policy   *string
	Tags     map[string]string
}

func (spec *KeySpec) Valid() (bool, error) {
	return true, nil
}

func (spec *KeySpec) CreateKeyInput() awskms.CreateKeyInput {
	tags := compileTags(spec.Tags)

	return awskms.CreateKeyInput{
		BypassPolicyLockoutSafetyCheck: awssdk.Bool(false),
		CustomKeyStoreId:               nil,
		CustomerMasterKeySpec:          awssdk.String(spec.KeyType.String()),
		Description:                    nil,
		KeyUsage:                       awssdk.String(spec.KeyUsage.String()),
		Origin:                         nil,
		Policy:                         spec.Policy,
		Tags:                           tags,
	}
}

func (spec *KeySpec) CreateAliasInput(alias, id string) awskms.CreateAliasInput {
	return awskms.CreateAliasInput{
		AliasName:   awssdk.String(alias),
		TargetKeyId: awssdk.String(id),
	}
}
