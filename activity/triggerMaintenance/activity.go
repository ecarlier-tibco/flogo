/*
 * Copyright © 2017. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package triggerMaintenance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

const (

	// Live Apps config
	tibcoAccountURL = "https://sso-ext.tibco.com:443"
	tenantID        = "BPM"

	// Live Apps application config
	activityID      = "ZZhUPxJOmEeeiw8XG50L7jw"
	applicationName = "Maintenance1"
	activityName    = "Task"
	processName     = "EnterMaintenanceData1"
	processLabel    = "Trigger Case on Alert"

	// config params
	deviceIDConfigParam  = "DeviceID"
	latitudeConfigParam  = "DeviceLatitude"
	longitudeConfigParam = "DeviceLongitude"

	// params
	alertParam                = "alert"
	ovResult                  = "result"
	usernameParam             = "username"
	passwordParam             = "password"
	regionParam               = "region"
	accountIDParam            = "accountID"
	sandboxIDParam            = "sandboxID"
	applicationIDParam        = "applicationID"
	applicationCreatorIDParam = "applicationCreatorID"
	applicationVersionParam   = "applicationVersion"
)

var username string
var password string
var region string
var accountID string
var liveAppsURL string

var activityLog = logger.GetLogger("tibco-activity-maintenance-trigger")

type TriggerMaintenanceActivity struct {
	metadata *activity.Metadata
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &TriggerMaintenanceActivity{metadata: metadata}
}

func (a *TriggerMaintenanceActivity) Metadata() *activity.Metadata {
	return a.metadata
}
func (a *TriggerMaintenanceActivity) Eval(context activity.Context) (done bool, err error) {
	activityLog.Info("Executing Live Apps Trigger Maintenance Case activity")

	// sequence:
	// read data from activity inputs
	// authenticate and get access token
	// login to get Atmosphere session cookie
	// (Atmosphere cookie needs to be on subsequent calls)
	// run caseCreatory() and get process id
	// updateProcess to release creator task with data
	// output the process id that we get back from the updateProcess call

	// Validates that the connection parameters have been set
	username = context.GetInput(usernameParam).(string)
	if username == "" {
		return false, activity.NewError("Live Apps connection user name is not configured", "LIVEAPPS-CON-2000", nil)
	}
	password = context.GetInput(passwordParam).(string)
	if password == "" {
		return false, activity.NewError("Live Apps connection password is not configured", "LIVEAPPS-CON-2000", nil)
	}
	region = context.GetInput(regionParam).(string)
	if region == "" {
		return false, activity.NewError("Live Apps connection region is not configured", "LIVEAPPS-CON-2000", nil)
	}
	accountID = context.GetInput(accountIDParam).(string)
	if accountID == "" {
		return false, activity.NewError("Live Apps connection account ID is not configured", "LIVEAPPS-CON-2000", nil)
	}

	if region == "eu" || region == "EU" {
		liveAppsURL = "https://eu.liveapps.cloud.tibco.com"
	} else if region == "us" || region == "US" {
		liveAppsURL = "https://liveapps.cloud.tibco.com"
	} else if region == "au" || region == "AU" {
		liveAppsURL = "https://au.liveapps.cloud.tibco.com"
	} else {
		return false, activity.NewError("Live Apps region (EU/US/AU) is invalid", "LIVEAPPS-REGION-3000", nil)
	}

	// read data from activity inputs
	typedVal, ok := data.GetGlobalScope().GetAttr(deviceIDConfigParam)
	if !ok {
		errorMsg := fmt.Sprintf("Config Param Device ID not set: '%s'", deviceIDConfigParam)
		activityLog.Error(errorMsg)
		return false, activity.NewError(errorMsg, "", nil)
	}
	deviceIDVal := typedVal.Value.(string)

	typedVal, ok = data.GetGlobalScope().GetAttr(latitudeConfigParam)
	if !ok {
		errorMsg := fmt.Sprintf("Config Param Device Latitude not set: '%s'", latitudeConfigParam)
		activityLog.Error(errorMsg)
		return false, activity.NewError(errorMsg, "", nil)
	}
	latitudeStr := typedVal.Value.(string)
	latitude, err := strconv.ParseFloat(latitudeStr, 64)
	if err != nil {
		return false, activity.NewError("Config Param Device Latitude is not a proper number", "LIVEAPPS-NUMBER-3000", nil)
	}

	typedVal, ok = data.GetGlobalScope().GetAttr(longitudeConfigParam)
	if !ok {
		errorMsg := fmt.Sprintf("Config Param Device Longitude not set: '%s'", longitudeConfigParam)
		activityLog.Error(errorMsg)
		return false, activity.NewError(errorMsg, "", nil)
	}
	longitudeStr := typedVal.Value.(string)
	longitude, err := strconv.ParseFloat(longitudeStr, 64)
	if err != nil {
		return false, activity.NewError("onfig Param Device Longitude is not a proper number", "LIVEAPPS-NUMBER-3000", nil)
	}

	alert, _ := context.GetInput(alertParam).(string)

	// ok now make the API calls
	token := getToken()
	atmosphereCookie := performLogin(token, context)
	id := startPF(atmosphereCookie, context)
	result := updatePF(atmosphereCookie, context, id, deviceIDVal, alert, latitude, longitude)

	// we are done so return the result from updating the pageflow
	context.SetOutput(ovResult, result)

	return true, nil
}

func getToken() string {
	method := "POST"
	uri := tibcoAccountURL + "/as/token.oauth2?username=" + username + "&password=" + password + "&client_id=ropc_ipass&grant_type=password"

	var reqBody io.Reader

	contentType := "application/json; charset=UTF-8"
	reqBody = nil
	req, err := http.NewRequest(method, uri, reqBody)
	if reqBody != nil {
		req.Header.Set("Content-Type", contentType)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		activityLog.Error(err.Error())
	}
	defer resp.Body.Close()

	type Message struct {
		AccessToken string `json:"access_token"`
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		activityLog.Error(err.Error())
	}

	var s = new(Message)
	err2 := json.Unmarshal(body, &s)

	if err2 != nil {
		activityLog.Error("Error from unmarshal Tibco Account auth response:", err2)
	}
	return s.AccessToken
}

func performLogin(token string, context activity.Context) []*http.Cookie {
	method := "POST"
	uri := liveAppsURL + ":443/idm/v1/login-oauth"

	var reqBody io.Reader

	contentType := "application/x-www-form-urlencoded"

	data := url.Values{}
	data.Set("TenantId", tenantID)
	data.Add("AccessToken", token)
	accountID, ok := context.GetInput(accountIDParam).(string)

	if !ok {
		errorMsg := fmt.Sprintf("Param AccountID not set: '%s'", accountIDParam)
		activityLog.Error(errorMsg)
		return nil
	}
	data.Add("AccountId", accountID)
	reqBody = bytes.NewBufferString(data.Encode())

	req, err := http.NewRequest(method, uri, reqBody)

	if reqBody != nil {
		req.Header.Set("Content-Type", contentType)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// save the cookies as these are needed for subsequent API calls
	cookie := resp.Cookies()
	return cookie
}

func startPF(cookies []*http.Cookie, context activity.Context) string {

	sandboxID, ok := context.GetInput(sandboxIDParam).(string)

	if !ok {
		errorMsg := fmt.Sprintf("Param SandboxID not set: '%s'", sandboxIDParam)
		activityLog.Error(errorMsg)
		return ""
	}

	method := "POST"
	uri := liveAppsURL + ":443/pageflow/start?$sandbox=" + sandboxID

	contentType := "application/json"

	type CreatorMessage struct {
		ID              int      `json:"id"`
		Name            string   `json:"name"`
		Label           string   `json:"label"`
		Version         int      `json:"version"`
		ApplicationID   int      `json:"applicationId"`
		ApplicationName string   `json:"applicationName"`
		ActivityID      string   `json:"activityId"`
		ActivityName    string   `json:"activityName"`
		Roles           []string `json:"roles"`
	}

	creator := new(CreatorMessage)

	applicationCreatorID, _ := context.GetInput(applicationCreatorIDParam).(int)
	applicationID, _ := context.GetInput(applicationIDParam).(int)
	applicationVersion, _ := context.GetInput(applicationVersionParam).(int)

	creator.ID = applicationCreatorID
	creator.Name = processName
	creator.Label = processLabel
	creator.Version = applicationVersion
	creator.ApplicationID = applicationID
	creator.ApplicationName = applicationName
	creator.ActivityID = activityID
	creator.ActivityName = activityName
	creator.Roles = []string{}

	// create JSON representation of creator payload

	bodyStr, err1 := json.Marshal(*creator)
	if err1 != nil {
		activityLog.Error("Error creating start creator payload from input data: " + err1.Error())
	}

	req, err := http.NewRequest(method, uri, bytes.NewBuffer(bodyStr))
	req.Header.Set("Content-Type", contentType)

	cookieLen := len(cookies)
	atmosCookie := new(http.Cookie)

	// find the atmosphere session cookie in the passed cookies attach it to this request
	for i := 0; i < cookieLen; i++ {
		if cookies[i].Name == "AtmosphereSession" {
			atmosCookie = cookies[i]
		}
	}
	req.AddCookie(atmosCookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		activityLog.Error(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		activityLog.Error(err.Error())
	}

	// to understand the response we create a struct then unmarshal the JSON response to that struct
	// hint: there are various web tools that generate struct from JSON and vice versa
	type AppID struct {
		AppID          int `json:"appId"`
		SandboxID      int `json:"sandboxId"`
		SubscriptionID int `json:"subscriptionId"`
	}

	type Message struct {
		ID            string `json:"id"`
		Name          string `json:"name"`
		Version       int    `json:"version"`
		ActivityName  string `json:"activityName"`
		ApplicationID AppID  `json:"applicationId"`
		OldCaseState  string `json:"oldCaseState"`
	}

	var s = new(Message)
	err2 := json.Unmarshal(body, &s)

	if err2 != nil {
		activityLog.Info("Error from unmarshal startPf response:", err2)
	}
	activityLog.Info("Case Creator Started")

	return s.ID
}

func updatePF(cookies []*http.Cookie, context activity.Context, id string, deviceIdValue string, alertValue string, latValue float64, lngValue float64) string {

	method := "POST"
	sandboxID, _ := context.GetInput(sandboxIDParam).(string)

	uri := liveAppsURL + "/pageflow/update?$sandbox=" + sandboxID

	contentType := "application/json"

	// we need to create the object for the updatePF action
	type Attrib struct {
		Op    string      `json:"op"`
		Path  string      `json:"path"`
		Rank  int         `json:"rank"`
		Value interface{} `json:"value"`
	}

	type Maintenance struct {
		MaintenanceCaseData []Attrib `json:"Maintenance1"`
	}

	// note: data is actually the JSON string representation of Maintenance Case
	type Req1 struct {
		Data string `json:"data"`
		ID   string `json:"id"`
	}

	deviceId := new(Attrib)
	deviceId.Op = "add"
	deviceId.Path = "/DeviceID_v1/"
	deviceId.Rank = 0
	deviceId.Value = deviceIdValue

	alert := new(Attrib)
	alert.Op = "add"
	alert.Path = "/AlertDescription_v1/"
	alert.Rank = 1
	alert.Value = alertValue

	location := new(Attrib)
	location.Op = "add"
	location.Path = "/SiteLocation_v1/"
	location.Rank = 2
	var e struct{}
	location.Value = e

	lat := new(Attrib)
	lat.Op = "add"
	lat.Path = "/SiteLocation_v1/Lat_v1/"
	lat.Rank = 0
	lat.Value = latValue

	lng := new(Attrib)
	lng.Op = "add"
	lng.Path = "/SiteLocation_v1/Lng_v1/"
	lng.Rank = 1
	lng.Value = lngValue

	maintenanceCase := new(Maintenance)

	maintenanceCase.MaintenanceCaseData = make([]Attrib, 5)

	maintenanceCase.MaintenanceCaseData[0] = *deviceId
	maintenanceCase.MaintenanceCaseData[1] = *alert
	maintenanceCase.MaintenanceCaseData[2] = *location
	maintenanceCase.MaintenanceCaseData[3] = *lat
	maintenanceCase.MaintenanceCaseData[4] = *lng

	requestObj := new(Req1)

	// Data is actually a JSON string inside this JSON payload
	datastr, err1 := json.Marshal(*maintenanceCase)
	if err1 != nil {
		activityLog.Error("Error creating maintenance case from input data: " + err1.Error())
	}

	activityLog.Info("Posted Data: ", bytes.NewBuffer(datastr).String())
	requestObj.Data = bytes.NewBuffer(datastr).String()
	requestObj.ID = id

	jsonstr, err := json.Marshal(requestObj)
	if err != nil {
		activityLog.Error("Error creating maintenance case json string: " + err.Error())
	}

	bodyStr := jsonstr

	req, err := http.NewRequest(method, uri, bytes.NewBuffer(bodyStr))
	req.Header.Set("Content-Type", contentType)

	cookieLen := len(cookies)
	atmosCookie := new(http.Cookie)

	for i := 0; i < cookieLen; i++ {
		if cookies[i].Name == "AtmosphereSession" {
			atmosCookie = cookies[i]
		}
	}

	req.AddCookie(atmosCookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		activityLog.Error("DO ERR: ", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		activityLog.Error("BODY ERR: ", err.Error())
	}

	// to understand the response we create a struct then unmarshal the JSON response to that struct
	// hint: attribute names should start with upperclass letter so they are exported
	// hint: the json tag after the definition tells go the real json attribute name
	// hint: there are various web tools that generate struct from JSON and vice versa
	type UpdateResp struct {
		UpdatedInstID string `json:"updatedInstId"`
	}

	activityLog.Info("RESP BODY: ", bytes.NewBuffer(body).String())

	var s = new(UpdateResp)
	err2 := json.Unmarshal(body, &s)
	if err2 != nil {
		activityLog.Error("Error from unmarshal updatePF response::", err2)
	}

	activityLog.Info("Maintenance Case Creator complete: ID " + s.UpdatedInstID)

	return s.UpdatedInstID
}
