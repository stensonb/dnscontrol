package spaceship

import (
	"context"
	"fmt"
	"time"

	"github.com/DNSControl/dnscontrol/v4/models"
	"github.com/DNSControl/dnscontrol/v4/pkg/diff2"
	"github.com/libdns/libdns"
	spaceship "github.com/libdns/spaceship"
)

func (n *spaceshipProvider) GetZoneRecords(dc *models.DomainConfig) (models.Records, error) {
	p := &spaceship.Provider{
		APIKey:    n.ApiKey,
		APISecret: n.ApiSecret,
	}

	records, err := p.GetRecords(context.Background(), dc.Name)
	if err != nil {
		return nil, fmt.Errorf("spaceship API error: %w", err)
	}

	existingRecords := make([]*models.RecordConfig, 0, len(records))
	for _, r := range records {
		rc := &models.RecordConfig{
			Type:     r.RR().Type,
			TTL:      uint32(r.RR().TTL.Seconds()),
			Original: r,
		}

		rc.SetLabel(r.RR().Name, dc.Name)

		err := rc.PopulateFromString(r.RR().Type, r.RR().Data, dc.Name)
		if err != nil {
			return nil, fmt.Errorf("error parsing record from spaceship: %w", err)
		}

		existingRecords = append(existingRecords, rc)
	}

	return existingRecords, nil
}

func (n *spaceshipProvider) GetZoneRecordsCorrections(dc *models.DomainConfig, existing models.Records) ([]*models.Correction, int, error) {
	changes, actualChangeCount, err := diff2.ByRecord(existing, dc, nil)
	if err != nil {
		return nil, 0, err
	}

	p := &spaceship.Provider{
		APIKey:    n.ApiKey,
		APISecret: n.ApiSecret,
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
}

func (n *spaceshipProvider) GetNameservers(domainName string) ([]*models.Nameserver, error) {
	return models.ToNameservers([]string{
		"ns1.spaceship.com",
		"ns2.spaceship.com",
	})
}
