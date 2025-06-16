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
	_ datasource.DataSource              = &RoleDataSource{}
	_ datasource.DataSourceWithConfigure = &RoleDataSource{}
)

// RoleDataSource defines the data source implementation.
type RoleDataSource struct {
	client *kinde_client.Client
}

// RoleDataSourceModel describes the data source data model.
type RoleDataSourceModel struct {
	RoleId        types.String `tfsdk:"role_id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Key           types.String `tfsdk:"key"`
	IsDefaultRole types.Bool   `tfsdk:"is_default_role"`
	Permissions   types.List   `tfsdk:"permissions"`
	Scopes        types.List   `tfsdk:"scopes"`
}

// NewRoleDataSource is a helper function to simplify the provider implementation.
func NewRoleDataSource() datasource.DataSource {
	return &RoleDataSource{}
}

// Metadata returns the data source type name.
func (d *RoleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

// Schema defines the schema for the data source.
func (d *RoleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get information about a Kinde role.",
		Attributes: map[string]schema.Attribute{
			"role_id": schema.StringAttribute{
				Description: "Unique identifier of the role.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Human-readable name of the role.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the role.",
				Computed:    true,
			},
			"key": schema.StringAttribute{
				Description: "Key of the role.",
				Computed:    true,
			},
			"is_default_role": schema.BoolAttribute{
				Description: "Whether this is a default role.",
				Computed:    true,
			},
			"permissions": schema.ListAttribute{
				Description: "List of permission IDs associated with this role.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"scopes": schema.ListAttribute{
				Description: "List of scope IDs associated with this role.",
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *RoleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *RoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RoleDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	role, err := d.client.GetRole(ctx, data.RoleId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Role",
			"Could not read role ID "+data.RoleId.ValueString()+": "+err.Error(),
		)
		return
	}

	// Get permissions
	permissions, err := d.client.GetRolePermissions(ctx, data.RoleId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Role Permissions",
			"Could not read role permissions: "+err.Error(),
		)
		return
	}

	// Get scopes
	roleScopes, err := d.client.GetRoleScopes(ctx, data.RoleId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading role scopes",
			"Could not read role scopes: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	data.RoleId = types.StringValue(role.Id)
	data.Name = types.StringValue(role.Name)
	data.Description = types.StringValue(role.Description)
	data.Key = types.StringValue(role.Key)
	data.IsDefaultRole = types.BoolValue(role.IsDefaultRole)

	// Convert permissions to list
	if len(permissions) > 0 {
		permissionValues := make([]attr.Value, len(permissions))
		for i, p := range permissions {
			permissionValues[i] = types.StringValue(p)
		}
		data.Permissions = types.ListValueMust(types.StringType, permissionValues)
	} else {
		data.Permissions = types.ListValueMust(types.StringType, []attr.Value{})
	}

	// Convert scopes to list
	if len(roleScopes) > 0 {
		scopeValues := make([]attr.Value, len(roleScopes))
		for i, s := range roleScopes {
			scopeValues[i] = types.StringValue(s.Id)
		}
		data.Scopes = types.ListValueMust(types.StringType, scopeValues)
	} else {
		data.Scopes = types.ListValueMust(types.StringType, []attr.Value{})
	}

	// Save data into state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
