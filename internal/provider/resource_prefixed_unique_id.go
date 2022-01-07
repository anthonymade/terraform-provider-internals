package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func prefixedUniqueId() *schema.Resource {
	return &schema.Resource{
		Description: `
The resource ` + "`internals_prefixed_unique_id`" + ` exposes the internal sdk helper function
` + "`PrefixedUniqueId`" + ` which is used by resources such as aws_s3_bucket to name a bucket from
the bucket_prefix argument.

The prefix will be appended with ` + fmt.Sprintf("%v", resource.UniqueIDSuffixLength) + `
characters.

This resource was created to solve the issue of circular dependencies i.e. when using bucket_prefix and
requiring a logging configuration to use the bucket name for the target_prefix.

To use this resource, define it and use the ` + "`id`" + ` attribute to set resource name/id instead
of using the built in *_prefix arguments.
`,

		CreateContext: resourcePrefixedUniqueIdCreate,
		ReadContext:   resourcePrefixedUniqueIdRead,
		DeleteContext: resourcePrefixedUniqueIdDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePrefixedUniqueIdImport,
		},

		Schema: map[string]*schema.Schema{
			"prefix": {
				Description: "Id prefix",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"id": {
				Description: "The generated id",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourcePrefixedUniqueIdCreate(ctx context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	prefix := d.Get("prefix").(string)
	id := resource.PrefixedUniqueId(prefix)
	d.SetId(id)
	tflog.Trace(ctx, fmt.Sprintf("created an internals_prefixed_unique_id resource with id %v", id))
	return nil
}

func resourcePrefixedUniqueIdRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePrefixedUniqueIdDelete(ctx context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	tflog.Trace(ctx, fmt.Sprintf("removing an internals_prefixed_unique_id resource with id %v", d.Get("prefix").(string)))
	d.SetId("")
	return nil
}

func resourcePrefixedUniqueIdImport(ctx context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	tflog.Trace(ctx, fmt.Sprintf("importing an internals_prefixed_unique_id resource with id %v", d.Get("id").(string)))
	return []*schema.ResourceData{d}, nil
}
