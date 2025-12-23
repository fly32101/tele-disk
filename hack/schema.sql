-- Schema for tele-disk (MySQL 5.7)

CREATE TABLE IF NOT EXISTS `users` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `username` VARCHAR(64) NOT NULL COMMENT '登录用户名',
  `password_hash` VARCHAR(255) NOT NULL COMMENT '密码哈希',
  `is_admin` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否管理员',
  `quota_bytes` BIGINT NOT NULL DEFAULT 0 COMMENT '配额，0 表示不限',
  `used_bytes` BIGINT NOT NULL DEFAULT 0 COMMENT '已用容量',
  `file_count` INT NOT NULL DEFAULT 0 COMMENT '文件数量',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用，0停用',
  `last_login_at` DATETIME DEFAULT NULL COMMENT '上次登录时间',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_users_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

CREATE TABLE IF NOT EXISTS `telegram_accounts` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` BIGINT NOT NULL COMMENT '所属用户ID（无外键）',
  `name` VARCHAR(64) NOT NULL COMMENT '账号名称',
  `bot_username` VARCHAR(64) NOT NULL COMMENT 'Bot 用户名（@xxx）',
  `bot_token_enc` TEXT NOT NULL COMMENT '加密后的 Bot Token（对称加密，可逆）',
  `bot_id` BIGINT DEFAULT NULL COMMENT 'Bot ID',
  `chat_id` BIGINT DEFAULT NULL COMMENT '存储文件用的 chat_id',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用，0停用',
  `is_default` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否默认上传账号',
  `last_error` TEXT COMMENT '最近一次错误信息',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_tg_user_bot` (`user_id`, `bot_username`),
  KEY `idx_tg_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Telegram 账号表';

CREATE TABLE IF NOT EXISTS `files` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` BIGINT NOT NULL COMMENT '归属用户ID（无外键）',
  `telegram_account_id` BIGINT DEFAULT NULL COMMENT '使用的 Telegram 账号ID（无外键）',
  `telegram_file_id` VARCHAR(255) NOT NULL COMMENT 'Telegram file_id（下载用）',
  `telegram_file_unique_id` VARCHAR(255) COMMENT 'Telegram file_unique_id',
  `telegram_message_id` BIGINT COMMENT '发送到 Telegram 的消息ID',
  `telegram_chat_id` BIGINT COMMENT '发送到的 chat_id',
  `file_name` VARCHAR(255) NOT NULL COMMENT '原始文件名',
  `size_bytes` BIGINT NOT NULL COMMENT '文件大小（字节）',
  `mime_type` VARCHAR(255) COMMENT 'MIME 类型',
  `md5` CHAR(32) COMMENT 'MD5 去重',
  `sha256` CHAR(64) COMMENT 'sha256 校验/去重',
  `storage_provider` VARCHAR(32) NOT NULL DEFAULT 'telegram' COMMENT '存储提供方标识，默认 telegram',
  `object_key` VARCHAR(512) NOT NULL COMMENT '存储对象 key/路径，默认使用 Telegram file_id',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1=ok,2=uploading,3=failed',
  `is_deleted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否删除（软删）',
  `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '软删时间',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_files_user_object` (`user_id`, `object_key`),
  UNIQUE KEY `uk_files_user_md5` (`user_id`, `md5`),
  KEY `idx_files_user_status` (`user_id`, `status`),
  KEY `idx_files_checksum` (`sha256`),
  KEY `idx_files_object` (`object_key`),
  KEY `idx_files_fileid` (`telegram_file_id`),
  KEY `idx_files_account` (`telegram_account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件表';

-- 可选：访问日志表（如需审计/限流），无外键
-- CREATE TABLE IF NOT EXISTS `file_access_logs` (
--   `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
--   `file_id` BIGINT NOT NULL COMMENT '文件ID',
--   `user_id` BIGINT DEFAULT NULL COMMENT '访问用户ID',
--   `token` CHAR(36) COMMENT '使用的 token',
--   `ip` VARCHAR(64) COMMENT '来源 IP',
--   `user_agent` TEXT COMMENT 'UA',
--   `status_code` INT COMMENT '返回状态码',
--   `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '访问时间',
--   PRIMARY KEY (`id`),
--   KEY `idx_access_file` (`file_id`, `created_at`)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件访问日志表';
