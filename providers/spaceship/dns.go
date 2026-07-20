package spaceship

import (
	"context"
	"fmt"
	"strconv"
	//"time"

	"github.com/DNSControl/dnscontrol/v4/models"
	//"github.com/DNSControl/dnscontrol/v4/pkg/diff2"

	"github.com/namecheap/go-spaceship-sdk/client"
)

func (n *spaceshipProvider) GetZoneRecords(dc *models.DomainConfig) (models.Records, error) {
	client, err := client.NewClient(n.BaseURL, n.ApiKey, n.ApiSecret)
	if err != nil {
		return models.Records{}, err
	}

	records, err := client.GetDNSRecords(context.Background(), dc.Name)
	if err != nil {
		return models.Records{}, err
	}

	return toRecordConfigArray(dc.Name, records)
}

func toRecordConfigArray(domain string, records []client.DNSRecord) (models.Records, error) {
	result := make([]*models.RecordConfig, len(records))

	for i, rec := range records {
		rc := &models.RecordConfig{
			Type:     rec.Type,
			TTL:      uint32(rec.TTL),
			Original: rec, // Store original record for reference/debugging
		}

		// Set the record label (@, sub, etc.) relative to the domain origin
		rc.SetLabel(rec.Name, domain)

		// Set record targets based on type
		switch rec.Type {
		case "MX":
			// Spaceship provides Priority either as a separate field or parsed int
			priority := uint16(*rec.Priority)
			if err := rc.SetTargetMX(priority, rec.Value); err != nil {
				return nil, fmt.Errorf("failed to set MX target for %s: %w", rec.Name, err)
			}

		case "SRV":
			// If SRV parameters are stored in dedicated fields:
			priority := uint16(*rec.Priority)
			weight := uint16(*rec.Weight)

			// Spaceship API is polymorphic here
			// it accepts and returns the port either as a JSON string (e.g. "_443", "*") or as a JSON number
			// TODO: TESTING
			var port uint16
			if rec.Port.String == nil {
				port = uint16(*rec.Port.Int)
			} else {
				// TODO: this won't work with "_443" or "*" -- what do we do here?!
				val, err := strconv.ParseUint(*rec.Port.String, 10, 16)
				if err != nil {
					return nil, fmt.Errorf("failed to convert %v to uint: %w", *rec.Port.String, err)
				}

				// Cast the resulting uint64 into a uint16
				port = uint16(val)
			}

			if err := rc.SetTargetSRV(priority, weight, port, rec.Value); err != nil {
				return nil, fmt.Errorf("failed to set SRV target for %s: %w", rec.Name, err)
			}

		case "CAA":
			// If CAA tag/flag are separate or combined in Value
			tag := rec.Tag
			if tag == "" {
				tag = "issue" // fallback/default if applicable
			}
			flag := uint8(*rec.Flag)

			if err := rc.SetTargetCAA(flag, tag, rec.Value); err != nil {
				return nil, fmt.Errorf("failed to set CAA target for %s: %w", rec.Name, err)
			}

		case "TXT":
			// SetTargetTXT handles quote stripping and multi-string TXT chunks
			if err := rc.SetTargetTXT(rec.Value); err != nil {
				return nil, fmt.Errorf("failed to set TXT target for %s: %w", rec.Name, err)
			}

		default:
			// A, AAAA, CNAME, NS, PTR, etc.
			if err := rc.SetTarget(rec.Value); err != nil {
				return nil, fmt.Errorf("failed to set target for %s (%s): %w", rec.Name, rec.Type, err)
			}
		}

		result[i] = rc
	}

	return result, nil
}

func (n *spaceshipProvider) GetZoneRecordsCorrections(dc *models.DomainConfig, existing models.Records) ([]*models.Correction, int, error) {
	return []*models.Correction{}, 0, nil
/*
	changes, actualChangeCount, err := diff2.ByRecord(existing, dc, nil)
	if err != nil {
		return nil, 0, err
	}

	p := &spaceship.Provider{
		APIKey:    n.ApiKey,
		APISecret: n.ApiSecret,
		BaseURL:   n.BaseURL,
	}

	var corrections []*models.Correction
	for _, change := range changes {
		description := change.MsgsJoined

		var libdnsRecords []libdns.Record
		for _, rc := range change.New {
			libdnsRecords = append(libdnsRecords, libdns.RR{
				Type: rc.Type,
				Name: rc.GetLabel(),
				Data: rc.GetTargetField(),
				TTL:  time.Duration(rc.TTL) * time.Second,
			})
		}

		switch change.Type {
		case diff2.CREATE:
			corrections = append(corrections, &models.Correction{
				Msg: description,
				F: func() error {
					_, err := p.AppendRecords(context.Background(), dc.Name, libdnsRecords)
					return err
				},
			})
		case diff2.DELETE:
			var toDelete []libdns.Record
			for _, rc := range change.Old {
				if r, ok := rc.Original.(libdns.Record); ok {
					toDelete = append(toDelete, r)
				}
			}
			corrections = append(corrections, &models.Correction{
				Msg: description,
				F: func() error {
					_, err := p.DeleteRecords(context.Background(), dc.Name, toDelete)
					return err
				},
			})
		case diff2.CHANGE:
			corrections = append(corrections, &models.Correction{
				Msg: description,
				F: func() error {
					_, err := p.SetRecords(context.Background(), dc.Name, libdnsRecords)
					return err
				},
			})
		}
	}

	return corrections, actualChangeCount, nil
*/
}

func (n *spaceshipProvider) GetNameservers(domainName string) ([]*models.Nameserver, error) {
	client, err := client.NewClient(n.BaseURL, n.ApiKey, n.ApiSecret)
	if err != nil {
		return []*models.Nameserver{}, err
	}

	di, err := client.GetDomainInfo(context.Background(), domainName)
	if err != nil {
		return []*models.Nameserver{}, err
	}

	return models.ToNameservers(di.Nameservers.Hosts)
}
