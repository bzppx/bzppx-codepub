
![brand](/docs/guide_res/images/brand.png)   
暴走皮皮虾之代码发布系统,是现代的持续集成发布系统,由后台管理系统和agent两部分组成,一个运行着的agent就是一个节点,本系统并不是造轮子,是"鸟枪"到"大炮"的创新,对"前朝遗老"的革命.

[![stable](https://img.shields.io/badge/stable-stable-green.svg)](https://github.com/bzppx/bzppx-codepub/) [![license](https://img.shields.io/github/license/bzppx/bzppx-codepub.svg?style=plastic)]() [![download_count](https://img.shields.io/github/downloads/bzppx/bzppx-codepub/total.svg?style=plastic)](https://github.com/bzppx/bzppx-codepub/releases) [![download](https://img.shields.io/github/release/bzppx/bzppx-codepub.svg?style=plastic)](https://github.com/bzppx/bzppx-codepub/releases)   

# 特点
- 基于golang编写,站在巨人肩膀上,充分利用golang的协程,channel还有高并发的特点,甩其它基于虚拟机或者解释性语言编写的发布系统几条街.
- 极速部署,其实部署是不存在的,本系统只需要下载二进制文件执行即可,想用立刻就能用.
- GIT仓库支持,可以远程拉取GIT代码发布到集群节点上.
- 集群发布,一个GIT仓库可以配置发布到N个节点的集群.
- 命令支持,发布代码的前后,都可以自己定义执行一段系统命令,比如:重启程序.
- 构建发布,java,c,c++,golang等编写的程序需要先编译"发布包",然后分发"发布包"到节点集群.
- 封版管理,特有的封版本功能,同时具备封版白名单.公司在一些特殊时期,比如电商公司在某些促销或其它节日活动的时候,为了保证线上服务,往往需要对整个发布代码功能暂停,那么封版功能就十分好用了.
- 高性能,发布代码的速度极快,对系统环境的要求很低.
- 跨平台,Linux,windows,mac,树梅派,路由器等等...
- 人性化的权限控制,一般一个部门的老大是超级管理员角色,老大只需要管理谁是管理员即可,管理员只需要管理用户即可,用户能发布代码.每层的权限系统都做控制.
- 界面优美,交互简单,更符合现在的大众web审美.
- 公告功能,有时候有些重要事情需要告知使用发布系统的开发者,那么公告功能就很好的解决了这个问题.
- 集成外部登录,本系统支持通过外部系统认证用户,比如与公司的LDAP登录融合,只需要根据我们的开发文档花费几十分钟写个HTTP API接口即可.

# 安装
## 1. codepub 安装

打开 https://github.com/bzppx/bzppx-codepub/releases 找到对应平台的版本下载编译好的压缩包

```
# 创建目录
$ mkdir codepub
$ cd codepub
# 以 linux amd64 为例，下载版本 0.8 压缩包
$ wget https://github.com/bzppx/bzppx-codepub/releases/download/v0.8/bzppx-codepub-linux-amd64.tar.gz
# 解压到当前目录
$ tar -zxvf bzppx-codepub-linux-amd64.tar.gz
# 进入程序安装目录
$ cd install
# 执行安装程序，默认端口为 8090，指定其他端口加参数 --port=8087
$ ./install
# 浏览器访问 http://ip:8090 进入安装界面，完成安装配置
# Ctrl + C 停止 install 程序, 启动 codepub 管理后台
$ cd ..
$ ./codepub --conf conf/codepub.conf
```

## 2. codepub-agent 安装
请查看 https://github.com/bzppx/bzppx-agent-codepub

## 3. nginx 配置反向代理
```
upstream frontends {
    server 127.0.0.1:8088; # codepub 监听的ip:port
}
server {
    listen      80;
    server_name codepub.com www.codepub.com;
    location / {
        proxy_pass_header Server;
        proxy_set_header Host $http_host;
        proxy_redirect off;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Scheme $scheme;
        proxy_pass http://frontends;
    }
    # 静态资源交由nginx管理
    location /static {
        root        /www/bzppx-codepub; # codepub 的根目录
        expires     1d;
        add_header  Cache-Control public;
        access_log  off;
    }
}
```

# 同类软件对比
| - | 语言 | 部署 | 稳定性 | 系统要求 | 平台覆盖 | 发布速度 | 发布配置 | 邮件通知 | 封版 | 权限 | 公告 | 界面
| :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---:|  :---:
| `Jenkins` | java | 复杂 | 中 | 高 | 一般 | 很慢 | 灵活| √ | × | √ | × | 丑陋
| `CodePub` | golang | 1分钟 | 高 | 低 | 广泛 | 极快 | 十分灵活 | √ | √ | √(更细) | √ | 优美

# 系统介绍

## 1.用户功能

| - | 用户管理 | 设置管理员 | 发布代码 | 回滚代码 | 封版 | 公告管理 | 项目管理 | 节点管理
| ------ | ------ | ------ | ------ | ------ | ------ | ------ | ------ | ------
| `超级管理员` | √ | √ | √ | √ | √ | √ | √ | √
| `管理员` | √ | × | √ | √ | √ | √ | √ | √ | √
| `普通用户` | × | × | √ | √ | × | × | × | ×

## 2.系统界面,先睹为快

### 2.1 安装
![install](/docs/guide_res/images/install.png)
### 2.2 登录
![login](/docs/guide_res/images/login.png)
### 2.2 面板
![login](/docs/guide_res/images/index.png)
### 2.3 添加节点
![login](/docs/guide_res/images/add-node.png)
### 2.4 添加项目
![login](/docs/guide_res/images/add-project.png)
### 2.5 发布代码
![login](/docs/guide_res/images/publish.png)
### 2.6 节点进度
![login](/docs/guide_res/images/task.png)

# 开发

环境要求：go 1.8
```
$ git clone https://github.com/bzppx/bzppx-codepub.git
$ cd bzppx-codepub
$ go build ./
```

# 反馈

欢迎提交意见和代码 https://github.com/bzppx/bzppx-codepub/issues
官方交流 QQ 群：547481058

## License

MIT

谢谢
---
Create By BZPPX