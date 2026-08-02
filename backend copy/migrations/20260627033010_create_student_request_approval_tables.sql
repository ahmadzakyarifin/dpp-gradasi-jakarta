-- +goose Up
-- +goose StatementBegin

CREATE TABLE student_requests (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    request_number VARCHAR(100) NOT NULL,
    student_id BIGINT UNSIGNED NOT NULL,
    guardian_id BIGINT UNSIGNED NULL,
    invoice_id BIGINT UNSIGNED NULL,
    requested_by BIGINT UNSIGNED NULL,
    request_type ENUM('financial_aid', 'scholarship', 'installment', 'waiver', 'refund', 'discount', 'invoice_adjustment', 'other') NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT NULL,
    requested_amount DECIMAL(14,2) NULL,
    approved_amount DECIMAL(14,2) NULL,
    status ENUM('pending', 'reviewed', 'approved', 'rejected', 'cancelled') NOT NULL DEFAULT 'pending',
    approved_by BIGINT UNSIGNED NULL,
    approved_at DATETIME(3) NULL,
    rejected_reason TEXT NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    CONSTRAINT uq_student_requests_request_number UNIQUE (request_number),
    CONSTRAINT fk_student_requests_student FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE RESTRICT,
    CONSTRAINT fk_student_requests_guardian FOREIGN KEY (guardian_id) REFERENCES guardians(id) ON DELETE SET NULL,
    CONSTRAINT fk_student_requests_invoice FOREIGN KEY (invoice_id) REFERENCES invoices(id) ON DELETE SET NULL,
    CONSTRAINT fk_student_requests_requested_by FOREIGN KEY (requested_by) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT fk_student_requests_approved_by FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT chk_student_requests_amounts CHECK ((requested_amount IS NULL OR requested_amount >= 0) AND (approved_amount IS NULL OR approved_amount >= 0))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_student_requests_student_id ON student_requests(student_id);
CREATE INDEX idx_student_requests_guardian_id ON student_requests(guardian_id);
CREATE INDEX idx_student_requests_invoice_id ON student_requests(invoice_id);
CREATE INDEX idx_student_requests_status ON student_requests(status);
CREATE INDEX idx_student_requests_request_type ON student_requests(request_type);

CREATE TABLE approval_actions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    student_request_id BIGINT UNSIGNED NOT NULL,
    actor_id BIGINT UNSIGNED NULL,
    action ENUM('submitted', 'reviewed', 'approved', 'rejected', 'cancelled', 'revised') NOT NULL,
    note TEXT NULL,
    metadata JSON NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    CONSTRAINT fk_approval_actions_student_request FOREIGN KEY (student_request_id) REFERENCES student_requests(id) ON DELETE CASCADE,
    CONSTRAINT fk_approval_actions_actor FOREIGN KEY (actor_id) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT chk_approval_actions_metadata CHECK (metadata IS NULL OR JSON_VALID(metadata))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_approval_actions_student_request_id ON approval_actions(student_request_id);
CREATE INDEX idx_approval_actions_actor_id ON approval_actions(actor_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS approval_actions;
DROP TABLE IF EXISTS student_requests;
-- +goose StatementEnd
