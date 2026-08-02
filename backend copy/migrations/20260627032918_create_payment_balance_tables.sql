-- +goose Up
-- +goose StatementBegin

CREATE TABLE payment_methods (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(150) NOT NULL,
    provider VARCHAR(100) NULL,
    is_active TINYINT(1) NOT NULL DEFAULT 1,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    CONSTRAINT uq_payment_methods_code UNIQUE (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE payments (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    payment_number VARCHAR(100) NOT NULL,
    student_id BIGINT UNSIGNED NOT NULL,
    guardian_id BIGINT UNSIGNED NULL,
    payment_method_id BIGINT UNSIGNED NOT NULL,
    amount DECIMAL(14,2) NOT NULL,
    deposit_applied DECIMAL(14,2) NOT NULL DEFAULT 0,
    payment_date DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    status ENUM('pending', 'success', 'failed', 'cancelled', 'refunded') NOT NULL DEFAULT 'pending',
    reference_number VARCHAR(150) NULL,
    notes TEXT NULL,
    received_by BIGINT UNSIGNED NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL,
    CONSTRAINT uq_payments_payment_number UNIQUE (payment_number),
    CONSTRAINT fk_payments_student FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE RESTRICT,
    CONSTRAINT fk_payments_guardian FOREIGN KEY (guardian_id) REFERENCES guardians(id) ON DELETE SET NULL,
    CONSTRAINT fk_payments_payment_method FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id) ON DELETE RESTRICT,
    CONSTRAINT fk_payments_received_by FOREIGN KEY (received_by) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT chk_payments_amount CHECK (amount >= 0 AND deposit_applied >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_payments_student_id ON payments(student_id);
CREATE INDEX idx_payments_guardian_id ON payments(guardian_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_payment_date ON payments(payment_date);
CREATE INDEX idx_payments_deleted_at ON payments(deleted_at);

CREATE TABLE payment_allocations (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    payment_id BIGINT UNSIGNED NOT NULL,
    invoice_id BIGINT UNSIGNED NOT NULL,
    invoice_installment_id BIGINT UNSIGNED NULL,
    amount DECIMAL(14,2) NOT NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    CONSTRAINT fk_payment_allocations_payment FOREIGN KEY (payment_id) REFERENCES payments(id) ON DELETE CASCADE,
    CONSTRAINT fk_payment_allocations_invoice FOREIGN KEY (invoice_id) REFERENCES invoices(id) ON DELETE RESTRICT,
    CONSTRAINT fk_payment_allocations_installment FOREIGN KEY (invoice_installment_id) REFERENCES invoice_installments(id) ON DELETE SET NULL,
    CONSTRAINT chk_payment_allocations_amount CHECK (amount > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_payment_allocations_payment_id ON payment_allocations(payment_id);
CREATE INDEX idx_payment_allocations_invoice_id ON payment_allocations(invoice_id);

CREATE TABLE payment_gateway_transactions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    payment_id BIGINT UNSIGNED NULL,
    student_id BIGINT UNSIGNED NOT NULL,
    gateway VARCHAR(100) NOT NULL DEFAULT 'midtrans',
    order_id VARCHAR(150) NOT NULL,
    transaction_id VARCHAR(150) NULL,
    payment_type VARCHAR(100) NULL,
    gross_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
    status ENUM('pending', 'challenge', 'success', 'settlement', 'capture', 'deny', 'expire', 'cancel', 'failure', 'refund') NOT NULL DEFAULT 'pending',
    redirect_url TEXT NULL,
    snap_token TEXT NULL,
    raw_response JSON NULL,
    expired_at DATETIME(3) NULL,
    paid_at DATETIME(3) NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    CONSTRAINT uq_payment_gateway_transactions_order_id UNIQUE (order_id),
    CONSTRAINT fk_payment_gateway_transactions_payment FOREIGN KEY (payment_id) REFERENCES payments(id) ON DELETE SET NULL,
    CONSTRAINT fk_payment_gateway_transactions_student FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE RESTRICT,
    CONSTRAINT chk_payment_gateway_transactions_raw_response CHECK (raw_response IS NULL OR JSON_VALID(raw_response))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_payment_gateway_transactions_student_id ON payment_gateway_transactions(student_id);
CREATE INDEX idx_payment_gateway_transactions_payment_id ON payment_gateway_transactions(payment_id);
CREATE INDEX idx_payment_gateway_transactions_status ON payment_gateway_transactions(status);

CREATE TABLE payment_webhooks (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    gateway_transaction_id BIGINT UNSIGNED NULL,
    gateway VARCHAR(100) NOT NULL DEFAULT 'midtrans',
    order_id VARCHAR(150) NULL,
    external_event_id VARCHAR(150) NULL,
    event_type VARCHAR(100) NULL,
    signature_key TEXT NULL,
    raw_payload JSON NOT NULL,
    processing_status ENUM('pending', 'processed', 'failed', 'ignored') NOT NULL DEFAULT 'pending',
    processed_at DATETIME(3) NULL,
    error_message TEXT NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    CONSTRAINT fk_payment_webhooks_gateway_transaction FOREIGN KEY (gateway_transaction_id) REFERENCES payment_gateway_transactions(id) ON DELETE SET NULL,
    CONSTRAINT uq_payment_webhooks_external_event_id UNIQUE (external_event_id),
    CONSTRAINT chk_payment_webhooks_raw_payload CHECK (JSON_VALID(raw_payload))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_payment_webhooks_gateway_transaction_id ON payment_webhooks(gateway_transaction_id);
CREATE INDEX idx_payment_webhooks_order_id ON payment_webhooks(order_id);
CREATE INDEX idx_payment_webhooks_processing_status ON payment_webhooks(processing_status);

CREATE TABLE student_balances (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    student_id BIGINT UNSIGNED NOT NULL,
    balance DECIMAL(14,2) NOT NULL DEFAULT 0,
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    CONSTRAINT uq_student_balances_student UNIQUE (student_id),
    CONSTRAINT fk_student_balances_student FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
    CONSTRAINT chk_student_balances_balance CHECK (balance >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE student_balance_mutations (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    student_id BIGINT UNSIGNED NOT NULL,
    payment_id BIGINT UNSIGNED NULL,
    invoice_id BIGINT UNSIGNED NULL,
    mutation_type VARCHAR(50) NOT NULL,
    direction ENUM('in', 'out') NOT NULL,
    amount DECIMAL(14,2) NOT NULL,
    balance_before DECIMAL(14,2) NOT NULL,
    balance_after DECIMAL(14,2) NOT NULL,
    description TEXT NULL,
    created_by BIGINT UNSIGNED NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    CONSTRAINT fk_student_balance_mutations_student FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE RESTRICT,
    CONSTRAINT fk_student_balance_mutations_payment FOREIGN KEY (payment_id) REFERENCES payments(id) ON DELETE SET NULL,
    CONSTRAINT fk_student_balance_mutations_invoice FOREIGN KEY (invoice_id) REFERENCES invoices(id) ON DELETE SET NULL,
    CONSTRAINT fk_student_balance_mutations_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT chk_student_balance_mutations_amount CHECK (amount > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_student_balance_mutations_student_id ON student_balance_mutations(student_id);
CREATE INDEX idx_student_balance_mutations_payment_id ON student_balance_mutations(payment_id);
CREATE INDEX idx_student_balance_mutations_invoice_id ON student_balance_mutations(invoice_id);
CREATE INDEX idx_student_balance_mutations_created_at ON student_balance_mutations(created_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS student_balance_mutations;
DROP TABLE IF EXISTS student_balances;
DROP TABLE IF EXISTS payment_webhooks;
DROP TABLE IF EXISTS payment_gateway_transactions;
DROP TABLE IF EXISTS payment_allocations;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS payment_methods;
-- +goose StatementEnd
