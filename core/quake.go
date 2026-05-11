package core

import (
	"QuakeAPI/log"
	"QuakeAPI/model"
	"QuakeAPI/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

type QuakeInterface interface {
	GetUserInfo(key string)
	GetServiceInfo(key string, query string, total int, start int) (int, string)
}

type Core struct {
}

var httpClient utils.HttpClient

const quakeAPIBaseURL = "https://quake.360.cn"

func init() {
	httpClient = utils.HttpClient{}
}

func (c Core) GetUserInfo(key string) {
	url := quakeAPIBaseURL + "/api/v3/user/info"
	data := make(map[string]string)
	headers := make(map[string]string)
	headers["X-QuakeToken"] = key
	headers["Content-Type"] = "application/json"
	res := httpClient.DoGet(url, data, headers)
	var userInfo model.UserInfo
	err := json.Unmarshal(res, &userInfo)
	if err != nil {
		log.Log("unmarshal error:"+err.Error(), log.ERROR)
		return
	}
	if fmt.Sprintf("%v", userInfo.Code) != "0" {
		log.Log("Error API Key", log.ERROR)
		return
	}
	log.Log("Connect Success", log.INFO)
	log.Log("Your Name Is "+userInfo.Data.User.Username, log.INFO)
	log.Log("Your Email Is "+userInfo.Data.User.Email, log.INFO)
	log.Log("Your Phone Is "+userInfo.Data.MobilePhone, log.INFO)
	var roles bytes.Buffer
	for _, role := range userInfo.Data.Role {
		roles.WriteString(role.Fullname + " ")
	}
	log.Log("Your Role Is "+roles.String(), log.INFO)
}

func (c Core) GetServiceInfo(key string, query string, total int, start int) (int, string) {
	url := quakeAPIBaseURL + "/api/v3/search/quake_service"

	bodyMap := make(map[string]interface{})
	bodyMap["query"] = query
	bodyMap["start"] = start
	bodyMap["size"] = total
	bodyMap["ignore_cache"] = false
	bodyMap["latest"] = true
	bodyMap["include"] = []string{
		"ip",
		"port",
		"hostname",
		"transport",
		"asn",
		"org",
		"service.name",
		"location.country_cn",
		"location.province_cn",
		"location.city_cn",
		"service.http.host",
		"service.http.title",
		"service.http.server",
	}

	jsonBody, err := json.Marshal(bodyMap)
	if err != nil {
		log.Log("marshal error:"+err.Error(), log.ERROR)
		return 0, ""
	}
	data := make(map[string]string)
	data["_raw_body"] = string(jsonBody)

	headers := make(map[string]string)
	headers["X-QuakeToken"] = key
	headers["Content-Type"] = "application/json"

	log.Log(fmt.Sprintf("Requesting Data (start: %d, size: %d)......", start, total), log.INFO)
	res := httpClient.DoPost(url, data, headers)
	if len(res) == 0 {
		log.Log("empty response", log.ERROR)
		return 0, ""
	}

	var apiResponse model.APIResponse
	if err = json.Unmarshal(res, &apiResponse); err != nil {
		log.Log("raw response: "+string(res), log.ERROR)
		log.Log("unmarshal error:"+err.Error(), log.ERROR)
		return 0, ""
	}

	if fmt.Sprintf("%v", apiResponse.Code) != "0" {
		log.Log(fmt.Sprintf("API Error: %v - %s", apiResponse.Code, apiResponse.Message), log.ERROR)
		return 0, ""
	}

	var serviceInfo model.ServiceInfo
	err = json.Unmarshal(res, &serviceInfo)
	if err != nil {
		log.Log("raw response: "+string(res), log.ERROR)
		log.Log("unmarshal error:"+err.Error(), log.ERROR)
		return 0, ""
	}

	result := bytes.Buffer{}
	for _, value := range serviceInfo.Data {
		output := formatServiceLine(value)
		result.WriteString(output + "\n")
	}

	return len(serviceInfo.Data), result.String()
}

func formatServiceLine(value model.ServiceData) string {
	port := strconv.Itoa(value.Port)
	serviceName := value.Service.Name
	if serviceName == "" {
		serviceName = "-"
	}

	domain := firstNonEmpty(value.Domain, value.Service.HTTP.Host, value.Hostname, "-")
	title := firstNonEmpty(value.Service.HTTP.Title, "-")
	location := joinNonEmpty("/", value.Location.CountryCn, value.Location.ProvinceCn, value.Location.CityCn)
	if location == "" {
		location = "-"
	}

	return fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s", value.IP, port, serviceName, domain, title, location)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func joinNonEmpty(separator string, values ...string) string {
	var result bytes.Buffer
	for _, value := range values {
		if value == "" {
			continue
		}
		if result.Len() > 0 {
			result.WriteString(separator)
		}
		result.WriteString(value)
	}
	return result.String()
}
