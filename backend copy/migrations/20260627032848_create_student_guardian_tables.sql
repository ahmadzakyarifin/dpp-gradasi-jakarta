-- +goose Up
-- +goose StatementBegin

CREATE TABLE students (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NULL,
    cohort_id BIGINT UNSIGNED NULL,
    current_active_class_id BIGINT UNSIGNED NULL,
    current_major_id BIGINT UNSIGNED NULL,
    nis VARCHAR(50) NOT NULL,
    nisn VARCHAR(50) NULL,
    nik VARCHAR(30) NULL,
    name VARCHAR(150) NOT NULL,
    gender ENUM('male', 'female') NULL,
    religion VARCHAR(50) NULL,
    birth_place VARCHAR(100) NULL,
    birth_date DATE NULL,
    address TEXT NULL,
    rt VARCHAR(10) NULL,
    rw VARCHAR(10) NULL,
    village VARCHAR(100) NULL,
    district VARCHAR(100) NULL,
    city VARCHAR(100) NULL,
    province VARCHAR(100) NULL,
    email VARCHAR(150) NULL,
    phone VARCHAR(30) NULL,
    entry_year SMALLINT UNSIGNED NULL,
    status ENUM('active', 'inactive', 'graduated', 'transferred', 'dropped_out') NOT NULL DEFAULT 'active',
    photo_path VARCHAR(255) NULL,
    description TEXT NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL,
    CONSTRAINT uq_students_nis UNIQUE (nis),
    CONSTRAINT uq_students_nisn UNIQUE (nisn),
    CONSTRAINT uq_students_nik UNIQUE (nik),
    CONSTRAINT fk_students_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT fk_students_cohort FOREIGN KEY (cohort_id) REFERENCES cohorts(id) ON DELETE RESTRICT,
    CONSTRAINT fk_students_current_active_class FOREIGN KEY (current_active_class_id) REFERENCES active_classes(id) ON DELETE SET NULL,
    CONSTRAINT fk_students_current_major FOREIGN KEY (current_major_id) REFERENCES majors(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_students_user_id ON students(user_id);
CREATE INDEX idx_students_cohort_id ON students(cohort_id);
CREATE INDEX idx_students_current_active_class_id ON students(current_active_class_id);
CREATE INDEX idx_students_current_major_id ON students(current_major_id);
CREATE INDEX idx_students_entry_year ON students(entry_year);
CREATE INDEX idx_students_status ON students(status);
CREATE INDEX idx_students_deleted_at ON students(deleted_at);

CREATE TABLE guardians (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NULL,
    name VARCHAR(150) NOT NULL,
    phone VARCHAR(30) NULL,
    email VARCHAR(150) NULL,
    nik VARCHAR(30) NULL,
    education VARCHAR(100) NULL,
    occupation VARCHAR(150) NULL,
    income_range VARCHAR(100) NULL,
    address TEXT NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL,
    CONSTRAINT uq_guardians_user_id UNIQUE (user_id),
    CONSTRAINT uq_guardians_nik UNIQUE (nik),
    CONSTRAINT fk_guardians_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_guardians_user_id ON guardians(user_id);
CREATE INDEX idx_guardians_phone ON guardians(phone);
CREATE INDEX idx_guardians_deleted_at ON guardians(deleted_at);

CREATE TABLE student_guardians (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    student_id BIGINT UNSIGNED NOT NULL,
    guardian_id BIGINT UNSIGNED NOT NULL,
    relationship ENUM('father', 'mother', 'guardian', 'grandfather', 'grandmother', 'other') NOT NULL,
    is_primary_contact TINYINT(1) NOT NULL DEFAULT 0,
    can_login TINYINT(1) NOT NULL DEFAULT 0,
    can_receive_notification TINYINT(1) NOT NULL DEFAULT 1,
    can_make_payment TINYINT(1) NOT NULL DEFAULT 1,
    can_view_invoice TINYINT(1) NOT NULL DEFAULT 1,
    can_open_support_ticket TINYINT(1) NOT NULL DEFAULT 1,
    is_emergency_contact TINYINT(1) NOT NULL DEFAULT 0,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    CONSTRAINT fk_student_guardians_student FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
    CONSTRAINT fk_student_guardians_guardian FOREIGN KEY (guardian_id) REFERENCES guardians(id) ON DELETE CASCADE,
    CONSTRAINT uq_student_guardians UNIQUE (student_id, guardian_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_student_guardians_student_id ON student_guardians(student_id);
CREATE INDEX idx_student_guardians_guardian_id ON student_guardians(guardian_id);
CREATE INDEX idx_student_guardians_primary ON student_guardians(student_id, is_primary_contact);

CREATE TABLE class_memberships (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    student_id BIGINT UNSIGNED NOT NULL,
    active_class_id BIGINT UNSIGNED NOT NULL,
    academic_year_id BIGINT UNSIGNED NOT NULL,
    semester_id BIGINT UNSIGNED NULL,
    attendance_number SMALLINT UNSIGNED NULL,
    start_date DATE NULL,
    end_date DATE NULL,
    status ENUM('active', 'moved', 'completed', 'graduated', 'inactive') NOT NULL DEFAULT 'active',
    note TEXT NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    CONSTRAINT fk_class_memberships_student FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
    CONSTRAINT fk_class_memberships_active_class FOREIGN KEY (active_class_id) REFERENCES active_classes(id) ON DELETE RESTRICT,
    CONSTRAINT fk_class_memberships_academic_year FOREIGN KEY (academic_year_id) REFERENCES academic_years(id) ON DELETE RESTRICT,
    CONSTRAINT fk_class_memberships_semester FOREIGN KEY (semester_id) REFERENCES semesters(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_class_memberships_student_id ON class_memberships(student_id);
CREATE INDEX idx_class_memberships_active_class_id ON class_memberships(active_class_id);
CREATE INDEX idx_class_memberships_academic_year_id ON class_memberships(academic_year_id);
CREATE INDEX idx_class_memberships_semester_id ON class_memberships(semester_id);
CREATE INDEX idx_class_memberships_status ON class_memberships(status);
CREATE INDEX idx_class_memberships_period_status ON class_memberships(student_id, academic_year_id, semester_id, status);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS class_memberships;
DROP TABLE IF EXISTS student_guardians;
DROP TABLE IF EXISTS guardians;
DROP TABLE IF EXISTS students;
-- +goose StatementEnd
