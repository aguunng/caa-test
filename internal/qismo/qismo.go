package qismo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	neturl "net/url"

	"caa-test/internal/qismo/request"
	"caa-test/internal/qismo/response"
)

type Qismo struct {
	client    httpClient
	url       string
	appID     string
	secretKey string
}

func New(client httpClient, url, appID, secretKey string) *Qismo {
	return &Qismo{
		client:    client,
		url:       url,
		appID:     appID,
		secretKey: secretKey,
	}
}

func (q *Qismo) headers() map[string]string {
	return map[string]string{
		"Qiscus-App-Id":     q.appID,
		"Qiscus-Secret-Key": q.secretKey,
	}
}

func (q *Qismo) AssignAgent(ctx context.Context, params neturl.Values) error {
	url := fmt.Sprintf("%s/api/v1/admin/service/assign_agent?%s", q.url, params.Encode())
	payload, _ := json.Marshal(map[string]interface{}{})

	var response response.AgentsResponse

	err := q.client.Call(ctx, http.MethodPost, url, bytes.NewBuffer(payload), q.headers(), &response)
	if err != nil {
		return err
	}

	return nil
}

func (q *Qismo) Agents(ctx context.Context) (*response.AgentsResponse, error) {
	params := neturl.Values{}
	params.Set("limit", "100")

	url := fmt.Sprintf("%s/api/v1/admin/agents?%s", q.url, params.Encode())

	var response response.AgentsResponse

	err := q.client.Call(ctx, http.MethodGet, url, nil, q.headers(), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (q *Qismo) Rooms(ctx context.Context, params neturl.Values) (*response.RoomsResponse, error) {
	url := fmt.Sprintf("%s/api/v2/customer_rooms?%s", q.url, params.Encode())

	var response response.RoomsResponse

	err := q.client.Call(ctx, http.MethodGet, url, nil, q.headers(), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (q *Qismo) AgentByRoomID(ctx context.Context, roomId string) (*response.AgentsRoomResponse, error) {
	params := neturl.Values{}
	params.Set("room_id", roomId)

	url := fmt.Sprintf("%s/api/v2/admin/service/available_agents?%s", q.url, params.Encode())
	var response response.AgentsRoomResponse

	err := q.client.Call(ctx, http.MethodGet, url, nil, q.headers(), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (q *Qismo) AgentDetail(ctx context.Context, agentId string) (*response.AgentDetailResponse, error) {
	url := fmt.Sprintf("%s/api/v2/admin/agent/%s", q.url, agentId)
	var response response.AgentDetailResponse

	err := q.client.Call(ctx, http.MethodGet, url, nil, q.headers(), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (q *Qismo) AssignChannelToAgent(ctx context.Context, agentId string, request request.AgentUpdatedRequest) error {
	url := fmt.Sprintf("%s/api/v2/admin/agent/%s/update", q.url, agentId)
	payload, _ := json.Marshal(request)

	var response response.AgentDetailResponse

	err := q.client.Call(ctx, http.MethodPost, url, bytes.NewBuffer(payload), q.headers(), &response)
	if err != nil {
		return err
	}

	return nil
}
