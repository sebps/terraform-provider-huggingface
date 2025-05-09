package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type Endpoint struct {
	Namespace            types.String        `tfsdk:"namespace"`
	Name                 types.String        `tfsdk:"name"`
	Type                 types.String        `tfsdk:"type"`
	CloudProvider        types.Object        `tfsdk:"cloud_provider"`
	Compute              types.Object        `tfsdk:"compute"`
	Model                types.Object        `tfsdk:"model"`
	Tags                 basetypes.ListValue `tfsdk:"tags"`
	CacheHttpResponses   types.Bool          `tfsdk:"cache_http_responses"`
	ExperimentalFeatures types.Object        `tfsdk:"experimental_features"`
	PrivateService       types.Object        `tfsdk:"private_service"`
	Route                types.Object        `tfsdk:"route"`
	Status               types.Object        `tfsdk:"status"`
}

type EndpointCloudProvider struct {
	Vendor types.String `tfsdk:"vendor"`
	Region types.String `tfsdk:"region"`
}

func (m EndpointCloudProvider) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"vendor": types.StringType,
		"region": types.StringType,
	}
}

type EndpointCompute struct {
	ID           types.String `tfsdk:"id"`
	Accelerator  types.String `tfsdk:"accelerator"`
	InstanceType types.String `tfsdk:"instance_type"`
	InstanceSize types.String `tfsdk:"instance_size"`
	Scaling      types.Object `tfsdk:"scaling"`
}

func (e EndpointCompute) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":            types.StringType,
		"accelerator":   types.StringType,
		"instance_type": types.StringType,
		"instance_size": types.StringType,
		"scaling": types.ObjectType{
			AttrTypes: EndpointComputeScaling{}.AttributeTypes(),
		},
	}
}

type EndpointComputeScaling struct {
	MinReplica         types.Int32  `tfsdk:"min_replica"`
	MaxReplica         types.Int32  `tfsdk:"max_replica"`
	Measure            types.Object `tfsdk:"measure"`
	Metric             types.Object `tfsdk:"metric"`
	ScaleToZeroTimeout types.Number `tfsdk:"scale_to_zero_timeout"`
	Threshold          types.Number `tfsdk:"threshold"`
}

func (e EndpointComputeScaling) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"min_replica": types.NumberType,
		"max_replica": types.NumberType,
		"measure": types.ObjectType{
			AttrTypes: EndpointComputeScalingMeasure{}.AttributeTypes(),
		},
		"metric":                types.StringType,
		"scale_to_zero_timeout": types.NumberType,
		"threshold":             types.NumberType,
	}
}

type EndpointComputeScalingMeasure struct {
	HardwareUsage   types.Number `json:"hardware_usage"`
	PendingRequests types.Number `json:"pending_requests"`
}

func (e EndpointComputeScalingMeasure) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"hardware_usage":   types.NumberType,
		"pending_requests": types.NumberType,
	}
}

type Model struct {
	Repository types.String `tfsdk:"repository"`
	Framework  types.String `tfsdk:"framework"`
	Task       types.String `tfsdk:"task"`
	Image      types.Object `tfsdk:"image"`
}

func (m Model) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"repository": types.StringType,
		"framework":  types.StringType,
		"task":       types.StringType,
	}
}

type ModelImage struct {
	HuggingFace       types.Object `tfsdk:"huggingface"`
	HuggingFaceNeuron types.Object `tfsdk:"huggingface_neuron"`
	TGI               types.Object `tfsdk:"tgi"`
	TGINeuron         types.Object `tfsdk:"tgi_neuron"`
	TEI               types.Object `tfsdk:"tei"`
	LlamaCpp          types.Object `tfsdk:"llamacpp"`
	Custom            types.Object `tfsdk:"custom"`
}

func (m ModelImage) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"huggingface": types.ObjectType{},
		"huggingface_neuron": types.ObjectType{
			AttrTypes: ModelImageHuggingfaceNeuron{}.AttributeTypes(),
		},
		"tgi": types.ObjectType{
			AttrTypes: ModelImageTgi{}.AttributeTypes(),
		},
		"tgi_neuron": types.ObjectType{
			AttrTypes: ModelImageTgiNeuron{}.AttributeTypes(),
		},
		"tei": types.ObjectType{
			AttrTypes: ModelImageTei{}.AttributeTypes(),
		},
		"llamacpp": types.ObjectType{
			AttrTypes: ModelImageLlamacpp{}.AttributeTypes(),
		},
		"custom": types.ObjectType{
			AttrTypes: ModelImageCustom{}.AttributeTypes(),
		},
	}
}

type ModelImageHuggingface struct{}

type ModelImageHuggingfaceNeuron struct {
	BatchSize      types.Int32  `tfsdk:"batch_size"`
	NeuronCache    types.String `tfsdk:"neuron_cache"`
	SequenceLength types.Int32  `tfsdk:"sequence_length"`
}

func (m ModelImageHuggingfaceNeuron) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"batch_size":      types.Int32Type,
		"neuron_cache":    types.StringType,
		"sequence_length": types.Int32Type,
	}
}

type ModelImageTgi struct {
	HealthRoute           types.String `tfsdk:"health_route"`
	Port                  types.Int32  `tfsdk:"port"`
	Url                   types.String `tfsdk:"url"`
	MaxBatchPrefillTokens types.Int32  `tfsdk:"max_batch_prefill_tokens"`
	MaxBatchTotalTokens   types.Int32  `tfsdk:"max_batch_total_tokens"`
	MaxInputLength        types.Int32  `tfsdk:"max_input_length"`
	MaxTotalTokens        types.Int32  `tfsdk:"max_total_tokens"`
	DisableCustomKernels  types.Bool   `tfsdk:"disable_custom_kernels"`
	Quantize              types.String `tfsdk:"quantize"`
}

func (m ModelImageTgi) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"health_route":             types.StringType,
		"port":                     types.Int32Type,
		"url":                      types.StringType,
		"max_batch_prefill_tokens": types.Int32Type,
		"max_batch_total_tokens":   types.Int32Type,
		"max_input_length":         types.Int32Type,
		"max_total_tokens":         types.Int32Type,
		"disable_custom_kernels":   types.BoolType,
		"quantize":                 types.StringType,
	}
}

type ModelImageTgiNeuron struct {
	HealthRoute           types.String `tfsdk:"health_route"`
	Port                  types.Int32  `tfsdk:"port"`
	Url                   types.String `tfsdk:"url"`
	MaxBatchPrefillTokens types.Int32  `tfsdk:"max_batch_prefill_tokens"`
	MaxBatchTotalTokens   types.Int32  `tfsdk:"max_batch_total_tokens"`
	MaxInputLength        types.Int32  `tfsdk:"max_input_length"`
	MaxTotalTokens        types.Int32  `tfsdk:"max_total_tokens"`
	HfAutoCastType        types.String `tfsdk:"hf_auto_cast_type"`
	HfNumCores            types.Int32  `tfsdk:"hf_num_cores"`
}

func (m ModelImageTgiNeuron) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"health_route":             types.StringType,
		"port":                     types.Int32Type,
		"url":                      types.StringType,
		"max_batch_prefill_tokens": types.Int32Type,
		"max_batch_total_tokens":   types.Int32Type,
		"max_input_length":         types.Int32Type,
		"max_total_tokens":         types.Int32Type,
		"hf_auto_cast_type":        types.StringType,
		"hf_num_cores":             types.Int32Type,
	}
}

type ModelImageTei struct {
	HealthRoute           types.String `tfsdk:"health_route"`
	Port                  types.Int32  `tfsdk:"port"`
	URL                   types.String `tfsdk:"url"`
	MaxBatchTokens        types.Int32  `tfsdk:"max_batch_tokens"`
	MaxConcurrentRequests types.Int32  `tfsdk:"max_concurrent_requests"`
	Pooling               types.String `tfsdk:"pooling"`
}

func (m ModelImageTei) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"healthRoute":           types.Int32Type,
		"port":                  types.Int32Type,
		"url":                   types.Int32Type,
		"maxBatchTokens":        types.Int32Type,
		"maxConcurrentRequests": types.Int32Type,
		"pooling":               types.Int32Type,
	}
}

type ModelImageLlamacpp struct {
	HealthRoute types.String `json:"health_route"`
	Port        types.Int32  `json:"port"`
	URL         types.String `json:"url"`
	CtxSize     types.Int32  `json:"ctx_size"`
	Mode        types.String `json:"mode"`
	ModelPath   types.String `json:"model_path"`
	NGpuLayers  types.Int32  `json:"n_gpu_layers"`
	NParallel   types.Int32  `json:"n_parallel"`
	Pooling     types.String `json:"pooling"`
	ThreadsHttp types.Int32  `json:"threads_http"`
	Variant     types.String `json:"variant"`
}

func (m ModelImageLlamacpp) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"health_route": types.StringType,
		"port":         types.Int32Type,
		"url":          types.StringType,
		"ctx_size":     types.Int32Type,
		"mode":         types.StringType,
		"model_path":   types.StringType,
		"n_gpu_layers": types.Int32Type,
		"n_parallel":   types.Int32Type,
		"pooling":      types.StringType,
		"threads_http": types.Int32Type,
		"variant":      types.StringType,
	}
}

type ModelImageCustom struct {
	HealthRoute types.String `json:"health_route"`
	Port        types.Int32  `json:"port"`
	URL         types.String `json:"url"`
	Credentials types.Object `json:"credentials"`
}

func (m ModelImageCustom) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"health_route": types.StringType,
		"port":         types.Int32Type,
		"url":          types.StringType,
		"credentials": types.ObjectType{
			AttrTypes: Credentials{}.AttributeTypes(),
		},
	}
}

type Credentials struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (m Credentials) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"username": types.StringType,
		"password": types.StringType,
	}
}

type ExperimentalFeatures struct {
	CacheHTTPResponses types.Bool   `tfsdk:"cache_http_responses"`
	KVRouter           types.Object `tfsdk:"kv_router"`
}

func (e ExperimentalFeatures) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cache_http_responses": types.BoolType,
		"kv_router": types.ObjectType{
			AttrTypes: KvRouter{}.AttributeTypes(),
		},
	}
}

type KvRouter struct {
	Tag types.String `tfsdk:"tag"`
}

func (k KvRouter) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"tag": types.StringType,
	}
}

type PrivateService struct {
	AccountID types.String `tfsdk:"account_id"`
	Shared    types.Bool   `tfsdk:"shared"`
}

func (p PrivateService) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"account_id": types.StringType,
		"shared":     types.BoolType,
	}
}

type Route struct {
	Domain types.String `tfsdk:"domain"`
	Path   types.String `tfsdk:"path"`
}

func (r Route) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"domain": types.StringType,
		"path":   types.StringType,
	}
}

type Status struct {
	CreatedAt     types.String `tfsdk:"created_at"`
	CreatedBy     types.Object `tfsdk:"created_by"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
	UpdatedBy     types.Object `tfsdk:"updated_by"`
	State         types.String `tfsdk:"state"`
	Message       types.String `tfsdk:"message"`
	ReadyReplica  types.Int32  `tfsdk:"ready_replica"`
	TargetReplica types.Int32  `tfsdk:"target_replica"`
	ErrorMessage  types.String `tfsdk:"error_message"`
	Url           types.String `tfsdk:"url"`
	Private       types.Object `tfsdk:"private"`
}

func (s Status) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"created_at": types.StringType,
		"created_by": types.ObjectType{
			AttrTypes: User{}.AttributeTypes(),
		},
		"updated_at": types.StringType,
		"updated_by": types.ObjectType{
			AttrTypes: User{}.AttributeTypes(),
		},
		"state":          types.StringType,
		"message":        types.StringType,
		"ready_replica":  types.NumberType,
		"target_replica": types.NumberType,
		"error_message":  types.StringType,
		"url":            types.StringType,
		"private": types.ObjectType{
			AttrTypes: Private{}.AttributeTypes(),
		},
	}
}

type User struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (u User) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	}
}

type Private struct {
	ServiceName types.String `tfsdk:"service_name"`
}

func (p Private) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"service_name": types.StringType,
	}
}
