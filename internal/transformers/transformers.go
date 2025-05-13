package transformers

import (
	"context"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	huggingface "github.com/sebps/huggingface-client/client"
	"github.com/sebps/terraform-provider-huggingface/internal/models"
	"github.com/sebps/terraform-provider-huggingface/internal/states"
)

func FromModelToProvider(
	ctx context.Context,
	input *states.EndpointResourceState,
) (output huggingface.Endpoint) {
	output = huggingface.Endpoint{
		Name: input.Name.ValueString(),
		Type: huggingface.EndpointType(input.Type.ValueString()),
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
	cloudProviderAttributes := input.CloudProvider.Attributes()
	if vendor, ok := cloudProviderAttributes["vendor"]; ok && !vendor.IsNull() {
		tfVendor, _ := vendor.ToTerraformValue(ctx)
		tfVendor.As(&output.Provider.Vendor)
	}
	if region, ok := cloudProviderAttributes["region"]; ok && !region.IsNull() {
		tfRegion, _ := region.ToTerraformValue(ctx)
		tfRegion.As(&output.Provider.Region)
	}

	// Compute
	if !input.Compute.IsNull() {
		computeAttributes := input.Compute.Attributes()
		if accelerator, ok := computeAttributes["accelerator"]; ok && !accelerator.IsNull() {
			tfAccelerator, _ := accelerator.ToTerraformValue(ctx)
			var hfAccelerator string
			tfAccelerator.As(&hfAccelerator)
			output.Compute.Accelerator = huggingface.AcceleratorType(hfAccelerator)
		}
		if instanceType, ok := computeAttributes["instance_type"]; ok && !instanceType.IsNull() {
			tfInstanceType, _ := instanceType.ToTerraformValue(ctx)
			tfInstanceType.As(&output.Compute.InstanceType)
		}
		if instanceSize, ok := computeAttributes["instance_size"]; ok && !instanceSize.IsNull() {
			tfInstanceType, _ := instanceSize.ToTerraformValue(ctx)
			tfInstanceType.As(&output.Compute.InstanceSize)
		}
		if scaling, ok := computeAttributes["scaling"]; ok && !scaling.IsNull() {
			tfScaling, _ := scaling.ToTerraformValue(ctx)
			var scalingAttributes map[string]tftypes.Value
			tfScaling.As(&scalingAttributes)

			if tfMinReplica, ok := scalingAttributes["min_replica"]; ok && !tfMinReplica.IsNull() {
				var minReplicaBigFloat big.Float
				tfMinReplica.As(&minReplicaBigFloat)
				minReplicaInt, _ := minReplicaBigFloat.Int(nil)
				output.Compute.Scaling.MinReplica = int(minReplicaInt.Int64())
			} else {
				output.Compute.Scaling.MinReplica = 0
			}
			if tfMaxReplica, ok := scalingAttributes["max_replica"]; ok && !tfMaxReplica.IsNull() {
				var maxReplicaBigFloat big.Float
				tfMaxReplica.As(&maxReplicaBigFloat)
				maxReplicaInt, _ := maxReplicaBigFloat.Int(nil)
				output.Compute.Scaling.MaxReplica = int(maxReplicaInt.Int64())
			} else {
				output.Compute.Scaling.MaxReplica = 1
			}
			if tfMetric, ok := scalingAttributes["metric"]; ok && !tfMetric.IsNull() {
				tfMetric.As(&output.Compute.Scaling.Metric)
			}
			if tfScaleToZeroTimeout, ok := scalingAttributes["scale_to_zero_timeout"]; ok && !tfScaleToZeroTimeout.IsNull() && tfScaleToZeroTimeout.IsKnown() {
				var scaleToZeroTimeoutBigFloat big.Float
				tfScaleToZeroTimeout.As(&scaleToZeroTimeoutBigFloat)
				scaleToZeroTimeoutInt, _ := scaleToZeroTimeoutBigFloat.Int(nil)
				scaleToZeroTimeoutIntPrimitive := int(scaleToZeroTimeoutInt.Int64())
				output.Compute.Scaling.ScaleToZeroTimeout = &scaleToZeroTimeoutIntPrimitive
			}
			if tfThreshold, ok := scalingAttributes["threshold"]; ok && !tfThreshold.IsNull() && tfThreshold.IsKnown() {
				var thresholdBigFloat big.Float
				tfThreshold.As(&thresholdBigFloat)
				thresholdFloatPrimitive, _ := thresholdBigFloat.Float64()
				output.Compute.Scaling.Threshold = &thresholdFloatPrimitive
			}
			if tfMeasure, ok := scalingAttributes["measure"]; ok && !tfMeasure.IsNull() {
				var measureAttributes map[string]tftypes.Value
				tfMeasure.As(&measureAttributes)

				if tfHardwareUsage, ok := measureAttributes["hardware_usage"]; ok && !tfHardwareUsage.IsNull() {
					var hardwareUsageBigFloat big.Float
					tfHardwareUsage.As(&hardwareUsageBigFloat)
					hardwareUsageFloatPrimitive, _ := hardwareUsageBigFloat.Float64()
					output.Compute.Scaling.Measure.HardwareUsage = &hardwareUsageFloatPrimitive
				}
				if tfPendingRequests, ok := measureAttributes["pending_requests"]; ok && !tfPendingRequests.IsNull() {
					var pendingRequestsBigFloat big.Float
					tfPendingRequests.As(&pendingRequestsBigFloat)
					pendintRequestsFloatPrimitive, _ := pendingRequestsBigFloat.Float64()
					output.Compute.Scaling.Measure.PendingRequests = &pendintRequestsFloatPrimitive
				}
			}
		}
	}

	// Model
	if !input.Model.IsNull() {
		modelAttributes := input.Model.Attributes()
		if repository, ok := modelAttributes["repository"]; ok && !repository.IsNull() {
			tfRepository, _ := repository.ToTerraformValue(ctx)
			tfRepository.As(&output.Model.Repository)
		}
		if framework, ok := modelAttributes["framework"]; ok && !framework.IsNull() {
			tfFramework, _ := framework.ToTerraformValue(ctx)
			var hfFramework string
			tfFramework.As(&hfFramework)
			output.Model.Framework = huggingface.EndpointFramework(hfFramework)
		}
		if task, ok := modelAttributes["task"]; ok && !task.IsNull() {
			tfTask, _ := task.ToTerraformValue(ctx)
			tfTask.As(&output.Model.Task)
		}
		if image, ok := modelAttributes["image"]; ok && !image.IsNull() {
			tfImage, _ := image.ToTerraformValue(ctx)
			var imageAttributes map[string]tftypes.Value
			tfImage.As(&imageAttributes)

			if tfHuggingface, ok := imageAttributes["huggingface"]; ok && !tfHuggingface.IsNull() && tfHuggingface.IsKnown() {
				output.Model.Image.HuggingFace = &huggingface.HuggingFaceImage{}
			}
			if tfHuggingfaceNeuron, ok := imageAttributes["huggingface_neuron"]; ok && !tfHuggingfaceNeuron.IsNull() && tfHuggingfaceNeuron.IsKnown() {
				output.Model.Image.HuggingFaceNeuron = &huggingface.HuggingFaceNeuronImage{}
				var huggingfaceNeuronAttributes map[string]tftypes.Value
				tfHuggingfaceNeuron.As(&huggingfaceNeuronAttributes)

				if tfBatchSize, ok := huggingfaceNeuronAttributes["batch_size"]; ok && !tfBatchSize.IsNull() {
					var batchSizeBigFloat big.Float
					tfBatchSize.As(&batchSizeBigFloat)
					batchSizeInt, _ := batchSizeBigFloat.Int(nil)
					batchSizeIntPrimitive := int(batchSizeInt.Int64())
					output.Model.Image.HuggingFaceNeuron.BatchSize = &batchSizeIntPrimitive
				}
				if tfNeuronCache, ok := huggingfaceNeuronAttributes["neuron_cache"]; ok && !tfNeuronCache.IsNull() {
					tfNeuronCache.As(&output.Model.Image.HuggingFaceNeuron.NeuronCache)
				}
				if tfSequenceLength, ok := huggingfaceNeuronAttributes["sequence_length"]; ok && !tfSequenceLength.IsNull() {
					var sequenceLengthBigFloat big.Float
					tfSequenceLength.As(&sequenceLengthBigFloat)
					sequenceLengthInt, _ := sequenceLengthBigFloat.Int(nil)
					sequenceLengthIntPrimitive := int(sequenceLengthInt.Int64())
					output.Model.Image.HuggingFaceNeuron.SequenceLength = &sequenceLengthIntPrimitive
				}
			}
			if tfTgi, ok := imageAttributes["tgi"]; ok && !tfTgi.IsNull() && tfTgi.IsKnown() {
				output.Model.Image.TGI = &huggingface.TGIImage{}
				var tgiAttributes map[string]tftypes.Value
				tfTgi.As(&tgiAttributes)

				if tfDisableCustomKernels, ok := tgiAttributes["disable_custom_kernels"]; ok && !tfDisableCustomKernels.IsNull() {
					tfDisableCustomKernels.As(&output.Model.Image.TGI.DisableCustomKernels)
				}
				if tfHealthRoute, ok := tgiAttributes["health_route"]; ok && !tfHealthRoute.IsNull() {
					tfHealthRoute.As(&output.Model.Image.TGI.HealthRoute)
				}
				if tfMaxBatchPrefillTokens, ok := tgiAttributes["max_batch_prefill_tokens"]; ok && !tfMaxBatchPrefillTokens.IsNull() {
					tfMaxBatchPrefillTokens.As(&output.Model.Image.TGI.MaxBatchPrefillTokens)
				}
				if tfMaxBatchTotalTokens, ok := tgiAttributes["max_batch_total_tokens"]; ok && !tfMaxBatchTotalTokens.IsNull() {
					tfMaxBatchTotalTokens.As(&output.Model.Image.TGI.MaxBatchTotalTokens)
				}
				if tfMaxInputLength, ok := tgiAttributes["max_input_length"]; ok && !tfMaxInputLength.IsNull() {
					var maxInputLengthBigFloat big.Float
					tfMaxInputLength.As(&maxInputLengthBigFloat)
					maxInputLengthInt, _ := maxInputLengthBigFloat.Int(nil)
					maxInputLengthIntPrimitive := int(maxInputLengthInt.Int64())
					output.Model.Image.TGI.MaxInputLength = &maxInputLengthIntPrimitive
				}
				if tfMaxTotalTokens, ok := tgiAttributes["max_total_tokens"]; ok && !tfMaxTotalTokens.IsNull() {
					var maxTotalTokensBigFloat big.Float
					tfMaxTotalTokens.As(&maxTotalTokensBigFloat)
					maxTotalTokensInt, _ := maxTotalTokensBigFloat.Int(nil)
					maxTotalTokensIntPrimitive := int(maxTotalTokensInt.Int64())
					output.Model.Image.TGI.MaxTotalTokens = &maxTotalTokensIntPrimitive
				}
				if tfPort, ok := tgiAttributes["port"]; ok && !tfPort.IsNull() {
					var portBigFloat big.Float
					tfPort.As(&portBigFloat)
					portInt, _ := portBigFloat.Int(nil)
					portIntPrimitive := int(portInt.Int64())
					output.Model.Image.TGI.Port = portIntPrimitive
				}
				if tfQuantize, ok := tgiAttributes["quantize"]; ok && !tfQuantize.IsNull() {
					var quantizeString string
					tfQuantize.As(&quantizeString)
					quantizeType := huggingface.QuantizeType(quantizeString)
					output.Model.Image.TGI.Quantize = &quantizeType
				}
				if tfUrl, ok := tgiAttributes["url"]; ok && !tfUrl.IsNull() {
					tfUrl.As(&output.Model.Image.TGI.URL)
				}
			}
			if tfTgiNeuron, ok := imageAttributes["tgi_neuron"]; ok && !tfTgiNeuron.IsNull() && tfTgiNeuron.IsKnown() {
				output.Model.Image.TGINeuron = &huggingface.TGINeuronImage{}
				var tgiNeuronAttributes map[string]tftypes.Value
				tfTgiNeuron.As(&tgiNeuronAttributes)

				if tfHealthRoute, ok := tgiNeuronAttributes["health_route"]; ok && !tfHealthRoute.IsNull() {
					var healthRoute string
					tfHealthRoute.As(&healthRoute)
					output.Model.Image.TGINeuron.HealthRoute = &healthRoute
				}
				if tfHfAutoCastType, ok := tgiNeuronAttributes["hf_auto_cast_type"]; ok && !tfHfAutoCastType.IsNull() {
					var hfAutoCastType string
					tfHfAutoCastType.As(&hfAutoCastType)
					hfHfAutoCastType := huggingface.AutoCastType(hfAutoCastType)
					output.Model.Image.TGINeuron.HfAutoCastType = &hfHfAutoCastType
				}
				if tfHfNumCores, ok := tgiNeuronAttributes["hf_num_cores"]; ok && !tfHfNumCores.IsNull() {
					var hfNumCoresBigFloat big.Float
					tfHfNumCores.As(&hfNumCoresBigFloat)
					hfNumCoresInt, _ := hfNumCoresBigFloat.Int(nil)
					hfNumCoresIntPrimitive := int(hfNumCoresInt.Int64())
					output.Model.Image.TGINeuron.HfNumCores = &hfNumCoresIntPrimitive
				}
				if tfMaxBatchPrefillTokens, ok := tgiNeuronAttributes["max_batch_prefill_tokens"]; ok && !tfMaxBatchPrefillTokens.IsNull() {
					var maxBatchPrefillTokensBigFloat big.Float
					tfMaxBatchPrefillTokens.As(&maxBatchPrefillTokensBigFloat)
					hfMaxBatchPrefillTokens, _ := maxBatchPrefillTokensBigFloat.Int(nil)
					maxBatchPrefillTokensIntPrimitive := int(hfMaxBatchPrefillTokens.Int64())
					output.Model.Image.TGINeuron.MaxBatchPrefillTokens = &maxBatchPrefillTokensIntPrimitive
				}
				if tfMaxBatchTotalTokens, ok := tgiNeuronAttributes["max_batch_total_tokens"]; ok && !tfMaxBatchTotalTokens.IsNull() {
					var maxBatchTotalTokensBigFloat big.Float
					tfMaxBatchTotalTokens.As(&maxBatchTotalTokensBigFloat)
					hfMaxBatchTotalTokens, _ := maxBatchTotalTokensBigFloat.Int(nil)
					maxBatchTotalTokensIntPrimitive := int(hfMaxBatchTotalTokens.Int64())
					output.Model.Image.TGINeuron.MaxBatchTotalTokens = &maxBatchTotalTokensIntPrimitive
				}
				if tfMaxInputLength, ok := tgiNeuronAttributes["max_input_length"]; ok && !tfMaxInputLength.IsNull() {
					var maxInputLengthBigFloat big.Float
					tfMaxInputLength.As(&maxInputLengthBigFloat)
					hfMaxInputLength, _ := maxInputLengthBigFloat.Int(nil)
					hfMaxInputLengthIntPrimitive := int(hfMaxInputLength.Int64())
					output.Model.Image.TGINeuron.MaxInputLength = &hfMaxInputLengthIntPrimitive
				}
				if tfMaxTotalTokens, ok := tgiNeuronAttributes["max_total_tokens"]; ok && !tfMaxTotalTokens.IsNull() {
					var maxTotalTokensBigFloat big.Float
					tfMaxTotalTokens.As(&maxTotalTokensBigFloat)
					hfMaxTotalTokens, _ := maxTotalTokensBigFloat.Int(nil)
					maxTotalTokensIntPrimitive := int(hfMaxTotalTokens.Int64())
					output.Model.Image.TGINeuron.MaxTotalTokens = &maxTotalTokensIntPrimitive
				}
				if tfPort, ok := tgiNeuronAttributes["port"]; ok && !tfPort.IsNull() {
					var portBigFloat big.Float
					tfPort.As(&portBigFloat)
					portInt, _ := portBigFloat.Int(nil)
					portIntPrimitive := int(portInt.Int64())
					output.Model.Image.TGINeuron.Port = portIntPrimitive
				}
				if tfUrl, ok := tgiNeuronAttributes["url"]; ok && !tfUrl.IsNull() {
					tfUrl.As(&output.Model.Image.TGINeuron.URL)
				}
			}
			if tfTei, ok := imageAttributes["tei"]; ok && !tfTei.IsNull() && tfTei.IsKnown() {
				output.Model.Image.TEI = &huggingface.TEIImage{}
				var teiAttributes map[string]tftypes.Value
				tfTei.As(&teiAttributes)

				if tfHealthRoute, ok := teiAttributes["health_route"]; ok && !tfHealthRoute.IsNull() {
					var healthRoute string
					tfHealthRoute.As(&healthRoute)
					output.Model.Image.TEI.HealthRoute = &healthRoute
				}
				if tfMaxBatchTokens, ok := teiAttributes["max_batch_tokens"]; ok && !tfMaxBatchTokens.IsNull() {
					var maxBatchTokensBigFloat big.Float
					tfMaxBatchTokens.As(&maxBatchTokensBigFloat)
					hfMaxBatchTokens, _ := maxBatchTokensBigFloat.Int(nil)
					maxBatchTokensIntPrimitive := int(hfMaxBatchTokens.Int64())
					output.Model.Image.TEI.MaxBatchTokens = &maxBatchTokensIntPrimitive
				}
				if tfMaxConcurrentRequests, ok := teiAttributes["max_concurrent_requests"]; ok && !tfMaxConcurrentRequests.IsNull() {
					var maxConcurrentRequestsBigFloat big.Float
					tfMaxConcurrentRequests.As(&maxConcurrentRequestsBigFloat)
					hfMaxConcurrentRequests, _ := maxConcurrentRequestsBigFloat.Int(nil)
					maxConcurrentRequestsIntPrimitive := int(hfMaxConcurrentRequests.Int64())
					output.Model.Image.TEI.MaxConcurrentRequests = &maxConcurrentRequestsIntPrimitive
				}
				if tfPooling, ok := teiAttributes["pooling"]; ok && !tfPooling.IsNull() {
					var poolingType string
					tfPooling.As(&poolingType)
					hfPoolingType := huggingface.PoolingType(poolingType)
					output.Model.Image.TEI.Pooling = &hfPoolingType
				}
				if tfPort, ok := teiAttributes["port"]; ok && !tfPort.IsNull() {
					var portBigFloat big.Float
					tfPort.As(&portBigFloat)
					portInt, _ := portBigFloat.Int(nil)
					portIntPrimitive := int(portInt.Int64())
					output.Model.Image.TEI.Port = portIntPrimitive
				}
				if tfUrl, ok := teiAttributes["url"]; ok && !tfUrl.IsNull() {
					tfUrl.As(&output.Model.Image.TEI.URL)
				}
			}
			if tfLlamacpp, ok := imageAttributes["llamacpp"]; ok && !tfLlamacpp.IsNull() && tfLlamacpp.IsKnown() {
				output.Model.Image.LlamaCpp = &huggingface.LlamaCppImage{}
				var llamaCppImageAttributes map[string]tftypes.Value
				tfLlamacpp.As(&llamaCppImageAttributes)

				if tfCtxSize, ok := llamaCppImageAttributes["ctx_size"]; ok && !tfCtxSize.IsNull() {
					var ctxSizeBigFloat big.Float
					tfCtxSize.As(&ctxSizeBigFloat)
					hfCtxSize, _ := ctxSizeBigFloat.Int(nil)
					ctxSizeIntPrimitive := int(hfCtxSize.Int64())
					output.Model.Image.LlamaCpp.CtxSize = ctxSizeIntPrimitive
				}
				if tfHealthRoute, ok := llamaCppImageAttributes["health_route"]; ok && !tfHealthRoute.IsNull() {
					var healthRoute string
					tfHealthRoute.As(&healthRoute)
					output.Model.Image.LlamaCpp.HealthRoute = &healthRoute
				}
				if tfMode, ok := llamaCppImageAttributes["mode"]; ok && !tfMode.IsNull() {
					var modeType string
					tfMode.As(&modeType)
					hfModeType := huggingface.ModelMode(modeType)
					output.Model.Image.LlamaCpp.Mode = &hfModeType
				}
				if tfModelPath, ok := llamaCppImageAttributes["model_path"]; ok && !tfModelPath.IsNull() {
					tfModelPath.As(&output.Model.Image.LlamaCpp.ModelPath)
				}
				if tfNGpuLayers, ok := llamaCppImageAttributes["n_gpu_layers"]; ok && !tfNGpuLayers.IsNull() {
					var nGpuLayersBigFloat big.Float
					tfNGpuLayers.As(&nGpuLayersBigFloat)
					hfNGpuLayers, _ := nGpuLayersBigFloat.Int(nil)
					nGpuLayersIntPrimitive := int(hfNGpuLayers.Int64())
					output.Model.Image.LlamaCpp.NGpuLayers = nGpuLayersIntPrimitive
				}
				if tfNParallel, ok := llamaCppImageAttributes["n_parallel"]; ok && !tfNParallel.IsNull() {
					var nParallelBigFloat big.Float
					tfNParallel.As(&nParallelBigFloat)
					hfNParallel, _ := nParallelBigFloat.Int(nil)
					nParallelIntPrimitive := int(hfNParallel.Int64())
					output.Model.Image.LlamaCpp.NParallel = nParallelIntPrimitive
				}
				if tfPooling, ok := llamaCppImageAttributes["pooling"]; ok && !tfPooling.IsNull() {
					var poolingType string
					tfPooling.As(&poolingType)
					hfPoolingType := huggingface.PoolingType(poolingType)
					output.Model.Image.LlamaCpp.Pooling = &hfPoolingType
				}
				if tfPort, ok := llamaCppImageAttributes["port"]; ok && !tfPort.IsNull() {
					var portBigFloat big.Float
					tfPort.As(&portBigFloat)
					portInt, _ := portBigFloat.Int(nil)
					portIntPrimitive := int(portInt.Int64())
					output.Model.Image.LlamaCpp.Port = portIntPrimitive
				}
				if tfThreadsHttp, ok := llamaCppImageAttributes["threads_http"]; ok && !tfThreadsHttp.IsNull() {
					var threadsHttpBigFloat big.Float
					tfThreadsHttp.As(&threadsHttpBigFloat)
					threadsHttpInt, _ := threadsHttpBigFloat.Int(nil)
					threadsHttpIntPrimitive := int(threadsHttpInt.Int64())
					output.Model.Image.LlamaCpp.ThreadsHttp = &threadsHttpIntPrimitive
				}
				if tfUrl, ok := llamaCppImageAttributes["url"]; ok && !tfUrl.IsNull() {
					tfUrl.As(&output.Model.Image.LlamaCpp.URL)
				}
				if tfVariant, ok := llamaCppImageAttributes["variant"]; ok && !tfVariant.IsNull() {
					tfVariant.As(&output.Model.Image.LlamaCpp.Variant)
				}
			}
			if tfCustom, ok := imageAttributes["custom"]; ok && !tfCustom.IsNull() && tfCustom.IsKnown() {
				output.Model.Image.Custom = &huggingface.CustomImage{}
				var customImageAttributes map[string]tftypes.Value
				tfCustom.As(&customImageAttributes)

				if tfCredentials, ok := customImageAttributes["credentials"]; ok && !tfCredentials.IsNull() {
					output.Model.Image.Custom.Credentials = &huggingface.Credentials{}
					var credentialsAttributes map[string]tftypes.Value
					tfCredentials.As(&credentialsAttributes)

					if tfUsername, ok := credentialsAttributes["username"]; ok && !tfUsername.IsNull() {
						tfUsername.As(&output.Model.Image.Custom.Credentials.Username)
					}
					if tfPassword, ok := credentialsAttributes["password"]; ok && !tfPassword.IsNull() {
						tfPassword.As(&output.Model.Image.Custom.Credentials.Password)
					}
				}
				if tfHealthRoute, ok := customImageAttributes["health_route"]; ok && !tfHealthRoute.IsNull() {
					var healthRoute string
					tfHealthRoute.As(&healthRoute)
					output.Model.Image.LlamaCpp.HealthRoute = &healthRoute
				}
				if tfPort, ok := customImageAttributes["port"]; ok && !tfPort.IsNull() {
					var portBigFloat big.Float
					tfPort.As(&portBigFloat)
					portInt, _ := portBigFloat.Int(nil)
					portIntPrimitive := int(portInt.Int64())
					output.Model.Image.Custom.Port = portIntPrimitive
				}
				if tfUrl, ok := customImageAttributes["url"]; ok && !tfUrl.IsNull() {
					tfUrl.As(&output.Model.Image.Custom.URL)
				}
			}
		}
	}

	if !input.Tags.IsNull() {
		tags, _ := input.Tags.ToTerraformValue(ctx)
		tags.As(&output.Tags)
	}
	if !input.CacheHttpResponses.IsNull() {
		cacheHttpResponse := input.CacheHttpResponses.ValueBool()
		output.CacheHttpResponses = &cacheHttpResponse
	}
	if !input.ExperimentalFeatures.IsNull() {
		output.ExperimentalFeatures = &huggingface.ExperimentalFeatures{}
		experimentalFeaturesAttributes := input.ExperimentalFeatures.Attributes()

		if cacheHttpResponse, ok := experimentalFeaturesAttributes["cache_http_response"]; ok && !cacheHttpResponse.IsNull() {
			tfCacheHttpResponse, _ := cacheHttpResponse.ToTerraformValue(ctx)
			tfCacheHttpResponse.As(output.ExperimentalFeatures.CacheHttpResponses)
		}
		if kvRouter, ok := experimentalFeaturesAttributes["kv_router"]; ok && !kvRouter.IsNull() {
			output.ExperimentalFeatures.KvRouter = &huggingface.KvRouter{}

			tfKvRouter, _ := kvRouter.ToTerraformValue(ctx)
			var kvRouterAttributes map[string]tftypes.Value
			tfKvRouter.As(&kvRouterAttributes)

			if tfTag, ok := kvRouterAttributes["tag"]; ok {
				tfTag.As(output.ExperimentalFeatures.KvRouter.Tag)
			}
		}
	}

	if !input.PrivateService.IsNull() && !input.PrivateService.IsUnknown() {
		output.PrivateService = &huggingface.EndpointPrivateService{}
		privateServiceAttributes := input.PrivateService.Attributes()

		if accountId, ok := privateServiceAttributes["account_id"]; ok && !accountId.IsNull() {
			tfAccountId, _ := accountId.ToTerraformValue(ctx)
			tfAccountId.As(output.PrivateService.AccountID)
		}
		if shared, ok := privateServiceAttributes["shared"]; ok && !shared.IsNull() {
			tfShared, _ := shared.ToTerraformValue(ctx)
			tfShared.As(output.PrivateService.Shared)
		}
	}
	if !input.Route.IsNull() && !input.Route.IsUnknown() {
		output.Route = &huggingface.RouteSpec{}
		routeAttributes := input.Route.Attributes()

		if domain, ok := routeAttributes["domain"]; ok && !domain.IsNull() {
			tfDomain, _ := domain.ToTerraformValue(ctx)
			tfDomain.As(output.Route.Domain)
		}
		if path, ok := routeAttributes["path"]; ok && !path.IsNull() {
			tfPath, _ := path.ToTerraformValue(ctx)
			tfPath.As(output.Route.Path)
		}
	}

	return
}

func FromPlanToEndpointUpdate(
	ctx context.Context,
	input *states.EndpointResourceState,
) (output huggingface.EndpointUpdate) {
	output = huggingface.EndpointUpdate{}

	// Root properties
	if !input.Type.IsNull() {
		typ := huggingface.EndpointType(input.Type.ValueString())
		output.Type = &typ
	}

	// Compute
	if !input.Compute.IsNull() {
		output.Compute = &huggingface.EndpointComputeUpdate{}

		computeAttributes := input.Compute.Attributes()
		if accelerator, ok := computeAttributes["accelerator"]; ok && !accelerator.IsNull() {
			tfAccelerator, _ := accelerator.ToTerraformValue(ctx)
			var hfAccelerator string
			tfAccelerator.As(&hfAccelerator)
			pAccelerator := huggingface.AcceleratorType(hfAccelerator)
			output.Compute.Accelerator = &pAccelerator
		}
		if instanceType, ok := computeAttributes["instance_type"]; ok && !instanceType.IsNull() {
			tfInstanceType, _ := instanceType.ToTerraformValue(ctx)
			tfInstanceType.As(&output.Compute.InstanceType)
		}
		if instanceSize, ok := computeAttributes["instance_size"]; ok && !instanceSize.IsNull() {
			tfInstanceType, _ := instanceSize.ToTerraformValue(ctx)
			tfInstanceType.As(&output.Compute.InstanceSize)
		}
		if scaling, ok := computeAttributes["scaling"]; ok && !scaling.IsNull() {
			output.Compute.Scaling = &huggingface.EndpointScalingUpdate{}

			tfScaling, _ := scaling.ToTerraformValue(ctx)
			var scalingAttributes map[string]tftypes.Value
			tfScaling.As(&scalingAttributes)

			if tfMinReplica, ok := scalingAttributes["min_replica"]; ok && !tfMinReplica.IsNull() {
				var minReplicaBigFloat big.Float
				tfMinReplica.As(&minReplicaBigFloat)
				minReplicaInt, _ := minReplicaBigFloat.Int(nil)
				pMinReplica := int(minReplicaInt.Int64())
				output.Compute.Scaling.MinReplica = &pMinReplica
			} else {
				pMinReplica := 0
				output.Compute.Scaling.MinReplica = &pMinReplica
			}
			if tfMaxReplica, ok := scalingAttributes["max_replica"]; ok && !tfMaxReplica.IsNull() {
				var maxReplicaBigFloat big.Float
				tfMaxReplica.As(&maxReplicaBigFloat)
				maxReplicaInt, _ := maxReplicaBigFloat.Int(nil)
				pMaxReplica := int(maxReplicaInt.Int64())
				output.Compute.Scaling.MaxReplica = &pMaxReplica
			} else {
				pMaxReplica := 1
				output.Compute.Scaling.MaxReplica = &pMaxReplica
			}
			if tfMetric, ok := scalingAttributes["metric"]; ok && !tfMetric.IsNull() {
				tfMetric.As(&output.Compute.Scaling.Metric)
			}
			if tfScaleToZeroTimeout, ok := scalingAttributes["scale_to_zero_timeout"]; ok && !tfScaleToZeroTimeout.IsNull() && tfScaleToZeroTimeout.IsKnown() {
				var scaleToZeroTimeoutBigFloat big.Float
				tfScaleToZeroTimeout.As(&scaleToZeroTimeoutBigFloat)
				scaleToZeroTimeoutInt, _ := scaleToZeroTimeoutBigFloat.Int(nil)
				scaleToZeroTimeoutIntPrimitive := int(scaleToZeroTimeoutInt.Int64())
				output.Compute.Scaling.ScaleToZeroTimeout = &scaleToZeroTimeoutIntPrimitive
			}
			if tfThreshold, ok := scalingAttributes["threshold"]; ok && !tfThreshold.IsNull() && tfThreshold.IsKnown() {
				var thresholdBigFloat big.Float
				tfThreshold.As(&thresholdBigFloat)
				thresholdFloatPrimitive, _ := thresholdBigFloat.Float64()
				output.Compute.Scaling.Threshold = &thresholdFloatPrimitive
			}
			if tfMeasure, ok := scalingAttributes["measure"]; ok && !tfMeasure.IsNull() {
				output.Compute.Scaling.Measure = &huggingface.ScalingMeasure{}

				var measureAttributes map[string]tftypes.Value
				tfMeasure.As(&measureAttributes)

				if tfHardwareUsage, ok := measureAttributes["hardware_usage"]; ok && !tfHardwareUsage.IsNull() {
					var hardwareUsageBigFloat big.Float
					tfHardwareUsage.As(&hardwareUsageBigFloat)
					hardwareUsageFloatPrimitive, _ := hardwareUsageBigFloat.Float64()
					output.Compute.Scaling.Measure.HardwareUsage = &hardwareUsageFloatPrimitive
				}
				if tfPendingRequests, ok := measureAttributes["pending_requests"]; ok && !tfPendingRequests.IsNull() {
					var pendingRequestsBigFloat big.Float
					tfPendingRequests.As(&pendingRequestsBigFloat)
					pendintRequestsFloatPrimitive, _ := pendingRequestsBigFloat.Float64()
					output.Compute.Scaling.Measure.PendingRequests = &pendintRequestsFloatPrimitive
				}
			}
		}
	}

	// Model
	if !input.Model.IsNull() {
		output.Model = &huggingface.EndpointModelUpdate{}

		modelAttributes := input.Model.Attributes()
		if repository, ok := modelAttributes["repository"]; ok && !repository.IsNull() {
			tfRepository, _ := repository.ToTerraformValue(ctx)
			tfRepository.As(&output.Model.Repository)
		}
		if framework, ok := modelAttributes["framework"]; ok && !framework.IsNull() {
			tfFramework, _ := framework.ToTerraformValue(ctx)
			var hfFramework string
			tfFramework.As(&hfFramework)
			pFramework := huggingface.EndpointFramework(hfFramework)
			output.Model.Framework = &pFramework
		}
		if task, ok := modelAttributes["task"]; ok && !task.IsNull() {
			tfTask, _ := task.ToTerraformValue(ctx)
			tfTask.As(&output.Model.Task)
		}
		if image, ok := modelAttributes["image"]; ok && !image.IsNull() {
			output.Model.Image = &huggingface.EndpointModelImage{}

			tfImage, _ := image.ToTerraformValue(ctx)
			var imageAttributes map[string]tftypes.Value
			tfImage.As(&imageAttributes)

			if tfHuggingface, ok := imageAttributes["huggingface"]; ok && !tfHuggingface.IsNull() && tfHuggingface.IsKnown() {
				output.Model.Image.HuggingFace = &huggingface.HuggingFaceImage{}
			}
			if tfHuggingfaceNeuron, ok := imageAttributes["huggingface_neuron"]; ok && !tfHuggingfaceNeuron.IsNull() && tfHuggingfaceNeuron.IsKnown() {
				output.Model.Image.HuggingFaceNeuron = &huggingface.HuggingFaceNeuronImage{}
				var huggingfaceNeuronAttributes map[string]tftypes.Value
				tfHuggingfaceNeuron.As(&huggingfaceNeuronAttributes)

				if tfBatchSize, ok := huggingfaceNeuronAttributes["batch_size"]; ok && !tfBatchSize.IsNull() {
					var batchSizeBigFloat big.Float
					tfBatchSize.As(&batchSizeBigFloat)
					batchSizeInt, _ := batchSizeBigFloat.Int(nil)
					batchSizeIntPrimitive := int(batchSizeInt.Int64())
					output.Model.Image.HuggingFaceNeuron.BatchSize = &batchSizeIntPrimitive
				}
				if tfNeuronCache, ok := huggingfaceNeuronAttributes["neuron_cache"]; ok && !tfNeuronCache.IsNull() {
					tfNeuronCache.As(&output.Model.Image.HuggingFaceNeuron.NeuronCache)
				}
				if tfSequenceLength, ok := huggingfaceNeuronAttributes["sequence_length"]; ok && !tfSequenceLength.IsNull() {
					var sequenceLengthBigFloat big.Float
					tfSequenceLength.As(&sequenceLengthBigFloat)
					sequenceLengthInt, _ := sequenceLengthBigFloat.Int(nil)
					sequenceLengthIntPrimitive := int(sequenceLengthInt.Int64())
					output.Model.Image.HuggingFaceNeuron.SequenceLength = &sequenceLengthIntPrimitive
				}
			}
			if tfTgi, ok := imageAttributes["tgi"]; ok && !tfTgi.IsNull() && tfTgi.IsKnown() {
				output.Model.Image.TGI = &huggingface.TGIImage{}
				var tgiAttributes map[string]tftypes.Value
				tfTgi.As(&tgiAttributes)

				if tfDisableCustomKernels, ok := tgiAttributes["disable_custom_kernels"]; ok && !tfDisableCustomKernels.IsNull() {
					tfDisableCustomKernels.As(&output.Model.Image.TGI.DisableCustomKernels)
				}
				if tfHealthRoute, ok := tgiAttributes["health_route"]; ok && !tfHealthRoute.IsNull() {
					tfHealthRoute.As(&output.Model.Image.TGI.HealthRoute)
				}
				if tfMaxBatchPrefillTokens, ok := tgiAttributes["max_batch_prefill_tokens"]; ok && !tfMaxBatchPrefillTokens.IsNull() {
					tfMaxBatchPrefillTokens.As(&output.Model.Image.TGI.MaxBatchPrefillTokens)
				}
				if tfMaxBatchTotalTokens, ok := tgiAttributes["max_batch_total_tokens"]; ok && !tfMaxBatchTotalTokens.IsNull() {
					tfMaxBatchTotalTokens.As(&output.Model.Image.TGI.MaxBatchTotalTokens)
				}
				if tfMaxInputLength, ok := tgiAttributes["max_input_length"]; ok && !tfMaxInputLength.IsNull() {
					var maxInputLengthBigFloat big.Float
					tfMaxInputLength.As(&maxInputLengthBigFloat)
					maxInputLengthInt, _ := maxInputLengthBigFloat.Int(nil)
					maxInputLengthIntPrimitive := int(maxInputLengthInt.Int64())
					output.Model.Image.TGI.MaxInputLength = &maxInputLengthIntPrimitive
				}
				if tfMaxTotalTokens, ok := tgiAttributes["max_total_tokens"]; ok && !tfMaxTotalTokens.IsNull() {
					var maxTotalTokensBigFloat big.Float
					tfMaxTotalTokens.As(&maxTotalTokensBigFloat)
					maxTotalTokensInt, _ := maxTotalTokensBigFloat.Int(nil)
					maxTotalTokensIntPrimitive := int(maxTotalTokensInt.Int64())
					output.Model.Image.TGI.MaxTotalTokens = &maxTotalTokensIntPrimitive
				}
				if tfPort, ok := tgiAttributes["port"]; ok && !tfPort.IsNull() {
					var portBigFloat big.Float
					tfPort.As(&portBigFloat)
					portInt, _ := portBigFloat.Int(nil)
					portIntPrimitive := int(portInt.Int64())
					output.Model.Image.TGI.Port = portIntPrimitive
				}
				if tfQuantize, ok := tgiAttributes["quantize"]; ok && !tfQuantize.IsNull() {
					var quantizeString string
					tfQuantize.As(&quantizeString)
					quantizeType := huggingface.QuantizeType(quantizeString)
					output.Model.Image.TGI.Quantize = &quantizeType
				}
				if tfUrl, ok := tgiAttributes["url"]; ok && !tfUrl.IsNull() {
					tfUrl.As(&output.Model.Image.TGI.URL)
				}
			}
			if tfTgiNeuron, ok := imageAttributes["tgi_neuron"]; ok && !tfTgiNeuron.IsNull() && tfTgiNeuron.IsKnown() {
				output.Model.Image.TGINeuron = &huggingface.TGINeuronImage{}
				var tgiNeuronAttributes map[string]tftypes.Value
				tfTgiNeuron.As(&tgiNeuronAttributes)

				if tfHealthRoute, ok := tgiNeuronAttributes["health_route"]; ok && !tfHealthRoute.IsNull() {
					var healthRoute string
					tfHealthRoute.As(&healthRoute)
					output.Model.Image.TGINeuron.HealthRoute = &healthRoute
				}
				if tfHfAutoCastType, ok := tgiNeuronAttributes["hf_auto_cast_type"]; ok && !tfHfAutoCastType.IsNull() {
					var hfAutoCastType string
					tfHfAutoCastType.As(&hfAutoCastType)
					hfHfAutoCastType := huggingface.AutoCastType(hfAutoCastType)
					output.Model.Image.TGINeuron.HfAutoCastType = &hfHfAutoCastType
				}
				if tfHfNumCores, ok := tgiNeuronAttributes["hf_num_cores"]; ok && !tfHfNumCores.IsNull() {
					var hfNumCoresBigFloat big.Float
					tfHfNumCores.As(&hfNumCoresBigFloat)
					hfNumCoresInt, _ := hfNumCoresBigFloat.Int(nil)
					hfNumCoresIntPrimitive := int(hfNumCoresInt.Int64())
					output.Model.Image.TGINeuron.HfNumCores = &hfNumCoresIntPrimitive
				}
				if tfMaxBatchPrefillTokens, ok := tgiNeuronAttributes["max_batch_prefill_tokens"]; ok && !tfMaxBatchPrefillTokens.IsNull() {
					var maxBatchPrefillTokensBigFloat big.Float
					tfMaxBatchPrefillTokens.As(&maxBatchPrefillTokensBigFloat)
					hfMaxBatchPrefillTokens, _ := maxBatchPrefillTokensBigFloat.Int(nil)
					maxBatchPrefillTokensIntPrimitive := int(hfMaxBatchPrefillTokens.Int64())
					output.Model.Image.TGINeuron.MaxBatchPrefillTokens = &maxBatchPrefillTokensIntPrimitive
				}
				if tfMaxBatchTotalTokens, ok := tgiNeuronAttributes["max_batch_total_tokens"]; ok && !tfMaxBatchTotalTokens.IsNull() {
					var maxBatchTotalTokensBigFloat big.Float
					tfMaxBatchTotalTokens.As(&maxBatchTotalTokensBigFloat)
					hfMaxBatchTotalTokens, _ := maxBatchTotalTokensBigFloat.Int(nil)
					maxBatchTotalTokensIntPrimitive := int(hfMaxBatchTotalTokens.Int64())
					output.Model.Image.TGINeuron.MaxBatchTotalTokens = &maxBatchTotalTokensIntPrimitive
				}
				if tfMaxInputLength, ok := tgiNeuronAttributes["max_input_length"]; ok && !tfMaxInputLength.IsNull() {
					var maxInputLengthBigFloat big.Float
					tfMaxInputLength.As(&maxInputLengthBigFloat)
					hfMaxInputLength, _ := maxInputLengthBigFloat.Int(nil)
					hfMaxInputLengthIntPrimitive := int(hfMaxInputLength.Int64())
					output.Model.Image.TGINeuron.MaxInputLength = &hfMaxInputLengthIntPrimitive
				}
				if tfMaxTotalTokens, ok := tgiNeuronAttributes["max_total_tokens"]; ok && !tfMaxTotalTokens.IsNull() {
					var maxTotalTokensBigFloat big.Float
					tfMaxTotalTokens.As(&maxTotalTokensBigFloat)
					hfMaxTotalTokens, _ := maxTotalTokensBigFloat.Int(nil)
					maxTotalTokensIntPrimitive := int(hfMaxTotalTokens.Int64())
					output.Model.Image.TGINeuron.MaxTotalTokens = &maxTotalTokensIntPrimitive
				}
				if tfPort, ok := tgiNeuronAttributes["port"]; ok && !tfPort.IsNull() {
					var portBigFloat big.Float
					tfPort.As(&portBigFloat)
					portInt, _ := portBigFloat.Int(nil)
					portIntPrimitive := int(portInt.Int64())
					output.Model.Image.TGINeuron.Port = portIntPrimitive
				}
				if tfUrl, ok := tgiNeuronAttributes["url"]; ok && !tfUrl.IsNull() {
					tfUrl.As(&output.Model.Image.TGINeuron.URL)
				}
			}
			if tfTei, ok := imageAttributes["tei"]; ok && !tfTei.IsNull() && tfTei.IsKnown() {
				output.Model.Image.TEI = &huggingface.TEIImage{}
				var teiAttributes map[string]tftypes.Value
				tfTei.As(&teiAttributes)

				if tfHealthRoute, ok := teiAttributes["health_route"]; ok && !tfHealthRoute.IsNull() {
					var healthRoute string
					tfHealthRoute.As(&healthRoute)
					output.Model.Image.TEI.HealthRoute = &healthRoute
				}
				if tfMaxBatchTokens, ok := teiAttributes["max_batch_tokens"]; ok && !tfMaxBatchTokens.IsNull() {
					var maxBatchTokensBigFloat big.Float
					tfMaxBatchTokens.As(&maxBatchTokensBigFloat)
					hfMaxBatchTokens, _ := maxBatchTokensBigFloat.Int(nil)
					maxBatchTokensIntPrimitive := int(hfMaxBatchTokens.Int64())
					output.Model.Image.TEI.MaxBatchTokens = &maxBatchTokensIntPrimitive
				}
				if tfMaxConcurrentRequests, ok := teiAttributes["max_concurrent_requests"]; ok && !tfMaxConcurrentRequests.IsNull() {
					var maxConcurrentRequestsBigFloat big.Float
					tfMaxConcurrentRequests.As(&maxConcurrentRequestsBigFloat)
					hfMaxConcurrentRequests, _ := maxConcurrentRequestsBigFloat.Int(nil)
					maxConcurrentRequestsIntPrimitive := int(hfMaxConcurrentRequests.Int64())
					output.Model.Image.TEI.MaxConcurrentRequests = &maxConcurrentRequestsIntPrimitive
				}
				if tfPooling, ok := teiAttributes["pooling"]; ok && !tfPooling.IsNull() {
					var poolingType string
					tfPooling.As(&poolingType)
					hfPoolingType := huggingface.PoolingType(poolingType)
					output.Model.Image.TEI.Pooling = &hfPoolingType
				}
				if tfPort, ok := teiAttributes["port"]; ok && !tfPort.IsNull() {
					var portBigFloat big.Float
					tfPort.As(&portBigFloat)
					portInt, _ := portBigFloat.Int(nil)
					portIntPrimitive := int(portInt.Int64())
					output.Model.Image.TEI.Port = portIntPrimitive
				}
				if tfUrl, ok := teiAttributes["url"]; ok && !tfUrl.IsNull() {
					tfUrl.As(&output.Model.Image.TEI.URL)
				}
			}
			if tfLlamacpp, ok := imageAttributes["llamacpp"]; ok && !tfLlamacpp.IsNull() && tfLlamacpp.IsKnown() {
				output.Model.Image.LlamaCpp = &huggingface.LlamaCppImage{}
				var llamaCppImageAttributes map[string]tftypes.Value
				tfLlamacpp.As(&llamaCppImageAttributes)

				if tfCtxSize, ok := llamaCppImageAttributes["ctx_size"]; ok && !tfCtxSize.IsNull() {
					var ctxSizeBigFloat big.Float
					tfCtxSize.As(&ctxSizeBigFloat)
					hfCtxSize, _ := ctxSizeBigFloat.Int(nil)
					ctxSizeIntPrimitive := int(hfCtxSize.Int64())
					output.Model.Image.LlamaCpp.CtxSize = ctxSizeIntPrimitive
				}
				if tfHealthRoute, ok := llamaCppImageAttributes["health_route"]; ok && !tfHealthRoute.IsNull() {
					var healthRoute string
					tfHealthRoute.As(&healthRoute)
					output.Model.Image.LlamaCpp.HealthRoute = &healthRoute
				}
				if tfMode, ok := llamaCppImageAttributes["mode"]; ok && !tfMode.IsNull() {
					var modeType string
					tfMode.As(&modeType)
					hfModeType := huggingface.ModelMode(modeType)
					output.Model.Image.LlamaCpp.Mode = &hfModeType
				}
				if tfModelPath, ok := llamaCppImageAttributes["model_path"]; ok && !tfModelPath.IsNull() {
					tfModelPath.As(&output.Model.Image.LlamaCpp.ModelPath)
				}
				if tfNGpuLayers, ok := llamaCppImageAttributes["n_gpu_layers"]; ok && !tfNGpuLayers.IsNull() {
					var nGpuLayersBigFloat big.Float
					tfNGpuLayers.As(&nGpuLayersBigFloat)
					hfNGpuLayers, _ := nGpuLayersBigFloat.Int(nil)
					nGpuLayersIntPrimitive := int(hfNGpuLayers.Int64())
					output.Model.Image.LlamaCpp.NGpuLayers = nGpuLayersIntPrimitive
				}
				if tfNParallel, ok := llamaCppImageAttributes["n_parallel"]; ok && !tfNParallel.IsNull() {
					var nParallelBigFloat big.Float
					tfNParallel.As(&nParallelBigFloat)
					hfNParallel, _ := nParallelBigFloat.Int(nil)
					nParallelIntPrimitive := int(hfNParallel.Int64())
					output.Model.Image.LlamaCpp.NParallel = nParallelIntPrimitive
				}
				if tfPooling, ok := llamaCppImageAttributes["pooling"]; ok && !tfPooling.IsNull() {
					var poolingType string
					tfPooling.As(&poolingType)
					hfPoolingType := huggingface.PoolingType(poolingType)
					output.Model.Image.LlamaCpp.Pooling = &hfPoolingType
				}
				if tfPort, ok := llamaCppImageAttributes["port"]; ok && !tfPort.IsNull() {
					var portBigFloat big.Float
					tfPort.As(&portBigFloat)
					portInt, _ := portBigFloat.Int(nil)
					portIntPrimitive := int(portInt.Int64())
					output.Model.Image.LlamaCpp.Port = portIntPrimitive
				}
				if tfThreadsHttp, ok := llamaCppImageAttributes["threads_http"]; ok && !tfThreadsHttp.IsNull() {
					var threadsHttpBigFloat big.Float
					tfThreadsHttp.As(&threadsHttpBigFloat)
					threadsHttpInt, _ := threadsHttpBigFloat.Int(nil)
					threadsHttpIntPrimitive := int(threadsHttpInt.Int64())
					output.Model.Image.LlamaCpp.ThreadsHttp = &threadsHttpIntPrimitive
				}
				if tfUrl, ok := llamaCppImageAttributes["url"]; ok && !tfUrl.IsNull() {
					tfUrl.As(&output.Model.Image.LlamaCpp.URL)
				}
				if tfVariant, ok := llamaCppImageAttributes["variant"]; ok && !tfVariant.IsNull() {
					tfVariant.As(&output.Model.Image.LlamaCpp.Variant)
				}
			}
			if tfCustom, ok := imageAttributes["custom"]; ok && !tfCustom.IsNull() && tfCustom.IsKnown() {
				output.Model.Image.Custom = &huggingface.CustomImage{}
				var customImageAttributes map[string]tftypes.Value
				tfCustom.As(&customImageAttributes)

				if tfCredentials, ok := customImageAttributes["credentials"]; ok && !tfCredentials.IsNull() {
					output.Model.Image.Custom.Credentials = &huggingface.Credentials{}
					var credentialsAttributes map[string]tftypes.Value
					tfCredentials.As(&credentialsAttributes)

					if tfUsername, ok := credentialsAttributes["username"]; ok && !tfUsername.IsNull() {
						tfUsername.As(&output.Model.Image.Custom.Credentials.Username)
					}
					if tfPassword, ok := credentialsAttributes["password"]; ok && !tfPassword.IsNull() {
						tfPassword.As(&output.Model.Image.Custom.Credentials.Password)
					}
				}
				if tfHealthRoute, ok := customImageAttributes["health_route"]; ok && !tfHealthRoute.IsNull() {
					var healthRoute string
					tfHealthRoute.As(&healthRoute)
					output.Model.Image.LlamaCpp.HealthRoute = &healthRoute
				}
				if tfPort, ok := customImageAttributes["port"]; ok && !tfPort.IsNull() {
					var portBigFloat big.Float
					tfPort.As(&portBigFloat)
					portInt, _ := portBigFloat.Int(nil)
					portIntPrimitive := int(portInt.Int64())
					output.Model.Image.Custom.Port = portIntPrimitive
				}
				if tfUrl, ok := customImageAttributes["url"]; ok && !tfUrl.IsNull() {
					tfUrl.As(&output.Model.Image.Custom.URL)
				}
			}
		}
	}

	if !input.Tags.IsNull() {
		tags, _ := input.Tags.ToTerraformValue(ctx)
		tags.As(&output.Tags)
	}
	if !input.ExperimentalFeatures.IsNull() {
		output.ExperimentalFeatures = &huggingface.ExperimentalFeatures{}
		experimentalFeaturesAttributes := input.ExperimentalFeatures.Attributes()

		if cacheHttpResponse, ok := experimentalFeaturesAttributes["cache_http_response"]; ok && !cacheHttpResponse.IsNull() {
			tfCacheHttpResponse, _ := cacheHttpResponse.ToTerraformValue(ctx)
			tfCacheHttpResponse.As(output.ExperimentalFeatures.CacheHttpResponses)
		}
		if kvRouter, ok := experimentalFeaturesAttributes["kv_router"]; ok && !kvRouter.IsNull() {
			output.ExperimentalFeatures.KvRouter = &huggingface.KvRouter{}

			tfKvRouter, _ := kvRouter.ToTerraformValue(ctx)
			var kvRouterAttributes map[string]tftypes.Value
			tfKvRouter.As(&kvRouterAttributes)

			if tfTag, ok := kvRouterAttributes["tag"]; ok {
				tfTag.As(output.ExperimentalFeatures.KvRouter.Tag)
			}
		}
	}

	if !input.Route.IsNull() && !input.Route.IsUnknown() {
		output.Route = &huggingface.RouteSpec{}
		routeAttributes := input.Route.Attributes()

		if domain, ok := routeAttributes["domain"]; ok && !domain.IsNull() {
			tfDomain, _ := domain.ToTerraformValue(ctx)
			tfDomain.As(output.Route.Domain)
		}
		if path, ok := routeAttributes["path"]; ok && !path.IsNull() {
			tfPath, _ := path.ToTerraformValue(ctx)
			tfPath.As(output.Route.Path)
		}
	}

	return
}

func FromProviderToModel(
	ctx context.Context,
	input *huggingface.EndpointWithStatus,
) (output models.Endpoint, diags diag.Diagnostics) {
	var diag diag.Diagnostics

	// Root
	output.Name = types.StringValue(input.Name)
	output.Type = types.StringValue(string(input.Type))

	if input.CacheHttpResponses != nil {
		output.CacheHttpResponses = types.BoolValue(*input.CacheHttpResponses)
	}

	// 	Cloud Provider
	endpointCloudProvider := models.EndpointCloudProvider{
		Vendor: types.StringValue(input.Provider.Vendor),
		Region: types.StringValue(input.Provider.Region),
	}
	output.CloudProvider, diag = types.ObjectValueFrom(ctx, endpointCloudProvider.AttributeTypes(), endpointCloudProvider)
	diags.Append(diag...)
	if diags.HasError() {
		return
	}

	// Compute
	endpointCompute := models.EndpointCompute{
		ID:           types.StringValue(*input.Compute.ID),
		Accelerator:  types.StringValue(string(input.Compute.Accelerator)),
		InstanceType: types.StringValue(input.Compute.InstanceType),
		InstanceSize: types.StringValue(input.Compute.InstanceSize),
	}

	endpointComputeScaling := models.EndpointComputeScaling{
		MinReplica: types.Int32Value(int32(input.Compute.Scaling.MinReplica)),
		MaxReplica: types.Int32Value(int32(input.Compute.Scaling.MaxReplica)),
	}

	if input.Compute.Scaling.ScaleToZeroTimeout != nil {
		endpointComputeScaling.ScaleToZeroTimeout = types.Int32Value(int32(*input.Compute.Scaling.ScaleToZeroTimeout))
	}
	if input.Compute.Scaling.Threshold != nil {
		threshold := *input.Compute.Scaling.Threshold
		endpointComputeScaling.Threshold = types.Float64Value(threshold)
	}
	if input.Compute.Scaling.Metric != nil {
		endpointScalingMetric := *input.Compute.Scaling.Metric
		endpointComputeScaling.Metric = types.StringValue(string(endpointScalingMetric))
	}
	if input.Compute.Scaling.Measure != nil {
		endpointComputeScalingMeasure := models.EndpointComputeScalingMeasure{}

		if input.Compute.Scaling.Measure.HardwareUsage != nil {
			hardwareUsage := *input.Compute.Scaling.Measure.HardwareUsage
			endpointComputeScalingMeasure.HardwareUsage = types.Float64Value(hardwareUsage)
		}
		if input.Compute.Scaling.Measure.PendingRequests != nil {
			pendingRequests := *input.Compute.Scaling.Measure.PendingRequests
			endpointComputeScalingMeasure.PendingRequests = types.Float64Value(pendingRequests)
		}

		endpointComputeScaling.Measure, diag = types.ObjectValueFrom(ctx, endpointComputeScalingMeasure.AttributeTypes(), endpointComputeScalingMeasure)
		diags.Append(diag...)
		if diags.HasError() {
			return
		}
	}

	endpointCompute.Scaling, diag = types.ObjectValueFrom(ctx, endpointComputeScaling.AttributeTypes(), endpointComputeScaling)
	diags.Append(diag...)
	if diags.HasError() {
		return
	}

	output.Compute, diag = types.ObjectValueFrom(ctx, endpointCompute.AttributeTypes(), endpointCompute)
	diags.Append(diag...)
	if diags.HasError() {
		return
	}

	// Model
	endpointModel := models.Model{
		Repository: types.StringValue(input.Model.Repository),
		Framework:  types.StringValue(string(input.Model.Framework)),
		Task:       types.StringValue(string(input.Model.Task)),
	}

	modelImage := models.ModelImage{}

	if input.Model.Image.HuggingFace != nil {
		huggingFaceImage := models.ModelImageHuggingface{}

		modelImage.HuggingFace, diag = types.ObjectValueFrom(ctx, huggingFaceImage.AttributeTypes(), huggingFaceImage)
		diags.Append(diag...)
		if diags.HasError() {
			return
		}
	} else {
		modelImage.HuggingFace, diag = types.ObjectValueFrom(ctx, models.ModelImageHuggingface{}.AttributeTypes(), models.ModelImageHuggingface{})
		diags.Append(diag...)
		if diags.HasError() {
			return
		}
	}

	if input.Model.Image.HuggingFaceNeuron != nil {
		huggingFaceNeuronImage := models.ModelImageHuggingfaceNeuron{}

		if input.Model.Image.HuggingFaceNeuron.BatchSize != nil {
			huggingFaceNeuronImage.BatchSize = types.Int32Value(int32(*input.Model.Image.HuggingFaceNeuron.BatchSize))
		}

		huggingFaceNeuronImage.NeuronCache = types.StringValue(input.Model.Image.HuggingFaceNeuron.NeuronCache)

		if input.Model.Image.HuggingFaceNeuron.SequenceLength != nil {
			huggingFaceNeuronImage.SequenceLength = types.Int32Value(int32(*input.Model.Image.HuggingFaceNeuron.SequenceLength))
		}

		modelImage.HuggingFaceNeuron, diag = types.ObjectValueFrom(ctx, huggingFaceNeuronImage.AttributeTypes(), huggingFaceNeuronImage)
		diags.Append(diag...)
		if diags.HasError() {
			return
		}
	} else {
		modelImage.HuggingFaceNeuron, diag = types.ObjectValueFrom(ctx, models.ModelImageHuggingfaceNeuron{}.AttributeTypes(), models.ModelImageHuggingfaceNeuron{
			NeuronCache:    types.StringNull(),
			BatchSize:      types.Int32Null(),
			SequenceLength: types.Int32Null(),
		})
		diags.Append(diag...)
		if diags.HasError() {
			return
		}
	}

	if input.Model.Image.TGI != nil {
		huggingFaceTgiImage := models.ModelImageTgi{}

		if input.Model.Image.TGI.HealthRoute != nil {
			huggingFaceTgiImage.HealthRoute = types.StringValue(*input.Model.Image.TGI.HealthRoute)
		}

		huggingFaceTgiImage.Port = types.Int32Value(int32(input.Model.Image.TGI.Port))
		huggingFaceTgiImage.Url = types.StringValue(input.Model.Image.TGI.URL)

		if input.Model.Image.TGI.MaxBatchPrefillTokens != nil {
			huggingFaceTgiImage.MaxBatchPrefillTokens = types.Int32Value(int32(*input.Model.Image.TGI.MaxBatchPrefillTokens))
		}
		if input.Model.Image.TGI.MaxBatchTotalTokens != nil {
			huggingFaceTgiImage.MaxBatchTotalTokens = types.Int32Value(int32(*input.Model.Image.TGI.MaxBatchTotalTokens))
		}
		if input.Model.Image.TGI.MaxInputLength != nil {
			huggingFaceTgiImage.MaxInputLength = types.Int32Value(int32(*input.Model.Image.TGI.MaxInputLength))
		}
		if input.Model.Image.TGI.MaxTotalTokens != nil {
			huggingFaceTgiImage.MaxTotalTokens = types.Int32Value(int32(*input.Model.Image.TGI.MaxTotalTokens))
		}

		huggingFaceTgiImage.DisableCustomKernels = types.BoolValue(input.Model.Image.TGI.DisableCustomKernels)

		if input.Model.Image.TGI.Quantize != nil {
			huggingFaceTgiImage.Quantize = types.StringValue(string(*input.Model.Image.TGI.Quantize))
		}

		modelImage.TGI, diag = types.ObjectValueFrom(ctx, huggingFaceTgiImage.AttributeTypes(), huggingFaceTgiImage)
		diags.Append(diag...)
		if diags.HasError() {
			return
		}
	} else {
		modelImage.TGI, diag = types.ObjectValueFrom(ctx, models.ModelImageTgi{}.AttributeTypes(), models.ModelImageTgi{
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
		diags.Append(diag...)
		if diags.HasError() {
			return
		}
	}

	if input.Model.Image.TGINeuron != nil {
		huggingFaceTgiNeuronImage := models.ModelImageTgiNeuron{}

		if input.Model.Image.TGINeuron.HealthRoute != nil {
			huggingFaceTgiNeuronImage.HealthRoute = types.StringValue(*input.Model.Image.TGINeuron.HealthRoute)
		}

		huggingFaceTgiNeuronImage.Port = types.Int32Value(int32(input.Model.Image.TGINeuron.Port))
		huggingFaceTgiNeuronImage.Url = types.StringValue(input.Model.Image.TGINeuron.URL)

		if input.Model.Image.TGINeuron.MaxBatchPrefillTokens != nil {
			huggingFaceTgiNeuronImage.MaxBatchPrefillTokens = types.Int32Value(int32(*input.Model.Image.TGINeuron.MaxBatchPrefillTokens))
		}
		if input.Model.Image.TGINeuron.MaxBatchTotalTokens != nil {
			huggingFaceTgiNeuronImage.MaxBatchTotalTokens = types.Int32Value(int32(*input.Model.Image.TGI.MaxBatchTotalTokens))
		}
		if input.Model.Image.TGINeuron.MaxInputLength != nil {
			huggingFaceTgiNeuronImage.MaxInputLength = types.Int32Value(int32(*input.Model.Image.TGINeuron.MaxInputLength))
		}
		if input.Model.Image.TGINeuron.MaxTotalTokens != nil {
			huggingFaceTgiNeuronImage.MaxTotalTokens = types.Int32Value(int32(*input.Model.Image.TGINeuron.MaxTotalTokens))
		}
		if input.Model.Image.TGINeuron.MaxTotalTokens != nil {
			huggingFaceTgiNeuronImage.MaxTotalTokens = types.Int32Value(int32(*input.Model.Image.TGI.MaxTotalTokens))
		}
		if input.Model.Image.TGINeuron.HfAutoCastType != nil {
			huggingFaceTgiNeuronImage.HfAutoCastType = types.StringValue(string(*input.Model.Image.TGINeuron.HfAutoCastType))
		}
		if input.Model.Image.TGINeuron.HfNumCores != nil {
			huggingFaceTgiNeuronImage.HfNumCores = types.Int32Value(int32(*input.Model.Image.TGINeuron.HfNumCores))
		}

		modelImage.TGINeuron, diag = types.ObjectValueFrom(ctx, huggingFaceTgiNeuronImage.AttributeTypes(), huggingFaceTgiNeuronImage)
		diags.Append(diag...)
		if diags.HasError() {
			return
		}
	} else {
		modelImage.TGINeuron, diag = types.ObjectValueFrom(ctx, models.ModelImageTgiNeuron{}.AttributeTypes(), models.ModelImageTgiNeuron{
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
		diags.Append(diag...)
		if diags.HasError() {
			return
		}
	}

	if input.Model.Image.TEI != nil {
		huggingFaceTeiNeuronImage := models.ModelImageTei{}

		if input.Model.Image.TEI.HealthRoute != nil {
			huggingFaceTeiNeuronImage.HealthRoute = types.StringValue(*input.Model.Image.TGINeuron.HealthRoute)
		}

		huggingFaceTeiNeuronImage.Port = types.Int32Value(int32(input.Model.Image.TEI.Port))
		huggingFaceTeiNeuronImage.URL = types.StringValue(input.Model.Image.TEI.URL)

		if input.Model.Image.TEI.MaxBatchTokens != nil {
			huggingFaceTeiNeuronImage.MaxBatchTokens = types.Int32Value(int32(*input.Model.Image.TEI.MaxBatchTokens))
		}
		if input.Model.Image.TEI.MaxConcurrentRequests != nil {
			huggingFaceTeiNeuronImage.MaxConcurrentRequests = types.Int32Value(int32(*input.Model.Image.TEI.MaxConcurrentRequests))
		}
		if input.Model.Image.TEI.Pooling != nil {
			huggingFaceTeiNeuronImage.Pooling = types.StringValue(string(*input.Model.Image.TEI.Pooling))
		}

		modelImage.TEI, diag = types.ObjectValueFrom(ctx, huggingFaceTeiNeuronImage.AttributeTypes(), huggingFaceTeiNeuronImage)
		diags.Append(diag...)
		if diags.HasError() {
			return
		}
	} else {
		modelImage.TEI, diag = types.ObjectValueFrom(ctx, models.ModelImageTei{}.AttributeTypes(), models.ModelImageTei{
			HealthRoute:           types.StringNull(),
			Port:                  types.Int32Null(),
			URL:                   types.StringNull(),
			MaxBatchTokens:        types.Int32Null(),
			MaxConcurrentRequests: types.Int32Null(),
			Pooling:               types.StringNull(),
		})
		diags.Append(diag...)
		if diags.HasError() {
			return
		}
	}

	if input.Model.Image.LlamaCpp != nil {
		llamaCppImage := models.ModelImageLlamacpp{}

		if input.Model.Image.LlamaCpp.HealthRoute != nil {
			llamaCppImage.HealthRoute = types.StringValue(*input.Model.Image.LlamaCpp.HealthRoute)
		}

		llamaCppImage.Port = types.Int32Value(int32(input.Model.Image.LlamaCpp.Port))
		llamaCppImage.URL = types.StringValue(input.Model.Image.LlamaCpp.URL)
		llamaCppImage.CtxSize = types.Int32Value(int32(input.Model.Image.LlamaCpp.CtxSize))

		if input.Model.Image.LlamaCpp.Mode != nil {
			llamaCppImage.Mode = types.StringValue(string(*input.Model.Image.LlamaCpp.Mode))
		}

		llamaCppImage.ModelPath = types.StringValue(string(input.Model.Image.LlamaCpp.ModelPath))
		llamaCppImage.NGpuLayers = types.Int32Value(int32(input.Model.Image.LlamaCpp.NGpuLayers))
		llamaCppImage.NParallel = types.Int32Value(int32(input.Model.Image.LlamaCpp.NParallel))

		if input.Model.Image.LlamaCpp.Pooling != nil {
			llamaCppImage.Pooling = types.StringValue(string(*input.Model.Image.LlamaCpp.Pooling))
		}
		if input.Model.Image.LlamaCpp.ThreadsHttp != nil {
			llamaCppImage.ThreadsHttp = types.Int32Value(int32(*input.Model.Image.LlamaCpp.ThreadsHttp))
		}
		if input.Model.Image.LlamaCpp.Variant != nil {
			llamaCppImage.Variant = types.StringValue(string(*input.Model.Image.LlamaCpp.Variant))
		}

		modelImage.LlamaCpp, diag = types.ObjectValueFrom(ctx, llamaCppImage.AttributeTypes(), llamaCppImage)
		diags.Append(diag...)
		if diags.HasError() {
			return
		}
	} else {
		modelImage.LlamaCpp, diag = types.ObjectValueFrom(ctx, models.ModelImageLlamacpp{}.AttributeTypes(), models.ModelImageLlamacpp{
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
		diags.Append(diag...)
		if diags.HasError() {
			return
		}
	}

	if input.Model.Image.Custom != nil {
		customImage := models.ModelImageCustom{}

		if input.Model.Image.Custom.HealthRoute != nil {
			customImage.HealthRoute = types.StringValue(*input.Model.Image.Custom.HealthRoute)
		}

		customImage.Port = types.Int32Value(int32(input.Model.Image.Custom.Port))

		customImage.URL = types.StringValue(input.Model.Image.Custom.URL)

		if input.Model.Image.Custom.Credentials != nil {
			credentials := models.Credentials{
				Username: types.StringValue(string(input.Model.Image.Custom.Credentials.Username)),
				Password: types.StringValue(string(*input.Model.Image.Custom.Credentials.Password)),
			}

			customImage.Credentials, diag = types.ObjectValueFrom(ctx, credentials.AttributeTypes(), credentials)
			diags.Append(diag...)
			if diags.HasError() {
				return
			}
		}

		modelImage.Custom, diag = types.ObjectValueFrom(ctx, customImage.AttributeTypes(), customImage)
		diags.Append(diag...)
		if diags.HasError() {
			return
		}
	} else {
		mCredentials := models.Credentials{
			Username: types.StringNull(),
			Password: types.StringNull(),
		}

		var credentials types.Object
		credentials, diag = types.ObjectValueFrom(ctx, mCredentials.AttributeTypes(), mCredentials)
		diags.Append(diag...)
		if diags.HasError() {
			return
		}

		modelImage.Custom, diag = types.ObjectValueFrom(ctx, models.ModelImageCustom{}.AttributeTypes(), models.ModelImageCustom{
			HealthRoute: types.StringNull(),
			Port:        types.Int32Null(),
			URL:         types.StringNull(),
			Credentials: credentials,
		})
		diags.Append(diag...)
		if diags.HasError() {
			return
		}
	}

	endpointModel.Image, diag = types.ObjectValueFrom(ctx, modelImage.AttributeTypes(), modelImage)
	diags.Append(diag...)
	if diags.HasError() {
		return
	}

	output.Model, diag = types.ObjectValueFrom(ctx, endpointModel.AttributeTypes(), endpointModel)
	diags.Append(diag...)
	if diags.HasError() {
		return
	}

	// Tags
	output.Tags, diag = types.ListValueFrom(ctx, types.StringType, input.Tags)
	diags.Append(diag...)
	if diags.HasError() {
		return
	}

	if input.CacheHttpResponses != nil {
		output.CacheHttpResponses = types.BoolValue(*input.CacheHttpResponses)
	} else {
		output.CacheHttpResponses = types.BoolValue(false)
	}

	// Experimental Features
	var experimentalFeatures models.ExperimentalFeatures
	if input.ExperimentalFeatures != nil {
		experimentalFeatures = models.ExperimentalFeatures{
			CacheHTTPResponses: types.BoolValue(input.ExperimentalFeatures.CacheHttpResponses),
		}

		var kvRouter models.KvRouter
		if input.ExperimentalFeatures.KvRouter != nil {
			kvRouter = models.KvRouter{
				Tag: types.StringValue(input.ExperimentalFeatures.KvRouter.Tag),
			}
		} else {
			kvRouter = models.KvRouter{
				Tag: types.StringValue(""),
			}
		}

		experimentalFeatures.KVRouter, diag = types.ObjectValueFrom(ctx, kvRouter.AttributeTypes(), kvRouter)
		diags.Append(diag...)
		if diags.HasError() {
			return
		}
	} else {
		experimentalFeatures = models.ExperimentalFeatures{
			CacheHTTPResponses: types.BoolValue(false),
		}

		kvRouter := models.KvRouter{
			Tag: types.StringValue(""),
		}

		experimentalFeatures.KVRouter, diag = types.ObjectValueFrom(ctx, kvRouter.AttributeTypes(), kvRouter)
		diags.Append(diag...)
		if diags.HasError() {
			return
		}
	}

	output.ExperimentalFeatures, diag = types.ObjectValueFrom(ctx, experimentalFeatures.AttributeTypes(), experimentalFeatures)
	diags.Append(diag...)
	if diags.HasError() {
		return
	}

	// Private Service
	var endpointPrivateService models.PrivateService
	if input.PrivateService != nil {
		endpointPrivateService = models.PrivateService{
			AccountID: types.StringValue(input.PrivateService.AccountID),
			Shared:    types.BoolValue(input.PrivateService.Shared),
		}
	} else {
		endpointPrivateService = models.PrivateService{
			AccountID: types.StringValue(""),
			Shared:    types.BoolValue(false),
		}
	}

	output.PrivateService, diag = types.ObjectValueFrom(ctx, endpointPrivateService.AttributeTypes(), endpointPrivateService)
	diags.Append(diag...)
	if diags.HasError() {
		return
	}

	// Route
	var endpointRoute models.Route
	if input.Route != nil {
		endpointRoute = models.Route{
			Domain: types.StringValue(input.Route.Domain),
			Path:   types.StringValue(input.Route.Path),
		}
	} else {
		endpointRoute = models.Route{
			Domain: types.StringValue(""),
			Path:   types.StringValue(""),
		}
	}

	output.Route, diag = types.ObjectValueFrom(ctx, endpointRoute.AttributeTypes(), endpointRoute)
	diags.Append(diag...)
	if diags.HasError() {
		return
	}

	// Status
	endpointStatus := models.Status{
		CreatedAt:     types.StringValue(input.Status.CreatedAt.String()),
		UpdatedAt:     types.StringValue(input.Status.UpdatedAt.String()),
		State:         types.StringValue(string(input.Status.State)),
		Message:       types.StringValue(input.Status.Message),
		ReadyReplica:  types.Int32Value(int32(input.Status.ReadyReplica)),
		TargetReplica: types.Int32Value(int32(input.Status.TargetReplica)),
	}

	endpointStatusCreatedBy := models.User{
		Id:   types.StringValue(input.Status.CreatedBy.ID),
		Name: types.StringValue(input.Status.CreatedBy.Name),
	}

	endpointStatus.CreatedBy, diag = types.ObjectValueFrom(ctx, endpointStatusCreatedBy.AttributeTypes(), endpointStatusCreatedBy)
	diags.Append(diag...)
	if diags.HasError() {
		return
	}

	endpointStatusUpdatedBy := models.User{
		Id:   types.StringValue(input.Status.UpdatedBy.ID),
		Name: types.StringValue(input.Status.UpdatedBy.Name),
	}

	endpointStatus.UpdatedBy, diag = types.ObjectValueFrom(ctx, endpointStatusUpdatedBy.AttributeTypes(), endpointStatusUpdatedBy)
	diags.Append(diag...)
	if diags.HasError() {
		return
	}

	if input.Status.ErrorMessage != nil {
		endpointStatus.ErrorMessage = types.StringValue(*input.Status.ErrorMessage)
	} else {
		endpointStatus.ErrorMessage = types.StringValue("")
	}

	if input.Status.URL != nil {
		endpointStatus.Url = types.StringValue(*input.Status.URL)
	} else {
		endpointStatus.Url = types.StringValue("")
	}

	var endpointStatusPrivate models.Private
	if input.Status.Private != nil && input.Status.Private.ServiceName != nil {
		endpointStatusPrivate = models.Private{
			ServiceName: types.StringValue(*input.Status.Private.ServiceName),
		}
	} else {
		endpointStatusPrivate = models.Private{
			ServiceName: types.StringValue(""),
		}
	}

	endpointStatus.Private, diag = types.ObjectValueFrom(ctx, endpointStatusPrivate.AttributeTypes(), endpointStatusPrivate)
	diags.Append(diag...)
	if diags.HasError() {
		return
	}

	output.Status, diag = types.ObjectValueFrom(ctx, endpointStatus.AttributeTypes(), endpointStatus)
	diags.Append(diag...)
	if diags.HasError() {
		return
	}

	return
}
