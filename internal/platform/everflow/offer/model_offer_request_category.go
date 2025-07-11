package offer

type CategoryInfo struct {
	NetworkCategoryId int    `json:"network_category_id"`
	NetworkId         int    `json:"network_id"`
	Name              string `json:"name"`
	Status            string `json:"status"`
	TimeCreated       int    `json:"time_created"`
	TimeSaved         int    `json:"time_saved"`
}
