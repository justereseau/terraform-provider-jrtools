package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = (*toolsProvider)(nil)

type toolsProvider struct {
	version string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &toolsProvider{version: version}
	}
}

func (p *toolsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "tools"
	resp.Version = p.version
}

func (p *toolsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Utility resources for Terraform workflows.",
	}
}

func (p *toolsProvider) Configure(context.Context, provider.ConfigureRequest, *provider.ConfigureResponse) {
}

func (p *toolsProvider) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewWoVersionResource,
	}
}

func (p *toolsProvider) DataSources(context.Context) []func() datasource.DataSource {
	return nil
}
