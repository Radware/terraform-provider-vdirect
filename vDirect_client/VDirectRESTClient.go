package vdirect

/*
 Code was generated at Sun Jan 20 15:23:37 IST 2019
 vDirect version: 4.6.0
*/
import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const _STATUS string = "status"
const _SUCCESS string = "success"
const _COMPLETE string = "complete"
const _URI string = "uri"

func convertToString(val interface{}) string {
	t := reflect.TypeOf(val).String()
	switch t {
	case "string":
		return val.(string)
	case "int32":
		return strconv.Itoa(val.(int))
	case "int":
		return strconv.Itoa(val.(int))
	default:
		return val.(string)
	}
}

func mapToQuery(args map[string]interface{}) string {
	sb := strings.Builder{}
	if args != nil {
		sb.WriteString("?")
		for key, value := range args {
			if value != nil && reflect.TypeOf(value).Name() == "string" && len(value.(string)) > 0 {
				sb.WriteString(key)
				sb.WriteString("=")
				sb.WriteString(convertToString(value))
				sb.WriteString("&")
			}
		}
	}
	result := sb.String()
	var suffix = ""
	if len(result) == 1 {
		suffix = "?"
	} else {
		suffix = "&"
	}
	result = strings.TrimSuffix(result, suffix)
	return result
}

type M map[string]interface{}

func unmarshal(statusCode int, reason string, payload string, bytes []byte, data interface{}) VDirectClientResponse {
	if err := json.Unmarshal(bytes, &data); err != nil {
		return VDirectClientResponse{statusCode, reason, err.Error(), nil}
	}
	return VDirectClientResponse{statusCode, reason, payload, data}
}

func getRequestBody(payload interface{}, notJson bool) io.Reader {
	if payload != nil {
		if notJson {
			return bytes.NewBuffer(payload.([]byte))
		} else {
			jsonStr, _ := json.Marshal(payload)
			return bytes.NewBuffer(jsonStr)
		}
	} else {
		return nil
	}
}

func (v VDirectClientResponse) ToString() string {
	return fmt.Sprintf("statusCode: %d, reason: %s dataStr: %s, data: %s", v.StatusCode, v.Reason, v.DataStr, v.Data)
}

func responseToMap(response *http.Response) *M {
	var data M
	contents, _ := ioutil.ReadAll(response.Body)
	contentsStr := string(contents)
	bytes := []byte(contentsStr)
	json.Unmarshal(bytes, &data)
	return &data
}

func handleAsyncResponse(response *http.Response, err error, notJson bool, client *http.Client) VDirectClientResponse {
	var asyncResponse = *responseToMap(response)
	var uri = asyncResponse[_URI].(string)
	for i := 0; i < asyncOperationTimeOut; i++ {
		r, _ := client.Get(uri)
		var statusResponse = *responseToMap(r)
		if statusResponse[_COMPLETE].(bool) {
			if statusResponse[_SUCCESS].(bool) {
				return VDirectClientResponse{r.StatusCode, r.Status, "", statusResponse}
			} else {
				return VDirectClientResponse{-1, response.Status, "", statusResponse}
			}
		} else {
			time.Sleep(time.Duration(1) * time.Second)
		}
	}
	return VDirectClientResponse{-1, "timeout", "timeout", nil}
}

func handleResponse(response *http.Response, err error, notJson bool) VDirectClientResponse {
	if err != nil {
		return VDirectClientResponse{-1, err.Error(), "", nil}
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return VDirectClientResponse{-1, err.Error(), "", nil}
		} else {
			if notJson || response.StatusCode == 204 {
				return VDirectClientResponse{response.StatusCode, response.Status, "", contents}
			} else {
				contentsStr := string(contents)
				bytes := []byte(contentsStr)
				if strings.HasPrefix(contentsStr, "[") {
					// list of maps
					var data []M
					return unmarshal(response.StatusCode, response.Status, contentsStr, bytes, data)
				} else {
					// a map
					var data M
					return unmarshal(response.StatusCode, response.Status, contentsStr, bytes, data)
				}
			}
		}
	}
}

var HTTP_VERBS = []string{"GET", "DELETE", "PUT", "POST"}

func isValidHTTPVerb(verb string) bool {
	for _, v := range HTTP_VERBS {
		if v == verb {
			return true
		}
	}
	return false
}

/**
  Generic http call
*/
func call(details *connectionDetails,
	client *http.Client,
	verb string,
	urlPostfix string,
	payload interface{},
	headers map[string]string,
	async bool,
	notJson bool) VDirectClientResponse {
	//var url = fmt.Sprintf("https://%s:2189/api/%s", details.address, urlPostfix)
	var url = fmt.Sprintf("http://%s:2188/api/%s", details.address, urlPostfix)
	upperVerb := strings.ToUpper(verb)
	if isValidHTTPVerb(upperVerb) {
		req, err := http.NewRequest(upperVerb, url, getRequestBody(payload, notJson))
		if err != nil {
			return VDirectClientResponse{-1, err.Error(), "", nil}
		}
		if headers != nil {
			for key, value := range headers {
				req.Header.Set(key, value)
			}
		}
		req.SetBasicAuth(details.user, details.password)
		response, err := client.Do(req)
		if err != nil {
			return VDirectClientResponse{-1, err.Error(), "", nil}
		}
		if response.StatusCode == 202 && waitForAsyncOperation && async {
			return handleAsyncResponse(response, err, notJson, client)
		} else {
			return handleResponse(response, err, notJson)
		}

	} else {
		panic(fmt.Sprintf("HTTP verb %s is not supported", verb))
	}
}

// private
type connectionDetails struct {
	address  string
	user     string
	password string
}

type VDirectClientResponse struct {
	// http status code  (200)
	StatusCode int
	// http reason string ("OK")
	Reason string
	// returned data as string
	DataStr string
	// returned data as [map | list of maps] (most of the times)
	Data interface{}
}

type ClientConfig struct {
	// wait for async operations [True]
	WaitForAsyncOperation bool
	// How many seconds to wait for async operation [60]
	AsyncOperationTimeOut int
	// SSL context verification [True]
	Verify bool
	// HTTP timeout in seconds [60]
	HTTPTimeOut int
}

func NewClientConfig(waitForAsyncOperation bool, asyncOperationTimeOut int, verify bool, HTTPTimeOut int) *ClientConfig {
	cc := new(ClientConfig)
	cc.WaitForAsyncOperation = waitForAsyncOperation
	cc.AsyncOperationTimeOut = asyncOperationTimeOut
	cc.Verify = verify
	cc.HTTPTimeOut = HTTPTimeOut
	return cc
}

/* -------------------- */
/* Generated code below */
/* -------------------- */
type ADC struct {
	details *connectionDetails
	client  *http.Client
}

func (c *ADC) RunTemplate(data interface{}, adcName string, template string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["template"] = template

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.template-parameters+json"
	finalPath := fmt.Sprintf("adc/%s/config/", adcName) + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *ADC) UpdateGroups(data interface{}, adcName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	finalPath := fmt.Sprintf("adc/%s/config/", adcName)
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *ADC) Control1(adcName string, action string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["action"] = action

	var headers map[string]string
	finalPath := fmt.Sprintf("adc/%s/config/", adcName) + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, nil, headers, false, false)

}
func (c *ADC) Update(data interface{}, adcName string, configureDevice bool) VDirectClientResponse {
	args := make(map[string]interface{})
	args["configureDevice"] = configureDevice

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.adc+json"
	finalPath := fmt.Sprintf("adc/%s/", adcName) + mapToQuery(args)
	return call(c.details, c.client, "PUT", finalPath, &data, headers, true, false)

}
func (c *ADC) Control2(adcName string, reboot string, action string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["reboot"] = reboot
	args["action"] = action

	var headers map[string]string
	finalPath := fmt.Sprintf("adc/%s/", adcName) + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, nil, headers, false, false)

}
func (c *ADC) List(name string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["name"] = name

	var headers map[string]string
	finalPath := "adc/" + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *ADC) Delete(adcName string, action string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["action"] = action

	var headers map[string]string
	finalPath := fmt.Sprintf("adc/%s/?action=unregister", adcName) // + mapToQuery(args)
	return call(c.details, c.client, "DELETE", finalPath, nil, headers, true, false)

}
func (c *ADC) GetConfigurationLastCaptured(adcName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("adc/%s/configLastCaptured/", adcName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *ADC) GetConfiguration(adcName string, diff string, q string, prop string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["diff"] = diff
	args["q"] = q
	args["prop"] = prop

	var headers map[string]string
	finalPath := fmt.Sprintf("adc/%s/config/", adcName) + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *ADC) Get(adcName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("adc/%s/", adcName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *ADC) ControlDevice(adcName string, action string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["action"] = action

	var headers map[string]string
	finalPath := fmt.Sprintf("adc/%s/device/", adcName) + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, nil, headers, false, false)

}
func (c *ADC) Create(data interface{}, validate string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["validate"] = validate

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.adc+json"
	finalPath := "adc/" //+ mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, &data, headers, true, false)

}
func (c *ADC) UpdateConfiguration(data interface{}, adcName string, property string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	finalPath := fmt.Sprintf("adc/%s/config/?prop=%s", adcName, property)
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}

type AppWall struct {
	details *connectionDetails
	client  *http.Client
}

func (c *AppWall) Get(appWallName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("appWall/%s/", appWallName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *AppWall) ControlDevice(appWallName string, action string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["action"] = action

	var headers map[string]string
	finalPath := fmt.Sprintf("appWall/%s/device/", appWallName) + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, nil, headers, false, false)

}
func (c *AppWall) Create(data interface{}, validate string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["validate"] = validate

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.appwall+json"
	finalPath := "appWall/" + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *AppWall) Update(data interface{}, appWallName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.appwall+json"
	finalPath := fmt.Sprintf("appWall/%s/", appWallName)
	return call(c.details, c.client, "PUT", finalPath, &data, headers, false, false)

}
func (c *AppWall) List() VDirectClientResponse {

	var headers map[string]string
	finalPath := "appWall/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *AppWall) Delete(appWallName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("appWall/%s/", appWallName)
	return call(c.details, c.client, "DELETE", finalPath, nil, headers, false, false)

}

type Backup struct {
	details *connectionDetails
	client  *http.Client
}

func (c *Backup) DeleteBackup(name string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("backup/%s/", name)
	return call(c.details, c.client, "DELETE", finalPath, nil, headers, false, false)

}
func (c *Backup) Restore(name string, target string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["target"] = target

	var headers map[string]string
	finalPath := fmt.Sprintf("backup/%s/", name) + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, nil, headers, false, false)

}
func (c *Backup) Upload(data interface{}) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-zip-compressed"
	finalPath := "backup/"
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, true)

}
func (c *Backup) CleanOrCreate(comment string, target string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["comment"] = comment
	args["target"] = target

	var headers map[string]string
	finalPath := "backup/" + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, nil, headers, false, false)

}
func (c *Backup) GetArchive(name string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("backup/%s/archive/", name)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Backup) List() VDirectClientResponse {

	var headers map[string]string
	finalPath := "backup/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Backup) ListTargets() VDirectClientResponse {

	var headers map[string]string
	finalPath := "backup/targets/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Backup) GetBackup(name string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("backup/%s/", name)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}

type Catalog struct {
	details *connectionDetails
	client  *http.Client
}

func (c *Catalog) GetCatalogItemInstances(instanceType string, name string, _type string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("catalog/%s/%s/%s/", _type, name, instanceType)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Catalog) GetCatalogItem(name string, _type string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("catalog/%s/%s/", _type, name)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Catalog) GetCatalog() VDirectClientResponse {

	var headers map[string]string
	finalPath := "catalog/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}

type Container struct {
	details *connectionDetails
	client  *http.Client
}

func (c *Container) Get(containerName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("container/%s/", containerName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Container) Create0(data interface{}, validate string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["validate"] = validate

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.container+json"
	finalPath := "container/" // + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *Container) Update(data interface{}, containerName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.container+json"
	finalPath := fmt.Sprintf("container/%s/", containerName)
	return call(c.details, c.client, "PUT", finalPath, &data, headers, false, false)

}
func (c *Container) Control(containerName string, action string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["action"] = action

	var headers map[string]string
	finalPath := fmt.Sprintf("container/%s/device/", containerName) + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, nil, headers, false, false)

}
func (c *Container) Create1(data interface{}, containerName string, action string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["action"] = action

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.adc+json"
	finalPath := fmt.Sprintf("container/%s/", containerName) + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, &data, headers, true, false)

}
func (c *Container) List(name string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["name"] = name

	var headers map[string]string
	finalPath := "container/" + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Container) Delete(containerName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("container/%s/", containerName)
	return call(c.details, c.client, "DELETE", finalPath, nil, headers, false, false)

}
func (c *Container) GetSupportedVersions(containerName string, name string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["name"] = name

	var headers map[string]string
	finalPath := fmt.Sprintf("container/%s/adcVersion/", containerName) + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Container) ListAdcs(containerName string, name string, includeRegistered bool, includeUnregistered bool, includeMissing bool) VDirectClientResponse {
	args := make(map[string]interface{})
	args["name"] = name
	args["includeRegistered"] = includeRegistered
	args["includeUnregistered"] = includeUnregistered
	args["includeMissing"] = includeMissing

	var headers map[string]string
	finalPath := fmt.Sprintf("container/%s/adc/", containerName) + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Container) GetCapacity(containerName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("container/%s/capacity/", containerName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}

type ContainerDriver struct {
	details *connectionDetails
	client  *http.Client
}

func (c *ContainerDriver) Get(name string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("containerDriver/%s/", name)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *ContainerDriver) List() VDirectClientResponse {

	var headers map[string]string
	finalPath := "containerDriver/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *ContainerDriver) ListParameters(name string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("containerDriver/%s/parameters/", name)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}

type ContainerPool struct {
	details *connectionDetails
	client  *http.Client
}

func (c *ContainerPool) Get(containerPoolName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("resource/containerPool/%s/", containerPoolName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *ContainerPool) Create(data interface{}) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.container-resource-pool+json"
	finalPath := "resource/containerPool/"
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *ContainerPool) Update(data interface{}, containerPoolName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.container-resource-pool+json"
	finalPath := fmt.Sprintf("resource/containerPool/%s/", containerPoolName)
	return call(c.details, c.client, "PUT", finalPath, &data, headers, false, false)

}
func (c *ContainerPool) List(name string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["name"] = name

	var headers map[string]string
	finalPath := "resource/containerPool/" + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *ContainerPool) Delete(containerPoolName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("resource/containerPool/%s/", containerPoolName)
	return call(c.details, c.client, "DELETE", finalPath, nil, headers, false, false)

}
func (c *ContainerPool) GetCapacity(containerPoolName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("resource/containerPool/%s/capacity/", containerPoolName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}

type Credentials struct {
	details *connectionDetails
	client  *http.Client
}

func (c *Credentials) GetProtocols(includeStandard bool) VDirectClientResponse {
	args := make(map[string]interface{})
	args["includeStandard"] = includeStandard

	var headers map[string]string
	finalPath := "credentials/protocols/" + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Credentials) ListCredentials(service string, protocol string, host string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["service"] = service
	args["protocol"] = protocol
	args["host"] = host

	var headers map[string]string
	finalPath := "credentials/" + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Credentials) Update(data interface{}) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.credentials+json"
	finalPath := "credentials/"
	return call(c.details, c.client, "PUT", finalPath, &data, headers, false, false)

}
func (c *Credentials) Create(data interface{}) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.credentials+json"
	finalPath := "credentials/"
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *Credentials) GetServices(includeStandard bool) VDirectClientResponse {
	args := make(map[string]interface{})
	args["includeStandard"] = includeStandard

	var headers map[string]string
	finalPath := "credentials/services/" + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Credentials) Delete() VDirectClientResponse {

	var headers map[string]string
	finalPath := "credentials/"
	return call(c.details, c.client, "DELETE", finalPath, nil, headers, false, false)

}

type DefensePro struct {
	details *connectionDetails
	client  *http.Client
}

func (c *DefensePro) Get(defenseProName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("defensePro/%s/", defenseProName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *DefensePro) ControlDevice(defenseProName string, action string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["action"] = action

	var headers map[string]string
	finalPath := fmt.Sprintf("defensePro/%s/device/", defenseProName) + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, nil, headers, false, false)

}
func (c *DefensePro) Update(data interface{}, defenseProName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.defensepro+json"
	finalPath := fmt.Sprintf("defensePro/%s/", defenseProName)
	return call(c.details, c.client, "PUT", finalPath, &data, headers, false, false)

}
func (c *DefensePro) Create(data interface{}, validate string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["validate"] = validate

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.defensepro+json"
	finalPath := "defensePro/" + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *DefensePro) List() VDirectClientResponse {

	var headers map[string]string
	finalPath := "defensePro/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *DefensePro) Delete(defenseProName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("defensePro/%s/", defenseProName)
	return call(c.details, c.client, "DELETE", finalPath, nil, headers, false, false)

}

type DeviceCollection struct {
	details *connectionDetails
	client  *http.Client
}

func (c *DeviceCollection) Expand(name string, deviceType string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["deviceType"] = deviceType

	var headers map[string]string
	finalPath := fmt.Sprintf("deviceCollection/%s/devices/", name) + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *DeviceCollection) Get(name string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("deviceCollection/%s/", name)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *DeviceCollection) Create(data interface{}) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	finalPath := "deviceCollection/"
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *DeviceCollection) Update(data interface{}, name string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	finalPath := fmt.Sprintf("deviceCollection/%s/", name)
	return call(c.details, c.client, "PUT", finalPath, &data, headers, false, false)

}
func (c *DeviceCollection) List() VDirectClientResponse {

	var headers map[string]string
	finalPath := "deviceCollection/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *DeviceCollection) Delete(name string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("deviceCollection/%s/", name)
	return call(c.details, c.client, "DELETE", finalPath, nil, headers, false, false)

}

type Events struct {
	details *connectionDetails
	client  *http.Client
}

func (c *Events) PostJSONEvent(data interface{}, eventType string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["eventType"] = eventType

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	finalPath := "events/" + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *Events) GetMessageQueue() VDirectClientResponse {

	var headers map[string]string
	finalPath := "events//sse/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Events) PostFormEvent(data interface{}, eventType string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["eventType"] = eventType

	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	finalPath := "events/" + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *Events) Get() VDirectClientResponse {

	var headers map[string]string
	finalPath := "events//sseindex/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Events) PostEvent() VDirectClientResponse {

	var headers map[string]string
	finalPath := "events/"
	return call(c.details, c.client, "POST", finalPath, nil, headers, false, false)

}

type HA struct {
	details *connectionDetails
	client  *http.Client
}

func (c *HA) Sleep(brb int64) VDirectClientResponse {
	args := make(map[string]interface{})
	args["brb"] = brb

	var headers map[string]string
	finalPath := "ha/" + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, nil, headers, false, false)

}
func (c *HA) SetHaConfig(data interface{}) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.ha-configuration+json"
	finalPath := "ha/config/"
	return call(c.details, c.client, "PUT", finalPath, &data, headers, false, false)

}
func (c *HA) GetHaStatus() VDirectClientResponse {

	var headers map[string]string
	finalPath := "ha/status/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *HA) Recover() VDirectClientResponse {

	var headers map[string]string
	finalPath := "ha/recover/"
	return call(c.details, c.client, "POST", finalPath, nil, headers, false, false)

}
func (c *HA) GetHaConfig() VDirectClientResponse {

	var headers map[string]string
	finalPath := "ha/config/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *HA) IsActive() VDirectClientResponse {

	var headers map[string]string
	finalPath := "ha/active/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *HA) ChangeHaStatus(data interface{}) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.ha-status+json"
	finalPath := "ha/"
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}

type IpAddress struct {
	details *connectionDetails
	client  *http.Client
}

func (c *IpAddress) Release(ipAddressName string, resource string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["resource"] = resource

	var headers map[string]string
	finalPath := fmt.Sprintf("resource/ipAddress/%s/pool/", ipAddressName) + mapToQuery(args)
	return call(c.details, c.client, "DELETE", finalPath, nil, headers, false, false)

}
func (c *IpAddress) Get(ipAddressName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("resource/ipAddress/%s/", ipAddressName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *IpAddress) List2(ipAddressName string, resource string, owner string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["resource"] = resource
	args["owner"] = owner

	var headers map[string]string
	finalPath := fmt.Sprintf("resource/ipAddress/%s/pool/", ipAddressName) + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *IpAddress) Create5(name string, start string, end string, gateway string, mask string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["name"] = name
	args["start"] = start
	args["end"] = end
	args["gateway"] = gateway
	args["mask"] = mask

	var headers map[string]string
	finalPath := "resource/ipAddress/" + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, nil, headers, false, false)

}
func (c *IpAddress) AcquireFromFormData(data interface{}, ipAddressName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	finalPath := fmt.Sprintf("resource/ipAddress/%s/", ipAddressName)
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *IpAddress) Update(data interface{}, ipAddressName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.ip-pool+json"
	finalPath := fmt.Sprintf("resource/ipAddress/%s/", ipAddressName)
	return call(c.details, c.client, "PUT", finalPath, &data, headers, false, false)

}
func (c *IpAddress) Create1(data interface{}, name string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["name"] = name

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.ip-pool+json"
	finalPath := "resource/ipAddress/" + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *IpAddress) List3(name string, resource string, owner string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["name"] = name
	args["resource"] = resource
	args["owner"] = owner

	var headers map[string]string
	finalPath := "resource/ipAddress/" + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *IpAddress) Acquire0(data interface{}, ipAddressName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.resource+json"
	finalPath := fmt.Sprintf("resource/ipAddress/%s/", ipAddressName)
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *IpAddress) Delete(ipAddressName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("resource/ipAddress/%s/", ipAddressName)
	return call(c.details, c.client, "DELETE", finalPath, nil, headers, false, false)

}
func (c *IpAddress) Acquire4(ipAddressName string, comment string, owner string, reserve bool, resource string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["comment"] = comment
	args["owner"] = owner
	args["reserve"] = reserve
	args["resource"] = resource

	var headers map[string]string
	finalPath := fmt.Sprintf("resource/ipAddress/%s/", ipAddressName) + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, nil, headers, false, false)

}

type ISL struct {
	details *connectionDetails
	client  *http.Client
}

func (c *ISL) Release(islName string, resource string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["resource"] = resource

	var headers map[string]string
	finalPath := fmt.Sprintf("resource/isl/%s/pool/", islName) + mapToQuery(args)
	return call(c.details, c.client, "DELETE", finalPath, nil, headers, false, false)

}
func (c *ISL) List3(name string, resource string, owner string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["name"] = name
	args["resource"] = resource
	args["owner"] = owner

	var headers map[string]string
	finalPath := "resource/isl/" + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *ISL) Get(islName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("resource/isl/%s/", islName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *ISL) AcquireFromFormData(data interface{}, islName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	finalPath := fmt.Sprintf("resource/isl/%s/", islName)
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *ISL) Update(data interface{}, islName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.isl-pool+json"
	finalPath := fmt.Sprintf("resource/isl/%s/", islName)
	return call(c.details, c.client, "PUT", finalPath, &data, headers, false, false)

}
func (c *ISL) Create1(data interface{}, name string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["name"] = name

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.isl-pool+json"
	finalPath := "resource/isl/" + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *ISL) Create3(name string, min string, max string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["name"] = name
	args["min"] = min
	args["max"] = max

	var headers map[string]string
	finalPath := "resource/isl/" + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, nil, headers, false, false)

}
func (c *ISL) Acquire0(data interface{}, islName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.resource+json"
	finalPath := fmt.Sprintf("resource/isl/%s/", islName)
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *ISL) List2(islName string, resource string, owner string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["resource"] = resource
	args["owner"] = owner

	var headers map[string]string
	finalPath := fmt.Sprintf("resource/isl/%s/pool/", islName) + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *ISL) Acquire4(islName string, comment string, owner string, reserve bool, resource string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["comment"] = comment
	args["owner"] = owner
	args["reserve"] = reserve
	args["resource"] = resource

	var headers map[string]string
	finalPath := fmt.Sprintf("resource/isl/%s/", islName) + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, nil, headers, false, false)

}
func (c *ISL) Delete(islName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("resource/isl/%s/", islName)
	return call(c.details, c.client, "DELETE", finalPath, nil, headers, false, false)

}

type ManagedObject struct {
	details *connectionDetails
	client  *http.Client
}

func (c *ManagedObject) GetObject(name string, _type string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("managed-object/%s/%s/", _type, name)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *ManagedObject) Get(_type string, name string, related bool, id string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["type"] = _type
	args["name"] = name
	args["related"] = related
	args["id"] = id

	var headers map[string]string
	finalPath := "managed-object/" + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *ManagedObject) GetId(_type string, name string, id string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["type"] = _type
	args["name"] = name
	args["id"] = id

	var headers map[string]string
	finalPath := "managed-object/id/" + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *ManagedObject) ListObjects(_type string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("managed-object/%s/", _type)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}

type Message struct {
	details *connectionDetails
	client  *http.Client
}

func (c *Message) Get(messageId string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("message/%s/", messageId)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Message) GetMessageEntity(messageId string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("message/%s/entity/", messageId)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Message) Delete(messageId string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("message/%s/", messageId)
	return call(c.details, c.client, "DELETE", finalPath, nil, headers, false, false)

}
func (c *Message) GetMessages() VDirectClientResponse {

	var headers map[string]string
	finalPath := "message/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}

type Network struct {
	details *connectionDetails
	client  *http.Client
}

func (c *Network) Get(networkName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("resource/network/%s/", networkName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Network) Create(data interface{}, name string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["name"] = name

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.network+json"
	finalPath := "resource/network/" + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *Network) Update(data interface{}, networkName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.network+json"
	finalPath := fmt.Sprintf("resource/network/%s/", networkName)
	return call(c.details, c.client, "PUT", finalPath, &data, headers, false, false)

}
func (c *Network) List(name string, vlan string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["name"] = name
	args["vlan"] = vlan

	var headers map[string]string
	finalPath := "resource/network/" + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Network) Delete(networkName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("resource/network/%s/", networkName)
	return call(c.details, c.client, "DELETE", finalPath, nil, headers, false, false)

}

type Oper struct {
	details *connectionDetails
	client  *http.Client
}

func (c *Oper) GetListeners() VDirectClientResponse {

	var headers map[string]string
	finalPath := "oper/syslog/listener/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Oper) ListSessions() VDirectClientResponse {

	var headers map[string]string
	finalPath := "oper/sessions/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Oper) GetConverters() VDirectClientResponse {

	var headers map[string]string
	finalPath := "oper/syslog/converter/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Oper) AuthConfig() VDirectClientResponse {

	var headers map[string]string
	finalPath := "oper/auth/"
	return call(c.details, c.client, "POST", finalPath, nil, headers, false, false)

}
func (c *Oper) DownloadApplicationLogs() VDirectClientResponse {

	var headers map[string]string
	finalPath := "oper/logs/server/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Oper) GetConfig0() VDirectClientResponse {

	var headers map[string]string
	finalPath := "oper/config/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Oper) GetServerLogLevel() VDirectClientResponse {

	var headers map[string]string
	finalPath := "oper/logs/server/level/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Oper) SetLogLevelFromFormData(data interface{}) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	finalPath := "oper/logs/level/"
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *Oper) RequestRoute(data interface{}) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	finalPath := "oper/proxy/"
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *Oper) ListLocks() VDirectClientResponse {

	var headers map[string]string
	finalPath := "oper/locks/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Oper) RestartServer() VDirectClientResponse {

	var headers map[string]string
	finalPath := "oper/reset/"
	return call(c.details, c.client, "POST", finalPath, nil, headers, false, false)

}
func (c *Oper) GetServicesTable() VDirectClientResponse {

	var headers map[string]string
	finalPath := "oper/service/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Oper) GetLastMessage() VDirectClientResponse {

	var headers map[string]string
	finalPath := "oper/message/last/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Oper) ProxyAlive() VDirectClientResponse {

	var headers map[string]string
	finalPath := "oper/proxy/ping/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Oper) GetWorkflowsTable() VDirectClientResponse {

	var headers map[string]string
	finalPath := "oper/workflow/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Oper) GetLogLevel() VDirectClientResponse {

	var headers map[string]string
	finalPath := "oper/logs/level/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Oper) GetJsonServerLogs() VDirectClientResponse {

	var headers map[string]string
	finalPath := "oper/logs/server/preview/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Oper) PutConfig(data interface{}) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	finalPath := "oper/proxy/config/"
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *Oper) PutListeners(data interface{}) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.syslog+json"
	finalPath := "oper/syslog/listener/"
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *Oper) Test(data interface{}, action string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["action"] = action

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.syslog+json"
	finalPath := "oper/syslog/" + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *Oper) SetLogLevel(level string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["level"] = level

	var headers map[string]string
	finalPath := "oper/logs/level/" + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, nil, headers, false, false)

}
func (c *Oper) GetLogs() VDirectClientResponse {

	var headers map[string]string
	finalPath := "oper/logs/preview/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Oper) Control(lock string, owner string, action string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["lock"] = lock
	args["owner"] = owner
	args["action"] = action

	var headers map[string]string
	finalPath := "oper/locks/" + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, nil, headers, false, false)

}
func (c *Oper) EventListenerConfig() VDirectClientResponse {

	var headers map[string]string
	finalPath := "oper/events/"
	return call(c.details, c.client, "POST", finalPath, nil, headers, false, false)

}
func (c *Oper) SetServerLogLevel(level string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["level"] = level

	var headers map[string]string
	finalPath := "oper/logs/server/level/" + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, nil, headers, false, false)

}
func (c *Oper) List(_type string, format string, deviceType string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["type"] = _type
	args["format"] = format
	args["deviceType"] = deviceType

	var headers map[string]string
	finalPath := "oper/inventory/" + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Oper) GetConfig1(key string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["key"] = key

	var headers map[string]string
	finalPath := "oper/proxy/config/" + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Oper) GetServerLogs() VDirectClientResponse {

	var headers map[string]string
	finalPath := "oper/logs/server/preview/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Oper) SaveConfig(data interface{}) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	finalPath := "oper/config/"
	return call(c.details, c.client, "PUT", finalPath, &data, headers, false, false)

}
func (c *Oper) GetJsonLogs() VDirectClientResponse {

	var headers map[string]string
	finalPath := "oper/logs/preview/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Oper) DownloadAuditLogs() VDirectClientResponse {

	var headers map[string]string
	finalPath := "oper/logs/audit/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Oper) GetAdcLogs() VDirectClientResponse {

	var headers map[string]string
	finalPath := "oper/logs/adc/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Oper) SetServerLogLevelFromFormData(data interface{}) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	finalPath := "oper/logs/serverlevel/"
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *Oper) ListDeviceLocks() VDirectClientResponse {

	var headers map[string]string
	finalPath := "oper/locks/device/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Oper) PutConverters(data interface{}) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.syslog+json"
	finalPath := "oper/syslog/converter/"
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *Oper) GetLicensedCapacity() VDirectClientResponse {

	var headers map[string]string
	finalPath := "oper/adcLicense/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}

type RBAC struct {
	details *connectionDetails
	client  *http.Client
}

func (c *RBAC) ListUsers() VDirectClientResponse {

	var headers map[string]string
	finalPath := "rbac/user/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *RBAC) Get1(allowedParentOf string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["allowedParentOf"] = allowedParentOf

	var headers map[string]string
	finalPath := "rbac/role/" + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *RBAC) ListGroups() VDirectClientResponse {

	var headers map[string]string
	finalPath := "rbac/group/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *RBAC) Create(data interface{}) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.role+json"
	finalPath := "rbac/role/"
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *RBAC) Update(data interface{}, roleName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.role+json"
	finalPath := fmt.Sprintf("rbac/role/%s/", roleName)
	return call(c.details, c.client, "PUT", finalPath, &data, headers, false, false)

}
func (c *RBAC) ListPermissions(role string, permission string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["role"] = role
	args["permission"] = permission

	var headers map[string]string
	finalPath := "rbac/permission/" + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *RBAC) Get0(roleName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("rbac/role/%s/", roleName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *RBAC) Delete(roleName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("rbac/role/%s/", roleName)
	return call(c.details, c.client, "DELETE", finalPath, nil, headers, false, false)

}

type Runnable struct {
	details *connectionDetails
	client  *http.Client
}

func (c *Runnable) GetAvailableActions(name string, _type string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("runnable/%s/%s/", _type, name)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Runnable) GetActionInfo(name string, _type string, action string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("runnable/%s/%s/%s/", _type, name, action)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Runnable) Run(data interface{}, name string, _type string, action string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	finalPath := fmt.Sprintf("runnable/%s/%s/%s/", _type, name, action)
	return call(c.details, c.client, "POST", finalPath, &data, headers, true, false)

}
func (c *Runnable) GetTypes() VDirectClientResponse {

	var headers map[string]string
	finalPath := "runnable/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Runnable) GetRunnableObjects(_type string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("runnable/%s/", _type)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Runnable) GetSubInstanceActionInfo(instanceName string, instanceType string, name string, _type string, action string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("runnable/%s/%s/%s/%s/%s/", _type, name, instanceType, instanceName, action)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Runnable) RunSubInstanceAction(data interface{}, instanceName string, instanceType string, name string, _type string, action string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	finalPath := fmt.Sprintf("runnable/%s/%s/%s/%s/%s/", _type, name, instanceType, instanceName, action)
	return call(c.details, c.client, "POST", finalPath, &data, headers, true, false)

}
func (c *Runnable) GetCatalog() VDirectClientResponse {

	var headers map[string]string
	finalPath := "runnable/catalog/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}

type Scheduled struct {
	details *connectionDetails
	client  *http.Client
}

func (c *Scheduled) Get(scheduledName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("scheduled/%s/", scheduledName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Scheduled) Update(data interface{}, scheduledName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.scheduled-job+json"
	finalPath := fmt.Sprintf("scheduled/%s/", scheduledName)
	return call(c.details, c.client, "PUT", finalPath, &data, headers, false, false)

}
func (c *Scheduled) Create(data interface{}) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.scheduled-job+json"
	finalPath := "scheduled/"
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *Scheduled) Control(scheduledName string, action string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["action"] = action

	var headers map[string]string
	finalPath := fmt.Sprintf("scheduled/%s/", scheduledName) + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, nil, headers, false, false)

}
func (c *Scheduled) List() VDirectClientResponse {

	var headers map[string]string
	finalPath := "scheduled/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Scheduled) Delete(scheduledName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("scheduled/%s/", scheduledName)
	return call(c.details, c.client, "DELETE", finalPath, nil, headers, false, false)

}

type Service struct {
	details *connectionDetails
	client  *http.Client
}

func (c *Service) DeleteHistory(serviceName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("service/%s/history/", serviceName)
	return call(c.details, c.client, "DELETE", finalPath, nil, headers, false, false)

}
func (c *Service) GetHistory(serviceName string, format string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["format"] = format

	var headers map[string]string
	finalPath := fmt.Sprintf("service/%s/history/", serviceName) + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Service) CleanHistory(name string, tenant string, clean bool) VDirectClientResponse {
	args := make(map[string]interface{})
	args["name"] = name
	args["tenant"] = tenant
	args["clean"] = clean

	var headers map[string]string
	finalPath := "service/" + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, nil, headers, false, false)

}
func (c *Service) GetSpecification(serviceName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("service/%s/specification/", serviceName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Service) DeleteService(serviceName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("service/%s/", serviceName)
	return call(c.details, c.client, "DELETE", finalPath, nil, headers, true, false)

}
func (c *Service) RunAction(serviceName string, action string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["action"] = action

	var headers map[string]string
	finalPath := fmt.Sprintf("service/%s/", serviceName) + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, nil, headers, true, false)

}
func (c *Service) UpdateService(data interface{}, serviceName string, action string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["action"] = action

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.adc-service+json"
	finalPath := fmt.Sprintf("service/%s/", serviceName) + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, &data, headers, true, false)

}
func (c *Service) Get(serviceName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("service/%s/", serviceName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Service) Create(data interface{}, name string, tenant string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["name"] = name
	args["tenant"] = tenant

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.adc-service-specification+json"
	finalPath := "service/" + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *Service) List(name string, includeDeleted bool, deletedOnly bool, usingResourcePoolName string, usingResourcePoolId string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["name"] = name
	args["includeDeleted"] = includeDeleted
	args["deletedOnly"] = deletedOnly
	args["usingResourcePoolName"] = usingResourcePoolName
	args["usingResourcePoolId"] = usingResourcePoolId

	var headers map[string]string
	finalPath := "service/" + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Service) UpdateSpecification(data interface{}, serviceName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.adc-service-specification+json"
	finalPath := fmt.Sprintf("service/%s/specification/", serviceName)
	return call(c.details, c.client, "PUT", finalPath, &data, headers, true, false)

}
func (c *Service) FixService(data interface{}, serviceName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.adc-service+json"
	finalPath := fmt.Sprintf("service/%s/", serviceName)
	return call(c.details, c.client, "PUT", finalPath, &data, headers, true, false)

}

type Session struct {
	details *connectionDetails
	client  *http.Client
}

func (c *Session) Get(Cookie string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["Cookie"] = Cookie

	var headers map[string]string
	finalPath := "session/" + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Session) Update(data interface{}, Cookie string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["Cookie"] = Cookie

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.session+json"
	finalPath := "session/" + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *Session) Create(data interface{}, Cookie string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["Cookie"] = Cookie

	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	finalPath := "session/" + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *Session) Delete(Cookie string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["Cookie"] = Cookie

	var headers map[string]string
	finalPath := "session/" + mapToQuery(args)
	return call(c.details, c.client, "DELETE", finalPath, nil, headers, false, false)

}

type Status struct {
	details *connectionDetails
	client  *http.Client
}

func (c *Status) GetStats() VDirectClientResponse {

	var headers map[string]string
	finalPath := "status/stats/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Status) GetServerStats() VDirectClientResponse {

	var headers map[string]string
	finalPath := "status/serverStats/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Status) GetResult(token string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["token"] = token

	var headers map[string]string
	finalPath := "status/result/" + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Status) GetStatus(token string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["token"] = token

	var headers map[string]string
	finalPath := "status/" + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, true, false)

}

type Template struct {
	details *connectionDetails
	client  *http.Client
}

func (c *Template) RunTemplate(data interface{}, templateName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.template-parameters+json"
	finalPath := fmt.Sprintf("template/%s/", templateName)
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *Template) CreateFromFormData(data interface{}) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	finalPath := "template/"
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *Template) UploadSource(data interface{}, templateName string, failIfInvalid bool) VDirectClientResponse {
	args := make(map[string]interface{})
	args["failIfInvalid"] = failIfInvalid

	headers := make(map[string]string)
	headers["Content-Type"] = "text/x-velocity"
	finalPath := fmt.Sprintf("template/%s/source/", templateName) + mapToQuery(args)
	return call(c.details, c.client, "PUT", finalPath, &data, headers, false, true)

}
func (c *Template) UploadSourceFromFormData(data interface{}, templateName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	finalPath := fmt.Sprintf("template/%s/source/", templateName)
	return call(c.details, c.client, "PUT", finalPath, &data, headers, false, false)

}
func (c *Template) Get(templateName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("template/%s/", templateName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Template) Update(data interface{}, templateName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.template+json"
	finalPath := fmt.Sprintf("template/%s/", templateName)
	return call(c.details, c.client, "PUT", finalPath, &data, headers, false, false)

}
func (c *Template) DownloadSource(templateName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("template/%s/source/", templateName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Template) CreateFromSource(data interface{}, name string, tenant string, failIfInvalid bool) VDirectClientResponse {
	args := make(map[string]interface{})
	args["name"] = name
	args["tenant"] = tenant
	args["failIfInvalid"] = failIfInvalid

	headers := make(map[string]string)
	headers["Content-Type"] = "text/x-velocity"
	finalPath := "template/" + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, true)

}
func (c *Template) List(name string, display string, device string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["name"] = name
	args["display"] = display
	args["device"] = device

	var headers map[string]string
	finalPath := "template/" + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Template) Delete(templateName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("template/%s/", templateName)
	return call(c.details, c.client, "DELETE", finalPath, nil, headers, false, false)

}
func (c *Template) GetIcon(templateName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("template/%s/icon/", templateName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}

type Tenant struct {
	details *connectionDetails
	client  *http.Client
}

func (c *Tenant) Get0(include string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["include"] = include

	var headers map[string]string
	finalPath := "tenant/" + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Tenant) Create(data interface{}) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.tenant+json"
	finalPath := "tenant/"
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *Tenant) Update(data interface{}, tenantName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.tenant+json"
	finalPath := fmt.Sprintf("tenant/%s/", tenantName)
	return call(c.details, c.client, "PUT", finalPath, &data, headers, false, false)

}
func (c *Tenant) Get1(tenantName string, include string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["include"] = include

	var headers map[string]string
	finalPath := fmt.Sprintf("tenant/%s/", tenantName) + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Tenant) Delete(tenantName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("tenant/%s/", tenantName)
	return call(c.details, c.client, "DELETE", finalPath, nil, headers, false, false)

}

type Triggered struct {
	details *connectionDetails
	client  *http.Client
}

func (c *Triggered) Control1(data interface{}, triggeredName string, action string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["action"] = action

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	finalPath := fmt.Sprintf("triggered/%s/", triggeredName) + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *Triggered) Get(triggeredName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("triggered/%s/", triggeredName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Triggered) Create(data interface{}) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.triggered-job+json"
	finalPath := "triggered/"
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *Triggered) Update(data interface{}, triggeredName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.triggered-job+json"
	finalPath := fmt.Sprintf("triggered/%s/", triggeredName)
	return call(c.details, c.client, "PUT", finalPath, &data, headers, false, false)

}
func (c *Triggered) Control0(data interface{}, action string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["action"] = action

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	finalPath := "triggered/" + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *Triggered) Control(triggeredName string, action string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["action"] = action

	var headers map[string]string
	finalPath := fmt.Sprintf("triggered/%s/", triggeredName) + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, nil, headers, false, false)

}
func (c *Triggered) List() VDirectClientResponse {

	var headers map[string]string
	finalPath := "triggered/"
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Triggered) Delete(triggeredName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("triggered/%s/", triggeredName)
	return call(c.details, c.client, "DELETE", finalPath, nil, headers, false, false)

}

type VRRP struct {
	details *connectionDetails
	client  *http.Client
}

func (c *VRRP) Release(vrrpName string, resource string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["resource"] = resource

	var headers map[string]string
	finalPath := fmt.Sprintf("resource/vrrp/%s/pool/", vrrpName) + mapToQuery(args)
	return call(c.details, c.client, "DELETE", finalPath, nil, headers, false, false)

}
func (c *VRRP) List3(name string, resource string, owner string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["name"] = name
	args["resource"] = resource
	args["owner"] = owner

	var headers map[string]string
	finalPath := "resource/vrrp/" + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *VRRP) Get(vrrpName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("resource/vrrp/%s/", vrrpName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *VRRP) AcquireFromFormData(data interface{}, vrrpName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	finalPath := fmt.Sprintf("resource/vrrp/%s/", vrrpName)
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *VRRP) Update(data interface{}, vrrpName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.vrrp-pool+json"
	finalPath := fmt.Sprintf("resource/vrrp/%s/", vrrpName)
	return call(c.details, c.client, "PUT", finalPath, &data, headers, false, false)

}
func (c *VRRP) Create0(data interface{}) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.vrrp-pool+json"
	finalPath := "resource/vrrp/"
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *VRRP) Create1(name string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["name"] = name

	var headers map[string]string
	finalPath := "resource/vrrp/" + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, nil, headers, false, false)

}
func (c *VRRP) List2(vrrpName string, resource string, owner string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["resource"] = resource
	args["owner"] = owner

	var headers map[string]string
	finalPath := fmt.Sprintf("resource/vrrp/%s/pool/", vrrpName) + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *VRRP) Acquire0(data interface{}, vrrpName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.resource+json"
	finalPath := fmt.Sprintf("resource/vrrp/%s/", vrrpName)
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *VRRP) Delete(vrrpName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("resource/vrrp/%s/", vrrpName)
	return call(c.details, c.client, "DELETE", finalPath, nil, headers, false, false)

}
func (c *VRRP) Acquire4(vrrpName string, comment string, owner string, reserve bool, resource string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["comment"] = comment
	args["owner"] = owner
	args["reserve"] = reserve
	args["resource"] = resource

	var headers map[string]string
	finalPath := fmt.Sprintf("resource/vrrp/%s/", vrrpName) + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, nil, headers, false, false)

}

type Workflow struct {
	details *connectionDetails
	client  *http.Client
}

func (c *Workflow) UpdateWorkflow(data interface{}, workflowName string, actionName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.template-parameters+json"
	finalPath := fmt.Sprintf("workflow/%s/action/%s/", workflowName, actionName)
	return call(c.details, c.client, "POST", finalPath, &data, headers, true, false)

}
func (c *Workflow) DeleteHistory(workflowName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("workflow/%s/history/", workflowName)
	return call(c.details, c.client, "DELETE", finalPath, nil, headers, false, false)

}
func (c *Workflow) GetWorkflow(workflowName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("workflow/%s/", workflowName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Workflow) GetHistory(workflowName string, format string, actionId string, debug bool) VDirectClientResponse {
	args := make(map[string]interface{})
	args["format"] = format
	args["actionId"] = actionId
	args["debug"] = debug

	var headers map[string]string
	finalPath := fmt.Sprintf("workflow/%s/history/", workflowName) + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Workflow) CleanHistory(clean bool) VDirectClientResponse {
	args := make(map[string]interface{})
	args["clean"] = clean

	var headers map[string]string
	finalPath := "workflow/" + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, nil, headers, false, false)

}
func (c *Workflow) GetActionLog(workflowName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("workflow/%s/actionLog/", workflowName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Workflow) DeleteWorkflow(workflowName string, remove bool) VDirectClientResponse {
	args := make(map[string]interface{})
	args["remove"] = remove

	var headers map[string]string
	finalPath := fmt.Sprintf("workflow/%s/", workflowName) + mapToQuery(args)
	return call(c.details, c.client, "DELETE", finalPath, nil, headers, true, false)

}
func (c *Workflow) GetActionInfo(workflowName string, actionName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("workflow/%s/action/%s/", workflowName, actionName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Workflow) GetParameters(workflowName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("workflow/%s/parameters/", workflowName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *Workflow) List(name string, _type string, includeDeleted bool, deletedOnly bool, usingResourceName string, usingResourceId string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["name"] = name
	args["type"] = _type
	args["includeDeleted"] = includeDeleted
	args["deletedOnly"] = deletedOnly
	args["usingResourceName"] = usingResourceName
	args["usingResourceId"] = usingResourceId

	var headers map[string]string
	finalPath := "workflow/" + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}

type WorkflowTemplate struct {
	details *connectionDetails
	client  *http.Client
}

func (c *WorkflowTemplate) GetDescriptor(workflowTemplateName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("workflowTemplate/%s/descriptor/", workflowTemplateName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *WorkflowTemplate) CreateWorkflow(data interface{}, workflowTemplateName string, name string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["name"] = name

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.template-parameters+json"
	finalPath := fmt.Sprintf("workflowTemplate/%s/", workflowTemplateName) + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, &data, headers, true, false)

}
func (c *WorkflowTemplate) UpdateArchive(data interface{}, workflowTemplateName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-zip-compressed"
	finalPath := fmt.Sprintf("workflowTemplate/%s/archive/", workflowTemplateName)
	return call(c.details, c.client, "PUT", finalPath, &data, headers, false, true)

}
func (c *WorkflowTemplate) GetActionInfo(workflowTemplateName string, actionName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("workflowTemplate/%s/action/%s/", workflowTemplateName, actionName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *WorkflowTemplate) List(name string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["name"] = name

	var headers map[string]string
	finalPath := "workflowTemplate/" + mapToQuery(args)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *WorkflowTemplate) GetFile(fileName string, workflowTemplateName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("workflowTemplate/%s/file/%s/", workflowTemplateName, fileName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *WorkflowTemplate) UpdateWorkflowTemplate(data interface{}, workflowTemplateName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/vnd.com.radware.vdirect.workflow-template+json"
	finalPath := fmt.Sprintf("workflowTemplate/%s/", workflowTemplateName)
	return call(c.details, c.client, "PUT", finalPath, &data, headers, false, false)

}
func (c *WorkflowTemplate) CreateTemplate(data interface{}) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "multipart/form-data"
	finalPath := "workflowTemplate/"
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *WorkflowTemplate) GetIcon(workflowTemplateName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("workflowTemplate/%s/icon/", workflowTemplateName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *WorkflowTemplate) UpdateFile(data interface{}, fileName string, workflowTemplateName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "text/x-groovy"
	finalPath := fmt.Sprintf("workflowTemplate/%s/file/%s/", workflowTemplateName, fileName)
	return call(c.details, c.client, "PUT", finalPath, &data, headers, false, false)

}
func (c *WorkflowTemplate) DeleteWorkflowTemplate(workflowTemplateName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("workflowTemplate/%s/", workflowTemplateName)
	return call(c.details, c.client, "DELETE", finalPath, nil, headers, false, false)

}
func (c *WorkflowTemplate) CreateTemplateFromDescriptor(data interface{}, tenant string, failIfInvalid bool) VDirectClientResponse {
	args := make(map[string]interface{})
	args["tenant"] = tenant
	args["failIfInvalid"] = failIfInvalid

	headers := make(map[string]string)
	headers["Content-Type"] = "application/xml"
	finalPath := "workflowTemplate/" + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, true)

}
func (c *WorkflowTemplate) CreateTemplateFromArchive(data interface{}, validate bool, failIfInvalid bool, tenant string) VDirectClientResponse {
	args := make(map[string]interface{})
	args["validate"] = validate
	args["failIfInvalid"] = failIfInvalid
	args["tenant"] = tenant

	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-zip-compressed"
	finalPath := "workflowTemplate/" + mapToQuery(args)
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, true)

}
func (c *WorkflowTemplate) UpdateDescriptor(data interface{}, workflowTemplateName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/xml"
	finalPath := fmt.Sprintf("workflowTemplate/%s/descriptor/", workflowTemplateName)
	return call(c.details, c.client, "PUT", finalPath, &data, headers, false, true)

}
func (c *WorkflowTemplate) GetArchive(workflowTemplateName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("workflowTemplate/%s/archive/", workflowTemplateName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}
func (c *WorkflowTemplate) UpdateTemplate(data interface{}, workflowTemplateName string) VDirectClientResponse {

	headers := make(map[string]string)
	headers["Content-Type"] = "multipart/form-data"
	finalPath := fmt.Sprintf("workflowTemplate/%s/archive/", workflowTemplateName)
	return call(c.details, c.client, "POST", finalPath, &data, headers, false, false)

}
func (c *WorkflowTemplate) GetWorkflowTemplate(workflowTemplateName string) VDirectClientResponse {

	var headers map[string]string
	finalPath := fmt.Sprintf("workflowTemplate/%s/", workflowTemplateName)
	return call(c.details, c.client, "GET", finalPath, nil, headers, false, false)

}

type Client struct {
	ADC              ADC
	AppWall          AppWall
	Backup           Backup
	Catalog          Catalog
	Container        Container
	ContainerDriver  ContainerDriver
	ContainerPool    ContainerPool
	Credentials      Credentials
	DefensePro       DefensePro
	DeviceCollection DeviceCollection
	Events           Events
	HA               HA
	IpAddress        IpAddress
	ISL              ISL
	ManagedObject    ManagedObject
	Message          Message
	Network          Network
	Oper             Oper
	RBAC             RBAC
	Runnable         Runnable
	Scheduled        Scheduled
	Service          Service
	Session          Session
	Status           Status
	Template         Template
	Tenant           Tenant
	Triggered        Triggered
	VRRP             VRRP
	Workflow         Workflow
	WorkflowTemplate WorkflowTemplate
}

var waitForAsyncOperation = true
var asyncOperationTimeOut = 60

func NewClient(address string, user string, password string, config *ClientConfig) Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !config.Verify},
	}
	waitForAsyncOperation = config.WaitForAsyncOperation
	asyncOperationTimeOut = config.AsyncOperationTimeOut
	// see here for http timeouts https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	client := &http.Client{Timeout: time.Duration(config.HTTPTimeOut) * time.Second, Transport: tr}
	conn := connectionDetails{address, user, password}
	return Client{ADC{&conn, client}, AppWall{&conn, client}, Backup{&conn, client}, Catalog{&conn, client}, Container{&conn, client}, ContainerDriver{&conn, client}, ContainerPool{&conn, client}, Credentials{&conn, client}, DefensePro{&conn, client}, DeviceCollection{&conn, client}, Events{&conn, client}, HA{&conn, client}, IpAddress{&conn, client}, ISL{&conn, client}, ManagedObject{&conn, client}, Message{&conn, client}, Network{&conn, client}, Oper{&conn, client}, RBAC{&conn, client}, Runnable{&conn, client}, Scheduled{&conn, client}, Service{&conn, client}, Session{&conn, client}, Status{&conn, client}, Template{&conn, client}, Tenant{&conn, client}, Triggered{&conn, client}, VRRP{&conn, client}, Workflow{&conn, client}, WorkflowTemplate{&conn, client}}
}
