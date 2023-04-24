CREATE TABLE `insomnia_project_linting_rule` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `organization_id` INT UNSIGNED NOT NULL,
    `project_id` INT UNSIGNED NOT NULL,
    `creator_uid` INT UNSIGNED NOT NULL,
    `rule` TEXT NOT NULL COMMENT 'custom linting rule string',
    `rule_category` TINYINT UNSIGNED NOT NULL COMMENT '1:json_format, 2:yaml format',
    `create_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_organ_project` (`organization_id`, `project_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='save linting rule created';

CREATE TABLE `insomnia_project_apply_rule` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `organization_id` INT UNSIGNED NOT NULL,
    `project_id` INT UNSIGNED NOT NULL,
    `rule_id` BIGINT UNSIGNED NOT NULL COMMENT 'id from insomnia_project_linting_rule',
    `opt_uid` INT UNSIGNED NOT NULL COMMENT 'apply rule user_id',
    `create_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE `uni_org_project` (`organization_id`, `project_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='bind rule with project';