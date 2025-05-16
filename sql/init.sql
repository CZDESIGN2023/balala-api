/*!40101 SET NAMES utf8 */;

SET FOREIGN_KEY_CHECKS = 0;

-- 创建一个名为 balala 的数据库
CREATE DATABASE IF NOT EXISTS balala;
-- 使用 balala 数据库
USE balala;

-- ----------------------------
-- Table structure for business_id
-- ----------------------------
DROP TABLE IF EXISTS `business_id`;
CREATE TABLE `business_id`
(
    `id`          bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `business_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '业务id',
    `max_id`      bigint unsigned DEFAULT NULL COMMENT '最大id',
    `step`        bigint unsigned DEFAULT NULL COMMENT '步长',
    `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '描述',
    `create_time` datetime                                               DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime                                               DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `table_name`  varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_business_id` (`business_id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='分布式自增主键';

-- ----------------------------
-- Table structure for config
-- ----------------------------
DROP TABLE IF EXISTS `config`;
CREATE TABLE `config`
(
    `id`            bigint unsigned NOT NULL AUTO_INCREMENT,
    `config_name`   varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '参数名称',
    `config_key`    varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '参数键名',
    `config_value`  text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '参数键值',
    `config_status` tinyint unsigned DEFAULT NULL COMMENT '状态 0 未定义 1 启用',
    `remark`        varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '备注',
    `created_at`    bigint unsigned DEFAULT NULL COMMENT '创建时间',
    `updated_at`    bigint unsigned DEFAULT NULL COMMENT '更新时间',
    `deleted_at`    bigint unsigned DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `idx_config_key` (`config_key`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='系统配置表';

-- ----------------------------
-- Table structure for file_info
-- ----------------------------
DROP TABLE IF EXISTS `file_info`;
CREATE TABLE `file_info`
(
    `id`            bigint unsigned NOT NULL AUTO_INCREMENT,
    `hash`          varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT '' COMMENT '文件hash',
    `name`          varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '参数键名',
    `typ`           tinyint unsigned DEFAULT NULL COMMENT '文件類型 (1-图片，2-视频，3-音频，4-文本)',
    `size`          bigint unsigned DEFAULT NULL COMMENT '文件大小',
    `uri`           varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '前端讀取檔案用, 動態產生: /path/name?sign=hmac-sha256簽名(附加過期timestamp)',
    `cover`         varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '封面',
    `pwd`           varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT NULL COMMENT '访问密码',
    `status`        tinyint unsigned DEFAULT NULL COMMENT '0-未定義, 1-初始化, 2-上传中, 3-成功, 4-失败, 5-处理中（转码等）, 6-待删',
    `owner`         bigint unsigned DEFAULT NULL COMMENT '上传者id',
    `meta`          json DEFAULT NULL COMMENT '文件元数据',
    `upload_typ`    tinyint unsigned DEFAULT NULL COMMENT '上傳服務器類型 (0-local 1-fastdfs，2-s3)',
    `upload_domain` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT NULL COMMENT '上傳服務器',
    `upload_md5`    varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT NULL COMMENT '上傳完成md5',
    `upload_path`   varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '上傳完成檔案路徑',
    `created_at`    bigint unsigned DEFAULT NULL COMMENT '创建时间',
    `updated_at`    bigint unsigned DEFAULT NULL COMMENT '更新时间',
    `deleted_at`    bigint unsigned DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for notify
-- ----------------------------
DROP TABLE IF EXISTS `notify`;
CREATE TABLE `notify`
(
    `id`         bigint unsigned NOT NULL AUTO_INCREMENT,
    `space_id`   bigint unsigned NOT NULL COMMENT '空间id',
    `user_id`    bigint NOT NULL COMMENT '发布评论的用户id',
    `typ`        bigint NOT NULL COMMENT '事件类型',
    `doc`        text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '内容',
    `created_at` bigint DEFAULT NULL,
    `updated_at` bigint DEFAULT NULL,
    `deleted_at` bigint DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    KEY          `notify_user_id_IDX` (`user_id`,`typ`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='通知快照';

-- ----------------------------
-- Table structure for oper_log
-- ----------------------------
DROP TABLE IF EXISTS `oper_log`;
CREATE TABLE `oper_log`
(
    `id`             bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '日志主键',
    `title`          varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT '' COMMENT '模块标题',
    `business_type`  int                                                            DEFAULT '0' COMMENT '业务类型（0其它 1新增 2修改 3删除）',
    `method`         varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT '' COMMENT '方法名称',
    `request_method` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci   DEFAULT '' COMMENT '请求方式',
    `module_type`    int                                                            DEFAULT NULL COMMENT '功能模块类型',
    `module_id`      bigint unsigned DEFAULT NULL COMMENT '功能模块对应的数据ID',
    `operator_type`  int                                                            DEFAULT '0' COMMENT '操作类别（0其它 1后台用户 2手机端用户）',
    `oper_id`        bigint unsigned DEFAULT NULL COMMENT '操作人id',
    `oper_name`      varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci   DEFAULT '' COMMENT '操作人员',
    `oper_nickname`  varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci   DEFAULT '' COMMENT '操作人员昵称',
    `oper_url`       varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT '' COMMENT '请求URL',
    `oper_ip`        varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci   DEFAULT '' COMMENT '主机地址',
    `oper_location`  varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT '' COMMENT '操作地点',
    `oper_param`     text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '请求参数',
    `oper_msg`       varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '错误消息',
    `oper_time_at`   bigint unsigned DEFAULT NULL COMMENT '操作时间',
    `created_at`     bigint unsigned DEFAULT NULL COMMENT '创建时间',
    `updated_at`     bigint unsigned DEFAULT NULL COMMENT '更新时间',
    `deleted_at`     bigint unsigned DEFAULT NULL COMMENT '删除时间',
    `space_id`       bigint unsigned DEFAULT '0' COMMENT '空间id',
    `show_type`      int                                                            DEFAULT NULL COMMENT '展示类型',
    `space_name`     varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '空间名',
    `module_flag`    bigint DEFAULT NULL COMMENT '位标记',
    PRIMARY KEY (`id`) USING BTREE,
    KEY              `oper_log_space_id_IDX` (`space_id`) USING BTREE,
    KEY              `oper_log_module_type_IDX` (`module_type`,`module_id`) USING BTREE,
    KEY              `oper_log_oper_id_IDX` (`oper_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for space
-- ----------------------------
DROP TABLE IF EXISTS `space`;
CREATE TABLE `space`
(
    `id`           bigint unsigned NOT NULL AUTO_INCREMENT,
    `user_id`      bigint unsigned NOT NULL COMMENT '创建用户id',
    `space_guid`   varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci           DEFAULT NULL COMMENT '空间Guid',
    `space_name`   varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '空间名称',
    `space_status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '空间状态;0:禁用,1:正常,2:未验证',
    `remark`       varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci          DEFAULT '' COMMENT '备注',
    `describe`     text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '描述信息',
    `notify`       tinyint                                                                DEFAULT '1' COMMENT '0不通知 1通知',
    `created_at`   bigint unsigned DEFAULT NULL COMMENT '创建时间',
    `updated_at`   bigint unsigned DEFAULT NULL COMMENT '更新时间',
    `deleted_at`   bigint unsigned DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `space_UN` (`user_id`,`space_guid`) USING BTREE,
    KEY            `space_user_id_IDX` (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for space_config
-- ----------------------------
DROP TABLE IF EXISTS `space_config`;
CREATE TABLE `space_config`
(
    `id`                bigint unsigned NOT NULL AUTO_INCREMENT,
    `space_id`          bigint DEFAULT NULL,
    `working_day`       longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
    `created_at`        bigint DEFAULT NULL,
    `updated_at`        bigint DEFAULT NULL,
    `deleted_at`        bigint DEFAULT NULL,
    `comment_deletable` int DEFAULT NULL,
    `comment_deletable_when_archived` int DEFAULT NULL COMMENT '任务归档时评论是否可删除;0:否,1:是',
    `comment_show_pos` int DEFAULT NULL COMMENT '评论展示位置;0:底部,1:单独tab',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `space_id_idx` (`space_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='空间配置表项目';

-- ----------------------------
-- Table structure for space_file_info
-- ----------------------------
DROP TABLE IF EXISTS `space_file_info`;
CREATE TABLE `space_file_info`
(
    `id`           bigint unsigned NOT NULL AUTO_INCREMENT,
    `space_id`     bigint unsigned NOT NULL COMMENT '项目空间id',
    `file_info_id` bigint unsigned NOT NULL COMMENT '标签id',
    `file_uri`     varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '文件uri访问路径',
    `file_name`    varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '文件名称, 冗余',
    `file_size`    bigint unsigned DEFAULT NULL COMMENT '文件大小',
    `source_type`  bigint unsigned NOT NULL COMMENT '关联类型 1: 任务',
    `source_id`    bigint unsigned NOT NULL COMMENT '关联id',
    `status`       tinyint DEFAULT NULL COMMENT '状态 0:未知 1:使用中 2:已删除',
    `created_at`   bigint unsigned DEFAULT NULL COMMENT '创建时间',
    `updated_at`   bigint unsigned DEFAULT NULL COMMENT '更新时间',
    `deleted_at`   bigint unsigned DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `space_tag_relation_UN` (`space_id`,`file_info_id`,`source_type`,`source_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for space_member
-- ----------------------------
DROP TABLE IF EXISTS `space_member`;
CREATE TABLE `space_member`
(
    `id`              bigint unsigned NOT NULL AUTO_INCREMENT,
    `space_id`        bigint unsigned NOT NULL COMMENT '项目空间id',
    `user_id`         bigint unsigned NOT NULL COMMENT '创建用户id',
    `role_id`         int unsigned NOT NULL COMMENT '关联角色id 1: 项目空间管理员 2: 项目空间成员',
    `remark`          varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '备注',
    `ranking`         bigint  NOT NULL                                              DEFAULT '0' COMMENT '排序值',
    `notify`          tinyint NOT NULL                                              DEFAULT '1' COMMENT '是否通知',
    `created_at`      bigint unsigned DEFAULT NULL COMMENT '创建时间',
    `updated_at`      bigint unsigned DEFAULT NULL COMMENT '更新时间',
    `deleted_at`      bigint unsigned DEFAULT NULL COMMENT '删除时间',
    `history_role_id` int unsigned DEFAULT '0' COMMENT '最新一次历史的角色id',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `space_UN` (`user_id`,`space_id`) USING BTREE,
    KEY               `space_member_space_id_IDX` (`space_id`) USING BTREE,
    KEY               `space_member_user_id_IDX` (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for space_tag
-- ----------------------------
DROP TABLE IF EXISTS `space_tag`;
CREATE TABLE `space_tag`
(
    `id`         bigint unsigned NOT NULL AUTO_INCREMENT,
    `space_id`   bigint unsigned NOT NULL COMMENT '创建用户id',
    `tag_guid`   varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci          DEFAULT NULL COMMENT '标签Guid',
    `tag_name`   varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_as_cs NOT NULL DEFAULT '' COMMENT '空间名称',
    `tag_status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT 'Tag状态;0:禁用,1:正常,2:未验证',
    `created_at` bigint unsigned DEFAULT NULL COMMENT '创建时间',
    `updated_at` bigint unsigned DEFAULT NULL COMMENT '更新时间',
    `deleted_at` bigint unsigned DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `space_tag_UN` (`space_id`,`tag_name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for space_work_item_comment
-- ----------------------------
DROP TABLE IF EXISTS `space_work_item_comment`;
CREATE TABLE `space_work_item_comment`
(
    `id`               bigint unsigned NOT NULL AUTO_INCREMENT,
    `user_id`          bigint DEFAULT NULL,
    `work_item_id`     bigint DEFAULT NULL,
    `content`          longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
    `created_at`       bigint DEFAULT NULL,
    `updated_at`       bigint DEFAULT NULL,
    `deleted_at`       bigint DEFAULT NULL,
    `refer_user_ids`   longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
    `dont_delete`      bigint DEFAULT NULL,
    `reply_comment_id` bigint DEFAULT NULL,
    `emojis`           json   DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    KEY                `work_item_id_IDX` (`work_item_id`) USING BTREE,
    KEY                `work_item_id_created_time_IDX` (`work_item_id`,`created_at`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for space_work_item_flow_role_v2
-- ----------------------------
DROP TABLE IF EXISTS `space_work_item_flow_role_v2`;
CREATE TABLE `space_work_item_flow_role_v2`
(
    `id`                 bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `space_id`           bigint unsigned DEFAULT NULL COMMENT '空间id',
    `work_item_id`       bigint unsigned DEFAULT NULL COMMENT '任务id',
    `flow_id`            bigint unsigned DEFAULT NULL COMMENT '工作项流程Id',
    `flow_template_id`   bigint unsigned DEFAULT NULL COMMENT '流程模版id',
    `work_item_role_id`  bigint                                                        DEFAULT NULL COMMENT '节点的关联角色id',
    `work_item_role_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '节点的关联角色key',
    `directors`          json                                                          DEFAULT NULL COMMENT '节点负责人',
    `created_at`         bigint unsigned DEFAULT NULL COMMENT '创建时间',
    `updated_at`         bigint unsigned DEFAULT NULL COMMENT '更新时间',
    `deleted_at`         bigint unsigned DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='工作流程角色数据';

-- ----------------------------
-- Table structure for space_work_item_flow_v2
-- ----------------------------
DROP TABLE IF EXISTS `space_work_item_flow_v2`;
CREATE TABLE `space_work_item_flow_v2`
(
    `id`                    bigint unsigned NOT NULL AUTO_INCREMENT,
    `space_id`              bigint unsigned DEFAULT '0' COMMENT '空间id',
    `work_item_id`          bigint unsigned DEFAULT NULL COMMENT '任务id',
    `flow_mode`             varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT NULL COMMENT '流程模式',
    `flow_node_uuid`        varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '节点uuid',
    `flow_node_status`      tinyint unsigned NOT NULL DEFAULT '1' COMMENT '任务状态; 0:未定义 , 1: 未开启 2: 进行中 3:已完成',
    `flow_node_code`        varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT 'P0' COMMENT '节点编码',
    `flow_node_passed`      tinyint                                                       DEFAULT NULL COMMENT '节点是否通过',
    `flow_node_reached`     tinyint                                                       DEFAULT NULL COMMENT '节点是否到达',
    `flow_mode_version`     varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '[已废弃]模式版本',
    `flow_mode_code`        varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '[已废弃]模式编码',
    `directors`             json                                                          DEFAULT NULL COMMENT '节点负责人',
    `flow_id`               bigint unsigned DEFAULT NULL COMMENT '工作项流程Id',
    `flow_key`              varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '工作项流程Key',
    `flow_template_id`      bigint unsigned DEFAULT NULL COMMENT '流程模版id',
    `flow_template_version` bigint unsigned DEFAULT NULL COMMENT '工作项流程模版版本',
    `start_at`              bigint unsigned DEFAULT NULL COMMENT '开始时间',
    `finish_at`             bigint unsigned DEFAULT NULL COMMENT '完成时间',
    `work_item_role_id`     bigint                                                        DEFAULT '0' COMMENT '节点的关联角色id',
    `work_item_role_key`    varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '节点的关联角色key',
    `created_at`            bigint unsigned DEFAULT NULL COMMENT '创建时间',
    `updated_at`            bigint unsigned DEFAULT NULL COMMENT '更新时间',
    `deleted_at`            bigint unsigned DEFAULT NULL COMMENT '删除时间',
    `plan_start_at`         bigint unsigned DEFAULT NULL COMMENT '计划排期时间-开始',
    `plan_complete_at`      bigint unsigned DEFAULT NULL COMMENT '计划排期时间-结束',
    PRIMARY KEY (`id`) USING BTREE,
    KEY                     `work_item_id_idx` (`work_item_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for space_work_item_v2
-- ----------------------------
DROP TABLE IF EXISTS `space_work_item_v2`;
CREATE TABLE `space_work_item_v2`
(
    `id`                    bigint unsigned NOT NULL AUTO_INCREMENT,
    `pid`                   bigint unsigned DEFAULT NULL COMMENT '上级id 非0表示当前为子任务',
    `space_id`              bigint unsigned NOT NULL COMMENT '项目空间id',
    `work_object_id`        bigint unsigned NOT NULL COMMENT '项目工作项id',
    `user_id`               bigint unsigned NOT NULL COMMENT '创建用户id',
    `flow_id`               bigint unsigned DEFAULT NULL COMMENT '工作项流程Id',
    `flow_key`              varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci           DEFAULT NULL COMMENT '工作项流程Key',
    `work_item_flow_id`     bigint unsigned DEFAULT NULL COMMENT '工作项流程Id',
    `work_item_guid`        varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci            DEFAULT NULL COMMENT '工作项(工作模块)Guid',
    `work_item_name`        varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `work_item_status`      varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NOT NULL COMMENT '任务状态Val',
    `work_item_status_key`  varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NOT NULL DEFAULT '' COMMENT '任务状态Key',
    `work_item_status_id`   bigint unsigned NOT NULL DEFAULT '0' COMMENT '任务状态Id',
    `work_item_type_id`     bigint                                                                  DEFAULT NULL COMMENT '工作项类型Id',
    `work_item_type_key`    varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci           DEFAULT NULL COMMENT '工作项类型key',
    `doc`                   json                                                                    DEFAULT NULL,
    `flow_template_id`      bigint unsigned DEFAULT NULL COMMENT '工作项流程模版id',
    `flow_template_version` bigint unsigned DEFAULT NULL COMMENT '工作项流程模版版本',
    `work_item_type`        tinyint                                                                 DEFAULT NULL COMMENT '[已废弃]工作项类型 ',
    `flow_mode`             varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci            DEFAULT NULL COMMENT '流程模式 stateflow | workflow',
    `work_item_flow_key`    varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci           DEFAULT NULL COMMENT '工作项流程Key',
    `flow_mode_version`     varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci           DEFAULT NULL COMMENT '[已废弃] 模式版本',
    `flow_mode_code`        varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci           DEFAULT NULL COMMENT '[已废弃] 模式编码',
    `child_num`             int                                                                     DEFAULT '0' COMMENT '子任务数量',
    `last_status_at`        bigint unsigned DEFAULT '0' COMMENT '历史状态更新时间',
    `last_status_key`       varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NOT NULL COMMENT '历史状态Key',
    `last_status_id`        bigint unsigned NOT NULL DEFAULT '0' COMMENT '历史状态Id',
    `last_status`           varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci           DEFAULT NULL COMMENT '历史状态',
    `is_restart`            tinyint unsigned DEFAULT '0' COMMENT '是否为重启任务',
    `restart_at`            bigint unsigned DEFAULT '0' COMMENT '重启时间',
    `icon_flags`            int unsigned NOT NULL DEFAULT '0',
    `restart_user_id`       bigint unsigned DEFAULT '0' COMMENT '重启用户id',
    `comment_num`           int unsigned NOT NULL DEFAULT '0' COMMENT '评论数',
    `resume_at`             bigint unsigned DEFAULT '0' COMMENT '任务恢复时间',
    `version_id`            bigint unsigned DEFAULT '0' COMMENT '版本信息id',
    `reason`                json                                                                    DEFAULT NULL COMMENT '原因',
    `count_at`              bigint DEFAULT NULL COMMENT '计时时间',
    `created_at`            bigint unsigned DEFAULT NULL COMMENT '创建时间',
    `updated_at`            bigint unsigned DEFAULT NULL COMMENT '更新时间',
    `deleted_at`            bigint unsigned DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`) USING BTREE,
    KEY                     `space_work_item_space_id_IDX` (`space_id`) USING BTREE,
    KEY                     `space_work_item_user_id_IDX` (`user_id`) USING BTREE,
    KEY                     `space_work_item_work_object_id_IDX` (`work_object_id`) USING BTREE,
    KEY                     `space_work_item_pid_IDX` (`pid`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for space_work_object
-- ----------------------------
DROP TABLE IF EXISTS `space_work_object`;
CREATE TABLE `space_work_object`
(
    `id`                 bigint unsigned NOT NULL AUTO_INCREMENT,
    `space_id`           bigint unsigned NOT NULL COMMENT '项目空间id',
    `user_id`            bigint unsigned NOT NULL COMMENT '创建用户id',
    `work_object_guid`   varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci          DEFAULT NULL COMMENT '工作项(工作模块)Guid',
    `work_object_name`   varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_as_cs NOT NULL DEFAULT '' COMMENT '工作项(工作模块)名称',
    `work_object_status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '空间状态;0:禁用,1:正常,2:未验证',
    `remark`             varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci         DEFAULT '' COMMENT '备注',
    `describe`           varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci         DEFAULT '' COMMENT '描述信息',
    `ranking`            bigint                                                                DEFAULT '0' COMMENT '排序值',
    `created_at`         bigint unsigned DEFAULT NULL COMMENT '创建时间',
    `updated_at`         bigint unsigned DEFAULT NULL COMMENT '更新时间',
    `deleted_at`         bigint unsigned DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `space_work_object_2_UN` (`work_object_name`,`user_id`,`space_id`) USING BTREE,
    UNIQUE KEY `space_work_object_UN` (`user_id`,`space_id`,`work_object_guid`) USING BTREE,
    KEY                  `space_work_object_space_id_IDX` (`space_id`) USING BTREE,
    KEY                  `space_work_object_user_id_IDX` (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for space_work_version
-- ----------------------------
DROP TABLE IF EXISTS `space_work_version`;
CREATE TABLE `space_work_version`
(
    `id`             bigint unsigned NOT NULL AUTO_INCREMENT,
    `space_id`       bigint unsigned NOT NULL COMMENT '项目空间id',
    `version_key`    varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '工作项（版本） KEY',
    `version_name`   varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '工作项(工作模块)名称',
    `version_status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '空间状态;0:禁用,1:正常,2:未验证',
    `remark`         varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci         DEFAULT '' COMMENT '备注',
    `ranking`        bigint                                                                DEFAULT '0' COMMENT '排序值',
    `created_at`     bigint unsigned DEFAULT NULL COMMENT '创建时间',
    `updated_at`     bigint unsigned DEFAULT NULL COMMENT '更新时间',
    `deleted_at`     bigint unsigned DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='项目空间-版本信息表';

-- ----------------------------
-- Table structure for third_pf_account
-- ----------------------------
DROP TABLE IF EXISTS `third_pf_account`;
CREATE TABLE `third_pf_account`
(
    `id`              bigint NOT NULL AUTO_INCREMENT,
    `user_id`         bigint                                                        DEFAULT NULL COMMENT '用户id',
    `pf_code`         tinyint                                                       DEFAULT NULL COMMENT '平台枚举 3IM',
    `pf_name`         varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '平台名',
    `pf_user_key`     varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '平台用户唯一标识符',
    `pf_user_id`      bigint                                                        DEFAULT NULL COMMENT '平台用户id',
    `pf_user_name`    varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '平台用户名',
    `pf_user_account` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '平台用户账号',
    `created_at`      bigint                                                        DEFAULT NULL COMMENT '创建时间',
    `notify`          tinyint                                                       DEFAULT '1' COMMENT '通知',
    `updated_at`      bigint                                                        DEFAULT NULL,
    `deleted_at`      bigint                                                        DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `user_id_pf_code` (`user_id`,`pf_code`) USING BTREE,
    UNIQUE KEY `pf_user_key_pf_code_idx` (`pf_user_key`,`pf_code`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`
(
    `id`              bigint unsigned NOT NULL AUTO_INCREMENT,
    `user_name`       varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_as_cs           DEFAULT NULL COMMENT '用户名',
    `mobile`          varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci           DEFAULT NULL COMMENT '中国手机不带国家代码，国际手机号格式为：国家代码-手机号',
    `user_nickname`   varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NOT NULL DEFAULT '' COMMENT '用户昵称',
    `user_pinyin`     varchar(127) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci          DEFAULT NULL COMMENT '用户昵称拼音',
    `user_password`   varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '登录密码;cmf_password加密',
    `user_salt`       char(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci     NOT NULL COMMENT '加密盐',
    `user_status`     tinyint unsigned NOT NULL DEFAULT '1' COMMENT '用户状态;0:禁用,1:正常,2:未验证',
    `user_email`      varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci          DEFAULT NULL COMMENT '用户登录邮箱',
    `sex`             tinyint                                                                DEFAULT '0' COMMENT '性别;0:保密,1:男,2:女',
    `avatar`          varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户头像',
    `remark`          varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci          DEFAULT '' COMMENT '备注',
    `describe`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci          DEFAULT '' COMMENT '描述信息',
    `last_login_ip`   varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci           DEFAULT '' COMMENT '最后登录ip',
    `last_login_time` bigint unsigned DEFAULT NULL COMMENT '最后登录时间',
    `role`            tinyint                                                       NOT NULL DEFAULT '0' COMMENT '角色 0普通成员 50系统管理员 100超级管理员',
    `created_at`      bigint unsigned DEFAULT NULL COMMENT '创建时间',
    `updated_at`      bigint unsigned DEFAULT NULL COMMENT '更新时间',
    `deleted_at`      bigint unsigned DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `user_login` (`user_name`,`deleted_at`) USING BTREE,
    UNIQUE KEY `email` (`user_email`) USING BTREE,
    UNIQUE KEY `mobile` (`mobile`) USING BTREE,
    KEY               `user_nickname` (`user_nickname`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for user_login_log
-- ----------------------------
DROP TABLE IF EXISTS `user_login_log`;
CREATE TABLE `user_login_log`
(
    `id`                  bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '访问ID',
    `login_user_id`       bigint unsigned DEFAULT NULL COMMENT '登陆账号ID',
    `login_user_name`     varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT '' COMMENT '登录账号',
    `login_user_nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT '' COMMENT '登录账号昵称',
    `ipaddr`              varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT '' COMMENT '登录IP地址',
    `login_location`      varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '登录地点',
    `browser`             varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '浏览器类型',
    `os`                  varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT '' COMMENT '操作系统',
    `status`              tinyint                                                       DEFAULT '0' COMMENT '登录状态（0成功 1失败）',
    `msg`                 varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '提示消息',
    `module`              varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT '' COMMENT '登录模块',
    `login_at`            bigint unsigned DEFAULT NULL COMMENT '登录时间',
    `created_at`          bigint unsigned DEFAULT NULL COMMENT '创建时间',
    `updated_at`          bigint unsigned DEFAULT NULL COMMENT '更新时间',
    `deleted_at`          bigint unsigned DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for work_flow
-- ----------------------------
DROP TABLE IF EXISTS `work_flow`;
CREATE TABLE `work_flow`
(
    `id`                bigint unsigned NOT NULL AUTO_INCREMENT,
    `uuid`              varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '唯一识别码',
    `user_id`           bigint                                                        DEFAULT NULL COMMENT '创建人',
    `space_id`          bigint unsigned NOT NULL DEFAULT '0' COMMENT '空间id',
    `work_item_type_id` bigint unsigned NOT NULL COMMENT '工作项id',
    `name`              varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '工作流名称',
    `key`               varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'key',
    `flow_mode`         varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NOT NULL COMMENT '流程模式',
    `version`           int unsigned NOT NULL DEFAULT '0' COMMENT '当前最大版本',
    `is_sys`            tinyint                                                       DEFAULT '0' COMMENT '是否系统预设',
    `status`            tinyint                                                       DEFAULT NULL COMMENT '状态 0:禁用 1:启用',
    `ranking`           bigint unsigned DEFAULT NULL COMMENT '排序',
    `last_template_id`  bigint                                                        DEFAULT NULL COMMENT '最大版本流程模版id',
    `created_at`        bigint unsigned DEFAULT NULL COMMENT '创建时间',
    `updated_at`        bigint unsigned DEFAULT NULL COMMENT '更新时间',
    `deleted_at`        bigint unsigned DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`) USING BTREE,
    KEY                 `work_item_id_idx` (`work_item_type_id`) USING BTREE,
    KEY                 `space_id_idx` (`space_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='流程模版表-空间';

-- ----------------------------
-- Table structure for work_flow_template
-- ----------------------------
DROP TABLE IF EXISTS `work_flow_template`;
CREATE TABLE `work_flow_template`
(
    `id`                bigint unsigned NOT NULL AUTO_INCREMENT,
    `uuid`              varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '唯一识别码',
    `user_id`           bigint  DEFAULT NULL COMMENT '创建人',
    `space_id`          bigint unsigned NOT NULL DEFAULT '0' COMMENT '空间id',
    `work_item_type_id` bigint unsigned NOT NULL COMMENT '工作项id',
    `work_flow_id`      bigint unsigned NOT NULL DEFAULT '0' COMMENT '工作流id',
    `flow_mode`         varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NOT NULL COMMENT '流程模式',
    `version`           int unsigned NOT NULL DEFAULT '1' COMMENT '版本',
    `setting`           json    DEFAULT NULL COMMENT '配置信息',
    `status`            tinyint DEFAULT NULL COMMENT '状态: 0 禁言 1启用',
    `created_at`        bigint unsigned DEFAULT NULL COMMENT '创建时间',
    `updated_at`        bigint unsigned DEFAULT NULL COMMENT '更新时间',
    `deleted_at`        bigint unsigned DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`) USING BTREE,
    KEY                 `work_item_id_idx` (`work_item_type_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='流程模版配置表';

-- ----------------------------
-- Table structure for work_item_role
-- ----------------------------
DROP TABLE IF EXISTS `work_item_role`;
CREATE TABLE `work_item_role`
(
    `id`                bigint unsigned NOT NULL AUTO_INCREMENT,
    `uuid`              varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '唯一识别码',
    `space_id`          bigint unsigned NOT NULL COMMENT '项目空间id',
    `user_id`           bigint                                                        DEFAULT NULL COMMENT '创建人',
    `work_item_type_id` bigint unsigned NOT NULL COMMENT '创建用户id',
    `key`               varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `name`              varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT NULL COMMENT '角色名称',
    `status`            tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态;0:禁用,1:正常,2:未验证',
    `ranking`           bigint                                                        DEFAULT '0' COMMENT '排序值',
    `is_sys`            tinyint                                                       DEFAULT NULL COMMENT '系统预设 0:否 1:是',
    `flow_scope`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '可作用的流程范围  空(表示全部) |state_flow | work_flow',
    `created_at`        bigint unsigned DEFAULT NULL COMMENT '创建时间',
    `updated_at`        bigint unsigned DEFAULT NULL COMMENT '更新时间',
    `deleted_at`        bigint unsigned DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='工作项角色-空间';

-- ----------------------------
-- Table structure for work_item_status
-- ----------------------------
DROP TABLE IF EXISTS `work_item_status`;
CREATE TABLE `work_item_status`
(
    `id`                bigint unsigned NOT NULL AUTO_INCREMENT,
    `uuid`              varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '唯一识别码',
    `space_id`          bigint unsigned NOT NULL COMMENT '项目空间id',
    `user_id`           bigint                                                        DEFAULT NULL COMMENT '创建人',
    `work_item_type_id` bigint unsigned NOT NULL COMMENT '工作项类型id',
    `key`               varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `val`               varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `name`              varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT NULL COMMENT '角色名称',
    `status_type`       tinyint                                                       DEFAULT NULL COMMENT '状态类型 1:起始 2:过程 3:归档',
    `status`            tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态;0:禁用,1:正常,2:未验证',
    `ranking`           bigint                                                        DEFAULT '0' COMMENT '排序值',
    `created_at`        bigint unsigned DEFAULT NULL COMMENT '创建时间',
    `updated_at`        bigint unsigned DEFAULT NULL COMMENT '更新时间',
    `deleted_at`        bigint unsigned DEFAULT NULL COMMENT '删除时间',
    `is_sys`            tinyint                                                       DEFAULT NULL COMMENT '是否系统预设 0:否 1:是',
    `flow_scope`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '可作用的流程范围  空(表示全部) |state_flow | work_flow',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='工作项状态-全局';

-- ----------------------------
-- Table structure for work_item_type
-- ----------------------------
DROP TABLE IF EXISTS `work_item_type`;
CREATE TABLE `work_item_type`
(
    `id`         bigint unsigned NOT NULL AUTO_INCREMENT,
    `uuid`       varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '唯一标识',
    `space_id`   bigint unsigned DEFAULT '0' COMMENT '空间id',
    `user_id`    bigint                                                                 DEFAULT NULL COMMENT '创建人',
    `name`       varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NOT NULL DEFAULT '' COMMENT '名称',
    `key`        varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NOT NULL DEFAULT '' COMMENT '编码[设置后不允许修改]',
    `flow_mode`  varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci          DEFAULT NULL COMMENT '工作流模式',
    `ranking`    bigint unsigned DEFAULT '0' COMMENT '排序',
    `is_sys`     tinyint                                                                DEFAULT '0' COMMENT '是否系统预设',
    `status`     tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态;0:禁用,1:正常,2:未验证',
    `created_at` bigint unsigned DEFAULT NULL COMMENT '创建时间',
    `updated_at` bigint unsigned DEFAULT NULL COMMENT '更新时间',
    `deleted_at` bigint unsigned DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='任务类型-空间';

-- ----------------------------
-- Table structure for space_global_view
-- ----------------------------
DROP TABLE IF EXISTS `space_global_view`;
CREATE TABLE `space_global_view`
(
    `id`           bigint NOT NULL AUTO_INCREMENT,
    `space_id`     bigint                                                         DEFAULT NULL,
    `key`          varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT NULL,
    `name`         varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT NULL,
    `type`         tinyint                                                        DEFAULT NULL COMMENT '1系统2全局3用户',
    `query_config` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
    `table_config` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
    `created_at`   bigint                                                         DEFAULT NULL,
    `updated_at`   bigint                                                         DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY            `space_id` (`space_id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Table structure for space_user_view
-- ----------------------------
DROP TABLE IF EXISTS `space_user_view`;
CREATE TABLE `space_user_view`
(
    `id`           bigint NOT NULL AUTO_INCREMENT,
    `user_id`      bigint                                                         DEFAULT NULL,
    `space_id`     bigint                                                         DEFAULT NULL,
    `key`          varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT NULL,
    `name`         varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT NULL,
    `type`         tinyint                                                        DEFAULT NULL COMMENT '1系统2全局3用户',
    `outer_id`     bigint                                                         DEFAULT NULL COMMENT '当type为1|2时指向view_space',
    `query_config` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
    `table_config` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
    `ranking`      int                                                            DEFAULT NULL,
    `status`       tinyint                                                        DEFAULT NULL,
    `created_at`   bigint                                                         DEFAULT NULL,
    `updated_at`   bigint                                                         DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY            `space_id_user_id` (`space_id`,`user_id`) USING BTREE,
    KEY            `user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Table structure for user_config
-- ----------------------------
DROP TABLE IF EXISTS `user_config`;
CREATE TABLE `user_config`
(
    `id`         bigint NOT NULL AUTO_INCREMENT,
    `user_id`    bigint                                                        DEFAULT NULL,
    `key`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
    `value`      varchar(1023)                                                 DEFAULT NULL,
    `created_at` bigint                                                        DEFAULT NULL,
    `updated_at` bigint                                                        DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY          `user_id_key` (`user_id`,`key`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

SET
FOREIGN_KEY_CHECKS = 1;
