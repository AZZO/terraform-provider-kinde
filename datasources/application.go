package datasources

import (
	"context"
	"fmt"

	"github.com/AZZO/terraform-provider-kinde/kinde_client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &ApplicationDataSource{}
	_ datasource.DataSourceWithConfigure = &ApplicationDataSource{}
)

// ApplicationDataSource defines the data source implementation.
type ApplicationDataSource struct {
	client *kinde_client.Client
}

// ApplicationDataSourceModel describes the data source data model.
type ApplicationDataSourceModel struct {
	ApplicationId types.String `tfsdk:"application_id"`
	Name          types.String `tfsdk:"name"`
	Type          types.String `tfsdk:"type"`
	ClientId      types.String `tfsdk:"client_id"`
	ClientSecret  types.String `tfsdk:"client_secret"`
	LogoutUris    types.List   `tfsdk:"logout_uris"`
	RedirectUris  types.List   `tfsdk:"redirect_uris"`
}

// NewApplicationDataSource is a helper function to simplify the provider implementation.
func NewApplicationDataSource() datasource.DataSource {
	return &ApplicationDataSource{}
}

// Metadata returns the data source type name.
func (d *ApplicationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

// Schema defines the schema for the data source.
func (d *ApplicationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get information about a Kinde application.",
		Attributes: map[string]schema.Attribute{
			"application_id": schema.StringAttribute{
				Description: "Unique identifier of the application.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Human-readable name of the application.",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "Type of the application.",
				Computed:    true,
			},
			"client_id": schema.StringAttribute{
				Description: "Client ID of the application.",
				Optional:    true,
				Computed:    true,
			},
			"client_secret": schema.StringAttribute{
				Description: "Client Secret of the application.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
			},
			"logout_uris": schema.ListAttribute{
				Description: "List of logout URIs for the application.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"redirect_uris": schema.ListAttribute{
				Description: "List of redirect URIs for the application.",
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *ApplicationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *ApplicationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ApplicationDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	application, err := d.client.GetApplication(ctx, data.ApplicationId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Application",
			"Could not read application ID "+data.ApplicationId.ValueString()+": "+err.Error(),
		)
		return
	}

	// Get callbacks
	callbacks, err := d.client.GetApplicationCallbacks(ctx, data.ApplicationId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading application callbacks",
			"Could not read application callbacks: "+err.Error(),
		)
		return
	}

	// Convert []string to types.List
	logoutUrisList := make([]attr.Value, len(callbacks.LogoutUris))
	for i, uri := range callbacks.LogoutUris {
		logoutUrisList[i] = types.StringValue(uri)
	}

	redirectUrisList := make([]attr.Value, len(callbacks.RedirectUris))
	for i, uri := range callbacks.RedirectUris {
		redirectUrisList[i] = types.StringValue(uri)
	}

	// Map response body to schema and populate Computed attribute values
	data.ApplicationId = types.StringValue(application.Id)
	data.Name = types.StringValue(application.Name)
	data.Type = types.StringValue(application.Type)
	data.ClientId = types.StringValue(application.ClientId)
	data.ClientSecret = types.StringValue(application.ClientSecret)
	data.LogoutUris = types.ListValueMust(types.StringType, logoutUrisList)
	data.RedirectUris = types.ListValueMust(types.StringType, redirectUrisList)

	// Save data into state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
