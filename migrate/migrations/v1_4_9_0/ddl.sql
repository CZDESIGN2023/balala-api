CREATE TABLE `space_user_view`
(
    `id`           bigint NOT NULL AUTO_INCREMENT,
    `user_id`      bigint                                                         DEFAULT NULL,
    `space_id`     bigint                                                         DEFAULT NULL,
    `key`          varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT NULL,
    `name`         varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT NULL,
    `type`         tinyint                                                        DEFAULT NULL COMMENT '1系统2全局3用户',
    `outer_id`     bigint                                                         DEFAULT NULL COMMENT '当type为1|2时指向view_space',
    `query_config` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
    `table_config` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
    `ranking`      int                                                            DEFAULT NULL,
    `status`       tinyint                                                        DEFAULT NULL,
    `created_at`   bigint                                                         DEFAULT NULL,
    `updated_at`   bigint                                                         DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY            `space_id_user_id` (`space_id`,`user_id`) USING BTREE,
    KEY            `user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `space_global_view`
(
    `id`           bigint NOT NULL AUTO_INCREMENT,
    `space_id`     bigint                                                         DEFAULT NULL,
    `key`          varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT NULL,
    `name`         varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT NULL,
    `type`         tinyint                                                        DEFAULT NULL COMMENT '1系统2全局3用户',
    `query_config` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
    `table_config` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
    `created_at`   bigint                                                         DEFAULT NULL,
    `updated_at`   bigint                                                         DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY            `space_id` (`space_id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;