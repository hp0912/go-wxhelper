SET NAMES utf8mb4;

-- AI 助手表
CREATE TABLE IF NOT EXISTS `t_ai_assistant` (
  `id`           varchar(36)  NOT NULL COMMENT '主键',
  `created_at`   datetime  DEFAULT NULL COMMENT '创建时间',
  `name`         varchar(50)  NOT NULL COMMENT '名称',
  `personality`  text DEFAULT NULL COMMENT '人设',
  `model`        varchar(50)  NOT NULL COMMENT '使用的模型',
  `enable`       tinyint(1)   NOT NULL DEFAULT 1 COMMENT '是否启用',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='AI助手';

-- 好友 / 群表
CREATE TABLE IF NOT EXISTS `t_friend` (
  `wxid`             varchar(64) NOT NULL COMMENT '微信原始Id',
  `custom_account`   varchar(64) DEFAULT NULL COMMENT '微信号',
  `nickname`         varchar(128) DEFAULT NULL COMMENT '昵称',
  `pinyin`           varchar(128) DEFAULT NULL COMMENT '昵称拼音首字母',
  `pinyin_all`       varchar(128) DEFAULT NULL COMMENT '昵称全拼',
  `last_active`      datetime     DEFAULT NULL COMMENT '最后活跃时间',
  `enable_ai`        tinyint(1)   NOT NULL DEFAULT 0 COMMENT '是否使用AI',
  `ai_model`         varchar(64) DEFAULT NULL COMMENT 'AI模型',
  `prompt`           text         COMMENT '提示词',
  `enable_chat_rank` tinyint(1)   NOT NULL DEFAULT 0 COMMENT '是否使用聊天排行',
  `enable_welcome`   tinyint(1)   NOT NULL DEFAULT 0 COMMENT '是否启用迎新',
  `enable_summary`   tinyint(1)   NOT NULL DEFAULT 0 COMMENT '是否启用总结',
  `enable_news`      tinyint(1)   NOT NULL DEFAULT 0 COMMENT '是否启用新闻',
  `clear_member`     tinyint(1)   NOT NULL DEFAULT 0 COMMENT '清理成员天数',
  `is_ok`            tinyint(1)   NOT NULL DEFAULT 0 COMMENT '是否正常',
  PRIMARY KEY (`wxid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='好友 / 群表';

-- 群成员表
CREATE TABLE IF NOT EXISTS `t_group_user` (
  `group_id`       varchar(64) NOT NULL COMMENT '群Id',
  `wxid`           varchar(64) NOT NULL COMMENT '成员微信Id',
  `account`        varchar(64) DEFAULT NULL COMMENT '账号',
  `head_image`     varchar(512) DEFAULT NULL COMMENT '头像',
  `nickname`       varchar(128) DEFAULT NULL COMMENT '昵称',
  `is_member`      tinyint(1)   NOT NULL DEFAULT 0 COMMENT '是否群成员',
  `is_admin`       tinyint(1)   NOT NULL DEFAULT 0 COMMENT '是否群主',
  `join_time`      datetime     DEFAULT NULL COMMENT '加入时间',
  `last_active`    datetime     DEFAULT NULL COMMENT '最后活跃时间',
  `leave_time`     datetime     DEFAULT NULL COMMENT '离开时间',
  `skip_chat_rank` tinyint(1)   NOT NULL DEFAULT 0 COMMENT '是否跳过聊天排行',
  PRIMARY KEY (`group_id`,`wxid`),
  KEY `idx_group_user_last_active` (`group_id`,`last_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='群成员';

-- 消息表
CREATE TABLE IF NOT EXISTS `t_message` (
  `msg_id`               bigint       NOT NULL COMMENT '消息Id',
  `create_time`          int          DEFAULT NULL COMMENT '发送时间戳',
  `create_at`            datetime     DEFAULT NULL COMMENT '发送时间',
  `type`                 int          DEFAULT NULL COMMENT '消息类型',
  `content`              text         COMMENT '内容',
  `display_full_content` text         COMMENT '显示的完整内容',
  `from_user`            varchar(64) DEFAULT NULL COMMENT '发送者',
  `group_user`           varchar(64) DEFAULT NULL COMMENT '群成员',
  `to_user`              varchar(64) DEFAULT NULL COMMENT '接收者',
  `raw`                  text         COMMENT '原始通知字符串',
  PRIMARY KEY (`msg_id`),
  KEY `idx_message_type` (`type`),
  KEY `idx_message_group_user` (`group_user`),
  KEY `idx_message_from_time` (`from_user`,`create_time`),
  KEY `idx_message_to_time` (`to_user`,`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='消息';

-- 插件数据表
CREATE TABLE IF NOT EXISTS `t_plugin_data` (
  `user_id`     varchar(64)  NOT NULL COMMENT '用户Id',
  `plugin_code` varchar(64)  NOT NULL COMMENT '插件编码',
  `data`        text         COMMENT '数据(JSON等)',
  PRIMARY KEY (`user_id`,`plugin_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='插件数据';