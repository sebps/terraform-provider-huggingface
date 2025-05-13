// Copyright (c) sebps.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	huggingface "github.com/sebps/huggingface-client/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &huggingfaceProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &huggingfaceProvider{
			version: version,
		}
	}
}

// huggingfaceProvider is the provider implementation.
type huggingfaceProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *huggingfaceProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "huggingface"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *huggingfaceProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"hf_token": schema.StringAttribute{
				Description: "Hugging face token from the Access Token section",
				Optional:    false,
				Sensitive:   true,
				Required:    true,
			},
		},
	}
}

// Configure prepares a huggingface API client for data sources and resources.
func (p *huggingfaceProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Huggingface provider")

	// Retrieve provider data from configuration
	var config hashicupsProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.HfToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("hf_token"),
			"Unknown Hugging Face Token",
			"The provider cannot create the HuggingFace API client as there is an unknown configuration value for the Hugging Face Token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the HASHICUPS_HOST environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	hfToken := os.Getenv("HF_TOKEN")

	if !config.HfToken.IsNull() {
		hfToken = config.HfToken.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if hfToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Hugging Face Token",
			"The provider cannot create the HuggingFace API client as there is a missing or empty value for the Hugging Face Token. "+
				"Set the host value in the configuration or use the HF_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// logging
	ctx = tflog.SetField(ctx, "huggingface_token", hfToken)
	tflog.Debug(ctx, "Creating Huggingface client")

	// Create a new Hugging Face client using the configuration values
	// h := "http://localhost:8080"
	client, err := huggingface.NewClient(nil, &hfToken)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Hugging Face API Client",
			"An unexpected error occurred when creating the Hugging Face API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Hugging Face Client Error: "+err.Error(),
		)
		return
	}

	// Make the HuggingFace client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Huggingface client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *huggingfaceProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewEndpointsDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *huggingfaceProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewEndpointsResource,
	}
}

// hashicupsProviderModel maps provider schema data to a Go type.
type hashicupsProviderModel struct {
	HfToken types.String `tfsdk:"hf_token"`
}
