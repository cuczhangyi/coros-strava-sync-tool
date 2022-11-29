package utils

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
	"github.com/tal-tech/go-zero/core/logx"
)

//https://developers.strava.com/docs/authentication/#refreshingexpiredaccesstokens



func UploadItem() {
	file_path := gstr.TrimRight(g.Cfg().GetString("coros.file_save_path"),"/")+"/"+gtime.Now().Format("Ymd")+"/"

	if !gfile.Exists(file_path){
		gfile.Mkdir(file_path)
	}


	if StravaAuthToken == "" {
		logx.Errorf("StravaAuthToken is empty when upload item")
		return
	}

	checktime := gtime.NewFromTimeStamp(ExpireAt)
	if gtime.Now().After(checktime ){
		logx.Errorf("StravaAuthToken is empty when upload item")
		return
	}

	//file_path := gstr.TrimRight(g.Cfg().GetString("coros.file_save_path"),"/")+"/"+gtime.Now().Format("Ymd")+"/"

	if !gfile.Exists(file_path){
		gfile.Mkdir(file_path)
	}
	
	
	uploader := newUploader(StravaAuthToken)
	files, err := ioutil.ReadDir(file_path)
	if err != nil {
		logx.Errorf("read dir error: %v", err)
		return
	}

	var wg sync.WaitGroup
	log.Printf("Processing %d files\n", len(files))

	for _, f := range files {
		if f.IsDir() || strings.HasPrefix(f.Name(), ".") || !strings.HasSuffix(f.Name(), ".fit") {
			log.Printf("Ignoring %s\n", f.Name())
			continue
		}
		wg.Add(1)
		go func(fname string) {
			defer wg.Done()
			filePath := path.Join(file_path, fname)
			f, err := os.Open(filePath)
			if err != nil {
				log.Printf("%s - open: %s", fname, err)
				return
			}
			aid, err := uploader.Upload(fname, f)
			if err != nil {
				logx.Errorf(
					"%s - Activity created, you can view it at http://www.strava.com/activities/%d",
					fname, aid)
				return
			} else {
				defer gfile.Remove(filePath)
			}
			logx.Infof(
				"%s - Activity created, you can view it at http://www.strava.com/activities/%d",
				fname, aid)
		}(f.Name())
	}
	wg.Wait()
}
