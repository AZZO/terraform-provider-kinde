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

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &ScaffoldingProvider{}
var _ provider.ProviderWithFunctions = &ScaffoldingProvider{}

// ScaffoldingProvider defines the provider implementation.
type ScaffoldingProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ScaffoldingProviderModel describes the provider data model.
type ScaffoldingProviderModel struct {
	IssuerUrl    types.String `tfsdk:"endpoint"`
	ClientId     types.String `tfsdk:"client-id"`
	ClientSecret types.String `tfsdk:"client-secret"`
}

func (p *ScaffoldingProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "scaffolding"
	resp.Version = p.version
}

func (p *ScaffoldingProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"issuer-url": schema.StringAttribute{},
			"client-id":  schema.StringAttribute{},
			"client-secret": schema.StringAttribute{
				Sensitive: true,
			},
		},
	}
}

func (p *ScaffoldingProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ScaffoldingProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	if data.IssuerUrl.IsNull() {
		resp.Diagnostics.AddAttributeError(path.Root("issuer-url"),
			"Unknown Issuer URL",
			"Issuer URL must be provided")
	}
	if data.ClientId.IsNull() {
		resp.Diagnostics.AddAttributeError(path.Root("client-id"),
			"Unknown Client ID",
			"Client ID must be provided")
	}
	if data.ClientSecret.IsNull() {
		resp.Diagnostics.AddAttributeError(path.Root("client-secret"),
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

func (p *ScaffoldingProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *ScaffoldingProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewApplicationDataSource,
	}
}

func (p *ScaffoldingProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ScaffoldingProvider{
			version: version,
		}
	}
}
