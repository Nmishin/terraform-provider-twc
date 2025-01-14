package provider

import (
    "context"
    "terraform-provider-twc/internal/resource_server"

    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = (*serverResource)(nil)

func NewServerResource() resource.Resource {
    return &serverResource{}
}

type serverResource struct{}

func (r *serverResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_server"
}

func (r *serverResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = resource_server.ServerResourceSchema(ctx)
}

func (r *serverResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data resource_server.ServerModel

    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    resp.Diagnostics.Append(callServerAPI(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *serverResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data resource_server.ServerModel

    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *serverResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data resource_server.ServerModel

    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    resp.Diagnostics.Append(callServerAPI(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *serverResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data resource_server.ServerModel

    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Typically this method would contain logic that makes an HTTP call to a remote API, and then stores
// computed results back to the data model. For example purposes, this function just sets all unknown
// Pet values to null to avoid data consistency errors.
func callServerAPI(ctx context.Context, server *resource_server.ServerModel) diag.Diagnostics {
    if server.Id.IsUnknown() {
        server.Id = types.Int64Null()
    }

    if server.Status.IsUnknown() {
        server.Status = types.StringNull()
    }

    if server.Tags.IsUnknown() {
        server.Tags = types.ListNull(resource_server.TagsValue{}.Type(ctx))
    } else if !server.Tags.IsNull() {
        var tags []resource_server.TagsValue
        diags := server.Tags.ElementsAs(ctx, &tags, false)
        if diags.HasError() {
            return diags
        }

        for i := range tags {
            if tags[i].Id.IsUnknown() {
                tags[i].Id = types.Int64Null()
            }

            if tags[i].Name.IsUnknown() {
                tags[i].Name = types.StringNull()
            }
        }

        server.Tags, diags = types.ListValueFrom(ctx, resource_server.TagsValue{}.Type(ctx), tags)
        if diags.HasError() {
            return diags
        }
    }

    if server.Category.IsUnknown() {
        server.Category = resource_server.NewCategoryValueNull()
    } else if !server.Category.IsNull() {
        if server.Category.Id.IsUnknown() {
            server.Category.Id = types.Int64Null()
        }

        if server.Category.Name.IsUnknown() {
            server.Category.Name = types.StringNull()
        }
    }

    return nil
}

