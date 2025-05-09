package provider

import (
	"context"
	"fmt"

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
