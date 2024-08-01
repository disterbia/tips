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
	AdminLogin(email string, password string) (string, error)
	SendInquire(inquire InquireRequest) (string, error)
	GetMyInquires(id uint, page uint, startDate, endDate string) ([]InquireResponse, error)
	AnswerInquire(answer InquireReplyRequest) (string, error)
	GetAllInquires(id uint, page uint, startDate, endDate string) ([]InquireResponse, error)
	RemoveInquire(id uint, uid uint) (string, error)
	RemoveReply(id uint, uid uint) (string, error)
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

func (service *inquireService) AdminLogin(email string, password string) (string, error) {
	var u model.User
	if err := service.db.Where("email=? AND phone=?", email, password).First(&u).Error; err != nil {
		return "", err
	}

	if u.RoleID != uint(ADMINROLE) {
		return "", errors.New("not admin")
	}

	// 새로운 JWT 토큰 생성
	tokenString, err := generateJWT(u)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}
func (service *inquireService) GetMyInquires(id uint, page uint, startDate, endDate string) ([]InquireResponse, error) {

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
			replies = append(replies, InquireReplyResponse{Id: w.ID, InquireId: w.InquireID, Content: w.Content, ReplyType: w.ReplyType, Created: w.Model.CreatedAt.Format("2006-01-02 15:04:05")})
		}
		inquireResponses = append(inquireResponses, InquireResponse{Id: v.ID, Email: v.Email, Title: v.Title, Content: v.Content,
			CreatedAt: v.Model.CreatedAt.Format("2006-01-02 15:04:05"), Replies: replies})
	}

	return inquireResponses, nil
}

func (service *inquireService) GetAllInquires(id uint, page uint, startDate, endDate string) ([]InquireResponse, error) {

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
			replies = append(replies, InquireReplyResponse{Id: w.ID, InquireId: w.InquireID, Content: w.Content, ReplyType: w.ReplyType, Created: w.Model.CreatedAt.Format("2006-01-02 15:04:05")})
		}
		inquireResponses = append(inquireResponses, InquireResponse{Id: v.ID, Email: v.Email, Title: v.Title, Content: v.Content,
			CreatedAt: v.Model.CreatedAt.Format("2006-01-02 15:04:05"), Replies: replies})
	}

	return inquireResponses, nil
}

func (service *inquireService) SendInquire(request InquireRequest) (string, error) {

	// 유효성 검사기 생성
	validate := validator.New()

	// 이메일 검증
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

func (service *inquireService) AnswerInquire(request InquireReplyRequest) (string, error) {
	var inquire model.Inquire

	if err := service.db.Preload("User").First(&inquire, request.InquireId).Error; err != nil {
		return "", errors.New("db error")
	}

	if request.ReplyType { // true = 답변
		if inquire.User.RoleID != uint(ADMINROLE) {
			return "", errors.New("unauthorized: user is not an admin")
		}
	} else { // 추가문의
		if inquire.User.ID != inquire.Uid {
			return "", errors.New("unauthorized: illegal user")
		}
		if inquire.User.RoleID == uint(ADMINROLE) {
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
		Email:          inquire.Email,                                              // 받는 사람의 이메일
		CreatedAt:      inquire.Model.CreatedAt.Format("2006-01-02 15:04:05"),      // 문의 생성 날짜
		Title:          inquire.Title,                                              // 이메일 제목
		Content:        inquire.Content,                                            // 문의 내용
		ReplyContent:   inquireReply.Content,                                       // 답변 내용
		ReplyCreatedAt: inquireReply.Model.CreatedAt.Format("2006-01-02 15:04:05"), // 답변 생성 날짜
	})

	if err != nil {
		tx.Rollback()
		log.Printf("Failed to send email: %v", err)
		return "", errors.New("Failed to send email")
	}

	log.Printf("send email: %v", reponse)
	tx.Commit()

	return "200", nil
}

func (service *inquireService) RemoveInquire(id uint, uid uint) (string, error) {

	var inquire model.Inquire
	if err := service.db.Where("id = ? AND uid = ?", id, uid).Delete(&inquire).Error; err != nil {
		return "", errors.New("db error")
	}
	return "200", nil
}

func (service *inquireService) RemoveReply(id uint, uid uint) (string, error) {

	var inquireReply model.InquireReply
	if err := service.db.Where("id = ? AND uid = ?", id, uid).Delete(&inquireReply).Error; err != nil {
		return "", errors.New("db error")
	}

	return "200", nil
}
