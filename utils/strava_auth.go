package utils

import (
	"fmt"
	"time"

	"github.com/cuczhangyi/coros-strava/internal/svc"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/encoding/gurl"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/tal-tech/go-zero/core/logx"
)

// AuthHandler provides an url to direct the user to as well as
// an http.HandlerFunc to handle the redirect from the remote host.


var StravaAccessToken string 

var StravaRefreshToken string 

var StravaTokenExpireAt int64



const(
	strava_base_url = "https://www.strava.com/api/v3"
)

func GetAuthUrl() string{
	
	basePath := strava_base_url
	callbackUrl := svc.SCtx.Config.StravaInfo.CallBackUrl
	callbackUrl = gurl.Encode(callbackUrl)
	scope := "activity:write"
	path := fmt.Sprintf("%s/oauth/authorize?client_id=%d&response_type=code&redirect_uri=%s&scope=%v", basePath, svc.SCtx.Config.StravaInfo.ClientId, callbackUrl, scope)
	path += "&state=" + "state1"
	path += "&approval_prompt=force"
	return path 
}

func AuthGetTokenByCode(code string )  error{

	basePath := strava_base_url
	url := fmt.Sprintf("%s/oauth/token?client_id=%d&client_secret=%s&code=%s&grant_type=authorization_code", basePath, svc.SCtx.Config.StravaInfo.ClientId, svc.SCtx.Config.StravaInfo.ClientSecret, code)
	
	resp,err:= g.Client().Post(url)
	if err != nil{
		logx.Errorf("AuthGetTokenByCode error:%v code is %s", err,code)
		return nil 
	}
	defer resp.Close()

	respStr := resp.ReadAllString()

	jsonObj ,err:= gjson.LoadContent(respStr)
	if err != nil{
		logx.Errorf("AuthGetTokenByCode error:%v gjson.LoadContent %s", err,respStr)
		return err 
	}

	StravaAccessToken = jsonObj.GetString("access_token")
	StravaRefreshToken = jsonObj.GetString("refresh_token")
	StravaTokenExpireAt = jsonObj.GetInt64("expires_at")

	if StravaAccessToken == "" || StravaRefreshToken == "" || StravaTokenExpireAt == 0{
		logx.Errorf("AuthGetTokenByCode return error no token", )
		return err 
	}
	return nil
}


func RefreshStravaToken() error{

	// 	curl -X POST https://www.strava.com/api/v3/oauth/token \
	//   -d client_id=ReplaceWithClientID \
	//   -d client_secret=ReplaceWithClientSecret \
	//   -d grant_type=refresh_token \
	//   -d refresh_token=ReplaceWithRefreshToken
	
	if StravaTokenExpireAt == 0 || StravaRefreshToken == ""{
		fmt.Println("RefreshStravaToken error no token")
		return	gerror.New("no token")
	}
	checktime := gtime.NewFromTimeStamp(StravaTokenExpireAt).Add(-5 * time.Minute)
	checkTimeStr := checktime.String()
	fmt.Println("RefreshStravaToken checkTimeStr " +  checkTimeStr)
	if !gtime.Now().After(checktime){
		fmt.Println("RefreshStravaToken checkTime is error")
		return nil
	}

	if gtime.Now ().After(gtime.NewFromTimeStamp(StravaTokenExpireAt)){
		fmt.Println("strava token expired , please restart this app and reauth")
		return gerror.New("strava token expired , please restart this app and reauth")
	}

	url := fmt.Sprintf("https://www.strava.com/api/v3/oauth/token?client_id=%d&client_secret=%s&grant_type=refresh_token&refresh_token=%s", svc.SCtx.Config.StravaInfo.ClientId, svc.SCtx.Config.StravaInfo.ClientSecret, StravaRefreshToken)

	resp,err:= g.Client().Post(url)
	if err != nil{
		return err 
	}
	defer resp.Close()

	respStr := resp.ReadAllString()

	jsonObj,err := gjson.LoadContent(respStr)
	if err != nil{
		return err 
   	}
	
	AccessToken1 := jsonObj.GetString("access_token")
	RefreshToken1 := jsonObj.GetString("refresh_token")
	ExpireAt1 := jsonObj.GetInt64("expires_at")

	if AccessToken1 == "" || RefreshToken1 == "" || ExpireAt1 == 0{
		logx.Errorf("AuthGetTokenByCode return error no token", )
		return err 
	}
	StravaAccessToken  = AccessToken1
	StravaRefreshToken = RefreshToken1
	StravaTokenExpireAt = ExpireAt1
	fmt.Println("strava token refresh success")
	return nil
}