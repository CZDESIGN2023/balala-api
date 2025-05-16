ALTER TABLE `ods_witem_d`
    ADD COLUMN `count_at` bigint NULL AFTER `child_num`;

ALTER TABLE `space_work_item_v2`
    ADD COLUMN `count_at` bigint NULL COMMENT '计时时间' AFTER `reason`;

