## 食用方式
0. 新建一个文件夹
```shell
mkdir wechat-hook # 名字随便写
cd wechat-hook
```
1. 创建配置文件`config.yaml`
```shell
mkdir config # 先创建一个文件夹保存配置文件，文件名不要变
vim config.yaml # 编辑配置文件，内容如下
```
```yaml
# 微信HOOK配置
wechat:
  # 微信HOOK接口地址
  host: wechat:19088
  # 是否在启动的时候自动设置hook服务的回调
  autoSetCallback: true
  # 回调IP，如果是Docker运行，本参数必填，如果Docker修改了映射，格式为 ip:port，如果使用项目提供的docker-compsoe.yaml文件启动，可以不写
  callback: 

# 数据库
mysql:
  host: mysql
  port: 3306
  user: wechat
  password: wechat
  db: wechat

task:
  enable: false
  syncFriends:
      enable: true
      cron: '0 * * * *'
  waterGroup:
      enable: true
      cron: '30 9 * * *'
      # 需要发送水群排行榜的群Id
      groups:
        - '11111@chatroom' # 自行配置
      # 不计入统计范围的用户Id
      blacklist:
        - 'wxid_xxxx' # 自行配置
```

2. 创建`docker-compose.yaml`文件
```yaml
version: '3'

services:
  wechat:
    image: lxh01/wxhelper-docker:3.9.5.81
    container_name: gw-wechat
    restart: unless-stopped
    volumes:
      - ./data/wechat:/home/app/.wine/drive_c/users/app/Documents/WeChat\ Files
    ports:
      - "8080:8080"
      - "19088:19088"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:19088/api/checkLogin"]
      interval: 60s
      timeout: 10s
      retries: 5


  mysql:
    image: mysql:8
    container_name: gw-db
    restart: unless-stopped
    depends_on:
      - wechat
    environment:
      - MYSQL_ROOT_PASSWORD=wechat
      - MYSQL_USER=wechat
      - MYSQL_PASSWORD=wechat
      - MYSQL_DATABASE=wechat
    volumes:
      - ./data/db:/var/lib/mysql


  wxhelper:
    image: gitee.ltd/lxh/go-wxhelper:latest
    container_name: gw-service
    restart: unless-stopped
    depends_on:
      - mysql
    volumes:
      # 配置文件请参阅项目根目录的config.yaml文件
      - ./config/config.yaml:/app/config.yaml
    ports:
      - "19099:19099"

```

3. 启动  
`这玩意儿有点儿不完善，先启动wechat，确定启动起来了再启动剩余的`
```shell
# 以下命令选个能用的就行
docker-compose up -d # 老版本
docker compose up -d # 新版本

# 分顺序启动
docker compose up -d wechat # 1. 启动微信
docker compose up -d # 启动剩余的
```