package datasources

import (
	"context"
	"fmt"

	"github.com/AZZO/terraform-provider-kinde/kinde_client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &ApiDataSource{}
	_ datasource.DataSourceWithConfigure = &ApiDataSource{}
)

// ApiDataSource defines the data source implementation.
type ApiDataSource struct {
	client *kinde_client.Client
}

// ApiDataSourceModel describes the data source data model.
type ApiDataSourceModel struct {
	ApiId    types.String `tfsdk:"api_id"`
	Name     types.String `tfsdk:"name"`
	Audience types.String `tfsdk:"audience"`
}

// NewApiDataSource is a helper function to simplify the provider implementation.
func NewApiDataSource() datasource.DataSource {
	return &ApiDataSource{}
}

// Metadata returns the data source type name.
func (d *ApiDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api"
}

// Schema defines the schema for the data source.
func (d *ApiDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get information about a Kinde API.",
		Attributes: map[string]schema.Attribute{
			"api_id": schema.StringAttribute{
				Description: "Unique identifier of the API.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Human-readable name of the API.",
				Computed:    true,
			},
			"audience": schema.StringAttribute{
				Description: "The audience for the API.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *ApiDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*kinde_client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *kinde_client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *ApiDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ApiDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	api, err := d.client.GetApi(ctx, data.ApiId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading API",
			"Could not read API ID "+data.ApiId.ValueString()+": "+err.Error(),
		)
		return
	}

	// Set the data source attributes
	data.ApiId = types.StringValue(api.Id)
	data.Name = types.StringValue(api.Name)
	data.Audience = types.StringValue(api.Audience)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
