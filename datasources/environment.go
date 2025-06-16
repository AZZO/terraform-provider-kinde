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
	_ datasource.DataSource              = &EnvironmentDataSource{}
	_ datasource.DataSourceWithConfigure = &EnvironmentDataSource{}
)

// EnvironmentDataSource defines the data source implementation.
type EnvironmentDataSource struct {
	client *kinde_client.Client
}

// EnvironmentDataSourceModel describes the data source data model.
type EnvironmentDataSourceModel struct {
	EnvironmentId   types.String `tfsdk:"environment_id"`
	Code            types.String `tfsdk:"code"`
	Name            types.String `tfsdk:"name"`
	IsDefault       types.Bool   `tfsdk:"is_default"`
	IsLive          types.Bool   `tfsdk:"is_live"`
	KindeDomain     types.String `tfsdk:"kinde_domain"`
	CustomDomain    types.String `tfsdk:"custom_domain"`
	Logo            types.String `tfsdk:"logo"`
	LogoDark        types.String `tfsdk:"logo_dark"`
	FaviconSvg      types.String `tfsdk:"favicon_svg"`
	FaviconFallback types.String `tfsdk:"favicon_fallback"`
	CreatedOn       types.String `tfsdk:"created_on"`
}

// NewEnvironmentDataSource is a helper function to simplify the provider implementation.
func NewEnvironmentDataSource() datasource.DataSource {
	return &EnvironmentDataSource{}
}

// Metadata returns the data source type name.
func (d *EnvironmentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

// Schema defines the schema for the data source.
func (d *EnvironmentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get information about a Kinde environment.",
		Attributes: map[string]schema.Attribute{
			"environment_id": schema.StringAttribute{
				Description: "Unique identifier of the environment.",
				Required:    true,
			},
			"code": schema.StringAttribute{
				Description: "Code identifier of the environment.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Human-readable name of the environment.",
				Computed:    true,
			},
			"is_default": schema.BoolAttribute{
				Description: "Whether this is the default environment.",
				Computed:    true,
			},
			"is_live": schema.BoolAttribute{
				Description: "Whether this is a live environment.",
				Computed:    true,
			},
			"kinde_domain": schema.StringAttribute{
				Description: "The Kinde domain for this environment.",
				Computed:    true,
			},
			"custom_domain": schema.StringAttribute{
				Description: "The custom domain for this environment, if configured.",
				Computed:    true,
			},
			"logo": schema.StringAttribute{
				Description: "URL of the environment's logo.",
				Computed:    true,
			},
			"logo_dark": schema.StringAttribute{
				Description: "URL of the environment's dark mode logo.",
				Computed:    true,
			},
			"favicon_svg": schema.StringAttribute{
				Description: "URL of the environment's SVG favicon.",
				Computed:    true,
			},
			"favicon_fallback": schema.StringAttribute{
				Description: "URL of the environment's fallback favicon.",
				Computed:    true,
			},
			"created_on": schema.StringAttribute{
				Description: "Timestamp when the environment was created.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *EnvironmentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *EnvironmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EnvironmentDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	environment, err := d.client.GetEnvironment(ctx, data.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Environment",
			"Could not read environment ID "+data.EnvironmentId.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	data.EnvironmentId = types.StringValue(data.EnvironmentId.ValueString())
	data.Code = types.StringValue(environment.Code)
	data.Name = types.StringValue(environment.Name)
	data.IsDefault = types.BoolValue(environment.IsDefault)
	data.IsLive = types.BoolValue(environment.IsLive)
	data.KindeDomain = types.StringValue(environment.KindeDomain)
	data.CustomDomain = types.StringValue(environment.CustomDomain)
	data.Logo = types.StringValue(environment.Logo)
	data.LogoDark = types.StringValue(environment.LogoDark)
	data.FaviconSvg = types.StringValue(environment.FaviconSvg)
	data.FaviconFallback = types.StringValue(environment.FaviconFallback)
	data.CreatedOn = types.StringValue(environment.CreatedOn)

	// Save data into state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
