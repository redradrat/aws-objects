package rds

// Returns a "sane" defaulted InstanceSpec
func SanePostgres(name, subnetGroupName, instanceClass, user, pass string, tags map[string]string,
	securityGroupIds []string) InstanceSpec {
	return InstanceSpec{
		AutoMinorVersionUpgrade:    true,
		AvailabilityZone:           "",
		BackupRetentionPeriod:      14,
		DBInstanceClass:            instanceClass,
		DBName:                     name,
		DBSubnetGroupName:          subnetGroupName,
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
		RestorationEnabled:         true,
		Storage:                    InstanceStorageSpec{},
		Tags:                       tags,
		VpcSecurityGroupIds:        securityGroupIds,
	}
}
