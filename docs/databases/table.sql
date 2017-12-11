-- -------------
-- 用户表
-- -------------
DROP TABLE IF EXISTS `cp_user`;
CREATE TABLE IF NOT EXISTS `cp_user` (
  `user_id` int(10) NOT NULL AUTO_INCREMENT COMMENT '用户表主键',
  `username` varchar(50) NOT NULL DEFAULT '' COMMENT '用户名',
  `given_name` varchar(50) NOT NULL DEFAULT '' COMMENT '姓名',
  `password` char(32) NOT NULL DEFAULT '' COMMENT '密码',
  `email` varchar(50) NOT NULL DEFAULT '' COMMENT 'email地址',
  `mobile` char(13) NOT NULL DEFAULT '' COMMENT '手机号',
  `last_ip` varchar(15) NOT NULL DEFAULT '' COMMENT '最后登录ip',
  `last_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '最后登录时间',
  `role` tinyint(3) NOT NULL DEFAULT '0' COMMENT '0,普通用户; 1超级管理员; 2管理员',
  `status` tinyint(3) NOT NULL DEFAULT '0' COMMENT '状态，0：正常 -1：删除',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '更新时间',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='用户表';


-- -------------
-- 节点表
-- -------------
DROP TABLE IF EXISTS `cp_node`;
CREATE TABLE IF NOT EXISTS `cp_node` (
  `node_id` int(10) NOT NULL AUTO_INCREMENT COMMENT '节点信息表主键id',
  `ip` varchar(15) NOT NULL DEFAULT '' COMMENT '节点主机IP',
  `port` int(10) NOT NULL  DEFAULT '0' COMMENT '节点主机端口',
  `comment` varchar(30) NOT NULL DEFAULT '' COMMENT '备注',
  `status` tinyint(3) NOT NULL DEFAULT '0' COMMENT '状态，0：正常，-1：删除',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '更新时间',
  PRIMARY KEY (`node_id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='节点信息表';


-- -------------
-- 节点组表
-- -------------
DROP TABLE IF EXISTS `cp_nodes`;
CREATE TABLE IF NOT EXISTS `cp_nodes` (
  `nodes_id` int(10) NOT NULL AUTO_INCREMENT COMMENT '节点组表主键id',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '节点组名称',
  `comment` varchar(100) NOT NULL DEFAULT '' COMMENT '备注',
  `status` tinyint(3) NOT NULL DEFAULT '0' COMMENT '状态，0：正常，-1：删除',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '更新时间',
  PRIMARY KEY (`nodes_id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='节点组表';


-- -------------
-- 节点-节点组关系表
-- -------------
DROP TABLE IF EXISTS `cp_node_nodes`;
CREATE TABLE IF NOT EXISTS `cp_node_nodes` (
  `node_nodes_id` int(10) NOT NULL AUTO_INCREMENT COMMENT '节点和组的关系表主键id',
  `node_id` int(10) NOT NULL DEFAULT '' COMMENT '节点ID',
  `nodes_id` int(10) NOT NULL DEFAULT '' COMMENT '节点组ID',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`node_nodes_id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='节点和组的关系表' ;



-- -------------
-- 模块表
-- -------------
DROP TABLE IF EXISTS `cp_module`;
CREATE TABLE IF NOT EXISTS `cp_module` (
  `module_id` int(10) NOT NULL AUTO_INCREMENT COMMENT '代码模块表主键id',
  `repository_url` varchar(300) DEFAULT '' NOT NULL COMMENT 'git仓库地址，支持https和ssh',
  `branch` varchar(50) NOT NULL DEFAULT '' COMMENT '发布分支',
  `ssh_key` text NOT NULL COMMENT 'ssh私钥内容',
  `ssh_key_salt` text NOT NULL COMMENT 'ssh私钥密码',
  `ssh_password` varchar(50) NOT NULL DEFAULT '' COMMENT 'ssh密码',
  `https_username` varchar(50) NOT NULL DEFAULT '' COMMENT 'https认证用户名',
  `https_password` varchar(50) NOT NULL DEFAULT '' COMMENT 'https认证密码',
  `code_path` varchar(200) NOT NULL DEFAULT '' COMMENT '代码在节点上的存储路径',
  `code_user` varchar(50) NOT NULL DEFAULT '' COMMENT '节点上存储代码目录的所有者',
  `pre_command` text NOT NULL COMMENT '执行的前置命令',
  `pre_command_exec_type` tinyint(3) NOT NULL COMMENT '前置命令执行类型，1：同步执行,遇到错误停止;2：同步执行,遇到错误继续;3：异步执行',
  `pre_command_exec_timeout` int(10) NOT NULL COMMENT '前置命令执行超时时间，单位秒',
  `post_command` text NOT NULL COMMENT '执行的后置命令',
  `post_command_exec_type` tinyint(3) NOT NULL COMMENT '后置命令执行类型，1：同步执行,遇到错误停止;2：同步执行,遇到错误继续;3：异步执行',
  `post_command_exec_timeout` int(10) NOT NULL COMMENT '后置命令执行超时时间，单位秒',
  `exec_user` varchar(30) NOT NULL DEFAULT '' COMMENT '执行前置或后置命令的用户',
  `comment` varchar(200) NOT NULL DEFAULT '' COMMENT '备注',
  `user_id` int(10) NOT NULL DEFAULT '' COMMENT '用户表主键id',
  `status` tinyint(3) NOT NULL DEFAULT '0' COMMENT '状态，0：正常 -1：删除',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '更新时间',
  PRIMARY KEY (`module_id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='模块表';

-- -------------
-- 模块组表
-- -------------
DROP TABLE IF EXISTS `cp_modules`;
CREATE TABLE IF NOT EXISTS `cp_modules` (
  `modules_id` int(10) NOT NULL AUTO_INCREMENT COMMENT '模块组表主键id',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '模块组名称',
  `comment` varchar(100) NOT NULL DEFAULT '' COMMENT '备注',
  `status` tinyint(3) NOT NULL DEFAULT '0' COMMENT '状态，0：正常，-1：删除',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '更新时间',
  PRIMARY KEY (`modules_id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='模块组表';



-- -------------
-- 模块-模块组关系表
-- -------------
DROP TABLE IF EXISTS `cp_module_modules`;
CREATE TABLE IF NOT EXISTS `cp_module_modules` (
  `module_modules_id` int(10) NOT NULL AUTO_INCREMENT COMMENT '节点和组的关系表主键id',
  `module_id` int(10) NOT NULL DEFAULT '' COMMENT '模块ID',
  `modules_id` int(10) NOT NULL DEFAULT '' COMMENT '模块组ID',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`node_nodes_id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='节点和组的关系表' ;




-- -------------
-- 发布任务表
-- -------------
DROP TABLE IF EXISTS `cp_task`;
CREATE TABLE IF NOT EXISTS `cp_task` (
  `task_id` int(10) NOT NULL AUTO_INCREMENT COMMENT '发布代码任务主键ID',
  `module_id` int(10) NOT NULL COMMENT '模块ID',
  `sha1_id` varchar(200) NOT NULL DEFAULT '' COMMENT 'git commit ID',
  `comment` text NOT NULL COMMENT '发布备注',
  `user_id` int(10) NOT NULL DEFAULT '0' COMMENT '用户表主键id',
  `is_published` tinyint(3) NOT NULL DEFAULT '0' COMMENT '是否已经发布:0否，1是',
  `publish_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '发布时间',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`task_id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='发布任务表';


-- -------------
-- 任务日志表
-- -------------
DROP TABLE IF EXISTS `cp_task_log`;
CREATE TABLE IF NOT EXISTS `ac_task_log` (
  `task_log_id` int(10) NOT NULL AUTO_INCREMENT COMMENT '发布代码任务日志表主键ID',
  `response` text NOT NULL COMMENT '客户端返回的结果json数据',
  `task_id` int(10) NOT NULL DEFAULT '0' COMMENT '发布代码任务表主键ID',
  `client_id` int(10) NOT NULL DEFAULT '0' COMMENT '客户端表主键ID',
  `is_success` tinyint(3) NOT NULL DEFAULT '0' COMMENT '是否发布成功:1成功，2失败',
  `is_published` tinyint(3) NOT NULL DEFAULT '0' COMMENT '是否发布过：0未发布过，1发布过',
  `user_id` int(10) NOT NULL DEFAULT '0' COMMENT '用户表主键ID',
  `rollback_id` varchar(100) NOT NULL DEFAULT '' COMMENT '回滚用的sha1_id',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '更新时间',
  PRIMARY KEY (`task_log_id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='任务日志表';



-- -------------
-- 行为日志表
-- -------------
DROP TABLE IF EXISTS `cp_log`;
CREATE TABLE IF NOT EXISTS `cp_log` (
  `log_id` int(10) NOT NULL AUTO_INCREMENT COMMENT '日志表主键ID',
  `admin_id` int(10) NOT NULL DEFAULT '' COMMENT '用户表主键ID',
  `message` text NOT NULL COMMENT '日志内容',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`log_id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='行为日志表';


-- -------------
-- 节点-模块关系表
-- -------------
DROP TABLE IF EXISTS `cp_node_module`;
CREATE TABLE IF NOT EXISTS `cp_node_module` (
  `node_module_id` int(10) NOT NULL AUTO_INCREMENT COMMENT '节点和模块关系表主键',
  `node_id` int(10) NOT NULL DEFAULT '0' COMMENT '节点ID',
  `module_id` int(10) NOT NULL DEFAULT '0' COMMENT '模块ID',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`node_module_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='节点和模块关系表';


-- -------------
-- 用户-模块关系表
-- -------------
DROP TABLE IF EXISTS `cp_user_module`;
CREATE TABLE IF NOT EXISTS `cp_user_module` (
  `user_module_id` int(10) NOT NULL AUTO_INCREMENT COMMENT '用户和模块关系表主键',
  `user_id` int(10) NOT NULL DEFAULT '0' COMMENT '用户主键id',
  `module_id` int(10) NOT NULL DEFAULT '0' COMMENT '模块主键id',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`user_module_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户和模块关系表';
