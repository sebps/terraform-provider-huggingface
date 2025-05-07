package provider

import (
	"context"
	"fmt"

	"github.com/sebps/terraform-provider-huggingface/internal/models"
	"github.com/sebps/terraform-provider-huggingface/internal/states"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	huggingface "github.com/sebps/huggingface-client/client"
)

func NewEndpointsDataSource() datasource.DataSource {
	return &endpointsDataSource{}
}

type endpointsDataSource struct {
	client *huggingface.Client
}

func (d *endpointsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_endpoints"
}

func (d *endpointsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"namespace": schema.StringAttribute{
				Required: true,
			},
			"endpoints": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"namespace": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
							Optional: true,
						},
						"cloud_provider": schema.SingleNestedAttribute{
							Computed: true,
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"vendor": schema.StringAttribute{
									Computed: true,
									Optional: true,
								},
								"region": schema.StringAttribute{
									Computed: true,
									Optional: true,
								},
							},
						},
						"compute": schema.SingleNestedAttribute{
							Computed: true,
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"accelerator": schema.StringAttribute{
									Computed: true,
									Optional: true,
								},
								"id": schema.StringAttribute{
									Computed: true,
									Optional: true,
								},
								// "instance_type": schema.StringAttribute{
								// 	Computed: true,
								// },
								// "instance_size": schema.StringAttribute{
								// 	Computed: true,
								// },
								// "scaling": schema.MapNestedAttribute{
								// 	Computed: true,
								// 	NestedObject: schema.NestedAttributeObject{
								// 		Attributes: map[string]schema.Attribute{
								// 			"min_replica": schema.Int32Attribute{
								// 				Computed: true,
								// 			},
								// 			"max_replica": schema.Int32Attribute{
								// 				Computed: true,
								// 			},
								// 			"measure": schema.MapNestedAttribute{
								// 				Computed: true,
								// 				NestedObject: schema.NestedAttributeObject{
								// 					Attributes: map[string]schema.Attribute{
								// 						"hardware_usage": schema.Float64Attribute{
								// 							Computed: true,
								// 						},
								// 						"pending_requests": schema.Float64Attribute{
								// 							Computed: true,
								// 						},
								// 					},
								// 				},
								// 			},
								// 			"metric": schema.StringAttribute{
								// 				Computed: true,
								// 			},
								// 			"scale_to_zero_timeout": schema.Int32Attribute{
								// 				Computed: true,
								// 			},
								// 			"threshold": schema.Float64Attribute{
								// 				Computed: true,
								// 			},
								// 		},
								// 	},
								// },
							},
						},
						"model": schema.SingleNestedAttribute{
							Computed: true,
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"repository": schema.StringAttribute{
									Computed: true,
									Optional: true,
								},
								"framework": schema.StringAttribute{
									Computed: true,
									Optional: true,
								},
								"task": schema.StringAttribute{
									Computed: true,
									Optional: true,
								},
								// "image": schema.SingleNestedAttribute{
								// 	Computed: true,
								// 	Attributes: map[string]schema.Attribute{
								// 		"huggingface": schema.SingleNestedAttribute{
								// 			Computed:   true,
								// 			Attributes: map[string]schema.Attribute{},
								// 		},
								// 		"huggingface_neuron": schema.SingleNestedAttribute{
								// 			Computed: true,
								// 			Attributes: map[string]schema.Attribute{
								// 				"batch_size": schema.Int32Attribute{
								// 					Computed: true,
								// 				},
								// 				"neuron_cache": schema.StringAttribute{
								// 					Computed: true,
								// 				},
								// 				"sequence_length": schema.Int32Attribute{
								// 					Computed: true,
								// 				},
								// 			},
								// 		},
								// 		"tgi": schema.SingleNestedAttribute{
								// 			Computed: true,
								// 			Attributes: map[string]schema.Attribute{
								// 				"health_route": schema.StringAttribute{
								// 					Computed: true,
								// 				},
								// 				"port": schema.Int32Attribute{
								// 					Computed: true,
								// 				},
								// 				"url": schema.StringAttribute{
								// 					Computed: true,
								// 				},
								// 				"max_batch_prefill_tokens": schema.Int32Attribute{
								// 					Computed: true,
								// 				},
								// 				"max_batch_total_tokens": schema.Int32Attribute{
								// 					Computed: true,
								// 				},
								// 				"max_input_length": schema.Int32Attribute{
								// 					Computed: true,
								// 				},
								// 				"max_total_tokens": schema.Int32Attribute{
								// 					Computed: true,
								// 				},
								// 				"disable_custom_kernels": schema.BoolAttribute{
								// 					Computed: true,
								// 				},
								// 				"quantize": schema.StringAttribute{
								// 					Computed: true,
								// 				},
								// 			},
								// 		},
								// 		"tgi_neuron": schema.SingleNestedAttribute{
								// 			Computed: true,
								// 			Attributes: map[string]schema.Attribute{
								// 				"health_route": schema.StringAttribute{
								// 					Computed: true,
								// 				},
								// 				"port": schema.Int32Attribute{
								// 					Computed: true,
								// 				},
								// 				"url": schema.StringAttribute{
								// 					Computed: true,
								// 				},
								// 				"max_batch_prefill_tokens": schema.Int32Attribute{
								// 					Computed: true,
								// 				},
								// 				"max_batch_total_tokens": schema.Int32Attribute{
								// 					Computed: true,
								// 				},
								// 				"max_input_length": schema.Int32Attribute{
								// 					Computed: true,
								// 				},
								// 				"max_total_tokens": schema.Int32Attribute{
								// 					Computed: true,
								// 				},
								// 				"hf_auto_cast_type": schema.StringAttribute{
								// 					Computed: true,
								// 				},
								// 				"hf_num_cores": schema.Int32Attribute{
								// 					Computed: true,
								// 				},
								// 			},
								// 		},
								// 		"tei": schema.SingleNestedAttribute{
								// 			Computed: true,
								// 			Attributes: map[string]schema.Attribute{
								// 				"health_route": schema.StringAttribute{
								// 					Computed: true,
								// 				},
								// 				"port": schema.Int32Attribute{
								// 					Computed: true,
								// 				},
								// 				"url": schema.StringAttribute{
								// 					Computed: true,
								// 				},
								// 				"max_batch_tokens": schema.Int32Attribute{
								// 					Computed: true,
								// 				},
								// 				"max_concurrent_requests": schema.Int32Attribute{
								// 					Computed: true,
								// 				},
								// 				"pooling": schema.StringAttribute{
								// 					Computed: true,
								// 				},
								// 			},
								// 		},
								// 		"llamacpp": schema.SingleNestedAttribute{
								// 			Computed: true,
								// 			Attributes: map[string]schema.Attribute{
								// 				"health_route": schema.StringAttribute{
								// 					Computed: true,
								// 				},
								// 				"port": schema.Int32Attribute{
								// 					Computed: true,
								// 				},
								// 				"url": schema.StringAttribute{
								// 					Computed: true,
								// 				},
								// 				"ctx_size": schema.Int32Attribute{
								// 					Computed: true,
								// 				},
								// 				"mode": schema.StringAttribute{
								// 					Computed: true,
								// 				},
								// 				"model_path": schema.StringAttribute{
								// 					Computed: true,
								// 				},
								// 				"n_gpu_layers": schema.Int32Attribute{
								// 					Computed: true,
								// 				},
								// 				"n_parallel": schema.Int32Attribute{
								// 					Computed: true,
								// 				},
								// 				"pooling": schema.StringAttribute{
								// 					Computed: true,
								// 				},
								// 				"threads_http": schema.Int32Attribute{
								// 					Computed: true,
								// 				},
								// 				"variant": schema.StringAttribute{
								// 					Computed: true,
								// 				},
								// 			},
								// 		},
								// 		"custom": schema.SingleNestedAttribute{
								// 			Computed: true,
								// 			Attributes: map[string]schema.Attribute{
								// 				"url": schema.StringAttribute{
								// 					Computed: true,
								// 				},
								// 				"health_route": schema.StringAttribute{
								// 					Computed: true,
								// 				},
								// 				"port": schema.Int32Attribute{
								// 					Computed: true,
								// 				},
								// 				"credentials": schema.MapNestedAttribute{
								// 					Computed: true,
								// 					NestedObject: schema.NestedAttributeObject{
								// 						Attributes: map[string]schema.Attribute{
								// 							"username": schema.StringAttribute{
								// 								Computed: true,
								// 							},
								// 							"password": schema.StringAttribute{
								// 								Computed: true,
								// 							},
								// 						},
								// 					},
								// 				},
								// 			},
								// 		},
								// 	},
								// },
								// "instance_size": schema.StringAttribute{
								// 	Computed: true,
								// },
								// "scaling": schema.SingleNestedAttribute{
								// 	Computed: true,
								// 	Attributes: map[string]schema.Attribute{
								// 		"min_replica": schema.Int32Attribute{
								// 			Computed: true,
								// 		},
								// 		"max_replica": schema.Int32Attribute{
								// 			Computed: true,
								// 		},
								// 		"measure": schema.SingleNestedAttribute{
								// 			Computed: true,
								// 			Attributes: map[string]schema.Attribute{
								// 				"hardware_usage": schema.Float64Attribute{
								// 					Computed: true,
								// 				},
								// 				"pending_requests": schema.Float64Attribute{
								// 					Computed: true,
								// 				},
								// 			},
								// 		},
								// 		"metric": schema.StringAttribute{
								// 			Computed: true,
								// 		},
								// 		"scale_to_zero_timeout": schema.Int32Attribute{
								// 			Computed: true,
								// 		},
								// 		"threshold": schema.Float64Attribute{
								// 			Computed: true,
								// 		},
								// 	},
								// },
							},
						},
						"tags": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Optional:    true,
						},
						"cache_http_responses": schema.BoolAttribute{
							Computed: true,
							Optional: true,
						},
						"experimental_features": schema.SingleNestedAttribute{
							Computed: true,
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"cache_http_responses": schema.BoolAttribute{
									Computed: true,
									Optional: true,
								},
								"kv_router": schema.SingleNestedAttribute{
									Computed: true,
									Optional: true,
									Attributes: map[string]schema.Attribute{
										"tag": schema.StringAttribute{
											Computed: true,
											Optional: true,
										},
									},
								},
							},
						},
						"private_service": schema.SingleNestedAttribute{
							Computed: true,
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"account_id": schema.StringAttribute{
									Computed: true,
									Optional: true,
								},
								"shared": schema.BoolAttribute{
									Computed: true,
									Optional: true,
								},
							},
						},
						"route": schema.SingleNestedAttribute{
							Computed: true,
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"domain": schema.StringAttribute{
									Computed: true,
									Optional: true,
								},
								"path": schema.StringAttribute{
									Computed: true,
									Optional: true,
								},
							},
						},
						"status": schema.SingleNestedAttribute{
							Computed: true,
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"created_at": schema.StringAttribute{
									Computed: true,
									Optional: true,
								},
								"created_by": schema.SingleNestedAttribute{
									Computed: true,
									Optional: true,
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Computed: true,
											Optional: true,
										},
										"name": schema.StringAttribute{
											Computed: true,
											Optional: true,
										},
									},
								},
								"updated_at": schema.StringAttribute{
									Computed: true,
									Optional: true,
								},
								"updated_by": schema.SingleNestedAttribute{
									Computed: true,
									Optional: true,
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Computed: true,
											Optional: true,
										},
										"name": schema.StringAttribute{
											Computed: true,
											Optional: true,
										},
									},
								},
								"state": schema.StringAttribute{
									Computed: true,
									Optional: true,
								},
								"message": schema.StringAttribute{
									Computed: true,
									Optional: true,
								},
								"ready_replica": schema.NumberAttribute{
									Computed: true,
									Optional: true,
								},
								"target_replica": schema.NumberAttribute{
									Computed: true,
									Optional: true,
								},
								"error_message": schema.StringAttribute{
									Computed: true,
									Optional: true,
								},
								"url": schema.StringAttribute{
									Computed: true,
									Optional: true,
								},
								"private": schema.SingleNestedAttribute{
									Computed: true,
									Optional: true,
									Attributes: map[string]schema.Attribute{
										"service_name": schema.StringAttribute{
											Computed: true,
											Optional: true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

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
		endpointState := models.Endpoint{
			Namespace: state.Namespace,
			Name:      types.StringValue(endpoint.Name),
			Type:      types.StringValue(string(endpoint.Type)),
		}

		// 	CloudProvider        types.Object        `tfsdk:"cloud_provider"`
		endpointCloudProvider := models.EndpointCloudProvider{
			Vendor: types.StringValue(endpoint.Provider.Vendor),
			Region: types.StringValue(endpoint.Provider.Region),
		}
		endpointState.CloudProvider, diags = types.ObjectValueFrom(ctx, endpointCloudProvider.AttributeTypes(), endpointCloudProvider)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Compute              types.Object        `tfsdk:"compute"`
		endpointCompute := models.EndpointCompute{
			ID:          types.StringValue(*endpoint.Compute.ID),
			Accelerator: types.StringValue(string(endpoint.Compute.Accelerator)),
		}
		endpointState.Compute, diags = types.ObjectValueFrom(ctx, endpointCompute.AttributeTypes(), endpointCompute)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Model                types.Object        `tfsdk:"model"`
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

		// Tags                 basetypes.ListValue `tfsdk:"tags"`
		endpointState.Tags, diags = types.ListValueFrom(ctx, types.StringType, endpoint.Tags)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		if endpoint.CacheHttpResponses != nil {
			endpointState.CacheHttpResponses = types.BoolValue(*endpoint.CacheHttpResponses)
		}

		// ExperimentalFeatures types.Object        `tfsdk:"experimental_features"`
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

		// PrivateService       types.Object        `tfsdk:"private_service"`
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

		// Route                types.Object        `tfsdk:"route"`
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

		// Status               types.Object        `tfsdk:"status"`
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

		// Add deeper nested mapping (compute, model, experimental_features, etc.) as needed
		state.Endpoints = append(state.Endpoints, endpointState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Configure adds the provider configured client to the data source.
func (d *endpointsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	fmt.Println("into endpointsDataSource.Configure")

	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*huggingface.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *huggingface.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}
