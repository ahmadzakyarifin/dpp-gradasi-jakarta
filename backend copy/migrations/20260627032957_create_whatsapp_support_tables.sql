-- +goose Up
-- +goose StatementBegin

CREATE TABLE whatsapp_configurations (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    session_name VARCHAR(100) NOT NULL,
    display_name VARCHAR(150) NULL,
    phone_number VARCHAR(30) NULL,
    status ENUM('connected', 'disconnected', 'connecting', 'error') NOT NULL DEFAULT 'disconnected',
    is_active TINYINT(1) NOT NULL DEFAULT 1,
    last_connected_at DATETIME(3) NULL,
    last_disconnected_at DATETIME(3) NULL,
    last_error TEXT NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    CONSTRAINT uq_whatsapp_configurations_session UNIQUE (session_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE whatsapp_webhooks (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    configuration_id BIGINT UNSIGNED NULL,
    session_name VARCHAR(100) NULL,
    event_type VARCHAR(100) NOT NULL,
    message_id VARCHAR(200) NULL,
    from_number VARCHAR(30) NULL,
    to_number VARCHAR(30) NULL,
    message_text TEXT NULL,
    raw_payload JSON NOT NULL,
    processing_status ENUM('pending', 'processed', 'failed', 'ignored') NOT NULL DEFAULT 'pending',
    processed_at DATETIME(3) NULL,
    error_message TEXT NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    CONSTRAINT fk_whatsapp_webhooks_configuration FOREIGN KEY (configuration_id) REFERENCES whatsapp_configurations(id) ON DELETE SET NULL,
    CONSTRAINT chk_whatsapp_webhooks_raw_payload CHECK (JSON_VALID(raw_payload))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_whatsapp_webhooks_configuration_id ON whatsapp_webhooks(configuration_id);
CREATE INDEX idx_whatsapp_webhooks_event_type ON whatsapp_webhooks(event_type);
CREATE INDEX idx_whatsapp_webhooks_processing_status ON whatsapp_webhooks(processing_status);

CREATE TABLE customer_support_tickets (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    ticket_number VARCHAR(100) NOT NULL,
    student_id BIGINT UNSIGNED NULL,
    guardian_id BIGINT UNSIGNED NULL,
    created_by BIGINT UNSIGNED NULL,
    assigned_to BIGINT UNSIGNED NULL,
    subject VARCHAR(200) NOT NULL,
    status ENUM('open', 'in_progress', 'resolved', 'closed') NOT NULL DEFAULT 'open',
    priority ENUM('low', 'normal', 'high', 'urgent') NOT NULL DEFAULT 'normal',
    source ENUM('web', 'whatsapp', 'system') NOT NULL DEFAULT 'web',
    closed_at DATETIME(3) NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    CONSTRAINT uq_customer_support_tickets_ticket UNIQUE (ticket_number),
    CONSTRAINT fk_customer_support_tickets_student FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE SET NULL,
    CONSTRAINT fk_customer_support_tickets_guardian FOREIGN KEY (guardian_id) REFERENCES guardians(id) ON DELETE SET NULL,
    CONSTRAINT fk_customer_support_tickets_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT fk_customer_support_tickets_assigned_to FOREIGN KEY (assigned_to) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_customer_support_tickets_student_id ON customer_support_tickets(student_id);
CREATE INDEX idx_customer_support_tickets_guardian_id ON customer_support_tickets(guardian_id);
CREATE INDEX idx_customer_support_tickets_status ON customer_support_tickets(status);
CREATE INDEX idx_customer_support_tickets_assigned_to ON customer_support_tickets(assigned_to);

CREATE TABLE customer_support_messages (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    ticket_id BIGINT UNSIGNED NOT NULL,
    sender_user_id BIGINT UNSIGNED NULL,
    sender_type ENUM('admin', 'guardian', 'student', 'system') NOT NULL,
    message_type ENUM('text', 'image', 'file', 'system') NOT NULL DEFAULT 'text',
    message TEXT NOT NULL,
    attachment_id BIGINT UNSIGNED NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    CONSTRAINT fk_customer_support_messages_ticket FOREIGN KEY (ticket_id) REFERENCES customer_support_tickets(id) ON DELETE CASCADE,
    CONSTRAINT fk_customer_support_messages_sender_user FOREIGN KEY (sender_user_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_customer_support_messages_ticket_id ON customer_support_messages(ticket_id);
CREATE INDEX idx_customer_support_messages_sender_user_id ON customer_support_messages(sender_user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS customer_support_messages;
DROP TABLE IF EXISTS customer_support_tickets;
DROP TABLE IF EXISTS whatsapp_webhooks;
DROP TABLE IF EXISTS whatsapp_configurations;
-- +goose StatementEnd
