package provider

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestModelComputeInvalidSalt(t *testing.T) {
	cases := map[string]string{
		"not base64": "not-base64!!!",
		"empty":      "",
	}
	for name, salt := range cases {
		t.Run(name, func(t *testing.T) {
			m := &woVersionModel{Salt: types.StringValue(salt)}
			if err := m.compute("secret"); err == nil {
				t.Fatalf("expected error for salt %q", salt)
			}
		})
	}
}

func TestModelComputeValidSalt(t *testing.T) {
	rawSalt := []byte("0123456789abcdef")
	m := &woVersionModel{Salt: types.StringValue(base64.StdEncoding.EncodeToString(rawSalt))}
	if err := m.compute("secret"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantHash, wantVersion := fingerprint("secret", rawSalt)
	if m.Hash.ValueString() != wantHash {
		t.Fatalf("hash = %s, want %s", m.Hash.ValueString(), wantHash)
	}
	if m.Version.ValueInt64() != wantVersion {
		t.Fatalf("version = %d, want %d", m.Version.ValueInt64(), wantVersion)
	}
}

// Read and Delete are deliberate no-ops: the fingerprint has no backing
// system to read from or delete. These guard against that remaining true.
func TestReadAndDeleteAreNoops(t *testing.T) {
	r := &woVersionResource{}

	var readResp resource.ReadResponse
	r.Read(context.Background(), resource.ReadRequest{}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read produced diagnostics: %v", readResp.Diagnostics)
	}

	var deleteResp resource.DeleteResponse
	r.Delete(context.Background(), resource.DeleteRequest{}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete produced diagnostics: %v", deleteResp.Diagnostics)
	}
}
