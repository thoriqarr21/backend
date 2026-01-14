CREATE TABLE tvm_report_histories (
    id BIGSERIAL PRIMARY KEY,
    tvm_report_id BIGINT NOT NULL,
    from_status VARCHAR(30),
    to_status VARCHAR(30) NOT NULL,
    changed_by BIGINT NOT NULL,
    role VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_report
        FOREIGN KEY (tvm_report_id)
        REFERENCES tvm_reports(id)
        ON DELETE CASCADE
);
