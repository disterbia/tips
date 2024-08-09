package core

import (
	"context"
	"fcm-service/model"
	"log"
	"strconv"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/robfig/cron/v3"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

var firebaseClient *messaging.Client

func StartCentralCronScheduler(db *gorm.DB) {
	initializeFirebase()
	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc("0 * * * * *", func() {
		sendPendingNotifications(db)
	})
	if err != nil {
		log.Fatalf("Failed to create cron job: %v", err)
	}
	c.Start()
}
func initializeFirebase() {
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile("./firebase-adminkey.json"))
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
		return
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
		return
	}
	firebaseClient = client
}

func sendPendingNotifications(db *gorm.DB) {
	now := time.Now()
	var notifications []model.Notification
	db.Where("(start_at IS NULL OR start_at <= ?) AND (end_at IS NULL OR end_at >= ?)", now, now).Find(&notifications)

	for _, notification := range notifications {
		if shouldSendNotification(now, notification) {
			go sendMedicationReminder(context.Background(), notification, db)
		}
	}
}

func shouldSendNotification(now time.Time, notification model.Notification) bool {
	currentWeekday := int(now.Weekday())

	alarmTime, _ := time.Parse("15:04", notification.Time)

	// 알림이 설정된 요일 배열과 현재 요일을 비교합니다.
	for _, weekday := range notification.Weekdays {
		if int(weekday) == currentWeekday && now.Format("15:04") == alarmTime.Format("15:04") {
			return true
		}
	}

	return false
}

func sendMedicationReminder(ctx context.Context, notification model.Notification, db *gorm.DB) {
	var user model.User
	if err := db.First(&user, "id = ?", notification.Uid).Error; err != nil {
		log.Printf("Failed to find user: %v\n", err)
		return
	}

	notification_count := calculateNotificationCount(db, notification.Uid)
	log.Println(user.FCMToken)
	var title string
	switch notification.Type {
	case uint(MEDICINENOTIFICATION):
		title = "약물 복용"
	case uint(EXERCISENOTIFICATION):
		title = "운동 시간"
	}
	fcm := &messaging.Message{
		Data: map[string]string{
			"uid":                strconv.FormatUint(uint64(notification.Uid), 10),
			"type":               strconv.FormatUint(uint64(notification.Type), 10),
			"notification_count": strconv.FormatUint(uint64(notification_count), 10),
			"time":               notification.CreatedAt.Format("2006-01-02"),
			"parent_id":          strconv.FormatUint(uint64(notification.ParentId), 10),
		},
		Notification: &messaging.Notification{
			Title: title,
			Body:  notification.Body,
		},
		Token: user.FCMToken,
	}

	response, err := firebaseClient.Send(ctx, fcm)
	if err != nil {
		log.Printf("error sending message: %v\n", err)
		return
	}
	log.Printf("Successfully sent message: %s\n", response)

	message := model.Message{
		Uid:      notification.Uid,
		Type:     notification.Type,
		Body:     notification.Body,
		ParentID: notification.ParentId,
		IsRead:   false,
	}
	if err := db.Create(&message).Error; err != nil {
		log.Printf("error create message: %v\n", err)
	}

}

func calculateNotificationCount(db *gorm.DB, uid uint) uint {
	var count int64
	db.Model(&model.Message{}).Where("uid = ? AND is_read = ?", uid, false).Count(&count)
	return uint(count)
}
