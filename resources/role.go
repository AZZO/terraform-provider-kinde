package resources

import (
	"context"
	"fmt"

	"github.com/AZZO/terraform-provider-kinde/kinde_client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &RoleResource{}
	_ resource.ResourceWithConfigure = &RoleResource{}
)

// RoleResource is the resource implementation.
type RoleResource struct {
	client *kinde_client.Client
}

// RoleResourceModel describes the resource data model.
type RoleResourceModel struct {
	RoleId        types.String `tfsdk:"role_id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Key           types.String `tfsdk:"key"`
	IsDefaultRole types.Bool   `tfsdk:"is_default_role"`
	Permissions   types.List   `tfsdk:"permissions"`
	Scopes        types.List   `tfsdk:"scopes"`
}

// NewRoleResource is a helper function to simplify the provider implementation.
func NewRoleResource() resource.Resource {
	return &RoleResource{}
}

// Configure adds the provider configured client to the resource.
func (r *RoleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*kinde_client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *kinde_client.Client, got: %T.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Metadata returns the resource type name.
func (r *RoleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

// Schema defines the schema for the resource.
func (r *RoleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Kinde Role.",
		Attributes: map[string]schema.Attribute{
			"role_id": schema.StringAttribute{
				Description: "Unique identifier for the role.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description:         "The name of the role.",
				Required:            true,
				MarkdownDescription: "The name of the role.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			"description": schema.StringAttribute{
				Description:         "The description of the role.",
				Optional:            true,
				MarkdownDescription: "The description of the role.",
			},
			"key": schema.StringAttribute{
				Description:         "The key of the role.",
				Required:            true,
				MarkdownDescription: "The key of the role.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"is_default_role": schema.BoolAttribute{
				Description:         "Whether this is a default role.",
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether this is a default role.",
			},
			"permissions": schema.ListAttribute{
				Description:         "List of permission IDs associated with this role.",
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of permission IDs associated with this role.",
			},
			"scopes": schema.ListAttribute{
				Description:         "List of API scope IDs associated with this role.",
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of API scope IDs associated with this role.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *RoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RoleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create role
	role, err := r.client.CreateRole(ctx, plan.Name.ValueString(), plan.Description.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating role",
			"Could not create role, unexpected error: "+err.Error(),
		)
		return
	}

	// Convert types.List to []string
	var permissions []string
	if !plan.Permissions.IsNull() && !plan.Permissions.IsUnknown() {
		diags := plan.Permissions.ElementsAs(ctx, &permissions, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Update role permissions if provided
	if len(permissions) > 0 {
		err = r.client.UpdateRolePermissions(ctx, role.Id, permissions)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating role permissions",
				"Could not update role permissions, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Convert []string back to types.List
	permissionsList := make([]attr.Value, len(permissions))
	for i, permission := range permissions {
		permissionsList[i] = types.StringValue(permission)
	}

	plan = RoleResourceModel{
		RoleId:        types.StringValue(role.Id),
		Name:          types.StringValue(role.Name),
		Description:   types.StringValue(role.Description),
		Key:           types.StringValue(role.Key),
		IsDefaultRole: types.BoolValue(role.IsDefaultRole),
		Permissions:   types.ListValueMust(types.StringType, permissionsList),
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *RoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RoleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	role, err := r.client.GetRole(ctx, state.RoleId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Role",
			"Could not read role: "+err.Error(),
		)
		return
	}

	state.Name = types.StringValue(role.Name)
	state.Description = types.StringValue(role.Description)
	state.Key = types.StringValue(role.Key)
	state.IsDefaultRole = types.BoolValue(role.IsDefaultRole)

	// Get permissions
	permissions, err := r.client.GetRolePermissions(ctx, state.RoleId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Role Permissions",
			"Could not read role permissions: "+err.Error(),
		)
		return
	}
	if len(permissions) > 0 {
		permissionValues := make([]attr.Value, len(permissions))
		for i, p := range permissions {
			permissionValues[i] = types.StringValue(p)
		}
		state.Permissions = types.ListValueMust(types.StringType, permissionValues)
	} else {
		state.Permissions = types.ListValueMust(types.StringType, []attr.Value{})
	}

	// Get scopes
	scopes, err := r.client.GetRoleScopes(ctx, state.RoleId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Role Scopes",
			"Could not read role scopes: "+err.Error(),
		)
		return
	}
	if len(scopes) > 0 {
		scopeValues := make([]attr.Value, len(scopes))
		for i, s := range scopes {
			scopeValues[i] = types.StringValue(s.Id)
		}
		state.Scopes = types.ListValueMust(types.StringType, scopeValues)
	} else {
		state.Scopes = types.ListValueMust(types.StringType, []attr.Value{})
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *RoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan RoleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update role
	role, err := r.client.UpdateRole(ctx, plan.RoleId.ValueString(), plan.Name.ValueString(), plan.Description.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating role",
			"Could not update role, unexpected error: "+err.Error(),
		)
		return
	}

	// Convert types.List to []string
	var permissions []string
	if !plan.Permissions.IsNull() && !plan.Permissions.IsUnknown() {
		diags := plan.Permissions.ElementsAs(ctx, &permissions, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Update role permissions if provided
	if len(permissions) > 0 {
		err = r.client.UpdateRolePermissions(ctx, role.Id, permissions)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating role permissions",
				"Could not update role permissions, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Convert []string back to types.List
	permissionsList := make([]attr.Value, len(permissions))
	for i, permission := range permissions {
		permissionsList[i] = types.StringValue(permission)
	}

	plan = RoleResourceModel{
		RoleId:        types.StringValue(role.Id),
		Name:          types.StringValue(role.Name),
		Description:   types.StringValue(role.Description),
		Key:           types.StringValue(role.Key),
		IsDefaultRole: types.BoolValue(role.IsDefaultRole),
		Permissions:   types.ListValueMust(types.StringType, permissionsList),
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *RoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RoleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteRole(ctx, data.RoleId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting role",
			"Could not delete role: "+err.Error(),
		)
		return
	}
}
