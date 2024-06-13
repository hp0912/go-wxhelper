### 离线安装 docker

#### 下载 Docker 安装包

下载地址：https://download.docker.com/linux/static/stable/x86_64/docker-24.0.7.tgz

文件大小约：64MB

#### 上传到服务器

```
[root@uat-0001 ~]# ll -h
总用量 82M
-rw-r--r--  1 root root  64M 12月 20 20:54 docker-24.0.7.tgz
```

#### 解压软件包

```
[root@uat-0001 ~]# tar -zxvf docker-24.0.7.tgz
docker/
docker/docker
docker/docker-init
docker/dockerd
docker/runc
docker/ctr
docker/containerd-shim-runc-v2
docker/containerd
docker/docker-proxy
```

#### 将 docker 目录下文件复制到/usr/bin 目录下

```
[root@uat-0001 ~]# ls docker
containerd  containerd-shim-runc-v2  ctr  docker  dockerd  docker-init  docker-proxy  runc
[root@uat-0001 ~]# cp docker/* /usr/bin/
```

#### 创建 docker.service 文件

进入到/usr/lib/systemd/system/目录下，我们编辑创建 docker.service 文件，用于管理 docker 服务，复制黏贴如下内容即可。

```
[root@uat-0001 ~]# vim /usr/lib/systemd/system/docker.service
```

```ini
[Unit]
Description=Docker Application Container Engine
Documentation=http://docs.docker.com
After=network.target docker.socket
[Service]
Type=notify
EnvironmentFile=-/run/flannel/docker
WorkingDirectory=/usr/local/bin
ExecStart=/usr/bin/dockerd \
                -H tcp://0.0.0.0:4243 \
                -H unix:///var/run/docker.sock \
                --selinux-enabled=false \
                --log-opt max-size=100m
ExecReload=/bin/kill -s HUP $MAINPID
LimitNOFILE=infinity
LimitNPROC=infinity
LimitCORE=infinity
TimeoutStartSec=0
Delegate=yes
KillMode=process
Restart=on-failure
[Install]
WantedBy=multi-user.target
```

#### 重新加载 daemon-reload

```
systemctl daemon-reload
```

#### 启动 docker

```
systemctl start docker
```

#### 查看服务启动状态

```
systemctl status docker
```

#### 查看 Docker 版本

```
[root@uat-0001 ~]# docker version
Client:
 Version:           24.0.7
 API version:       1.43
 Go version:        go1.20.10
 Git commit:        afdd53b
 Built:             Thu Oct 26 09:04:00 2023
 OS/Arch:           linux/amd64
 Context:           default

Server: Docker Engine - Community
 Engine:
  Version:          24.0.7
  API version:      1.43 (minimum version 1.12)
  Go version:       go1.20.10
  Git commit:       311b9ff
  Built:            Thu Oct 26 09:05:28 2023
  OS/Arch:          linux/amd64
  Experimental:     false
 containerd:
  Version:          v1.7.6
  GitCommit:        091922f03c2762540fd057fba91260237ff86acb
 runc:
  Version:          1.1.9
  GitCommit:        v1.1.9-0-gccaecfc
 docker-init:
  Version:          0.19.0
  GitCommit:        de40ad0
```
