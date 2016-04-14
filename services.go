package cfclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type servicesResponse struct {
	Count     int                `json:"total_results"`
	Pages     int                `json:"total_pages"`
	Resources []servicesResource `json:"resources"`
}

type servicesResource struct {
	Meta   Meta    `json:"metadata"`
	Entity Service `json:"entity"`
}

type Service struct {
	Guid  string `json:"guid"`
	Label string `json:"label"`
	c     *Client
}

func (c *Client) ListServices(filter Filters) (services []Service, err error) {
	var serviceResp servicesResponse

	path := "/v2/services"
	filterPath := filter.ToParam()
	if filterPath != "" {
		path = path + "?" + filterPath
	}
	fmt.Println(path)
	r := c.newRequest("GET", path)
	resp, err := c.doRequest(r)
	if err != nil {
		return nil, fmt.Errorf("Error requesting services: %v", err)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading services request: %v", err)
	}

	err = json.Unmarshal(resBody, &serviceResp)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling services: %v", err)
	}
	for _, service := range serviceResp.Resources {
		service.Entity.Guid = service.Meta.Guid
		service.Entity.c = c
		services = append(services, service.Entity)
	}
	return services, nil
}
