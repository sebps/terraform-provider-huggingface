package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/sebps/terraform-provider-huggingface/internal/states"
	"github.com/sebps/terraform-provider-huggingface/internal/transformers"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &endpointsResource{}
	_ resource.ResourceWithConfigure = &endpointsResource{}
)

// Read resource information.
func (r *endpointsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var plan states.EndpointResourceState
	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id, name, namespace string
	if !plan.ID.IsNull() && !plan.ID.IsUnknown() {
		// extract id
		id = plan.ID.ValueString()
		idParts := strings.Split(id, "/")

		tflog.Info(ctx, "ID received", map[string]any{"id": id})

		// extract name and namespace from id
		namespace = idParts[0]
		name = idParts[1]
	} else {
		tflog.Info(ctx, "ID not received")

		// extract name
		name = plan.Name.ValueString()

		// extract namespace
		namespace = plan.Namespace.ValueString()
	}

	// Get refreshed endpoint value from Huggingface
	endpoint, err := r.client.GetEndpoint(namespace, name)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Huggingface Endpoint",
			"Could not read Huggingface Endpoint "+plan.Namespace.ValueString()+"/"+plan.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	updatedPlan, diags := transformers.FromProviderToModel(ctx, endpoint)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// inject namespace
	updatedPlan.Namespace = types.StringValue(namespace)

	// inject id
	id = fmt.Sprintf("%s/%s", namespace, name)
	updatedPlan.ID = types.StringValue(id)

	// Set refreshed state
	diags = resp.State.Set(ctx, &updatedPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
