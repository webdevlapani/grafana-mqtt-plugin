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
