package resources

import (
	"context"
	"fmt"
	"regexp"

	"github.com/AZZO/terraform-provider-kinde/kinde_client"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
	ApplicationId types.String   `tfsdk:"application_id"`
	Name          types.String   `tfsdk:"name"`
	Type          types.String   `tfsdk:"type"`
	ClientId      types.String   `tfsdk:"client_id"`
	ClientSecret  types.String   `tfsdk:"client_secret"`
	LogoutUris    []types.String `tfsdk:"logout_uris"`
	RedirectUris  []types.String `tfsdk:"redirect_uris"`
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
			"logout_uris": schema.ListAttribute{
				Description: "List of logout URIs for the application. Must be valid URLs.",
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.List{
					listvalidator.SizeAtMost(100),
					listvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^https?:\/\/[^\s\/$.?#].[^\s]*$`),
							"must be a valid URL",
						),
					),
				},
			},
			"redirect_uris": schema.ListAttribute{
				Description: "List of redirect URIs for the application. Must be valid URLs.",
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.List{
					listvalidator.SizeAtMost(100),
					listvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^https?:\/\/[^\s\/$.?#].[^\s]*$`),
							"must be a valid URL",
						),
					),
				},
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

	// Convert types.String slices to []string
	logoutUris := make([]string, len(plan.LogoutUris))
	for i, uri := range plan.LogoutUris {
		logoutUris[i] = uri.ValueString()
	}

	redirectUris := make([]string, len(plan.RedirectUris))
	for i, uri := range plan.RedirectUris {
		redirectUris[i] = uri.ValueString()
	}

	// Update the application with URIs if provided
	if len(logoutUris) > 0 || len(redirectUris) > 0 {
		callbacks := kinde_client.Callbacks{
			LogoutUris:   logoutUris,
			RedirectUris: redirectUris,
		}
		_, err = r.client.UpdateCallbacks(ctx, application.Id, callbacks)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating application URIs",
				"Could not update application URIs, unexpected error: "+err.Error(),
			)
			return
		}
	}

	plan = ApplicationResourceModel{
		ApplicationId: types.StringValue(application.Id),
		Name:          types.StringValue(application.Name),
		Type:          types.StringValue(application.Type),
		ClientId:      types.StringValue(application.ClientId),
		ClientSecret:  types.StringValue(application.ClientSecret),
		LogoutUris:    plan.LogoutUris,
		RedirectUris:  plan.RedirectUris,
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

	// Get callbacks
	callbacks, err := r.client.GetCallbacks(ctx, state.ApplicationId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading application callbacks",
			"Could not read application callbacks: "+err.Error(),
		)
		return
	}

	// Convert []string to []types.String
	logoutUris := make([]types.String, len(callbacks.LogoutUris))
	for i, uri := range callbacks.LogoutUris {
		logoutUris[i] = types.StringValue(uri)
	}

	redirectUris := make([]types.String, len(callbacks.RedirectUris))
	for i, uri := range callbacks.RedirectUris {
		redirectUris[i] = types.StringValue(uri)
	}

	state.ApplicationId = types.StringValue(application.Id)
	state.Name = types.StringValue(application.Name)
	state.Type = types.StringValue(application.Type)
	state.ClientId = types.StringValue(application.ClientId)
	state.ClientSecret = types.StringValue(application.ClientSecret)
	state.LogoutUris = logoutUris
	state.RedirectUris = redirectUris

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *ApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ApplicationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ApplicationResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update application name if changed
	if plan.Name.ValueString() != state.Name.ValueString() {
		application, err := r.client.UpdateApplication(ctx, state.ApplicationId.ValueString(), plan.Name.ValueString(), nil, nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating application",
				"Could not update application, unexpected error: "+err.Error(),
			)
			return
		}
		plan.ApplicationId = types.StringValue(application.Id)
		plan.Name = types.StringValue(application.Name)
		plan.Type = types.StringValue(application.Type)
		plan.ClientId = types.StringValue(application.ClientId)
		plan.ClientSecret = types.StringValue(application.ClientSecret)
	}

	// Convert types.String slices to []string
	logoutUris := make([]string, len(plan.LogoutUris))
	for i, uri := range plan.LogoutUris {
		logoutUris[i] = uri.ValueString()
	}

	redirectUris := make([]string, len(plan.RedirectUris))
	for i, uri := range plan.RedirectUris {
		redirectUris[i] = uri.ValueString()
	}

	// Update callbacks if changed
	callbacks := kinde_client.Callbacks{
		LogoutUris:   logoutUris,
		RedirectUris: redirectUris,
	}
	_, err := r.client.UpdateCallbacks(ctx, state.ApplicationId.ValueString(), callbacks)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating application callbacks",
			"Could not update application callbacks, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ApplicationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteApplication(ctx, state.ApplicationId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting application",
			"Could not delete application, unexpected error: "+err.Error(),
		)
		return
	}
}
