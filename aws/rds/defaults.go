package rds

import "github.com/redradrat/cloud-objects/aws"

// Returns a "sane" defaulted InstanceSpec
func SanePostgres(name, subnetGroupName, instanceClass, user, pass string, tags map[string]string,
	securityGroupIds []string) InstanceSpec {
	return InstanceSpec{
		AutoMinorVersionUpgrade:    true,
		AvailabilityZone:           "",
		BackupRetentionPeriod:      14,
		DBInstanceClass:            instanceClass,
		DBName:                     name,
		DBSubnetGroupName:          aws.CloudObjectResource("SG", subnetGroupName),
		Engine:                     PostgreSQLInstanceDBEngine,
		EngineVersion:              "12.2",
		MasterUserPassword:         pass,
		MasterUsername:             user,
		Monitoring:                 nil,
		PerformanceInsights:        nil,
		Port:                       5432,
		PreferredBackupWindow:      "01:00-02:00",
		PreferredMaintenanceWindow: "Sun:02:00-Sun:03:00",
		PubliclyAccessible:         false,
		RestorationDisabled:        false,
		Storage: InstanceStorageSpec{
			AllocatedStorage:    20,
			MaxAllocatedStorage: 30,
			StorageEncrypted:    true,
			StorageType:         GP2InstanceStorageType,
		},
		Tags:                tags,
		VpcSecurityGroupIds: securityGroupIds,
	}
}
