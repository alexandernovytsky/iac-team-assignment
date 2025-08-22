package sdk

import (
	"context"

	"github.com/alexandernovytsky/iac-assignment/sdk/config"
	"github.com/alexandernovytsky/iac-assignment/sdk/gen"
	internal "github.com/alexandernovytsky/iac-assignment/sdk/internal/rest-client"
)

// CreateAlertRequestV3 defines the structure for creating/updating an alert
type CreateAlertRequestV3 gen.V3AlertDefProperties

// CreateAlertResponseV3 defines the structure for alert responses
type CreateAlertResponseV3 gen.V3CreateAlertDefResponse

// AlertsClient provides methods for interacting with the Alerts API
type AlertsClient struct {
	rest *internal.RestClient
	path string
}

// NewAlertsClient creates a new client for Alerts
func NewAlertsV3Client(restClient *internal.RestClient) *AlertsClient {
	return &AlertsClient{
		rest: restClient,
		path: config.GetResourcePath(config.V3, config.AlertsPath),
	}
}

// Create creates a new alert
func (c *AlertsClient) Create(ctx context.Context, req *CreateAlertRequestV3) (*CreateAlertResponseV3, error) {
	var resp CreateAlertResponseV3
	err := c.rest.Post(ctx, c.path, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
