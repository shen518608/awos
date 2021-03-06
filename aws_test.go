package awos

/**
AccessKeyId=${accessKeyId} AccessKeySecret=${accessKeySecret} Endpoint=${endpoint} Bucket=${bucket} go test -v aws_test.go aws.go
*/

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"strconv"
	"testing"
)

const (
	AWSGuid         = "test123"
	AWSContent      = "aaaaaa"
	AWSExpectLength = 6
	AWSExpectHead   = 1
)

var (
	awsClient *AWS
)

func init() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:           aws.String("cn-north-1"),
		DisableSSL:       aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(os.Getenv("AccessKeyId"), os.Getenv("AccessKeySecret"), ""),
		Endpoint:         aws.String(os.Getenv("Endpoint")),
		S3ForcePathStyle: aws.Bool(true),
	}))

	service := s3.New(sess)

	awsClient = &AWS{
		BucketName: os.Getenv("Bucket"),
		Client:     service,
	}
}

func TestAWS_Put(t *testing.T) {
	meta := make(map[string]string)
	meta["head"] = strconv.Itoa(AWSExpectHead)
	meta["length"] = strconv.Itoa(AWSExpectLength)

	err := awsClient.Put(AWSGuid, AWSContent, meta)
	if err != nil {
		t.Log("aws put error", err)
		t.Fail()
	}
}

func TestAWS_Head(t *testing.T) {
	attributes := make([]string, 0)
	attributes = append(attributes, "head")
	var res map[string]string
	var err error
	var head int
	var length int

	res, err = awsClient.Head(AWSGuid, attributes)
	if err != nil {
		t.Log("aws head error", err)
		t.Fail()
	}

	head, err = strconv.Atoi(res["head"])
	if err != nil || head != 1 {
		t.Log("aws get head fail, res:", res, "err:", err)
		t.Fail()
	}

	attributes = append(attributes, "length")
	res, err = awsClient.Head(AWSGuid, attributes)
	if err != nil {
		t.Log("aws head error", err)
		t.Fail()
	}

	head, err = strconv.Atoi(res["head"])
	length, err = strconv.Atoi(res["length"])
	if err != nil || head != AWSExpectHead || length != AWSExpectLength {
		t.Log("aws get head fail, res:", res, "err:", err)
		t.Fail()
	}
}

func TestAWS_Get(t *testing.T) {
	res, err := awsClient.Get(AWSGuid)
	if err != nil || res != AWSContent {
		t.Log("aws get AWSContent fail, res:", res, "err:", err)
		t.Fail()
	}
}

func TestAWS_ListObject(t *testing.T) {
	res, err := awsClient.ListObject(AWSGuid, AWSGuid[0:4], "", 10, "")
	if err != nil || len(res) == 0 {
		t.Log("aws list objects fail, res:", res, "err:", err)
		t.Fail()
	}
}

func TestAWS_Del(t *testing.T) {
	err := awsClient.Del(AWSGuid)
	if err != nil {
		t.Log("aws del key fail, err:", err)
		t.Fail()
	}
}

func TestAWS_GetNotExist(t *testing.T) {
	res1, err := awsClient.Get(AWSGuid + "123")
	if res1 != "" || err != nil {
		t.Log("aws get not exist key fail, res:", res1, "err:", err)
		t.Fail()
	}

	attributes := make([]string, 0)
	attributes = append(attributes, "head")
	res2, err := awsClient.Head(AWSGuid+"123", attributes)
	if res2 != nil || err != nil {
		t.Log("aws head not exist key fail, res:", res2, "err:", err, err.Error())
		t.Fail()
	}
}
