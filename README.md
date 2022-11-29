# coros-strava
## 一个自动同步coros 活动 到strava 的工具

## 使用方法
### 配置方法
1. 在/etc/config.yaml 输入您的coros 帐号，密码
2. 在/etc/corosstrata-api.yaml 输入您的strava clientId 和 screct (如何获取 client 和 screct[https://developers.strava.com/docs/getting-started/])
3. 如果您有回调域名，请误必在strava 的回调域名处填写好。
3. 如果您是运行在有域名的服务器上，请在在/etc/corosstrata-api.yaml 配置上带有域名的的callurl ， 如果本地运行，请保持配置中的localhost （回调路径是 /strava/callback）

### 使用流程
1. 在程序运行后会输出一个地址，请您第一次人工访问该地址，然后点击授权
2. 程序会在收到strava 的回调授权之后， 开始定时刷新token ，并且每天会自动获取前一天的coros 数据，自动上传。



