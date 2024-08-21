package room

import (
	"caa-test/internal/api/resp"
	"caa-test/internal/config"
	"caa-test/internal/entity"
	"caa-test/internal/qismo/request"
	"encoding/json"
	"fmt"
	"net/http"
	neturl "net/url"
	"reflect"

	"github.com/go-playground/validator/v10"
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

// @Summary Webhook CAA
// @Description Handles the CAA webhook request with data received from the request body.
// @Tags webhook
// @Accept  json
// @Produce  json
// @Param   data  body   request.WebhookCaaRequest  true  "Request body for CAA webhook"
// @Success 200  {object}  resp.HTTPSuccess  "Successfully caa webhook"
// @Failure 400  {object}  resp.HTTPError  "Bad Request"
// @Failure 500  {object}  resp.HTTPError  "Internal Server Error"
// @Router /caa [post]
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

	resp.WriteJSON(w, http.StatusOK, resp.HTTPSuccess{
		Message: "Successfully send webhook",
	})
}

// @Summary Mark Webhook as Resolved
// @Description Marks the webhook as resolved by processing the data provided in the request body and assigning an agent.
// @Tags webhook
// @Accept  json
// @Produce  json
// @Param   data  body   request.WebhookMarkResolvedRequest  true  "Request body to mark the webhook as resolved"
// @Success 200  {object}  resp.HTTPSuccess  "Successfully marked the webhook as resolved"
// @Failure 400  {object}  resp.HTTPError  "Bad Request"
// @Failure 500  {object}  resp.HTTPError  "Internal Server Error"
// @Router /mark_as_resolved [post]
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

	resp.WriteJSON(w, http.StatusOK, resp.HTTPSuccess{
		Message: "Successfully send webhook",
	})
}

// @Summary Get Customer Rooms
// @Description Retrieves the list of available customer rooms.
// @Tags rooms
// @Produce  json
// @Success 200  {object}  response.RoomsResponse  "Successfully retrieved list of rooms"
// @Failure 500  {object}  resp.HTTPError  "Internal Server Error"
// @Router /rooms [get]
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

// @Summary Get First Unserved Room ID
// @Description Retrieves the ID of the first unserved customer room.
// @Tags rooms
// @Produce  json
// @Success 200  {object}  resp.HTTPSuccess  "Successfully retrieved the first unserved room ID"
// @Failure 500  {object}  resp.HTTPError  "Internal Server Error"
// @Router /first_room [get]
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

// @Summary Update Maximum Customer Limit
// @Description Updates the maximum number of customers allowed in the configuration.
// @Tags configuration
// @Accept  json
// @Produce  json
// @Param   data  body   entity.Config  true  "Payload to update max customer"
// @Success 200  {object}  resp.HTTPSuccess  "Successfully updated max customer"
// @Failure 400  {object}  resp.HTTPError  "Bad Request"
// @Failure 500  {object}  resp.HTTPError  "Internal Server Error"
// @Router /update-max-customer [post]
func (h *httpHandler) UpdateMaxCustomerHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req entity.Config

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.WriteJSON(w, http.StatusBadRequest, resp.HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		errorResponse := resp.ValidationErrorResponse{
			StatusCode: http.StatusBadRequest,
			Errors:     formatValidationError(err, req),
		}
		resp.WriteJSON(w, http.StatusBadRequest, errorResponse)
		return
	}

	appConfig, err := config.ReadConfig()
	if err != nil {
		log.Ctx(ctx).Error().Msgf("failed read config: %s", err.Error())
		resp.WriteJSONFromError(w, err)
		return
	}

	appConfig.MaxCustomer = req.MaxCustomer

	if err := config.WriteConfig(appConfig); err != nil {
		log.Error().Msgf("failed to write config: %s", err.Error())
		http.Error(w, "failed to write config", http.StatusInternalServerError)
		return
	}

	resp.WriteJSON(w, http.StatusOK, resp.HTTPSuccess{
		Message: "Succesfully update max customer",
	})
}

func formatValidationError(err error, req interface{}) map[string][]string {
	validationErrors := make(map[string][]string)
	typ := reflect.TypeOf(req)

	for _, e := range err.(validator.ValidationErrors) {
		fieldName := e.Field()
		jsonTag := getJSONTag(fieldName, typ)
		var errorMsg string

		switch e.Tag() {
		case "required":
			errorMsg = fmt.Sprintf("%s field is required.", jsonTag)
		case "min":
			errorMsg = fmt.Sprintf("%s field must be at least %s.", jsonTag, e.Param())
		case "max":
			errorMsg = fmt.Sprintf("%s field cannot be greater than %s.", jsonTag, e.Param())
		case "email":
			errorMsg = fmt.Sprintf("%s field must be a valid email address.", jsonTag)
		case "len":
			errorMsg = fmt.Sprintf("%s field must be exactly %s characters long.", jsonTag, e.Param())
		case "lte":
			errorMsg = fmt.Sprintf("%s field must be less than or equal to %s.", jsonTag, e.Param())
		case "gte":
			errorMsg = fmt.Sprintf("%s field must be greater than or equal to %s.", jsonTag, e.Param())
		case "oneof":
			errorMsg = fmt.Sprintf("%s field must be one of the following values: %s.", jsonTag, e.Param())
		default:
			errorMsg = fmt.Sprintf("%s field failed validation for rule '%s'.", jsonTag, e.Tag())
		}

		validationErrors[fieldName] = append(validationErrors[fieldName], errorMsg)
	}
	return validationErrors
}

func getJSONTag(fieldName string, typ reflect.Type) string {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Name == fieldName {
			jsonTag := field.Tag.Get("json")
			if jsonTag != "" {
				return jsonTag
			}
		}
	}
	return fieldName
}
