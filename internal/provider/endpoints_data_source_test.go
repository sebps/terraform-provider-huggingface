package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEndpointsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "huggingface_endpoints" "test" { 
					namespace = "<YOUR_NAMESPACE>" 
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of coffees returned
					resource.TestCheckResourceAttr("data.huggingface_endpoints.test", "endpoints.#", "1"),
					// Verify the first coffee to ensure all attributes are set
					resource.TestCheckResourceAttr("data.huggingface_endpoints.test", "endpoints.0.name", "test-terraform-0"),
					resource.TestCheckResourceAttr("data.huggingface_endpoints.test", "endpoints.0.type", "protected"),
					resource.TestCheckResourceAttr("data.huggingface_endpoints.test", "endpoints.0.cloud_provider.vendor", "aws"),
					resource.TestCheckResourceAttr("data.huggingface_endpoints.test", "endpoints.0.cloud_provider.region", "us-east-1"),
					resource.TestCheckResourceAttr("data.huggingface_endpoints.test", "endpoints.0.compute.accelerator", "cpu"),
				),
			},
		},
	})
}
