package room

import (
	"caa-test/internal/config"
	"caa-test/internal/entity"
	"caa-test/internal/qismo/request"
	"caa-test/internal/qismo/response"
	"context"
	"fmt"
	neturl "net/url"
	"sort"
	"strconv"

	"github.com/rs/zerolog/log"
)

type Omnichannel interface {
	Rooms(ctx context.Context, params neturl.Values) (*response.RoomsResponse, error)
	Agents(ctx context.Context) (*response.AgentsResponse, error)
	AssignAgent(ctx context.Context, params neturl.Values) error
	SearchCandidateAgent(ctx context.Context, roomID string) (*response.AgentsRoomResponse, error)
	AgentDetail(ctx context.Context, agentId string) (*response.AgentDetailResponse, error)
	AssignChannelToAgent(ctx context.Context, agentId string, channels request.AgentUpdatedRequest) error
}

type Repository interface {
	GetFirstUnassignedCustomerToday(ctx context.Context) (*entity.AgentAllocationQueue, error)
	AssignAgentToCustomer(queue *entity.AgentAllocationQueue) error
	GetRoomQueueByRoomId(roomId string) (*entity.AgentAllocationQueue, error)
	AddToQueue(queue *entity.AgentAllocationQueue) error
}

type Service struct {
	repo Repository
	omni Omnichannel
	cfg  *config.Config
}

func NewService(repo Repository, omni Omnichannel, cfg *config.Config) *Service {
	return &Service{
		repo: repo,
		omni: omni,
		cfg:  cfg,
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

func (s *Service) AssignAgentFromCaa(ctx context.Context, request *request.WebhookCaaRequest) error {
	l := log.Ctx(ctx).
		With().
		Str("func", "room.Service.AssignAgentFromCaa").
		Logger()

	var roomId string

	if request.AppID != s.cfg.Qiscus.AppID {
		return &roomError{500, "failed assign with wrong app id"}
	}

	err := s.AddToQueue(request.RoomID)
	if err != nil {
		l.Error().Msgf("failed saved customer to queue : %s", err.Error())
	}

	roomQueueDetail, _ := s.repo.GetRoomQueueByRoomId(request.RoomID)

	if len(roomQueueDetail.AgentID) > 0 {
		return &roomError{500, "current room id has been assigned to agent"}
	}

	unAssignedQueue, err := s.repo.GetFirstUnassignedCustomerToday(ctx)

	// Check if queue is empty to set with room id from webhook
	if err == nil {
		roomId = unAssignedQueue.RoomID
	} else {
		roomId = request.RoomID
	}

	agentId := s.AvailableAgentId(ctx, roomId)

	if agentId == 0 {
		return fmt.Errorf("failed get available agent to assign with room id : %s", request.RoomID)
	}

	err = s.AssignAgent(ctx, strconv.Itoa(agentId), roomId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) AssignAgent(ctx context.Context, agentId string, roomId string) error {
	params := neturl.Values{}

	params.Set("agent_id", agentId)
	params.Set("room_id", roomId)

	err := s.omni.AssignAgent(ctx, params)
	if err != nil {
		return err
	}

	queueUpdate := &entity.AgentAllocationQueue{
		RoomID:  roomId,
		AgentID: agentId,
	}

	err = s.repo.AssignAgentToCustomer(queueUpdate)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) AssignAgentFromResolved(ctx context.Context) error {
	unAssignedQueue, _ := s.repo.GetFirstUnassignedCustomerToday(ctx)

	if unAssignedQueue == nil {
		return &roomError{500, "empty queue customer for assign to agent"}
	}

	agentId := s.AvailableAgentId(ctx, unAssignedQueue.AgentID)

	if agentId == 0 {
		return fmt.Errorf("failed get available agent to assign with room id : %s", unAssignedQueue.RoomID)
	}

	err := s.AssignAgent(ctx, strconv.Itoa(agentId), unAssignedQueue.RoomID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) AvailableAgentId(ctx context.Context, roomId string) int {
	l := log.Ctx(ctx).
		With().
		Str("func", "room.Service.AvailableAgentIds").
		Logger()

	candidateAgents, err := s.omni.SearchCandidateAgent(ctx, roomId)
	if err != nil {
		l.Error().Msgf("unable get candidate agent in room %s : %s", roomId, err.Error())
		return 0
	}

	appConfig, err := config.ReadConfig()
	if err != nil {
		log.Ctx(ctx).Error().Msgf("failed read config: %s", err.Error())
		appConfig.MaxCustomer = 2
	}

	agentId := FilterAgentRoomAvailable(candidateAgents, appConfig)

	if agentId == 0 {
		log.Ctx(ctx).Error().Msgf("failed get available agent with config max customer : %d", appConfig.MaxCustomer)
	}

	return agentId
}

func (s *Service) AddToQueue(roomId string) error {
	queue := &entity.AgentAllocationQueue{
		RoomID: roomId,
	}
	err := s.repo.AddToQueue(queue)

	return err
}

func FilterAgentRoomAvailable(agents *response.AgentsRoomResponse, appConfig *config.AppConfig) int {
	// Filter available agent and handle customer < 2
	var filteredAgents []struct {
		ID                   int
		CurrentCustomerCount int
	}

	for _, item := range agents.Data.Agents {
		if item.IsAvailable && item.CurrentCustomerCount < appConfig.MaxCustomer {
			filteredAgents = append(filteredAgents, struct {
				ID                   int
				CurrentCustomerCount int
			}{
				ID:                   item.ID,
				CurrentCustomerCount: item.CurrentCustomerCount,
			})
		}
	}

	// Sort filtered agents by CurrentCustomerCount in ascending order
	sort.Slice(filteredAgents, func(i, j int) bool {
		return filteredAgents[i].CurrentCustomerCount < filteredAgents[j].CurrentCustomerCount
	})

	// Return the ID of the first agent in the sorted list
	if len(filteredAgents) > 0 {
		return filteredAgents[0].ID
	}

	return 0
}
