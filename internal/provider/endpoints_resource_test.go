package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEndpointsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
					resource "huggingface_endpoint" "test" {
						namespace = "<YOUR_NAMESPACE>"
						name      = "test-terraform-1"
						type      = "protected"

						compute = {
							accelerator   = "cpu"
							instance_type = "intel-icl"
							instance_size = "x4"
							scaling = {
								min_replica           = 0
								max_replica           = 1
								metric                = "hardwareUsage"

								measure = {
									hardware_usage = 10
								}
							}
						}

						model = {
							framework  = "pytorch"
							repository = "openai-community/gpt2"
							task       = "text-generation"
							image = {
								huggingface = {}
							}
						}

						cloud_provider = {
							region = "us-east-1"
							vendor = "aws"
						}					
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("huggingface_endpoint.test", "name", "test-terraform-1"),
					resource.TestCheckResourceAttr("huggingface_endpoint.test", "type", "protected"),
					resource.TestCheckResourceAttr("huggingface_endpoint.test", "cloud_provider.vendor", "aws"),
					resource.TestCheckResourceAttr("huggingface_endpoint.test", "cloud_provider.region", "us-east-1"),
					resource.TestCheckResourceAttr("huggingface_endpoint.test", "compute.scaling.min_replica", "0"),
					resource.TestCheckResourceAttr("huggingface_endpoint.test", "compute.scaling.max_replica", "1"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "huggingface_endpoint.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does not exist in the HashiCups
				// API, therefore there is no value for it during import.
				ImportStateVerifyIgnore: []string{"status"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
					resource "huggingface_endpoint" "test" {
						namespace = "<YOUR_NAMESPACE>"
						name      = "test-terraform-1"
						type      = "protected"

						compute = {
							accelerator   = "cpu"
							instance_type = "intel-icl"
							instance_size = "x4"
							scaling = {
								min_replica           = 0
								max_replica           = 2
								metric                = "hardwareUsage"

								measure = {
									hardware_usage = 10
								}
							}
						}

						model = {
							framework  = "pytorch"
							repository = "openai-community/gpt2"
							task       = "text-generation"
							image = {
								huggingface = {}
							}
						}

						cloud_provider = {
							region = "us-east-1"
							vendor = "aws"
						}
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("huggingface_endpoint.test", "name", "test-terraform-1"),
					resource.TestCheckResourceAttr("huggingface_endpoint.test", "type", "protected"),
					resource.TestCheckResourceAttr("huggingface_endpoint.test", "cloud_provider.vendor", "aws"),
					resource.TestCheckResourceAttr("huggingface_endpoint.test", "cloud_provider.region", "us-east-1"),
					resource.TestCheckResourceAttr("huggingface_endpoint.test", "compute.scaling.min_replica", "0"),
					resource.TestCheckResourceAttr("huggingface_endpoint.test", "compute.scaling.max_replica", "2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
