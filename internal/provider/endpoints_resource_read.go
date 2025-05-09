package provider

import (
	"context"

	"github.com/sebps/terraform-provider-huggingface/internal/states"

	"github.com/sebps/terraform-provider-huggingface/internal/models"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &endpointsResource{}
	_ resource.ResourceWithConfigure = &endpointsResource{}
)

// Read resource information.
func (r *endpointsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var endpointState states.EndpointResourceState
	diags := req.State.Get(ctx, &endpointState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed endpoint value from Huggingface
	endpoint, err := r.client.GetEndpoint(endpointState.Namespace.ValueString(), endpointState.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Huggingface Endpoint",
			"Could not read Huggingface Endpoint "+endpointState.Namespace.ValueString()+"/"+endpointState.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	// 	Cloud Provider
	endpointCloudProvider := models.EndpointCloudProvider{
		Vendor: types.StringValue(endpoint.Provider.Vendor),
		Region: types.StringValue(endpoint.Provider.Region),
	}
	endpointState.CloudProvider, diags = types.ObjectValueFrom(ctx, endpointCloudProvider.AttributeTypes(), endpointCloudProvider)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Compute
	endpointCompute := models.EndpointCompute{
		ID:          types.StringValue(*endpoint.Compute.ID),
		Accelerator: types.StringValue(string(endpoint.Compute.Accelerator)),
	}
	endpointState.Compute, diags = types.ObjectValueFrom(ctx, endpointCompute.AttributeTypes(), endpointCompute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Model
	endpointModel := models.Model{
		Repository: types.StringValue(endpoint.Model.Repository),
		Framework:  types.StringValue(string(endpoint.Model.Framework)),
		Task:       types.StringValue(string(endpoint.Model.Task)),
	}
	endpointState.Model, diags = types.ObjectValueFrom(ctx, endpointModel.AttributeTypes(), endpointModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Tags
	endpointState.Tags, diags = types.ListValueFrom(ctx, types.StringType, endpoint.Tags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if endpoint.CacheHttpResponses != nil {
		endpointState.CacheHttpResponses = types.BoolValue(*endpoint.CacheHttpResponses)
	} else {
		endpointState.CacheHttpResponses = types.BoolValue(false)
	}

	// Experimental Features
	var experimentalFeatures models.ExperimentalFeatures
	if endpoint.ExperimentalFeatures != nil {
		experimentalFeatures = models.ExperimentalFeatures{
			CacheHTTPResponses: types.BoolValue(endpoint.ExperimentalFeatures.CacheHttpResponses),
		}

		var kvRouter models.KvRouter
		if endpoint.ExperimentalFeatures.KvRouter != nil {
			kvRouter = models.KvRouter{
				Tag: types.StringValue(endpoint.ExperimentalFeatures.KvRouter.Tag),
			}
		} else {
			kvRouter = models.KvRouter{
				Tag: types.StringValue(""),
			}
		}
		experimentalFeatures.KVRouter, diags = types.ObjectValueFrom(ctx, kvRouter.AttributeTypes(), kvRouter)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		experimentalFeatures = models.ExperimentalFeatures{
			CacheHTTPResponses: types.BoolValue(false),
		}

		kvRouter := models.KvRouter{
			Tag: types.StringValue(""),
		}
		experimentalFeatures.KVRouter, diags = types.ObjectValueFrom(ctx, kvRouter.AttributeTypes(), kvRouter)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	endpointState.ExperimentalFeatures, diags = types.ObjectValueFrom(ctx, experimentalFeatures.AttributeTypes(), experimentalFeatures)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Private Service
	var endpointPrivateService models.PrivateService
	if endpoint.PrivateService != nil {
		endpointPrivateService = models.PrivateService{
			AccountID: types.StringValue(endpoint.PrivateService.AccountID),
			Shared:    types.BoolValue(endpoint.PrivateService.Shared),
		}
	} else {
		endpointPrivateService = models.PrivateService{
			AccountID: types.StringValue(""),
			Shared:    types.BoolValue(false),
		}
	}
	endpointState.PrivateService, diags = types.ObjectValueFrom(ctx, endpointPrivateService.AttributeTypes(), endpointPrivateService)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Route
	var endpointRoute models.Route
	if endpoint.Route != nil {
		endpointRoute = models.Route{
			Domain: types.StringValue(endpoint.Route.Domain),
			Path:   types.StringValue(endpoint.Route.Path),
		}
	} else {
		endpointRoute = models.Route{
			Domain: types.StringValue(""),
			Path:   types.StringValue(""),
		}
	}
	endpointState.Route, diags = types.ObjectValueFrom(ctx, endpointRoute.AttributeTypes(), endpointRoute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Status
	endpointStatus := models.Status{
		CreatedAt:     types.StringValue(endpoint.Status.CreatedAt.String()),
		UpdatedAt:     types.StringValue(endpoint.Status.UpdatedAt.String()),
		State:         types.StringValue(string(endpoint.Status.State)),
		Message:       types.StringValue(endpoint.Status.Message),
		ReadyReplica:  types.Int32Value(int32(endpoint.Status.ReadyReplica)),
		TargetReplica: types.Int32Value(int32(endpoint.Status.TargetReplica)),
	}
	endpointStatusCreatedBy := models.User{
		Id:   types.StringValue(endpoint.Status.CreatedBy.ID),
		Name: types.StringValue(endpoint.Status.CreatedBy.Name),
	}
	endpointStatus.CreatedBy, diags = types.ObjectValueFrom(ctx, endpointStatusCreatedBy.AttributeTypes(), endpointStatusCreatedBy)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpointStatusUpdatedBy := models.User{
		Id:   types.StringValue(endpoint.Status.UpdatedBy.ID),
		Name: types.StringValue(endpoint.Status.UpdatedBy.Name),
	}
	endpointStatus.UpdatedBy, diags = types.ObjectValueFrom(ctx, endpointStatusUpdatedBy.AttributeTypes(), endpointStatusUpdatedBy)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if endpoint.Status.ErrorMessage != nil {
		endpointStatus.ErrorMessage = types.StringValue(*endpoint.Status.ErrorMessage)
	} else {
		endpointStatus.ErrorMessage = types.StringValue("")
	}

	if endpoint.Status.URL != nil {
		endpointStatus.Url = types.StringValue(*endpoint.Status.URL)
	} else {
		endpointStatus.Url = types.StringValue("")
	}

	var endpointStatusPrivate models.Private
	if endpoint.Status.Private != nil && endpoint.Status.Private.ServiceName != nil {
		endpointStatusPrivate = models.Private{
			ServiceName: types.StringValue(*endpoint.Status.Private.ServiceName),
		}
	} else {
		endpointStatusPrivate = models.Private{
			ServiceName: types.StringValue(""),
		}
	}
	endpointStatus.Private, diags = types.ObjectValueFrom(ctx, endpointStatusPrivate.AttributeTypes(), endpointStatusPrivate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpointState.Status, diags = types.ObjectValueFrom(ctx, endpointStatus.AttributeTypes(), endpointStatus)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &endpointState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
