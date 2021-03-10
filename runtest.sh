#!/bin/bash

export ARM_CLIENT_ID="46c9fd62-54f7-4ffd-93f9-db1014dc6340"
export ARM_CLIENT_SECRET="nMQ_LGVGxjoY2E5ofn1.LygDI-6hEBC-VQ"
export ARM_SUBSCRIPTION_ID="22dad228-206b-434c-be66-d19228f01226"
export ARM_TENANT_ID="72f988bf-86f1-41af-91ab-2d7cd011db47"
export ARM_ENVIRONMENT="public"
export ARM_TEST_LOCATION="eastus"
export ARM_TEST_LOCATION_ALT="eastus2"
export ARM_TEST_LOCATION_ALT2="westus2"

make acctests SERVICE='mssql' TESTARGS="-run=TestAccMsSqlDatabase_geoBackupPolicy" TESTTIMEOUT='60m'