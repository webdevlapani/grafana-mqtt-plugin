package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
	//"github.com/aws/aws-sdk-go/service/s3"
	//"github.com/aws/aws-sdk-go/service/iam"
)

func (ds *MQTTDatasource) authenticate(ctx context.Context, region string) (*iot.IoT, error) {
	config := aws.Config{
		Region: aws.String(region),
	}

	if ds.cfg.AWSAccessKey != "" && ds.cfg.AWSSecretKey != "" {
		config.Credentials = credentials.NewStaticCredentials(ds.cfg.AWSAccessKey, ds.cfg.AWSSecretKey, "")
	}

	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	return iot.New(sess, &config), nil
}

func throw(writer http.ResponseWriter, status int, msg string, err string) {
	response := make(map[string]string)

	if msg == "" {
		msg = "Something went wrong!"
	}
	if err == "" {
		err = msg
	}

	response["message"] = msg
	response["error"] = err

	writer.WriteHeader(status)

	responseBytes, errp := json.Marshal(response)
	if errp != nil {
		log.DefaultLogger.Error("HandleRequest", "error", errp)
		return
	}
	_, _ = writer.Write(responseBytes)
}

func (ds *MQTTDatasource) createStorage(ctx context.Context, region string) error {
	if _, ok := ds.cfg.Storage[region]; ok {
		return nil
	}

	config := aws.Config{
		Region: aws.String(region),
	}

	if ds.cfg.AWSAccessKey != "" && ds.cfg.AWSSecretKey != "" {
		config.Credentials = credentials.NewStaticCredentials(ds.cfg.AWSAccessKey, ds.cfg.AWSSecretKey, "")
	}

	//sess, err := session.NewSession()
	//if err != nil {
	//	return err
	//}

	//iotSvc := iot.New(sess, &config)
	//s3Svc := s3.New(sess, &config)
	//iamSvc := iam.New(sess, &config)

	// 1. create s3 bucket <ds.cfg.S3BucketPrefix>.<region> exists using s3Svc (https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#S3.CreateBucketWithContext)
	// 	1.1. if exists then skip
	// 2. check if role S3.<ds.cfg.S3BucketPrefix>.<region> exists using iotSvc (https://docs.aws.amazon.com/sdk-for-go/api/service/iot/#IoT.GetTopicRuleWithContext)
	//	2.1. if not
	//		2.1.1. get arn of IAM role with name S3.<ds.cfg.S3BucketPrefix>.Role using iamSvc (https://docs.aws.amazon.com/sdk-for-go/api/service/iam/#IAM.GetRoleWithContext)
	//		2.1.2. create iot rule with name S3.<ds.cfg.S3BucketPrefix>.<region> using iotSvc (https://docs.aws.amazon.com/sdk-for-go/api/service/iot/#IoT.CreateTopicRuleWithContext)
	//			Tags: {
	//				Service: ds.cfg.Service
	//				Zone: ds.cfg.Zone
	//				Environment: ds.cfg.Env
	//			}
	//			Payload: {
	//				"awsIotSqlVersion": "2016-03-23",
	//				"sql": "SELECT * FROM '#'",
	//				"ruleDisabled": false,
	//				"actions": [{
	//					"s3": {
	//						"roleArn": <role.arn>,
	//						"bucketName": <ds.cfg.S3BucketPrefix>.<region>,
	//						"key": ds.cfg.Service + "/" + ds.cfg.Env + "/MQTT/${topic()}/${timestamp()}",
	//					}
	//				}]
	//			}


	// Note: if IAM role S3.<ds.cfg.S3BucketPrefix>.<region>.Role does not exist in your account, create a role for IoT and attach the following policy
	//{
	//    "Version": "2012-10-17",
	//    "Statement": {
	//        "Effect": "Allow",
	//        "Action": "s3:PutObject",
	//        "Resource": "arn:aws:s3:::<ds.cfg.S3BucketPrefix>.*"
	//    }
	//}

	// Uncomment this line after implmentation and testing of the above
	//ds.cfg.Storage[region] = 1

	return nil
}
