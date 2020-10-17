package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMApplicationInsightsWorkbook_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_workbook_test", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsWorkbookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationInsightsWorkbook_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsWorkbookExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationInsightsWorkbook_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_workbook_test", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsWorkbookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationInsightsWorkbook_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsWorkbookExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMApplicationInsightsWorkbook_requiresImport),
		},
	})
}

func TestAccAzureRMApplicationInsightsWorkbook_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_workbook_test", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsWorkbookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationInsightsWorkbook_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsWorkbookExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationInsightsWorkbook_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_workbook_test", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsWorkbookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationInsightsWorkbook_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsWorkbookExists(data.ResourceName),
					//resource.TestCheckResourceAttr(data.ResourceName, "geo_locations.#", "1"),
					//resource.TestCheckResourceAttr(data.ResourceName, "frequency", "300"),
					//resource.TestCheckResourceAttr(data.ResourceName, "timeout", "30"),
				),
			},
			{
				Config: testAccAzureRMApplicationInsightsWorkbook_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsWorkbookExists(data.ResourceName),
					//resource.TestCheckResourceAttr(data.ResourceName, "geo_locations.#", "2"),
					//resource.TestCheckResourceAttr(data.ResourceName, "frequency", "900"),
					//resource.TestCheckResourceAttr(data.ResourceName, "timeout", "120"),
				),
			},
			{
				Config: testAccAzureRMApplicationInsightsWorkbook_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsWorkbookExists(data.ResourceName),
					//resource.TestCheckResourceAttr(data.ResourceName, "geo_locations.#", "1"),
					//resource.TestCheckResourceAttr(data.ResourceName, "frequency", "300"),
					//resource.TestCheckResourceAttr(data.ResourceName, "timeout", "30"),
				),
			},
		},
	})
}

func testCheckAzureRMApplicationInsightsWorkbookDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).AppInsights.WorkbookClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_application_insights_workbook" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Application Insights Workbook still exists:\n%#v", resp.ApplicationInsightsComponentProperties)
		}
	}

	return nil
}

func testCheckAzureRMApplicationInsightsWorkbookExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).AppInsights.WorkbookClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for App Insights: %s", name)
		}

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on appInsightsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Application Insights Workbook '%q' (resource group: '%q') does not exist", name, resourceGroup)
		}

		return nil
	}
}

func TestAccAzureRMApplicationInsightsWorkbook_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_workbook", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsWorkbookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationInsightsWorkbook_complete(data, "web"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsWorkbookExists(data.ResourceName),
					//resource.TestCheckResourceAttr(data.ResourceName, "application_type", "web"),
					//resource.TestCheckResourceAttr(data.ResourceName, "retention_in_days", "120"),
					//resource.TestCheckResourceAttr(data.ResourceName, "sampling_percentage", "50"),
					//resource.TestCheckResourceAttr(data.ResourceName, "daily_data_cap_in_gb", "50"),
					//resource.TestCheckResourceAttr(data.ResourceName, "daily_data_cap_notifications_disabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Hello", "World"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMApplicationInsightsWorkbook_basic(data acceptance.TestData, applicationType string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_application_insights_workbook" "test" {
  name                    = "acctestappinsightsworkbook-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  application_insights_id = azurerm_application_insights.test.id
  kind                    = "User"
  category                = "category%d"
  user_id                 = # TODO: AzureAD User ID fake
  workbook_tags           = ["tag1"]

  serialized_data = "{}" # TODO

  lifecycle {
    ignore_changes = ["tags"]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, applicationType)
}

func testAccAzureRMApplicationInsightsWorkbook_requiresImport(data acceptance.TestData, applicationType string) string {
	template := testAccAzureRMApplicationInsightsWorkbook_basic(data, applicationType)
	return fmt.Sprintf(`
%s

resource "azurerm_application_insights_web_test" "import" {
  name                    = azurerm_application_insights_web_test.test.name
  location                = azurerm_application_insights_web_test.test.location
  resource_group_name     = azurerm_application_insights_web_test.test.resource_group_name
  application_insights_id = azurerm_application_insights_web_test.test.application_insights_id
  kind                    = azurerm_application_insights_web_test.test.kind
  category                = azurerm_application_insights_web_test.test.category
  serialized_data         = azurerm_application_insights_web_test.test.serialized_data
  user_id                 = azurerm_application_insights_web_test.test.user_id
  workbook_tags           = azurerm_application_insights_web_test.test.workbook_tags
}
`, template)
}

func testAccAzureRMApplicationInsightsWorkbook_complete(data acceptance.TestData, applicationType string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_application_insights_workbook" "test" {
  name                    = "acctestappinsightsworkbook-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  application_insights_id = azurerm_application_insights.test.id
  kind                    = "User"
  category                = "category%d"
  user_id                 = # TODO: AzureAD User ID fake
  workbook_tags           = ["tag1"]

  serialized_data = "{}" # TODO

  lifecycle {
    ignore_changes = ["tags"]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, applicationType)
}
