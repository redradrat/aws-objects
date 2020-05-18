package kms

func EncryptSymmetric() KeySpec {
	return KeySpec{
		KeyUsage: EncryptDecryptKeyUsage,
		KeyType:  SymmetricDefaultKeyType,
	}
}
