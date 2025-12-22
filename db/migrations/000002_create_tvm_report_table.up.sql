CREATE TABLE `tvm_reports` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `tvm_code` varchar(50) NOT NULL,
  `location` varchar(255) NOT NULL,
  `issue_type` varchar(100) NOT NULL,
  `description` text NOT NULL,
  `status` varchar(20) NOT NULL DEFAULT 'pending',
  `priority` varchar(20) DEFAULT 'normal',
  `image_url` varchar(500) DEFAULT NULL,
  `reported_by` bigint unsigned NOT NULL,
  `resolved_by` bigint unsigned DEFAULT NULL,
  `resolved_at` datetime(3) DEFAULT NULL,
  `notes` text,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;