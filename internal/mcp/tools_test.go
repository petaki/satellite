package mcp

import (
	"context"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/petaki/satellite/internal/models"
)

type stubProbeRepo struct {
	probes  []models.Probe
	deleted []models.Probe
}

func (s *stubProbeRepo) FindAll() ([]models.Probe, error) { return s.probes, nil }
func (s *stubProbeRepo) FindLatestValues(models.Probe, int) ([]any, *time.Time, error) {
	return nil, nil, nil
}
func (s *stubProbeRepo) HasHeartbeat(models.Probe) (bool, error) { return false, nil }
func (s *stubProbeRepo) SetHeartbeat(models.Probe, int) error    { return nil }
func (s *stubProbeRepo) Delete(p models.Probe) error {
	s.deleted = append(s.deleted, p)
	return nil
}

func req(probe string) mcp.CallToolRequest {
	var r mcp.CallToolRequest
	r.Params.Name = "delete_probe"
	r.Params.Arguments = map[string]any{"probe": probe}
	return r
}

func TestDeleteProbeRejectsNamesThatAreNotProbes(t *testing.T) {
	for _, name := range []string{"*", "*:*", "web-0?", "[a-z]*", "nope"} {
		repo := &stubProbeRepo{probes: []models.Probe{"web-01", "db-01"}}
		h := &handler{probeRepository: repo}

		res, err := h.deleteProbe(context.Background(), req(name))
		if err != nil {
			t.Fatalf("%q: unexpected err: %v", name, err)
		}
		if !res.IsError {
			t.Errorf("%q: expected an error result, got success", name)
		}
		if len(repo.deleted) != 0 {
			t.Errorf("%q: Delete reached the repository with %v", name, repo.deleted)
		}
	}
}

func TestDeleteProbeAllowsKnownName(t *testing.T) {
	repo := &stubProbeRepo{probes: []models.Probe{"web-01", "db-01"}}
	h := &handler{probeRepository: repo}

	res, err := h.deleteProbe(context.Background(), req("web-01"))
	if err != nil {
		t.Fatal(err)
	}
	if res.IsError {
		t.Fatal("expected success for a real probe")
	}
	if len(repo.deleted) != 1 || repo.deleted[0] != "web-01" {
		t.Fatalf("expected delete of web-01, got %v", repo.deleted)
	}
}
