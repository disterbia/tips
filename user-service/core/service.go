// /user-service/service/service.go

package core

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"user-service/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type UserService interface {
	snsLogin(request LoginRequest) (string, error)
	phoneLogin(request PhoneLoginRequest) (string, error)
	autoLogin(request AutoLoginRequest) (string, error)
	verifyAuthCode(number, code string) (string, error)
	sendAuthCodeForSingin(number string) (string, error)
	sendAuthCodeForLogin(number string) (string, error)
	updateUser(request UserRequest) (string, error)
	getUser(id uint) (UserResponse, error)
	linkEmail(uid uint, idToken string) (string, error)
	removeUser(id uint) (string, error)
	getVersion() (AppVersionResponse, error)
	getPolices() ([]PoliceResponse, error)
	exchangeCodeForToken(code string) (*TokenResponse, error)
}

type userService struct {
	db        *gorm.DB
	s3svc     *s3.S3
	bucket    string
	bucketUrl string
}

func NewUserService(db *gorm.DB, s3svc *s3.S3, bucket string, bucketUrl string) UserService {
	return &userService{db: db, s3svc: s3svc, bucket: bucket, bucketUrl: bucketUrl}
}

type PublicKey struct {
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type JWKS struct {
	Keys []PublicKey `json:"keys"`
}

func (service *userService) snsLogin(request LoginRequest) (string, error) {
	if request.FCMToken == "" || request.DeviceID == "" {
		log.Println("check")
		return "", errors.New("check fcm_token,device_id")
	}
	iss, err := decodeJwt(request.IdToken)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var email string
	var snsType uint

	if strings.Contains(iss, "kakao") { // 카카오
		snsType = uint(KAKAO)
		if email, err = kakaoLogin(request); err != nil {
			log.Println(err)
			return "", err
		}
	} else if strings.Contains(iss, "google") { // 구글
		snsType = uint(GOOGLE)
		if email, err = googleLogin(request); err != nil {
			log.Println(err)
			return "", err
		}
	} else if strings.Contains(iss, "apple") { // 애플
		snsType = uint(APPLE)
		if email, err = appleLogin(request); err != nil {
			log.Println(err)
			return "", err
		}
	}

	var user model.User
	if err := service.db.Where("email = ? ", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 연동 로그인
			// 연동 이메일 목록에 없다면 유저생성 (같은번호 있으면 db에서 생성안됨. 같은번호 없으면 회원가입과 같음)
			var linkedEmail model.LinkedEmail
			if err := service.db.Where("email = ?", email).First(&linkedEmail).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// 유효성 검사 수행
					if err := validateSignIn(request); err != nil {
						log.Println("-2")
						return "", errors.New("-2")
					}
					now := time.Now()
					thirtyMinutesAgo := now.Add(-30 * time.Minute)

					if err := service.db.Where("target = ? AND created_at >= ?", request.Phone, thirtyMinutesAgo).Last(&model.VerifiedTarget{}).Error; err != nil {
						if errors.Is(err, gorm.ErrRecordNotFound) {
							log.Println("-1")
							return "", errors.New("-1") // 인증해야함
						}
						log.Println("db error")
						return "", errors.New("db error")
					}

					birthday, err := time.Parse("2006-01-02", request.Birthday)
					if err != nil {
						log.Println("db error-2")
						return "", errors.New("-2")
					}

					user = model.User{Name: request.Name, Email: &email, SnsType: snsType, DeviceID: request.DeviceID, FCMToken: request.FCMToken, Phone: request.Phone, Gender: request.Gender,
						Birthday: birthday, UserType: request.UserType}
					if err := service.db.Create(&user).Error; err != nil {
						if strings.Contains(err.Error(), "duplicate") {
							if err := service.db.Where("phone = ? ", request.Phone).First(&user).Error; err != nil {
								log.Println("db error2")
								return "", errors.New("db error2")
							}
							log.Println("snstype")
							return "", errors.New(strconv.Itoa(int(user.SnsType))) // 이미 가입된 번호
						}
						log.Println("db error3")
						return "", errors.New("db error3")
					}
				} else {
					log.Println("db error4")
					return "", errors.New("db error4")

				}
			} else {
				// 있다면 해당 이메일의 uid로 조회
				if err := service.db.Where("id = ?", linkedEmail.Uid).First(&user).Error; err != nil {
					log.Println("db error5")
					return "", errors.New("db error5")
				}
			}

			if err := service.db.Model(&user).Updates(model.User{FCMToken: request.FCMToken, DeviceID: request.DeviceID}).Error; err != nil {
				log.Println("db error6")
				return "", errors.New("db error6")
			}

		} else {
			log.Println("db error7")
			return "", errors.New("db error7")
		}
	} else {
		if err := service.db.Model(&user).Updates(model.User{FCMToken: request.FCMToken, DeviceID: request.DeviceID}).Error; err != nil {
			log.Println("db error8")
			return "", errors.New("db error8")
		}
	}

	// JWT 토큰 생성
	tokenString, err := generateJWT(user)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (service *userService) phoneLogin(request PhoneLoginRequest) (string, error) {
	if request.FCMToken == "" || request.DeviceID == "" {
		log.Println("check")
		return "", errors.New("check fcm_token,device_id")
	}
	now := time.Now()
	threeMinutesAgo := now.Add(-3 * time.Minute)
	var verify model.VerifiedTarget

	if err := service.db.Where("target = ? AND created_at >= ?", request.Phone, threeMinutesAgo).Last(&verify).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("-1")
			return "", errors.New("-1") // 인증해야함
		}
		log.Println("db error")
		return "", errors.New("db error")
	}

	if err := service.db.Unscoped().Delete(&verify).Error; err != nil {
		log.Println("db error2")
		return "", errors.New("db error2")
	}

	var user model.User

	if err := service.db.Where("phone = ? ", request.Phone).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			// 유효성 검사 수행
			if err := validatePhoneSignIn(request); err != nil {
				log.Println("-2")
				return "", errors.New("-2")
			}

			birthday, err := time.Parse("2006-01-02", request.Birthday)
			if err != nil {
				log.Println("-2")
				return "", errors.New("-2")
			}

			user = model.User{Name: request.Name, DeviceID: request.DeviceID, FCMToken: request.FCMToken, Phone: request.Phone, Gender: request.Gender,
				Birthday: birthday, UserType: request.UserType}
			if err := service.db.Create(&user).Error; err != nil {
				log.Println("db error3")
				return "", errors.New("db error3")
			}
		} else {
			log.Println("db error4")
			return "", errors.New("db error4")
		}
	} else {
		if user.SnsType != 0 {
			log.Println("snstype")
			return "", errors.New(strconv.Itoa(int(user.SnsType))) //이미 가입된 번호
		}

		if err := service.db.Model(&user).Updates(model.User{FCMToken: request.FCMToken, DeviceID: request.DeviceID}).Error; err != nil {
			log.Println("db error4")
			return "", errors.New("db error4")
		}
	}

	// JWT 토큰 생성
	tokenString, err := generateJWT(user)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return tokenString, nil

}

func (service *userService) autoLogin(request AutoLoginRequest) (string, error) {
	if request.FcmToken == "" || request.DeviceId == "" {
		log.Println("check")
		return "", errors.New("check fcm_token,device_id")
	}
	// 데이터베이스에서 사용자 조회
	var u model.User
	if err := service.db.Where("id = ?", request.Id).First(&u).Error; err != nil {
		log.Println("db error")
		return "", errors.New("db error")
	}
	// 새로운 JWT 토큰 생성
	tokenString, err := generateJWT(u)
	if err != nil {
		return "", err
	}

	if err := service.db.Model(&u).Updates(model.User{FCMToken: request.FcmToken, DeviceID: request.DeviceId}).Error; err != nil {
		log.Println("db error2")
		return "", errors.New("db error2")
	}
	return tokenString, nil
}

func (service *userService) sendAuthCodeForSingin(number string) (string, error) {
	err := validatePhoneNumber(number)
	if err != nil {
		return "", err
	}

	//존재하는 번호인지 체크
	result := service.db.Debug().Where("phone = ?", number).Find(&model.User{})
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Println("db error")
			return "", errors.New("db error")
		}

	} else if result.RowsAffected > 0 {
		// 레코드가 존재할 때
		log.Println("-1")
		return "", errors.New("-1")
	}

	code, err := sendCode(number)

	if err != nil {
		return "", err
	}

	if err := service.db.Create(&model.AuthCode{Phone: number, Code: code}).Error; err != nil {
		return "", err
	}
	return "200", nil
}

func (service *userService) sendAuthCodeForLogin(number string) (string, error) {

	if err := validatePhoneNumber(number); err != nil {
		return "", err
	}

	code, err := sendCode(number)

	if err != nil {
		return "", err
	}

	if err := service.db.Create(&model.AuthCode{Phone: number, Code: code}).Error; err != nil {
		return "", err
	}
	return "200", nil
}

func (service *userService) verifyAuthCode(number, code string) (string, error) {
	now := time.Now()
	threeMinutesAgo := now.Add(-3 * time.Minute)
	var authCode model.AuthCode

	if err := service.db.Where("phone = ? AND created_at >= ? ", number, threeMinutesAgo).Last(&authCode).Error; err != nil {
		log.Println("db error")
		return "", errors.New("db error")
	}
	if authCode.Code != code {
		log.Println("-1")
		return "", errors.New("-1")
	}

	tx := service.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	if err := tx.Where("phone = ?", authCode.Phone).Unscoped().Delete(&model.AuthCode{}).Error; err != nil {
		tx.Rollback()
		log.Println("db error3")
		return "", errors.New("db error3")
	}

	if err := tx.Create(&model.VerifiedTarget{Target: authCode.Phone}).Error; err != nil {
		tx.Rollback()
		log.Println("db error2")
		return "", errors.New("db error2")
	}

	tx.Commit()
	return "200", nil
}

func (service *userService) updateUser(request UserRequest) (string, error) {
	var birtday time.Time
	// 유효성 검사 수행

	if birth, err := validateDate(request.Birthday); err != nil {
		log.Println(err)
		return "", err
	} else {
		birtday = birth
	}

	if err := validateSignInForUpdate(request); err != nil {
		log.Println(err)
		return "", err
	}

	var user model.User
	if err := service.db.Where("id= ? ", request.ID).First(&user).Error; err != nil {
		log.Println("db error")
		return "", errors.New("db error")
	}
	if user.Phone != request.Phone {
		now := time.Now()
		threeMinutesAgo := now.Add(-3 * time.Minute)
		var verify model.VerifiedTarget

		if err := service.db.Where("target = ? AND created_at >= ?", request.Phone, threeMinutesAgo).Last(&verify).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Println("-1")
				return "", errors.New("-1") // 인증해야함
			}
			log.Println("db error1")
			return "", errors.New("db error1")
		}

		if err := service.db.Unscoped().Delete(&verify).Error; err != nil {
			log.Println("db error2")
			return "", errors.New("db error2")
		}

	}

	imageChan := make(chan model.Image)
	errorChan := make(chan error)

	if request.ProfileImage != "" {
		go func() {
			defer close(imageChan)
			defer close(errorChan)
			imgData, err := base64.StdEncoding.DecodeString(request.ProfileImage)
			if err != nil {
				errorChan <- err
				return
			}

			contentType, ext, err := getImageFormat(imgData)
			if err != nil {
				errorChan <- err
				return
			}

			var wg sync.WaitGroup
			var resizedImg, thumbnailData []byte

			// 이미지 크기 조정 (10MB 제한)
			wg.Add(1)
			go func() {
				defer wg.Done()
				if len(imgData) > 10*1024*1024 {
					resizedImg, err = reduceImageSize(imgData)
					if err != nil {
						errorChan <- err
						return
					}
				} else {
					resizedImg = imgData
				}
			}()

			// 썸네일 이미지 생성
			wg.Add(1)
			go func() {
				defer wg.Done()
				thumbnailData, err = createThumbnail(imgData)
				if err != nil {
					errorChan <- err
					return
				}
			}()

			wg.Wait()

			// S3에 이미지 및 썸네일 업로드
			fileName, thumbnailFileName, err := uploadImagesToS3(resizedImg, thumbnailData, contentType, ext, service.s3svc, service.bucket, service.bucketUrl, strconv.FormatUint(uint64(request.ID), 10))
			if err != nil {
				errorChan <- err
				return
			}
			imageChan <- model.Image{
				Uid:          request.ID,
				Url:          fileName,
				ThumbnailUrl: thumbnailFileName,
				ParentId:     request.ID,
				Type:         uint(PROFILEIMAGETYPE),
			}
		}()
	}

	// 트랜잭션 시작
	tx := service.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	//유저 정보 업데이트
	user.Name = request.Name
	user.Phone = request.Phone
	user.Birthday = birtday
	user.UserType = request.UserType
	user.Gender = request.Gender
	if err := tx.Save(&user).Error; err != nil {
		log.Println(err.Error())
		tx.Rollback()
		if request.ProfileImage != "" {
			// 이미 업로드된 파일들을 S3에서 삭제
			go func() {
				select {
				case image := <-imageChan:
					deleteFromS3(image.Url, service.s3svc, service.bucket, service.bucketUrl)
					deleteFromS3(image.ThumbnailUrl, service.s3svc, service.bucket, service.bucketUrl)
				case <-errorChan:
					log.Println(err)
				}
			}()
		}
		log.Println("db error3")
		return "", errors.New("db error3")
	}

	if request.ProfileImage != "" {
		select {
		case image := <-imageChan:
			// 기존 이미지 레코드 논리삭제
			if err := service.db.Debug().Where("parent_id = ? AND type =?", user.ID, PROFILEIMAGETYPE).Delete(&model.Image{}).Error; err != nil {
				log.Println(err.Error())
				tx.Rollback()
				go func() {
					deleteFromS3(image.Url, service.s3svc, service.bucket, service.bucketUrl)
					deleteFromS3(image.ThumbnailUrl, service.s3svc, service.bucket, service.bucketUrl)
				}()
				log.Println("db error4")
				return "", errors.New("db error4")
			}

			// 이미지 레코드 생성
			if err := tx.Create(&image).Error; err != nil {
				log.Println(err)
				tx.Rollback()
				go func() {
					deleteFromS3(image.Url, service.s3svc, service.bucket, service.bucketUrl)
					deleteFromS3(image.ThumbnailUrl, service.s3svc, service.bucket, service.bucketUrl)
				}()
				log.Println("db error5")
				return "", errors.New("db error5")
			}
		case err := <-errorChan:
			tx.Rollback()
			log.Println(err)
			return "", err
		}
	}

	tx.Commit()
	return "200", nil
}

func (service *userService) linkEmail(uid uint, idToken string) (string, error) {
	iss, err := decodeJwt(idToken)

	if err != nil {
		return "", err
	}

	if strings.Contains(iss, "kakao") { //카카오
		jwks, err := getKakaoPublicKeys()
		if err != nil {
			return "", err
		}

		parsedToken, err := verifyKakaoTokenSignature(idToken, jwks)
		if err != nil {
			return "", err
		}

		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
			email, ok := claims["email"].(string)
			if !ok {
				log.Println("claims")
				return "", errors.New("email not found in token claims")
			}
			if err := saveLinkedEmail(uid, email, service, uint(KAKAO)); err != nil {
				return "", err
			}
		}

	} else if strings.Contains(iss, "google") { // 구글
		email, err := validateGoogleIDToken(idToken)
		if err != nil {
			return "", err
		}
		if err := saveLinkedEmail(uid, email, service, uint(GOOGLE)); err != nil {
			return "", err
		}

	} else if strings.Contains(iss, "apple") { // 애플
		jwks, err := getApplePublicKeys()
		if err != nil {
			return "", err
		}

		parsedToken, err := verifyAppleIDToken(idToken, jwks)
		if err != nil {
			return "", err
		}

		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
			email, ok := claims["email"].(string)
			if !ok {
				log.Println("claims")
				return "", errors.New("email not found in token claims")
			}
			if err := saveLinkedEmail(uid, email, service, uint(APPLE)); err != nil {
				return "", err
			}
		}
	} else {
		log.Println("invalid")
		return "", errors.New("invalid snsType")
	}
	return "200", nil
}

func saveLinkedEmail(uid uint, email string, service *userService, snsType uint) error {
	var user model.User
	if err := service.db.Where("id = ? ", uid).First(&user).Error; err != nil {
		log.Println("db error")
		return errors.New("db error")
	}
	if *user.Email == email {
		log.Println("wrong")
		return errors.New("wrong request")
	}

	linkedEmail := model.LinkedEmail{Email: email, Uid: uid, SnsType: snsType}

	result := service.db.Where(linkedEmail).First(&model.LinkedEmail{})

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 레코드가 존재하지 않으면 새 레코드 생성
		err := service.db.Create(&linkedEmail).Error
		if err != nil {
			log.Println("db erro2")
			return errors.New("db error2")
		}
	} else if result.Error != nil {
		log.Println("db error3")
		return errors.New("db error3")
	} else {
		// 레코드가 존재하면 삭제
		if err := service.db.Where(linkedEmail).Unscoped().Delete(&model.LinkedEmail{}).Error; err != nil {
			log.Println("db error")
			return errors.New("db error4")
		}
	}

	return nil
}

func (service *userService) getUser(id uint) (UserResponse, error) {
	var user model.User
	result := service.db.Debug().Preload("ProfileImages", "type = ?", PROFILEIMAGETYPE).
		Preload("LinkedEmails").First(&user, id)
	if result.Error != nil {
		log.Println("db error")
		return UserResponse{}, errors.New("db error")
	}

	var userResponse UserResponse

	if len(user.ProfileImages) > 0 && user.ProfileImages[0].Url != "" { // gorm은 nil이 아닌 빈 슬라이스로 매핑함
		urlkey := extractKeyFromUrl(user.ProfileImages[0].Url, service.bucket, service.bucketUrl)
		thumbnailUrlkey := extractKeyFromUrl(user.ProfileImages[0].ThumbnailUrl, service.bucket, service.bucketUrl)
		// 사전 서명된 URL을 생성
		url, _ := service.s3svc.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(service.bucket),
			Key:    aws.String(urlkey),
		})
		thumbnailUrl, _ := service.s3svc.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(service.bucket),
			Key:    aws.String(thumbnailUrlkey),
		})
		urlStr, err := url.Presign(5 * time.Second) // URL은 5초 동안 유효
		if err != nil {
			return UserResponse{}, err
		}
		thumbnailUrlStr, err := thumbnailUrl.Presign(5 * time.Second) // URL은 5초 동안 유효 CachedNetworkImage 에서 캐싱해서 쓰면됨
		if err != nil {
			return UserResponse{}, err
		}
		userResponse.ProfileImage.Url = urlStr // 사전 서명된 URL로 업데이트
		userResponse.ProfileImage.ThumbnailUrl = thumbnailUrlStr
	}
	var linkedEmails []LinkedResponse
	for _, v := range user.LinkedEmails {
		linkedEmail := LinkedResponse{SnsType: v.SnsType, Email: v.Email}
		linkedEmails = append(linkedEmails, linkedEmail)
	}
	userResponse.Birthday = user.Birthday.Format("2006-01-02")
	userResponse.Gender = user.Gender
	userResponse.Name = user.Name
	userResponse.Phone = user.Phone
	userResponse.SnsType = user.SnsType
	userResponse.LinkedEmails = linkedEmails
	userResponse.CreatedAt = user.CreatedAt.Format("2006-01-02")
	userResponse.UserType = user.UserType

	return userResponse, nil
}

func (service *userService) removeUser(id uint) (string, error) {
	if err := service.db.Where("id = ?", id).Delete(&model.User{}).Error; err != nil {
		log.Println("db error")
		return "", errors.New("db error")
	}
	return "200", nil
}

func (service *userService) getVersion() (AppVersionResponse, error) {
	var version model.AppVersion
	if err := service.db.Last(&version).Error; err != nil {
		log.Println("db error")
		return AppVersionResponse{}, errors.New("db error")
	}

	versionResponse := AppVersionResponse{LatestVersion: version.LatestVersion, AndroidLink: version.AndroidLink, IosLink: version.IosLink}
	return versionResponse, nil
}

func (service *userService) getPolices() ([]PoliceResponse, error) {
	var polices []model.Police
	if err := service.db.Where("is_last = true").Find(&polices).Error; err != nil {
		log.Println("db error")
		return nil, errors.New("db error")
	}

	var policeResponse []PoliceResponse

	for _, v := range polices {
		policeResponse = append(policeResponse, PoliceResponse{Title: v.Title, PoliceType: v.PoliceType, Body: v.Body})
	}

	return policeResponse, nil
}

// ExchangeCodeForToken은 애플 서버와 통신해 Authorization Code를 토큰으로 교환합니다.
func (s *userService) exchangeCodeForToken(code string) (*TokenResponse, error) {
	clientID := os.Getenv("APPLE_CLIENT_ID")
	keyID := os.Getenv("APPLE_KEY_ID")
	teamID := os.Getenv("APPLE_TEAM_ID")
	privateKey := os.Getenv("APPLE_PRIVATE_KEY")

	// client_secret 생성
	clientSecret, err := GenerateClientSecret(keyID, teamID, clientID, privateKey)
	if err != nil {
		return nil, err
	}

	// 요청 데이터 생성
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", "https://haruharulab.com/user/apple/callback")

	// POST 요청 생성
	req, err := http.NewRequest("POST", "https://appleid.apple.com/auth/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 요청 전송
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 응답 확인
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to exchange token, status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	// 응답 파싱
	var tokenResponse TokenResponse
	tokenResponse.Code = code
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}
