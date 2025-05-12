package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/sebps/terraform-provider-huggingface/internal/states"
	"github.com/sebps/terraform-provider-huggingface/internal/transformers"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Create creates the resource and sets the initial Terraform state.
func (r *endpointsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
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
	} else {
		// extract namespace
		namespace = plan.Namespace.ValueString()
	}

	// Define endpoint to create from plan
	endpointToCreate := transformers.FromPlanToEndpoint(ctx, &plan)

	// Create new endpoint
	endpointCreated, err := r.client.CreateEndpoint(namespace, endpointToCreate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating endpoint",
			"Could not create endpoint, unexpected error: "+err.Error(),
		)
		return
	}

	// Map back plan from updated endpoint
	updatedPlan, diags := transformers.FromEndpointToPlan(ctx, endpointCreated)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// inject namespace
	updatedPlan.Namespace = types.StringValue(namespace)

	// inject name
	name = endpointCreated.Name

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
