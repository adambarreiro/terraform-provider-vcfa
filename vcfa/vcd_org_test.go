//go:build org || tm || ALL || functional

package vcfa

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVcdTmOrg(t *testing.T) {
	preTestChecks(t)
	skipIfNotSysAdmin(t)

	var params = StringMap{
		"Testname": t.Name(),

		"Tags": "tm",
	}
	testParamsNotEmpty(t, params)

	configText1 := templateFill(testAccVcdTmOrgStep1, params)
	params["FuncName"] = t.Name() + "-step2"
	configText2 := templateFill(testAccVcdTmOrgStep2, params)
	params["FuncName"] = t.Name() + "-step3"
	configText3 := templateFill(testAccVcdTmOrgStep3DS, params)

	debugPrintf("#[DEBUG] CONFIGURATION step1: %s\n", configText1)
	debugPrintf("#[DEBUG] CONFIGURATION step2: %s\n", configText2)
	debugPrintf("#[DEBUG] CONFIGURATION step3: %s\n", configText3)
	if vcfaShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configText1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vcfa_org.test", "name", t.Name()),
					resource.TestCheckResourceAttr("vcfa_org.test", "display_name", "terraform-test"),
					resource.TestCheckResourceAttr("vcfa_org.test", "description", "terraform test"),
					resource.TestCheckResourceAttr("vcfa_org.test", "is_enabled", "true"),
					resource.TestCheckResourceAttr("vcfa_org.test", "is_subprovider", "false"),
					resource.TestMatchResourceAttr("vcfa_org.test", "managed_by_id", regexp.MustCompile("^urn:vcloud:org:")),
					resource.TestCheckResourceAttr("vcfa_org.test", "managed_by_name", "System"),
					resource.TestCheckResourceAttr("vcfa_org.test", "is_classic_tenant", "false"),
				),
			},
			{
				Config: configText2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vcfa_org.test", "name", t.Name()+"-updated"),
					resource.TestCheckResourceAttr("vcfa_org.test", "display_name", "terraform-test"),
					resource.TestCheckResourceAttr("vcfa_org.test", "description", ""),
					resource.TestCheckResourceAttr("vcfa_org.test", "is_enabled", "false"),
					resource.TestCheckResourceAttr("vcfa_org.test", "is_subprovider", "false"),
					resource.TestMatchResourceAttr("vcfa_org.test", "managed_by_id", regexp.MustCompile("^urn:vcloud:org:")),
					resource.TestCheckResourceAttr("vcfa_org.test", "managed_by_name", "System"),
					resource.TestCheckResourceAttr("vcfa_org.test", "is_classic_tenant", "false"),
				),
			},
			{
				Config: configText3,
				Check: resource.ComposeTestCheckFunc(
					resourceFieldsEqual("vcfa_org.test", "data.vcfa_org.test", nil),
				),
			},
			{
				ResourceName:      "vcfa_org.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     params["Testname"].(string),
			},
		},
	})

	postTestChecks(t)
}

const testAccVcdTmOrgStep1 = `
resource "vcfa_org" "test" {
  name         = "{{.Testname}}"
  display_name = "terraform-test"
  description  = "terraform test"
  is_enabled   = true
}
`

const testAccVcdTmOrgStep2 = `
resource "vcfa_org" "test" {
  name         = "{{.Testname}}-updated"
  display_name = "terraform-test"
  description  = ""
  is_enabled   = false
}
`

const testAccVcdTmOrgStep3DS = testAccVcdTmOrgStep1 + `
data "vcfa_org" "test" {
  name = vcfa_org.test.name
}
`

func TestAccVcdTmOrgSubProvider(t *testing.T) {
	preTestChecks(t)

	skipIfNotSysAdmin(t)

	var params = StringMap{
		"Testname": t.Name(),

		"Tags": "tm",
	}
	testParamsNotEmpty(t, params)

	configText1 := templateFill(testAccVcdTmOrgSubproviderStep1, params)
	params["FuncName"] = t.Name() + "-step2"
	configText2 := templateFill(testAccVcdTmOrgSubproviderStep2, params)

	debugPrintf("#[DEBUG] CONFIGURATION step1: %s\n", configText1)
	debugPrintf("#[DEBUG] CONFIGURATION step2: %s\n", configText2)
	if vcfaShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configText1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vcfa_org.test", "name", t.Name()),
					resource.TestCheckResourceAttr("vcfa_org.test", "display_name", "terraform-test"),
					resource.TestCheckResourceAttr("vcfa_org.test", "description", "terraform test"),
					resource.TestCheckResourceAttr("vcfa_org.test", "is_enabled", "true"),
					resource.TestCheckResourceAttr("vcfa_org.test", "is_subprovider", "true"),
					resource.TestMatchResourceAttr("vcfa_org.test", "managed_by_id", regexp.MustCompile("^urn:vcloud:org:")),
					resource.TestCheckResourceAttr("vcfa_org.test", "managed_by_name", "System"),
					resource.TestCheckResourceAttr("vcfa_org.test", "is_classic_tenant", "false"),
				),
			},
			{
				Config: configText2,
				Check: resource.ComposeTestCheckFunc(
					resourceFieldsEqual("vcfa_org.test", "data.vcfa_org.test", nil),
				),
			},
			{
				ResourceName:      "vcfa_org.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     params["Testname"].(string),
			},
		},
	})

	postTestChecks(t)
}

const testAccVcdTmOrgSubproviderStep1 = `
resource "vcfa_org" "test" {
  name           = "{{.Testname}}"
  display_name   = "terraform-test"
  description    = "terraform test"
  is_enabled     = true
  is_subprovider = true
}
`

const testAccVcdTmOrgSubproviderStep2 = testAccVcdTmOrgSubproviderStep1 + `
data "vcfa_org" "test" {
  name = vcfa_org.test.name
}
`

// TestAccVcdTmOrgClassicTenant tests a Tenant Manager Organization configured as "Classic Tenant"
func TestAccVcdTmOrgClassicTenant(t *testing.T) {
	preTestChecks(t)
	skipIfNotSysAdmin(t)

	var params = StringMap{
		"Testname": t.Name(),
		"Tags":     "tm",
	}
	testParamsNotEmpty(t, params)

	configText1 := templateFill(testAccVcdTmOrgClassicStep1, params)
	params["FuncName"] = t.Name() + "-step2"
	configText2 := templateFill(testAccVcdTmOrgClassicStep2, params)

	debugPrintf("#[DEBUG] CONFIGURATION step1: %s\n", configText1)
	debugPrintf("#[DEBUG] CONFIGURATION step2: %s\n", configText2)
	if vcfaShortTest {
		t.Skip(acceptanceTestsSkipped)
		return
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configText1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vcfa_org.test", "name", t.Name()),
					resource.TestCheckResourceAttr("vcfa_org.test", "display_name", "terraform-test"),
					resource.TestCheckResourceAttr("vcfa_org.test", "description", "terraform test"),
					resource.TestCheckResourceAttr("vcfa_org.test", "is_enabled", "true"),
					resource.TestCheckResourceAttr("vcfa_org.test", "is_subprovider", "false"),
					resource.TestMatchResourceAttr("vcfa_org.test", "managed_by_id", regexp.MustCompile("^urn:vcloud:org:")),
					resource.TestCheckResourceAttr("vcfa_org.test", "managed_by_name", "System"),
					resource.TestCheckResourceAttr("vcfa_org.test", "is_classic_tenant", "true"),
				),
			},
			{
				Config: configText2,
				Check: resource.ComposeTestCheckFunc(
					resourceFieldsEqual("vcfa_org.test", "data.vcfa_org.test", nil),
				),
			},
			{
				ResourceName:      "vcfa_org.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     params["Testname"].(string),
			},
		},
	})

	postTestChecks(t)
}

const testAccVcdTmOrgClassicStep1 = `
resource "vcfa_org" "test" {
  name              = "{{.Testname}}"
  display_name      = "terraform-test"
  description       = "terraform test"
  is_enabled        = true
  is_classic_tenant = true
}

resource "vcfa_org" "test2" {
  name              = "{{.Testname}}2"
  display_name      = "terraform-test"
  description       = "terraform test"
  is_enabled        = true
  is_classic_tenant = true
}
`

const testAccVcdTmOrgClassicStep2 = testAccVcdTmOrgClassicStep1 + `
data "vcfa_org" "test" {
  name = vcfa_org.test.name
}
`