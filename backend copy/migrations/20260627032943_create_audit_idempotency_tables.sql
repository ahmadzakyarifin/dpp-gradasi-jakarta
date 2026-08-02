-- +goose Up
-- +goose StatementBegin

CREATE TABLE activity_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    actor_id BIGINT UNSIGNED NULL,
    actor_name VARCHAR(150) NULL,
    actor_role VARCHAR(100) NULL,

    action VARCHAR(100) NOT NULL,

    status ENUM(
        'success',
        'failed',
        'warning'
    ) NOT NULL DEFAULT 'success',

    entity_type VARCHAR(100) NOT NULL,
    entity_id BIGINT UNSIGNED NULL,
    entity_label VARCHAR(200) NULL,

    risk_level ENUM(
        'low',
        'medium',
        'high'
    ) NOT NULL DEFAULT 'low',

    description TEXT NULL,

    ip_address VARCHAR(45) NULL,
    user_agent TEXT NULL,

    metadata JSON NULL,

    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),

    CONSTRAINT fk_activity_logs_actor
        FOREIGN KEY (actor_id)
        REFERENCES users(id)
        ON DELETE SET NULL,

    CONSTRAINT chk_activity_logs_metadata_json
        CHECK (
            metadata IS NULL
            OR JSON_VALID(metadata)
        )
) ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_activity_logs_actor_id
ON activity_logs(actor_id);

CREATE INDEX idx_activity_logs_action
ON activity_logs(action);

CREATE INDEX idx_activity_logs_status
ON activity_logs(status);

CREATE INDEX idx_activity_logs_entity
ON activity_logs(entity_type, entity_id);

CREATE INDEX idx_activity_logs_risk_level
ON activity_logs(risk_level);

CREATE INDEX idx_activity_logs_created_at
ON activity_logs(created_at);

CREATE INDEX idx_activity_logs_action_created_at
ON activity_logs(action, created_at);

CREATE INDEX idx_activity_logs_status_created_at
ON activity_logs(status, created_at);

CREATE INDEX idx_activity_logs_risk_created_at
ON activity_logs(risk_level, created_at);

CREATE TABLE idempotency_keys (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `key` VARCHAR(200) NOT NULL,
    user_id BIGINT UNSIGNED NULL,
    request_method VARCHAR(20) NOT NULL,
    request_path VARCHAR(500) NOT NULL,
    request_hash CHAR(64) NOT NULL,
    response_status INT NULL,
    response_body JSON NULL,
    status ENUM('processing', 'completed', 'failed')
        NOT NULL DEFAULT 'processing',
    locked_until DATETIME(3) NULL,
    expires_at DATETIME(3) NOT NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),


    CONSTRAINT uq_idempotency_keys_key UNIQUE (`key`),

    CONSTRAINT fk_idempotency_keys_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_idempotency_keys_user_id
ON idempotency_keys(user_id);

CREATE INDEX idx_idempotency_keys_expires_at
ON idempotency_keys(expires_at);

CREATE INDEX idx_idempotency_keys_status
ON idempotency_keys(status);

CREATE INDEX idx_idempotency_keys_locked_until
ON idempotency_keys(locked_until);

CREATE INDEX idx_idempotency_keys_status_locked_until
ON idempotency_keys(status, locked_until);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin


DROP TABLE IF EXISTS idempotency_keys;
DROP TABLE IF EXISTS activity_logs;

-- +goose StatementEnd