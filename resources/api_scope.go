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
	_ resource.Resource              = &ApiScopeResource{}
	_ resource.ResourceWithConfigure = &ApiScopeResource{}
)

// ApiScopeResource is the resource implementation.
type ApiScopeResource struct {
	client *kinde_client.Client
}

// ApiScopeResourceModel describes the resource data model.
type ApiScopeResourceModel struct {
	ScopeId     types.String `tfsdk:"scope_id"`
	ApiId       types.String `tfsdk:"api_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

// NewApiScopeResource is a helper function to simplify the provider implementation.
func NewApiScopeResource() resource.Resource {
	return &ApiScopeResource{}
}

// Configure adds the provider configured client to the resource.
func (r *ApiScopeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ApiScopeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_scope"
}

// Schema defines the schema for the resource.
func (r *ApiScopeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Kinde API Scope.",
		Attributes: map[string]schema.Attribute{
			"scope_id": schema.StringAttribute{
				Description: "Unique identifier for the scope.",
				Computed:    true,
			},
			"api_id": schema.StringAttribute{
				Description:         "The ID of the API this scope belongs to.",
				Required:            true,
				MarkdownDescription: "The ID of the API this scope belongs to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "The name of the scope.",
				Required:            true,
				MarkdownDescription: "The name of the scope.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			"description": schema.StringAttribute{
				Description:         "The description of the scope.",
				Optional:            true,
				MarkdownDescription: "The description of the scope.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *ApiScopeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ApiScopeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	scope, err := r.client.CreateApiScope(
		ctx,
		plan.ApiId.ValueString(),
		plan.Name.ValueString(),
		plan.Description.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating API scope",
			"Could not create API scope, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ScopeId = types.StringValue(scope.Id)
	plan.Name = types.StringValue(scope.Name)
	plan.Description = types.StringValue(scope.Description)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *ApiScopeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ApiScopeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get all scopes for the API
	scopes, err := r.client.GetApiScopes(ctx, state.ApiId.ValueString())
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
		if scope.Id == state.ScopeId.ValueString() {
			foundScope = &scope
			break
		}
	}

	if foundScope == nil {
		resp.Diagnostics.AddError(
			"Error Reading API Scope",
			"Could not find API scope with ID "+state.ScopeId.ValueString(),
		)
		return
	}

	state.ScopeId = types.StringValue(foundScope.Id)
	state.Name = types.StringValue(foundScope.Name)
	state.Description = types.StringValue(foundScope.Description)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *ApiScopeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ApiScopeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Since scopes are immutable, we need to delete the old one and create a new one
	err := r.client.DeleteApiScope(ctx, plan.ApiId.ValueString(), plan.ScopeId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting old API scope",
			"Could not delete old API scope: "+err.Error(),
		)
		return
	}

	scope, err := r.client.CreateApiScope(
		ctx,
		plan.ApiId.ValueString(),
		plan.Name.ValueString(),
		plan.Description.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating new API scope",
			"Could not create new API scope: "+err.Error(),
		)
		return
	}

	plan.ScopeId = types.StringValue(scope.Id)
	plan.Name = types.StringValue(scope.Name)
	plan.Description = types.StringValue(scope.Description)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ApiScopeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ApiScopeResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteApiScope(ctx, data.ApiId.ValueString(), data.ScopeId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting API scope",
			"Could not delete API scope: "+err.Error(),
		)
		return
	}
}
