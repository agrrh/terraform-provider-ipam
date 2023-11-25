// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	goipam "github.com/metal-stack/go-ipam"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AllocationResource{}
var _ resource.ResourceWithImportState = &AllocationResource{}

func NewAllocationResource() resource.Resource {
	return &AllocationResource{}
}

// AllocationResource defines the resource implementation.
type AllocationResource struct {
	ipam goipam.Ipamer
}

// AllocationResourceModel describes the resource data model.
type AllocationResourceModel struct {
	PoolId types.String `tfsdk:"pool_id"`
	Size   types.Int64  `tfsdk:"size"`
	Id     types.String `tfsdk:"id"`
	CIDR   types.String `tfsdk:"cidr"`
}

func (r *AllocationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_allocation"
}

func (r *AllocationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "An IPAM Allocation represents an sub-ranges from previously defined IP Pool",

		Attributes: map[string]schema.Attribute{
			"pool_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Pool identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"size": schema.Int64Attribute{
				MarkdownDescription: "Size of allocation",
				Required:            true,
				Validators: []validator.Int64{
					int64validator.Between(8, 32),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Allocation identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cidr": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Allocation CIDR",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *AllocationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	ipam, ok := req.ProviderData.(goipam.Ipamer)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.ipam = ipam
}

func (r *AllocationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AllocationResourceModel

	mutex.Lock()
	defer mutex.Unlock()

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	prefix, err := r.ipam.AcquireChildPrefix(ctx, data.PoolId.ValueString(), uint8(data.Size.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource", fmt.Sprintf("... details ... %s", err))
	}

	data.Id = types.StringValue(prefix.Cidr)
	data.CIDR = types.StringValue(prefix.Cidr)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AllocationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AllocationResourceModel

	mutex.RLock()
	defer mutex.RUnlock()

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	parentPrefix := r.ipam.PrefixFrom(ctx, data.PoolId.ValueString())

	if parentPrefix == nil {
		return
	}

	// FIXME
	// data.CIDR = types.StringValue(prefix.Cidr)
	// data.Id = types.StringValue(prefix.Cidr)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AllocationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AllocationResourceModel

	mutex.Lock()
	defer mutex.Unlock()

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AllocationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AllocationResourceModel

	mutex.Lock()
	defer mutex.Unlock()

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	child := &goipam.Prefix{
		Cidr:       data.CIDR.ValueString(),
		ParentCidr: data.PoolId.ValueString(),
	}
	err := r.ipam.ReleaseChildPrefix(ctx, child)
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource", fmt.Sprintf("... details ... %s", err))
	}
}

func (r *AllocationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
