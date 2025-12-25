-- 手工迁移脚本：角色-权限关联表（多对多映射）
-- 作用：关联角色表与权限表，实现RBAC模型中角色对权限的聚合管理
-- 适配项目：智算博弈场（GAME模块）、灵感卡片（CARD模块）
-- 执行依赖：需先执行 role.sql（角色表）、permissions.sql（权限表）
USE ddz;

-- 1. 创建角色-权限关联表
CREATE TABLE IF NOT EXISTS role_permissions (
    -- 关联字段：角色ID（关联roles表主键）
                                                role_id BIGINT UNSIGNED NOT NULL COMMENT '角色ID（关联roles表id）',
    -- 关联字段：权限ID（关联permissions表主键）
                                                perm_id BIGINT UNSIGNED NOT NULL COMMENT '权限ID（关联permissions表id）',
    -- 审计字段：创建时间（精确到毫秒，与现有表格式对齐）
                                                created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '关联创建时间',
    -- 审计字段：更新时间（精确到毫秒，用于记录权限调整时间）
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '关联更新时间',
    -- 主键约束：角色-权限组合唯一（避免重复绑定）
    PRIMARY KEY (role_id, perm_id),
    -- 外键约束：角色ID关联roles表，角色删除时级联删除关联记录
    CONSTRAINT fk_role_perm_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    -- 外键约束：权限ID关联permissions表，权限删除时级联删除关联记录
    CONSTRAINT fk_role_perm_perm FOREIGN KEY (perm_id) REFERENCES permissions(id) ON DELETE CASCADE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT '角色-权限关联表（实现角色与权限的多对多映射）';


-- 2. 初始化角色-权限绑定（贴合项目业务场景）
-- 说明：按角色职责分配权限，避免硬编码重复授权，支持后续动态调整
INSERT IGNORE INTO role_permissions (role_id, perm_id)
-- 子查询：关联roles表与permissions表，按角色名批量绑定权限
SELECT
    r.id AS role_id,        -- 角色ID（来自roles表）
    p.id AS perm_id         -- 权限ID（来自permissions表）
FROM roles r
         CROSS JOIN permissions p
WHERE
   -- （1）超级管理员（admin）：拥有所有模块权限（全量授权）
    (r.name = 'admin')

   -- （2）开发者（developer）：专注智算博弈场Bot开发+灵感卡片高级功能
   OR (r.name = 'developer' AND p.module IN ('GAME', 'CARD')
    AND p.perm_code NOT IN ('SYS_USER_MANAGE', 'SYS_ROLE_MANAGE', 'SYS_PERM_MANAGE'))

   -- （3）普通用户（viewer）：基础功能权限（人机对战、卡片生成、查看排行榜）
   OR (r.name = 'viewer' AND p.perm_code IN (
                                             'GAME_PVE',          -- 人机对战（智算博弈场核心体验）
                                             'GAME_RANK_VIEW',    -- 排行榜查看（激励机制配套）
                                             'CARD_GENERATE',     -- 灵感卡片生成（基础创作功能）
                                             'CARD_TOKEN_VIEW'    -- Token消耗查看（计量透明化）
    ))

   -- （4）运营人员（operator）：活动管理+数据查看（无开发/系统权限）
   OR (r.name = 'operator' AND p.perm_code IN (
                                               'GAME_COMPUTE_CLAIM',-- 算力券发放（激励活动操作）
                                               'CARD_HISTORY_VIEW', -- 卡片历史查看（内容审计）
                                               'GAME_REPLAY_VIEW'   -- 对战复盘查看（问题排查）
    ));


-- 3. 索引优化（提升权限校验查询效率）
-- 场景：后端校验用户权限时，常通过角色ID查询关联权限（如：SELECT perm_code FROM ... WHERE role_id = ?）
CREATE INDEX idx_role_perm_role ON role_permissions(role_id);
-- 场景：统计权限关联的角色数量时，通过权限ID查询（如：SELECT COUNT(role_id) FROM ... WHERE perm_id = ?）
CREATE INDEX idx_role_perm_perm ON role_permissions(perm_id);