package rpc

import (
	"context"

	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum-optimism/optimism/op-service/eth"
)

type conductor interface {
	Pause(ctx context.Context) error
	Resume(ctx context.Context) error
	Paused() bool
	Stopped() bool

	Leader(ctx context.Context) bool
	LeaderWithID(ctx context.Context) (string, string)
	AddServerAsVoter(ctx context.Context, id string, addr string) error
	AddServerAsNonvoter(ctx context.Context, id string, addr string) error
	RemoveServer(ctx context.Context, id string) error
	TransferLeader(ctx context.Context) error
	TransferLeaderToServer(ctx context.Context, id string, addr string) error
	CommitUnsafePayload(ctx context.Context, payload *eth.ExecutionPayload) error
}

// APIBackend is the backend implementation of the API.
// TODO: (https://github.com/ethereum-optimism/protocol-quest/issues/45) Add metrics tracer here.
// TODO: (https://github.com/ethereum-optimism/protocol-quest/issues/44) add tests after e2e setup.
type APIBackend struct {
	log log.Logger
	con conductor
}

// NewAPIBackend creates a new APIBackend instance.
func NewAPIBackend(log log.Logger, con conductor) *APIBackend {
	return &APIBackend{
		log: log,
		con: con,
	}
}

var _ API = (*APIBackend)(nil)

// Active implements API.
func (api *APIBackend) Active(_ context.Context) (bool, error) {
	return !api.con.Stopped() && !api.con.Paused(), nil
}

// AddServerAsNonvoter implements API.
func (api *APIBackend) AddServerAsNonvoter(ctx context.Context, id string, addr string) error {
	return api.con.AddServerAsNonvoter(ctx, id, addr)
}

// AddServerAsVoter implements API.
func (api *APIBackend) AddServerAsVoter(ctx context.Context, id string, addr string) error {
	return api.con.AddServerAsVoter(ctx, id, addr)
}

// CommitUnsafePayload implements API.
func (api *APIBackend) CommitUnsafePayload(ctx context.Context, payload *eth.ExecutionPayload) error {
	return api.con.CommitUnsafePayload(ctx, payload)
}

// Leader implements API, returns true if current conductor is leader of the cluster.
func (api *APIBackend) Leader(ctx context.Context) (bool, error) {
	return api.con.Leader(ctx), nil
}

// LeaderWithID implements API, returns the leader's server ID and address (not necessarily the current conductor).
func (api *APIBackend) LeaderWithID(ctx context.Context) (*ServerInfo, error) {
	id, addr := api.con.LeaderWithID(ctx)
	return &ServerInfo{
		ID:   id,
		Addr: addr,
	}, nil
}

// Pause implements API.
func (api *APIBackend) Pause(ctx context.Context) error {
	return api.con.Pause(ctx)
}

// RemoveServer implements API.
func (api *APIBackend) RemoveServer(ctx context.Context, id string) error {
	return api.con.RemoveServer(ctx, id)
}

// Resume implements API.
func (api *APIBackend) Resume(ctx context.Context) error {
	return api.con.Resume(ctx)
}

// TransferLeader implements API.
func (api *APIBackend) TransferLeader(ctx context.Context) error {
	return api.con.TransferLeader(ctx)
}

// TransferLeaderToServer implements API.
func (api *APIBackend) TransferLeaderToServer(ctx context.Context, id string, addr string) error {
	return api.con.TransferLeaderToServer(ctx, id, addr)
}
