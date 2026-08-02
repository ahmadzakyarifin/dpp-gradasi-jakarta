-- +goose Up
-- +goose StatementBegin

CREATE TABLE academic_years (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active TINYINT(1) NOT NULL DEFAULT 0,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL,
    CONSTRAINT uq_academic_years_name UNIQUE (name),
    CONSTRAINT chk_academic_years_date CHECK (end_date >= start_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_academic_years_is_active ON academic_years(is_active);
CREATE INDEX idx_academic_years_deleted_at ON academic_years(deleted_at);

CREATE TABLE cohorts (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    description TEXT NULL,
    is_active TINYINT(1) NOT NULL DEFAULT 1,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL,
    CONSTRAINT uq_cohorts_name UNIQUE (name),
    CONSTRAINT chk_cohorts_date CHECK (end_date >= start_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_cohorts_start_date ON cohorts(start_date);
CREATE INDEX idx_cohorts_is_active ON cohorts(is_active);
CREATE INDEX idx_cohorts_deleted_at ON cohorts(deleted_at);

CREATE TABLE semesters (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    academic_year_id BIGINT UNSIGNED NOT NULL,
    name ENUM('ganjil', 'genap') NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active TINYINT(1) NOT NULL DEFAULT 0,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL,
    CONSTRAINT uq_semesters_academic_year_name UNIQUE (academic_year_id, name),
    CONSTRAINT fk_semesters_academic_year FOREIGN KEY (academic_year_id) REFERENCES academic_years(id) ON DELETE CASCADE,
    CONSTRAINT chk_semesters_date CHECK (end_date >= start_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_semesters_academic_year_id ON semesters(academic_year_id);
CREATE INDEX idx_semesters_is_active ON semesters(is_active);

CREATE TABLE majors (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(30) NULL,
    name VARCHAR(150) NOT NULL,
    is_active TINYINT(1) NOT NULL DEFAULT 1,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL,
    CONSTRAINT uq_majors_code UNIQUE (code),
    CONSTRAINT uq_majors_name UNIQUE (name),
    CONSTRAINT chk_majors_code_format CHECK (code REGEXP '^[A-Z0-9]{2,12}$')
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_majors_is_active ON majors(is_active);
CREATE INDEX idx_majors_deleted_at ON majors(deleted_at);

CREATE TABLE class_templates (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    major_id BIGINT UNSIGNED NULL,
    name VARCHAR(100) NOT NULL,
    grade_level TINYINT UNSIGNED NOT NULL,
    description TEXT NULL,
    is_active TINYINT(1) NOT NULL DEFAULT 1,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL,
    CONSTRAINT uq_class_templates_major_grade_name UNIQUE (major_id, grade_level, name),
    CONSTRAINT fk_class_templates_major FOREIGN KEY (major_id) REFERENCES majors(id) ON DELETE SET NULL,
    CONSTRAINT chk_class_templates_grade_level CHECK (grade_level BETWEEN 1 AND 12)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_class_templates_major_id ON class_templates(major_id);
CREATE INDEX idx_class_templates_grade_level ON class_templates(grade_level);
CREATE INDEX idx_class_templates_is_active ON class_templates(is_active);
CREATE INDEX idx_class_templates_deleted_at ON class_templates(deleted_at);

CREATE TABLE active_classes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    academic_year_id BIGINT UNSIGNED NOT NULL,
    class_template_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(120) NOT NULL,
    homeroom_number VARCHAR(30) NULL,
    homeroom_teacher_name VARCHAR(150) NULL,
    room VARCHAR(100) NULL,
    capacity SMALLINT UNSIGNED NOT NULL DEFAULT 0,
    student_count SMALLINT UNSIGNED NOT NULL DEFAULT 0,
    is_active TINYINT(1) NOT NULL DEFAULT 1,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL,
    CONSTRAINT uq_active_classes_year_name UNIQUE (academic_year_id, name),
    CONSTRAINT uq_active_classes_year_template_homeroom UNIQUE (academic_year_id, class_template_id, homeroom_number),
    CONSTRAINT fk_active_classes_academic_year FOREIGN KEY (academic_year_id) REFERENCES academic_years(id) ON DELETE RESTRICT,
    CONSTRAINT fk_active_classes_template FOREIGN KEY (class_template_id) REFERENCES class_templates(id) ON DELETE RESTRICT,
    CONSTRAINT chk_active_classes_student_count CHECK (student_count <= capacity)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_active_classes_academic_year_id ON active_classes(academic_year_id);
CREATE INDEX idx_active_classes_template_id ON active_classes(class_template_id);
CREATE INDEX idx_active_classes_is_active ON active_classes(is_active);
CREATE INDEX idx_active_classes_deleted_at ON active_classes(deleted_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS active_classes;
DROP TABLE IF EXISTS class_templates;
DROP TABLE IF EXISTS majors;
DROP TABLE IF EXISTS semesters;
DROP TABLE IF EXISTS cohorts;
DROP TABLE IF EXISTS academic_years;
-- +goose StatementEnd
