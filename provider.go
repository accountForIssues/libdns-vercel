package vercel

import (
	"context"
	"strings"

	"github.com/libdns/libdns"
)

// Provider implements the libdns interfaces for Vercel
type Provider struct {
	// AuthAPIToken is the Vercel Authentication Token - see https://vercel.com/docs/api#api-basics/authentication
	AuthAPIToken string `json:"auth_api_token,omitempty"`
	// Optional, TeamId is the Vercel Team ID - see https://vercel.com/docs/rest-api#introduction/api-basics/authentication/accessing-resources-owned-by-a-team
	TeamId string `json:"team_id,omitempty"`
}

// GetRecords lists all the records in the zone.
func (p *Provider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	records, err := getAllRecords(ctx, p.AuthAPIToken, unFQDN(zone), p.TeamId)
	if err != nil {
		return nil, err
	}

	return records, nil
}

// AppendRecords adds records to the zone. It returns the records that were added.
func (p *Provider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	var appendedRecords []libdns.Record

	for _, record := range records {
		newRecord, err := createRecord(ctx, p.AuthAPIToken, unFQDN(zone), record, p.TeamId)
		if err != nil {
			return nil, err
		}
		appendedRecords = append(appendedRecords, newRecord)
	}

	return appendedRecords, nil
}

// DeleteRecords deletes the records from the zone.
func (p *Provider) DeleteRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	for _, record := range records {
		err := deleteRecord(ctx, unFQDN(zone), p.AuthAPIToken, record, p.TeamId)
		if err != nil {
			return nil, err
		}
	}

	return records, nil
}

// SetRecords sets the records in the zone, either by updating existing records
// or creating new ones. It returns the updated records.
func (p *Provider) SetRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	var setRecords []libdns.Record

	for _, record := range records {
		setRecord, err := createOrUpdateRecord(ctx, p.AuthAPIToken, unFQDN(zone), record, p.TeamId)
		if err != nil {
			return setRecords, err
		}
		setRecords = append(setRecords, setRecord)
	}

	return setRecords, nil
}

// unFQDN trims any trailing "." from fqdn. Vercel's API does not use FQDNs.
func unFQDN(fqdn string) string {
	return strings.TrimSuffix(fqdn, ".")
}

// Interface guards
var (
	_ libdns.RecordGetter   = (*Provider)(nil)
	_ libdns.RecordAppender = (*Provider)(nil)
	_ libdns.RecordSetter   = (*Provider)(nil)
	_ libdns.RecordDeleter  = (*Provider)(nil)
)
