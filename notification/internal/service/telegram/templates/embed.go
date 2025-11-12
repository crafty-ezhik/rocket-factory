package templates

import "embed"

//go:embed order_paid_notification.tmpl order_assembled_notification.tmpl
var FS embed.FS
