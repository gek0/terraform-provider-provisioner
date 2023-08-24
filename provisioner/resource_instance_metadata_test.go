package provisioner

import (
	"fmt"
	"testing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccProvisionerInstanceMetadataBasic(t *testing.T) {
	ad_domain := "acme.local"
	instance_name := "vm-gpu-loc-1"

	resource.Test(t, resource.TestCase {
		PreCheck:		func() { testAccPreCheck(t) },
		Providers:		testAccProviders,
		CheckDestroy:	testAccCheckProvisionerInstanceMetadataDestroy,
		Steps:			[]resource.TestStep {
			{
				Config: testAccCheckProvisionerInstanceMetadataConfigBasic(ad_domain, instance_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProvisionerInstanceMetadataExists("provisioner_instance_metadata.new"),
					resource.TestCheckResourceAttr(
						"provisioner_instance_metadata.new", "ad_domain", ad_domain,
					),
					resource.TestCheckResourceAttr(
						"provisioner_instance_metadata.new", "instance_name", instance_name,
					),
				),
			},
		},
	})
}

func testAccCheckProvisionerInstanceMetadataConfigBasic(ad_domain string, instance_name string) string {
	return fmt.Sprintf(`
	resource "provisioner_instance_metadata" "new" {
		ad_domain     = %s
		instance_name = %s
	}
	`, ad_domain, instance_name)
}

func testAccCheckProvisionerInstanceMetadataExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No InstanceMetadataID set")
		}

		return nil
	}
}

func testAccCheckProvisionerInstanceMetadataDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "provisioner_instance_metadata" {
			continue
		}

		// get ID from state since upstream API does not manage it
		if rs.Primary.ID != "" {
			return fmt.Errorf("InstanceMetadataID stll exists")
		}
	}

	return nil
}