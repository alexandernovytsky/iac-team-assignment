package sdk

import (
	"context"
	"fmt"

	"github.com/alexandernovytsky/iac-assignment/sdk/config"
	"github.com/alexandernovytsky/iac-assignment/sdk/gen"
	internal "github.com/alexandernovytsky/iac-assignment/sdk/internal/rest-client"
)

// WebhookRequest defines the structure for creating/updating a webhook
type CreateWebhookRequestV1 gen.V1CreateOutgoingWebhookRequest

// WebhookResponse defines the structure for webhook responses
type CreateWebhookResponseV1 gen.V1CreateOutgoingWebhookResponse

type GetWebhookResponseV1 gen.V1GetOutgoingWebhookResponse

// WebhooksClient provides methods for interacting with the Webhooks API
type WebhooksClient struct {
	restClient *internal.RestClient
	path       string
}

// NewWebhooksClient creates a new client for Webhooks
func NewWebhooksV1Client(restClient *internal.RestClient) *WebhooksClient {
	return &WebhooksClient{
		restClient: restClient,
		path:       config.GetResourcePath(config.V1, config.WebhooksPath),
	}
}

// Create creates a new webhook
func (c *WebhooksClient) Create(ctx context.Context, req *CreateWebhookRequestV1) (*CreateWebhookResponseV1, error) {
	var resp CreateWebhookResponseV1
	err := c.restClient.Post(ctx, c.path, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *WebhooksClient) Get(ctx context.Context, id string) (*GetWebhookResponseV1, error) {
	var resp GetWebhookResponseV1
	err := c.restClient.Get(ctx, fmt.Sprintf("%s/%s", c.path, id), &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
