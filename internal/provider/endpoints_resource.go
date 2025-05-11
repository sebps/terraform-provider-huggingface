package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	huggingface "github.com/sebps/huggingface-client/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &endpointsResource{}
	_ resource.ResourceWithConfigure = &endpointsResource{}
)

func NewEndpointsResource() resource.Resource {
	return &endpointsResource{}
}

type endpointsResource struct {
	client *huggingface.Client
}

// Metadata returns the resource type name.
func (r *endpointsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_endpoint"
}

// Schema defines the schema for the resource.
func (r *endpointsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"namespace": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"type": schema.StringAttribute{
				Required: true,
			},
			"cloud_provider": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"vendor": schema.StringAttribute{
						Required: true,
					},
					"region": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"compute": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"accelerator": schema.StringAttribute{
						Required: true,
					},
					"id": schema.StringAttribute{
						Computed: true,
						Optional: true,
					},
					"instance_type": schema.StringAttribute{
						Required: true,
					},
					"instance_size": schema.StringAttribute{
						Required: true,
					},
					"scaling": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"min_replica": schema.Int32Attribute{
								Required: true,
							},
							"max_replica": schema.Int32Attribute{
								Required: true,
							},
							"measure": schema.SingleNestedAttribute{
								Required: true,
								Attributes: map[string]schema.Attribute{
									"hardware_usage": schema.Float64Attribute{
										Computed: true,
										Optional: true,
									},
									"pending_requests": schema.Float64Attribute{
										Optional: true,
									},
								},
							},
							"metric": schema.StringAttribute{
								Computed: true,
								Optional: true,
							},
							"scale_to_zero_timeout": schema.Int32Attribute{
								Computed: true,
								Optional: true,
							},
							"threshold": schema.Float64Attribute{
								Computed: true,
								Optional: true,
							},
						},
					},
				},
			},
			"model": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"repository": schema.StringAttribute{
						Required: true,
					},
					"framework": schema.StringAttribute{
						Required: true,
					},
					"task": schema.StringAttribute{
						Required: true,
					},
					"image": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"huggingface": schema.SingleNestedAttribute{
								Computed:   true,
								Optional:   true,
								Attributes: map[string]schema.Attribute{},
							},
							"huggingface_neuron": schema.SingleNestedAttribute{
								Computed: true,
								Optional: true,
								Attributes: map[string]schema.Attribute{
									"batch_size": schema.Int32Attribute{
										Computed: true,
										Optional: true,
									},
									"neuron_cache": schema.StringAttribute{
										Computed: true,
										Optional: true,
									},
									"sequence_length": schema.Int32Attribute{
										Computed: true,
										Optional: true,
									},
								},
							},
							"tgi": schema.SingleNestedAttribute{
								Computed: true,
								Optional: true,
								Attributes: map[string]schema.Attribute{
									"health_route": schema.StringAttribute{
										Computed: true,
										Optional: true,
									},
									"port": schema.Int32Attribute{
										Computed: true,
										Optional: true,
									},
									"url": schema.StringAttribute{
										Required: true,
									},
									"max_batch_prefill_tokens": schema.Int32Attribute{
										Computed: true,
										Optional: true,
									},
									"max_batch_total_tokens": schema.Int32Attribute{
										Computed: true,
										Optional: true,
									},
									"max_input_length": schema.Int32Attribute{
										Computed: true,
										Optional: true,
									},
									"max_total_tokens": schema.Int32Attribute{
										Computed: true,
										Optional: true,
									},
									"disable_custom_kernels": schema.BoolAttribute{
										Computed: true,
										Optional: true,
									},
									"quantize": schema.StringAttribute{
										Computed: true,
										Optional: true,
									},
								},
							},
							"tgi_neuron": schema.SingleNestedAttribute{
								Computed: true,
								Optional: true,
								Attributes: map[string]schema.Attribute{
									"health_route": schema.StringAttribute{
										Computed: true,
										Optional: true,
									},
									"port": schema.Int32Attribute{
										Computed: true,
										Optional: true,
									},
									"url": schema.StringAttribute{
										Required: true,
									},
									"max_batch_prefill_tokens": schema.Int32Attribute{
										Computed: true,
										Optional: true,
									},
									"max_batch_total_tokens": schema.Int32Attribute{
										Computed: true,
										Optional: true,
									},
									"max_input_length": schema.Int32Attribute{
										Computed: true,
										Optional: true,
									},
									"max_total_tokens": schema.Int32Attribute{
										Computed: true,
										Optional: true,
									},
									"hf_auto_cast_type": schema.StringAttribute{
										Computed: true,
										Optional: true,
									},
									"hf_num_cores": schema.Int32Attribute{
										Computed: true,
										Optional: true,
									},
								},
							},
							"tei": schema.SingleNestedAttribute{
								Computed: true,
								Optional: true,
								Attributes: map[string]schema.Attribute{
									"health_route": schema.StringAttribute{
										Computed: true,
										Optional: true,
									},
									"port": schema.Int32Attribute{
										Computed: true,
										Optional: true,
									},
									"url": schema.StringAttribute{
										Required: true,
									},
									"max_batch_tokens": schema.Int32Attribute{
										Computed: true,
										Optional: true,
									},
									"max_concurrent_requests": schema.Int32Attribute{
										Computed: true,
										Optional: true,
									},
									"pooling": schema.StringAttribute{
										Computed: true,
										Optional: true,
									},
								},
							},
							"llamacpp": schema.SingleNestedAttribute{
								Computed: true,
								Optional: true,
								Attributes: map[string]schema.Attribute{
									"health_route": schema.StringAttribute{
										Computed: true,
										Optional: true,
									},
									"port": schema.Int32Attribute{
										Computed: true,
										Optional: true,
									},
									"url": schema.StringAttribute{
										Required: true,
									},
									"ctx_size": schema.Int32Attribute{
										Required: true,
									},
									"mode": schema.StringAttribute{
										Computed: true,
										Optional: true,
									},
									"model_path": schema.StringAttribute{
										Required: true,
									},
									"n_gpu_layers": schema.Int32Attribute{
										Computed: true,
										Optional: true,
									},
									"n_parallel": schema.Int32Attribute{
										Required: true,
									},
									"pooling": schema.StringAttribute{
										Computed: true,
										Optional: true,
									},
									"threads_http": schema.Int32Attribute{
										Required: true,
									},
									"variant": schema.StringAttribute{
										Computed: true,
										Optional: true,
									},
								},
							},
							"custom": schema.SingleNestedAttribute{
								Computed: true,
								Optional: true,
								Attributes: map[string]schema.Attribute{
									"url": schema.StringAttribute{
										Required: true,
									},
									"health_route": schema.StringAttribute{
										Computed: true,
									},
									"port": schema.Int32Attribute{
										Computed: true,
									},
									"credentials": schema.SingleNestedAttribute{
										Computed: true,
										Attributes: map[string]schema.Attribute{
											"username": schema.StringAttribute{
												Required: true,
											},
											"password": schema.StringAttribute{
												Required: true,
											},
										},
									},
								},
							},
						},
					},
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
				Attributes: map[string]schema.Attribute{
					"created_at": schema.StringAttribute{
						Computed: true,
					},
					"created_by": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Computed: true,
							},
							"name": schema.StringAttribute{
								Computed: true,
							},
						},
					},
					"updated_at": schema.StringAttribute{
						Computed: true,
					},
					"updated_by": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Computed: true,
							},
							"name": schema.StringAttribute{
								Computed: true,
							},
						},
					},
					"state": schema.StringAttribute{
						Computed: true,
					},
					"message": schema.StringAttribute{
						Computed: true,
					},
					"ready_replica": schema.NumberAttribute{
						Computed: true,
					},
					"target_replica": schema.NumberAttribute{
						Computed: true,
					},
					"error_message": schema.StringAttribute{
						Computed: true,
					},
					"url": schema.StringAttribute{
						Computed: true,
					},
					"private": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"service_name": schema.StringAttribute{
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *endpointsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}
