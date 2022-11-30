package utils

import (
	"encoding/json"
	"fmt"

	"github.com/cavaliergopher/grab/v3"
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/guid"
)


type LoginReq struct {
	Account string `json:"account"`
	AccountType int `json:"accountType"`
	Pwd string `json:"pwd"`
}


type ActiveListResp struct {
	APICode string `json:"apiCode"`
	Data ActiveData `json:"data"`
	Message string `json:"message"`
	Result string `json:"result"`
}
type DataList struct {
	AdjustedPace int `json:"adjustedPace"`
	Ascent int `json:"ascent"`
	Avg5X10S int `json:"avg5x10s"`
	AvgCadence int `json:"avgCadence"`
	AvgHr int `json:"avgHr"`
	AvgPower int `json:"avgPower"`
	AvgSpeed int `json:"avgSpeed"`
	AvgStrkRate int `json:"avgStrkRate"`
	Best int `json:"best"`
	Best500M int `json:"best500m"`
	BestKm int `json:"bestKm"`
	BestLen int `json:"bestLen"`
	BodyTemperature int `json:"bodyTemperature"`
	Cadence int `json:"cadence"`
	Calorie int `json:"calorie"`
	Date int `json:"date"`
	Descent int `json:"descent"`
	Device string `json:"device"`
	Distance float64 `json:"distance"`
	DownhillDesc int `json:"downhillDesc"`
	DownhillDist float64 `json:"downhillDist"`
	DownhillTime int `json:"downhillTime"`
	EndTime int `json:"endTime"`
	EndTimezone int `json:"endTimezone"`
	HasMessage int `json:"hasMessage"`
	ImageURL string `json:"imageUrl"`
	ImageURLType int `json:"imageUrlType"`
	IsShowMs int `json:"isShowMs"`
	LabelID string `json:"labelId"`
	Lengths int `json:"lengths"`
	Max2S int `json:"max2s"`
	MaxGrade int `json:"maxGrade"`
	MaxSlope int `json:"maxSlope"`
	MaxSpeed int `json:"maxSpeed"`
	Mode int `json:"mode"`
	Name string `json:"name"`
	Np int `json:"np"`
	Pitch int `json:"pitch"`
	Sets int `json:"sets"`
	SpeedType int `json:"speedType"`
	SportType int `json:"sportType"`
	StartTime int `json:"startTime"`
	StartTimezone int `json:"startTimezone"`
	Step int `json:"step"`
	SubMode int `json:"subMode"`
	Swolf int `json:"swolf"`
	Total float64 `json:"total"`
	TotalDescent int `json:"totalDescent"`
	TotalReps int `json:"totalReps"`
	TotalTime int `json:"totalTime"`
	TrainingLoad int `json:"trainingLoad"`
	UnitType int `json:"unitType"`
	WaterTemperature int `json:"waterTemperature"`
	WorkoutTime int `json:"workoutTime"`
	savePath  string `json:"save_path,omitempty"`
}
type ActiveData struct {
	Count int `json:"count"`
	DataList []*DataList `json:"dataList"`
	PageNumber int `json:"pageNumber"`
	TotalPage int `json:"totalPage"`
}

type FileDownloadUrlResp struct {
	APICode string `json:"apiCode"`
	Data FileURLData  `json:"data"`
	Message string `json:"message"`
	Result string `json:"result"`
}
type FileURLData struct {
	FileURL string `json:"fileUrl"`
}

func formatDate(date string) (string,error) {
	dateStr := gstr.Trim(date)
	dateItem := gtime.NewFromStr(dateStr)
	if dateItem == nil {
		return "",gerror.New("date format error")
	}
	return dateItem.Format("Ymd"),nil
}


func login(user_name, password string) (string, error) {
	//{"account":"123@123.com","accountType":2,"pwd":"base64xxx"}
	// Get the token
	pwdBase64 ,_:= gmd5.EncryptString(password)
	url := corosURL + "/account/login"
	req := &LoginReq{
		Account: user_name,
		Pwd: pwdBase64,
		AccountType: 2, //默认为2 
	}
	jsonStr ,_:= json.Marshal(req)
	resp,err := g.Client().Post(url ,jsonStr );
	if err != nil {
		return "", err
	}
	defer resp.Close()
	respStr := resp.ReadAllString()

	jsonObj ,err := gjson.LoadContent(respStr,true)
	if err != nil {
		return "", gerror.New("login failed")
	}
	retMsg := jsonObj.GetString("message")
	if retMsg != "OK" {
		return "", gerror.New("login failed")
	}
	retCode := jsonObj.GetString("result")
	if retCode != "0000" {
		return "", gerror.New("login failed")
	}

	token := jsonObj.GetString("data.accessToken")
	if token == "" {
		return "", gerror.New("no token in response")
	}
	return token,nil
}

func getActivesByHttp(strartDate string, endDate string, token string, pageCurrent int64) ([]*DataList,int64,error){
	url := corosURL + "/activity/query?size=50&pageNumber=%d&modeList=&startDay=%s&endDay=%s"	
	url = fmt.Sprintf(url,pageCurrent , strartDate ,endDate )
	resp,err:= g.Client().SetHeader("accesstoken",token).Get(url)
	if err != nil {
		return  nil, 0, err
	}
	defer resp.Close()

	retStr := resp.ReadAllString()
	jsonObj ,err := gjson.LoadContent(retStr,true)
	if err != nil {
		return  nil, 0,err
	}

	retMsg := jsonObj.GetString("message")
	if retMsg != "OK" {
		return nil, 0,  gerror.New("load_active failed")
	}
	retCode := jsonObj.GetString("result")
	if retCode != "0000" {
		return nil, 0, gerror.New("load_active failed")
	}

	retStruct := &ActiveListResp{}

	err = jsonObj.Struct(retStruct )
	if err != nil {
		return nil, 0,err
	}

	nextPage := pageCurrent
	if pageCurrent < gconv.Int64(retStruct.Data.TotalPage) {
		nextPage += 1
	}

	
	return retStruct.Data.DataList,nextPage,nil
}


func getAllActives(strartDate string, endDate string, token string) ([]*DataList, error ){
	// Get the activities
	//https://teamcnapi.coros.com/activity/query?size=20&pageNumber=1&modeList=
	// beginDate,err := formatDate(strartDate)
	// if err != nil {
	// 	return nil, err
	// }
	// stopDate,err := formatDate(endDate)
	// if err != nil {
	// 	return nil, err
	// }
	
	beginDate := strartDate
	stopDate := endDate	
	
	curPage := 1
	allActives := make([]*DataList,0)
	var errWhile error
	for{
		data,nextPage,err:= getActivesByHttp(beginDate,stopDate,token,int64(curPage))
		if err != nil {
			
			errWhile = err
			break
		}

		

		allActives = append(allActives,data...)
		if nextPage == int64(curPage){
			break
		}
	}	
	if errWhile != nil {
		return nil,errWhile
	}
	return allActives ,nil 
}

func getActiveFitFileUrl(activeId string , token string ) (string,error){
	//https://teamcnapi.coros.com/activity/detail/download?labelId=447852161146060800&sportType=100&fileType=4

	url := corosURL + "/activity/detail/download?labelId=%s&sportType=100&fileType=4"
	url = fmt.Sprintf(url,activeId)
	resp,err:= g.Client().SetHeader("accesstoken",token).Post(url)
	if err != nil {
		 return "",err
	}

	if err != nil {
		return "",err
	}
	defer resp.Close()

	retStr := resp.ReadAllString()
	jsonObj ,err := gjson.LoadContent(retStr,true)
	if err != nil {
		return "",err
	}

	retMsg := jsonObj.GetString("message")
	if retMsg != "OK" {
		return "",  gerror.New("load_active_file failed")
	}
	retCode := jsonObj.GetString("result")
	if retCode != "0000" {
		return "", gerror.New("load_active_file failed")
	}

	return jsonObj.GetString("data.fileUrl"),nil

}

func downloadFitFile(url string) (string,error){
	itemName := guid.S()
	file_path := gstr.TrimRight(g.Cfg().GetString("coros.file_save_path"),"/")+"/"+gtime.Now().Format("Ymd")+"/"
	if !gfile.Exists(file_path){
		gfile.Mkdir(file_path)
	}

	savePath := gstr.TrimRight(file_path, "/") + "/" + itemName + ".fit"
	_, err := grab.Get(savePath, url)
	if err != nil {
		return "", gerror.New("downloadFitFile failed:"+ url )
	}
	return savePath, nil	
}


func DownloadAllCorosFiles(user_name, password string, strartDate string, endDate string) ([]*DataList, error ){

	if user_name == "" {
		user_name = g.Cfg().GetString("coros.user_email")
	}
	if password == ""{
		password = g.Cfg().GetString("coros.password")
	}

	token,err:= login(user_name,password)
	if err != nil {
		fmt.Printf("coros_login failed %v",err)
		return nil, err 
	}

	datas,err:= getAllActives(strartDate,endDate,token)
	if err != nil {
		return nil, err
	}

	for index,data := range datas {
		url,err :=getActiveFitFileUrl(data.LabelID,token)
		if err != nil{
			fmt.Println("getActiveFitFileUrl failed "+data.LabelID + "" +err.Error() )
			continue
		}
		savePath,err := downloadFitFile(url)
		if err != nil{
			fmt.Println("downloadFitFile failed "+data.LabelID + "" +err.Error())
			continue
		}
		datas[index].savePath = savePath
	}
	return datas,nil 
}
