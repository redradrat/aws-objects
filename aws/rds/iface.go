package rds

import "github.com/aws/aws-sdk-go/service/iam/iamiface"

// RDS defines interactions with the rds cloud objects
type RDS interface {
	Create(svc iamiface.IAMAPI) error
	Update(svc iamiface.IAMAPI) error
	Delete(svc iamiface.IAMAPI) error
	SyncStatus(svc iamiface.IAMAPI)
}
