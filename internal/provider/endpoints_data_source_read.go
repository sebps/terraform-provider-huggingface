package provider

import (
	"context"

	"github.com/sebps/terraform-provider-huggingface/internal/states"
	"github.com/sebps/terraform-provider-huggingface/internal/transformers"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Read refreshes the Terraform state with the latest data.
func (d *endpointsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state states.EndpointsDataSourceState

	// Get configuration into the model
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	namespace := state.Namespace.ValueString()
	endpoints, err := d.client.ListEndpoints(namespace, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Huggingface Endpoints",
			"Could not Read Huggingface Endpoints for namespace "+state.Namespace.ValueString()+" : "+err.Error(),
		)
		return
	}

	for _, endpoint := range endpoints {
		endpointState, diags := transformers.FromProviderToModel(ctx, &endpoint)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Add deeper nested mapping (compute, model, experimental_features, etc.) as needed
		state.Endpoints = append(state.Endpoints, endpointState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
