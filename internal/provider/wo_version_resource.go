package provider

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/crypto/argon2"
)

// OWASP-recommended Argon2id parameters. Changing them changes every hash.
const (
	argonTime    = 2
	argonMemory  = 19 * 1024
	argonThreads = 1
	argonKeyLen  = 32
	saltLen      = 16
)

var (
	_ resource.Resource               = (*woVersionResource)(nil)
	_ resource.ResourceWithModifyPlan = (*woVersionResource)(nil)
)

type woVersionResource struct{}

type woVersionModel struct {
	ValueWO types.String `tfsdk:"value_wo"`
	Salt    types.String `tfsdk:"salt"`
	Hash    types.String `tfsdk:"hash"`
	Version types.Int64  `tfsdk:"version"`
}

func NewWoVersionResource() resource.Resource {
	return &woVersionResource{}
}

func (r *woVersionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_wo_version"
}

func (r *woVersionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Derives a stable, non-sensitive version from a write-only value, for use as a `*_wo_version` argument. " +
			"The version changes only when the value changes.",
		Attributes: map[string]schema.Attribute{
			"value_wo": schema.StringAttribute{
				Description: "Write-only value to fingerprint. Accepts ephemeral values; never stored in plan or state.",
				Required:    true,
				WriteOnly:   true,
			},
			"salt": schema.StringAttribute{
				Description: "Random base64 salt generated on create.",
				Computed:    true,
			},
			"hash": schema.StringAttribute{
				Description: "Hex-encoded Argon2id hash of `value_wo`.",
				Computed:    true,
			},
			"version": schema.Int64Attribute{
				Description: "First 8 hex characters of `hash` as an integer (0 to 4294967295).",
				Computed:    true,
			},
		},
	}
}

// ModifyPlan recomputes the fingerprint with the salt from state so that
// changes are known at plan time. On create the salt doesn't exist yet, so
// the outputs stay unknown until apply.
func (r *woVersionResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}

	var state woVersionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	var value types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("value_wo"), &value)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan := woVersionModel{ValueWO: types.StringNull(), Salt: state.Salt}
	if value.IsUnknown() {
		plan.Hash = types.StringUnknown()
		plan.Version = types.Int64Unknown()
	} else if err := plan.compute(value.ValueString()); err != nil {
		resp.Diagnostics.AddError("Failed to hash value_wo", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
}

func (r *woVersionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		resp.Diagnostics.AddError("Failed to generate salt", err.Error())
		return
	}
	r.apply(ctx, req.Config, base64.StdEncoding.EncodeToString(salt), &resp.State, &resp.Diagnostics)
}

func (r *woVersionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state woVersionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.apply(ctx, req.Config, state.Salt.ValueString(), &resp.State, &resp.Diagnostics)
}

func (r *woVersionResource) Read(context.Context, resource.ReadRequest, *resource.ReadResponse) {}

func (r *woVersionResource) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {
}

func (r *woVersionResource) apply(ctx context.Context, config tfsdk.Config, salt string, state *tfsdk.State, diags *diag.Diagnostics) {
	var value types.String
	diags.Append(config.GetAttribute(ctx, path.Root("value_wo"), &value)...)
	if diags.HasError() {
		return
	}

	m := woVersionModel{ValueWO: types.StringNull(), Salt: types.StringValue(salt)}
	if err := m.compute(value.ValueString()); err != nil {
		diags.AddError("Failed to hash value_wo", err.Error())
		return
	}
	diags.Append(state.Set(ctx, &m)...)
}

func (m *woVersionModel) compute(value string) error {
	salt, err := base64.StdEncoding.DecodeString(m.Salt.ValueString())
	if err != nil {
		return fmt.Errorf("decoding salt: %w", err)
	}
	if len(salt) == 0 {
		return fmt.Errorf("empty salt")
	}
	hash, version := fingerprint(value, salt)
	m.Hash = types.StringValue(hash)
	m.Version = types.Int64Value(version)
	return nil
}

func fingerprint(value string, salt []byte) (string, int64) {
	sum := argon2.IDKey([]byte(value), salt, argonTime, argonMemory, argonThreads, argonKeyLen)
	return hex.EncodeToString(sum), int64(binary.BigEndian.Uint32(sum[:4]))
}
