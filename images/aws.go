package images

import (
	"context"
    "fmt"
    "os"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3PutObjectAPI defines the interface for the PutObject function.
// We use this interface to test the function using a mocked service.
type S3PutObjectAPI interface {
	PutObject(ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

func PutFile(
	c context.Context, api S3PutObjectAPI, input *s3.PutObjectInput) (
	*s3.PutObjectOutput, error) {
	return api.PutObject(c, input)
}


func (i *WebImage) UploadImagesToSpaces() (err error) {

    customResolver := aws.EndpointResolverWithOptionsFunc(
        func(service, region string, options ...interface{}) (aws.Endpoint, error) {
            return aws.Endpoint{URL: i.doEndpointUrl}, nil
        })

    cfg, err := config.LoadDefaultConfig(context.TODO(),
        config.WithEndpointResolverWithOptions(customResolver),
        config.WithCredentialsProvider(
            credentials.NewStaticCredentialsProvider(
                i.doAccessKey, i.doSecret, "")),
    )

    if err != nil {
        panic("AWS configuration error, " + err.Error())
    }

    client := s3.NewFromConfig(cfg, func(o *s3.Options) {
        o.Region = i.doRegionName
    })

    for localPath, uploadPath := range(i.newSizes){


        file, err := os.Open(localPath)
        if err != nil {
            fmt.Println("Unable to open file ", localPath)
            return err
        }

        defer file.Close()

        input := &s3.PutObjectInput{
            Bucket:       &i.doBucket,
            Key:          &uploadPath,
            Body:         file,
            CacheControl: &i.doCacheControl,
            ContentType:  &i.doContentType,
            ACL:          "public-read",
        }

        _, err = PutFile(context.TODO(), client, input)
        if err != nil {
            fmt.Println("Error uploading the file")
            fmt.Println(err)
            return err
        }
    }

    return err
}
