package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/sebps/terraform-provider-huggingface/internal/states"
)

// Delete deletes the resource and removes the Terraform state on success.
func (r *endpointsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from plan
	var state states.EndpointResourceState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id, name, namespace string
	if !state.ID.IsNull() {
		// extract id
		id = state.ID.ValueString()
		idParts := strings.Split(id, "/")

		// extract name and namespace from id
		namespace = idParts[0]
		name = idParts[1]
	} else {
		// extract name
		name = state.Name.ValueString()

		// extract namespace
		namespace = state.Namespace.ValueString()
	}

	// Delete endpoint
	err := r.client.DeleteEndpoint(namespace, name)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting endpoint",
			"Could not delete endpoint, unexpected error: "+err.Error(),
		)
		return
	}
}
