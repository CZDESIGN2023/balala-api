ALTER TABLE `file_info`
    MODIFY COLUMN `name` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '参数键名',
    MODIFY COLUMN `upload_path` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '上傳完成檔案路徑';

ALTER TABLE `space_file_info`
    MODIFY COLUMN `file_uri` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '文件uri访问路径',
    MODIFY COLUMN `file_name` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '文件名称, 冗余';