UPDATE space_user_view suv
    JOIN space_global_view sgv
ON suv.space_id = sgv.space_id AND suv.key = sgv.key
    SET suv.outer_id = sgv.id, suv.type = 2
WHERE suv.key IN ('processing', 'expired');

-- 更新关闭终止任务的当前负责人
UPDATE space_work_item_v2 swi
    JOIN (
    SELECT
    work_item_id,
    JSON_ARRAYAGG(director) AS directors
    FROM (
    SELECT DISTINCT
    work_item_id,
    jt.director
    FROM space_work_item_flow_v2 swif
    JOIN JSON_TABLE(
    swif.directors,
    '$[*]' COLUMNS (director VARCHAR(255) PATH '$')
    ) AS jt
    WHERE swif.directors IS NOT NULL
    ) extracted_directors
    GROUP BY work_item_id
    ) flow_data ON swi.id = flow_data.work_item_id
    SET swi.doc = JSON_SET(
        swi.doc,
        '$.directors',
        flow_data.directors
        )
WHERE swi.work_item_status_key IN ('close', 'terminated');