ALTER TABLE `oper_log`
    ADD COLUMN `show_type` int NULL COMMENT '展示类型' AFTER `space_id`;