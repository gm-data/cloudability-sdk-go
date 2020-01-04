package cloudability

import (
	"strconv"
)


type viewsEndpoint struct {
	*cloudabilityV3Endpoint
}

func newViewsEndpoint(apikey string) *viewsEndpoint {
	e := &viewsEndpoint{newCloudabilityV3Endpoint(apikey)}
	e.EndpointPath = "/v3/views/"
	return e
}

type ViewFilter struct {
	Field string `json:"field"`
	Comparator string `json:"comparator"`
	Value string `json:"value"`
}

type View struct {
	Id string `json:"id"`
	Title string `json:"title"`
	SharedWithUsers []string `json:"sharedWithUsers"`
	SharedWithOrganization bool `json:"sharedWithOrganization"`
	OwnerId string `json:"ownerId"`
	Filters []ViewFilter `json:"filters"`
}

func (e viewsEndpoint) GetViews() ([]View, error) {
	var views []View
	err := e.get("", &views)
	return views, err
}

func (e viewsEndpoint) GetView(id int) (*View, error) {
	var view View
	err := e.get(strconv.Itoa(id), &view)
	return &view, err
}