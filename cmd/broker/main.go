package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/alo/kafscale/pkg/broker"
	"github.com/alo/kafscale/pkg/protocol"
)

type handler struct {
	apiVersions []protocol.ApiVersion
	cluster     metadataState
}

type metadataState struct {
	brokers      []protocol.MetadataBroker
	controllerID int32
	topics       []protocol.MetadataTopic
	clusterID    *string
}

func (h *handler) Handle(header *protocol.RequestHeader, req protocol.Request) ([]byte, error) {
	switch req.(type) {
	case *protocol.ApiVersionsRequest:
		return protocol.EncodeApiVersionsResponse(&protocol.ApiVersionsResponse{
			CorrelationID: header.CorrelationID,
			ErrorCode:     0,
			Versions:      h.apiVersions,
		})
	case *protocol.MetadataRequest:
		return protocol.EncodeMetadataResponse(&protocol.MetadataResponse{
			CorrelationID: header.CorrelationID,
			Brokers:       h.cluster.brokers,
			ClusterID:     h.cluster.clusterID,
			ControllerID:  h.cluster.controllerID,
			Topics:        h.cluster.topics,
		})
	default:
		return nil, ErrUnsupportedAPI
	}
}

var ErrUnsupportedAPI = fmt.Errorf("unsupported api")

func defaultHandler() *handler {
	clusterID := "kafscale-cluster"
	return &handler{
		apiVersions: []protocol.ApiVersion{
			{APIKey: protocol.APIKeyApiVersion, MinVersion: 0, MaxVersion: 0},
			{APIKey: protocol.APIKeyMetadata, MinVersion: 0, MaxVersion: 0},
		},
		cluster: metadataState{
			controllerID: 1,
			clusterID:    &clusterID,
			brokers: []protocol.MetadataBroker{
				{NodeID: 1, Host: "localhost", Port: 19092},
			},
			topics: []protocol.MetadataTopic{
				{
					ErrorCode: 0,
					Name:      "orders",
					Partitions: []protocol.MetadataPartition{
						{
							ErrorCode:      0,
							PartitionIndex: 0,
							LeaderID:       1,
							ReplicaNodes:   []int32{1},
							ISRNodes:       []int32{1},
						},
					},
				},
			},
		},
	}
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	srv := &broker.Server{
		Addr:    ":19092",
		Handler: defaultHandler(),
	}
	if err := srv.ListenAndServe(ctx); err != nil {
		log.Fatalf("broker server error: %v", err)
	}
	srv.Wait()
}
