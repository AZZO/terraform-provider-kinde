package kinde

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
	_ datasource.DataSource              = &applicationDataSource{}
	_ datasource.DataSourceWithConfigure = &applicationDataSource{}
)

func NewApplicationDataSource() datasource.DataSource {
	return &applicationDataSource{}
}

// applicationDataSource defines the data source implementation.
type applicationDataSource struct {
	client *kinde_client.Client
}

// applicationDataSourceModel describes the data source data model.
type applicationDataSourceModel struct {
	ApplicationId types.String `tfsdk:"application_id"`
	Name          types.String `tfsdk:"name"`
	Type          types.String `tfsdk:"type"`
	ClientId      types.String `tfsdk:"client_id"`
	ClientSecret  types.String `tfsdk:"client_secret"`
}

func (d *applicationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (d *applicationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"application_id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Human-readable name of the application",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of the application",
				Computed:            true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "Client ID",
				Optional:            true,
				Computed:            true,
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "Client Secret",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
		},
	}
}

func (d *applicationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*kinde_client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *applicationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data applicationDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	detail, err := d.client.GetApplication(data.ApplicationId.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Could not fetch Application", err.Error())
		return
	}

	data.ApplicationId = types.StringValue(detail.Id)
	data.Name = types.StringValue(detail.Name)
	data.Type = types.StringValue(detail.Type)
	data.ClientId = types.StringValue(detail.ClientId)
	data.ClientSecret = types.StringValue(detail.ClientSecret)

	// Save data into state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
