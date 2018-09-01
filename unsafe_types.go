package rundeck

import "time"

// ScheduledExecution is an entry in the scheduled_execution table
type ScheduledExecution struct {
	ID                         int64      `gorm:"column:id"`
	Version                    int64      `gorm:"column:version"`
	ArgString                  *string    `gorm:"column:arg_string"`
	DateCreated                time.Time  `gorm:"column:date_created"`
	DayOfMonth                 *string    `gorm:"column:day_of_month"`
	DayOfWeek                  *string    `gorm:"column:day_of_week"`
	DefaultTab                 *string    `gorm:"column:default_tab"`
	Description                *string    `gorm:"column:description"`
	DoNodeDispatch             *bool      `gorm:"column:do_nodedispatch"`
	ExecCount                  *int64     `gorm:"column:exec_count"`
	ExecutionEnabled           *bool      `gorm:"column:execution_enabled"`
	Filter                     *string    `gorm:"column:filter"`
	GroupPath                  *string    `gorm:"column:group_path"`
	Hour                       *string    `gorm:"column:hour"`
	JobName                    string     `gorm:"column:job_name"`
	LastUpdated                time.Time  `gorm:"column:last_updated"`
	LogOutputThreshold         *string    `gorm:"column:log_output_threshold"`
	LogOutputThresholdAction   *string    `gorm:"column:log_output_threshold_action"`
	LogOutputThresholdStatus   *string    `gorm:"column:log_output_threshold_status"`
	LogLevel                   *string    `gorm:"column:loglevel"`
	Minute                     **string   `gorm:"column:minute"`
	Month                      *string    `gorm:"column:month"`
	MultipleExecutions         *bool      `gorm:"column:multiple_executions"`
	NextExecution              *time.Time `gorm:"column:next_execution"`
	NodeExclude                *string    `gorm:"column:node_exclude"`
	NodeExcludeName            *string    `gorm:"column:node_exclude_name"`
	NodeExcludeOSArch          *string    `gorm:"column:node_exclude_os_arch"`
	NodeExcludeOSFamily        *string    `gorm:"column:node_exclude_os_family"`
	NodeExcludeOSName          *string    `gorm:"column:node_exclude_os_name"`
	NodeExcludeOSVersion       *string    `gorm:"column:node_exclude_os_version"`
	NodeExcludePrecedence      *bool      `gorm:"column:node_exclude_precedence"`
	NodeExcludeTags            *string    `gorm:"column:node_exclude_tags"`
	NodeFilterEditable         *bool      `gorm:"column:node_filter_editable"`
	NodeInclude                *string    `gorm:"column:node_include"`
	NodeIncludeName            *string    `gorm:"column:node_include_name"`
	NodeIncludeOSArch          *string    `gorm:"column:node_include_os_arch"`
	NodeIncludeOSFamily        *string    `gorm:"column:node_include_os_family"`
	NodeIncludeOSName          *string    `gorm:"column:node_include_os_name"`
	NodeIncludeOSVersion       *string    `gorm:"column:node_include_os_version"`
	NodeIncludePrecedence      *bool      `gorm:"column:node_include_precedence"`
	NodeIncludeTags            *string    `gorm:"column:node_include_tags"`
	NodeKeepGoing              *bool      `gorm:"column:node_keepgoing"`
	NodeRankAttribute          *string    `gorm:"column:node_rank_attribute"`
	NodeRankOrderAscending     *bool      `gorm:"column:node_rank_order_ascending"`
	NodeThreadcount            *int       `gorm:"column:node_threadcount"`
	NodeThreadcountDynamic     *string    `gorm:"column:node_threadcount_dynamic"`
	NodesSelectedByDefault     *bool      `gorm:"column:nodes_selected_by_default"`
	NotifyAvgDurationThreshold *string    `gorm:"column:notify_avg_duration_threshold"`
	OrchestratorID             *int64     `gorm:"column:orchestrator_id"`
	Project                    string     `gorm:"column:project"`
	RefExecCount               *int64     `gorm:"column:ref_exec_count"`
	Retry                      *string    `gorm:"column:retry"`
	RetryDelay                 *string    `gorm:"column:retry_delay"`
	ScheduleEnabled            *bool      `gorm:"column:schedule_enabled"`
	Scheduled                  bool       `gorm:"column:scheduled"`
	Seconds                    *string    `gorm:"column:seconds"`
	ServerNodeUUID             *string    `gorm:"column:server_nodeuuid"`
	SuccessOnEmptyNodeFilter   *bool      `gorm:"column:success_on_empty_node_filter"`
	TimeZone                   *string    `gorm:"column:time_zone"`
	Timeout                    *string    `gorm:"column:timeout"`
	TotalTime                  *int64     `gorm:"column:total_time"`
	RDUser                     *string    `gorm:"column:rduser"`
	UserRoleList               *string    `gorm:"column:user_role_list"`
	UUID                       *string    `gorm:"column:uuid"`
	WorkflowID                 *int64     `gorm:"column:workflow_id"`
	Year                       *string    `gorm:"column:year"`
	MaxMultipleExecutions      *string    `gorm:"column:max_multiple_executions"`
}

// TableName ...
func (ScheduledExecution) TableName() string {
	return "scheduled_execution"
}
