package room

import (
	"caa-test/internal/entity"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repo {
	return &repo{
		db: db,
	}
}

func (r *repo) GetFirstUnassignedCustomerToday(ctx context.Context) (*entity.AgentAllocationQueue, error) {
	var queue entity.AgentAllocationQueue
	if r.db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	today := time.Now().Format("2006-01-02")

	err := r.db.
		WithContext(ctx).
		Where("agent_id IS NULL AND DATE(created_at) = ?", today).
		Order("created_at ASC").
		First(&queue).Error

	if err != nil {
		return nil, err
	}

	return &queue, nil
}

func (r *repo) AssignAgentToCustomer(queue *entity.AgentAllocationQueue) error {
	return r.db.Model(&entity.AgentAllocationQueue{}).
		Where("room_id = ?", queue.RoomID).
		Updates(entity.AgentAllocationQueue{
			AgentID: queue.AgentID,
		}).Error
}

func (r *repo) AddToQueue(queue *entity.AgentAllocationQueue) error {
	var existingQueue entity.AgentAllocationQueue
	if err := r.db.Where("room_id = ?", queue.RoomID).First(&existingQueue).Error; err == nil {
		return fmt.Errorf("room id already in queue")
	}

	// Added new queue when empty
	return r.db.Create(queue).Error
}

func (r *repo) GetRoomQueueByRoomId(roomId string) (*entity.AgentAllocationQueue, error) {
    var queue entity.AgentAllocationQueue
    
	err := r.db.
		Where("room_id = ?", roomId).
		First(&queue).Error

	// Tangani error dari query
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Kembalikan nil jika tidak ada record ditemukan
			return nil, nil
		}
		// Kembalikan error lainnya
		return nil, err
	}

	return &queue, nil
}
