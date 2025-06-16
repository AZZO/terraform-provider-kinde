package resources

import (
	"context"
	"fmt"

	"github.com/AZZO/terraform-provider-kinde/kinde_client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &ApplicationResource{}
	_ resource.ResourceWithConfigure = &ApplicationResource{}
)

// ApplicationResource is the resource implementation.
type ApplicationResource struct {
	client *kinde_client.Client
}

// ApplicationResourceModel describes the resource data model.
type ApplicationResourceModel struct {
	ApplicationId types.String `tfsdk:"application_id"`
	Name          types.String `tfsdk:"name"`
	Type          types.String `tfsdk:"type"`
	ClientId      types.String `tfsdk:"client_id"`
	ClientSecret  types.String `tfsdk:"client_secret"`
}

// NewApplicationResource is a helper function to simplify the provider implementation.
func NewApplicationResource() resource.Resource {
	return &ApplicationResource{}
}

// Configure adds the provider configured client to the resource.
func (r *ApplicationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ApplicationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

// Schema defines the schema for the resource.
func (r *ApplicationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Kinde application.",
		Attributes: map[string]schema.Attribute{
			"application_id": schema.StringAttribute{
				Description: "Unique identifier for the application.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Human-readable name of the application.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "Type of the application. Must be one of: reg (Regular), m2m (Machine to Machine), spa (Single Page Application).",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("reg", "m2m", "spa"),
				},
			},
			"client_id": schema.StringAttribute{
				Description: "Client ID for the application.",
				Optional:    true,
				Computed:    true,
			},
			"client_secret": schema.StringAttribute{
				Description: "Client Secret for the application.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *ApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ApplicationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	application, err := r.client.CreateApplication(ctx, plan.Name.ValueString(), plan.Type.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating application",
			"Could not create application, unexpected error: "+err.Error(),
		)
		return
	}

	plan = ApplicationResourceModel{
		ApplicationId: types.StringValue(application.Id),
		Name:          types.StringValue(application.Name),
		Type:          types.StringValue(application.Type),
		ClientId:      types.StringValue(application.ClientId),
		ClientSecret:  types.StringValue(application.ClientSecret),
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *ApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ApplicationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	application, err := r.client.GetApplication(ctx, state.ApplicationId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading application",
			"Could not read application ID "+state.ApplicationId.ValueString()+": "+err.Error(),
		)
		return
	}

	state.ApplicationId = types.StringValue(application.Id)
	state.Name = types.StringValue(application.Name)
	state.Type = types.StringValue(application.Type)
	state.ClientId = types.StringValue(application.ClientId)
	state.ClientSecret = types.StringValue(application.ClientSecret)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *ApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// TODO: Implement update functionality
	resp.Diagnostics.AddWarning(
		"Update not implemented",
		"Update functionality is not yet implemented for this resource.",
	)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// TODO: Implement delete functionality
	resp.Diagnostics.AddWarning(
		"Delete not implemented",
		"Delete functionality is not yet implemented for this resource.",
	)
}
