package provider

import (
	"context"
	"math/big"

	"github.com/sebps/terraform-provider-huggingface/internal/states"

	"github.com/sebps/terraform-provider-huggingface/internal/models"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	huggingface "github.com/sebps/huggingface-client/client"
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

	// Define new endpoint
	endpointToCreate := huggingface.Endpoint{
		Name: plan.Name.ValueString(),
		Type: huggingface.EndpointType(plan.Type.ValueString()),
		Compute: huggingface.EndpointCompute{
			Scaling: huggingface.EndpointScaling{
				Measure: &huggingface.ScalingMeasure{},
			},
		},
		Model: huggingface.EndpointModel{
			Image: huggingface.EndpointModelImage{},
			Task:  huggingface.EndpointTask("text-generation"),
		},
		Provider: huggingface.EndpointProvider{},
	}

	// Cloud Provider
	cloudProviderAttributes := plan.CloudProvider.Attributes()
	if vendor, ok := cloudProviderAttributes["vendor"]; ok && !vendor.IsNull() {
		tfVendor, _ := vendor.ToTerraformValue(ctx)
		tfVendor.As(&endpointToCreate.Provider.Vendor)
	}
	if region, ok := cloudProviderAttributes["region"]; ok && !region.IsNull() {
		tfRegion, _ := region.ToTerraformValue(ctx)
		tfRegion.As(&endpointToCreate.Provider.Region)
	}

	// Compute
	if !plan.Compute.IsNull() {
		computeAttributes := plan.Compute.Attributes()
		if accelerator, ok := computeAttributes["accelerator"]; ok && !accelerator.IsNull() {
			tfAccelerator, _ := accelerator.ToTerraformValue(ctx)
			var hfAccelerator string
			tfAccelerator.As(&hfAccelerator)
			endpointToCreate.Compute.Accelerator = huggingface.AcceleratorType(hfAccelerator)
		}
		if instanceType, ok := computeAttributes["instance_type"]; ok && !instanceType.IsNull() {
			tfInstanceType, _ := instanceType.ToTerraformValue(ctx)
			tfInstanceType.As(&endpointToCreate.Compute.InstanceType)
		}
		if instanceSize, ok := computeAttributes["instance_size"]; ok && !instanceSize.IsNull() {
			tfInstanceType, _ := instanceSize.ToTerraformValue(ctx)
			tfInstanceType.As(&endpointToCreate.Compute.InstanceSize)
		}
		if scaling, ok := computeAttributes["scaling"]; ok && !scaling.IsNull() {
			tfScaling, _ := scaling.ToTerraformValue(ctx)
			var scalingAttributes map[string]tftypes.Value
			tfScaling.As(&scalingAttributes)

			if tfMinReplica, ok := scalingAttributes["min_replica"]; ok && !tfMinReplica.IsNull() {
				var minReplicaBigFloat big.Float
				tfMinReplica.As(&minReplicaBigFloat)
				maxReplicaInt, _ := minReplicaBigFloat.Int(nil)
				endpointToCreate.Compute.Scaling.MinReplica = int(maxReplicaInt.Int64())
			} else {
				endpointToCreate.Compute.Scaling.MinReplica = 0
			}
			if tfMaxReplica, ok := scalingAttributes["max_replica"]; ok && !tfMaxReplica.IsNull() {
				var maxReplicaBigFloat big.Float
				tfMaxReplica.As(&maxReplicaBigFloat)
				maxReplicaInt, _ := maxReplicaBigFloat.Int(nil)
				endpointToCreate.Compute.Scaling.MaxReplica = int(maxReplicaInt.Int64())
			} else {
				endpointToCreate.Compute.Scaling.MaxReplica = 1
			}
			if tfMetric, ok := scalingAttributes["metric"]; ok && !tfMetric.IsNull() {
				tfMetric.As(&endpointToCreate.Compute.Scaling.Metric)
			}
			if tfScaleToZeroTimeout, ok := scalingAttributes["scale_to_zero_timeout"]; ok && !tfScaleToZeroTimeout.IsNull() && tfScaleToZeroTimeout.IsKnown() {
				var scaleToZeroTimeoutBigFloat big.Float
				tfScaleToZeroTimeout.As(&scaleToZeroTimeoutBigFloat)
				scaleToZeroTimeoutInt, _ := scaleToZeroTimeoutBigFloat.Int(nil)
				scaleToZeroTimeoutIntPrimitive := int(scaleToZeroTimeoutInt.Int64())
				endpointToCreate.Compute.Scaling.ScaleToZeroTimeout = &scaleToZeroTimeoutIntPrimitive
			}
			if tfThreshold, ok := scalingAttributes["threshold"]; ok && !tfThreshold.IsNull() && tfThreshold.IsKnown() {
				var thresholdBigFloat big.Float
				tfThreshold.As(&thresholdBigFloat)
				thresholdFloatPrimitive, _ := thresholdBigFloat.Float64()
				endpointToCreate.Compute.Scaling.Threshold = &thresholdFloatPrimitive
			}
			if tfMeasure, ok := scalingAttributes["measure"]; ok && !tfMeasure.IsNull() {
				var measureAttributes map[string]tftypes.Value
				tfMeasure.As(&measureAttributes)

				if tfHardwareUsage, ok := measureAttributes["hardware_usage"]; ok && !tfHardwareUsage.IsNull() {
					var hardwareUsageBigFloat big.Float
					tfHardwareUsage.As(&hardwareUsageBigFloat)
					hardwareUsageFloatPrimitive, _ := hardwareUsageBigFloat.Float64()
					endpointToCreate.Compute.Scaling.Measure.HardwareUsage = &hardwareUsageFloatPrimitive
				}
				if tfPendingRequests, ok := measureAttributes["pending_requests"]; ok && !tfPendingRequests.IsNull() {
					var pendingRequestsBigFloat big.Float
					tfPendingRequests.As(&pendingRequestsBigFloat)
					pendintRequestsFloatPrimitive, _ := pendingRequestsBigFloat.Float64()
					endpointToCreate.Compute.Scaling.Measure.PendingRequests = &pendintRequestsFloatPrimitive
				}
			}
		}
	}

	// Model
	if !plan.Model.IsNull() {
		modelAttributes := plan.Model.Attributes()
		if repository, ok := modelAttributes["repository"]; ok && !repository.IsNull() {
			tfRepository, _ := repository.ToTerraformValue(ctx)
			tfRepository.As(&endpointToCreate.Model.Repository)
		}
		if framework, ok := modelAttributes["framework"]; ok && !framework.IsNull() {
			tfFramework, _ := framework.ToTerraformValue(ctx)
			var hfFramework string
			tfFramework.As(&hfFramework)
			endpointToCreate.Model.Framework = huggingface.EndpointFramework(hfFramework)
		}
		if task, ok := modelAttributes["task"]; ok && !task.IsNull() {
			tfTask, _ := task.ToTerraformValue(ctx)
			tfTask.As(&endpointToCreate.Model.Task)
		}
		if image, ok := modelAttributes["image"]; ok && !image.IsNull() {
			tfImage, _ := image.ToTerraformValue(ctx)
			var imageAttributes map[string]tftypes.Value
			tfImage.As(&imageAttributes)

			if tfHuggingface, ok := imageAttributes["huggingface"]; ok && !tfHuggingface.IsNull() && tfHuggingface.IsKnown() {
				endpointToCreate.Model.Image.HuggingFace = &huggingface.HuggingFaceImage{}
			}
			if tfHuggingfaceNeuron, ok := imageAttributes["huggingface_neuron"]; ok && !tfHuggingfaceNeuron.IsNull() && tfHuggingfaceNeuron.IsKnown() {
				endpointToCreate.Model.Image.HuggingFaceNeuron = &huggingface.HuggingFaceNeuronImage{}
				var huggingfaceNeuronAttributes map[string]tftypes.Value
				tfHuggingfaceNeuron.As(&huggingfaceNeuronAttributes)

				if tfBatchSize, ok := huggingfaceNeuronAttributes["batch_size"]; ok && !tfBatchSize.IsNull() {
					var batchSizeBigFloat big.Float
					tfBatchSize.As(&batchSizeBigFloat)
					batchSizeInt, _ := batchSizeBigFloat.Int(nil)
					batchSizeIntPrimitive := int(batchSizeInt.Int64())
					endpointToCreate.Model.Image.HuggingFaceNeuron.BatchSize = &batchSizeIntPrimitive
				}
				if tfNeuronCache, ok := huggingfaceNeuronAttributes["neuron_cache"]; ok && !tfNeuronCache.IsNull() {
					tfNeuronCache.As(&endpointToCreate.Model.Image.HuggingFaceNeuron.NeuronCache)
				}
				if tfSequenceLength, ok := huggingfaceNeuronAttributes["sequence_length"]; ok && !tfSequenceLength.IsNull() {
					var sequenceLengthBigFloat big.Float
					tfSequenceLength.As(&sequenceLengthBigFloat)
					sequenceLengthInt, _ := sequenceLengthBigFloat.Int(nil)
					sequenceLengthIntPrimitive := int(sequenceLengthInt.Int64())
					endpointToCreate.Model.Image.HuggingFaceNeuron.SequenceLength = &sequenceLengthIntPrimitive
				}
			}
			if tfTgi, ok := imageAttributes["tgi"]; ok && !tfTgi.IsNull() && tfTgi.IsKnown() {
				endpointToCreate.Model.Image.TGI = &huggingface.TGIImage{}
				var tgiAttributes map[string]tftypes.Value
				tfTgi.As(&tgiAttributes)

				if tfDisableCustomKernels, ok := tgiAttributes["disable_custom_kernels"]; ok && !tfDisableCustomKernels.IsNull() {
					tfDisableCustomKernels.As(&endpointToCreate.Model.Image.TGI.DisableCustomKernels)
				}
				if tfHealthRoute, ok := tgiAttributes["health_route"]; ok && !tfHealthRoute.IsNull() {
					tfHealthRoute.As(&endpointToCreate.Model.Image.TGI.HealthRoute)
				}
				if tfMaxBatchPrefillTokens, ok := tgiAttributes["max_batch_prefill_tokens"]; ok && !tfMaxBatchPrefillTokens.IsNull() {
					tfMaxBatchPrefillTokens.As(&endpointToCreate.Model.Image.TGI.MaxBatchPrefillTokens)
				}
				if tfMaxBatchTotalTokens, ok := tgiAttributes["max_batch_total_tokens"]; ok && !tfMaxBatchTotalTokens.IsNull() {
					tfMaxBatchTotalTokens.As(&endpointToCreate.Model.Image.TGI.MaxBatchTotalTokens)
				}
				if tfMaxInputLength, ok := tgiAttributes["max_input_length"]; ok && !tfMaxInputLength.IsNull() {
					var maxInputLengthBigFloat big.Float
					tfMaxInputLength.As(&maxInputLengthBigFloat)
					maxInputLengthInt, _ := maxInputLengthBigFloat.Int(nil)
					maxInputLengthIntPrimitive := int(maxInputLengthInt.Int64())
					endpointToCreate.Model.Image.TGI.MaxInputLength = &maxInputLengthIntPrimitive
				}
				if tfMaxTotalTokens, ok := tgiAttributes["max_total_tokens"]; ok && !tfMaxTotalTokens.IsNull() {
					var maxTotalTokensBigFloat big.Float
					tfMaxTotalTokens.As(&maxTotalTokensBigFloat)
					maxTotalTokensInt, _ := maxTotalTokensBigFloat.Int(nil)
					maxTotalTokensIntPrimitive := int(maxTotalTokensInt.Int64())
					endpointToCreate.Model.Image.TGI.MaxTotalTokens = &maxTotalTokensIntPrimitive
				}
				if tfPort, ok := tgiAttributes["port"]; ok && !tfPort.IsNull() {
					var portBigFloat big.Float
					tfPort.As(&portBigFloat)
					portInt, _ := portBigFloat.Int(nil)
					portIntPrimitive := int(portInt.Int64())
					endpointToCreate.Model.Image.TGI.Port = portIntPrimitive
				}
				if tfQuantize, ok := tgiAttributes["quantize"]; ok && !tfQuantize.IsNull() {
					var quantizeString string
					tfQuantize.As(&quantizeString)
					quantizeType := huggingface.QuantizeType(quantizeString)
					endpointToCreate.Model.Image.TGI.Quantize = &quantizeType
				}
				if tfUrl, ok := tgiAttributes["url"]; ok && !tfUrl.IsNull() {
					tfUrl.As(&endpointToCreate.Model.Image.TGI.URL)
				}
			}
			if tfTgiNeuron, ok := imageAttributes["tgi_neuron"]; ok && !tfTgiNeuron.IsNull() && tfTgiNeuron.IsKnown() {
				endpointToCreate.Model.Image.TGINeuron = &huggingface.TGINeuronImage{}
				var tgiNeuronAttributes map[string]tftypes.Value
				tfTgiNeuron.As(&tgiNeuronAttributes)

				if tfHealthRoute, ok := tgiNeuronAttributes["health_route"]; ok && !tfHealthRoute.IsNull() {
					var healthRoute string
					tfHealthRoute.As(&healthRoute)
					endpointToCreate.Model.Image.TGINeuron.HealthRoute = &healthRoute
				}
				if tfHfAutoCastType, ok := tgiNeuronAttributes["hf_auto_cast_type"]; ok && !tfHfAutoCastType.IsNull() {
					var hfAutoCastType string
					tfHfAutoCastType.As(&hfAutoCastType)
					hfHfAutoCastType := huggingface.AutoCastType(hfAutoCastType)
					endpointToCreate.Model.Image.TGINeuron.HfAutoCastType = &hfHfAutoCastType
				}
				if tfHfNumCores, ok := tgiNeuronAttributes["hf_num_cores"]; ok && !tfHfNumCores.IsNull() {
					var hfNumCoresBigFloat big.Float
					tfHfNumCores.As(&hfNumCoresBigFloat)
					hfNumCoresInt, _ := hfNumCoresBigFloat.Int(nil)
					hfNumCoresIntPrimitive := int(hfNumCoresInt.Int64())
					endpointToCreate.Model.Image.TGINeuron.HfNumCores = &hfNumCoresIntPrimitive
				}
				if tfMaxBatchPrefillTokens, ok := tgiNeuronAttributes["max_batch_prefill_tokens"]; ok && !tfMaxBatchPrefillTokens.IsNull() {
					var maxBatchPrefillTokensBigFloat big.Float
					tfMaxBatchPrefillTokens.As(&maxBatchPrefillTokensBigFloat)
					hfMaxBatchPrefillTokens, _ := maxBatchPrefillTokensBigFloat.Int(nil)
					maxBatchPrefillTokensIntPrimitive := int(hfMaxBatchPrefillTokens.Int64())
					endpointToCreate.Model.Image.TGINeuron.MaxBatchPrefillTokens = &maxBatchPrefillTokensIntPrimitive
				}
				if tfMaxBatchTotalTokens, ok := tgiNeuronAttributes["max_batch_total_tokens"]; ok && !tfMaxBatchTotalTokens.IsNull() {
					var maxBatchTotalTokensBigFloat big.Float
					tfMaxBatchTotalTokens.As(&maxBatchTotalTokensBigFloat)
					hfMaxBatchTotalTokens, _ := maxBatchTotalTokensBigFloat.Int(nil)
					maxBatchTotalTokensIntPrimitive := int(hfMaxBatchTotalTokens.Int64())
					endpointToCreate.Model.Image.TGINeuron.MaxBatchTotalTokens = &maxBatchTotalTokensIntPrimitive
				}
				if tfMaxInputLength, ok := tgiNeuronAttributes["max_input_length"]; ok && !tfMaxInputLength.IsNull() {
					var maxInputLengthBigFloat big.Float
					tfMaxInputLength.As(&maxInputLengthBigFloat)
					hfMaxInputLength, _ := maxInputLengthBigFloat.Int(nil)
					hfMaxInputLengthIntPrimitive := int(hfMaxInputLength.Int64())
					endpointToCreate.Model.Image.TGINeuron.MaxInputLength = &hfMaxInputLengthIntPrimitive
				}
				if tfMaxTotalTokens, ok := tgiNeuronAttributes["max_total_tokens"]; ok && !tfMaxTotalTokens.IsNull() {
					var maxTotalTokensBigFloat big.Float
					tfMaxTotalTokens.As(&maxTotalTokensBigFloat)
					hfMaxTotalTokens, _ := maxTotalTokensBigFloat.Int(nil)
					maxTotalTokensIntPrimitive := int(hfMaxTotalTokens.Int64())
					endpointToCreate.Model.Image.TGINeuron.MaxTotalTokens = &maxTotalTokensIntPrimitive
				}
				if tfPort, ok := tgiNeuronAttributes["port"]; ok && !tfPort.IsNull() {
					var portBigFloat big.Float
					tfPort.As(&portBigFloat)
					portInt, _ := portBigFloat.Int(nil)
					portIntPrimitive := int(portInt.Int64())
					endpointToCreate.Model.Image.TGINeuron.Port = portIntPrimitive
				}
				if tfUrl, ok := tgiNeuronAttributes["url"]; ok && !tfUrl.IsNull() {
					tfUrl.As(&endpointToCreate.Model.Image.TGINeuron.URL)
				}
			}
			if tfTei, ok := imageAttributes["tei"]; ok && !tfTei.IsNull() && tfTei.IsKnown() {
				endpointToCreate.Model.Image.TEI = &huggingface.TEIImage{}
				var teiAttributes map[string]tftypes.Value
				tfTei.As(&teiAttributes)

				if tfHealthRoute, ok := teiAttributes["health_route"]; ok && !tfHealthRoute.IsNull() {
					var healthRoute string
					tfHealthRoute.As(&healthRoute)
					endpointToCreate.Model.Image.TEI.HealthRoute = &healthRoute
				}
				if tfMaxBatchTokens, ok := teiAttributes["max_batch_tokens"]; ok && !tfMaxBatchTokens.IsNull() {
					var maxBatchTokensBigFloat big.Float
					tfMaxBatchTokens.As(&maxBatchTokensBigFloat)
					hfMaxBatchTokens, _ := maxBatchTokensBigFloat.Int(nil)
					maxBatchTokensIntPrimitive := int(hfMaxBatchTokens.Int64())
					endpointToCreate.Model.Image.TEI.MaxBatchTokens = &maxBatchTokensIntPrimitive
				}
				if tfMaxConcurrentRequests, ok := teiAttributes["max_concurrent_requests"]; ok && !tfMaxConcurrentRequests.IsNull() {
					var maxConcurrentRequestsBigFloat big.Float
					tfMaxConcurrentRequests.As(&maxConcurrentRequestsBigFloat)
					hfMaxConcurrentRequests, _ := maxConcurrentRequestsBigFloat.Int(nil)
					maxConcurrentRequestsIntPrimitive := int(hfMaxConcurrentRequests.Int64())
					endpointToCreate.Model.Image.TEI.MaxConcurrentRequests = &maxConcurrentRequestsIntPrimitive
				}
				if tfPooling, ok := teiAttributes["pooling"]; ok && !tfPooling.IsNull() {
					var poolingType string
					tfPooling.As(&poolingType)
					hfPoolingType := huggingface.PoolingType(poolingType)
					endpointToCreate.Model.Image.TEI.Pooling = &hfPoolingType
				}
				if tfPort, ok := teiAttributes["port"]; ok && !tfPort.IsNull() {
					var portBigFloat big.Float
					tfPort.As(&portBigFloat)
					portInt, _ := portBigFloat.Int(nil)
					portIntPrimitive := int(portInt.Int64())
					endpointToCreate.Model.Image.TEI.Port = portIntPrimitive
				}
				if tfUrl, ok := teiAttributes["url"]; ok && !tfUrl.IsNull() {
					tfUrl.As(&endpointToCreate.Model.Image.TEI.URL)
				}
			}
			if tfLlamacpp, ok := imageAttributes["llamacpp"]; ok && !tfLlamacpp.IsNull() && tfLlamacpp.IsKnown() {
				endpointToCreate.Model.Image.LlamaCpp = &huggingface.LlamaCppImage{}
				var llamaCppImageAttributes map[string]tftypes.Value
				tfLlamacpp.As(&llamaCppImageAttributes)

				if tfCtxSize, ok := llamaCppImageAttributes["ctx_size"]; ok && !tfCtxSize.IsNull() {
					var ctxSizeBigFloat big.Float
					tfCtxSize.As(&ctxSizeBigFloat)
					hfCtxSize, _ := ctxSizeBigFloat.Int(nil)
					ctxSizeIntPrimitive := int(hfCtxSize.Int64())
					endpointToCreate.Model.Image.LlamaCpp.CtxSize = ctxSizeIntPrimitive
				}
				if tfHealthRoute, ok := llamaCppImageAttributes["health_route"]; ok && !tfHealthRoute.IsNull() {
					var healthRoute string
					tfHealthRoute.As(&healthRoute)
					endpointToCreate.Model.Image.LlamaCpp.HealthRoute = &healthRoute
				}
				if tfMode, ok := llamaCppImageAttributes["mode"]; ok && !tfMode.IsNull() {
					var modeType string
					tfMode.As(&modeType)
					hfModeType := huggingface.ModelMode(modeType)
					endpointToCreate.Model.Image.LlamaCpp.Mode = &hfModeType
				}
				if tfModelPath, ok := llamaCppImageAttributes["model_path"]; ok && !tfModelPath.IsNull() {
					tfModelPath.As(&endpointToCreate.Model.Image.LlamaCpp.ModelPath)
				}
				if tfNGpuLayers, ok := llamaCppImageAttributes["n_gpu_layers"]; ok && !tfNGpuLayers.IsNull() {
					var nGpuLayersBigFloat big.Float
					tfNGpuLayers.As(&nGpuLayersBigFloat)
					hfNGpuLayers, _ := nGpuLayersBigFloat.Int(nil)
					nGpuLayersIntPrimitive := int(hfNGpuLayers.Int64())
					endpointToCreate.Model.Image.LlamaCpp.NGpuLayers = nGpuLayersIntPrimitive
				}
				if tfNParallel, ok := llamaCppImageAttributes["n_parallel"]; ok && !tfNParallel.IsNull() {
					var nParallelBigFloat big.Float
					tfNParallel.As(&nParallelBigFloat)
					hfNParallel, _ := nParallelBigFloat.Int(nil)
					nParallelIntPrimitive := int(hfNParallel.Int64())
					endpointToCreate.Model.Image.LlamaCpp.NParallel = nParallelIntPrimitive
				}
				if tfPooling, ok := llamaCppImageAttributes["pooling"]; ok && !tfPooling.IsNull() {
					var poolingType string
					tfPooling.As(&poolingType)
					hfPoolingType := huggingface.PoolingType(poolingType)
					endpointToCreate.Model.Image.LlamaCpp.Pooling = &hfPoolingType
				}
				if tfPort, ok := llamaCppImageAttributes["port"]; ok && !tfPort.IsNull() {
					var portBigFloat big.Float
					tfPort.As(&portBigFloat)
					portInt, _ := portBigFloat.Int(nil)
					portIntPrimitive := int(portInt.Int64())
					endpointToCreate.Model.Image.LlamaCpp.Port = portIntPrimitive
				}
				if tfThreadsHttp, ok := llamaCppImageAttributes["threads_http"]; ok && !tfThreadsHttp.IsNull() {
					var threadsHttpBigFloat big.Float
					tfThreadsHttp.As(&threadsHttpBigFloat)
					threadsHttpInt, _ := threadsHttpBigFloat.Int(nil)
					threadsHttpIntPrimitive := int(threadsHttpInt.Int64())
					endpointToCreate.Model.Image.LlamaCpp.ThreadsHttp = &threadsHttpIntPrimitive
				}
				if tfUrl, ok := llamaCppImageAttributes["url"]; ok && !tfUrl.IsNull() {
					tfUrl.As(&endpointToCreate.Model.Image.LlamaCpp.URL)
				}
				if tfVariant, ok := llamaCppImageAttributes["variant"]; ok && !tfVariant.IsNull() {
					tfVariant.As(&endpointToCreate.Model.Image.LlamaCpp.Variant)
				}
			}
			if tfCustom, ok := imageAttributes["custom"]; ok && !tfCustom.IsNull() && tfCustom.IsKnown() {
				endpointToCreate.Model.Image.Custom = &huggingface.CustomImage{}
				var customImageAttributes map[string]tftypes.Value
				tfCustom.As(&customImageAttributes)

				if tfCredentials, ok := customImageAttributes["credentials"]; ok && !tfCredentials.IsNull() {
					endpointToCreate.Model.Image.Custom.Credentials = &huggingface.Credentials{}
					var credentialsAttributes map[string]tftypes.Value
					tfCredentials.As(&credentialsAttributes)

					if tfUsername, ok := credentialsAttributes["username"]; ok && !tfUsername.IsNull() {
						tfUsername.As(&endpointToCreate.Model.Image.Custom.Credentials.Username)
					}
					if tfPassword, ok := credentialsAttributes["password"]; ok && !tfPassword.IsNull() {
						tfPassword.As(&endpointToCreate.Model.Image.Custom.Credentials.Password)
					}
				}
				if tfHealthRoute, ok := customImageAttributes["health_route"]; ok && !tfHealthRoute.IsNull() {
					var healthRoute string
					tfHealthRoute.As(&healthRoute)
					endpointToCreate.Model.Image.LlamaCpp.HealthRoute = &healthRoute
				}
				if tfPort, ok := customImageAttributes["port"]; ok && !tfPort.IsNull() {
					var portBigFloat big.Float
					tfPort.As(&portBigFloat)
					portInt, _ := portBigFloat.Int(nil)
					portIntPrimitive := int(portInt.Int64())
					endpointToCreate.Model.Image.Custom.Port = portIntPrimitive
				}
				if tfUrl, ok := customImageAttributes["url"]; ok && !tfUrl.IsNull() {
					tfUrl.As(&endpointToCreate.Model.Image.Custom.URL)
				}
			}
		}
	}

	if !plan.Tags.IsNull() {
		tags, _ := plan.Tags.ToTerraformValue(ctx)
		tags.As(&endpointToCreate.Tags)
	}

	if !plan.CacheHttpResponses.IsNull() {
		cacheHttpResponse := plan.CacheHttpResponses.ValueBool()
		endpointToCreate.CacheHttpResponses = &cacheHttpResponse
	}

	if !plan.ExperimentalFeatures.IsNull() {
		endpointToCreate.ExperimentalFeatures = &huggingface.ExperimentalFeatures{}
		experimentalFeaturesAttributes := plan.ExperimentalFeatures.Attributes()

		if cacheHttpResponse, ok := experimentalFeaturesAttributes["cache_http_response"]; ok && !cacheHttpResponse.IsNull() {
			tfCacheHttpResponse, _ := cacheHttpResponse.ToTerraformValue(ctx)
			tfCacheHttpResponse.As(endpointToCreate.ExperimentalFeatures.CacheHttpResponses)
		}
		if kvRouter, ok := experimentalFeaturesAttributes["kv_router"]; ok && !kvRouter.IsNull() {
			endpointToCreate.ExperimentalFeatures.KvRouter = &huggingface.KvRouter{}

			tfKvRouter, _ := kvRouter.ToTerraformValue(ctx)
			var kvRouterAttributes map[string]tftypes.Value
			tfKvRouter.As(&kvRouterAttributes)

			if tfTag, ok := kvRouterAttributes["tag"]; ok {
				tfTag.As(endpointToCreate.ExperimentalFeatures.KvRouter.Tag)
			}
		}
	}

	if !plan.PrivateService.IsNull() && !plan.PrivateService.IsUnknown() {
		endpointToCreate.PrivateService = &huggingface.EndpointPrivateService{}
		privateServiceAttributes := plan.PrivateService.Attributes()

		if accountId, ok := privateServiceAttributes["account_id"]; ok && !accountId.IsNull() {
			tfAccountId, _ := accountId.ToTerraformValue(ctx)
			tfAccountId.As(endpointToCreate.PrivateService.AccountID)
		}
		if shared, ok := privateServiceAttributes["shared"]; ok && !shared.IsNull() {
			tfShared, _ := shared.ToTerraformValue(ctx)
			tfShared.As(endpointToCreate.PrivateService.Shared)
		}
	}

	if !plan.Route.IsNull() && !plan.Route.IsUnknown() {
		endpointToCreate.Route = &huggingface.RouteSpec{}
		routeAttributes := plan.Route.Attributes()

		if domain, ok := routeAttributes["domain"]; ok && !domain.IsNull() {
			tfDomain, _ := domain.ToTerraformValue(ctx)
			tfDomain.As(endpointToCreate.Route.Domain)
		}
		if path, ok := routeAttributes["path"]; ok && !path.IsNull() {
			tfPath, _ := path.ToTerraformValue(ctx)
			tfPath.As(endpointToCreate.Route.Path)
		}
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
		ID:           types.StringValue(*endpoint.Compute.ID),
		Accelerator:  types.StringValue(string(endpoint.Compute.Accelerator)),
		InstanceType: types.StringValue(endpoint.Compute.InstanceType),
		InstanceSize: types.StringValue(endpoint.Compute.InstanceSize),
	}

	endpointComputeScaling := models.EndpointComputeScaling{
		MinReplica: types.Int32Value(int32(endpoint.Compute.Scaling.MinReplica)),
		MaxReplica: types.Int32Value(int32(endpoint.Compute.Scaling.MaxReplica)),
	}

	if endpoint.Compute.Scaling.ScaleToZeroTimeout != nil {
		endpointComputeScaling.ScaleToZeroTimeout = types.Int32Value(int32(*endpoint.Compute.Scaling.ScaleToZeroTimeout))
	}
	if endpoint.Compute.Scaling.Threshold != nil {
		threshold := *endpoint.Compute.Scaling.Threshold
		endpointComputeScaling.Threshold = types.Float64Value(threshold)
	}
	if endpoint.Compute.Scaling.Metric != nil {
		endpointScalingMetric := *endpoint.Compute.Scaling.Metric
		endpointComputeScaling.Metric = types.StringValue(string(endpointScalingMetric))
	}

	if endpoint.Compute.Scaling.Measure != nil {
		endpointComputeScalingMeasure := models.EndpointComputeScalingMeasure{}
		if endpoint.Compute.Scaling.Measure.HardwareUsage != nil {
			hardwareUsage := *endpoint.Compute.Scaling.Measure.HardwareUsage
			endpointComputeScalingMeasure.HardwareUsage = types.Float64Value(hardwareUsage)
		}
		if endpoint.Compute.Scaling.Measure.PendingRequests != nil {
			pendingRequests := *endpoint.Compute.Scaling.Measure.PendingRequests
			endpointComputeScalingMeasure.PendingRequests = types.Float64Value(pendingRequests)
		}

		endpointComputeScaling.Measure, diags = types.ObjectValueFrom(ctx, endpointComputeScalingMeasure.AttributeTypes(), endpointComputeScalingMeasure)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	endpointCompute.Scaling, diags = types.ObjectValueFrom(ctx, endpointComputeScaling.AttributeTypes(), endpointComputeScaling)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
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

	modelImage := models.ModelImage{}

	if endpoint.Model.Image.HuggingFace != nil {
		huggingFaceImage := models.ModelImageHuggingface{}

		modelImage.HuggingFace, diags = types.ObjectValueFrom(ctx, huggingFaceImage.AttributeTypes(), huggingFaceImage)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		modelImage.HuggingFace, diags = types.ObjectValueFrom(ctx, models.ModelImageHuggingface{}.AttributeTypes(), models.ModelImageHuggingface{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if endpoint.Model.Image.HuggingFaceNeuron != nil {
		huggingFaceNeuronImage := models.ModelImageHuggingfaceNeuron{}

		if endpoint.Model.Image.HuggingFaceNeuron.BatchSize != nil {
			huggingFaceNeuronImage.BatchSize = types.Int32Value(int32(*endpoint.Model.Image.HuggingFaceNeuron.BatchSize))
		}

		huggingFaceNeuronImage.NeuronCache = types.StringValue(endpoint.Model.Image.HuggingFaceNeuron.NeuronCache)

		if endpoint.Model.Image.HuggingFaceNeuron.SequenceLength != nil {
			huggingFaceNeuronImage.SequenceLength = types.Int32Value(int32(*endpoint.Model.Image.HuggingFaceNeuron.SequenceLength))
		}

		modelImage.HuggingFaceNeuron, diags = types.ObjectValueFrom(ctx, huggingFaceNeuronImage.AttributeTypes(), huggingFaceNeuronImage)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		modelImage.HuggingFaceNeuron, diags = types.ObjectValueFrom(ctx, models.ModelImageHuggingfaceNeuron{}.AttributeTypes(), models.ModelImageHuggingfaceNeuron{
			NeuronCache:    types.StringNull(),
			BatchSize:      types.Int32Null(),
			SequenceLength: types.Int32Null(),
		})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if endpoint.Model.Image.TGI != nil {
		huggingFaceTgiImage := models.ModelImageTgi{}

		if endpoint.Model.Image.TGI.HealthRoute != nil {
			// huggingFaceTgiImage.HealthRoute = types.StringValue(*endpoint.Model.Image.TGI.HealthRoute)
		}

		huggingFaceTgiImage.Port = types.Int32Value(int32(endpoint.Model.Image.TGI.Port))

		huggingFaceTgiImage.Url = types.StringValue(endpoint.Model.Image.TGI.URL)

		if endpoint.Model.Image.TGI.MaxBatchPrefillTokens != nil {
			huggingFaceTgiImage.MaxBatchPrefillTokens = types.Int32Value(int32(*endpoint.Model.Image.TGI.MaxBatchPrefillTokens))
		}
		if endpoint.Model.Image.TGI.MaxBatchTotalTokens != nil {
			huggingFaceTgiImage.MaxBatchTotalTokens = types.Int32Value(int32(*endpoint.Model.Image.TGI.MaxBatchTotalTokens))
		}
		if endpoint.Model.Image.TGI.MaxInputLength != nil {
			huggingFaceTgiImage.MaxInputLength = types.Int32Value(int32(*endpoint.Model.Image.TGI.MaxInputLength))
		}
		if endpoint.Model.Image.TGI.MaxTotalTokens != nil {
			huggingFaceTgiImage.MaxTotalTokens = types.Int32Value(int32(*endpoint.Model.Image.TGI.MaxTotalTokens))
		}

		huggingFaceTgiImage.DisableCustomKernels = types.BoolValue(endpoint.Model.Image.TGI.DisableCustomKernels)

		if endpoint.Model.Image.TGI.Quantize != nil {
			huggingFaceTgiImage.Quantize = types.StringValue(string(*endpoint.Model.Image.TGI.Quantize))
		}

		modelImage.TGI, diags = types.ObjectValueFrom(ctx, huggingFaceTgiImage.AttributeTypes(), huggingFaceTgiImage)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		modelImage.TGI, diags = types.ObjectValueFrom(ctx, models.ModelImageTgi{}.AttributeTypes(), models.ModelImageTgi{
			HealthRoute:           types.StringNull(),
			Port:                  types.Int32Null(),
			Url:                   types.StringNull(),
			MaxBatchPrefillTokens: types.Int32Null(),
			MaxBatchTotalTokens:   types.Int32Null(),
			MaxInputLength:        types.Int32Null(),
			MaxTotalTokens:        types.Int32Null(),
			DisableCustomKernels:  types.BoolNull(),
			Quantize:              types.StringNull(),
		})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if endpoint.Model.Image.TGINeuron != nil {
		huggingFaceTgiNeuronImage := models.ModelImageTgiNeuron{}

		if endpoint.Model.Image.TGINeuron.HealthRoute != nil {
			// huggingFaceTgiNeuronImage.HealthRoute = types.StringValue(*endpoint.Model.Image.TGINeuron.HealthRoute)
		}

		huggingFaceTgiNeuronImage.Port = types.Int32Value(int32(endpoint.Model.Image.TGINeuron.Port))

		huggingFaceTgiNeuronImage.Url = types.StringValue(endpoint.Model.Image.TGINeuron.URL)

		if endpoint.Model.Image.TGINeuron.MaxBatchPrefillTokens != nil {
			huggingFaceTgiNeuronImage.MaxBatchPrefillTokens = types.Int32Value(int32(*endpoint.Model.Image.TGINeuron.MaxBatchPrefillTokens))
		}
		if endpoint.Model.Image.TGINeuron.MaxBatchTotalTokens != nil {
			huggingFaceTgiNeuronImage.MaxBatchTotalTokens = types.Int32Value(int32(*endpoint.Model.Image.TGI.MaxBatchTotalTokens))
		}
		if endpoint.Model.Image.TGINeuron.MaxInputLength != nil {
			huggingFaceTgiNeuronImage.MaxInputLength = types.Int32Value(int32(*endpoint.Model.Image.TGINeuron.MaxInputLength))
		}
		if endpoint.Model.Image.TGINeuron.MaxTotalTokens != nil {
			huggingFaceTgiNeuronImage.MaxTotalTokens = types.Int32Value(int32(*endpoint.Model.Image.TGINeuron.MaxTotalTokens))
		}
		if endpoint.Model.Image.TGINeuron.MaxTotalTokens != nil {
			huggingFaceTgiNeuronImage.MaxTotalTokens = types.Int32Value(int32(*endpoint.Model.Image.TGI.MaxTotalTokens))
		}
		if endpoint.Model.Image.TGINeuron.HfAutoCastType != nil {
			huggingFaceTgiNeuronImage.HfAutoCastType = types.StringValue(string(*endpoint.Model.Image.TGINeuron.HfAutoCastType))
		}
		if endpoint.Model.Image.TGINeuron.HfNumCores != nil {
			huggingFaceTgiNeuronImage.HfNumCores = types.Int32Value(int32(*endpoint.Model.Image.TGINeuron.HfNumCores))
		}

		modelImage.TGINeuron, diags = types.ObjectValueFrom(ctx, huggingFaceTgiNeuronImage.AttributeTypes(), huggingFaceTgiNeuronImage)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		modelImage.TGINeuron, diags = types.ObjectValueFrom(ctx, models.ModelImageTgiNeuron{}.AttributeTypes(), models.ModelImageTgiNeuron{
			HealthRoute:           types.StringNull(),
			Port:                  types.Int32Null(),
			Url:                   types.StringNull(),
			MaxBatchPrefillTokens: types.Int32Null(),
			MaxBatchTotalTokens:   types.Int32Null(),
			MaxInputLength:        types.Int32Null(),
			MaxTotalTokens:        types.Int32Null(),
			HfAutoCastType:        types.StringNull(),
			HfNumCores:            types.Int32Null(),
		})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if endpoint.Model.Image.TEI != nil {
		huggingFaceTeiNeuronImage := models.ModelImageTei{}

		if endpoint.Model.Image.TEI.HealthRoute != nil {
			// huggingFaceTeiNeuronImage.HealthRoute = types.StringValue(*endpoint.Model.Image.TGINeuron.HealthRoute)
		}

		huggingFaceTeiNeuronImage.Port = types.Int32Value(int32(endpoint.Model.Image.TEI.Port))

		huggingFaceTeiNeuronImage.URL = types.StringValue(endpoint.Model.Image.TEI.URL)

		if endpoint.Model.Image.TEI.MaxBatchTokens != nil {
			huggingFaceTeiNeuronImage.MaxBatchTokens = types.Int32Value(int32(*endpoint.Model.Image.TEI.MaxBatchTokens))
		}
		if endpoint.Model.Image.TEI.MaxConcurrentRequests != nil {
			huggingFaceTeiNeuronImage.MaxConcurrentRequests = types.Int32Value(int32(*endpoint.Model.Image.TEI.MaxConcurrentRequests))
		}
		if endpoint.Model.Image.TEI.Pooling != nil {
			huggingFaceTeiNeuronImage.Pooling = types.StringValue(string(*endpoint.Model.Image.TEI.Pooling))
		}

		modelImage.TEI, diags = types.ObjectValueFrom(ctx, huggingFaceTeiNeuronImage.AttributeTypes(), huggingFaceTeiNeuronImage)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		modelImage.TEI, diags = types.ObjectValueFrom(ctx, models.ModelImageTei{}.AttributeTypes(), models.ModelImageTei{
			HealthRoute:           types.StringNull(),
			Port:                  types.Int32Null(),
			URL:                   types.StringNull(),
			MaxBatchTokens:        types.Int32Null(),
			MaxConcurrentRequests: types.Int32Null(),
			Pooling:               types.StringNull(),
		})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if endpoint.Model.Image.LlamaCpp != nil {
		llamaCppImage := models.ModelImageLlamacpp{}

		if endpoint.Model.Image.LlamaCpp.HealthRoute != nil {
			// llamaCppImage.HealthRoute = types.StringValue(*endpoint.Model.Image.LlamaCpp.HealthRoute)
		}

		llamaCppImage.Port = types.Int32Value(int32(endpoint.Model.Image.LlamaCpp.Port))

		llamaCppImage.URL = types.StringValue(endpoint.Model.Image.LlamaCpp.URL)

		llamaCppImage.CtxSize = types.Int32Value(int32(endpoint.Model.Image.LlamaCpp.CtxSize))

		if endpoint.Model.Image.LlamaCpp.Mode != nil {
			llamaCppImage.Mode = types.StringValue(string(*endpoint.Model.Image.LlamaCpp.Mode))
		}

		llamaCppImage.ModelPath = types.StringValue(string(endpoint.Model.Image.LlamaCpp.ModelPath))

		llamaCppImage.NGpuLayers = types.Int32Value(int32(endpoint.Model.Image.LlamaCpp.NGpuLayers))

		llamaCppImage.NParallel = types.Int32Value(int32(endpoint.Model.Image.LlamaCpp.NParallel))

		if endpoint.Model.Image.LlamaCpp.Pooling != nil {
			llamaCppImage.Pooling = types.StringValue(string(*endpoint.Model.Image.LlamaCpp.Pooling))
		}

		if endpoint.Model.Image.LlamaCpp.ThreadsHttp != nil {
			llamaCppImage.ThreadsHttp = types.Int32Value(int32(*endpoint.Model.Image.LlamaCpp.ThreadsHttp))
		}

		if endpoint.Model.Image.LlamaCpp.Variant != nil {
			llamaCppImage.Variant = types.StringValue(string(*endpoint.Model.Image.LlamaCpp.Variant))
		}

		modelImage.LlamaCpp, diags = types.ObjectValueFrom(ctx, llamaCppImage.AttributeTypes(), llamaCppImage)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		modelImage.LlamaCpp, diags = types.ObjectValueFrom(ctx, models.ModelImageLlamacpp{}.AttributeTypes(), models.ModelImageLlamacpp{
			HealthRoute: types.StringNull(),
			Port:        types.Int32Null(),
			URL:         types.StringNull(),
			CtxSize:     types.Int32Null(),
			Mode:        types.StringNull(),
			ModelPath:   types.StringNull(),
			NGpuLayers:  types.Int32Null(),
			NParallel:   types.Int32Null(),
			Pooling:     types.StringNull(),
			ThreadsHttp: types.Int32Null(),
			Variant:     types.StringNull(),
		})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if endpoint.Model.Image.Custom != nil {
		customImage := models.ModelImageCustom{}

		if endpoint.Model.Image.Custom.HealthRoute != nil {
			customImage.HealthRoute = types.StringValue(*endpoint.Model.Image.Custom.HealthRoute)
		}

		customImage.Port = types.Int32Value(int32(endpoint.Model.Image.Custom.Port))

		customImage.URL = types.StringValue(endpoint.Model.Image.Custom.URL)

		if endpoint.Model.Image.Custom.Credentials != nil {
			credentials := models.Credentials{
				Username: types.StringValue(string(endpoint.Model.Image.Custom.Credentials.Username)),
				Password: types.StringValue(string(*endpoint.Model.Image.Custom.Credentials.Password)),
			}

			customImage.Credentials, diags = types.ObjectValueFrom(ctx, credentials.AttributeTypes(), credentials)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
		}

		modelImage.Custom, diags = types.ObjectValueFrom(ctx, customImage.AttributeTypes(), customImage)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		mCredentials := models.Credentials{
			Username: types.StringNull(),
			Password: types.StringNull(),
		}
		credentials, diags := types.ObjectValueFrom(ctx, mCredentials.AttributeTypes(), mCredentials)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		modelImage.Custom, diags = types.ObjectValueFrom(ctx, models.ModelImageCustom{}.AttributeTypes(), models.ModelImageCustom{
			HealthRoute: types.StringNull(),
			Port:        types.Int32Null(),
			URL:         types.StringNull(),
			Credentials: credentials,
		})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	endpointModel.Image, diags = types.ObjectValueFrom(ctx, modelImage.AttributeTypes(), modelImage)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
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
