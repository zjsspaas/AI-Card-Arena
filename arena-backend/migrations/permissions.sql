-- 手工迁移脚本：创建权限表
-- 如需调整库名，请在执行前修改 `USE ddz;`
USE ddz;

CREATE TABLE IF NOT EXISTS permissions (
                                           id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
                                           perm_code VARCHAR(60) NOT NULL UNIQUE,
    perm_name VARCHAR(50) NOT NULL,
    module VARCHAR(30) NOT NULL,
    description VARCHAR(255),
    created_at DATETIME(3),
    updated_at DATETIME(3),
    PRIMARY KEY (id)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


-- 手工迁移脚本：创建角色-权限多对多关联表
CREATE TABLE IF NOT EXISTS role_permissions (
                                                role_id BIGINT UNSIGNED NOT NULL,
                                                perm_id BIGINT UNSIGNED NOT NULL,
                                                created_at DATETIME(3),
    updated_at DATETIME(3),
    PRIMARY KEY (role_id, perm_id),
    CONSTRAINT fk_role_permissions_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT fk_role_permissions_perm FOREIGN KEY (perm_id) REFERENCES permissions(id) ON DELETE CASCADE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


-- 可选：预置项目核心权限
INSERT IGNORE INTO permissions (perm_code, perm_name, module, description, created_at, updated_at)
VALUES
    ('GAME_PVE', '人机对战', 'GAME', '允许挑战平台内置Bot', NOW(), NOW()),
    ('GAME_PVP', '机机对战', 'GAME', '允许接入自定义Bot对战', NOW(), NOW()),
    ('GAME_API_ACCESS', 'Bot API接入', 'GAME', '允许调用Bot接入接口', NOW(), NOW()),
    ('CARD_GENERATE', '灵感卡片生成', 'CARD', '允许生成差异化方案', NOW(), NOW()),
    ('CARD_REFINE', '灵感卡片精炼', 'CARD', '允许优化已有方案', NOW(), NOW()),
    ('SYS_USER_MANAGE', '用户管理', 'SYS', '允许管理平台用户', NOW(), NOW());


-- 可选：给预设角色绑定权限（与roles.sql中的admin/developer/viewer对应）
INSERT IGNORE INTO role_permissions (role_id, perm_id, created_at, updated_at)
SELECT
    r.id AS role_id,
    p.id AS perm_id,
    NOW() AS created_at,
    NOW() AS updated_at
FROM roles r
         CROSS JOIN permissions p
WHERE
   -- 管理员：拥有所有权限
    (r.name = 'admin')
   -- 开发者：拥有GAME/CARD模块权限
   OR (r.name = 'developer' AND p.module IN ('GAME', 'CARD'))
   -- 查看者：仅拥有查看类权限
   OR (r.name = 'viewer' AND p.perm_code IN ('GAME_RANK_VIEW', 'CARD_HISTORY_VIEW'));gitjjjj