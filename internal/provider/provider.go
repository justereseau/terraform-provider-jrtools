package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = (*jrtoolsProvider)(nil)

type jrtoolsProvider struct {
	version string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &jrtoolsProvider{version: version}
	}
}

func (p *jrtoolsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "jrtools"
	resp.Version = p.version
}

func (p *jrtoolsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Utility resources for Terraform workflows.",
	}
}

func (p *jrtoolsProvider) Configure(context.Context, provider.ConfigureRequest, *provider.ConfigureResponse) {
}

func (p *jrtoolsProvider) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewWoVersionResource,
	}
}

func (p *jrtoolsProvider) DataSources(context.Context) []func() datasource.DataSource {
	return nil
}
