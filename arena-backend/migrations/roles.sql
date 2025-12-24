-- 手工迁移脚本：创建角色表和用户-角色多对多关联表
-- 如需调整库名，请在执行前修改 `USE ddz;`

USE ddz;

CREATE TABLE IF NOT EXISTS roles (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255),
    created_at DATETIME(3),
    updated_at DATETIME(3),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS user_roles (
    user_id BIGINT UNSIGNED NOT NULL,
    role_id BIGINT UNSIGNED NOT NULL,
    created_at DATETIME(3),
    updated_at DATETIME(3),
    PRIMARY KEY (user_id, role_id),
    CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- 可选：预置几个常用角色
INSERT IGNORE INTO roles (name, description, created_at, updated_at)
VALUES
    ('admin', '管理员', NOW(), NOW()),
    ('developer', '开发者', NOW(), NOW()),
    ('viewer', '查看者', NOW(), NOW());


