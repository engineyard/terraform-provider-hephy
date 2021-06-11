package hephy

import (
	//"fmt"

	"github.com/dchest/uniuri"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getAlphaNumericString(len int) string {
	firstPart := uniuri.NewLenChars(1, []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"))
	lastPart := uniuri.NewLen(len - 1)
	return firstPart + lastPart
}

func diagFromError(summary string, err error) diag.Diagnostic {
	return diag.Diagnostic{
		Severity: diag.Error,
		Summary:  summary,
		Detail:   err.Error(),
	}
}

func rollbackAttributes(d *schema.ResourceData, attributes ...string) {
	for _, attribute := range attributes {
		old, _ := d.GetChange(attribute)
		d.Set(attribute, old)
	}
}
