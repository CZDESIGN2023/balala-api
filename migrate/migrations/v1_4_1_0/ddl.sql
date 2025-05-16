ALTER TABLE `file_info`
    ADD COLUMN `meta` json NULL COMMENT '文件元数据' AFTER `owner`;