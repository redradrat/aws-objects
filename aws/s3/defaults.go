package s3

import awss3 "github.com/aws/aws-sdk-go/service/s3"

// Returns a "sane" defaulted BucketSpec
func SaneS3Bucket() BucketSpec {
	return BucketSpec{
		//Location:             "eu-central-1",
		ACL:                   awss3.BucketCannedACLPrivate,
		ObjectLock:            true,
		Versioning:            true,
		TransferAcceleration:  false,
		BlockPublicAcls:       false,
		IgnorePublicAcls:      true,
		BlockPublicPolicy:     true,
		RestrictPublicBuckets: true,
	}
}
