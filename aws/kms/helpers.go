package kms

import (
	awssdk "github.com/aws/aws-sdk-go/aws"
	awskms "github.com/aws/aws-sdk-go/service/kms"
)

const (
	SignVerifyKeyUsage     KeyUsage = "SIGN_VERIFY"
	EncryptDecryptKeyUsage KeyUsage = "ENCRYPT_DECRYPT"
)

type KeyUsage string

func (usage KeyUsage) String() string {
	return string(usage)
}

const (
	SymmetricDefaultKeyType KeyType = "SYMMETRIC_DEFAULT"
	RSA2048KeyType          KeyType = "RSA_2048"
	RSA3072KeyType          KeyType = "RSA_3072"
	RSA4096KeyType          KeyType = "RSA_4096"
	NISTP256KeyType         KeyType = "ECC_NIST_P256"
	NISTP521KeyType         KeyType = "ECC_NIST_P521"
)

type KeyType string

func (kt KeyType) String() string {
	return string(kt)
}

func compileTags(tagMap map[string]string) []*awskms.Tag {
	var tags []*awskms.Tag
	for k, v := range tagMap {
		tags = append(tags, &awskms.Tag{
			TagKey:   awssdk.String(k),
			TagValue: awssdk.String(v),
		})
	}
	return tags
}
