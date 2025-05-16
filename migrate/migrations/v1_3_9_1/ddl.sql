ALTER TABLE `oper_log`
    ADD COLUMN `module_flag` bigint NULL COMMENT '位标记' AFTER `space_name`;