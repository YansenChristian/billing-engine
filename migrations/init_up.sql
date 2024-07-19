CREATE TABLE `reminder_tab` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `entity_id` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `entity_type` tinyint(3) unsigned NOT NULL,
    `region` char(2) COLLATE utf8mb4_unicode_ci NOT NULL,
    `content` text COLLATE utf8mb4_unicode_ci NOT NULL,
    `notification_schedule` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `reminder_type` tinyint(3) unsigned NOT NULL,
    `notification_cron_expr` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL,
    `notification_interval` bigint(20) unsigned NOT NULL,
    `notification_start_at` bigint(20) unsigned NOT NULL,
    `notification_end_at` bigint(20) unsigned NOT NULL,
    `next_notify_at` bigint(20) unsigned NOT NULL,
    `created_at` bigint(20) unsigned NOT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_entityid` (`entity_id`),
    KEY `idx_nextnotifyat` (`next_notify_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 DEFAULT COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `subscriber_tab` (
    `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `entity_id` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `entity_type` tinyint UNSIGNED NOT NULL DEFAULT 0,
    `region` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL,
    `created_at` bigint(20) UNSIGNED NOT NULL,
    `deleted_at` bigint(20) UNSIGNED NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_idx_entityid_entitytype` (`entity_id`, `entity_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 DEFAULT COLLATE=utf8mb4_unicode_ci;