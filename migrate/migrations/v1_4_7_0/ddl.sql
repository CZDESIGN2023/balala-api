ALTER TABLE `ods_witem_status_d`
    ADD COLUMN `flow_scope` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '可作用的流程范围  空(表示全部) |state_flow | work_flow' AFTER `is_sys`;

ALTER TABLE `work_item_status`
    ADD COLUMN `flow_scope` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '可作用的流程范围  空(表示全部) |state_flow | work_flow' AFTER `is_sys`;

ALTER TABLE `work_item_role`
    ADD COLUMN `flow_scope` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '可作用的流程范围  空(表示全部) |state_flow | work_flow' AFTER `deleted_at`;


UPDATE `work_item_status` SET `flow_scope` = 'work_flow' WHERE `flow_scope` = '' AND `val` != '3';
UPDATE `work_item_role` SET `flow_scope` = 'work_flow' WHERE `flow_scope` = '';

UPDATE `space_work_item_v2` SET `flow_mode` = 'state_flow' WHERE `pid` != 0;
UPDATE `space_work_item_v2` SET `flow_mode` = 'work_flow' WHERE `pid` = 0 AND `flow_mode` = '';