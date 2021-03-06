package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// newDatasource returns datasource.ServeOpts.
func newDatasource() datasource.ServeOpts {
	// creates a instance manager for your plugin. The function passed
	// into `NewInstanceManger` is called when the instance is created
	// for the first time or when a datasource configuration changed.
	im := datasource.NewInstanceManager(newDataSourceInstance)
	ds := &MQTTDatasource{
		im: im,
		cfg: loadConfig(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/endpoint", ds.handleGetEndpoint)

	mux.HandleFunc("/certificates", ds.handleGetCertificates)
	mux.HandleFunc("/certificates/create", ds.handleCreateCertificate)
	mux.HandleFunc("/certificates/register", ds.handleRegisterCertificate)
	mux.HandleFunc("/certificates/revoke", ds.handleRevokeCertificate)
	mux.HandleFunc("/certificates/set-inactive", ds.handleCertificateSetActive)
	mux.HandleFunc("/certificates/set-active", ds.handleCertificateSetInactive)
	mux.HandleFunc("/certificates/delete", ds.handleDeleteCertificate)

	mux.HandleFunc("/ca", ds.handleGetCA)
	mux.HandleFunc("/ca/registration-code", ds.handleGetRegistrationCode)
	mux.HandleFunc("/ca/register", ds.handleRegisterCA)
	mux.HandleFunc("/ca/set-inactive", ds.handleCASetInactive)
	mux.HandleFunc("/ca/set-active", ds.handleCASetActive)
	mux.HandleFunc("/ca/enable-auto-registration", ds.handleCAEnableAutoRegistration)
	mux.HandleFunc("/ca/disable-auto-registration", ds.handleCADisableAutoRegistration)
	mux.HandleFunc("/ca/delete", ds.handleDeleteCA)

	return datasource.ServeOpts{
		QueryDataHandler:    ds,
		CallResourceHandler: httpadapter.New(mux),
		CheckHealthHandler:  ds,
	}
}

type MQTTDatasource struct {
	// The instance manager can help with lifecycle management
	// of datasource instances in plugins. It's not a requirements
	// but a best practice that we recommend that you follow.
	im  instancemgmt.InstanceManager
	cfg *Cfg
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifer).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (ds *MQTTDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	// create response struct
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := ds.query(ctx, q)

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	return response, nil
}

type queryModel struct {
	Format string `json:"format"`
}

func (ds *MQTTDatasource) query(ctx context.Context, query backend.DataQuery) backend.DataResponse {
	// Unmarshal the json into our queryModel
	var qm queryModel

	response := backend.DataResponse{}

	response.Error = json.Unmarshal(query.JSON, &qm)
	if response.Error != nil {
		return response
	}

	// Log a warning if `Format` is empty.
	if qm.Format == "" {
		log.DefaultLogger.Warn("format is empty. defaulting to time series")
	}

	// create data frame response
	frame := data.NewFrame("response")

	// add the time dimension
	frame.Fields = append(frame.Fields,
		data.NewField("time", nil, []time.Time{query.TimeRange.From, query.TimeRange.To}),
	)

	// add values
	frame.Fields = append(frame.Fields,
		data.NewField("values", nil, []int64{10, 20}),
	)

	// add the frames to the response
	response.Frames = append(response.Frames, frame)

	return response
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (ds *MQTTDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	var status = backend.HealthStatusOk
	var message = "Data source is working"

	if rand.Int()%2 == 0 {
		status = backend.HealthStatusError
		message = "randomized error"
	}

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}

type instanceSettings struct {
	httpClient *http.Client
}

func newDataSourceInstance(setting backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	return &instanceSettings{
		httpClient: &http.Client{},
	}, nil
}

func (s *instanceSettings) Dispose() {
	// Called before creatinga a new instance to allow plugin authors
	// to cleanup.
}
