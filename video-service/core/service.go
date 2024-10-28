package core

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"video-service/model"

	"gorm.io/gorm"
)

type AdminVideoService interface {
	getLevel1s(id uint) ([]VimeoLevel1, error)
	getLevel2s(id uint, projectId string) ([]VimeoLevel2, error)
	saveVideos(videoData VideoData) (string, error)
}

type adminVideoService struct {
	db *gorm.DB
}

func NewAdminVideoService(db *gorm.DB) AdminVideoService {
	return &adminVideoService{db: db}
}
func (service *adminVideoService) getLevel1s(id uint) ([]VimeoLevel1, error) {
	var user model.User
	err := service.db.Where("id = ?", id).Find(&user).Error
	if err != nil {
		log.Println("db error")
		return nil, errors.New("db error")
	}

	if user.RoleID != uint(SUPERROLE) {
		log.Println("deny")
		return nil, errors.New("deny")
	}

	apiURL := "https://api.vimeo.com/users/145953562/projects/14798949/items"

	// HTTP 클라이언트 생성
	client := &http.Client{}

	// 요청 생성
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Println(err)
		return nil, err // 에러 반환
	}

	// Vimeo API 토큰 설정
	req.Header.Add("Authorization", "Bearer 915b8388768a803e93bac552f36e81a8")
	req.Header.Add("Accept", "application/vnd.vimeo.*+json;version=3.4")

	// 요청 보내기
	resp, err := client.Do(req)
	if err != nil {
		return nil, err // 에러 반환
	}
	defer resp.Body.Close()

	// 응답 처리
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err // 에러 반환
	}

	var response VimeoResponse
	err = json.Unmarshal(body, &response) // body는 이미 byte slice
	if err != nil {
		return nil, err // 에러 반환
	}

	var vimeoData []VimeoLevel1
	for _, item := range response.Data {
		if item.Type == "folder" {
			// URI에서 프로젝트 ID 추출
			splitUri := strings.Split(item.Folder.Uri, "/")
			projectId := splitUri[len(splitUri)-1]

			// VimeoLevel1 구조체로 데이터 매핑
			vimeoData = append(vimeoData, VimeoLevel1{
				ProjectId: projectId,
				Name:      item.Folder.Name,
			})
		}
	}

	return vimeoData, nil // 결과 반환
}

func (service *adminVideoService) getLevel2s(id uint, projectId string) ([]VimeoLevel2, error) {
	var user model.User
	err := service.db.Where("id = ?", id).Find(&user).Error
	if err != nil {
		log.Println("db error")
		return nil, errors.New("db error")
	}

	if user.RoleID != uint(SUPERROLE) {
		log.Println("deny")
		return nil, errors.New("deny")
	}

	apiURL := "https://api.vimeo.com/users/145953562/projects/" + projectId + "/videos"

	// HTTP 클라이언트 생성
	client := &http.Client{}

	// 요청 생성
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err // 에러 반환
	}

	// Vimeo API 토큰 설정
	req.Header.Add("Authorization", "Bearer 915b8388768a803e93bac552f36e81a8")
	req.Header.Add("Accept", "application/vnd.vimeo.*+json;version=3.4")

	// 요청 보내기
	resp, err := client.Do(req)

	if err != nil {
		return nil, err // 에러 반환
	}
	defer resp.Body.Close()

	// 응답 처리
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err // 에러 반환
	}

	var response VimeoResponse2
	err = json.Unmarshal(body, &response) // body는 이미 byte slice
	if err != nil {
		return nil, err // 에러 반환
	}

	var videos []model.Video
	err = service.db.Where("project_id = ?", projectId).Find(&videos).Error
	if err != nil {
		log.Println("db error")
		return nil, errors.New("db error")
	}

	var vimeoData []VimeoLevel2
	for _, item := range response.Data {
		isActive := false
		// URI에서 프로젝트 ID 추출
		splitUri := strings.Split(item.Uri, "/")
		videoId := splitUri[len(splitUri)-1]
		for _, v := range videos {
			if v.VideoId == videoId {
				isActive = true
			}
		}
		vimeoData = append(vimeoData, VimeoLevel2{
			VideoId:  videoId,
			Name:     item.Name,
			IsActive: isActive,
		})

	}

	return vimeoData, nil // 결과 반환
}

func (service *adminVideoService) saveVideos(videoData VideoData) (string, error) {
	var user model.User
	err := service.db.Where("id = ?", videoData.Id).Find(&user).Error
	if err != nil {
		log.Println("db error")
		return "", errors.New("db error")
	}

	if user.RoleID != uint(SUPERROLE) {
		log.Println("deny")
		return "", errors.New("deny")
	}

	selectedVideos := videoData.SelectedVideos
	deselectedVideos := videoData.DeselectedVideos

	// HTTP 클라이언트 생성
	client := &http.Client{}

	var wg sync.WaitGroup
	videosChan := make(chan model.Video, len(selectedVideos))
	proIdChan := make(chan string, len(selectedVideos))
	for _, item := range selectedVideos {
		wg.Add(1)
		go func(item string) {
			defer wg.Done()

			apiURL := "https://api.vimeo.com/users/145953562/videos/" + item

			// 요청 생성
			req, err := http.NewRequest("GET", apiURL, nil)
			if err != nil {
				log.Println(err)
				return
			}

			// Vimeo API 토큰 설정
			req.Header.Add("Authorization", "Bearer 915b8388768a803e93bac552f36e81a8")
			req.Header.Add("Accept", "application/vnd.vimeo.*+json;version=3.4")

			// 요청 보내기
			resp, err := client.Do(req)
			if err != nil {
				log.Println(err)
				return
			}
			defer resp.Body.Close()

			// 응답 처리
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
				return
			}

			var response VimeoResponse3
			err = json.Unmarshal(body, &response) // body는 이미 byte slice
			if err != nil {
				log.Println(err)
				return
			}
			if response.Name == "" {
				log.Println("no id")
				return
			}

			splitUri := strings.Split(response.ParentFolder.Uri, "/")
			projectId := splitUri[len(splitUri)-1]

			videosChan <- model.Video{
				VideoId:      item,
				Name:         response.Name,
				Duration:     response.Duration,
				ProjectId:    projectId,
				ThumbnailUrl: response.Pictures.BaseLink,
				ProjectName:  response.ParentFolder.Name,
			}
			proIdChan <- projectId
		}(item)
	}

	wg.Wait()
	close(videosChan)
	close(proIdChan)

	var videos []model.Video
	var proIds []string
	for video := range videosChan {
		videos = append(videos, video)
	}

	for proId := range proIdChan {
		proIds = append(proIds, proId)
	}

	if len(proIds) > 0 {
		if err := service.db.Where("project_id IN ?", proIds).Delete(&model.Video{}).Error; err != nil {
			log.Println("db error2")
			return "", errors.New("db error2")
		}

		if len(videos) > 0 {
			if err := service.db.Create(&videos).Error; err != nil {
				log.Println("db error3")
				return "", errors.New("db error3")
			}
		}
	}

	// 해제된 비디오 처리
	if len(deselectedVideos) > 0 {
		if err := service.db.Where("video_id IN ?", deselectedVideos).Delete(&model.Video{}).Error; err != nil {
			log.Println("db error4")
			return "", errors.New("db error4")
		}
	}

	return "200", nil

}
