CREATE TABLE `dws_space_witem_1h` (
  `space_id` bigint DEFAULT NULL COMMENT '纬度-空间',
  `expire_num` int DEFAULT NULL COMMENT '过期数量',
  `num` int DEFAULT NULL COMMENT '任务数量',
  `todo_num` int DEFAULT NULL COMMENT '待办的任务数量',
  `complete_num` int DEFAULT NULL COMMENT '完成的任务数量',
  `close_num` int DEFAULT NULL COMMENT '关闭的任务数量',
  `abort_num` int DEFAULT NULL COMMENT '终止的任务数量',
  `start_date` datetime DEFAULT NULL COMMENT '汇总的开始时间',
  `end_date` datetime DEFAULT NULL COMMENT '汇总的结束时间',
  `_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`_id`),
  UNIQUE KEY `space_id` (`space_id`,`start_date`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='项目任务统计表/小时';

INSERT INTO dws_space_witem_1h (
    space_id,
    start_date,
    end_date,
    num,
    todo_num,
    complete_num,
    close_num,
    abort_num,
    expire_num
) SELECT
      space_id,
      start_date,
      MAX( end_date ) AS end_date,
      SUM( num ) AS num,
      SUM( todo_num ) AS todo_num,
      SUM( complete_num ) AS complete_num,
      SUM( close_num ) AS close_num,
      SUM( abort_num ) AS abort_num,
      SUM( expire_num ) expire_num
FROM
    dws_vers_witem_1h
GROUP BY
    space_id,
    start_date;


-- 清理数据
delete from dws_space_witem_1h where space_id not in (select id from space);
delete from dws_vers_witem_1h where space_id not in (select id from space);
delete from dws_mbr_witem_1h where space_id not in (select id from space);

-- 兼容旧数据
update space_work_item_v2 set flow_mode = "work_flow" where pid != 0;

UPDATE dwd_member
SET end_date = "2025-01-01 00:00:00"
WHERE
    ( space_id, user_id ) NOT IN ( SELECT space_id, user_id FROM space_member )
  AND end_date = "9999-12-31 00:00:00";