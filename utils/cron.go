package utils

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)


func GetCorosFileAndUpload() {

	if StravaAccessToken == ""{
		return	;
	}

	userName := g.Cfg().GetString("coros.user_email")
	password := g.Cfg().GetString("coros.password")
	//strartDate := gtime.Now().AddDate(0,0,-1).Format("Ymd")
	strartDate := gtime.Now().Format("Ymd")
	DownloadAllCorosFiles(userName,password,strartDate,strartDate)
	UploadItem();
}