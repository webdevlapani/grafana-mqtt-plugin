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

	err := ds.createStorage(req.Context(), region)
	if err != nil {
		throw(resp, 500, "Could not provision storage!", err.Error())
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

	err := ds.createStorage(req.Context(), region)
	if err != nil {
		throw(resp, 500, "Could not provision storage!", err.Error())
		return
	}

	// TODO read the following fields from POST json body
	//	topic          string. optional. "*" if empty
	//	client         string. optional. "*" if empty

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

	region := req.URL.Query().Get("region")
	if region == "" {
		throw(resp, 400, "Invalid or missing region!", "")
		return
	}

	err := ds.createStorage(req.Context(), region)
	if err != nil {
		throw(resp, 500, "Could not provision storage!", err.Error())
		return
	}

	//svc, err := ds.authenticate(req.Context(), region)
	//if err != nil {
	//	throw(resp, 500, "Could not create session!", err.Error())
	//}

	// TODO: list certificates which have policies for current orgId and dsId
	// Steps:
	//	1. List all policies (https://docs.aws.amazon.com/sdk-for-go/api/service/iot/#IoT.ListPoliciesWithContext)
	//	2. Filter relevant polices
	//	3. Get certificates for the policies (https://docs.aws.amazon.com/sdk-for-go/api/service/iot/#IoT.ListPolicyPrincipalsWithContext)
	//	4. Compile the data into response
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

	err := ds.createStorage(req.Context(), region)
	if err != nil {
		throw(resp, 500, "Could not provision storage!", err.Error())
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

	err := ds.createStorage(req.Context(), region)
	if err != nil {
		throw(resp, 500, "Could not provision storage!", err.Error())
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

	err := ds.createStorage(req.Context(), region)
	if err != nil {
		throw(resp, 500, "Could not provision storage!", err.Error())
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

	err := ds.createStorage(req.Context(), region)
	if err != nil {
		throw(resp, 500, "Could not provision storage!", err.Error())
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

	// TODO: delete the corresponding policy and certificate
	// Steps:
	//	1. List policy associated with the certificate (https://docs.aws.amazon.com/sdk-for-go/api/service/iot/#IoT.ListPrincipalPoliciesWithContext)
	//	2. Delete the policy if policy.tags.CA != 1 (https://docs.aws.amazon.com/sdk-for-go/api/service/iot/#IoT.DeletePolicyWithContext)
	//	3. Delete the certificate (https://docs.aws.amazon.com/sdk-for-go/api/service/iot/#IoT.DeleteCertificateWithContext)

	fmt.Fprintf(resp, "Hello, \"%s?id=%s\" from %q", req.URL.Path, id, region)
}

func (ds *MQTTDatasource) handleRegisterCertificate(resp http.ResponseWriter, req *http.Request) {
	// TODO: allow only post requests

	region := req.URL.Query().Get("region")
	if region == "" {
		throw(resp, 400, "Invalid or missing region!", "")
		return
	}

	err := ds.createStorage(req.Context(), region)
	if err != nil {
		throw(resp, 500, "Could not provision storage!", err.Error())
		return
	}

	// TODO read the following fields from POST json body
	//	certificate    string. required.
	//	ca_certificate string. optional.
	//	topic          string. optional. "*" if empty
	//	client         string. optional. "*" if empty

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
	// AWS IoT Certificate
	// 	if certificate and ca https://docs.aws.amazon.com/sdk-for-go/api/service/iot/#IoT.RegisterCertificateWithContext
	// 	if only certificate https://docs.aws.amazon.com/sdk-for-go/api/service/iot/#IoT.RegisterCertificateWithoutCAWithContext
	// Attach the created policy to certificate (https://docs.aws.amazon.com/sdk-for-go/api/service/iot/#IoT.AttachPrincipalPolicyWithContext)

	// TODO: enable transaction lock. Say if certificate creation fails, corresponding policy must be deleted

	fmt.Fprintf(resp, "Hello, %q from %q", req.URL.Path, region)
}

type CA struct {
	Id               string `json:"id"`
	Status           string `json:"status"`
	AutoRegistration string `json:"auto_registration"`
	Topic            string `json:"topic"`
	Client           string `json:"client"`
}

func (ds *MQTTDatasource) handleGetCA(resp http.ResponseWriter, req *http.Request) {
	// TODO: allow only get requests

	region := req.URL.Query().Get("region")
	if region == "" {
		throw(resp, 400, "Invalid or missing region!", "")
		return
	}

	err := ds.createStorage(req.Context(), region)
	if err != nil {
		throw(resp, 500, "Could not provision storage!", err.Error())
		return
	}

	//svc, err := ds.authenticate(req.Context(), region)
	//if err != nil {
	//	throw(resp, 500, "Could not create session!", err.Error())
	//}

	// TODO: list certificates for current orgId and dsId
	// Steps:
	//	1. List all ca certificates (https://docs.aws.amazon.com/sdk-for-go/api/service/iot/#IoT.ListCACertificatesWithContext)
	//	2. Filter by orgId and dsId
	//	4. Compile the data into response
	// Response: [{
	//	id: <certificate id>,
	//	topic: <tags.topic>,
	//	client: <tags.client>,
	//	auto_registraion: <auto_registraion>,
	//	status: <status>,
	// },{
	//	...
	// }]

	certificates := make([]CA, 0)
	responseBytes, err := json.Marshal(certificates)
	if err != nil {
		throw(resp, 500, "Could not list certificates!", err.Error())
	}

	_, _ = resp.Write(responseBytes)
}

func (ds *MQTTDatasource) handleGetRegistrationCode(resp http.ResponseWriter, req *http.Request) {
	// TODO: allow only get requests

	region := req.URL.Query().Get("region")
	if region == "" {
		throw(resp, 400, "Invalid or missing region!", "")
		return
	}

	err := ds.createStorage(req.Context(), region)
	if err != nil {
		throw(resp, 500, "Could not provision storage!", err.Error())
		return
	}

	//svc, err := ds.authenticate(req.Context(), region)
	//if err != nil {
	//	throw(resp, 500, "Could not create session!", err.Error())
	//}

	// TODO https://docs.aws.amazon.com/sdk-for-go/api/service/iot/#IoT.GetRegistrationCodeWithContext

	fmt.Fprintf(resp, "Hello, %q from %q", req.URL.Path, region)
}

func (ds *MQTTDatasource) handleRegisterCA(resp http.ResponseWriter, req *http.Request) {
	// TODO: allow only post requests

	region := req.URL.Query().Get("region")
	if region == "" {
		throw(resp, 400, "Invalid or missing region!", "")
		return
	}

	err := ds.createStorage(req.Context(), region)
	if err != nil {
		throw(resp, 500, "Could not provision storage!", err.Error())
		return
	}

	// TODO read the following fields from POST json body
	//	certificate    string. required.
	//	ca_certificate string. optional.
	//	topic          string. optional. "*" if empty
	//	client         string. optional. "*" if empty

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
	//		CA: 1
	//	}
	// AWS Register CA Certificate
	// 	https://docs.aws.amazon.com/sdk-for-go/api/service/iot/#IoT.RegisterCACertificateWithContext
	//	Tags: {
	//		Customer: orgId
	//		Datasource: dsId
	//		Service: ds.cfg.Service
	//		Zone: ds.cfg.Zone
	//		Environment: ds.cfg.Env
	//		Policy: <PolicyName>
	//		Topic: <topic>
	//		Client: <topic>
	//	}

	// TODO: enable transaction lock. Say if certificate creation fails, corresponding policy must be deleted


	fmt.Fprintf(resp, "Hello, %q from %q", req.URL.Path, region)
}

func (ds *MQTTDatasource) handleCASetInactive(resp http.ResponseWriter, req *http.Request) {
	// TODO: allow only patch requests

	region := req.URL.Query().Get("region")
	if region == "" {
		throw(resp, 400, "Invalid or missing region!", "")
		return
	}

	err := ds.createStorage(req.Context(), region)
	if err != nil {
		throw(resp, 500, "Could not provision storage!", err.Error())
		return
	}

	//svc, err := ds.authenticate(req.Context(), region)
	//if err != nil {
	//	throw(resp, 500, "Could not create session!", err.Error())
	//}

	// TODO https://docs.aws.amazon.com/sdk-for-go/api/service/iot/#IoT.UpdateCACertificateWithContext

	fmt.Fprintf(resp, "Hello, %q from %q", req.URL.Path, region)
}

func (ds *MQTTDatasource) handleCASetActive(resp http.ResponseWriter, req *http.Request) {
	// TODO: allow only patch requests

	region := req.URL.Query().Get("region")
	if region == "" {
		throw(resp, 400, "Invalid or missing region!", "")
		return
	}

	err := ds.createStorage(req.Context(), region)
	if err != nil {
		throw(resp, 500, "Could not provision storage!", err.Error())
		return
	}

	//svc, err := ds.authenticate(req.Context(), region)
	//if err != nil {
	//	throw(resp, 500, "Could not create session!", err.Error())
	//}

	// TODO https://docs.aws.amazon.com/sdk-for-go/api/service/iot/#IoT.UpdateCACertificateWithContext

	fmt.Fprintf(resp, "Hello, %q from %q", req.URL.Path, region)
}

func (ds *MQTTDatasource) handleCAEnableAutoRegistration(resp http.ResponseWriter, req *http.Request) {
	// TODO: allow only patch requests

	region := req.URL.Query().Get("region")
	if region == "" {
		throw(resp, 400, "Invalid or missing region!", "")
		return
	}

	err := ds.createStorage(req.Context(), region)
	if err != nil {
		throw(resp, 500, "Could not provision storage!", err.Error())
		return
	}

	//svc, err := ds.authenticate(req.Context(), region)
	//if err != nil {
	//	throw(resp, 500, "Could not create session!", err.Error())
	//}

	// TODO https://docs.aws.amazon.com/sdk-for-go/api/service/iot/#IoT.UpdateCACertificateWithContext

	fmt.Fprintf(resp, "Hello, %q from %q", req.URL.Path, region)
}

func (ds *MQTTDatasource) handleCADisableAutoRegistration(resp http.ResponseWriter, req *http.Request) {
	// TODO: allow only patch requests

	region := req.URL.Query().Get("region")
	if region == "" {
		throw(resp, 400, "Invalid or missing region!", "")
		return
	}

	err := ds.createStorage(req.Context(), region)
	if err != nil {
		throw(resp, 500, "Could not provision storage!", err.Error())
		return
	}

	//svc, err := ds.authenticate(req.Context(), region)
	//if err != nil {
	//	throw(resp, 500, "Could not create session!", err.Error())
	//}

	// TODO https://docs.aws.amazon.com/sdk-for-go/api/service/iot/#IoT.UpdateCACertificateWithContext

	fmt.Fprintf(resp, "Hello, %q from %q", req.URL.Path, region)
}

func (ds *MQTTDatasource) handleDeleteCA(resp http.ResponseWriter, req *http.Request) {
	// TODO: allow only delete requests

	region := req.URL.Query().Get("region")
	if region == "" {
		throw(resp, 400, "Invalid or missing region!", "")
		return
	}

	err := ds.createStorage(req.Context(), region)
	if err != nil {
		throw(resp, 500, "Could not provision storage!", err.Error())
		return
	}

	//svc, err := ds.authenticate(req.Context(), region)
	//if err != nil {
	//	throw(resp, 500, "Could not create session!", err.Error())
	//}

	// TODO:
	// 1. List certificates associated (https://docs.aws.amazon.com/sdk-for-go/api/service/iot/#IoT.ListCertificatesByCAWithContext)
	// 2. Delete these certiicates (https://docs.aws.amazon.com/sdk-for-go/api/service/iot/#IoT.DeleteCertificateWithContext)
	// 3. Delete corresponding policy for this CA
	// 4. Delete CA (https://docs.aws.amazon.com/sdk-for-go/api/service/iot/#IoT.DeleteCACertificateWithContext)

	fmt.Fprintf(resp, "Hello, %q from %q", req.URL.Path, region)
}
