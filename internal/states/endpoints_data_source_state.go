package states

import (
	"github.com/sebps/terraform-provider-huggingface/internal/models"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// endpointsDataSourceModel maps the data source schema data.
type EndpointsDataSourceState struct {
	Namespace types.String      `tfsdk:"namespace"`
	Endpoints []models.Endpoint `tfsdk:"endpoints"`
}
