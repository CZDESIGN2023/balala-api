/*!40101 SET NAMES utf8 */;

USE balala;

-- 初始化超管账号
INSERT INTO `balala`.`user` (`id`, `user_name`, `user_nickname`, `user_pinyin`, `user_password`, `user_salt`, `user_status`, `role`, `created_at`, `updated_at`) VALUES (1, 'super', 'super', ',super,super,', '32ad6548eda1551c862df680f32f1aa7', 'kCegTuVUFV', 1, 100, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- 初始化超管配置
INSERT INTO `balala`.`user_config` (`user_id`, `key`, `value`, `created_at`, `updated_at`) VALUES (1, 'notify_switch_global', '1', 1742364763, 1742364763);
INSERT INTO `balala`.`user_config` (`user_id`, `key`, `value`, `created_at`, `updated_at`) VALUES ( 1, 'notify_switch_third_platform', '1', 1742364763, 1742364763);
INSERT INTO `balala`.`user_config` (`user_id`, `key`, `value`, `created_at`, `updated_at`) VALUES ( 1, 'notify_switch_space', '1', 1742364763, 1742364763);


-- 初始化配置
INSERT INTO config (id, config_name, config_key, config_value, config_status, remark, created_at, updated_at, deleted_at) VALUES(1, '空间附件资源', 'space.file.domain', '/upload', 1, NULL, NULL, NULL, NULL);
INSERT INTO config (id, config_name, config_key, config_value, config_status, remark, created_at, updated_at, deleted_at) VALUES(2, '用户头像资源', 'user.avatar.domain', '/upload', 1, NULL, NULL, NULL, NULL);
INSERT INTO config (id, config_name, config_key, config_value, config_status, remark, created_at, updated_at, deleted_at) VALUES(3, '消息通知跳转链接地址', 'notify.redirect.domain', '', 1, NULL, NULL, NULL, NULL);
INSERT INTO config (id, config_name, config_key, config_value, config_status, remark, created_at, updated_at, deleted_at) VALUES(4, 'balala静态资源地址', 'balala.assect.domain', '', 1, NULL, NULL, NULL, NULL);
INSERT INTO config (id, config_name, config_key, config_value, config_status, remark, created_at, updated_at, deleted_at) VALUES(5, 'balala logo', 'balala.logo', '', 1, NULL, NULL, NULL, NULL);
INSERT INTO config (id, config_name, config_key, config_value, config_status, remark, created_at, updated_at, deleted_at) VALUES(6, 'balala 标题', 'balala.title', '', 1, NULL, NULL, NULL, NULL);
INSERT INTO config (id, config_name, config_key, config_value, config_status, remark, created_at, updated_at, deleted_at) VALUES(7, 'balala 注册入口', 'balala.register.entry', '0', 1, NULL, NULL, NULL, NULL);
INSERT INTO config (id, config_name, config_key, config_value, config_status, remark, created_at, updated_at, deleted_at) VALUES(8, 'balala 登录页背景', 'balala.bg', '', 1, NULL, NULL, NULL, NULL);
INSERT INTO config (id, config_name, config_key, config_value, config_status, remark, created_at, updated_at, deleted_at) VALUES(9, 'balala 附件大小', 'balala.attach', '{"value":"100","unit":"MB"}', 1, NULL, NULL, NULL, NULL);
INSERT INTO config (id, config_name, config_key, config_value, config_status, remark, created_at, updated_at, deleted_at) VALUES(10, 'balala 版本号', 'balala.version', '1.4.14.1', 1, NULL, NULL, NULL, NULL);

-- 初始化业务ID
INSERT INTO `business_id` (`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`, `table_name`) VALUES (1, 'work_flow', 0, 30, NULL, '2024-05-09 08:07:12', '2024-07-04 11:48:43', 'work_flow');
INSERT INTO `business_id` (`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`, `table_name`) VALUES (2, 'work_flow_template', 0, 30, NULL, '2024-05-09 08:07:46', '2024-07-04 11:48:43', 'work_flow_template');
INSERT INTO `business_id` (`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`, `table_name`) VALUES (3, 'work_item_status', 0, 30, NULL, '2024-05-09 09:16:15', '2024-07-04 11:48:43', 'work_item_status');
INSERT INTO `business_id` (`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`, `table_name`) VALUES (4, 'work_item_role', 0, 30, NULL, '2024-05-09 10:16:10', '2024-07-04 11:48:43', 'work_item_role');
INSERT INTO `business_id` (`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`, `table_name`) VALUES (5, 'space', 0, 30, NULL, '2024-05-10 12:32:24', '2024-07-04 11:48:43', 'space');
INSERT INTO `business_id` (`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`, `table_name`) VALUES (6, 'user', 0, 30, NULL, '2024-05-10 12:32:56', '2024-07-03 09:22:41', 'user');
INSERT INTO `business_id` (`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`, `table_name`) VALUES (7, 'space_member', 0, 30, NULL, '2024-05-11 02:01:35', '2024-07-04 11:48:43', 'space_member');
INSERT INTO `business_id` (`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`, `table_name`) VALUES (8, 'work_item_type', 0, 30, NULL, '2024-05-22 08:28:47', '2024-07-04 11:48:43', 'work_item_type');
INSERT INTO `business_id` (`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`, `table_name`) VALUES (9, 'space_work_object', 0, 30, NULL, '2024-05-29 08:29:03', '2024-07-04 11:48:43', 'space_work_object');
INSERT INTO `business_id` (`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`, `table_name`) VALUES (10, 'space_work_item', 0, 30, NULL, '2024-05-29 09:01:53', '2024-07-04 11:48:43', 'space_work_item_v2');
INSERT INTO `business_id` (`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`, `table_name`) VALUES (11, 'space_tag', 0, 30, NULL, '2024-06-07 10:50:51', '2024-07-04 11:48:43', 'space_tag');
INSERT INTO `business_id` (`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`, `table_name`) VALUES (12, 'file_info', 0, 30, NULL, '2024-06-12 03:43:22', '2024-07-04 06:59:51', 'file_info');
INSERT INTO `business_id` (`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`, `table_name`) VALUES (13, 'space_file_info', 0, 30, NULL, '2024-06-12 03:43:39', '2024-07-04 06:59:51', 'space_file_info');
INSERT INTO `business_id` (`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`, `table_name`) VALUES (14, 'space_work_version', 0, 30, NULL, '2024-06-27 02:42:51', '2024-07-04 11:48:43', 'space_work_version');