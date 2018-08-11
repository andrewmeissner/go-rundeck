package rundeck

// APIVersion24 is defaulted to the specified api version
const APIVersion24 = 24

// EnvRundeckToken sets the name of the environment variable to read
const EnvRundeckToken = "RUNDECK_TOKEN"

// EnvRundeckServerURL is the name of the environment variable for the server url
const EnvRundeckServerURL = "RUNDECK_SERVER_URL"

const localRundeckURL = "http://localhost:4440"

const ExecutionStatusRunning ExecutionStatus = "running"
const ExecutionStatusSucceeded ExecutionStatus = "succeeded"
const ExecutionStatusFailed ExecutionStatus = "failed"
const ExecutionStatusAborted ExecutionStatus = "aborted"
const ExecutionStatusTimedout ExecutionStatus = "timedout"
const ExecutionStatusFailedWithRetry ExecutionStatus = "failed-with-retry"
const ExecutionStatusScheduled ExecutionStatus = "scheduled"
const ExecutionStatusOther ExecutionStatus = "other"
const ExecutionTypeScheduled ExecutionType = "scheduled"
const ExecutionTypeUser ExecutionType = "user"
const ExecutionTypeUserScheduled ExecutionType = "user-scheduled"
const BooleanDefault Boolean = 0
const BooleanFalse Boolean = 1
const BooleanTrue Boolean = 2
const JobLogLevelDebug LogLevel = "DEBUG"
const JobLogLevelVerbose LogLevel = "VERBOSE"
const JobLogLevelInfo LogLevel = "INFO"
const JobLogLevelWarn LogLevel = "WARN"
const JobLogLevelError LogLevel = "ERROR"
const JobFormatXML JobFormat = "xml"
const JobFormatYAML JobFormat = "yaml"
const DuplicateOptionSkip DuplicateOption = "skip"
const DuplicateOptionCreate DuplicateOption = "create"
const DuplicateOptionUpdate DuplicateOption = "update"
const UUIDOptionPreserve UUIDOption = "preserve"
const UUIDOptionRemove UUIDOption = "remove"
const ToggleKindExecution ToggleKind = "execution"
const ToggleKindSchedule ToggleKind = "schedule"
const FileStateTemp FileState = "temp"
const FileStateDeleted FileState = "deleted"
const FileStateExpired FileState = "expired"
const FileStateRetained FileState = "retained"
const ExecutionModeActive ExecutionMode = "active"
const ExecutionModePassive ExecutionMode = "passive"
