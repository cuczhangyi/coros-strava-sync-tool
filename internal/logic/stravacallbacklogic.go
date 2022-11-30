package logic

import (
	"context"
	"fmt"

	"github.com/cuczhangyi/coros-strava/internal/svc"
	"github.com/cuczhangyi/coros-strava/internal/types"
	"github.com/cuczhangyi/coros-strava/utils"
	"github.com/gogf/gf/errors/gerror"
	strava "github.com/strava/go.strava"

	"github.com/tal-tech/go-zero/core/logx"
)

type StravaCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}




func NewStravaCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) StravaCallbackLogic {
	return StravaCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StravaCallbackLogic) StravaCallback(req types.StravaCallbackReq) (*types.StravaCallbackResp, error) {
	// todo: add your logic here and delete this line
	logx.Info("StravaCallbackLogic.StravaCallback")

	fmt.Println("StravaCallbackLogic.StravaCallback")
	if req.Error !=  "" {
		return nil, strava.OAuthAuthorizationDeniedErr
	}

	if req.Code == ""{
		logx.Errorf("strava callbacl no code ")
		return nil, gerror.New("strava callbacl no code ")
	}

	fmt.Println("StravaCallbackLogic.StravaCallback code is " +  req.Code)
	err:= utils.AuthGetTokenByCode(req.Code)
	if err != nil {
		logx.Errorf("AuthGetTokenByCode error: %v", err)
		return nil, gerror.New("AuthGetTokenByCode error")
	}
	
	return &types.StravaCallbackResp{}, nil
}
