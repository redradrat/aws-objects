package s3

import (
	"github.com/aws/aws-sdk-go/aws/client"
	awss3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/redradrat/cloud-objects/aws/kms"
	"github.com/redradrat/cloud-objects/cloudobject"
	"reflect"
	"testing"
)

func TestBucketSecrets_Map(t *testing.T) {
	tests := []struct {
		name string
		want map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secrets := BucketSecrets{}
			if got := secrets.Map(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucketSpec_CreateBucketInput(t *testing.T) {
	type fields struct {
		Location              string
		ACL                   string
		ObjectLock            bool
		Versioning            bool
		TransferAcceleration  bool
		BlockPublicAcls       bool
		IgnorePublicAcls      bool
		BlockPublicPolicy     bool
		RestrictPublicBuckets bool
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   awss3.CreateBucketInput
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := BucketSpec{
				Location:              tt.fields.Location,
				ACL:                   tt.fields.ACL,
				ObjectLock:            tt.fields.ObjectLock,
				Versioning:            tt.fields.Versioning,
				TransferAcceleration:  tt.fields.TransferAcceleration,
				BlockPublicAcls:       tt.fields.BlockPublicAcls,
				IgnorePublicAcls:      tt.fields.IgnorePublicAcls,
				BlockPublicPolicy:     tt.fields.BlockPublicPolicy,
				RestrictPublicBuckets: tt.fields.RestrictPublicBuckets,
			}
			if got := b.CreateBucketInput(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateBucketInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucketSpec_PutBucketAccelerationInput(t *testing.T) {
	type fields struct {
		Location              string
		ACL                   string
		ObjectLock            bool
		Versioning            bool
		TransferAcceleration  bool
		BlockPublicAcls       bool
		IgnorePublicAcls      bool
		BlockPublicPolicy     bool
		RestrictPublicBuckets bool
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   awss3.PutBucketAccelerateConfigurationInput
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := BucketSpec{
				Location:              tt.fields.Location,
				ACL:                   tt.fields.ACL,
				ObjectLock:            tt.fields.ObjectLock,
				Versioning:            tt.fields.Versioning,
				TransferAcceleration:  tt.fields.TransferAcceleration,
				BlockPublicAcls:       tt.fields.BlockPublicAcls,
				IgnorePublicAcls:      tt.fields.IgnorePublicAcls,
				BlockPublicPolicy:     tt.fields.BlockPublicPolicy,
				RestrictPublicBuckets: tt.fields.RestrictPublicBuckets,
			}
			if got := b.PutBucketAccelerationInput(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PutBucketAccelerationInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucketSpec_PutBucketAclInput(t *testing.T) {
	type fields struct {
		Location              string
		ACL                   string
		ObjectLock            bool
		Versioning            bool
		TransferAcceleration  bool
		BlockPublicAcls       bool
		IgnorePublicAcls      bool
		BlockPublicPolicy     bool
		RestrictPublicBuckets bool
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   awss3.PutBucketAclInput
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := BucketSpec{
				Location:              tt.fields.Location,
				ACL:                   tt.fields.ACL,
				ObjectLock:            tt.fields.ObjectLock,
				Versioning:            tt.fields.Versioning,
				TransferAcceleration:  tt.fields.TransferAcceleration,
				BlockPublicAcls:       tt.fields.BlockPublicAcls,
				IgnorePublicAcls:      tt.fields.IgnorePublicAcls,
				BlockPublicPolicy:     tt.fields.BlockPublicPolicy,
				RestrictPublicBuckets: tt.fields.RestrictPublicBuckets,
			}
			if got := b.PutBucketAclInput(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PutBucketAclInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucketSpec_PutBucketEncryptionInput(t *testing.T) {
	type fields struct {
		Location              string
		ACL                   string
		ObjectLock            bool
		Versioning            bool
		TransferAcceleration  bool
		BlockPublicAcls       bool
		IgnorePublicAcls      bool
		BlockPublicPolicy     bool
		RestrictPublicBuckets bool
	}
	type args struct {
		id  string
		key *kms.Key
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   awss3.PutBucketEncryptionInput
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := BucketSpec{
				Location:              tt.fields.Location,
				ACL:                   tt.fields.ACL,
				ObjectLock:            tt.fields.ObjectLock,
				Versioning:            tt.fields.Versioning,
				TransferAcceleration:  tt.fields.TransferAcceleration,
				BlockPublicAcls:       tt.fields.BlockPublicAcls,
				IgnorePublicAcls:      tt.fields.IgnorePublicAcls,
				BlockPublicPolicy:     tt.fields.BlockPublicPolicy,
				RestrictPublicBuckets: tt.fields.RestrictPublicBuckets,
			}
			if got := b.PutBucketEncryptionInput(tt.args.id, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PutBucketEncryptionInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucketSpec_PutBucketVersioningInput(t *testing.T) {
	type fields struct {
		Location              string
		ACL                   string
		ObjectLock            bool
		Versioning            bool
		TransferAcceleration  bool
		BlockPublicAcls       bool
		IgnorePublicAcls      bool
		BlockPublicPolicy     bool
		RestrictPublicBuckets bool
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   awss3.PutBucketVersioningInput
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := BucketSpec{
				Location:              tt.fields.Location,
				ACL:                   tt.fields.ACL,
				ObjectLock:            tt.fields.ObjectLock,
				Versioning:            tt.fields.Versioning,
				TransferAcceleration:  tt.fields.TransferAcceleration,
				BlockPublicAcls:       tt.fields.BlockPublicAcls,
				IgnorePublicAcls:      tt.fields.IgnorePublicAcls,
				BlockPublicPolicy:     tt.fields.BlockPublicPolicy,
				RestrictPublicBuckets: tt.fields.RestrictPublicBuckets,
			}
			if got := b.PutBucketVersioningInput(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PutBucketVersioningInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucketSpec_PutPublicAccessBlockInput(t *testing.T) {
	type fields struct {
		Location              string
		ACL                   string
		ObjectLock            bool
		Versioning            bool
		TransferAcceleration  bool
		BlockPublicAcls       bool
		IgnorePublicAcls      bool
		BlockPublicPolicy     bool
		RestrictPublicBuckets bool
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   awss3.PutPublicAccessBlockInput
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := BucketSpec{
				Location:              tt.fields.Location,
				ACL:                   tt.fields.ACL,
				ObjectLock:            tt.fields.ObjectLock,
				Versioning:            tt.fields.Versioning,
				TransferAcceleration:  tt.fields.TransferAcceleration,
				BlockPublicAcls:       tt.fields.BlockPublicAcls,
				IgnorePublicAcls:      tt.fields.IgnorePublicAcls,
				BlockPublicPolicy:     tt.fields.BlockPublicPolicy,
				RestrictPublicBuckets: tt.fields.RestrictPublicBuckets,
			}
			if got := b.PutPublicAccessBlockInput(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PutPublicAccessBlockInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucketSpec_Valid(t *testing.T) {
	type fields struct {
		Location              string
		ACL                   string
		ObjectLock            bool
		Versioning            bool
		TransferAcceleration  bool
		BlockPublicAcls       bool
		IgnorePublicAcls      bool
		BlockPublicPolicy     bool
		RestrictPublicBuckets bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := BucketSpec{
				Location:              tt.fields.Location,
				ACL:                   tt.fields.ACL,
				ObjectLock:            tt.fields.ObjectLock,
				Versioning:            tt.fields.Versioning,
				TransferAcceleration:  tt.fields.TransferAcceleration,
				BlockPublicAcls:       tt.fields.BlockPublicAcls,
				IgnorePublicAcls:      tt.fields.IgnorePublicAcls,
				BlockPublicPolicy:     tt.fields.BlockPublicPolicy,
				RestrictPublicBuckets: tt.fields.RestrictPublicBuckets,
			}
			got, err := b.Valid()
			if (err != nil) != tt.wantErr {
				t.Errorf("Valid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Valid() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucketStatus_ProviderID(t *testing.T) {
	type fields struct {
		Bucket    awss3.Bucket
		Encrypted bool
		ARN       string
	}
	tests := []struct {
		name   string
		fields fields
		want   cloudobject.ProviderID
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := BucketStatus{
				Bucket:    tt.fields.Bucket,
				Encrypted: tt.fields.Encrypted,
				ARN:       tt.fields.ARN,
			}
			if got := status.ProviderID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProviderID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucketStatus_String(t *testing.T) {
	type fields struct {
		Bucket    awss3.Bucket
		Encrypted bool
		ARN       string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := BucketStatus{
				Bucket:    tt.fields.Bucket,
				Encrypted: tt.fields.Encrypted,
				ARN:       tt.fields.ARN,
			}
			if got := status.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucket_Create(t *testing.T) {
	type fields struct {
		name    string
		status  BucketStatus
		session *awss3.S3
	}
	type args struct {
		spec cloudobject.CloudObjectSpec
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    cloudobject.Secrets
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bucket{
				name:    tt.fields.name,
				status:  tt.fields.status,
				session: tt.fields.session,
			}
			got, err := b.Create(tt.args.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucket_Delete(t *testing.T) {
	type fields struct {
		name    string
		status  BucketStatus
		session *awss3.S3
	}
	type args struct {
		in0 bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bucket{
				name:    tt.fields.name,
				status:  tt.fields.status,
				session: tt.fields.session,
			}
			if err := b.Delete(tt.args.in0); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBucket_Exists(t *testing.T) {
	type fields struct {
		name    string
		status  BucketStatus
		session *awss3.S3
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bucket{
				name:    tt.fields.name,
				status:  tt.fields.status,
				session: tt.fields.session,
			}
			got, err := b.Exists()
			if (err != nil) != tt.wantErr {
				t.Errorf("Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Exists() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucket_ID(t *testing.T) {
	type fields struct {
		name    string
		status  BucketStatus
		session *awss3.S3
	}
	tests := []struct {
		name   string
		fields fields
		want   cloudobject.ID
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bucket{
				name:    tt.fields.name,
				status:  tt.fields.status,
				session: tt.fields.session,
			}
			if got := b.ID(); got != tt.want {
				t.Errorf("ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucket_Read(t *testing.T) {
	type fields struct {
		name    string
		status  BucketStatus
		session *awss3.S3
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bucket{
				name:    tt.fields.name,
				status:  tt.fields.status,
				session: tt.fields.session,
			}
			if err := b.Read(); (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBucket_Status(t *testing.T) {
	type fields struct {
		name    string
		status  BucketStatus
		session *awss3.S3
	}
	tests := []struct {
		name   string
		fields fields
		want   cloudobject.Status
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bucket{
				name:    tt.fields.name,
				status:  tt.fields.status,
				session: tt.fields.session,
			}
			if got := b.Status(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Status() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucket_Update(t *testing.T) {
	type fields struct {
		name    string
		status  BucketStatus
		session *awss3.S3
	}
	type args struct {
		spec cloudobject.CloudObjectSpec
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    cloudobject.Secrets
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bucket{
				name:    tt.fields.name,
				status:  tt.fields.status,
				session: tt.fields.session,
			}
			got, err := b.Update(tt.args.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBucket(t *testing.T) {
	type args struct {
		name    string
		session client.ConfigProvider
	}
	tests := []struct {
		name    string
		args    args
		want    *Bucket
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBucket(tt.args.name, tt.args.session)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBucket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBucket() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSaneS3Bucket(t *testing.T) {
	tests := []struct {
		name string
		want BucketSpec
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SaneS3Bucket(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SaneS3Bucket() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ensureBucketConfig(t *testing.T) {
	type args struct {
		assertedSpec *BucketSpec
		b            *Bucket
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ensureBucketConfig(tt.args.assertedSpec, tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("ensureBucketConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_kmsKeySession(t *testing.T) {
	type args struct {
		b *Bucket
	}
	tests := []struct {
		name    string
		args    args
		want    *kms.Key
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := kmsKeySession(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("kmsKeySession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("kmsKeySession() got = %v, want %v", got, tt.want)
			}
		})
	}
}
