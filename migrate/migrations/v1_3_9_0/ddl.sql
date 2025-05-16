ALTER TABLE `oper_log`
    ADD COLUMN `space_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '空间名' AFTER `show_type`;