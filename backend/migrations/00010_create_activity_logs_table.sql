-- +goose Up
-- +goose StatementBegin
CREATE TABLE activity_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    actor_id INT NULL,
    actor_name VARCHAR(150) NULL,
    actor_role VARCHAR(100) NULL,

    action VARCHAR(100) NOT NULL,

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_activity_logs_actor_id ON activity_logs(actor_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_activity_logs_action ON activity_logs(action);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_activity_logs_entity ON activity_logs(entity_type, entity_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_activity_logs_risk_level ON activity_logs(risk_level);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_activity_logs_created_at ON activity_logs(created_at);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_activity_logs_action_created_at ON activity_logs(action, created_at);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_activity_logs_risk_created_at ON activity_logs(risk_level, created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS activity_logs;
-- +goose StatementEnd
