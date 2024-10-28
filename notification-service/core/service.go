package core

import (
	"errors"
	"log"
	"notification-service/model"
	"time"

	"gorm.io/gorm"
)

type NotificationService interface {
	SaveNotifications(r NotificationRequest) (string, error)
	RemoveNotifications(r NotificationRequest) (string, error)
	GetMessages(uid uint) ([]MessageResponse, error)
	RemoveMessages(ids []uint, uid uint) (string, error)
	ReadAll(uid uint) (string, error)
}

type notificationService struct {
	db *gorm.DB
}

func NewNotificationService(db *gorm.DB) NotificationService {
	return &notificationService{db: db}
}

func (service *notificationService) SaveNotifications(r NotificationRequest) (string, error) {

	tx := service.db.Begin()
	if err := tx.Where("parent_id = ? AND uid = ? AND type = ?", r.ParentID, r.Uid, r.Type).Unscoped().Delete(&model.Notification{}).Error; err != nil {
		tx.Rollback()
		log.Println("db error")
		return "", errors.New("db error")
	}

	var notifications []model.Notification

	var startAt *time.Time
	var endAt *time.Time
	if r.StartAt != "" {
		start, _ := time.Parse("2006-01-02", r.StartAt)
		startAt = &start
	}
	if r.EndAt != "" {
		end, _ := time.Parse("2006-01-02", r.EndAt)
		endAt = &end
	}

	for _, v := range r.Times {
		notifications = append(notifications, model.Notification{Uid: r.Uid, ParentId: r.ParentID, Type: r.Type, Body: r.Body, Time: v, Weekdays: uintSliceToInt64Array(r.Week),
			StartAt: startAt, EndAt: endAt})
	}

	if err := tx.Create(&notifications).Error; err != nil {
		tx.Rollback()
		log.Println("db error2")
		return "", errors.New("db error2")
	}

	tx.Commit()
	return "200", nil
}

func (service *notificationService) RemoveNotifications(r NotificationRequest) (string, error) {

	if err := service.db.Where("parent_id = ? AND uid = ? AND type = ?", r.ParentID, r.Uid, r.Type).Unscoped().Delete(&model.Notification{}).Error; err != nil {
		log.Println("db error")
		return "", errors.New("db error")
	}
	return "200", nil
}

func (service *notificationService) GetMessages(uid uint) ([]MessageResponse, error) {
	var messages []model.Message
	result := service.db.Where("uid = ? ", uid).Order("id DESC").Find(&messages)
	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}

	var response []MessageResponse
	for _, v := range messages {
		response = append(response, MessageResponse{Id: v.ID, Type: v.Type, Body: v.Body, IsRead: v.IsRead, ParentId: v.ParentID, CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05")})
	}

	return response, nil
}

func (service *notificationService) ReadAll(uid uint) (string, error) {
	result := service.db.Model(&model.Message{}).Where("uid = ?", uid).Select("is_read").Updates(map[string]interface{}{"is_read": true})
	if result.Error != nil {
		log.Println("db error")
		return "", errors.New("db error")
	}
	return "200", nil
}

func (service *notificationService) RemoveMessages(ids []uint, uid uint) (string, error) {

	result := service.db.Where("id IN ? AND uid= ?", ids, uid).Delete(&model.Message{})

	if result.Error != nil {
		log.Println("db error")
		return "", errors.New("db error")
	}
	return "200", nil
}
