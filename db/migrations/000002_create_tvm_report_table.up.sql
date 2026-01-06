-- CREATE TABLE `tvm_reports` (
--   `id` bigint unsigned NOT NULL AUTO_INCREMENT,
--   `tvm_code` varchar(50) NOT NULL,
--   `location` varchar(255) NOT NULL,
--   `issue_type` varchar(100) NOT NULL,
--   `description` text NOT NULL,
--   `status` varchar(20) NOT NULL DEFAULT 'pending',
--   `priority` varchar(20) DEFAULT 'normal',
--   `image_url` varchar(500) DEFAULT NULL,
--   `reported_by` bigint unsigned NOT NULL,
--   `resolved_by` bigint unsigned DEFAULT NULL,
--   `resolved_at` datetime(3) DEFAULT NULL,
--   `notes` text,
--   `created_at` datetime(3) DEFAULT NULL,
--   `updated_at` datetime(3) DEFAULT NULL,
--   PRIMARY KEY (`id`)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE tvm_reports (
    id BIGSERIAL PRIMARY KEY,
    barang_id BIGINT NOT NULL,
    tvm_code VARCHAR(50) NOT NULL,
    location VARCHAR(255) NOT NULL,
    issue_type VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    priority VARCHAR(20) DEFAULT 'normal',
    image_url VARCHAR(500) DEFAULT NULL,
    reported_by BIGINT NOT NULL,
    resolved_by BIGINT DEFAULT NULL,
    resolved_at TIMESTAMP(3) DEFAULT NULL,
    -- notes TEXT,
    created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP
);

-- Trigger untuk auto-update updated_at
CREATE OR REPLACE FUNCTION update_tvm_reports_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_tvm_reports_updated_at
BEFORE UPDATE ON tvm_reports
FOR EACH ROW
EXECUTE FUNCTION update_tvm_reports_updated_at();

-- Index untuk performa (opsional tapi disarankan)
CREATE INDEX idx_tvm_reports_tvm_code ON tvm_reports(tvm_code);
CREATE INDEX idx_tvm_reports_status ON tvm_reports(status);
CREATE INDEX idx_tvm_reports_reported_by ON tvm_reports(reported_by);