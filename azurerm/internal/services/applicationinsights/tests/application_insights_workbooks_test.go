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
	data := acceptance.BuildTestData(t, "azurerm_application_insights_workbook", "test")

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
	data := acceptance.BuildTestData(t, "azurerm_application_insights_workbook", "test")

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
	data := acceptance.BuildTestData(t, "azurerm_application_insights_workbook", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsWorkbookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationInsightsWorkbook_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsWorkbookExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "shared"),
					resource.TestCheckResourceAttr(data.ResourceName, "serialized_data", "{Book1:\"test\"}"),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "1.0"),
					// TODO: This is computed: resource.TestCheckResourceAttr(data.ResourceName, "workbook_id", "WBID"),
					resource.TestCheckResourceAttr(data.ResourceName, "category", "workbook"),
					//resource.TestCheckResourceAttr(data.ResourceName, "user_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr(data.ResourceName, "user_id", data.Client().Default.ClientID),
					resource.TestCheckResourceAttr(data.ResourceName, "source_resource_id", "/path/to/resource/id"),
					resource.TestCheckResourceAttr(data.ResourceName, "workbook_tags", "tag1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Hello", "World"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApplicationInsightsWorkbook_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_workbook", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsWorkbookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApplicationInsightsWorkbook_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsWorkbookExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "1.1"),
					resource.TestCheckResourceAttr(data.ResourceName, "category", "usage"),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "user"),
				),
			},
			{
				Config: testAccAzureRMApplicationInsightsWorkbook_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsWorkbookExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "2.0"),
					resource.TestCheckResourceAttr(data.ResourceName, "category", "tsg"),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "shared"),
				),
			},
			{
				Config: testAccAzureRMApplicationInsightsWorkbook_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApplicationInsightsWorkbookExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "1.1"),
					resource.TestCheckResourceAttr(data.ResourceName, "category", "usage"),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "user"),
				),
			},
		},
	})
}

func testCheckAzureRMApplicationInsightsWorkbookDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).AppInsights.WorkbooksClient
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
			return fmt.Errorf("Application Insights Workbook still exists:\n%#v", resp.WorkbookProperties)
		}
	}

	return nil
}

func testCheckAzureRMApplicationInsightsWorkbookExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).AppInsights.WorkbooksClient
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

func testAccAzureRMApplicationInsightsWorkbook_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-workbook-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMApplicationInsightsWorkbook_basic(data acceptance.TestData) string {
	init := testAccAzureRMApplicationInsightsWorkbook_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_application_insights_workbook" "test" {
  name                    = "acctestappinsightsworkbook-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  kind                    = "shared"
  category                = "workbook"
  user_id                 = "%s"
  workbook_tags           = ["tag1"]

  serialized_data = <<DATA
{
  "version": "Notebook/1.0",
  "items": [
    {
      "type": 1,
      "content": {
        "json": "test123"
      },
      "name": "text - 0"
    }
  ],
  "isLocked": false,
  "fallbackResourceIds": [
    ${azurerm_resource_group.test.id}
  ]
}
DATA

  lifecycle {
    ignore_changes = ["tags"]
  }
}
`, init, data.RandomInteger, data.Client().Default.ClientID)
}

func testAccAzureRMApplicationInsightsWorkbook_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApplicationInsightsWorkbook_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_application_insights_workbook" "import" {
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

func testAccAzureRMApplicationInsightsWorkbook_complete(data acceptance.TestData) string {
	init := testAccAzureRMApplicationInsightsWorkbook_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_application_insights_workbook" "test" {
  name                    = "acctestappinsightsworkbook-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  kind                    = "user"
  category                = "workbook"
  user_id                 = "%s"
  workbook_tags           = ["tag1"]
  source_resource_id      = "/path/to/resource/id"

  serialized_data = <<DATA
{
  "version": "Notebook/1.0",
  "items": [
    {
      "type": 1,
      "content": {
        "json": "test123"
      },
      "name": "text - 0"
    }
  ],
  "isLocked": false,
  "fallbackResourceIds": [
    ${azurerm_resource_group.test.id}
  ]
}
DATA

  tags {
    "hello": "world",
  }

  lifecycle {
    ignore_changes = ["tags"]
  }
}
`, init, data.RandomInteger, data.Client().Default.ClientID)
}
