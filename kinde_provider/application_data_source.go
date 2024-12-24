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
var _ datasource.DataSource = &ApplicationDataSource{}

func NewApplicationDataSource() datasource.DataSource {
	return &ApplicationDataSource{}
}

// ApplicationDataSource defines the data source implementation.
type ApplicationDataSource struct {
	client *kinde_client.Client
}

// ApplicationDataSourceModel describes the data source data model.
type ApplicationDataSourceModel struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Type         types.String `tfsdk:"type"`
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	RedirectUris types.List   `tfsdk:"redirect_uris"`
	LogoutUris   types.List   `tfsdk:"logout_uris"`
}

func (d *ApplicationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_example"
}

func (d *ApplicationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Human-readable name of the application",
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of the application",
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "Client ID",
				Optional:            true,
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "Client Secret",
				Optional:            true,
			},
			"redirect_uris": schema.ListAttribute{
				MarkdownDescription: "List of redirect URIs",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"logout_uris": schema.ListAttribute{
				MarkdownDescription: "List of logout URIs",
				ElementType:         types.StringType,
				Optional:            true,
			},
		},
	}
}

func (d *ApplicationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ApplicationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ApplicationDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	detail, err := d.client.GetApplication(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Could not fetch Application", err.Error())
		return
	}

	data.Id = types.StringValue(detail.Id)
	data.Name = types.StringValue(detail.Name)
	data.Type = types.StringValue(detail.Type)
	data.ClientId = types.StringValue(detail.ClientId)
	data.ClientSecret = types.StringValue(detail.ClientSecret)

	// TODO fetch these from additional endpoints
	data.RedirectUris = types.ListNull(types.StringType)
	data.LogoutUris = types.ListNull(types.StringType)

	// Save data into state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
