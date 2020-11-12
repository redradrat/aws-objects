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
