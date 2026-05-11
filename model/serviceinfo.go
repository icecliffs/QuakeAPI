package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type ServiceDataList []ServiceData

func (s *ServiceDataList) UnmarshalJSON(data []byte) error {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" || trimmed == "null" {
		*s = nil
		return nil
	}

	if strings.HasPrefix(trimmed, "[") {
		var items []ServiceData
		if err := json.Unmarshal(data, &items); err != nil {
			return err
		}
		*s = items
		return nil
	}

	if strings.HasPrefix(trimmed, "{") {
		var single ServiceData
		if err := json.Unmarshal(data, &single); err == nil && single.IP != "" {
			*s = []ServiceData{single}
			return nil
		}

		var obj map[string]interface{}
		if err := json.Unmarshal(data, &obj); err != nil {
			return err
		}
		return fmt.Errorf("unexpected object response in data: %v", obj)
	}

	return fmt.Errorf("unexpected data response: %s", trimmed)
}

type ServiceData struct {
	IP        string    `json:"ip"`
	Port      int       `json:"port"`
	Hostname  string    `json:"hostname"`
	Transport string    `json:"transport"`
	Asn       int       `json:"asn"`
	Org       string    `json:"org"`
	Time      time.Time `json:"time"`
	Domain    string    `json:"domain"`
	OsName    string    `json:"os_name"`
	OsVersion string    `json:"os_version"`
	IsIpv6    bool      `json:"is_ipv6"`

	Service struct {
		Name     string `json:"name"`
		Product  string `json:"product"`
		Banner   string `json:"banner"`
		Version  string `json:"version"`
		Response string `json:"response"`
		HTTP     struct {
			Title        string `json:"title"`
			StatusCode   int    `json:"status_code"`
			Server       string `json:"server"`
			Host         string `json:"host"`
			Body         string `json:"body"`
			XPoweredBy   string `json:"x_powered_by"`
			MetaKeywords string `json:"meta_keywords"`
			Favicon      struct {
				Hash     string `json:"hash"`
				Data     string `json:"data"`
				Location string `json:"location"`
			} `json:"favicon"`
			ICP struct {
				Licence string `json:"licence"`
			} `json:"icp"`
		} `json:"http"`
	} `json:"service"`

	Location struct {
		CountryCn  string    `json:"country_cn"`
		ProvinceCn string    `json:"province_cn"`
		CityCn     string    `json:"city_cn"`
		CountryEn  string    `json:"country_en"`
		Isp        string    `json:"isp"`
		Gps        []float64 `json:"gps"`
	} `json:"location"`

	Components []struct {
		ProductNameCn string `json:"product_name_cn"`
		ProductVendor string `json:"product_vendor"`
		Version       string `json:"version"`
	} `json:"components"`
}

type APIResponse struct {
	Code    interface{}     `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
	Meta    json.RawMessage `json:"meta"`
}

type ServiceInfo struct {
	Code    interface{}     `json:"code"`
	Message string          `json:"message"`
	Data    ServiceDataList `json:"data"`
	Meta    struct {
		Pagination struct {
			Count     int `json:"count"`
			PageIndex int `json:"page_index"`
			PageSize  int `json:"page_size"`
			Total     int `json:"total"`
		} `json:"pagination"`
		PaginationID string `json:"pagination_id"`
	} `json:"meta"`
}
