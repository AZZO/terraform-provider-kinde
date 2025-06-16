package resources

import (
	"context"
	"fmt"
	"regexp"

	"github.com/AZZO/terraform-provider-kinde/kinde_client"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	LogoutUris    types.List   `tfsdk:"logout_uris"`
	RedirectUris  types.List   `tfsdk:"redirect_uris"`
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

	// Convert types.List to []string
	var logoutUris []string
	if !plan.LogoutUris.IsNull() && !plan.LogoutUris.IsUnknown() {
		diags := plan.LogoutUris.ElementsAs(ctx, &logoutUris, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var redirectUris []string
	if !plan.RedirectUris.IsNull() && !plan.RedirectUris.IsUnknown() {
		diags := plan.RedirectUris.ElementsAs(ctx, &redirectUris, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Update the application with URIs if provided
	if len(logoutUris) > 0 || len(redirectUris) > 0 {
		callbacks := kinde_client.Callbacks{
			LogoutUris:   logoutUris,
			RedirectUris: redirectUris,
		}
		_, err = r.client.UpdateApplicationCallbacks(ctx, application.Id, callbacks)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating application URIs",
				"Could not update application URIs, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Convert []string back to types.List
	logoutUrisList := make([]attr.Value, len(logoutUris))
	for i, uri := range logoutUris {
		logoutUrisList[i] = types.StringValue(uri)
	}

	redirectUrisList := make([]attr.Value, len(redirectUris))
	for i, uri := range redirectUris {
		redirectUrisList[i] = types.StringValue(uri)
	}

	plan = ApplicationResourceModel{
		ApplicationId: types.StringValue(application.Id),
		Name:          types.StringValue(application.Name),
		Type:          types.StringValue(application.Type),
		ClientId:      types.StringValue(application.ClientId),
		ClientSecret:  types.StringValue(application.ClientSecret),
		LogoutUris:    types.ListValueMust(types.StringType, logoutUrisList),
		RedirectUris:  types.ListValueMust(types.StringType, redirectUrisList),
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *ApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ApplicationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get application details
	application, err := r.client.GetApplication(ctx, data.ApplicationId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read application, got error: %s", err))
		return
	}

	// Get application callbacks
	callbacks, err := r.client.GetApplicationCallbacks(ctx, data.ApplicationId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read application callbacks, got error: %s", err))
		return
	}

	// Convert []string to types.List
	logoutUrisList := make([]attr.Value, len(callbacks.LogoutUris))
	for i, uri := range callbacks.LogoutUris {
		logoutUrisList[i] = types.StringValue(uri)
	}

	redirectUrisList := make([]attr.Value, len(callbacks.RedirectUris))
	for i, uri := range callbacks.RedirectUris {
		redirectUrisList[i] = types.StringValue(uri)
	}

	// Update the model with the response data
	data.ApplicationId = types.StringValue(application.Id)
	data.Name = types.StringValue(application.Name)
	data.Type = types.StringValue(application.Type)
	data.ClientId = types.StringValue(application.ClientId)
	data.ClientSecret = types.StringValue(application.ClientSecret)
	data.LogoutUris = types.ListValueMust(types.StringType, logoutUrisList)
	data.RedirectUris = types.ListValueMust(types.StringType, redirectUrisList)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *ApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ApplicationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert types.List to []string
	var logoutUris []string
	if !data.LogoutUris.IsNull() && !data.LogoutUris.IsUnknown() {
		diags := data.LogoutUris.ElementsAs(ctx, &logoutUris, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var redirectUris []string
	if !data.RedirectUris.IsNull() && !data.RedirectUris.IsUnknown() {
		diags := data.RedirectUris.ElementsAs(ctx, &redirectUris, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Update application
	application, err := r.client.UpdateApplication(ctx, data.ApplicationId.ValueString(), data.Name.ValueString(), logoutUris, redirectUris)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update application, got error: %s", err))
		return
	}

	// Update callbacks
	callbacks := kinde_client.Callbacks{
		LogoutUris:   logoutUris,
		RedirectUris: redirectUris,
	}
	_, err = r.client.UpdateApplicationCallbacks(ctx, data.ApplicationId.ValueString(), callbacks)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update application callbacks, got error: %s", err))
		return
	}

	// Map response body to model
	data.ApplicationId = types.StringValue(application.Id)
	data.Name = types.StringValue(application.Name)
	data.Type = types.StringValue(application.Type)
	data.ClientId = types.StringValue(application.ClientId)
	data.ClientSecret = types.StringValue(application.ClientSecret)

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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
