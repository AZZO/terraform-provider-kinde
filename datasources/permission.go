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
	_ datasource.DataSource              = &PermissionDataSource{}
	_ datasource.DataSourceWithConfigure = &PermissionDataSource{}
)

// PermissionDataSource defines the data source implementation.
type PermissionDataSource struct {
	client *kinde_client.Client
}

// PermissionDataSourceModel describes the data source data model.
type PermissionDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Key         types.String `tfsdk:"key"`
}

// NewPermissionDataSource is a helper function to simplify the provider implementation.
func NewPermissionDataSource() datasource.DataSource {
	return &PermissionDataSource{}
}

// Metadata returns the data source type name.
func (d *PermissionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_permission"
}

// Schema defines the schema for the data source.
func (d *PermissionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get information about a Kinde permission.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier of the permission.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Human-readable name of the permission.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the permission.",
				Computed:    true,
			},
			"key": schema.StringAttribute{
				Description: "Key of the permission.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *PermissionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *PermissionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PermissionDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	permission, err := d.client.GetPermission(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Permission",
			"Could not read permission ID "+data.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	data.Id = types.StringValue(permission.Id)
	data.Name = types.StringValue(permission.Name)
	data.Description = types.StringValue(permission.Description)
	data.Key = types.StringValue(permission.Key)

	// Save data into state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
