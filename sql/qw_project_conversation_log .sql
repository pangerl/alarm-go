CREATE TABLE public.qw_project_conversation_log (
	id serial4 NOT NULL,
	project_name varchar(64) NOT NULL,
	project_id int4 NULL,
	corp_name varchar(64) NOT NULL,
	current_message_num int4 NULL,
	yesterday_message_num int4 NULL,
	create_by varchar(64) NULL,
	create_time timestamp NULL,
	update_by varchar(64) NULL,
	update_time timestamp NULL,
	inspection_date date NOT NULL,
    UNIQUE (corp_name, inspection_date),
	CONSTRAINT qw_project_conversation_log_pkey PRIMARY KEY (id)
);

COMMENT ON COLUMN public.qw_project_conversation_log.project_name IS '项目名称';
COMMENT ON COLUMN public.qw_project_conversation_log.corp_name IS '租户名称';
COMMENT ON COLUMN public.qw_project_conversation_log.current_message_num IS '当前会话数统计';
COMMENT ON COLUMN public.qw_project_conversation_log.yesterday_message_num IS '昨天会话数统计';
COMMENT ON COLUMN public.qw_project_conversation_log.create_by IS '记录创建者信息';
COMMENT ON COLUMN public.qw_project_conversation_log.create_time IS '记录创建时间';
COMMENT ON COLUMN public.qw_project_conversation_log.update_by IS '记录更新者信息';
COMMENT ON COLUMN public.qw_project_conversation_log.update_time IS '记录更新时间';
COMMENT ON COLUMN public.qw_project_conversation_log.inspection_date IS '检查日期，字符串格式，如 YYYY-MM-DD';