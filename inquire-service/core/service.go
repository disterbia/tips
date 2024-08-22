package core

import (
	"context"
	"errors"

	"inquire-service/model"
	"log"

	pb "inquire-service/proto"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type InquireService interface {
	sendInquire(inquire InquireRequest) (string, error)
	getMyInquires(id uint, page uint, startDate, endDate string) ([]InquireResponse, error)
	answerInquire(answer InquireReplyRequest) (string, error)
	getAllInquires(id uint, page uint, startDate, endDate string) ([]InquireResponse, error)
	removeInquire(id uint, uid uint) (string, error)
	removeReply(id uint, uid uint) (string, error)
}

type inquireService struct {
	db          *gorm.DB
	emailClient pb.EmailServiceClient
}

func NewInquireService(db *gorm.DB, conn *grpc.ClientConn) InquireService {
	emailClient := pb.NewEmailServiceClient(conn)
	return &inquireService{
		db:          db,
		emailClient: emailClient,
	}
}

func (service *inquireService) getMyInquires(id uint, page uint, startDate, endDate string) ([]InquireResponse, error) {

	if startDate != "" {
		if err := validateDate(startDate); err != nil {
			return nil, err
		}
	}
	if endDate != "" {
		if err := validateDate(endDate); err != nil {
			return nil, err
		}
	}

	pageSize := uint(30)
	var inquires []model.Inquire
	offset := page * pageSize

	query := service.db.Where("uid = ? ", id)
	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created_at <= ?", endDate+" 23:59:59")
	}
	query = query.Order("id DESC")
	result := query.Offset(int(offset)).Limit(int(pageSize)).Preload("Replies").Find(&inquires)

	if result.Error != nil {
		return nil, result.Error
	}

	var inquireResponses []InquireResponse

	for _, v := range inquires {
		var replies []InquireReplyResponse
		for _, w := range v.Replies {
			replies = append(replies, InquireReplyResponse{Id: w.ID, InquireId: w.InquireID, Content: w.Content, ReplyType: w.ReplyType, Created: w.CreatedAt.Format("2006-01-02 15:04:05")})
		}
		inquireResponses = append(inquireResponses, InquireResponse{Id: v.ID, Email: v.Email, Title: v.Title, Content: v.Content,
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"), Replies: replies})
	}

	return inquireResponses, nil
}

func (service *inquireService) getAllInquires(id uint, page uint, startDate, endDate string) ([]InquireResponse, error) {

	if startDate != "" {
		if err := validateDate(startDate); err != nil {
			return nil, err
		}
	}
	if endDate != "" {
		if err := validateDate(endDate); err != nil {
			return nil, err
		}
	}

	pageSize := uint(30)
	var inquires []model.Inquire
	offset := page * pageSize

	var user model.User
	if err := service.db.First(&user, id).Error; err != nil {
		return nil, errors.New("db error")
	}

	if user.RoleID != uint(ADMINROLE) {
		return nil, errors.New("unauthorized: user is not an admin")
	}

	query := service.db.Offset(int(offset)).Limit(int(pageSize))

	if startDate != "" {
		query = query.Where("created >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created <= ?", endDate)
	}
	query = query.Order("id DESC")
	if err := query.Preload("Replies").Find(&inquires).Error; err != nil {
		return nil, errors.New("db error2")
	}

	var inquireResponses []InquireResponse
	for _, v := range inquires {
		var replies []InquireReplyResponse
		for _, w := range v.Replies {
			replies = append(replies, InquireReplyResponse{Id: w.ID, InquireId: w.InquireID, Content: w.Content, ReplyType: w.ReplyType, Created: w.CreatedAt.Format("2006-01-02 15:04:05")})
		}
		inquireResponses = append(inquireResponses, InquireResponse{Id: v.ID, Email: v.Email, Title: v.Title, Content: v.Content,
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"), Replies: replies})
	}

	return inquireResponses, nil
}

func (service *inquireService) sendInquire(request InquireRequest) (string, error) {

	// 유효성 검사기 생성
	validate := validator.New()

	//유효성 검증
	if err := validate.Struct(request); err != nil {
		return "", err
	}

	inquire := model.Inquire{
		Uid:     request.Uid,
		Email:   request.Email,
		Content: request.Content,
		Title:   request.Title,
	}

	if err := service.db.Create(&inquire).Error; err != nil {
		return "", errors.New("db error")
	}

	return "200", nil
}

func (service *inquireService) answerInquire(request InquireReplyRequest) (string, error) {
	var inquire model.Inquire

	if err := service.db.Preload("User").First(&inquire, request.InquireId).Error; err != nil {
		return "", errors.New("db error")
	}

	if request.ReplyType { // true = 답변
		if inquire.User.RoleID != uint(SUPERROLE) {
			return "", errors.New("unauthorized: user is not an admin")
		}
	} else { // 추가문의
		if inquire.Uid != request.Uid {
			return "", errors.New("unauthorized: illegal user")
		}
		if inquire.User.RoleID == uint(ADMINROLE) || inquire.User.RoleID == uint(SUPERROLE) {
			return "", errors.New("unauthorized: can't admin ")
		}
	}

	inquireReply := model.InquireReply{
		Uid:       request.Uid,
		InquireID: request.InquireId,
		Content:   request.Content,
		ReplyType: request.ReplyType,
	}

	tx := service.db.Begin()
	if err := tx.Create(&inquireReply).Error; err != nil {
		tx.Rollback()
		return "", errors.New("db error")
	}

	reponse, err := service.emailClient.SendEmail(context.Background(), &pb.EmailRequest{
		Email:          inquire.Email,                                        // 받는 사람의 이메일
		CreatedAt:      inquire.CreatedAt.Format("2006-01-02 15:04:05"),      // 문의 생성 날짜
		Title:          inquire.Title,                                        // 이메일 제목
		Content:        inquire.Content,                                      // 문의 내용
		ReplyContent:   inquireReply.Content,                                 // 답변 내용
		ReplyCreatedAt: inquireReply.CreatedAt.Format("2006-01-02 15:04:05"), // 답변 생성 날짜
	})

	if err != nil {
		tx.Rollback()
		log.Printf("Failed to send email: %v", err)
		return "", errors.New("failed to send email")
	}

	log.Printf("send email: %v", reponse)
	tx.Commit()

	return "200", nil
}

func (service *inquireService) removeInquire(id uint, uid uint) (string, error) {

	var inquire model.Inquire
	if err := service.db.Where("id = ? AND uid = ?", id, uid).Delete(&inquire).Error; err != nil {
		return "", errors.New("db error")
	}
	return "200", nil
}

func (service *inquireService) removeReply(id uint, uid uint) (string, error) {

	var inquireReply model.InquireReply
	if err := service.db.Where("id = ? AND uid = ?", id, uid).Delete(&inquireReply).Error; err != nil {
		return "", errors.New("db error")
	}

	return "200", nil
}
