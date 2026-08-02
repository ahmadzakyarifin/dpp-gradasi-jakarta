-- +goose Up
-- +goose StatementBegin

CREATE TABLE billing_types (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(150) NOT NULL,
    description TEXT NULL,
    is_recurring TINYINT(1) NOT NULL DEFAULT 0,
    is_active TINYINT(1) NOT NULL DEFAULT 1,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL,
    CONSTRAINT uq_billing_types_code UNIQUE (code),
    CONSTRAINT uq_billing_types_name UNIQUE (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_billing_types_is_active ON billing_types(is_active);
CREATE INDEX idx_billing_types_deleted_at ON billing_types(deleted_at);

CREATE TABLE billing_rules (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    billing_type_id BIGINT UNSIGNED NOT NULL,
    academic_year_id BIGINT UNSIGNED NOT NULL,
    semester_id BIGINT UNSIGNED NULL,
    active_class_id BIGINT UNSIGNED NULL,
    cohort_id BIGINT UNSIGNED NULL,
    major_id BIGINT UNSIGNED NULL,
    name VARCHAR(150) NOT NULL,
    amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    frequency ENUM('once', 'monthly', 'semester', 'yearly') NOT NULL DEFAULT 'once',
    due_day TINYINT UNSIGNED NULL,
    allow_installment TINYINT(1) NOT NULL DEFAULT 0,
    installment_count SMALLINT UNSIGNED NULL,
    generate_notifications TINYINT(1) NOT NULL DEFAULT 1,
    is_active TINYINT(1) NOT NULL DEFAULT 1,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL,
    CONSTRAINT fk_billing_rules_billing_type FOREIGN KEY (billing_type_id) REFERENCES billing_types(id) ON DELETE RESTRICT,
    CONSTRAINT fk_billing_rules_academic_year FOREIGN KEY (academic_year_id) REFERENCES academic_years(id) ON DELETE RESTRICT,
    CONSTRAINT fk_billing_rules_semester FOREIGN KEY (semester_id) REFERENCES semesters(id) ON DELETE SET NULL,
    CONSTRAINT fk_billing_rules_active_class FOREIGN KEY (active_class_id) REFERENCES active_classes(id) ON DELETE SET NULL,
    CONSTRAINT fk_billing_rules_cohort FOREIGN KEY (cohort_id) REFERENCES cohorts(id) ON DELETE SET NULL,
    CONSTRAINT fk_billing_rules_major FOREIGN KEY (major_id) REFERENCES majors(id) ON DELETE SET NULL,
    CONSTRAINT chk_billing_rules_amount CHECK (amount >= 0),
    CONSTRAINT chk_billing_rules_installment CHECK (installment_count IS NULL OR installment_count > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_billing_rules_billing_type_id ON billing_rules(billing_type_id);
CREATE INDEX idx_billing_rules_academic_year_id ON billing_rules(academic_year_id);
CREATE INDEX idx_billing_rules_semester_id ON billing_rules(semester_id);
CREATE INDEX idx_billing_rules_active_class_id ON billing_rules(active_class_id);
CREATE INDEX idx_billing_rules_cohort_id ON billing_rules(cohort_id);
CREATE INDEX idx_billing_rules_major_id ON billing_rules(major_id);
CREATE INDEX idx_billing_rules_deleted_at ON billing_rules(deleted_at);

CREATE TABLE invoice_generation_batches (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    academic_year_id BIGINT UNSIGNED NOT NULL,
    semester_id BIGINT UNSIGNED NULL,
    billing_rule_id BIGINT UNSIGNED NULL,
    generated_by BIGINT UNSIGNED NULL,
    title VARCHAR(200) NOT NULL,
    status ENUM('processing', 'completed', 'failed', 'cancelled') NOT NULL DEFAULT 'processing',
    total_students INT UNSIGNED NOT NULL DEFAULT 0,
    total_invoices INT UNSIGNED NOT NULL DEFAULT 0,
    total_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    notes TEXT NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    completed_at DATETIME(3) NULL,
    CONSTRAINT fk_invoice_generation_batches_academic_year FOREIGN KEY (academic_year_id) REFERENCES academic_years(id) ON DELETE RESTRICT,
    CONSTRAINT fk_invoice_generation_batches_semester FOREIGN KEY (semester_id) REFERENCES semesters(id) ON DELETE SET NULL,
    CONSTRAINT fk_invoice_generation_batches_billing_rule FOREIGN KEY (billing_rule_id) REFERENCES billing_rules(id) ON DELETE SET NULL,
    CONSTRAINT fk_invoice_generation_batches_generated_by FOREIGN KEY (generated_by) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE invoices (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    invoice_number VARCHAR(100) NOT NULL,
    student_id BIGINT UNSIGNED NOT NULL,
    academic_year_id BIGINT UNSIGNED NOT NULL,
    semester_id BIGINT UNSIGNED NULL,
    billing_rule_id BIGINT UNSIGNED NULL,
    generation_batch_id BIGINT UNSIGNED NULL,
    title VARCHAR(200) NOT NULL,
    issue_date DATE NOT NULL DEFAULT (CURRENT_DATE),
    due_date DATE NULL,
    subtotal_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    discount_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    total_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    paid_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    remaining_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    status ENUM('draft', 'unpaid', 'partial', 'paid', 'overdue', 'cancelled', 'voided') NOT NULL DEFAULT 'unpaid',
    notes TEXT NULL,
    created_by BIGINT UNSIGNED NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL,
    CONSTRAINT uq_invoices_invoice_number UNIQUE (invoice_number),
    CONSTRAINT fk_invoices_student FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE RESTRICT,
    CONSTRAINT fk_invoices_academic_year FOREIGN KEY (academic_year_id) REFERENCES academic_years(id) ON DELETE RESTRICT,
    CONSTRAINT fk_invoices_semester FOREIGN KEY (semester_id) REFERENCES semesters(id) ON DELETE SET NULL,
    CONSTRAINT fk_invoices_billing_rule FOREIGN KEY (billing_rule_id) REFERENCES billing_rules(id) ON DELETE SET NULL,
    CONSTRAINT fk_invoices_generation_batch FOREIGN KEY (generation_batch_id) REFERENCES invoice_generation_batches(id) ON DELETE SET NULL,
    CONSTRAINT fk_invoices_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT chk_invoices_amounts CHECK (subtotal_amount >= 0 AND discount_amount >= 0 AND total_amount >= 0 AND paid_amount >= 0 AND remaining_amount >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_invoices_student_id ON invoices(student_id);
CREATE INDEX idx_invoices_academic_year_id ON invoices(academic_year_id);
CREATE INDEX idx_invoices_semester_id ON invoices(semester_id);
CREATE INDEX idx_invoices_billing_rule_id ON invoices(billing_rule_id);
CREATE INDEX idx_invoices_status ON invoices(status);
CREATE INDEX idx_invoices_due_date ON invoices(due_date);
CREATE INDEX idx_invoices_deleted_at ON invoices(deleted_at);

CREATE TABLE invoice_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    invoice_id BIGINT UNSIGNED NOT NULL,
    billing_type_id BIGINT UNSIGNED NULL,
    name VARCHAR(200) NOT NULL,
    item_type ENUM('charge', 'discount', 'waiver', 'scholarship', 'adjustment') NOT NULL DEFAULT 'charge',
    description TEXT NULL,
    quantity INT UNSIGNED NOT NULL DEFAULT 1,
    unit_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    CONSTRAINT fk_invoice_items_invoice FOREIGN KEY (invoice_id) REFERENCES invoices(id) ON DELETE CASCADE,
    CONSTRAINT fk_invoice_items_billing_type FOREIGN KEY (billing_type_id) REFERENCES billing_types(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_invoice_items_invoice_id ON invoice_items(invoice_id);

CREATE TABLE invoice_installments (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    invoice_id BIGINT UNSIGNED NOT NULL,
    installment_number SMALLINT UNSIGNED NOT NULL,
    due_date DATE NOT NULL,
    amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    paid_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    status ENUM('unpaid', 'partial', 'paid', 'overdue', 'cancelled') NOT NULL DEFAULT 'unpaid',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    CONSTRAINT fk_invoice_installments_invoice FOREIGN KEY (invoice_id) REFERENCES invoices(id) ON DELETE CASCADE,
    CONSTRAINT uq_invoice_installments_number UNIQUE (invoice_id, installment_number),
    CONSTRAINT chk_invoice_installments_amounts CHECK (amount >= 0 AND paid_amount >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_invoice_installments_invoice_id ON invoice_installments(invoice_id);
CREATE INDEX idx_invoice_installments_status ON invoice_installments(status);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS invoice_installments;
DROP TABLE IF EXISTS invoice_items;
DROP TABLE IF EXISTS invoices;
DROP TABLE IF EXISTS invoice_generation_batches;
DROP TABLE IF EXISTS billing_rules;
DROP TABLE IF EXISTS billing_types;
-- +goose StatementEnd
