UPDATE space_work_item_v2 t
    LEFT JOIN (
				SELECT
					work_item_id,
					JSON_ARRAYAGG( element ) AS node_directors
					FROM
							(
									SELECT DISTINCT
											work_item_id,
											element
									FROM
											space_work_item_flow_v2
									JOIN JSON_TABLE (
											directors,
											'$[*]' COLUMNS ( element VARCHAR ( 255 ) PATH '$' )
									) AS jt
					) AS t
					GROUP BY
					work_item_id
    ) AS d
ON t.id = d.work_item_id
SET doc = JSON_SET(doc, "$.node_directors", if(d.node_directors IS NULL,JSON_ARRAY(),d.node_directors) );
