package core

import (
	"context"
	"fcm-service/model"
	"log"
	"strconv"
	"sync"
	"sync/atomic"
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
	// 오전 10시마다 실행되는 스케줄러 추가
	_, err = c.AddFunc("0 0 10 * * *", func() {
		sendDailyNotifications(db, EMOTIONNOTIFICATION)
	})
	if err != nil {
		log.Fatalf("Failed to create 10AM cron job: %v", err)
	}

	// 오후 8시마다 실행되는 스케줄러 추가
	_, err = c.AddFunc("0 0 20 * * *", func() {
		sendDailyNotifications(db, EMOTIONNOTIFICATION)
		sendDailyNotifications(db, CHECKNOTIFICATION)
	})
	if err != nil {
		log.Fatalf("Failed to create 8PM cron job: %v", err)
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
	currentDate := now.Format("2006-01-02") // yyyy-mm-dd 형식의 날짜
	currentWeekday := int(now.Weekday())    // 현재 요일 (0 = Sunday, ..., 6 = Saturday)
	currentTime := now.Format("15:04")      // 현재 시간 (hh:mm 형식)

	var notifications []model.Notification

	// 필요한 알림만 조회 (StartAt, EndAt, 요일, 시간 조건 추가)
	if err := db.Preload("User").Debug().
		Where("(start_at IS NULL OR start_at <= ?) AND (end_at IS NULL OR end_at >= ?)", currentDate, currentDate).
		Where("? = ANY(weekdays)", currentWeekday).
		Where("time = ?", currentTime).
		Find(&notifications).Error; err != nil {
		log.Printf("Error fetching notifications: %v", err)
		return
	}

	// 최신 유저만 선택
	uniqueUsers, err := getLatestUsersByFCMToken(db)
	if err != nil {
		log.Printf("Error fetching unique users: %v\n", err)
		return
	}

	// 유저별로 FCMToken -> Uid 맵 생성
	latestUserMap := make(map[string]uint)
	for _, user := range uniqueUsers {
		latestUserMap[user.FcmToken] = user.Uid
	}

	// 2. 사용자별 미확인 알림 카운트를 한 번에 조회
	notificationCounts := getNotificationCounts(db)

	// 3. 메시지와 저장할 알림을 준비
	var messages []*messaging.Message
	var newMessages []model.Message

	for _, n := range notifications {
		// 중복된 FCM 토큰 필터링: 최신 유저만 메시지 전송
		if latestUserMap[n.User.FcmToken] != n.Uid {
			continue // 중복된 경우 스킵
		}
		notificationCount := notificationCounts[n.Uid]
		log.Println("aa:", n.User.FcmToken)
		// FCM 메시지 생성
		fcmMessage := &messaging.Message{
			Data: map[string]string{
				"uid":                strconv.FormatUint(uint64(n.Uid), 10),
				"type":               strconv.FormatUint(uint64(n.Type), 10),
				"notification_count": strconv.FormatUint(uint64(notificationCount), 10),
				"time":               currentDate,
				"parent_id":          strconv.FormatUint(uint64(n.ParentId), 10),
			},
			Notification: &messaging.Notification{
				Title: getTitle(n.Type),
				Body:  n.Body,
			},
			Token: n.User.FcmToken, // 조회한 FCMToken 사용
		}

		messages = append(messages, fcmMessage)

		// 성공 시 DB에 저장할 메시지 준비
		newMessages = append(newMessages, model.Message{
			Uid:      n.Uid,
			Type:     n.Type,
			Body:     n.Body,
			ParentID: n.ParentId,
			IsRead:   false,
		})

	}
	sendFCMMessages(db, messages, newMessages)
}

func sendDailyNotifications(db *gorm.DB, notificationType uint) {
	var title, body string
	parentID := 0
	today := time.Now().Format("2006-01-02")

	// 1. 알림 메시지 설정
	switch notificationType {
	case uint(EMOTIONNOTIFICATION):
		title = "기분 등록"
		body = "오늘 기분이 등록되지 않았습니다.\n기분을 등록해보세요."
	case uint(CHECKNOTIFICATION):
		title = "검사 시간"
		body = "오늘 운동 테스트 하셨나요?\n운동 등록과 함께 해보세요."
	}

	// 2. 오늘의 데이터가 없는 유저 조회
	var usersToNotify []struct {
		Uid      uint
		FcmToken string
	}

	if notificationType == uint(EMOTIONNOTIFICATION) {
		if err := db.Raw(`
			SELECT u.id AS uid, u.fcm_token as fcm_token
			FROM users u
			WHERE u.id NOT IN (
				SELECT e.uid FROM emotions e WHERE e.target_date = ?
			)
		`, today).Scan(&usersToNotify).Error; err != nil {
			return
		}
	} else if notificationType == uint(CHECKNOTIFICATION) {
		if err := db.Raw(`
			SELECT u.id AS uid, u.fcm_token as fcm_token
			FROM users u
			WHERE u.id NOT IN (
				SELECT f.uid FROM face_infos f WHERE DATE(f.created_at) = ?
			) AND u.id NOT IN (
				SELECT t.uid FROM tap_blink_scores t WHERE DATE(t.created_at) = ?)
		`, today, today).Scan(&usersToNotify).Error; err != nil {
			return
		}
	}

	if len(usersToNotify) == 0 {
		log.Printf("No users to notify for notification type %d\n", notificationType)
		return
	}

	uniqueUsers, err := getLatestUsersByFCMToken(db)
	if err != nil {
		log.Printf("Error fetching unique users: %v\n", err)
		return
	}

	latestUserMap := make(map[string]uint)
	for _, user := range uniqueUsers {
		latestUserMap[user.FcmToken] = user.Uid
	}

	// 3. 유저별 미확인 알림 수 조회
	notificationCounts := getNotificationCounts(db)

	// 4. FCM 메시지와 DB에 저장할 메시지 준비
	var messages []*messaging.Message
	var newMessages []model.Message

	for _, user := range usersToNotify {
		if latestUserMap[user.FcmToken] != user.Uid {
			continue // 중복된 FCM 토큰 스킵
		}
		notificationCount := notificationCounts[user.Uid]
		// FCM 메시지 생성
		message := &messaging.Message{
			Data: map[string]string{
				"uid":                strconv.FormatUint(uint64(user.Uid), 10),
				"type":               strconv.FormatUint(uint64(notificationType), 10),
				"notification_count": strconv.FormatUint(uint64(notificationCount), 10),
				"time":               today,
				"parent_id":          strconv.Itoa(parentID),
			},
			Notification: &messaging.Notification{
				Title: title,
				Body:  body,
			},
			Token: user.FcmToken,
		}
		messages = append(messages, message)

		// DB에 저장할 메시지 생성
		newMessages = append(newMessages, model.Message{
			Uid:      user.Uid,
			Type:     notificationType,
			Body:     body,
			ParentID: uint(parentID),
			IsRead:   false,
		})
	}
	sendFCMMessages(db, messages, newMessages)
}

func getNotificationCounts(db *gorm.DB) map[uint]uint {
	var results []struct {
		Uid   uint
		Count uint
	}
	if err := db.Table("messages").
		Select("uid, COUNT(*) as count").
		Where("is_read = ?", false).
		Group("uid").
		Scan(&results).Error; err != nil {
		return nil
	}

	counts := make(map[uint]uint)
	for _, result := range results {
		counts[result.Uid] = result.Count
	}
	return counts
}

func sendFCMMessages(db *gorm.DB, messages []*messaging.Message, newMessages []model.Message) {
	if len(messages) == 0 {
		return
	}

	ctx := context.Background()
	var successCount, failureCount int64
	var wg sync.WaitGroup

	// 성공한 메시지를 수집할 채널
	successCh := make(chan model.Message, len(messages))

	for i, msg := range messages {
		wg.Add(1)

		go func(i int, msg *messaging.Message) {
			defer wg.Done() // 고루틴 종료 시 WaitGroup에서 하나 제거

			// FCM 메시지 전송
			_, err := firebaseClient.Send(ctx, msg)
			if err != nil {
				log.Printf("Failed to send message %d: %v\n", i, err)
				atomic.AddInt64(&failureCount, 1)
				return // 실패 시 채널에 추가하지 않음
			}

			log.Printf("Successfully sent message %d: %s\n", i, msg.Token)
			atomic.AddInt64(&successCount, 1)

			// 성공한 메시지를 채널에 추가
			successCh <- newMessages[i]
		}(i, msg) // 고루틴에서 변수 캡처
	}

	// 모든 고루틴이 종료되면 채널을 닫음
	go func() {
		wg.Wait()
		close(successCh)
	}()

	// 채널에서 성공한 메시지들을 수집
	var successfulMessages []model.Message
	for msg := range successCh {
		successfulMessages = append(successfulMessages, msg)
	}

	// 성공한 메시지를 DB에 저장
	if len(successfulMessages) > 0 {
		if err := db.Create(&successfulMessages).Error; err != nil {
			log.Printf("Error saving successful messages: %v\n", err)
		} else {
			log.Printf("Successfully saved %d messages to the database\n", len(successfulMessages))
		}
	}

	log.Printf("Total Success: %d, Total Failure: %d\n", successCount, failureCount)
}

func getTitle(notificationType uint) string {
	switch notificationType {
	case uint(MEDICINENOTIFICATION):
		return "약물 복용"
	case uint(EXERCISENOTIFICATION):
		return "운동 시간"
	default:
		return "알림"
	}
}

func getLatestUsersByFCMToken(db *gorm.DB) ([]struct {
	Uid      uint
	FcmToken string
}, error) {
	var users []struct {
		Uid      uint
		FcmToken string
	}

	query := `
		SELECT u.id AS uid, u.fcm_token AS fcm_token
		FROM users u
		INNER JOIN (
			SELECT fcm_token, MAX(updated_at) AS latest_update
			FROM users
			WHERE fcm_token IS NOT NULL AND fcm_token != ''
			GROUP BY fcm_token
		) latest_users ON u.fcm_token = latest_users.fcm_token
		AND u.updated_at = latest_users.latest_update
	`

	if err := db.Raw(query).Scan(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
