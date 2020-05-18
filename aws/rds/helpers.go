package rds

import (
	awssdk "github.com/aws/aws-sdk-go/aws"
	awsrds "github.com/aws/aws-sdk-go/service/rds"
)

func compileTags(tagMap map[string]string) []*awsrds.Tag {
	var tags []*awsrds.Tag
	for k, v := range tagMap {
		tags = append(tags, &awsrds.Tag{
			Key:   awssdk.String(k),
			Value: awssdk.String(v),
		})
	}
	return tags
}
