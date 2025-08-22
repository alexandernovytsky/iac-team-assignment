package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/alexandernovytsky/iac-assignment/sdk"
	alerts "github.com/alexandernovytsky/iac-assignment/sdk/alerts/v3"
	"github.com/alexandernovytsky/iac-assignment/sdk/config"
	"github.com/alexandernovytsky/iac-assignment/sdk/gen"
	webhooks "github.com/alexandernovytsky/iac-assignment/sdk/webhooks/v1"
	"github.com/google/uuid"
)

func main() {
	// set config
	hookUrl := os.Getenv("CORALOGIX_WEBHOOK_URL")
	webhookId := uuid.New().String()

	// crate client generator
	clientCreator := sdk.NewClientCreator(
		config.RegionEU2,
		os.Getenv("CORALOGIX_API_KEY"),
		sdk.WithTimeout(10*time.Second),
		sdk.WithMaxRetries(3),
		sdk.WithBackoff(100*time.Millisecond),
	)

	ctx := context.Background()

	// create hook
	hc := clientCreator.Webhooks()
	if hc == nil {
		panic("failed to create webhooks client")
	}
	res, err := hc.Create(ctx, &webhooks.CreateWebhookRequestV1{
		Data: gen.V1OutgoingWebhookInputData{
			Type: gen.V1WebhookTypeGENERIC,
			GenericWebhook: &gen.V1GenericWebhookConfig{
				Method: gen.GenericWebhookConfigMethodTypeGET,
				Uuid:   webhookId,
				Headers: &map[string]string{
					"Content-Type": "application/json",
				},
			},
			Url:  &hookUrl,
			Name: "API Webhook - " + time.Now().Format("2006-01-02 15:04"),
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Webhook created successfully ID:", res.Id)

	// get external id of webhook for integration
	wh, err := hc.Get(ctx, res.Id)
	if err != nil {
		panic(err)
	}
	fmt.Println("Webhook fetched successfully, External ID:", wh.Webhook.ExternalId)

	// create alert
	alertsClient := clientCreator.Alerts()
	if alertsClient == nil {
		panic("failed to create alerts client")
	}

	errorLuceneQuery := "logRecord.severityNumber: 17"
	infoLuceneQuery := "logRecord.severityNumber: 9"
	timeWindowType := gen.LOGSRATIOTIMEWINDOWVALUEMINUTES10
	overridePriority := gen.ALERTDEFPRIORITYP2
	denominatorAlias := "Informative Logs"
	numeratorAlias := "Error Logs"
	threshold := float64(1.5)
	alert, err := alertsClient.Create(ctx, &alerts.CreateAlertRequestV3{
		Name:     "Error to Info Time Window Ratio - " + time.Now().Format("2006-01-02 15:04"),
		Priority: gen.ALERTDEFPRIORITYP2,
		Type:     gen.ALERTDEFTYPELOGSRATIOTHRESHOLD,
		LogsRatioThreshold: &gen.V3LogsRatioThresholdType{
			NumeratorAlias: &numeratorAlias,
			Numerator: gen.Alertsv3LogsFilter{
				SimpleFilter: &gen.V3LogsSimpleFilter{
					LuceneQuery: &errorLuceneQuery,
					LabelFilters: &gen.V3LabelFilters{
						ApplicationName: &[]gen.V3LabelFilterType{{Value: "sample-app"}},
					},
				},
			},
			DenominatorAlias: &denominatorAlias,
			Denominator: gen.Alertsv3LogsFilter{
				SimpleFilter: &gen.V3LogsSimpleFilter{
					LuceneQuery: &infoLuceneQuery,
					LabelFilters: &gen.V3LabelFilters{
						ApplicationName: &[]gen.V3LabelFilterType{{Value: "sample-app"}},
					},
				},
			},
			Rules: []gen.V3LogsRatioRules{{
				Condition: gen.V3LogsRatioCondition{
					Threshold: threshold,
					TimeWindow: gen.V3LogsRatioTimeWindow{
						LogsRatioTimeWindowSpecificValue: &timeWindowType,
					},
				},
				Override: gen.V3AlertDefOverride{
					Priority: &overridePriority,
				},
			},
			},
		},
		NotificationGroup: &gen.V3AlertDefNotificationGroup{
			Webhooks: &[]gen.V3AlertDefWebhooksSettings{
				{
					Integration: gen.Alertsv3IntegrationType{
						IntegrationId: &wh.Webhook.ExternalId,
					},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nAlert created successfully, Name: %s Type: %+v Enabled: %+v Webhook External ID: %d", alert.AlertDef.AlertDefProperties.Name, alert.AlertDef.AlertDefProperties.Type, *alert.AlertDef.AlertDefProperties.Enabled, *(*alert.AlertDef.AlertDefProperties.NotificationGroup.Webhooks)[0].Integration.IntegrationId)
}
