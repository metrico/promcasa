package webmessages

// Messages Codes System
const (
	PleaseTryAgain  = "Please Try Again"
	OperationFailed = "Operation Failed"
)

// Messages Codes for Users
const (
	UserRequestFailed             = "failed to get Users"
	UserSettingsFailed            = "failed to get user settings"
	UserProfileFailed             = "failed to get user profile"
	MappingHepSubFailed           = "failed to get hepsub mapping schema"
	InsertDashboardFailed         = "failed to create a new dashboard"
	GetDashboardFailed            = "failed to get dashboards"
	GetDashboardListFailed        = "failed to get of list of dashboards"
	GetDBNodeListFailed           = "failed to get list of db nodes"
	GetAuthTypeListFailed         = "failed to get list of auths"
	DeleteDashboardFailed         = "failed to delete the dashboard"
	MappingSchemaFailed           = "failed to get mapping schema"
	HepSubRequestFailed           = "failed to get hep sub schema"
	GetAgentSubFailed             = "failed to get agent sub"
	GetAuthTokenFailed            = "failed to get auth token"
	DeleteAdvancedAgainstFailed   = "failed to delete advanced setting"
	GetAdvancedAgainstFailed      = "failed to get advanced setting"
	MappingSchemaByUUIDFailed     = "failed to get mapping schema by uuid"
	MappingRecreateFailed         = "failed to recreated all mappings"
	DeleteMappingSchemaFailed     = "failed to delete mapping schema"
	SmartHepProfileFailed         = "failed to get smart hep profile"
	DashboardNotExists            = "dashboard for the user doesn't exist"
	HomeDashboardNotExists        = "home dashboard for the user doesn't exist"
	Unauthorized                  = "Unauthorized"
	IncorrectPassword             = "incorrect password"
	UserCreationFailed            = "failed to create User"
	UserDeleteionFailed           = "failed to delete User"
	SuccessfullyCreatedUser       = "successfully created user"
	UserRequestFormatIncorrect    = "request format is not correct"
	UserGroupInformationError     = "couldn't get usergroup information"
	ResyncFormatIncorrect         = "resync format is not correct"
	ResyncBadSourceNode           = "bad source node selected"
	FileRequestFormatNotFound     = "file has been not found"
	AgentStorageNotFound          = "agent storage has been not found"
	AgentNotFoundData             = "agent couldn't find the data"
	StenographerException         = "stenographer data couldn't be retrieved"
	StenographerProxyException    = "stenographer proxy is not availble"
	UserRequestBadValueInside     = "request has not permitted values"
	ShareLinkGeneratorError       = "share link generate error"
	ShareLinkNotExistsError       = "share link doesn't exist or expired"
	ShareLinkBadDataError         = "share link bad data error"
	UserReportZipGenerationError  = "zip error generation"
	UIVersionFileNotExistsError   = "error version file couldn't be found"
	SwaggerFileNotExistsError     = "swagger file couldn't be found"
	ExportConfigNotActiveError    = "export config is not active"
	DownloadLogFileNotActiveError = "download log is not allowed. Please activate it in your config"
	BackupFileCreateError         = "couldn't create a backup file"
	BadValidationForParamterers   = "couldn't validate parameters"
	AllGood                       = "this is successful response"
	GetConfigNodeListFailed       = "failed to get list of config nodes"
	GrafanaProcessingError        = "grafana returned"
)