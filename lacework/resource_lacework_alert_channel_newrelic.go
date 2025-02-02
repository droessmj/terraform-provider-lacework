package lacework

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lacework/go-sdk/api"
)

func resourceLaceworkAlertChannelNewRelic() *schema.Resource {
	return &schema.Resource{
		Create: resourceLaceworkAlertChannelNewRelicCreate,
		Read:   resourceLaceworkAlertChannelNewRelicRead,
		Update: resourceLaceworkAlertChannelNewRelicUpdate,
		Delete: resourceLaceworkAlertChannelNewRelicDelete,

		Importer: &schema.ResourceImporter{
			State: importLaceworkIntegration,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"intg_guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"account_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"insert_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"test_integration": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to test the integration of an alert channel upon creation and modification",
			},
			"created_or_updated_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_or_updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"org_level": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceLaceworkAlertChannelNewRelicCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		relic    = api.NewNewRelicAlertChannel(d.Get("name").(string),
			api.NewRelicChannelData{
				AccountID: d.Get("account_id").(int),
				InsertKey: d.Get("insert_key").(string),
			},
		)
	)
	if !d.Get("enabled").(bool) {
		relic.Enabled = 0
	}

	log.Printf("[INFO] Creating %s integration with data:\n%+v\n", api.NewRelicChannelIntegration, relic)
	response, err := lacework.Integrations.CreateNewRelicAlertChannel(relic)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateNewRelicAlertChannelResponse(&response)
	if err != nil {
		return err
	}

	integration := response.Data[0]
	d.SetId(integration.IntgGuid)
	d.Set("name", integration.Name)
	d.Set("intg_guid", integration.IntgGuid)
	d.Set("enabled", integration.Enabled == 1)
	d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
	d.Set("type_name", integration.TypeName)
	d.Set("org_level", integration.IsOrg == 1)

	if d.Get("test_integration").(bool) {
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.NewRelicChannelIntegration, d.Id())
		if err := VerifyAlertChannelAndRollback(d.Id(), lacework); err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.NewRelicChannelIntegration, d.Id())
	}

	log.Printf("[INFO] Created %s integration with guid %s\n", api.NewRelicChannelIntegration, integration.IntgGuid)
	return nil
}

func resourceLaceworkAlertChannelNewRelicRead(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Reading %s integration with guid %s\n", api.NewRelicChannelIntegration, d.Id())
	response, err := lacework.Integrations.GetNewRelicAlertChannel(d.Id())
	if err != nil {
		return err
	}

	for _, integration := range response.Data {
		if integration.IntgGuid == d.Id() {
			d.Set("name", integration.Name)
			d.Set("intg_guid", integration.IntgGuid)
			d.Set("enabled", integration.Enabled == 1)
			d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
			d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
			d.Set("type_name", integration.TypeName)
			d.Set("org_level", integration.IsOrg == 1)
			d.Set("account_id", integration.Data.AccountID)
			d.Set("insert_key", integration.Data.InsertKey)

			log.Printf("[INFO] Read %s integration with guid %s\n",
				api.NewRelicChannelIntegration, integration.IntgGuid)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceLaceworkAlertChannelNewRelicUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		lacework = meta.(*api.Client)
		relic    = api.NewNewRelicAlertChannel(d.Get("name").(string),
			api.NewRelicChannelData{
				AccountID: d.Get("account_id").(int),
				InsertKey: d.Get("insert_key").(string),
			},
		)
	)

	if !d.Get("enabled").(bool) {
		relic.Enabled = 0
	}

	relic.IntgGuid = d.Id()

	log.Printf("[INFO] Updating %s integration with data:\n%+v\n", api.NewRelicChannelIntegration, relic)
	response, err := lacework.Integrations.UpdateNewRelicAlertChannel(relic)
	if err != nil {
		return err
	}

	log.Println("[INFO] Verifying server response data")
	err = validateNewRelicAlertChannelResponse(&response)
	if err != nil {
		return err
	}

	integration := response.Data[0]
	d.Set("name", integration.Name)
	d.Set("intg_guid", integration.IntgGuid)
	d.Set("enabled", integration.Enabled == 1)
	d.Set("created_or_updated_time", integration.CreatedOrUpdatedTime)
	d.Set("created_or_updated_by", integration.CreatedOrUpdatedBy)
	d.Set("type_name", integration.TypeName)
	d.Set("org_level", integration.IsOrg == 1)

	if d.Get("test_integration").(bool) {
		log.Printf("[INFO] Testing %s integration for guid %s\n", api.NewRelicChannelIntegration, d.Id())
		if err := lacework.V2.AlertChannels.Test(d.Id()); err != nil {
			return err
		}
		log.Printf("[INFO] Tested %s integration with guid %s successfully\n", api.NewRelicChannelIntegration, d.Id())
	}

	log.Printf("[INFO] Updated %s integration with guid %s\n", api.NewRelicChannelIntegration, d.Id())
	return nil
}

func resourceLaceworkAlertChannelNewRelicDelete(d *schema.ResourceData, meta interface{}) error {
	lacework := meta.(*api.Client)

	log.Printf("[INFO] Deleting %s integration with guid %s\n", api.NewRelicChannelIntegration, d.Id())
	_, err := lacework.Integrations.Delete(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted %s integration with guid %s\n", api.NewRelicChannelIntegration, d.Id())
	return nil
}

func validateNewRelicAlertChannelResponse(response *api.NewRelicAlertChannelResponse) error {
	if len(response.Data) == 0 {
		msg := `
Unable to read sever response data. (empty 'data' field)

This was an unexpected behavior, verify that your integration has been
created successfully and report this issue to support@lacework.net
`
		return fmt.Errorf(msg)
	}

	if len(response.Data) > 1 {
		msg := `
There is more that one integration inside the server response data.

List of integrations:
`
		for _, integration := range response.Data {
			msg = msg + fmt.Sprintf("\t%s: %s\n", integration.IntgGuid, integration.Name)
		}
		msg = msg + unexpectedBehaviorMsg()
		return fmt.Errorf(msg)
	}

	return nil
}
