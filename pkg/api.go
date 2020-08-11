package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	//"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iot"
)

func (ds *MQTTDatasource) handleGetEndpoint(resp http.ResponseWriter, req *http.Request) {
	// TODO: allow only get requests

	region := req.URL.Query().Get("region")
	if region == "" {
		throw(resp, 400, "Invalid or missing region!", "")
		return
	}

	svc, err := ds.authenticate(req.Context(), region)
	if err != nil {
		throw(resp, 500, "Could not create session!", err.Error())
		return
	}

	// Get ATS endpoint of the region (https://docs.aws.amazon.com/sdk-for-go/api/service/iot/#IoT.DescribeEndpointWithContext)
	result, err := svc.DescribeEndpointWithContext(req.Context(), &iot.DescribeEndpointInput{EndpointType: aws.String("iot:Data-ATS")})
	if err != nil {
		throw(resp, 500, "Could not fetch MQTT host!", err.Error())
		return
	}

	_, _ = resp.Write([]byte(*result.EndpointAddress))
}

func (ds *MQTTDatasource) handleCreateCertificate(resp http.ResponseWriter, req *http.Request) {
	// TODO: allow only post requests

	region := req.URL.Query().Get("region")
	if region == "" {
		throw(resp, 400, "Invalid or missing region!", "")
		return
	}

	// TODO: get the following parameters (topic and client) from json body
	// following snippet reads them from url parameters and can be used for debugging
	//topic := req.URL.Query().Get("topic")
	//if topic == "" {
	//	topic = "*"
	//}
	//client := req.URL.Query().Get("client")
	//if client == "" {
	//	client = "*"
	//}

	//svc, err := ds.authenticate(req.Context(), region)
	//if err != nil {
	//	throw(resp, 500, "Could not create session!", err.Error())
	//}

	// TODO: create following resources with following information
	// AWS IoT Policy (https://docs.aws.amazon.com/sdk-for-go/api/service/iot/#IoT.CreatePolicyWithContext)
	// 	ctx = req.Context()
	//	PolicyDocument: {
	//	  "Version": "2012-10-17",
	//	  "Statement": [
	//	    {
	//	      "Effect": "Allow",
	//	      "Action": "iot:Publish",
	//	      "Resource": "arn:aws:iot:<region>:<ds.cfg.AWSAccountID>:topic/<orgId>/<dsId>/<topic>"
	//	    },
	//	    {
	//	      "Effect": "Allow",
	//	      "Action": "iot:Connect",
	//	      "Resource": "arn:aws:iot:<region>:<ds.cfg.AWSAccountID>:client/<orgId>/<dsId>/<client>"
	//	    }
	//	  ]
	//	}
	//	PolicyName:
	//		if ds.cfg.Env == "production" {
	//			ds.cfg.Service + "/" + orgId + "/" + dsId + "/" + timestamp
	//		} else {
	//			ds.cfg.Service + "/" + ds.cfg.Env + "/" + orgId + "/" + dsId + "/" + timestamp
	//		}
	//	Tags: {
	//		Customer: orgId
	//		Datasource: dsId
	//		Service: ds.cfg.Service
	//		Zone: ds.cfg.Zone
	//		Environment: ds.cfg.Env
	//	}
	// AWS IoT Certificate (https://docs.aws.amazon.com/sdk-for-go/api/service/iot/#IoT.CreateKeysAndCertificateWithContext)
	// 	ctx = req.Context()
	//	Active: true
	// Attach the created policy to certificate (https://docs.aws.amazon.com/sdk-for-go/api/service/iot/#IoT.AttachPrincipalPolicyWithContext)

	// TODO: enable transaction lock. Say if certificate creation fails, corresponding policy must be deleted

	fmt.Fprintf(resp, "Hello, %q from %q", req.URL.Path, region)
}

type Certificate struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Topic  string `json:"topic"`
	Client string `json:"client"`
}

func (ds *MQTTDatasource) handleGetCertificates(resp http.ResponseWriter, req *http.Request) {
	// TODO: allow only get requests

	//region := req.URL.Query().Get("region")
	//if region == "" {
	//	throw(resp, 400, "Invalid or missing region!", "")
	//	return
	//}

	//svc, err := ds.authenticate(req.Context(), region)
	//if err != nil {
	//	throw(resp, 500, "Could not create session!", err.Error())
	//}

	// TODO: list certificates which have policies attached with tag 'Datasource' set to dsId
	// Response: [{
	//	id: <certificate id>,
	//	topic: <policy topic prefix>,
	//	client: <policy client prefix>,
	//	status: (active|inactive|revoked)
	// },{
	//	...
	// }]

	certificates := make([]Certificate, 0)
	responseBytes, err := json.Marshal(certificates)
	if err != nil {
		throw(resp, 500, "Could not list certificates!", err.Error())
	}

	_, _ = resp.Write(responseBytes)
}

func (ds *MQTTDatasource) handleCertificateSetActive(resp http.ResponseWriter, req *http.Request) {
	// TODO: allow only patch requests

	region := req.URL.Query().Get("region")
	if region == "" {
		throw(resp, 400, "Invalid or missing region!", "")
		return
	}

	id := req.URL.Query().Get("id")
	if id == "" {
		throw(resp, 400, "Invalid or missing certificate id!", "")
		return
	}

	//svc, err := ds.authenticate(req.Context(), region)
	//if err != nil {
	//	throw(resp, 500, "Could not create session!", err.Error())
	//}

	// TODO: Set status of certificate with id = id to ACTIVE (https://docs.aws.amazon.com/sdk-for-go/api/service/iot/#IoT.UpdateCertificateWithContext)

	fmt.Fprintf(resp, "Hello, \"%s?id=%s\" from %q", req.URL.Path, id, region)
}

func (ds *MQTTDatasource) handleCertificateSetInactive(resp http.ResponseWriter, req *http.Request) {
	// TODO: allow only patch requests

	region := req.URL.Query().Get("region")
	if region == "" {
		throw(resp, 400, "Invalid or missing region!", "")
		return
	}

	id := req.URL.Query().Get("id")
	if id == "" {
		throw(resp, 400, "Invalid or missing certificate id!", "")
		return
	}

	//svc, err := ds.authenticate(req.Context(), region)
	//if err != nil {
	//	throw(resp, 500, "Could not create session!", err.Error())
	//}

	// TODO: Set status of certificate with id = id to INACTIVE (https://docs.aws.amazon.com/sdk-for-go/api/service/iot/#IoT.UpdateCertificateWithContext)

	fmt.Fprintf(resp, "Hello, \"%s?id=%s\" from %q", req.URL.Path, id, region)
}

func (ds *MQTTDatasource) handleRevokeCertificate(resp http.ResponseWriter, req *http.Request) {
	// TODO: allow only patch requests

	region := req.URL.Query().Get("region")
	if region == "" {
		throw(resp, 400, "Invalid or missing region!", "")
		return
	}

	id := req.URL.Query().Get("id")
	if id == "" {
		throw(resp, 400, "Invalid or missing certificate id!", "")
		return
	}

	//svc, err := ds.authenticate(req.Context(), region)
	//if err != nil {
	//	throw(resp, 500, "Could not create session!", err.Error())
	//}

	// TODO: Set status of certificate with id = id to REVOKED (https://docs.aws.amazon.com/sdk-for-go/api/service/iot/#IoT.UpdateCertificateWithContext)

	fmt.Fprintf(resp, "Hello, \"%s?id=%s\" from %q", req.URL.Path, id, region)
}

func (ds *MQTTDatasource) handleDeleteCertificate(resp http.ResponseWriter, req *http.Request) {
	// TODO: allow only delete requests

	region := req.URL.Query().Get("region")
	if region == "" {
		throw(resp, 400, "Invalid or missing region!", "")
		return
	}

	id := req.URL.Query().Get("id")
	if id == "" {
		throw(resp, 400, "Invalid or missing certificate id!", "")
		return
	}

	//svc, err := ds.authenticate(req.Context(), region)
	//if err != nil {
	//	throw(resp, 500, "Could not create session!", err.Error())
	//}

	// TODO: delete the certificate and corresponding policy

	fmt.Fprintf(resp, "Hello, \"%s?id=%s\" from %q", req.URL.Path, id, region)
}
