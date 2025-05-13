package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sebps/terraform-provider-huggingface/internal/states"
	"github.com/sebps/terraform-provider-huggingface/internal/transformers"
)

// Update updates the resource and sets the updated Terraform state on success.
func (r *endpointsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan states.EndpointResourceState
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id, name, namespace string
	if !plan.ID.IsNull() && !plan.ID.IsUnknown() {
		// extract id
		id = plan.ID.ValueString()
		idParts := strings.Split(id, "/")

		// extract name and namespace from id
		namespace = idParts[0]
		name = idParts[1]
	} else {
		// extract name
		name = plan.Name.ValueString()

		// extract namespace
		namespace = plan.Namespace.ValueString()
	}

	// Define endpoint to create from plan
	endpointToUpdate := transformers.FromPlanToEndpointUpdate(ctx, &plan)

	// Update endpoint
	endpointUpdated, err := r.client.UpdateEndpoint(namespace, name, endpointToUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating endpoint",
			"Could not update endpoint, unexpected error: "+err.Error(),
		)
		return
	}

	// Map back plan from created endpoint
	updatedPlan, diags := transformers.FromProviderToModel(ctx, endpointUpdated)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// inject namespace
	updatedPlan.Namespace = types.StringValue(namespace)

	// inject id
	id = fmt.Sprintf("%s/%s", namespace, name)
	updatedPlan.ID = types.StringValue(id)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, updatedPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
