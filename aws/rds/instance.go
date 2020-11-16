package rds

import (
	"fmt"

	awssdk "github.com/aws/aws-sdk-go/aws"
	awsarn "github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	awsrds "github.com/aws/aws-sdk-go/service/rds"

	"github.com/redradrat/cloud-objects/aws"
	"github.com/redradrat/cloud-objects/aws/kms"
	"github.com/redradrat/cloud-objects/cloudobject"
)

const (
	PreDeleteDBSnapshotTopic = "predelete"
	DBInstanceTopic          = "db"
)

// Instance represents the RDS Instance CloudObject
type Instance struct {
	name    string
	status  *InstanceStatus
	session *awsrds.RDS
}

// NewInstance returns a new RDS instance object
func NewInstance(name string, session client.ConfigProvider) (*Instance, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("given name is empty")
	}
	if len(name) > 200 {
		return nil, fmt.Errorf("given name is longer than 200 characters")
	}

	ins := Instance{
		name:    name,
		session: awsrds.New(session),
	}

	return &ins, nil
}

// Get the CloudObjectId for our Instance. Equals to Instance Name. This is not the AWS Id.
func (i *Instance) ID() cloudobject.ID {
	return cloudobject.ID(aws.CloudObjectResource(DBInstanceTopic, i.name))
}

// Create our RDS Instance for realsies
func (i *Instance) Create(spec cloudobject.CloudObjectSpec) (cloudobject.Secrets, error) {
	var err error

	// It's fair to assume, that we get an RDS InstanceSpec here.
	assertedSpec, ok := spec.(*InstanceSpec)
	if !ok {
		// If not, we're throwing an error here... ya done messed up.
		return nil, cloudobject.SpecInvalidError{Message: "got unsupported spec"}
	}

	// If the RDS Instance already exists, we're done here... you're trying to play us for a fool!
	exists, err := i.Exists()
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, nil
	}

	var snapshotFound bool
	var keyFound bool

	// Check whether snapshot already exists
	snapshotFound, err = snapshotExists(i)
	if err != nil {
		return nil, err
	}

	// Check whether encryption key already exists
	var key *kms.Key
	key, err = kmsKeySession(i)
	if err != nil {
		return nil, err
	}
	keyFound, err = key.Exists()
	if err != nil {
		return nil, err
	}

	restore := false
	// Set restore to true if both key and snapshot were found
	if keyFound && snapshotFound {
		restore = true
	}

	if restore {
		// If create mode is restore, but RestorationDisabled is true, we need to throw an error here
		if assertedSpec.RestorationDisabled {
			return nil, RestorationDisabledError{Message: fmt.Sprintf("creation without restoration triggered, "+
				"but key and snapshot exist for RDS Instance '%s'", i.ID().String())}
		}
		// As we found our preexisting key and snapshot, we just assume we need to restore our stuff
		input := assertedSpec.RestoreDBInstanceFromDBSnapshotInput(i.ID().String(), finalDBSnapshotName(i))
		_, err := i.session.RestoreDBInstanceFromDBSnapshot(&input)
		if err != nil {
			return nil, err
		}
	} else {
		// If snapshot was found, but no key, we need to error out.
		if snapshotFound {
			return nil, cloudobject.NotExistsError{Message: fmt.Sprintf(
				"RDS snaphshot with id '%s' already exists, but no key found", finalDBSnapshotName(i))}
		}

		// Let's create our key, if it doesn't already exist.
		if !keyFound {
			_, err := key.Create(&kms.KeySpec{
				KeyUsage: kms.EncryptDecryptKeyUsage,
				KeyType:  kms.SymmetricDefaultKeyType,
			})
			if err != nil {
				return nil, err
			}
		}

		// So now we should be good to go ahead with DB creation
		input := assertedSpec.CreateDBInstanceInput(i.ID().String())
		input.KmsKeyId = key.ID().StringPtr()
		_, err = i.session.CreateDBInstance(&input)
		if err != nil {
			return nil, err
		}
	}

	// re-trigger status update
	if err = i.Read(); err != nil {
		return nil, err
	}

	// TODO: create instance secrets

	return nil, nil
}

func kmsKeySession(i *Instance) (*kms.Key, error) {
	kmsSession, err := session.NewSession(&i.session.Config)
	if err != nil {
		return nil, err
	}
	key, err := kms.NewKey(i.name, kmsSession)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (i *Instance) Read() error {
	// Call AWS to describe our DB Instance
	out, err := i.session.DescribeDBInstances(&awsrds.DescribeDBInstancesInput{
		DBInstanceIdentifier: i.ID().StringPtr(),
	})
	if err != nil {
		if err.(awserr.Error).Code() == awsrds.ErrCodeDBInstanceNotFoundFault {
			return cloudobject.NotExistsError{Message: fmt.Sprintf("RDS DB Instance with id '%s' not found", i.ID().String())}
		}
		return err
	}
	// If our output DB list is 0, we didn't find any matches -> not exists
	if len(out.DBInstances) == 0 {
		return cloudobject.NotExistsError{Message: fmt.Sprintf("RDS DB Instance with id '%s' not found", i.ID().String())}
	}
	if len(out.DBInstances) < 1 {
		return cloudobject.AmbiguousIdentifierError{Message: fmt.Sprintf(
			"multiple RDS DB Instance with id '%s' found", i.ID().String())}
	}
	i.status = (*InstanceStatus)(out.DBInstances[0])

	return nil
}

func (i *Instance) Update(spec cloudobject.CloudObjectSpec) (cloudobject.Secrets, error) {
	assertedSpec, ok := spec.(*InstanceSpec)
	if !ok {
		return nil, cloudobject.SpecInvalidError{Message: "got unsupported spec"}
	}
	if err := i.Read(); err != nil {
		return nil, err
	}
	// Here we could copy the old instance before we read again, and compute a delta
	oldStatus := i.status

	if awssdk.BoolValue(oldStatus.MultiAZ) {
		if assertedSpec.DBSubnetGroupName != awssdk.StringValue(oldStatus.DBSubnetGroup.DBSubnetGroupName) {
			return nil, cloudobject.SpecInvalidError{Message: "modifying DBSubnetGroupName is not possible for MultiAZ instances"}
		}
	}

	if assertedSpec.Storage.StorageEncrypted != awssdk.BoolValue(oldStatus.StorageEncrypted) {
		return nil, cloudobject.SpecInvalidError{Message: "modifying Storage Encryption is not possible"}
	}


	input := assertedSpec.ModifyDBInstanceInput(i.ID().String())
	if _, err := i.session.ModifyDBInstance(&input); err != nil {
		return nil, err
	}

	// TODO: create instance secrets

	return nil, nil
}

// Delete deletes an Instance.
func (i *Instance) Delete(purge bool) error {
	exists, err := i.Exists()
	if err != nil {
		return err
	}
	if !exists {
		return cloudobject.NotExistsError{Message: fmt.Sprintf("cannot delete non-existing RDS instance '%s'",
			i.ID().String())}
	}

	// If status is already in 'deleting' then we can stop here
	if *i.status.DBInstanceStatus == "deleting" {
		return nil
	}

	// If status is not available, we shouldn't continue... we should only delete instances that are ready
	if *i.status.DBInstanceStatus != "available" {
		return cloudobject.NotReadyError{Message: fmt.Sprintf("cannot delete not-available RDS instance '%s'",
			i.ID().String())}
	}

	snapExists, err := snapshotExists(i)
	if err != nil {
		return err
	}

	// If snapshot already exists we need to throw an error, as we won't be able to backup
	if snapExists && !purge {
		_, err := i.session.DeleteDBSnapshot(&awsrds.DeleteDBSnapshotInput{
			DBSnapshotIdentifier: awssdk.String(finalDBSnapshotName(i)),
		})
		if err != nil {
			return err
		}
	}

	var skipfinalsnapshot bool
	var deletebackups bool
	// on purge, we don't want to store a final snapshot
	if purge {
		skipfinalsnapshot = true
		deletebackups = true
	}

	input := awsrds.DeleteDBInstanceInput{
		DBInstanceIdentifier:      i.ID().StringPtr(),
		DeleteAutomatedBackups:    awssdk.Bool(deletebackups),
		FinalDBSnapshotIdentifier: awssdk.String(finalDBSnapshotName(i)),
		SkipFinalSnapshot:         awssdk.Bool(skipfinalsnapshot),
	}
	// Let's do this... Let's actually delete the DB instance
	if _, err := i.session.ModifyDBInstance(&awsrds.ModifyDBInstanceInput{
		DeletionProtection:   awssdk.Bool(false),
		DBInstanceIdentifier: i.ID().StringPtr(),
	}); err != nil {
		return err
	}
	if _, err := i.session.DeleteDBInstance(&input); err != nil {
		if err.(awserr.Error).Code() != awsrds.ErrCodeDBInstanceNotFoundFault {
			return err
		}
	}

	// If purge we delete our encryption key also
	if purge {
		var key *kms.Key
		key, err = kmsKeySession(i)
		if err != nil {
			return err
		}
		err = key.Delete(purge)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Instance) ARN() *awsarn.ARN {
	if err := i.Read(); err != nil {
		return nil
	}
	arn := aws.MustParse(*i.status.DBInstanceArn)
	return &arn
}

func (i *Instance) Exists() (bool, error) {
	return cloudobject.Exists(i)
}

func (i *Instance) Status() cloudobject.Status {
	return i.status
}

////////////
/// SPEC ///
////////////

type InstanceSpec struct {
	// A value that indicates whether minor engine upgrades are applied automatically
	// to the DB instance during the maintenance window. By default, minor engine
	// upgrades are applied automatically.
	AutoMinorVersionUpgrade bool

	// The Availability Zone (AZ) where the database will be created. For information
	// on AWS Regions and Availability Zones, see Regions and Availability Zones
	// (https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.RegionsAndAvailabilityZones.html).
	//
	// Default: A random, system-chosen Availability Zone in the endpoint's AWS
	// Region.
	//
	// Example: us-east-1d
	//
	// Constraint: The AvailabilityZone parameter can't be specified if the DB instance
	// is a Multi-AZ deployment. The specified Availability Zone must be in the
	// same AWS Region as the current endpoint.
	//
	// If you're creating a DB instance in an RDS on VMware environment, specify
	// the identifier of the custom Availability Zone to create the DB instance
	// in.
	//
	// For more information about RDS on VMware, see the RDS on VMware User Guide.
	// (https://docs.aws.amazon.com/AmazonRDS/latest/RDSonVMwareUserGuide/rds-on-vmware.html)
	AvailabilityZone string

	// The number of days for which automated backups are retained. Setting this
	// parameter to a positive number enables backups. Setting this parameter to
	// 0 disables automated backups.
	//
	// Default: 1
	//
	// Constraints:
	//
	//    * Must be a value from 0 to 35
	//
	//    * Can't be set to 0 if the DB instance is a source to read replicas
	BackupRetentionPeriod int64

	// The compute and memory capacity of the DB instance, for example, db.m4.large.
	// Not all DB instance classes are available in all AWS Regions, or for all
	// database engines. For the full list of DB instance classes, and availability
	// for your engine, see DB Instance Class (https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.DBInstanceClass.html)
	// in the Amazon RDS User Guide.
	//
	// DBInstanceClass is a required field
	DBInstanceClass string

	// The initialized Database name. Default: 'default'
	DBName string

	// A DB subnet group to associate with this DB instance.
	//
	// If there is no DB subnet group, then it is a non-VPC DB instance.
	DBSubnetGroupName string

	// The name of the database engine to be used for this instance.
	Engine InstanceDBEngine

	// The version number of the database engine to use.
	EngineVersion string

	// The password for the master user. The password can include any printable
	// ASCII character except "/", """, or "@".
	MasterUserPassword string

	// The name for the master user.
	MasterUsername string

	// Monitoring defines a separete Monitoring role setup
	Monitoring *InstanceMonitoringSpec

	// Defines PerformanceInsights config if set
	PerformanceInsights *InstancePerformanceInsightsSpec

	// The port number on which the database accepts connections.
	Port int64

	// The daily time range during which automated backups are created if automated
	// backups are enabled, using the BackupRetentionPeriod parameter. For more
	// information, see The Backup Window (https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_WorkingWithAutomatedBackups.html#USER_WorkingWithAutomatedBackups.BackupWindow)
	// in the Amazon RDS User Guide.
	//
	// The default is a 30-minute window selected at random from an 8-hour block
	// of time for each AWS Region. To see the time blocks available, see Adjusting
	// the Preferred DB Instance Maintenance Window (https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_UpgradeDBInstance.Maintenance.html#AdjustingTheMaintenanceWindow)
	// in the Amazon RDS User Guide.
	//
	// Constraints:
	//
	//    * Must be in the format hh24:mi-hh24:mi.
	//
	//    * Must be in Universal Coordinated Time (UTC).
	//
	//    * Must not conflict with the preferred maintenance window.
	//
	//    * Must be at least 30 minutes.
	PreferredBackupWindow string

	// The time range each week during which system maintenance can occur, in Universal
	// Coordinated Time (UTC). For more information, see Amazon RDS Maintenance
	// Window (https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_UpgradeDBInstance.Maintenance.html#Concepts.DBMaintenance).
	//
	// Format: ddd:hh24:mi-ddd:hh24:mi
	//
	// The default is a 30-minute window selected at random from an 8-hour block
	// of time for each AWS Region, occurring on a random day of the week.
	//
	// Valid Days: Mon, Tue, Wed, Thu, Fri, Sat, Sun.
	//
	// Constraints: Minimum 30-minute window.
	PreferredMaintenanceWindow string

	// Defines whether the DB will have a public endpoint attached
	PubliclyAccessible bool

	// If true: throws an error when restoration procedure is auto-detected. (Key & Snapshot detected)
	RestorationDisabled bool

	Storage InstanceStorageSpec

	// Tags to assign to the DB instance.
	Tags map[string]string

	// A list of Amazon EC2 VPC security groups to associate with this DB instance.
	VpcSecurityGroupIds []string
}

type InstanceMonitoringSpec struct {

	// The interval, in seconds, between points when Enhanced Monitoring metrics
	// are collected for the DB instance. To disable collecting Enhanced Monitoring
	// metrics, specify 0. The default is 0.
	//
	// If MonitoringRoleArn is specified, then you must also set MonitoringInterval
	// to a value other than 0.
	//
	// Valid Values: 0, 1, 5, 10, 15, 30, 60
	MonitoringInterval int64

	// The ARN for the IAM role that permits RDS to send enhanced monitoring metrics
	// to Amazon CloudWatch Logs. For example, arn:aws:iam:123456789012:role/emaccess.
	// For information on creating a monitoring role, go to Setting Up and Enabling
	// Enhanced Monitoring (https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_Monitoring.OS.html#USER_Monitoring.OS.Enabling)
	// in the Amazon RDS User Guide.
	//
	// If MonitoringInterval is set to a value other than 0, then you must supply
	// a MonitoringRoleArn value.
	MonitoringRoleArn awsarn.ARN
}

type InstancePerformanceInsightsSpec struct {

	// The amount of time, in days, to retain Performance Insights data. Valid values
	// are 7 or 731 (2 years).
	PerformanceInsightsRetentionPeriod int64
}

type InstanceStorageSpec struct {

	// The amount of storage (in gibibytes) to allocate for the DB instance.
	//
	// MySQL
	//
	// Constraints to the amount of storage for each storage type are the following:
	//
	//    * General Purpose (SSD) storage (gp2): Must be an integer from 20 to 65536.
	//
	//    * Provisioned IOPS storage (io1): Must be an integer from 100 to 65536.
	//
	//    * Magnetic storage (standard): Must be an integer from 5 to 3072.
	//
	// MariaDB
	//
	// Constraints to the amount of storage for each storage type are the following:
	//
	//    * General Purpose (SSD) storage (gp2): Must be an integer from 20 to 65536.
	//
	//    * Provisioned IOPS storage (io1): Must be an integer from 100 to 65536.
	//
	//    * Magnetic storage (standard): Must be an integer from 5 to 3072.
	//
	// PostgreSQL
	//
	// Constraints to the amount of storage for each storage type are the following:
	//
	//    * General Purpose (SSD) storage (gp2): Must be an integer from 20 to 65536.
	//
	//    * Provisioned IOPS storage (io1): Must be an integer from 100 to 65536.
	//
	//    * Magnetic storage (standard): Must be an integer from 5 to 3072.
	AllocatedStorage int64

	// The amount of Provisioned IOPS (input/output operations per second) to be
	// initially allocated for the DB instance. For information about valid Iops
	// values, see Amazon RDS Provisioned IOPS Storage to Improve Performance (https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_Storage.html#USER_PIOPS)
	// in the Amazon RDS User Guide.
	//
	// Constraints: For MariaDB, MySQL, Oracle, and PostgreSQL DB instances, must
	// be a multiple between .5 and 50 of the storage amount for the DB instance.
	Iops int64

	// The upper limit to which Amazon RDS can automatically scale the storage of
	// the DB instance.
	MaxAllocatedStorage int64

	// A value that indicates whether the DB instance is encrypted.
	StorageEncrypted bool

	// Specifies the storage type to be associated with the DB instance.
	//
	// Valid values: standard | gp2 | io1
	StorageType InstanceStorageType
}

func (spec *InstanceSpec) Valid() (bool, error) {
	if spec.DBName == "" {
		return false, cloudobject.SpecInvalidError{Message: "DBName in spec is empty"}
	}

	if len(spec.MasterUserPassword) > 41 || len(spec.MasterUserPassword) < 1 {
		return false, cloudobject.SpecInvalidError{Message: "MasterUserPassword must be > 1 && <= 41 characters"}
	}

	if len(spec.MasterUsername) > 16 || len(spec.MasterUsername) < 1 {
		return false, cloudobject.SpecInvalidError{Message: "MasterUsername must be > 1 && <= 16 characters"}
	}

	if spec.Storage.StorageType == "" {
		return false, cloudobject.SpecInvalidError{Message: "StorageType in spec is empty"}
	}

	return true, nil
}

const (
	MySQLInstanceDBEngine      InstanceDBEngine = "mysql"
	PostgreSQLInstanceDBEngine InstanceDBEngine = "postgres"
	MariaDBInstanceDBEngine    InstanceDBEngine = "mariadb"
)

type InstanceDBEngine string

const (
	StandardInstanceStorageType InstanceStorageType = "standard"
	IO1InstanceStorageType      InstanceStorageType = "io1"
	GP2InstanceStorageType      InstanceStorageType = "gp2"
)

func (engine InstanceDBEngine) String() string {
	return string(engine)
}

type InstanceStorageType string

func (storageType InstanceStorageType) String() string {
	return string(storageType)
}

///////////////
/// HELPERS ///
///////////////

type InstanceStatus awsrds.DBInstance

func (status *InstanceStatus) String() string {
	return awsrds.DBInstance(*status).String()
}

func (status *InstanceStatus) ProviderID() cloudobject.ProviderID {
	return cloudobject.ProviderID{
		Type: cloudobject.AWSProvider,
		Value: *status.DBInstanceArn,
	}
}

type InstanceSecrets struct {
}

func (secrets InstanceSecrets) Map() map[string]string {
	out := make(map[string]string)
	// TODO: Fill map
	return out
}

func finalDBSnapshotName(i *Instance) string {
	return aws.CloudObjectResource(PreDeleteDBSnapshotTopic, i.name)
}

// Use to see if pre-delete snapshot exists
func snapshotExists(i *Instance) (bool, error) {
	out, err := i.session.DescribeDBSnapshots(&awsrds.DescribeDBSnapshotsInput{
		DBSnapshotIdentifier: awssdk.String(finalDBSnapshotName(i)),
		IncludePublic:        awssdk.Bool(false),
		IncludeShared:        awssdk.Bool(false),
	})
	if err != nil {
		if err.(awserr.Error).Code() == awsrds.ErrCodeDBSnapshotNotFoundFault {
			return false, nil
		}
		return false, err
	}
	// If our output DB list is greater than 1, we have an issue with our backup detector
	if len(out.DBSnapshots) < 1 {
		return false, cloudobject.AmbiguousIdentifierError{Message: fmt.Sprintf(
			"multiple RDS DB Instance with id '%s' found", i.ID().String())}
	}
	// If our output DB list is 0, we didn't find our snapshot
	if len(out.DBSnapshots) == 0 {
		return false, nil
	}

	return true, nil
}

// Use to see if DB key already exists
func keyExists(i *Instance, assertedSpec *InstanceSpec, kmsSession *session.Session) (bool, error) {
	key, err := kms.NewKey(i.name, kmsSession)
	if err != nil {
		return false, err
	}
	return key.Exists()
}

///////////////
/// AWS API ///
///////////////

// CreateDBInstanceInput returns the marshaled AWS Interface object of same name
func (spec *InstanceSpec) CreateDBInstanceInput(id string) awsrds.CreateDBInstanceInput {
	dbname := getDBName(spec)

	tags := compileTags(spec.Tags)

	out := awsrds.CreateDBInstanceInput{
		AutoMinorVersionUpgrade: awssdk.Bool(spec.AutoMinorVersionUpgrade),
		BackupRetentionPeriod:   awssdk.Int64(spec.BackupRetentionPeriod),
		CopyTagsToSnapshot:      awssdk.Bool(true),
		DBInstanceClass:         awssdk.String(spec.DBInstanceClass),
		DBInstanceIdentifier:    awssdk.String(id),
		DBName:                  awssdk.String(dbname),
		DBSubnetGroupName:       awssdk.String(spec.DBSubnetGroupName),
		DeletionProtection:      awssdk.Bool(true),
		Engine:                  awssdk.String(spec.Engine.String()),
		EngineVersion:           awssdk.String(spec.EngineVersion),
		// KmsKeyId we'll set on creation... there we have the key creation/discovery logic
		MasterUserPassword:         awssdk.String(spec.MasterUserPassword),
		MasterUsername:             awssdk.String(spec.MasterUsername),
		Port:                       awssdk.Int64(spec.Port),
		PreferredBackupWindow:      awssdk.String(spec.PreferredBackupWindow),
		PreferredMaintenanceWindow: awssdk.String(spec.PreferredMaintenanceWindow),
		PubliclyAccessible:         awssdk.Bool(spec.PubliclyAccessible),
		Tags:                       tags,
		VpcSecurityGroupIds:        awssdk.StringSlice(spec.VpcSecurityGroupIds),
	}

	out.StorageType = awssdk.String(spec.Storage.StorageType.String())
	out.AllocatedStorage = awssdk.Int64(spec.Storage.AllocatedStorage)
	out.MaxAllocatedStorage = awssdk.Int64(spec.Storage.MaxAllocatedStorage)
	out.StorageEncrypted = awssdk.Bool(spec.Storage.StorageEncrypted)
	if spec.Storage.StorageType == IO1InstanceStorageType {
		out.Iops = awssdk.Int64(spec.Storage.Iops)
	}

	if spec.Monitoring != nil {
		out.MonitoringInterval = awssdk.Int64(spec.Monitoring.MonitoringInterval)
		out.MonitoringRoleArn = awssdk.String(spec.Monitoring.MonitoringRoleArn.String())
	}

	if spec.PerformanceInsights != nil {
		out.EnablePerformanceInsights = awssdk.Bool(true)
		out.PerformanceInsightsRetentionPeriod = awssdk.Int64(spec.PerformanceInsights.PerformanceInsightsRetentionPeriod)
	}

	if spec.DBSubnetGroupName != "" {
		out.DBSubnetGroupName = awssdk.String(spec.DBSubnetGroupName)
	}

	if spec.AvailabilityZone != "" {
		out.AvailabilityZone = awssdk.String(spec.AvailabilityZone)
	} else {
		out.MultiAZ = awssdk.Bool(true)
	}
	return out
}

// RestoreDBInstanceFromDBSnapshotInput returns the marshaled AWS Interface object of same name
func (spec *InstanceSpec) RestoreDBInstanceFromDBSnapshotInput(id string, snapshotId string) awsrds.
	RestoreDBInstanceFromDBSnapshotInput {

	tags := compileTags(spec.Tags)

	out := awsrds.RestoreDBInstanceFromDBSnapshotInput{
		AutoMinorVersionUpgrade: awssdk.Bool(spec.AutoMinorVersionUpgrade),
		CopyTagsToSnapshot:      awssdk.Bool(true),
		DBInstanceClass:         awssdk.String(spec.DBInstanceClass),
		DBInstanceIdentifier:    awssdk.String(id),
		DBSnapshotIdentifier:    awssdk.String(snapshotId),
		DBSubnetGroupName:       awssdk.String(spec.DBSubnetGroupName),
		DeletionProtection:      awssdk.Bool(true),
		Engine:                  awssdk.String(spec.Engine.String()),
		Port:                    awssdk.Int64(spec.Port),
		PubliclyAccessible:      awssdk.Bool(spec.PubliclyAccessible),
		Tags:                    tags,
		VpcSecurityGroupIds:     awssdk.StringSlice(spec.VpcSecurityGroupIds),
	}

	out.StorageType = awssdk.String(spec.Storage.StorageType.String())
	if spec.Storage.StorageType == IO1InstanceStorageType {
		out.Iops = awssdk.Int64(spec.Storage.Iops)
	}

	if spec.AvailabilityZone != "" {
		out.AvailabilityZone = awssdk.String(spec.AvailabilityZone)
	} else {
		out.MultiAZ = awssdk.Bool(true)
	}

	return out
}

// ModifyDBInstanceInput returns the marshaled AWS Interface object of same name
func (spec *InstanceSpec) ModifyDBInstanceInput(id string) awsrds.
	ModifyDBInstanceInput {
	out := awsrds.ModifyDBInstanceInput{
		ApplyImmediately:           awssdk.Bool(true),
		AutoMinorVersionUpgrade:    awssdk.Bool(spec.AutoMinorVersionUpgrade),
		BackupRetentionPeriod:      awssdk.Int64(spec.BackupRetentionPeriod),
		CopyTagsToSnapshot:         awssdk.Bool(true),
		DBInstanceClass:            awssdk.String(spec.DBInstanceClass),
		DBInstanceIdentifier:       awssdk.String(id),
		DeletionProtection:         awssdk.Bool(true),
		EngineVersion:              awssdk.String(spec.EngineVersion),
		MasterUserPassword:         awssdk.String(spec.MasterUserPassword),
		PreferredBackupWindow:      awssdk.String(spec.PreferredBackupWindow),
		PreferredMaintenanceWindow: awssdk.String(spec.PreferredMaintenanceWindow),
		PubliclyAccessible:         awssdk.Bool(spec.PubliclyAccessible),
	}

	out.StorageType = awssdk.String(spec.Storage.StorageType.String())
	out.AllocatedStorage = awssdk.Int64(spec.Storage.AllocatedStorage)
	out.MaxAllocatedStorage = awssdk.Int64(spec.Storage.MaxAllocatedStorage)
	if spec.Storage.StorageType == IO1InstanceStorageType {
		out.Iops = awssdk.Int64(spec.Storage.Iops)
	}

	if spec.Monitoring != nil {
		out.MonitoringInterval = awssdk.Int64(spec.Monitoring.MonitoringInterval)
		out.MonitoringRoleArn = awssdk.String(spec.Monitoring.MonitoringRoleArn.String())
	}

	if spec.PerformanceInsights != nil {
		out.EnablePerformanceInsights = awssdk.Bool(true)
		out.PerformanceInsightsRetentionPeriod = awssdk.Int64(spec.PerformanceInsights.PerformanceInsightsRetentionPeriod)
	}

	if spec.PerformanceInsights != nil {
		out.EnablePerformanceInsights = awssdk.Bool(true)
		out.PerformanceInsightsRetentionPeriod = awssdk.Int64(spec.PerformanceInsights.PerformanceInsightsRetentionPeriod)
	}

	if spec.AvailabilityZone == "" {
		out.MultiAZ = awssdk.Bool(true)
	} else {
		out.MultiAZ = awssdk.Bool(false)
	}

	if len(spec.VpcSecurityGroupIds) != 0 {
		out.VpcSecurityGroupIds = awssdk.StringSlice(spec.VpcSecurityGroupIds)
	}

	return out
}

func getDBName(spec *InstanceSpec) string {
	dbname := spec.DBName
	if dbname == "" {
		dbname = "default"
	}
	return dbname
}

//////////////
/// ERRORS ///
//////////////

type RestorationDisabledError struct {
	Message string
}

func (e RestorationDisabledError) Error() string {
	return e.Message
}
