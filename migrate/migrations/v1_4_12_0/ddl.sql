-- Step 0: 删除 ods_witem_d 中重复的记录，保留最新的记录
WITH RankedRecords AS (
    SELECT
        _id,
        ROW_NUMBER() OVER (PARTITION BY id, _op_ts ORDER BY _id DESC) AS rn
    FROM
        ods_witem_d
)
DELETE FROM ods_witem_d
WHERE _id IN (
    SELECT _id
    FROM RankedRecords
    WHERE rn > 1
);

-- Step 1: 更新 ods_witem_d 中已存在的记录
WITH LatestRecords AS (
    SELECT
        _id,
        id,
        ROW_NUMBER() OVER (PARTITION BY id ORDER BY _id DESC) AS rn
    FROM
        ods_witem_d
)
UPDATE
    ods_witem_d AS ods
    JOIN
    LatestRecords AS lr
ON
    ods._id = lr._id AND lr.rn = 1
    JOIN
    space_work_item_v2 AS swi
    ON
    ods.id = swi.id
    SET
        ods.pid = swi.pid,
        ods.space_id = swi.space_id,
        ods.work_object_id = swi.work_object_id,
        ods.user_id = swi.user_id,
        ods.flow_id = swi.flow_id,
        ods.flow_key = swi.flow_key,
        ods.work_item_flow_id = swi.work_item_flow_id,
        ods.work_item_guid = swi.work_item_guid,
        ods.work_item_name = swi.work_item_name,
        ods.work_item_status = swi.work_item_status,
        ods.work_item_status_key = swi.work_item_status_key,
        ods.work_item_status_id = swi.work_item_status_id,
        ods.work_item_type_id = swi.work_item_type_id,
        ods.work_item_type_key = swi.work_item_type_key,
        ods.doc = swi.doc,
        ods.flow_template_id = swi.flow_template_id,
        ods.flow_template_version = swi.flow_template_version,
        ods.work_item_type = swi.work_item_type,
        ods.flow_mode = swi.flow_mode,
        ods.work_item_flow_key = swi.work_item_flow_key,
        ods.flow_mode_version = swi.flow_mode_version,
        ods.flow_mode_code = swi.flow_mode_code,
        ods.child_num = swi.child_num,
        ods.last_status_at = swi.last_status_at,
        ods.last_status_key = swi.last_status_key,
        ods.last_status_id = swi.last_status_id,
        ods.last_status = swi.last_status,
        ods.is_restart = swi.is_restart,
        ods.restart_at = swi.restart_at,
        ods.icon_flags = swi.icon_flags,
        ods.restart_user_id = swi.restart_user_id,
        ods.comment_num = swi.comment_num,
        ods.resume_at = swi.resume_at,
        ods.version_id = swi.version_id,
        ods.reason = swi.reason,
        ods.count_at = swi.count_at,
        ods.created_at = swi.created_at,
        ods.updated_at = swi.updated_at,
        ods.deleted_at = swi.deleted_at;

-- Step 2: 插入 ods_witem_d 中不存在的新记录
INSERT INTO ods_witem_d (
    id, pid, space_id, work_object_id, user_id, flow_id, flow_key, work_item_flow_id,
    work_item_guid, work_item_name, work_item_status, work_item_status_key,
    work_item_status_id, work_item_type_id, work_item_type_key, doc, flow_template_id,
    flow_template_version, work_item_type, flow_mode, work_item_flow_key,
    flow_mode_version, flow_mode_code, child_num, last_status_at, last_status_key,
    last_status_id, last_status, is_restart, restart_at, icon_flags, restart_user_id,
    comment_num, resume_at, version_id, reason, count_at, created_at, updated_at, deleted_at,
    _op_ts
)
WITH ExistingRecords AS (
    SELECT DISTINCT id
    FROM ods_witem_d
)
SELECT
    swi.id, swi.pid, swi.space_id, swi.work_object_id, swi.user_id, swi.flow_id, swi.flow_key, swi.work_item_flow_id,
    swi.work_item_guid, swi.work_item_name, swi.work_item_status, swi.work_item_status_key,
    swi.work_item_status_id, swi.work_item_type_id, swi.work_item_type_key, swi.doc, swi.flow_template_id,
    swi.flow_template_version, swi.work_item_type, swi.flow_mode, swi.work_item_flow_key,
    swi.flow_mode_version, swi.flow_mode_code, swi.child_num, swi.last_status_at, swi.last_status_key,
    swi.last_status_id, swi.last_status, swi.is_restart, swi.restart_at, swi.icon_flags, swi.restart_user_id,
    swi.comment_num, swi.resume_at, swi.version_id, swi.reason, swi.count_at, swi.created_at, swi.updated_at, swi.deleted_at,
    swi.created_at
FROM
    space_work_item_v2 AS swi
        LEFT JOIN
    ExistingRecords AS er
    ON
        swi.id = er.id
WHERE
    er.id IS NULL;


-- Step 0: 删除 ods_witem_status_d 中重复的记录，保留最新的记录
WITH RankedRecords AS (
    SELECT
        _id,
        ROW_NUMBER() OVER (PARTITION BY id, _op_ts ORDER BY _id DESC) AS rn
    FROM
        ods_witem_status_d
)
DELETE FROM ods_witem_status_d
WHERE _id IN (
    SELECT _id
    FROM RankedRecords
    WHERE rn > 1
);

-- Step 1: 更新 ods_witem_status_d 中已存在的记录
WITH LatestRecords AS (
    SELECT
        id,
        _id,
        ROW_NUMBER() OVER (PARTITION BY id ORDER BY _id DESC) AS rn
    FROM
        ods_witem_status_d
)
UPDATE
    ods_witem_status_d AS ods
    JOIN
    LatestRecords AS lr
ON
    ods._id = lr._id AND lr.rn = 1
    JOIN
    work_item_status AS wis
    ON
    ods.id = wis.id
    SET
        ods.uuid = wis.uuid,
        ods.space_id = wis.space_id,
        ods.user_id = wis.user_id,
        ods.work_item_type_id = wis.work_item_type_id,
        ods.key = wis.key,
        ods.val = wis.val,
        ods.name = wis.name,
        ods.status_type = wis.status_type,
        ods.status = wis.status,
        ods.ranking = wis.ranking,
        ods.created_at = wis.created_at,
        ods.updated_at = wis.updated_at,
        ods.deleted_at = wis.deleted_at,
        ods.is_sys = wis.is_sys,
        ods.flow_scope = wis.flow_scope;

-- Step 2: 插入 ods_witem_status_d 中不存在的新记录
INSERT INTO ods_witem_status_d (
    id, uuid, space_id, user_id, work_item_type_id, `key`, val, name, status_type,
    status, ranking, created_at, updated_at, deleted_at, is_sys, flow_scope, _op_ts
)
SELECT
    wis.id, wis.uuid, wis.space_id, wis.user_id, wis.work_item_type_id, wis.`key`, wis.val, wis.name, wis.status_type,
    wis.status, wis.ranking, wis.created_at, wis.updated_at, wis.deleted_at, wis.is_sys, wis.flow_scope, wis.created_at
FROM
    work_item_status AS wis
        LEFT JOIN
    (
        SELECT
            id,
            _id,
            ROW_NUMBER() OVER (PARTITION BY id ORDER BY _id DESC) AS rn
        FROM
            ods_witem_status_d
    ) AS lr
    ON
        wis.id = lr.id AND lr.rn = 1
WHERE
    lr.id IS NULL;


-- Step 0: 删除 ods_member_d 中重复的记录，保留最新的记录
WITH RankedRecords AS (
    SELECT
        _id,
        ROW_NUMBER() OVER (PARTITION BY id, _op_ts ORDER BY _id DESC) AS rn
    FROM
        ods_member_d
)
DELETE FROM ods_member_d
WHERE _id IN (
    SELECT _id
    FROM RankedRecords
    WHERE rn > 1
);

-- Step 1: 更新 ods_member_d 中已存在的记录
WITH LatestRecords AS (
    SELECT
        id,
        _id,
        ROW_NUMBER() OVER (PARTITION BY id ORDER BY _id DESC) AS rn
    FROM
        ods_member_d
)
UPDATE
    ods_member_d AS ods
    JOIN
    LatestRecords AS lr
ON
    ods._id = lr._id AND lr.rn = 1
    JOIN
    space_member AS mem
    ON
    ods.id = mem.id
    SET
        ods.space_id = mem.space_id,
        ods.user_id = mem.user_id,
        ods.role_id = mem.role_id,
        ods.remark = mem.remark,
        ods.ranking = mem.ranking,
        ods.notify = mem.notify,
        ods.created_at = mem.created_at,
        ods.updated_at = mem.updated_at,
        ods.deleted_at = mem.deleted_at,
        ods.history_role_id = mem.history_role_id;

-- Step 2: 插入 ods_member_d 中不存在的新记录
INSERT INTO ods_member_d (
    id, space_id, user_id, role_id, remark, ranking, notify, created_at, updated_at, deleted_at, history_role_id, _op_ts
)
WITH LatestRecords AS (
    SELECT
        id,
        _id,
        ROW_NUMBER() OVER (PARTITION BY id ORDER BY _id DESC) AS rn
    FROM
        ods_member_d
)
SELECT
    mem.id, mem.space_id, mem.user_id, mem.role_id, mem.remark, mem.ranking, mem.notify, mem.created_at, mem.updated_at, mem.deleted_at, mem.history_role_id, mem.created_at
FROM
    space_member AS mem
        LEFT JOIN
    LatestRecords AS lr
    ON
        mem.id = lr.id AND lr.rn = 1
WHERE
    lr.id IS NULL;

-- 清空工作项表
TRUNCATE TABLE `dwd_witem`;

ALTER TABLE `dwd_witem`
    ADD COLUMN `work_item_type_key` varchar(255) NULL DEFAULT NULL COMMENT '任务类型key' AFTER `version_id`,
    ADD COLUMN `last_status_at` bigint NULL DEFAULT NULL COMMENT '状态变更时间' AFTER `work_item_type_key`,
    ADD COLUMN `node_directors` json NULL COMMENT '节点负责人' AFTER `directors`;

-- 清空工作项状态表
TRUNCATE TABLE `dim_witem_status`;

ALTER TABLE `dim_witem_status`
    ADD COLUMN `flow_scope` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '工作项状态Val' AFTER `status_type`;

-- 清空 dwd 成员表
TRUNCATE TABLE `dwd_member`;

-- 重新构建 dwd_member 和 dim_witem_status
UPDATE `job_variables` SET `variable_value` = 0 WHERE `job_name` = 'ods_to_dim_job:dim:ods_to_dim_witem_status_task:witemStatus';
UPDATE `job_variables` SET `variable_value` = 0 WHERE `job_name` = 'ods_to_dwd_job:dwd:ods_to_dwd_member_task:memberTask';

-- 清理 doc 中的 remark 和 describe
UPDATE ods_witem_d
SET doc = JSON_SET(doc, '$.remark', '', '$.describe', '')
WHERE JSON_EXTRACT(doc, '$.remark') IS NOT NULL
   OR JSON_EXTRACT(doc, '$.describe') IS NOT NULL;


-- 初始化 dwd_witem
INSERT INTO dwd_witem (
    work_item_id,
    space_id,
    user_id,
    status_id,
    object_id,
    version_id,
    work_item_type_key,
    last_status_at,
    priority,
    plan_start_at,
    plan_complete_at,
    directors,
    node_directors,
    participators,
    gmt_create,
    gmt_modified,
    start_date,
    end_date
)
SELECT
    swi.id AS work_item_id, -- 工作项id
    swi.space_id AS space_id, -- 空间id
    swi.user_id AS user_id, -- 创建人id
    swi.work_item_status_id AS status_id, -- 工作项状态id
    swi.work_object_id AS object_id, -- 模块id
    swi.version_id AS version_id, -- 版本id
    swi.work_item_type_key AS work_item_type_key, -- 任务类型key
    swi.last_status_at AS last_status_at, -- 状态变更时间
    JSON_UNQUOTE(JSON_EXTRACT(swi.doc, '$.priority')) AS priority, -- 优先级
    JSON_UNQUOTE(JSON_EXTRACT(swi.doc, '$.plan_start_at')) AS plan_start_at, -- 排期-开始时间
    JSON_UNQUOTE(JSON_EXTRACT(swi.doc, '$.plan_complete_at')) AS plan_complete_at, -- 排期-结束时间
    JSON_EXTRACT(swi.doc, '$.directors') AS directors, -- 负责人
    JSON_EXTRACT(swi.doc, '$.node_directors') AS node_directors, -- 节点负责人
    JSON_EXTRACT(swi.doc, '$.participators') AS participators, -- 参与人
    FROM_UNIXTIME(swi.created_at) AS gmt_create, -- 工作项创建时间
    FROM_UNIXTIME(swi.updated_at) AS gmt_modified, -- 工作项更新时间
    FROM_UNIXTIME(swi.created_at) AS start_date, -- 生效日期
    '9999-12-31 00:00:00' AS end_date -- 失效日期
FROM
    space_work_item_v2 swi;
