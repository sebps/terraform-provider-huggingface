package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Update updates the resource and sets the updated Terraform state on success.
func (r *endpointsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}
