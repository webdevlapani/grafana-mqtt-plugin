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

type CreateCertificate struct {
	Id          string `json:"id"`
	Status      string `json:"status"`
	Topic       string `json:"topic"`
	Client      string `json:"client"`
	PublicKey   string `json:"public_key"`
	PrivateKey  string `json:"private_key"`
	Certificate string `json:"certificate"`
	RootCA      string `json:"root_ca"`
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

	certificate := CreateCertificate{
		Id: "72dd3f3e77",
		Status: "ACTIVE",
		Topic: "/1/2/*",
		Client: "/1/2/*",
		PublicKey: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEApsssipafEUmLQZUtDtYd
nSmajTjtBnwqCIl6ucEjWg/gUvs+hiBegvlOmKjZS8gaU2mwUWcXGSlD31OPpcjq
H9R9hpLPyrfq1zFpT7sE/9soC7n6Th9S5VsTyxUMm9+e2YYVjfsqelpOKi0XSyBX
dI2vsulHl61TYotju7PDk5Lc7ZYUo6GyTtql+Hz5JR1KnumSktrOqOVWxbZpSPIr
GJ9ToudRyrPhIFTcWGkoBsQVBMJxvNe6V0b7gyRb414wr89jIO0VsPps6ronAAPu
+pLOJDO7SUIvG8nVxsa0CbTFJDUloobQBItbWs5YSsC25Q/RmSv8PVxblNlXFLXx
+QIDAQAB
-----END PUBLIC KEY-----`,
		PrivateKey: `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEApsssipafEUmLQZUtDtYdnSmajTjtBnwqCIl6ucEjWg/gUvs+
hiBegvlOmKjZS8gaU2mwUWcXGSlD31OPpcjqH9R9hpLPyrfq1zFpT7sE/9soC7n6
Th9S5VsTyxUMm9+e2YYVjfsqelpOKi0XSyBXdI2vsulHl61TYotju7PDk5Lc7ZYU
o6GyTtql+Hz5JR1KnumSktrOqOVWxbZpSPIrGJ9ToudRyrPhIFTcWGkoBsQVBMJx
vNe6V0b7gyRb414wr89jIO0VsPps6ronAAPu+pLOJDO7SUIvG8nVxsa0CbTFJDUl
oobQBItbWs5YSsC25Q/RmSv8PVxblNlXFLXx+QIDAQABAoIBADak88fHxv9j59Kp
q+Rjc7pMqgzAbK8mOKMpX2LCCvHzp5uoInjQ3AXu5bgQAXjZav6O7qwMqT2eDlV5
S+OVqlaZSDKxoJAapz6vOoBbliy4wSruWDoF+yOXLinnkIT0w1cinacxdV42fctF
kI8VXnGaBckIsmLX7yym3BrfryCGgoxFLw30dMJ0TgQSvU2JqVNqwsDiEy666y8K
YDtQqUGW0GO4SDARPzw4SI+3Z1AKJSmKGp88rD0fiKDYxUKZRYRsEuKuENxUxLfF
zdnkQkB6aJ3B9jjvHT6YoeID9aZm6c35KJtw4Jq4rVJy4VoZOlfeWIFDNTCOWh12
g4fG6t0CgYEA1kYpP9tw7zeiIPE8pAcsf8lZ071bn1R2/hW5Tsrjn9AN8v1EKGkz
EnjiUi0ki0IQkJeZcy1UDClkw0XCjxTTCWwTLKoJoCoVJcFxST50AgnHZlp4g4CG
kif3QounExKMzEuIB5aybI1Fcu66tDjphSe4YsEV1gelP7ANyj8BbrMCgYEAx0YN
fZM434c8tgD1do7OAKFQSfTbRBD7KKZiSo7Vdios17FBG4WpHZE3dqzAuyerSbiV
DRc7JZE3psAEqvDgKdTqTbqvAMui7/9MHSF7XlT//zV5HeBlaEGGz29G4isZIyaY
B810T96C8Qd0oynuIQW33LtOxDeOkAO9aNLcsqMCgYABrw117gCGMLa6cYrbcx77
ZhapnkxRBTXmKz+Iifmd8OGbLjhR5Pm8xGxq3uXxnjRJHpfbGtkVO2IKUssDmtNJ
uKqx6CgpNQtzf4CnZbE9rtv9Ruq5hdII5f2AbV6DvNqUZGeOP7XpOnb4Pz4CWowj
OrutMv078FVxGa4SD8qwFwKBgGeSqopdTc8ojE6Q2wQfH0VGkuONp7WOGey75hSY
fqxnKV2GXK/AXfDnPGurSJU9/hJYJOhj7bMN8l3yKbrrbadwacOyxyjjrrGNAPOX
JncWOORd17DGpA53GGmSjcYZ3nvdoGFV0SF+JpK+bEouDf4N6c2JcVwdADUsLHNi
PaF/AoGBALw1S9hDEy0kO8KouzM4N2gyGpZArFkynDEF7MdcXYzfh+r4EjY7eXe6
7JhM2bGNXS7/5BwxroXroqrnYX70WS11Q+ZKiB/iD3B0flFu1xPoFOTXvxLXVLTC
nB8B4tYXUKhBe7tdetdejC05lVLfxx9PRn9l6SWB0ZlMMwn7+pns
-----END RSA PRIVATE KEY-----`,
		Certificate: `-----BEGIN CERTIFICATE-----
MIIDWTCCAkGgAwIBAgIUD7znBhWBAx2UkHAejHvAYr0NvHcwDQYJKoZIhvcNAQEL
BQAwTTFLMEkGA1UECwxCQW1hem9uIFdlYiBTZXJ2aWNlcyBPPUFtYXpvbi5jb20g
SW5jLiBMPVNlYXR0bGUgU1Q9V2FzaGluZ3RvbiBDPVVTMB4XDTIwMDkwMzExMDIy
NloXDTQ5MTIzMTIzNTk1OVowHjEcMBoGA1UEAwwTQVdTIElvVCBDZXJ0aWZpY2F0
ZTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAKbLLIqWnxFJi0GVLQ7W
HZ0pmo047QZ8KgiJernBI1oP4FL7PoYgXoL5Tpio2UvIGlNpsFFnFxkpQ99Tj6XI
6h/UfYaSz8q36tcxaU+7BP/bKAu5+k4fUuVbE8sVDJvfntmGFY37KnpaTiotF0sg
V3SNr7LpR5etU2KLY7uzw5OS3O2WFKOhsk7apfh8+SUdSp7pkpLazqjlVsW2aUjy
KxifU6LnUcqz4SBU3FhpKAbEFQTCcbzXuldG+4MkW+NeMK/PYyDtFbD6bOq6JwAD
7vqSziQzu0lCLxvJ1cbGtAm0xSQ1JaKG0ASLW1rOWErAtuUP0Zkr/D1cW5TZVxS1
8fkCAwEAAaNgMF4wHwYDVR0jBBgwFoAUKw5fIbldApdyWIyoTSb3lnsTXAowHQYD
VR0OBBYEFAIqtKOavNt34+eyVm9GJkikruS/MAwGA1UdEwEB/wQCMAAwDgYDVR0P
AQH/BAQDAgeAMA0GCSqGSIb3DQEBCwUAA4IBAQCnEo12gBsTSz+6sA8NFJh940oP
8DgOTpYdo6WVfMKdrUXiUVc46QU+USMvMkgAjNr8g5l945u8JrWnqwa7q8y+2abw
h1c/Nmgx3pwXM2MSYvWP5CdvxlTIvxk9h24AbhqwQTlThkWLAKd5LSRfqhc/LD43
2RWDFrwCd7+Igakl+q8kZscuzzRtZUGgpQ7dTz3eOe2ELQWRU3SBC/2+33HoMdgN
WG/VAWFPrtQ6FlM3GyEBBDamxW/7boXnOv+vlrPH37B19Aoaz494kiQzN3/vj9bn
VeS/M8g9itrRh2oMcpmj9CdMDxt+qLoe86R0ru7JlGxXp0/YYO7iSJh5M7u0
-----END CERTIFICATE-----`,
		RootCA: `-----BEGIN CERTIFICATE-----
MIIDQTCCAimgAwIBAgITBmyfz5m/jAo54vB4ikPmljZbyjANBgkqhkiG9w0BAQsF
ADA5MQswCQYDVQQGEwJVUzEPMA0GA1UEChMGQW1hem9uMRkwFwYDVQQDExBBbWF6
b24gUm9vdCBDQSAxMB4XDTE1MDUyNjAwMDAwMFoXDTM4MDExNzAwMDAwMFowOTEL
MAkGA1UEBhMCVVMxDzANBgNVBAoTBkFtYXpvbjEZMBcGA1UEAxMQQW1hem9uIFJv
b3QgQ0EgMTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALJ4gHHKeNXj
ca9HgFB0fW7Y14h29Jlo91ghYPl0hAEvrAIthtOgQ3pOsqTQNroBvo3bSMgHFzZM
9O6II8c+6zf1tRn4SWiw3te5djgdYZ6k/oI2peVKVuRF4fn9tBb6dNqcmzU5L/qw
IFAGbHrQgLKm+a/sRxmPUDgH3KKHOVj4utWp+UhnMJbulHheb4mjUcAwhmahRWa6
VOujw5H5SNz/0egwLX0tdHA114gk957EWW67c4cX8jJGKLhD+rcdqsq08p8kDi1L
93FcXmn/6pUCyziKrlA4b9v7LWIbxcceVOF34GfID5yHI9Y/QCB/IIDEgEw+OyQm
jgSubJrIqg0CAwEAAaNCMEAwDwYDVR0TAQH/BAUwAwEB/zAOBgNVHQ8BAf8EBAMC
AYYwHQYDVR0OBBYEFIQYzIU07LwMlJQuCFmcx7IQTgoIMA0GCSqGSIb3DQEBCwUA
A4IBAQCY8jdaQZChGsV2USggNiMOruYou6r4lK5IpDB/G/wkjUu0yKGX9rbxenDI
U5PMCCjjmCXPI6T53iHTfIUJrU6adTrCC2qJeHZERxhlbI1Bjjt/msv0tadQ1wUs
N+gDS63pYaACbvXy8MWy7Vu33PqUXHeeE6V/Uq2V8viTO96LXFvKWlJbYK8U90vv
o/ufQJVtMVT8QtPHRh8jrdkPSHCa2XV4cdFyQzR1bldZwgJcJmApzyMZFo6IQ6XU
5MsI+yMRQ+hDKXJioaldXgjUkK642M4UwtBV8ob2xJNDd2ZhwLnoQdeXeGADbkpy
rqXRfboQnoZsG4q5WTP468SQvvG5
-----END CERTIFICATE-----`,
	}

	responseBytes, err := json.Marshal(certificate)
	if err != nil {
		throw(resp, 500, "Could not create certificate!", err.Error())
	}

	_, _ = resp.Write(responseBytes)
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

	//certificates := make([]Certificate, 0)

	// dummy certificates
	certificates := []Certificate{
		Certificate{Id: "798bff761e3d831831f0310b074b95a307e57b2dc83272bdd291d66df5ac8e64", Status: "ACTIVE", Topic: "/1/2/*", Client: "/1/2/*"},
		Certificate{Id: "91c2e97079d5c4688af516a58f7e4de03a9da40b5cdfefe8232473b0c5dcf1b0", Status: "INACTIVE", Topic: "/1/2/*", Client: "/1/2/*"},
		Certificate{Id: "0041fd2d647fb6835ea3ab6a0d13fd3c6c0aa7beb9bcb0191068cf081c3127b9", Status: "INACTIVE", Topic: "/1/2/*", Client: "/1/2/*"},
		Certificate{Id: "b5fcad63461f690c724bb76c506eefac06735275cac5bb7727195746d1b844e9", Status: "REVOKED", Topic: "/1/2/*", Client: "/1/2/*"},
        }

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
