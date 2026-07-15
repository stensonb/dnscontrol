package spaceship

import (
	"github.com/DNSControl/dnscontrol/v4/models"
//	"github.com/DNSControl/dnscontrol/v4/pkg/rejectif"
)

// AuditRecords returns a list of errors corresponding to the records
// that aren't supported by this provider.  If all records are
// supported, an empty list is returned.
//
// see providers/vultr/auditrecords.go for example
func AuditRecords(records []*models.RecordConfig) []error {
	return []error{}
}
