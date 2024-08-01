package core

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"math/big"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
	"user-service/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"google.golang.org/api/idtoken"
)

var jwtSecretKey = []byte("adapfit_mark")

func verifyJWT(c *fiber.Ctx) (uint, error) {
	// 헤더에서 JWT 토큰 추출
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return 0, errors.New("authorization header is required")
	}

	// 'Bearer ' 접두사 제거
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	id := uint((*claims)["id"].(float64))

	return id, nil
}

func generateJWT(user model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(), // 한달 유효 기간
	})

	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func decodeJwt(tokenString string) string {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		log.Println(err)
	}

	// MapClaims 타입으로 claims 확인
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// 'iss' 확인
		if iss, ok := claims["iss"].(string); ok {
			fmt.Println("Issuer (iss):", iss)
			return iss
		} else {
			fmt.Println("'iss' 이 없습니다.")
			return ""
		}
	} else {
		log.Println("클레임을 MapClaims로 변환할 수 없습니다.")
		return ""
	}
}

// func copyStruct(src, dst interface{}) error {
// 	srcVal := reflect.ValueOf(src)
// 	dstVal := reflect.ValueOf(dst).Elem()

// 	for i := 0; i < srcVal.NumField(); i++ {
// 		srcField := srcVal.Field(i)
// 		srcFieldName := srcVal.Type().Field(i).Name

// 		dstField := dstVal.FieldByName(srcFieldName)
// 		if dstField.IsValid() && dstField.Type() == srcField.Type() {
// 			dstField.Set(srcField)
// 		}
// 	}

// 	return nil
// }

// func copyStruct(input interface{}, output interface{}) error {
// 	jsonData, err := json.Marshal(input)
// 	if err != nil {
// 		return err
// 	}

// 	err = json.Unmarshal(jsonData, output)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func validateTime(timeStr string) error {
	if len(timeStr) != 5 {
		return errors.New("invalid time format, should be HH:MM")
	}
	_, err := time.Parse("15:04", timeStr)
	if err != nil {
		return errors.New("invalid time format, should be HH:MM")
	}
	return nil
}

func validatePhoneNumber(phone string) error {
	// 정규 표현식 패턴: 010으로 시작하며 총 11자리 숫자
	pattern := `^010\d{8}$`
	matched, err := regexp.MatchString(pattern, phone)
	if err != nil || !matched {
		return errors.New("invalid phone format, should be 01000000000")
	}
	return nil
}

func validateDate(dateStr string) (time.Time, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, errors.New("invalid date format, should be YYYY-MM-DD")
	}
	return date, nil
}
func validateSignIn(request LoginRequest) error {
	err := validatePhoneNumber(request.Phone)
	if err != nil {
		return err
	}
	name := strings.TrimSpace(request.Name)
	if utf8.RuneCountInString(name) > 5 || len(name) == 0 {
		return errors.New("invalid name")
	}
	return nil
}

func validatePhoneSignIn(request PhoneLoginRequest) error {
	err := validatePhoneNumber(request.Phone)
	if err != nil {
		return err
	}
	name := strings.TrimSpace(request.Name)
	if utf8.RuneCountInString(name) > 5 || len(name) == 0 {
		return errors.New("invalid name")
	}
	return nil
}

// Apple 공개키 가져오기
func getApplePublicKeys() (JWKS, error) {
	resp, err := http.Get("https://appleid.apple.com/auth/keys")
	if err != nil {
		return JWKS{}, err
	}
	defer resp.Body.Close()

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return JWKS{}, err
	}
	return jwks, nil
}

// Apple 공개키로 서명 검증
func verifyAppleIDToken(token string, jwks JWKS) (*jwt.Token, error) {
	kid, err := extractKidFromToken(token)
	if err != nil {
		return nil, err
	}

	var key *rsa.PublicKey
	for _, jwk := range jwks.Keys {
		if jwk.Kid == kid {
			nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
			if err != nil {
				return nil, err
			}
			eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
			if err != nil {
				return nil, err
			}

			n := big.NewInt(0).SetBytes(nBytes)
			e := big.NewInt(0).SetBytes(eBytes).Int64()
			key = &rsa.PublicKey{N: n, E: int(e)}
			break
		}
	}

	if key == nil {
		return nil, errors.New("appropriate public key not found")
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	return parsedToken, nil
}

// 카카오 공개키 가져오기
func getKakaoPublicKeys() (JWKS, error) {
	resp, err := http.Get("https://kauth.kakao.com/.well-known/jwks.json")
	if err != nil {
		return JWKS{}, err
	}
	defer resp.Body.Close()

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return JWKS{}, err
	}
	return jwks, nil
}

// 카카오 공개키로 서명 검증
func verifyKakaoTokenSignature(token string, jwks JWKS) (*jwt.Token, error) {
	kid, err := extractKidFromToken(token)
	if err != nil {
		return nil, err
	}

	var key *rsa.PublicKey
	for _, jwk := range jwks.Keys {
		if jwk.Kid == kid {
			nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
			if err != nil {
				return nil, err
			}
			eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
			if err != nil {
				return nil, err
			}

			n := big.NewInt(0).SetBytes(nBytes)
			e := big.NewInt(0).SetBytes(eBytes).Int64()
			key = &rsa.PublicKey{N: n, E: int(e)}
			break
		}
	}

	if key == nil {
		return nil, errors.New("appropriate public key not found")
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	return parsedToken, nil
}

// ID 토큰에서 kid 추출
func extractKidFromToken(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", errors.New("invalid token format")
	}
	headerPart := parts[0]
	headerJson, err := base64.RawURLEncoding.DecodeString(headerPart)
	if err != nil {
		return "", err
	}

	var header map[string]interface{}
	if err := json.Unmarshal(headerJson, &header); err != nil {
		return "", err
	}

	kid, ok := header["kid"].(string)
	if !ok {
		return "", errors.New("kid not found in token header")
	}
	return kid, nil
}

// Google ID 토큰을 검증하고 이메일을 반환
func validateGoogleIDToken(idToken string) (string, error) {
	log.Print("idToken: ", idToken)
	// idtoken 패키지를 사용하여 토큰 검증
	payload, err := idtoken.Validate(context.Background(), idToken, "390432007084-1hqslpiclba2hucb6hl41acecv1qekbt.apps.googleusercontent.com")
	if err != nil {
		log.Printf("Token validation error: %v", err)
		return "", err
	}

	// 이메일 추출
	email, ok := payload.Claims["email"].(string)
	if !ok {
		return "", errors.New("email claim not found in token")
	}

	return email, nil
}

func appleLogin(request LoginRequest) (string, error) {
	jwks, err := getApplePublicKeys()
	if err != nil {
		return "", err
	}

	parsedToken, err := verifyAppleIDToken(request.IdToken, jwks)
	if err != nil {
		return "", err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		email, ok := claims["email"].(string)
		if !ok {
			return "", errors.New("email not found in token claims")
		}

		return email, nil

	}
	return "", errors.New("invalid token")

}
func kakaoLogin(request LoginRequest) (string, error) {
	jwks, err := getKakaoPublicKeys()
	if err != nil {
		return "", err
	}

	parsedToken, err := verifyKakaoTokenSignature(request.IdToken, jwks)
	if err != nil {
		return "", err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		email, ok := claims["email"].(string)
		if !ok {
			return "", errors.New("email not found in token claims")
		}

		return email, nil
	}
	return "", errors.New("invalid token")

}

func googleLogin(request LoginRequest) (string, error) {
	email, err := validateGoogleIDToken(request.IdToken)
	if err != nil {
		return "", err
	}
	return email, nil
}

func deleteFromS3(fileKey string, s3Client *s3.S3, bucket string, bucketUrl string) error {

	// URL에서 객체 키 추출
	key := extractKeyFromUrl(fileKey, bucket, bucketUrl)
	log.Println("key", fileKey)

	_, err := s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	// 에러 발생 시 처리 로직
	if err != nil {
		fmt.Printf("Failed to delete object from S3: %s, error: %v\n", fileKey, err)
	}

	return err
}

// URL에서 S3 객체 키를 추출하는 함수
func extractKeyFromUrl(url, bucket string, bucketUrl string) string {
	prefix := fmt.Sprintf("https://%s.%s/", bucket, bucketUrl)
	return strings.TrimPrefix(url, prefix)
}

func uploadImagesToS3(imgData []byte, thumbnailData []byte, contentType string, ext string, s3Client *s3.S3, bucket string, bucketUrl string, uid string) (string, string, error) {
	// 이미지 파일 이름과 썸네일 파일 이름 생성
	imgFileName := "images/profile/" + uid + "/images/" + uuid.New().String() + ext
	thumbnailFileName := "images/profile/" + uid + "/thumbnails/" + uuid.New().String() + ext

	// S3에 이미지 업로드
	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(imgFileName),
		Body:        bytes.NewReader(imgData),
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return "", "", err
	}

	// S3에 썸네일 업로드
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(thumbnailFileName),
		Body:        bytes.NewReader(thumbnailData),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", "", err
	}

	// 업로드된 이미지와 썸네일의 URL 생성 및 반환
	imgURL := "https://" + bucket + "." + bucketUrl + "/" + imgFileName
	thumbnailURL := "https://" + bucket + "." + bucketUrl + "/" + thumbnailFileName

	return imgURL, thumbnailURL, nil
}

func reduceImageSize(data []byte) ([]byte, error) {
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	log.Println("image size: ", len(data))
	// 원본 이미지의 크기를 절반씩 줄이면서 10MB 이하로 만듦
	for len(data) > 10*1024*1024 {
		newWidth := img.Bounds().Dx() / 2
		newHeight := img.Bounds().Dy() / 2

		resizedImg := resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)

		var buf bytes.Buffer
		switch format {
		case "jpeg":
			err = jpeg.Encode(&buf, resizedImg, nil)
		case "png":
			err = png.Encode(&buf, resizedImg)
		case "gif":
			err = gif.Encode(&buf, resizedImg, nil)
		case "webp":
			// WebP 인코딩은 지원하지 않으므로 PNG 형식으로 인코딩
			err = png.Encode(&buf, resizedImg)
		// 여기에 필요한 다른 형식을 추가할 수 있습니다.
		default:
			log.Printf("Unsupported format: %s\n", format)
			return nil, err
		}
		if err != nil {
			return nil, err
		}

		data = buf.Bytes()
		img = resizedImg
	}

	return data, nil
}

func createThumbnail(data []byte) ([]byte, error) {
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	// 썸네일의 크기를 절반씩 줄이면서 1MB 이하로 만듦
	for {
		newWidth := img.Bounds().Dx() / 2
		newHeight := img.Bounds().Dy() / 2

		thumbnail := resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)

		var buf bytes.Buffer
		switch format {
		case "jpeg":
			err = jpeg.Encode(&buf, thumbnail, nil)
		case "png":
			err = png.Encode(&buf, thumbnail)
		case "gif":
			err = gif.Encode(&buf, thumbnail, nil)
		case "webp":
			err = png.Encode(&buf, thumbnail)
		default:
			log.Printf("Unsupported format: %s\n", format)
			return nil, err
		}
		if err != nil {
			return nil, err
		}

		thumbnailData := buf.Bytes()
		log.Println("thumbnailData size: ", len(thumbnailData))
		if len(thumbnailData) < 1024*1024 {
			return thumbnailData, nil
		}

		img = thumbnail
	}
}

func getImageFormat(imgData []byte) (contentType, extension string, err error) {
	_, format, err := image.DecodeConfig(bytes.NewReader(imgData))
	if err != nil {
		return "", "", err
	}
	switch format {
	case "jpeg":
		contentType = "image/jpeg"
		extension = ".jpg"
	case "png":
		contentType = "image/png"
		extension = ".png"
	case "gif":
		contentType = "image/gif"
		extension = ".gif"
	case "wepb":
		contentType = "image/wepb"
		extension = ".wepb"
	default:
		return "", "", fmt.Errorf("unsupported image format: %s", format)
	}

	return contentType, extension, nil
}
