CREATE TABLE `users_tab` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `created_at` bigint(20) unsigned NOT NULL,
    `updated_at` bigint(20) unsigned NOT NULL,
    `deleted_at` bigint(20) unsigned NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 DEFAULT COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `loan_requests_tab` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `user_id` bigint(20) unsigned NOT NULL,
    `loan_amount` decimal(25, 2) NOT NULL,
    `principal_paid_amount` decimal(25, 2) NOT NULL,
    `interest_paid_amount` decimal(25, 2) NOT NULL,
    `disbursement_time` bigint(20) unsigned NOT NULL,
    `tenure_value` int NOT NULL,
    `tenure_unit` tinyint unsigned NOT NULL,
    `status` tinyint unsigned NOT NULL,
    `annual_interest_rate` decimal(25, 2) NOT NULL,
    `created_at` bigint(20) unsigned NOT NULL,
    `updated_at` bigint(20) unsigned NOT NULL,
    `deleted_at` bigint(20) unsigned NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 DEFAULT COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `billings_tab` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `billing_id` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
    `loan_id` bigint(20) unsigned NOT NULL,
    `payment_id` bigint(20) unsigned NOT NULL,
    `recurring_index` int NOT NULL,
    `principal_amount` decimal(25, 2) NOT NULL,
    `interest_amount` decimal(25, 2) NOT NULL,
    `total_amount` decimal(25, 2) NOT NULL,
    `due_time` bigint(20) unsigned NOT NULL,
    `payment_completed_at` bigint(20) unsigned NOT NULL,
    `status` tinyint unsigned NOT NULL,
    `created_at` bigint(20) unsigned NOT NULL,
    `updated_at` bigint(20) unsigned NOT NULL,
    `deleted_at` bigint(20) unsigned NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uniq_idx_billingid` (`billing_id`),
    UNIQUE INDEX `uniq_idx_loanid_recurringindex` (`loan_id`,`recurring_index`),
    INDEX `idx_loanid_status_duetime` (`loan_id`,`status`,`due_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 DEFAULT COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `payments_tab` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `user_id` bigint(20) unsigned NOT NULL,
    `amount` decimal(25, 2) NOT NULL,
    `created_at` bigint(20) unsigned NOT NULL,
    `updated_at` bigint(20) unsigned NOT NULL,
    `deleted_at` bigint(20) unsigned NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 DEFAULT COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `loan_request_histories_tab` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `loan_id` bigint(20) unsigned NOT NULL,
    `principal_paid_amount` decimal(25, 2) NOT NULL,
    `interest_paid_amount` decimal(25, 2) NOT NULL,
    `status` tinyint unsigned NOT NULL,
    `created_at` bigint(20) unsigned NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 DEFAULT COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `billing_histories_tab` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `billing_id` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
    `payment_completed_at` bigint(20) unsigned NOT NULL,
    `status` tinyint unsigned NOT NULL,
    `created_at` bigint(20) unsigned NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 DEFAULT COLLATE=utf8mb4_unicode_ci;
