
![brand](/docs/guide_res/images/brand.png)   
暴走皮皮虾之代码发布系统,是现代的持续集成的发布系统,由后台管理系统和agent两部分组成,一个运行着的agent就是一个节点,本系统并不是造轮子,是"鸟枪"到"大炮"的创新,对"前朝遗老"的革命.

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


# 0. 系统开发中,请保持关注哟!

# 1. 一些对比

| - | 语言 | 部署 | 稳定性 | 系统要求 | 平台覆盖 | 发布速度 | 发布配置 | 邮件通知 | 封版 | 权限 | 公告
| :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---:
| `Jenkins` | java | 复杂 | 中 | 高 | 一般 | 很慢 | 灵活| √ | × | √ | ×
| `CodePub` | golang | 1分钟 | 高 | 低 | 广泛 | 极快 | 十分灵活 | √ | √ | √(更细) | √

# 2. 系统介绍

## 2.1 用户功能

| - | 用户管理 | 设置管理员 | 发布代码 | 回滚代码 | 封版 | 公告管理 | 模块管理 | 节点管理
| ------ | ------ | ------ | ------ | ------ | ------ | ------ | ------ | ------ 
| `超级管理员` | √ | √ | √ | √ | √ | √ | √ | √
| `管理员` | √ | × | √ | √ | √ | √ | √ | √ | √
| `普通用户` | × | × | √ | √ | × | × | × | ×

## 2.2 系统界面,先睹为快

### 2.2.1 登录
![login](/docs/guide_res/images/login.png)
### 2.2.2 添加代码模块
![login](/docs/guide_res/images/module-add.png)
