package main

import (
	"flag"
	"fmt"

	"github.com/cuczhangyi/coros-strava/internal/config"
	"github.com/cuczhangyi/coros-strava/internal/handler"
	"github.com/cuczhangyi/coros-strava/internal/svc"
	"github.com/cuczhangyi/coros-strava/utils"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcron"

	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/rest"
)


var configFile = flag.String("f", "etc/corosstrata-api.yaml", "the config file")
var gfcfgFile = flag.String("gffile", "config.yaml", "the gf config file") //这行新增
var gfcfgPath = flag.String("gfpath", "./etc", "the gf config file")//这行新增


func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)

	g.Cfg().SetPath(*gfcfgPath)
	g.Cfg().SetFileName(*gfcfgFile)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	
	
	url := utils.GetAuthUrl()
	fmt.Printf("please visit url to get the code %s \n", url)
	//utils.Upload("/Volumes/DATA/GoWorkSpace/src/github.com/cuczhangyi/coros-strava/test_fit_file");


	refreshTokenCron := g.Cfg().GetString("cron_plan.refresh_token")
	gcron.Add(refreshTokenCron,func() {
		utils.RefreshStravaToken()
	})	
	
	uploadCron := g.Cfg().GetString("cron_plan.upload_item")

	gcron.Add(uploadCron,func() {
		utils.UploadItem()
	})





	server.Start()
}
