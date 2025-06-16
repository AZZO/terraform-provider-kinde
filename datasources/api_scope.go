package datasources

import (
	"context"

	"github.com/AZZO/terraform-provider-kinde/kinde_client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &ApiScopeDataSource{}
	_ datasource.DataSourceWithConfigure = &ApiScopeDataSource{}
)

// ApiScopeDataSource is the data source implementation.
type ApiScopeDataSource struct {
	client *kinde_client.Client
}

// ApiScopeDataSourceModel describes the data source data model.
type ApiScopeDataSourceModel struct {
	ScopeId     types.String `tfsdk:"scope_id"`
	ApiId       types.String `tfsdk:"api_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

// NewApiScopeDataSource is a helper function to simplify the provider implementation.
func NewApiScopeDataSource() datasource.DataSource {
	return &ApiScopeDataSource{}
}

// Metadata returns the data source type name.
func (d *ApiScopeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_scope"
}

// Schema defines the schema for the data source.
func (d *ApiScopeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get information about a Kinde API scope.",
		Attributes: map[string]schema.Attribute{
			"scope_id": schema.StringAttribute{
				Description: "Unique identifier of the scope.",
				Required:    true,
			},
			"api_id": schema.StringAttribute{
				Description: "The ID of the API this scope belongs to.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Human-readable name of the scope.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the scope.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *ApiScopeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*kinde_client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected *kinde_client.Client, got: %T. Please report this issue to the provider developers.",
		)
		return
	}

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *ApiScopeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ApiScopeDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get all scopes for the API
	scopes, err := d.client.GetApiScopes(ctx, data.ApiId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading API Scopes",
			"Could not read API scopes: "+err.Error(),
		)
		return
	}

	// Find the specific scope
	var foundScope *kinde_client.ApiScopeResource
	for _, scope := range scopes {
		if scope.Id == data.ScopeId.ValueString() {
			foundScope = &scope
			break
		}
	}

	if foundScope == nil {
		resp.Diagnostics.AddError(
			"Error Reading API Scope",
			"Could not find API scope with ID "+data.ScopeId.ValueString(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	data.ScopeId = types.StringValue(foundScope.Id)
	data.Name = types.StringValue(foundScope.Name)
	data.Description = types.StringValue(foundScope.Description)

	// Save data into state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
