package core

type VimeoResponse struct {
	Data []struct {
		Type   string `json:"type"`
		Folder struct {
			Name string `json:"name"`
			Uri  string `json:"uri"`
		} `json:"folder"`
	} `json:"data"`
}

type VideoData struct {
	Id               uint     `json:"-"`
	SelectedVideos   []string `json:"selectedVideos"`
	DeselectedVideos []string `json:"deselectedVideos"`
}

type VimeoResponse2 struct {
	Data []struct {
		Name string `json:"name"`
		Uri  string `json:"uri"`
	} `json:"data"`
}

type VimeoResponse3 struct {
	Name     string `json:"name"`
	Uri      string `json:"uri"`
	Duration uint   `json:"duration"`
	Pictures struct {
		BaseLink string `json:"base_link"`
	} `json:"pictures"`
	ParentFolder struct {
		Uri  string `json:"uri"`
		Name string `json:"name"`
	} `json:"parent_folder"`
}

type VimeoLevel1 struct {
	ProjectId string `json:"project_id"`
	Name      string `json:"name"`
}

type VimeoLevel2 struct {
	Name     string `json:"name"`
	VideoId  string `json:"video_id"`
	IsActive bool   `json:"is_active"`
}

type SuccessResponse struct {
	Jwt string `json:"jwt"`
}
type ErrorResponse struct {
	Err string `json:"err"`
}

type BasicResponse struct {
	Code string `json:"code"`
}
