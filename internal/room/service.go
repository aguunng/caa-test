package room

import (
	"caa-test/internal/config"
	"caa-test/internal/entity"
	"caa-test/internal/qismo/request"
	"caa-test/internal/qismo/response"
	"context"
	"fmt"
	neturl "net/url"
	"strconv"

	"github.com/rs/zerolog/log"
)

type Omnichannel interface {
	Rooms(ctx context.Context, params neturl.Values) (*response.RoomsResponse, error)
	Agents(ctx context.Context) (*response.AgentsResponse, error)
	AssignAgent(ctx context.Context, params neturl.Values) error
	AgentByRoomID(ctx context.Context, roomID string) (*response.AgentsRoomResponse, error)
	AgentDetail(ctx context.Context, agentId string) (*response.AgentDetailResponse, error)
	AssignChannelToAgent(ctx context.Context, agentId string, channels request.AgentUpdatedRequest) error
}

type Service struct {
	omni Omnichannel
}

func NewService(omni Omnichannel) *Service {
	return &Service{
		omni: omni,
	}
}

func (s *Service) GetCustomerRoom(ctx context.Context) (*response.RoomsResponse, error) {
	l := log.Ctx(ctx).
		With().
		Str("func", "room.Service.GetCustomerRoom").
		Logger()

	result, err := s.omni.Rooms(ctx, neturl.Values{})
	if err != nil {
		l.Error().Msgf("unable get customer rooms : %s", err.Error())
		return nil, fmt.Errorf("failed to get customer name: %w", err)
	}

	return result, nil
}

func (s *Service) AssignAgent(ctx context.Context, request *request.WebhookCaaRequest) error {
	l := log.Ctx(ctx).
		With().
		Str("func", "room.Service.AssignAgent").
		Logger()

	params := neturl.Values{}

	firstRoom, err := s.FindFirstUnservedRoomId(ctx, neturl.Values{})
	if err != nil || firstRoom == nil {
		l.Error().Msgf("unable get fisrt customer room : %s", err.Error())
		return err
	}
	params.Set("room_id", firstRoom.ID)

	agentId, err := s.AvailableAgentIds(ctx, firstRoom)
	if err != nil {
		l.Error().Msgf("unable get availabale agent : %s", err.Error())
		return err
	}

	params.Set("agent_id", strconv.Itoa(agentId))

	err = s.omni.AssignAgent(ctx, params)
	if err != nil {
		l.Error().Msgf("unable assign agent : %s", err.Error())
		return err
	}

	return nil
}

func (s *Service) AvailableAgentIds(ctx context.Context, room *entity.Room) (int, error) {
	l := log.Ctx(ctx).
		With().
		Str("func", "room.Service.AvailableAgentIds").
		Logger()

	agentsInRoom, err := s.omni.AgentByRoomID(ctx, room.ID)
	if err != nil {
		l.Error().Msgf("unable geting agent data in room %s : %s", room.ID, err.Error())
		return 0, err
	}

	appConfig, err := config.ReadConfig()
	if err != nil {
		log.Ctx(ctx).Error().Msgf("failed read config: %s", err.Error())
		appConfig.MaxCustomer = 2
	}

	agentInRoomID := FilterAgentRoomAvailable(agentsInRoom, appConfig)

	if agentInRoomID > 0 {
		return agentInRoomID, nil
	}

	agents, err := s.omni.Agents(ctx)
	if err != nil {
		l.Error().Msgf("unable get agent data : %s", err.Error())
		return 0, err
	}

	agentId := FilterAgentsAvailableAssign(agents, appConfig)

	agentDetail, err := s.omni.AgentDetail(ctx, strconv.Itoa(agentId))
	if err != nil {
		l.Error().Msgf("unable ge agent detail : %s", err.Error())
		return 0, err
	}

	newChannel := response.UserChannel{
		ID:   room.ChannelID,
		Name: room.Source,
	}

	agentChannels := append(agentDetail.Data.Agent.UserChannels, newChannel)

	var mappedChannels []request.UserChannel
	for _, uc := range agentChannels {
		mappedChannels = append(mappedChannels, request.UserChannel{
			ChannelID: uc.ID,
			Source:    uc.Name,
		})
	}
	var userRoles []string
	for _, uc := range agentDetail.Data.Agent.UserRoles {
		userRoles = append(userRoles, strconv.Itoa(uc.ID))
	}

	if userRoles == nil {
		userRoles = []string{}
	}

	request := request.AgentUpdatedRequest{
		Channels:  mappedChannels,
		Name:      agentDetail.Data.Agent.Name,
		Email:     agentDetail.Data.Agent.Email,
		UserRoles: userRoles,
	}

	err = s.omni.AssignChannelToAgent(ctx, strconv.Itoa(agentDetail.Data.Agent.ID), request)

	if err != nil {
		l.Error().Msgf("unable assign agent %d to channels %d : %s", agentDetail.Data.Agent.ID, room.ChannelID, err.Error())
		return 0, err
	}

	return agentId, nil
}

func (s *Service) FindFirstUnservedRoomId(ctx context.Context, params neturl.Values) (*entity.Room, error) {
	l := log.Ctx(ctx).
		With().
		Str("func", "room.Service.FindFirstUnservedRoomId").
		Logger()

	params.Set("serve_status", "unserved")

	response, err := s.omni.Rooms(ctx, params)
	if err != nil {
		l.Error().Msgf("unable getting room data: %s", err.Error())
		return nil, err
	}

	if len(response.Data.CustomerRooms) >= 50 {
		params.Set("cursor_after", response.Meta.CursorAfter)
		return s.FindFirstUnservedRoomId(ctx, params)
	}

	for i := len(response.Data.CustomerRooms) - 1; i >= 0; i-- {
		room := response.Data.CustomerRooms[i]
		// Handle channel deleted source
		if room.Source != "uFppL" {
			return &entity.Room{
				ID:        room.RoomID,
				ChannelID: room.ChannelID,
				Source:    room.Source,
			}, nil
		}
	}

	return nil, nil
}

func FilterAgentRoomAvailable(agents *response.AgentsRoomResponse, appConfig *config.AppConfig) int {
	// Filter available agent and handle customer < 2
	var firstAgentId int
	for _, item := range agents.Data.Agents {
		if item.IsAvailable && item.CurrentCustomerCount < appConfig.MaxCustomer {
			firstAgentId = item.ID
		}
	}

	return firstAgentId
}
func FilterAgentsAvailableAssign(agents *response.AgentsResponse, appConfig *config.AppConfig) int {
	// Filter available agent and handle customer < 2 and not supervisor
	var firstAgentId int
	for _, item := range agents.Data.Agents.Data {
		if item.IsAvailable && item.CurrentCustomerCount < appConfig.MaxCustomer && !item.IsSupervisor {
			firstAgentId = item.ID
		}
	}

	return firstAgentId
}
