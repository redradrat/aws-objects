# cloud-objects ![Go](https://github.com/redradrat/cloud-objects/workflows/Go/badge.svg?branch=master)
Generic Cloud Service Object C(R)UD Representation

This repository hosts a generic abstraction of Cloud Service Objects into golang CRUD instances.

Currently following Objects are implemented:

**AWS**

| Amazon Web Service | Supported Resources |
| --- | --- |
| IAM | <ul><li>Group</li><li>Policy</li><li>PolicyAttachment</li><li>Role</li><li>User</li></ul> |
| RDS | <ul><li>DB Instance</li><li>DB SubnetGroup</li></ul> |
| KMS | <ul><li>Key</li></ul> |

### RDS

Following types are currently supported:
* DB Instance 
* DB Subnet Group

As this cloud object library attempts to provide "simple" C(R)UD interactions on these 
objects there is opinionated logic attached to RDS instance handling.

**RDS Instance**

On initial instance creation, an KMS Key will be created. If storage encryption has been 
enabled, this key will be used to encrypt. 

On deletion a final DB snapshot will be taken, and unless a user selects the "purge"
option, the key and snapshot will remain.

On creation of an instance with the same name, the library will check for existance of
key and snapshot and will restore the snapshot. If the user opts to disable the 
restoration option, an error will be given on creation. (A user will have to manually
clean the snapshot and maybe even the key. This is equal to selecting "purge" on delete.)

Summary:
* A DB will always store a snapshot on delete
* A DB will always restore if snapshot and encryption key detected
* A DB has to be purged to completely be deleted