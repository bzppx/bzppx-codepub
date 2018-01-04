-- Adminer 4.2.5 MySQL dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';


DROP TABLE IF EXISTS `cp_module`;
CREATE TABLE `cp_module` (
  `module_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '模块id',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '模块名称',
  `user_id` int(11) NOT NULL COMMENT 'user id',
  `modules_id` int(11) NOT NULL COMMENT '模块组 id',
  `repository_url` varchar(300) NOT NULL COMMENT 'git 仓库地址 https or ssh',
  `branch` varchar(50) NOT NULL COMMENT '分支',
  `ssh_key` text NOT NULL COMMENT 'ssh key ',
  `ssh_key_salt` text NOT NULL COMMENT 'ssh  key salt',
  `https_username` varchar(50) NOT NULL COMMENT 'https 用户名',
  `https_password` varchar(50) NOT NULL COMMENT 'https 密码',
  `code_path` varchar(200) NOT NULL COMMENT '代码发布目录',
  `code_dir_user` varchar(50) NOT NULL COMMENT '目录所属用户',
  `pre_command` text NOT NULL COMMENT '前置命令',
  `pre_command_exec_type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '前置命令执行方式, 1 同步执行，遇到错误停止;2 同步执行，遇到错误继续;3 异步执行',
  `pre_command_exec_timeout` int(11) NOT NULL DEFAULT '30' COMMENT '前置命令超时时间,单位秒',
  `post_command` text NOT NULL COMMENT '后置命令',
  `post_command_exec_type` int(11) NOT NULL DEFAULT '1' COMMENT '后置命令执行方式, 1 同步执行，遇到错误停止;2 同步执行，遇到错误继续;3 异步执行',
  `post_command_exec_timeout` int(11) NOT NULL DEFAULT '30' COMMENT '后置命令超时时间,单位秒',
  `exec_user` varchar(30) NOT NULL DEFAULT '' COMMENT '执行命令用户',
  `comment` varchar(200) NOT NULL DEFAULT '' COMMENT '备注',
  `is_delete` tinyint(1) NOT NULL DEFAULT '0' COMMENT '删除 0 否 1 是',
  `last_publish_time` int(11) NOT NULL DEFAULT '0' COMMENT '最后一次发布时间',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int(11) NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`module_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='代码模块表';


DROP TABLE IF EXISTS `cp_modules`;
CREATE TABLE `cp_modules` (
  `modules_id` int(10) NOT NULL AUTO_INCREMENT COMMENT '模块组表主键id',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '模块组名称',
  `comment` varchar(100) NOT NULL DEFAULT '' COMMENT '备注',
  `is_delete` tinyint(3) NOT NULL DEFAULT '0' COMMENT '是否删除 0 否 1 是',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`modules_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='模块组表';


DROP TABLE IF EXISTS `cp_module_node`;
CREATE TABLE `cp_module_node` (
  `module_node_id` int(10) NOT NULL AUTO_INCREMENT COMMENT '节点和模块关系表主键',
  `module_id` int(10) NOT NULL DEFAULT '0' COMMENT '模块ID',
  `node_id` int(10) NOT NULL DEFAULT '0' COMMENT '节点ID',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`module_node_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='模块和节点关系表';


DROP TABLE IF EXISTS `cp_node_nodes`;
CREATE TABLE `cp_node_nodes` (
  `node_nodes_id` int(10) NOT NULL AUTO_INCREMENT COMMENT '节点和节点组关系表主键',
  `nodes_id` int(10) NOT NULL DEFAULT '0' COMMENT '节点组ID',
  `node_id` int(10) NOT NULL DEFAULT '0' COMMENT '节点ID',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`node_nodes_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='节点和节点组关系表';


DROP TABLE IF EXISTS `cp_node`;
CREATE TABLE `cp_node` (
  `node_id` int(10) NOT NULL AUTO_INCREMENT COMMENT '节点信息表主键id',
  `ip` varchar(15) NOT NULL DEFAULT '' COMMENT '节点主机IP',
  `port` int(10) NOT NULL DEFAULT '0' COMMENT '节点主机端口',
  `comment` varchar(30) NOT NULL DEFAULT '' COMMENT '备注',
  `last_active_time` int(11) NOT NULL DEFAULT '0' COMMENT '最后存活时间',
  `is_delete` tinyint(3) NOT NULL DEFAULT '0' COMMENT '是否删除 0 否 1 是',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`node_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='节点信息表';


DROP TABLE IF EXISTS `cp_nodes`;
CREATE TABLE `cp_nodes` (
  `nodes_id` int(10) NOT NULL AUTO_INCREMENT COMMENT '节点组表主键id',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '节点组名称',
  `comment` varchar(100) NOT NULL DEFAULT '' COMMENT '备注',
  `is_delete` tinyint(3) NOT NULL DEFAULT '0' COMMENT '是否删除 0 否 1 是',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`nodes_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='节点组表';


DROP TABLE IF EXISTS `cp_task`;
CREATE TABLE `cp_task` (
  `task_id` int(10) NOT NULL AUTO_INCREMENT COMMENT '发布代码任务主键 id',
  `module_id` int(10) NOT NULL COMMENT '模块id',
  `sha1_id` varchar(200) NOT NULL DEFAULT '' COMMENT 'git commit id',
  `comment` text NOT NULL COMMENT '发布备注',
  `user_id` int(10) NOT NULL DEFAULT '0' COMMENT '用户id',
  `is_published` tinyint(3) NOT NULL DEFAULT '0' COMMENT '是否已经发布:0否，1是',
  `publish_time` int(11) NOT NULL DEFAULT '0' COMMENT '发布时间',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='发布任务表';


DROP TABLE IF EXISTS `cp_task_log`;
CREATE TABLE `cp_task_log` (
  `task_log_id` int(10) NOT NULL AUTO_INCREMENT COMMENT '发布代码任务日志表主键id',
  `task_id` int(10) NOT NULL DEFAULT '0' COMMENT '任务 id',
  `user_id` int(10) NOT NULL DEFAULT '0' COMMENT '用户 id',
  `module_id` int(10) NOT NULL DEFAULT '0' COMMENT '模块 id',
  `response` text NOT NULL COMMENT '客户端返回的结果 json 数据',
  `is_success` tinyint(3) NOT NULL DEFAULT '0' COMMENT '是否发布成功:1成功，2失败',
  `is_published` tinyint(3) NOT NULL DEFAULT '0' COMMENT '是否发布过：0未发布过，1发布过',
  `rollback_id` varchar(100) NOT NULL DEFAULT '' COMMENT '回滚用的sha1_id',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`task_log_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='任务日志表';


DROP TABLE IF EXISTS `cp_user`;
CREATE TABLE `cp_user` (
  `user_id` int(10) NOT NULL AUTO_INCREMENT COMMENT '用户表主键',
  `username` varchar(50) NOT NULL DEFAULT '' COMMENT '用户名',
  `given_name` varchar(50) NOT NULL DEFAULT '' COMMENT '姓名',
  `password` char(32) NOT NULL DEFAULT '' COMMENT '密码',
  `email` varchar(50) NOT NULL DEFAULT '' COMMENT '邮箱',
  `mobile` char(13) NOT NULL DEFAULT '' COMMENT '手机号',
  `last_ip` varchar(15) NOT NULL DEFAULT '' COMMENT '最后登录ip',
  `last_time` int(11) NOT NULL DEFAULT '0' COMMENT '最后登录时间',
  `role` tinyint(3) NOT NULL DEFAULT '0' COMMENT '1,普通用户;  2管理员;3超级管理员;',
  `is_delete` tinyint(3) NOT NULL DEFAULT '0' COMMENT '是否删除，0 否 1 是',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户表';


DROP TABLE IF EXISTS `cp_user_module`;
CREATE TABLE `cp_user_module` (
  `user_module_id` int(10) NOT NULL AUTO_INCREMENT COMMENT '用户和模块关系表 id',
  `user_id` int(10) NOT NULL DEFAULT '0' COMMENT '用户 id',
  `module_id` int(10) NOT NULL DEFAULT '0' COMMENT '模块 id',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`user_module_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户和模块对应关系表';

DROP TABLE IF EXISTS `cp_configure`;
CREATE TABLE IF NOT EXISTS `cp_configure` (
  `configure_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '配置表主键Id',
  `key` char(50) NOT NULL COMMENT '配置键',
  `value` text NOT NULL COMMENT '配置值',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
  `is_delete` int(11) NOT NULL DEFAULT '0' COMMENT '是否删除，0 否 1 是',
  PRIMARY KEY (`configure_id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='配置表';


-- ----------------------------------------------------------
-- 行为日志表
-- ----------------------------------------------------------
DROP TABLE IF EXISTS `cp_log`;
CREATE TABLE `cp_log` (
  `log_id` int(10) NOT NULL AUTO_INCREMENT COMMENT '日志id',
  `controller` char(100) NOT NULL DEFAULT '' COMMENT '控制器',
  `action` char(100) NOT NULL DEFAULT '' COMMENT '动作',
  `get` text NOT NULL COMMENT 'get参数',
  `post` text NOT NULL COMMENT 'post参数',
  `message` varchar(255) NOT NULL DEFAULT '' COMMENT '信息',
  `ip` char(100) NOT NULL DEFAULT '' COMMENT 'ip地址',
  `user_agent` char(200) NOT NULL DEFAULT '' COMMENT '用户代理',
  `referer` char(100) NOT NULL DEFAULT '' COMMENT 'referer',
  `user_id` int(10) NOT NULL DEFAULT '0' COMMENT '帐号id',
  `username` char(100) NOT NULL DEFAULT '' COMMENT '帐号名',
  `create_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`log_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='行为日志表';

-- 2017-12-28 09:25:16