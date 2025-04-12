package ancillaries

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func GetS3Client() *s3.Client {
	endpoint := "https://fra1.digitaloceanspaces.com"
	env := GetEnv()
	return s3.NewFromConfig(aws.Config{
		BaseEndpoint: &endpoint,
		Region:       "fra1",
		Credentials: credentials.NewStaticCredentialsProvider(
			env.Aws_Access_Key_Id,
			env.Aws_Secret_Access_Key,
			"",
		),
	})
}

func GetPresignClient() *s3.PresignClient {
	return s3.NewPresignClient(GetS3Client())
}

// BucketBasics encapsulates the Amazon Simple Storage Service (Amazon S3) actions
// used in the examples.
// It contains S3Client, an Amazon S3 service client that is used to perform bucket
// and object actions.
type S3Bucket struct {
	S3Client      *s3.Client
	Name          string
	PresignClient *s3.PresignClient
}

// ListObjects lists the objects in a bucket.
func (b *S3Bucket) ListObjects(ctx context.Context) ([]types.Object, error) {
	var err error
	var output *s3.ListObjectsV2Output
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(b.Name),
	}
	var objects []types.Object
	objectPaginator := s3.NewListObjectsV2Paginator(b.S3Client, input)
	for objectPaginator.HasMorePages() {
		output, err = objectPaginator.NextPage(ctx)
		if err != nil {
			var noBucket *types.NoSuchBucket
			if errors.As(err, &noBucket) {
				log.Printf("Bucket %s does not exist.\n", b.Name)
				err = noBucket
			}
			break
		} else {
			objects = append(objects, output.Contents...)
		}
	}
	return objects, err
}

func (b *S3Bucket) GetObject(key string) ([]byte, error) {
	output, err := b.S3Client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: &b.Name,
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	var data = []byte{}
	_, err = output.Body.Read(data)
	return data, err
}

func NewBucketConn() S3Bucket {
	return S3Bucket{
		Name:          "rashedstudio",
		S3Client:      GetS3Client(),
		PresignClient: GetPresignClient(),
	}
}

func (b *S3Bucket) GetUrl(key string) (string, error) {
	request, err := b.PresignClient.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(b.Name),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(3600 * int64(time.Second))
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to get %v:%v. Here's why: %v\n", b.Name, key, err)
	}
	return request.URL, err
}

var S3 = NewBucketConn()
