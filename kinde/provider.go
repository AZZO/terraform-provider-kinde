package kinde

import (
	"context"

	"github.com/AZZO/terraform-provider-kinde/kinde_client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure KindeProvider satisfies various provider interfaces.
var _ provider.Provider = &KindeProvider{}
var _ provider.ProviderWithFunctions = &KindeProvider{}

// KindeProvider defines the provider implementation.
type KindeProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// KindeProviderModel describes the provider data model.
type KindeProviderModel struct {
	IssuerUrl    types.String `tfsdk:"issuer_url"`
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

func (p *KindeProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "kinde"
	resp.Version = p.version
}

func (p *KindeProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"issuer_url": schema.StringAttribute{
				Required: true,
			},
			"client_id": schema.StringAttribute{
				Required: true,
			},
			"client_secret": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *KindeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data KindeProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	if data.IssuerUrl.IsNull() {
		resp.Diagnostics.AddAttributeError(path.Root("issuer_url"),
			"Unknown Issuer URL",
			"Issuer URL must be provided")
	}
	if data.ClientId.IsNull() {
		resp.Diagnostics.AddAttributeError(path.Root("client_id"),
			"Unknown Client ID",
			"Client ID must be provided")
	}
	if data.ClientSecret.IsNull() {
		resp.Diagnostics.AddAttributeError(path.Root("client_secret"),
			"Unknown Client SEcret",
			"Client Secret must be provided")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := kinde_client.NewClient(ctx, data.IssuerUrl.ValueString(), data.ClientId.ValueString(), data.ClientSecret.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Kinde API Client",
			"An unexpected error occurred when creating the Kinde API client: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *KindeProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewApplicationResource,
	}
}

func (p *KindeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewApplicationDataSource,
	}
}

func (p *KindeProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &KindeProvider{
			version: version,
		}
	}
}
