package provider

import (
	"context"
	"fmt"

	"github.com/sebps/terraform-provider-huggingface/internal/states"

	"github.com/sebps/terraform-provider-huggingface/internal/models"

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
						Computed: true,
					},
					"instance_size": schema.StringAttribute{
						Computed: true,
					},
					"scaling": schema.MapNestedAttribute{
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"min_replica": schema.Int32Attribute{
									Computed: true,
								},
								"max_replica": schema.Int32Attribute{
									Computed: true,
								},
								"measure": schema.MapNestedAttribute{
									Computed: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"hardware_usage": schema.Float64Attribute{
												Computed: true,
											},
											"pending_requests": schema.Float64Attribute{
												Computed: true,
											},
										},
									},
								},
								"metric": schema.StringAttribute{
									Computed: true,
								},
								"scale_to_zero_timeout": schema.Int32Attribute{
									Computed: true,
								},
								"threshold": schema.Float64Attribute{
									Computed: true,
								},
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
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"huggingface": schema.SingleNestedAttribute{
								Computed:   true,
								Attributes: map[string]schema.Attribute{},
							},
							"huggingface_neuron": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"batch_size": schema.Int32Attribute{
										Computed: true,
									},
									"neuron_cache": schema.StringAttribute{
										Computed: true,
									},
									"sequence_length": schema.Int32Attribute{
										Computed: true,
									},
								},
							},
							"tgi": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"health_route": schema.StringAttribute{
										Computed: true,
									},
									"port": schema.Int32Attribute{
										Computed: true,
									},
									"url": schema.StringAttribute{
										Computed: true,
									},
									"max_batch_prefill_tokens": schema.Int32Attribute{
										Computed: true,
									},
									"max_batch_total_tokens": schema.Int32Attribute{
										Computed: true,
									},
									"max_input_length": schema.Int32Attribute{
										Computed: true,
									},
									"max_total_tokens": schema.Int32Attribute{
										Computed: true,
									},
									"disable_custom_kernels": schema.BoolAttribute{
										Computed: true,
									},
									"quantize": schema.StringAttribute{
										Computed: true,
									},
								},
							},
							"tgi_neuron": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"health_route": schema.StringAttribute{
										Computed: true,
									},
									"port": schema.Int32Attribute{
										Computed: true,
									},
									"url": schema.StringAttribute{
										Computed: true,
									},
									"max_batch_prefill_tokens": schema.Int32Attribute{
										Computed: true,
									},
									"max_batch_total_tokens": schema.Int32Attribute{
										Computed: true,
									},
									"max_input_length": schema.Int32Attribute{
										Computed: true,
									},
									"max_total_tokens": schema.Int32Attribute{
										Computed: true,
									},
									"hf_auto_cast_type": schema.StringAttribute{
										Computed: true,
									},
									"hf_num_cores": schema.Int32Attribute{
										Computed: true,
									},
								},
							},
							"tei": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"health_route": schema.StringAttribute{
										Computed: true,
									},
									"port": schema.Int32Attribute{
										Computed: true,
									},
									"url": schema.StringAttribute{
										Computed: true,
									},
									"max_batch_tokens": schema.Int32Attribute{
										Computed: true,
									},
									"max_concurrent_requests": schema.Int32Attribute{
										Computed: true,
									},
									"pooling": schema.StringAttribute{
										Computed: true,
									},
								},
							},
							"llamacpp": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"health_route": schema.StringAttribute{
										Computed: true,
									},
									"port": schema.Int32Attribute{
										Computed: true,
									},
									"url": schema.StringAttribute{
										Computed: true,
									},
									"ctx_size": schema.Int32Attribute{
										Computed: true,
									},
									"mode": schema.StringAttribute{
										Computed: true,
									},
									"model_path": schema.StringAttribute{
										Computed: true,
									},
									"n_gpu_layers": schema.Int32Attribute{
										Computed: true,
									},
									"n_parallel": schema.Int32Attribute{
										Computed: true,
									},
									"pooling": schema.StringAttribute{
										Computed: true,
									},
									"threads_http": schema.Int32Attribute{
										Computed: true,
									},
									"variant": schema.StringAttribute{
										Computed: true,
									},
								},
							},
							"custom": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"url": schema.StringAttribute{
										Computed: true,
									},
									"health_route": schema.StringAttribute{
										Computed: true,
									},
									"port": schema.Int32Attribute{
										Computed: true,
									},
									"credentials": schema.MapNestedAttribute{
										Computed: true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"username": schema.StringAttribute{
													Computed: true,
												},
												"password": schema.StringAttribute{
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
					},
					"scaling": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"min_replica": schema.Int32Attribute{
								Computed: true,
							},
							"max_replica": schema.Int32Attribute{
								Computed: true,
							},
							"measure": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"hardware_usage": schema.Float64Attribute{
										Computed: true,
									},
									"pending_requests": schema.Float64Attribute{
										Computed: true,
									},
								},
							},
							"metric": schema.StringAttribute{
								Computed: true,
							},
							"scale_to_zero_timeout": schema.Int32Attribute{
								Computed: true,
							},
							"threshold": schema.Float64Attribute{
								Computed: true,
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
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *endpointsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan states.EndpointResourceState
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Define new endpoint
	endpointToCreate := huggingface.Endpoint{
		Name: plan.Name.ValueString(),
		Type: huggingface.EndpointType(plan.Type.ValueString()),
		Compute: huggingface.EndpointCompute{
			Scaling: huggingface.EndpointScaling{},
		},
		Model: huggingface.EndpointModel{
			Image: huggingface.EndpointModelImage{
				HuggingFace:       &huggingface.HuggingFaceImage{},
				HuggingFaceNeuron: &huggingface.HuggingFaceNeuronImage{},
				TGI:               &huggingface.TGIImage{},
				TGINeuron:         &huggingface.TGINeuronImage{},
				TEI:               &huggingface.TEIImage{},
				LlamaCpp:          &huggingface.LlamaCppImage{},
				Custom:            &huggingface.CustomImage{},
			},
			Task: huggingface.EndpointTask("text-generation"),
		},
		Provider: huggingface.EndpointProvider{},
	}

	// Compute
	computeAttributes := plan.Compute.Attributes()
	if accelerator, ok := computeAttributes["accelerator"]; ok {
		tfAccelerator, _ := accelerator.ToTerraformValue(ctx)
		tfAccelerator.As(&endpointToCreate.Compute.Accelerator)
	}
	if instanceType, ok := computeAttributes["instance_type"]; ok {
		tfInstanceType, _ := instanceType.ToTerraformValue(ctx)
		tfInstanceType.As(&endpointToCreate.Compute.InstanceType)
	}
	if instanceSize, ok := computeAttributes["instance_size"]; ok {
		tfInstanceType, _ := instanceSize.ToTerraformValue(ctx)
		tfInstanceType.As(&endpointToCreate.Compute.InstanceSize)
	}
	if scaling, ok := computeAttributes["scaling"]; ok {
		tfScaling, _ := scaling.ToTerraformValue(ctx)
		scalingModel := types.Object{}
		tfScaling.As(&scalingModel)

		scalingAttributes := scalingModel.Attributes()
		if minReplica, ok := scalingAttributes["min_replica"]; ok {
			tfMinReplica, _ := minReplica.ToTerraformValue(ctx)
			tfMinReplica.As(&endpointToCreate.Compute.Scaling.MinReplica)
		} else {
			endpointToCreate.Compute.Scaling.MinReplica = 0
		}
		if maxReplica, ok := scalingAttributes["max_replica"]; ok {
			tfMaxReplica, _ := maxReplica.ToTerraformValue(ctx)
			tfMaxReplica.As(&endpointToCreate.Compute.Scaling.MaxReplica)
		} else {
			endpointToCreate.Compute.Scaling.MaxReplica = 0
		}
	}

	// Model
	modelAttributes := plan.Model.Attributes()
	if repository, ok := modelAttributes["repository"]; ok {
		tfRepository, _ := repository.ToTerraformValue(ctx)
		tfRepository.As(&endpointToCreate.Model.Repository)
	}
	if framework, ok := modelAttributes["framework"]; ok {
		tfFramework, _ := framework.ToTerraformValue(ctx)
		tfFramework.As(&endpointToCreate.Model.Framework)
	}
	if task, ok := modelAttributes["task"]; ok {
		tfTask, _ := task.ToTerraformValue(ctx)
		tfTask.As(&endpointToCreate.Model.Task)
	}
	if image, ok := computeAttributes["image"]; ok {
		tfImage, _ := image.ToTerraformValue(ctx)
		imageModel := types.Object{}
		tfImage.As(&imageModel)

		imageAttributes := imageModel.Attributes()
		if huggingFace, ok := imageAttributes["huggingface"]; ok {
			tfHuggingface, _ := huggingFace.ToTerraformValue(ctx)
			tfHuggingface.As(&endpointToCreate.Model.Image.HuggingFace)
		}
		if huggingfaceNeuron, ok := imageAttributes["huggingface_neuron"]; ok {
			tfHuggingfaceNeuron, _ := huggingfaceNeuron.ToTerraformValue(ctx)
			tfHuggingfaceNeuronModel := types.Object{}
			tfHuggingfaceNeuron.As(tfHuggingfaceNeuronModel)

			huggingfaceNeuronAttributes := tfHuggingfaceNeuronModel.Attributes()
			if batchSize, ok := huggingfaceNeuronAttributes["batch_size"]; ok {
				tfBatchSize, _ := batchSize.ToTerraformValue(ctx)
				tfBatchSize.As(&endpointToCreate.Model.Image.HuggingFaceNeuron.BatchSize)
			}
			if neuronCache, ok := huggingfaceNeuronAttributes["neuron_cache"]; ok {
				tfNeuronCache, _ := neuronCache.ToTerraformValue(ctx)
				tfNeuronCache.As(&endpointToCreate.Model.Image.HuggingFaceNeuron.NeuronCache)
			}
			if sequenceLength, ok := huggingfaceNeuronAttributes["sequence_length"]; ok {
				tfSequenceLength, _ := sequenceLength.ToTerraformValue(ctx)
				tfSequenceLength.As(&endpointToCreate.Model.Image.HuggingFaceNeuron.SequenceLength)
			}
		}
		if tgi, ok := imageAttributes["tgi"]; ok {
			tfTgi, _ := tgi.ToTerraformValue(ctx)
			tfTgiModel := types.Object{}
			tfTgi.As(tfTgiModel)

			tgiAttributes := tfTgiModel.Attributes()
			if disableCustomKernels, ok := tgiAttributes["disable_custom_kernels"]; ok {
				tfDisableCustomKernels, _ := disableCustomKernels.ToTerraformValue(ctx)
				tfDisableCustomKernels.As(&endpointToCreate.Model.Image.TGI.DisableCustomKernels)
			}
			if healthRoute, ok := tgiAttributes["health_route"]; ok {
				tfHealthRoute, _ := healthRoute.ToTerraformValue(ctx)
				tfHealthRoute.As(&endpointToCreate.Model.Image.TGI.HealthRoute)
			}
			if maxBatchPrefillTokens, ok := tgiAttributes["max_batch_prefill_tokens"]; ok {
				tfMaxBatchPrefillTokens, _ := maxBatchPrefillTokens.ToTerraformValue(ctx)
				tfMaxBatchPrefillTokens.As(&endpointToCreate.Model.Image.TGI.MaxBatchPrefillTokens)
			}
			if maxBatchTotalTokens, ok := tgiAttributes["max_batch_total_tokens"]; ok {
				tfMaxBatchTotalTokens, _ := maxBatchTotalTokens.ToTerraformValue(ctx)
				tfMaxBatchTotalTokens.As(&endpointToCreate.Model.Image.TGI.MaxBatchTotalTokens)
			}
			if maxInputLength, ok := tgiAttributes["max_input_length"]; ok {
				tfMaxInputLength, _ := maxInputLength.ToTerraformValue(ctx)
				tfMaxInputLength.As(&endpointToCreate.Model.Image.TGI.MaxInputLength)
			}
			if maxTotalTokens, ok := tgiAttributes["max_total_tokens"]; ok {
				tfMaxTotalTokens, _ := maxTotalTokens.ToTerraformValue(ctx)
				tfMaxTotalTokens.As(&endpointToCreate.Model.Image.TGI.MaxTotalTokens)
			}
			if port, ok := tgiAttributes["port"]; ok {
				tfPort, _ := port.ToTerraformValue(ctx)
				tfPort.As(&endpointToCreate.Model.Image.TGI.Port)
			}
			if quantize, ok := tgiAttributes["quantize"]; ok {
				tfQuantize, _ := quantize.ToTerraformValue(ctx)
				tfQuantize.As(&endpointToCreate.Model.Image.TGI.Quantize)
			}
			if url, ok := tgiAttributes["url"]; ok {
				tfUrl, _ := url.ToTerraformValue(ctx)
				tfUrl.As(&endpointToCreate.Model.Image.TGI.URL)
			}
		}
		if tgiNeuron, ok := imageAttributes["tgi_neuron"]; ok {
			tfTgiNeuron, _ := tgiNeuron.ToTerraformValue(ctx)
			tfTgiNeuronModel := types.Object{}
			tfTgiNeuron.As(tfTgiNeuronModel)

			tgiNeuronAttributes := tfTgiNeuronModel.Attributes()
			if healthRoute, ok := tgiNeuronAttributes["health_route"]; ok {
				tfHealthRoute, _ := healthRoute.ToTerraformValue(ctx)
				tfHealthRoute.As(&endpointToCreate.Model.Image.TGINeuron.HealthRoute)
			}
			if hfAutoCastType, ok := tgiNeuronAttributes["hf_auto_cast_type"]; ok {
				tfHfAutoCastType, _ := hfAutoCastType.ToTerraformValue(ctx)
				tfHfAutoCastType.As(&endpointToCreate.Model.Image.TGINeuron.HfAutoCastType)
			}
			if hfNumCores, ok := tgiNeuronAttributes["hf_num_cores"]; ok {
				tfHfNumCores, _ := hfNumCores.ToTerraformValue(ctx)
				tfHfNumCores.As(&endpointToCreate.Model.Image.TGINeuron.HfNumCores)
			}
			if maxBatchPrefillTokens, ok := tgiNeuronAttributes["max_batch_prefill_tokens"]; ok {
				tfMaxBatchPrefillTokens, _ := maxBatchPrefillTokens.ToTerraformValue(ctx)
				tfMaxBatchPrefillTokens.As(&endpointToCreate.Model.Image.TGINeuron.MaxBatchPrefillTokens)
			}
			if maxBatchTotalTokens, ok := tgiNeuronAttributes["max_batch_total_tokens"]; ok {
				tfMaxBatchTotalTokens, _ := maxBatchTotalTokens.ToTerraformValue(ctx)
				tfMaxBatchTotalTokens.As(&endpointToCreate.Model.Image.TGINeuron.MaxBatchTotalTokens)
			}
			if maxInputLength, ok := tgiNeuronAttributes["max_input_length"]; ok {
				tfMaxInputLength, _ := maxInputLength.ToTerraformValue(ctx)
				tfMaxInputLength.As(&endpointToCreate.Model.Image.TGINeuron.MaxInputLength)
			}
			if maxTotalTokens, ok := tgiNeuronAttributes["max_total_tokens"]; ok {
				tfMaxTotalTokens, _ := maxTotalTokens.ToTerraformValue(ctx)
				tfMaxTotalTokens.As(&endpointToCreate.Model.Image.TGINeuron.MaxTotalTokens)
			}
			if port, ok := tgiNeuronAttributes["port"]; ok {
				tfPort, _ := port.ToTerraformValue(ctx)
				tfPort.As(&endpointToCreate.Model.Image.TGINeuron.Port)
			}
			if url, ok := tgiNeuronAttributes["url"]; ok {
				tfUrl, _ := url.ToTerraformValue(ctx)
				tfUrl.As(&endpointToCreate.Model.Image.TGINeuron.URL)
			}
		}
		if teiNeuron, ok := imageAttributes["tei_neuron"]; ok {
			tfTeiNeuron, _ := teiNeuron.ToTerraformValue(ctx)
			tfTeiNeuronModel := types.Object{}
			tfTeiNeuron.As(tfTeiNeuronModel)

			teiNeuronAttributes := tfTeiNeuronModel.Attributes()
			if healthRoute, ok := teiNeuronAttributes["health_route"]; ok {
				tfHealthRoute, _ := healthRoute.ToTerraformValue(ctx)
				tfHealthRoute.As(&endpointToCreate.Model.Image.TEI.HealthRoute)
			}
			if maxBatchTokens, ok := teiNeuronAttributes["max_batch_tokens"]; ok {
				tfMaxBatchTokens, _ := maxBatchTokens.ToTerraformValue(ctx)
				tfMaxBatchTokens.As(&endpointToCreate.Model.Image.TEI.MaxBatchTokens)
			}
			if maxConcurrentRequests, ok := teiNeuronAttributes["max_concurrent_requests"]; ok {
				tfMaxConcurrentRequests, _ := maxConcurrentRequests.ToTerraformValue(ctx)
				tfMaxConcurrentRequests.As(&endpointToCreate.Model.Image.TEI.MaxConcurrentRequests)
			}
			if pooling, ok := teiNeuronAttributes["pooling"]; ok {
				tfPooling, _ := pooling.ToTerraformValue(ctx)
				tfPooling.As(&endpointToCreate.Model.Image.TEI.Pooling)
			}
			if port, ok := teiNeuronAttributes["port"]; ok {
				tfPort, _ := port.ToTerraformValue(ctx)
				tfPort.As(&endpointToCreate.Model.Image.TEI.Port)
			}
			if url, ok := teiNeuronAttributes["url"]; ok {
				tfUrl, _ := url.ToTerraformValue(ctx)
				tfUrl.As(&endpointToCreate.Model.Image.TEI.URL)
			}
		}
		if llamacpp, ok := imageAttributes["llamacpp"]; ok {
			tfLlamacpp, _ := llamacpp.ToTerraformValue(ctx)
			tfLlamacppModel := types.Object{}
			tfLlamacpp.As(tfLlamacppModel)

			llamacppAttributes := tfLlamacppModel.Attributes()
			if ctxSize, ok := llamacppAttributes["ctx_size"]; ok {
				tfCtxSize, _ := ctxSize.ToTerraformValue(ctx)
				tfCtxSize.As(&endpointToCreate.Model.Image.LlamaCpp.CtxSize)
			}
			if healthRoute, ok := llamacppAttributes["health_route"]; ok {
				tfHealthRoute, _ := healthRoute.ToTerraformValue(ctx)
				tfHealthRoute.As(&endpointToCreate.Model.Image.LlamaCpp.HealthRoute)
			}
			if mode, ok := llamacppAttributes["mode"]; ok {
				tfMode, _ := mode.ToTerraformValue(ctx)
				tfMode.As(&endpointToCreate.Model.Image.LlamaCpp.Mode)
			}
			if modelPath, ok := llamacppAttributes["model_path"]; ok {
				tfModelPath, _ := modelPath.ToTerraformValue(ctx)
				tfModelPath.As(&endpointToCreate.Model.Image.LlamaCpp.ModelPath)
			}
			if nGpuLayers, ok := llamacppAttributes["n_gpu_layers"]; ok {
				tfNGpuLayers, _ := nGpuLayers.ToTerraformValue(ctx)
				tfNGpuLayers.As(&endpointToCreate.Model.Image.LlamaCpp.NGpuLayers)
			}
			if nParallel, ok := llamacppAttributes["n_parallel"]; ok {
				tfNParallel, _ := nParallel.ToTerraformValue(ctx)
				tfNParallel.As(&endpointToCreate.Model.Image.LlamaCpp.NParallel)
			}
			if pooling, ok := llamacppAttributes["pooling"]; ok {
				tfPooling, _ := pooling.ToTerraformValue(ctx)
				tfPooling.As(&endpointToCreate.Model.Image.LlamaCpp.Pooling)
			}
			if port, ok := llamacppAttributes["port"]; ok {
				tfPort, _ := port.ToTerraformValue(ctx)
				tfPort.As(&endpointToCreate.Model.Image.LlamaCpp.Port)
			}
			if threadsHttp, ok := llamacppAttributes["threads_http"]; ok {
				tfThreadsHttp, _ := threadsHttp.ToTerraformValue(ctx)
				tfThreadsHttp.As(&endpointToCreate.Model.Image.LlamaCpp.ThreadsHttp)
			}
			if url, ok := llamacppAttributes["url"]; ok {
				tfUrl, _ := url.ToTerraformValue(ctx)
				tfUrl.As(&endpointToCreate.Model.Image.LlamaCpp.URL)
			}
			if variant, ok := llamacppAttributes["variant"]; ok {
				tfVariant, _ := variant.ToTerraformValue(ctx)
				tfVariant.As(&endpointToCreate.Model.Image.LlamaCpp.Variant)
			}
		}
		if custom, ok := imageAttributes["custom"]; ok {
			tfCustom, _ := custom.ToTerraformValue(ctx)
			tfCustomModel := types.Object{}
			tfCustom.As(tfCustomModel)

			customAttributes := tfCustomModel.Attributes()
			if credentials, ok := customAttributes["credentials"]; ok {
				tfCustomCredentials, _ := credentials.ToTerraformValue(ctx)
				tfCustomCredentialsModel := types.Object{}
				tfCustomCredentials.As(tfCustomCredentialsModel)

				customCredentialsAttributes := tfCustomCredentialsModel.Attributes()
				if username, ok := customCredentialsAttributes["username"]; ok {
					tfUsername, _ := username.ToTerraformValue(ctx)
					tfUsername.As(&endpointToCreate.Model.Image.Custom.Credentials.Username)
				}
				if password, ok := customCredentialsAttributes["password"]; ok {
					tfPassword, _ := password.ToTerraformValue(ctx)
					tfPassword.As(&endpointToCreate.Model.Image.Custom.Credentials.Password)
				}
			}
			if healthRoute, ok := customAttributes["health_route"]; ok {
				tfHealthRoute, _ := healthRoute.ToTerraformValue(ctx)
				tfHealthRoute.As(&endpointToCreate.Model.Image.Custom.HealthRoute)
			}
			if port, ok := customAttributes["port"]; ok {
				tfPort, _ := port.ToTerraformValue(ctx)
				tfPort.As(&endpointToCreate.Model.Image.Custom.Port)
			}
			if url, ok := customAttributes["url"]; ok {
				tfUrl, _ := url.ToTerraformValue(ctx)
				tfUrl.As(&endpointToCreate.Model.Image.Custom.URL)
			}
		}
	}

	// Cloud Provider
	cloudProviderAttributes := plan.CloudProvider.Attributes()
	if vendor, ok := cloudProviderAttributes["vendor"]; ok {
		tfVendor, _ := vendor.ToTerraformValue(ctx)
		tfVendor.As(&endpointToCreate.Provider.Vendor)
	}
	if region, ok := cloudProviderAttributes["region"]; ok {
		tfRegion, _ := region.ToTerraformValue(ctx)
		tfRegion.As(&endpointToCreate.Provider.Region)
	}

	// Create new endpoint
	endpoint, err := r.client.CreateEndpoint(plan.Namespace.ValueString(), endpointToCreate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating endpoint",
			"Could not create endpoint, unexpected error: "+err.Error(),
		)
		return
	}

	// Map values back to state

	// 	Cloud Provider
	endpointCloudProvider := models.EndpointCloudProvider{
		Vendor: types.StringValue(endpoint.Provider.Vendor),
		Region: types.StringValue(endpoint.Provider.Region),
	}
	plan.CloudProvider, diags = types.ObjectValueFrom(ctx, endpointCloudProvider.AttributeTypes(), endpointCloudProvider)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Compute
	endpointCompute := models.EndpointCompute{
		ID:          types.StringValue(*endpoint.Compute.ID),
		Accelerator: types.StringValue(string(endpoint.Compute.Accelerator)),
	}
	plan.Compute, diags = types.ObjectValueFrom(ctx, endpointCompute.AttributeTypes(), endpointCompute)
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
	plan.Model, diags = types.ObjectValueFrom(ctx, endpointModel.AttributeTypes(), endpointModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Tags
	plan.Tags, diags = types.ListValueFrom(ctx, types.StringType, endpoint.Tags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if endpoint.CacheHttpResponses != nil {
		plan.CacheHttpResponses = types.BoolValue(*endpoint.CacheHttpResponses)
	} else {
		plan.CacheHttpResponses = types.BoolValue(false)
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
	plan.ExperimentalFeatures, diags = types.ObjectValueFrom(ctx, experimentalFeatures.AttributeTypes(), experimentalFeatures)
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
	plan.PrivateService, diags = types.ObjectValueFrom(ctx, endpointPrivateService.AttributeTypes(), endpointPrivateService)
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
	plan.Route, diags = types.ObjectValueFrom(ctx, endpointRoute.AttributeTypes(), endpointRoute)
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

	plan.Status, diags = types.ObjectValueFrom(ctx, endpointStatus.AttributeTypes(), endpointStatus)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

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

// Update updates the resource and sets the updated Terraform state on success.
func (r *endpointsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *endpointsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
