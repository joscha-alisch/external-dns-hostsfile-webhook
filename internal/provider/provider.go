package provider

import (
	"context"
	"fmt"
	"github.com/joscha-alisch/external-dns-hostsfile-webhook/internal/hostsfile"
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/plan"
)

type Provider struct {
	h *hostsfile.HostsFile
}

func New(hostsFilePath string) *Provider {
	return &Provider{
		h: &hostsfile.HostsFile{
			Path: hostsFilePath,
		},
	}
}

func (p *Provider) Records(ctx context.Context) ([]*endpoint.Endpoint, error) {
	records, err := p.h.Get()
	if err != nil {
		return nil, fmt.Errorf("unable to get records")
	}

	var result []*endpoint.Endpoint
	for _, record := range records {
		result = append(result, endpoint.NewEndpoint(record.DNSName, record.Type, record.Ip))
	}

	return result, nil
}

func (p *Provider) ApplyChanges(ctx context.Context, changes *plan.Changes) error {
	if !changes.HasChanges() {
		return nil
	}

	for _, delOp := range changes.Delete {
		p.h.Remove(delOp.RecordType, delOp.DNSName)
	}

	for _, createOp := range changes.Create {
		p.h.Set(createOp.RecordType, createOp.DNSName, createOp.Targets[0])
	}

	for _, updateOp := range changes.UpdateNew {
		p.h.Set(updateOp.RecordType, updateOp.DNSName, updateOp.Targets[0])
	}

	p.h.Flush()
	return nil
}

func (p *Provider) AdjustEndpoints(endpoints []*endpoint.Endpoint) ([]*endpoint.Endpoint, error) {
	var result []*endpoint.Endpoint
	for _, e := range endpoints {
		switch e.RecordType {
		case "A":
			result = append(result, e)
		}
	}
	return result, nil
}

func (p *Provider) GetDomainFilter() endpoint.DomainFilter {
	return endpoint.DomainFilter{
		Filters: nil,
	}
}