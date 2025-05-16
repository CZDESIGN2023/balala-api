-- 用户配置表
CREATE TABLE `user_config`
(
    `id`         bigint NOT NULL AUTO_INCREMENT,
    `user_id`    bigint                                                        DEFAULT NULL,
    `key`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
    `value`      varchar(1023)                                                 DEFAULT NULL,
    `created_at` bigint                                                        DEFAULT NULL,
    `updated_at` bigint                                                        DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY          `user_id_key` (`user_id`,`key`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 插入数据
INSERT INTO user_config (user_id, `key`, `value`, created_at, updated_at)
SELECT u.id                   AS user_id,    -- 用户 ID
       'notify_switch_global' AS `key`,      -- 配置项键
       '1'                    AS `value`,    -- 配置项值（默认开启）
       UNIX_TIMESTAMP()       AS created_at, -- 创建时间戳（毫秒）
       UNIX_TIMESTAMP()       AS updated_at  -- 更新时间戳（毫秒）
FROM user u;

INSERT INTO user_config (user_id, `key`, `value`, created_at, updated_at)
SELECT u.id                           AS user_id,    -- 用户 ID
       'notify_switch_third_platform' AS `key`,      -- 配置项键
       '1'                            AS `value`,    -- 配置项值（默认开启）
       UNIX_TIMESTAMP()               AS created_at, -- 创建时间戳（毫秒）
       UNIX_TIMESTAMP()               AS updated_at  -- 更新时间戳（毫秒）
FROM user u;

INSERT INTO user_config (user_id, `key`, `value`, created_at, updated_at)
SELECT u.id                           AS user_id,    -- 用户 ID
       'notify_switch_space' AS `key`,      -- 配置项键
       '1'                            AS `value`,    -- 配置项值（默认开启）
       UNIX_TIMESTAMP()               AS created_at, -- 创建时间戳（毫秒）
       UNIX_TIMESTAMP()               AS updated_at  -- 更新时间戳（毫秒）
FROM user u;


-- 添加字段
ALTER TABLE `third_pf_account`
    ADD COLUMN `notify` tinyint NULL DEFAULT 1 COMMENT '通知' AFTER `pf_user_account`;