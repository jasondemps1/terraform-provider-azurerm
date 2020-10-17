package applicationinsights

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2015-05-01/insights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

//func schemaWorkbookProperties() *schema.Schema {
//return &schema.Schema{
//Type:     schema.TypeList,
//MaxItems: 1,
//Optional: true,
//Computed: true,
//Elem: &schema.Resource{
//}
//}
//}

func resourceArmApplicationInsightsWorkbooks() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApplicationInsightsCreateUpdate,
		Read:   resourceArmApplicationInsightsRead,
		Update: resourceArmApplicationInsightsCreateUpdate,
		Delete: resourceArmApplicationInsightsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				//ValidateFunc: validation.NoZeroValues,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"application_insights_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"location": azure.SchemaLocation(),

			"tags": tags.Schema(),

			"kind": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(insights.SharedTypeKindUser),
					string(insights.SharedTypeKindShared),
				}, false),
			},

			"serialized_data": {
				Type:     schema.TypeString,
				Required: true,
				// TODO: Do we need?
				// ForceNew: true,
				ValidateFunc: validation.StringIsJSON,
				// TODO: Not sure if we need:
				//DiffSuppressFunc: structure.SuppressJsonDiff,
			},

			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"workbook_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},

			"category": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},

			"workbook_tags": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 0,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotWhiteSpace,
					//ValidateFunc:     validation.StringIsNotEmpty,
					//StateFunc:        location.StateFunc,
					//DiffSuppressFunc: location.DiffSuppressFunc,
				},
			},

			"user_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},

			"source_resource_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceArmApplicationInsightsWorkbooksCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.WorkbooksClient
	//billingClient := meta.(*clients.Client).AppInsights.BillingClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)

	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Application Insights Workbooks creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	appInsightsID := d.Get("application_insights_id").(string)

	id, err := azure.ParseAzureResourceID(appInsightsID)
	if err != nil {
		return err
	}

	appInsightsName := id.Path["components"]

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Application Insights Workbooks %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_application_insights_workbook", *existing.ID)
		}
	}

	//applicationType := d.Get("application_type").(string)
	//samplingPercentage := utils.Float(d.Get("sampling_percentage").(float64))
	//disableIpMasking := d.Get("disable_ip_masking").(bool)
	location := d.Get("location").(string)
	//tags := d.Get("tags").(map[string]*string)

	t := d.Get("tags").(map[string]interface{})
	tagKey := fmt.Sprintf("hidden-link:/subscriptions/%s/resourceGroups/%s/providers/microsoft.insights/components/%s", client.SubscriptionID, resGroup, appInsightsName)
	t[tagKey] = "Resource"

	kind := d.Get("kind").(string)
	serializedData := d.Get("serialized_data").(string)
	version := d.Get("version").(string)
	// Needed for flatten: workbookId := config["name"].(string)
	category := d.Get("category").(string)
	workbookTags := d.Get("workbook_tags").([]string)
	userID := d.Get("user_id").(string)
	sourceResourceID := d.Get("source_resource_id").(string)
	workbookID := d.Get("workbook_id").(string)
	//sharedTypeKind := insights.SharedTypeKind(fmt.Sprintf("SharedTypeKind%s", d.Get("kind").(string)))

	//propertiesRaw := d.Get("properties").([]interface{})
	//workbookProperties := expandProperties(propertiesRaw)

	workbook := insights.Workbook{
		Name:     &name,
		Location: &location,
		Kind:     insights.SharedTypeKind(kind),
		//WorkbookProperties: workbookProperties,
		WorkbookProperties: &insights.WorkbookProperties{
			Name:             &name,
			SerializedData:   &serializedData,
			Version:          &version,
			WorkbookID:       &workbookID,
			SharedTypeKind:   insights.SharedTypeKind(kind),
			Category:         &category,
			Tags:             &workbookTags,
			UserID:           &userID,
			SourceResourceID: &sourceResourceID,
		},
		Tags: tags.Expand(t),
	}

	//category :=

	//location := azure.NormalizeLocation(d.Get("location").(string))
	//t := d.Get("tags").(map[string]interface{})

	//workbookProperties := insights.WorkbookProperties{
	//Name:           &name,
	//SerializedData: &serializedData,
	//Version:        &version,
	////WorkbookID: , // TODO: How to generate a unique ID and string it?
	//SharedTypeKind: sharedTypeKind,

	//}

	//applicationInsightsComponentProperties := insights.ApplicationInsightsComponentProperties{
	//ApplicationID:      &name,
	//ApplicationType:    insights.ApplicationType(applicationType),
	//SamplingPercentage: samplingPercentage,
	//DisableIPMasking:   utils.Bool(disableIpMasking),
	//}

	//if v, ok := d.GetOk("retention_in_days"); ok {
	//applicationInsightsComponentProperties.RetentionInDays = utils.Int32(int32(v.(int)))
	//}

	//insightProperties := insights.ApplicationInsightsComponent{
	//Name:                                   &name,
	//Location:                               &location,
	//Kind:                                   &applicationType,
	//ApplicationInsightsComponentProperties: &applicationInsightsComponentProperties,
	//Tags:                                   tags.Expand(t),
	//}

	resp, err := client.CreateOrUpdate(ctx, resGroup, name, workbook)
	if err != nil {
		return fmt.Errorf("Error creating Application Insights Workbook %q (Resource Group %q): %+v", name, resGroup, err)
	}

	//billingRead, err := billingClient.Get(ctx, resGroup, name)
	//if err != nil {
	//return fmt.Errorf("Error read Application Insights Workbooks Billing Features %q (Resource Group %q): %+v", name, resGroup, err)
	//}

	//applicationInsightsComponentBillingFeatures := insights.ApplicationInsightsComponentBillingFeatures{
	//CurrentBillingFeatures: billingRead.CurrentBillingFeatures,
	//DataVolumeCap:          billingRead.DataVolumeCap,
	//}

	//if v, ok := d.GetOk("daily_data_cap_in_gb"); ok {
	//applicationInsightsComponentBillingFeatures.DataVolumeCap.Cap = utils.Float(v.(float64))
	//}

	//if v, ok := d.GetOk("daily_data_cap_notifications_disabled"); ok {
	//applicationInsightsComponentBillingFeatures.DataVolumeCap.StopSendNotificationWhenHitCap = utils.Bool(v.(bool))
	//}

	//if _, err = billingClient.Update(ctx, resGroup, name, applicationInsightsComponentBillingFeatures); err != nil {
	//return fmt.Errorf("Error update Application Insights Workbooks Billing Feature %q (Resource Group %q): %+v", name, resGroup, err)
	//}

	d.SetId(*resp.ID)

	return resourceArmApplicationInsightsWorkbooksRead(d, meta)
}

func resourceArmApplicationInsightsWorkbooksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.WorkbooksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading AzureRM Application Insights Workbooks '%s'", id)

	resGroup := id.ResourceGroup
	name := id.Path["workbooks"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Application Insights Workbook %q was not found in Resource Group %q - removing from state!", name, resGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM Application Insights Workbooks '%s': %+v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	appInsightsID := ""
	for i := range resp.Tags {
		if strings.HasPrefix(i, "hidden-link") {
			appInsightsID = strings.Split(i, ":")[1]
		}
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("kind", resp.Kind)
	d.Set("application_insights_id", appInsightsID)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.WorkbookProperties; props != nil {
		// It is possible that the root level `kind` in response is empty in some cases (see PR #8372 for more info)
		if resp.Kind == "" {
			d.Set("kind", props.SharedTypeKind)
		}
		d.Set("serialized_data", props.SerializedData)
		d.Set("version", props.Version)
		d.Set("workbook_id", props.WorkbookID)
		d.Set("category", props.Category)
		d.Set("user_id", props.UserID)
		d.Set("source_resource_id", props.SourceResourceID)
		d.Set("workbook_tags", props.Tags)

		//if config := props.Configuration; config != nil {
		//d.Set("configuration", config.WebTest)
		//}

		//if err := d.Set("workbook_tags", flattenApplicationInsightsWorkbookTags(props.Tags)); err != nil {
		//return fmt.Errorf("Error setting `workbook_tags`: %+v", err)
		//}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmApplicationInsightsWorkbooksDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.WorkbooksClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["workbooks"]

	log.Printf("[DEBUG] Deleting AzureRM Application Insights Workbooks '%s' (resource group '%s')", name, resGroup)

	resp, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("Error issuing AzureRM delete request for Application Insights Workbooks '%s': %+v", name, err)
	}

	return err
}

//func flattenApplicationInsightsWorkbookTags(input *[]insights.Workbook)

//func expandProperties(input []interface{}) *insights.WorkbookProperties {
//if len(input) == 0 {
//return nil
//}

//config := input[0].(map[string]interface{})

//name := config["name"].(string)
//serializedData := config["serialized_data"].(string)
//version := config["version"].(string)
//// Needed for flatten: workbookId := config["name"].(string)
//sharedTypeKind := config["kind"].(insights.SharedTypeKind)
//category := config["category"].(string)
//tags := config["tags"].([]string)
//userID := config["userId"].(string)
//sourceResourceID := config["sourceResourceId"].(string)

////keyData := ""
////if key, ok := linuxKeys[0].(map[string]interface{}); ok {
////keyData = key["key_data"].(string)
////}

//return &insights.WorkbookProperties{
//Name:             &name,
//SerializedData:   &serializedData,
//Version:          &version,
//SharedTypeKind:   sharedTypeKind,
//Category:         &category,
//Tags:             &tags,
//UserID:           &userID,
//SourceResourceID: &sourceResourceID,
//}
//}
