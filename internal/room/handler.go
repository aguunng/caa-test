package room

import (
	"caa-test/internal/api/resp"
	"caa-test/internal/config"
	"caa-test/internal/qismo/request"
	"encoding/json"
	"net/http"
	neturl "net/url"

	"github.com/rs/zerolog/log"
)

type httpHandler struct {
	svc *Service
}

func NewHttpHandler(svc *Service) *httpHandler {
	return &httpHandler{
		svc: svc,
	}
}

func (h *httpHandler) WebhookCaa(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req request.WebhookCaaRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		resp.WriteJSON(w, http.StatusBadRequest, resp.HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})

		return
	}

	if err := h.svc.AssignAgent(ctx, &req); err != nil {
		log.Ctx(ctx).Error().Msgf("failed assign agent to room : %s", err.Error())
		resp.WriteJSONFromError(w, err)
		return
	}

	resp.WriteJSON(w, http.StatusOK, "ok")
}

func (h *httpHandler) WebhookMarkResolved(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req request.WebhookMarkResolvedRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		resp.WriteJSON(w, http.StatusBadRequest, resp.HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})

		return
	}

	if err := h.svc.AssignAgent(ctx, nil); err != nil {
		log.Ctx(ctx).Error().Msgf("failed assign agent to room : %s", err.Error())
		resp.WriteJSONFromError(w, err)
		return
	}

	resp.WriteJSON(w, http.StatusOK, "ok")
}

func (h *httpHandler) GetRooms(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	room, err := h.svc.GetCustomerRoom(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("failed to get room: %s", err.Error())
		resp.WriteJSONFromError(w, err)
		return
	}

	resp.WriteJSON(w, http.StatusOK, room)
}

func (h *httpHandler) FirstCustomerRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	roomId, err := h.svc.FindFirstUnservedRoomId(ctx, neturl.Values{})
	if err != nil {
		log.Ctx(ctx).Error().Msgf("failed to get room: %s", err.Error())
		resp.WriteJSONFromError(w, err)
		return
	}

	resp.WriteJSON(w, http.StatusOK, resp.HTTPSuccess{
		Message: "Succesfully get data first room",
		Data: map[string]interface{}{
			"room_id": roomId,
		},
	})
}

func (h *httpHandler) UpdateMaxCustomerHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var requestPayload struct {
		MaxCustomer int `json:"max_customer"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	appConfig, err := config.ReadConfig()
	if err != nil {
		log.Ctx(ctx).Error().Msgf("failed read config: %s", err.Error())
		resp.WriteJSONFromError(w, err)
		return
	}

	appConfig.MaxCustomer = requestPayload.MaxCustomer

	if err := config.WriteConfig(appConfig); err != nil {
		log.Error().Msgf("failed to write config: %s", err.Error())
		http.Error(w, "failed to write config", http.StatusInternalServerError)
		return
	}

	resp.WriteJSON(w, http.StatusOK, resp.HTTPSuccess{
		Message: "Succesfully update max customer",
	})
}
