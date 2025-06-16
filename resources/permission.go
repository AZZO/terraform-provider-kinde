package resources

import (
	"context"
	"fmt"

	"github.com/AZZO/terraform-provider-kinde/kinde_client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &PermissionResource{}
	_ resource.ResourceWithConfigure = &PermissionResource{}
)

// PermissionResource is the resource implementation.
type PermissionResource struct {
	client *kinde_client.Client
}

// PermissionResourceModel describes the resource data model.
type PermissionResourceModel struct {
	PermissionId types.String `tfsdk:"permission_id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Key          types.String `tfsdk:"key"`
}

// NewPermissionResource is a helper function to simplify the provider implementation.
func NewPermissionResource() resource.Resource {
	return &PermissionResource{}
}

// Configure adds the provider configured client to the resource.
func (r *PermissionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *PermissionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_permission"
}

// Schema defines the schema for the resource.
func (r *PermissionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Kinde Permission.",
		Attributes: map[string]schema.Attribute{
			"permission_id": schema.StringAttribute{
				Description: "Unique identifier for the permission.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description:         "The name of the permission.",
				Required:            true,
				MarkdownDescription: "The name of the permission.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			"description": schema.StringAttribute{
				Description:         "The description of the permission.",
				Optional:            true,
				MarkdownDescription: "The description of the permission.",
			},
			"key": schema.StringAttribute{
				Description:         "The key of the permission.",
				Required:            true,
				MarkdownDescription: "The key of the permission.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *PermissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PermissionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permission, err := r.client.CreatePermission(
		ctx,
		plan.Name.ValueString(),
		plan.Key.ValueString(),
		plan.Description.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating permission",
			"Could not create permission, unexpected error: "+err.Error(),
		)
		return
	}

	// Update the plan with the response data
	plan.PermissionId = types.StringValue(permission.Id)
	plan.Name = types.StringValue(permission.Name)
	plan.Description = types.StringValue(permission.Description)
	plan.Key = types.StringValue(permission.Key)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *PermissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PermissionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permission, err := r.client.GetPermission(ctx, state.PermissionId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Permission",
			"Could not read permission ID "+state.PermissionId.ValueString()+": "+err.Error(),
		)
		return
	}

	state.PermissionId = types.StringValue(permission.Id)
	state.Name = types.StringValue(permission.Name)
	state.Description = types.StringValue(permission.Description)
	state.Key = types.StringValue(permission.Key)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *PermissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PermissionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permission, err := r.client.UpdatePermission(
		ctx,
		plan.PermissionId.ValueString(),
		plan.Name.ValueString(),
		plan.Description.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating permission",
			"Could not update permission, unexpected error: "+err.Error(),
		)
		return
	}

	// Update the plan with the response data
	plan.PermissionId = types.StringValue(permission.Id)
	plan.Name = types.StringValue(permission.Name)
	plan.Description = types.StringValue(permission.Description)
	plan.Key = types.StringValue(permission.Key)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *PermissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state PermissionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeletePermission(ctx, state.PermissionId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting permission",
			"Could not delete permission, unexpected error: "+err.Error(),
		)
		return
	}
}
