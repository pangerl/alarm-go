CREATE TABLE public.qw_project_doris_log (
	id serial4 NOT NULL,
	project_name varchar(64) NOT NULL,
	project_id int4 NULL,
	job_failure_count int4 NULL,
	staff_count int4 NULL,
	use_analyse_count int4 NULL,
	customer_group_count int4 NULL,
	online_be_nodes int4 NULL,
	total_be_nodes int4 NULL,
	create_by varchar(64) NULL,
	create_time timestamp NULL,
	update_by varchar(64) NULL,
	update_time timestamp NULL,
	inspection_date date NOT NULL,
    UNIQUE (project_name, inspection_date),
	CONSTRAINT qw_project_doris_log_pkey PRIMARY KEY (id)
);

COMMENT ON COLUMN public.qw_project_doris_log.project_name IS '项目名称';
COMMENT ON COLUMN public.qw_project_doris_log.job_failure_count IS 'Job失败数';
COMMENT ON COLUMN public.qw_project_doris_log.staff_count IS '员工统计表前一天增量数据';
COMMENT ON COLUMN public.qw_project_doris_log.use_analyse_count IS '使用分析表前一天增量数据';
COMMENT ON COLUMN public.qw_project_doris_log.customer_group_count IS '客户群统计表前一天增量数据';
COMMENT ON COLUMN public.qw_project_doris_log.online_be_nodes IS '当前在线的 BE 节点数量';
COMMENT ON COLUMN public.qw_project_doris_log.total_be_nodes IS 'BE 节点总数';
COMMENT ON COLUMN public.qw_project_doris_log.create_by IS '记录创建者信息';
COMMENT ON COLUMN public.qw_project_doris_log.create_time IS '记录创建时间';
COMMENT ON COLUMN public.qw_project_doris_log.update_by IS '记录更新者信息';
COMMENT ON COLUMN public.qw_project_doris_log.update_time IS '记录更新时间';
COMMENT ON COLUMN public.qw_project_doris_log.inspection_date IS '检查日期，字符串格式，如 YYYY-MM-DD';