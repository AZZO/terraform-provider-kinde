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
	_ resource.Resource              = &ApiResource{}
	_ resource.ResourceWithConfigure = &ApiResource{}
)

// ApiResource is the resource implementation.
type ApiResource struct {
	client *kinde_client.Client
}

// ApiResourceModel describes the resource data model.
type ApiResourceModel struct {
	ApiId    types.String `tfsdk:"api_id"`
	Name     types.String `tfsdk:"name"`
	Audience types.String `tfsdk:"audience"`
}

// NewApiResource is a helper function to simplify the provider implementation.
func NewApiResource() resource.Resource {
	return &ApiResource{}
}

// Configure adds the provider configured client to the resource.
func (r *ApiResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ApiResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api"
}

// Schema defines the schema for the resource.
func (r *ApiResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Kinde API.",
		Attributes: map[string]schema.Attribute{
			"api_id": schema.StringAttribute{
				Description: "Unique identifier for the API.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description:         "The name of the API.",
				Required:            true,
				MarkdownDescription: "The name of the API.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"audience": schema.StringAttribute{
				Description:         "The audience for the API. Must be between 1 and 64 characters.",
				Required:            true,
				MarkdownDescription: "The audience for the API. Must be between 1 and 64 characters.",
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
func (r *ApiResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ApiResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	api, err := r.client.CreateApi(
		ctx,
		plan.Name.ValueString(),
		plan.Audience.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating API",
			"Could not create API, unexpected error: "+err.Error(),
		)
		return
	}

	plan = ApiResourceModel{
		ApiId:    types.StringValue(api.Id),
		Name:     types.StringValue(api.Name),
		Audience: types.StringValue(api.Audience),
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *ApiResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ApiResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	api, err := r.client.GetApi(ctx, state.ApiId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading API",
			"Could not read API ID "+state.ApiId.ValueString()+": "+err.Error(),
		)
		return
	}

	state.ApiId = types.StringValue(api.Id)
	state.Name = types.StringValue(api.Name)
	state.Audience = types.StringValue(api.Audience)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ApiResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ApiResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteApi(ctx, data.ApiId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting API",
			"Could not delete API: "+err.Error(),
		)
		return
	}
}

// Update is not supported as APIs are immutable
func (r *ApiResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"APIs are immutable and cannot be updated. Any changes will result in the API being replaced.",
	)
}
