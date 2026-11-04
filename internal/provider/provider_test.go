package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"jrtools": providerserver.NewProtocol6WithError(New("test")()),
}

// Configure is a deliberate no-op: the provider has no schema to read
// configuration from and nothing to set up. This guards against that
// remaining true.
func TestProviderConfigureIsNoop(t *testing.T) {
	p := &jrtoolsProvider{version: "test"}

	var resp provider.ConfigureResponse
	p.Configure(context.Background(), provider.ConfigureRequest{}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Configure produced diagnostics: %v", resp.Diagnostics)
	}
}
