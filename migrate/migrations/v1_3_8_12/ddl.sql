ALTER TABLE `file_info`
    ADD COLUMN `cover` varchar(1024) NULL COMMENT '封面' AFTER `pwd`;