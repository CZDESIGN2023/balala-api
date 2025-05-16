UPDATE space_work_item_v2
SET doc = JSON_SET(doc, '$.node_directors', doc->'$.directors')
WHERE
    pid != 0
	AND JSON_TYPE(doc->'$.node_directors') = 'NULL';

UPDATE work_flow f
    JOIN (
    SELECT
    space_id,
    ranking,
    ROW_NUMBER() over ( PARTITION BY space_id ORDER BY ranking DESC ) AS rn
    FROM
    work_flow
    WHERE
    `status` = 0
    AND ranking < 10000000
    ) t ON f.space_id = t.space_id
    SET f.ranking = t.ranking + 1
WHERE
    f.`status` = 0
  AND f.`key` = 'issue'
  AND f.ranking > 10000000
  AND t.rn = 1;