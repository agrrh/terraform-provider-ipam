// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	goipam "github.com/metal-stack/go-ipam"
)

// Ensure IPAMProvider satisfies various provider interfaces.
var _ provider.Provider = &IPAMProvider{}

// Global storage-access lock to use within this plugin.
var mutex sync.RWMutex

// IPAMProvider defines the provider implementation.
type IPAMProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// IPAMProviderModel describes the provider data model.
type IPAMProviderModel struct {
	File types.String `tfsdk:"file"`
}

func (p *IPAMProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ipam"
	resp.Version = p.version
}

func (p *IPAMProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"file": schema.StringAttribute{
				MarkdownDescription: "Path to the file to store IPAM state in",
				Optional:            true,
			},
		},
	}
}

func (p *IPAMProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data IPAMProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	if data.File.IsNull() {
		data.File = types.StringValue("./default.ipam.json")
	}

	fileStorage := goipam.NewLocalFile(ctx, data.File.ValueString())
	ipam := goipam.NewWithStorage(fileStorage)

	// Example client configuration for data sources and resources
	resp.DataSourceData = ipam
	resp.ResourceData = ipam
}

func (p *IPAMProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewPoolResource,
		NewAllocationResource,
	}
}

func (p *IPAMProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// NewExampleDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &IPAMProvider{
			version: version,
		}
	}
}
