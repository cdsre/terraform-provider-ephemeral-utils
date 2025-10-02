// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &RevealerResource{}
var _ resource.ResourceWithImportState = &RevealerResource{}
var _ resource.ResourceWithModifyPlan = &RevealerResource{}

func NewRevealerResource() resource.Resource {
	return &RevealerResource{}
}

// RevealerResource defines the resource implementation.
type RevealerResource struct{}

// RevealerModel describes the resource data model.
type RevealerModel struct {
	DataWo types.String `tfsdk:"data_wo"`
	Data   types.String `tfsdk:"data"`
	Id     types.String `tfsdk:"id"`
}

func (r *RevealerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_revealer"
}

func (r *RevealerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Revealer resource",

		Attributes: map[string]schema.Attribute{
			"data_wo": schema.StringAttribute{
				MarkdownDescription: "A non sensitive attribute from an ephemeral resource",
				WriteOnly:           true,
				Required:            true,
			},
			"data": schema.StringAttribute{
				MarkdownDescription: "A persisted attribute from an ephemeral resource that is not sensitive.",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *RevealerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Nothing to configure
}

func (r *RevealerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RevealerModel
	var config RevealerModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Id = config.DataWo

	// set the data attribute to the value of data_wo to persist it and use it in non WO attributes in other resources
	data.Data = config.DataWo

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RevealerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RevealerModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RevealerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RevealerModel
	var config RevealerModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	data.Data = config.DataWo
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RevealerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RevealerModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *RevealerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(req.ID))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("data"), types.StringValue(req.ID))...)
}

// ModifyPlan ensures changes to write-only data_wo cause a planned change on computed data.
func (r *RevealerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// If creating (no prior state) or destroying (no plan), default behavior is fine
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	// Read config to access data_wo
	var config RevealerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If data_wo is known in config, project it into the planned value for data
	if !config.DataWo.IsNull() && !config.DataWo.IsUnknown() {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("data"), config.DataWo)...)
	}
}
