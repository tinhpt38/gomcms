package system

import (
	"context"

	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type initApi struct{}

const initOrderApi = system.InitOrderSystem + 1

// auto run
func init() {
	system.RegisterInit(initOrderApi, &initApi{})
}

func (i initApi) InitializerName() string {
	return sysModel.SysApi{}.TableName()
}

func (i *initApi) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&sysModel.SysApi{})
}

func (i *initApi) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&sysModel.SysApi{})
}

func (i *initApi) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	entities := []sysModel.SysApi{
		{ApiGroup: "jwt", Method: "POST", Path: "/jwt/jsonInBlacklist", Description: "Add JWT to blacklist (logout, required)"},

		{ApiGroup: "System User", Method: "DELETE", Path: "/user/deleteUser", Description: "Delete User"},
		{ApiGroup: "System User", Method: "POST", Path: "/user/admin_register", Description: "User Registration"},
		{ApiGroup: "System User", Method: "POST", Path: "/user/getUserList", Description: "Get User List"},
		{ApiGroup: "System User", Method: "PUT", Path: "/user/setUserInfo", Description: "Set User Information"},
		{ApiGroup: "System User", Method: "PUT", Path: "/user/setSelfInfo", Description: "Set Self Information (required)"},
		{ApiGroup: "System User", Method: "GET", Path: "/user/getUserInfo", Description: "Get Self Information (required)"},
		{ApiGroup: "System User", Method: "POST", Path: "/user/setUserAuthorities", Description: "Set User Authorities"},
		{ApiGroup: "System User", Method: "POST", Path: "/user/changePassword", Description: "Change Password (recommended)"},
		{ApiGroup: "System User", Method: "POST", Path: "/user/setUserAuthority", Description: "Modify User Role (required)"},
		{ApiGroup: "System User", Method: "POST", Path: "/user/resetPassword", Description: "Reset User Password"},

		{ApiGroup: "API", Method: "POST", Path: "/api/createApi", Description: "Create API"},
		{ApiGroup: "API", Method: "POST", Path: "/api/deleteApi", Description: "Delete API"},
		{ApiGroup: "API", Method: "POST", Path: "/api/updateApi", Description: "Update API"},
		{ApiGroup: "API", Method: "POST", Path: "/api/getApiList", Description: "Get API List"},
		{ApiGroup: "API", Method: "POST", Path: "/api/getAllApis", Description: "Get All APIs"},
		{ApiGroup: "API", Method: "POST", Path: "/api/getApiById", Description: "Get API Details"},
		{ApiGroup: "API", Method: "DELETE", Path: "/api/deleteApisByIds", Description: "Batch Delete APIs"},
		{ApiGroup: "API", Method: "GET", Path: "/api/syncApi", Description: "Get Pending Sync APIs"},
		{ApiGroup: "API", Method: "GET", Path: "/api/getApiGroups", Description: "Get Route Groups"},
		{ApiGroup: "API", Method: "POST", Path: "/api/enterSyncApi", Description: "Confirm Sync API"},
		{ApiGroup: "API", Method: "POST", Path: "/api/ignoreApi", Description: "Ignore API"},

		{ApiGroup: "Role", Method: "POST", Path: "/authority/copyAuthority", Description: "Copy Role"},
		{ApiGroup: "Role", Method: "POST", Path: "/authority/createAuthority", Description: "Create Role"},
		{ApiGroup: "Role", Method: "POST", Path: "/authority/deleteAuthority", Description: "Delete Role"},
		{ApiGroup: "Role", Method: "PUT", Path: "/authority/updateAuthority", Description: "Update Role Information"},
		{ApiGroup: "Role", Method: "POST", Path: "/authority/getAuthorityList", Description: "Get Role List"},
		{ApiGroup: "Role", Method: "POST", Path: "/authority/setDataAuthority", Description: "Set Role Resource Permissions"},

		{ApiGroup: "Casbin", Method: "POST", Path: "/casbin/updateCasbin", Description: "Change Role API Permissions"},
		{ApiGroup: "Casbin", Method: "POST", Path: "/casbin/getPolicyPathByAuthorityId", Description: "Get Permission List"},

		{ApiGroup: "Menu", Method: "POST", Path: "/menu/addBaseMenu", Description: "Add Menu"},
		{ApiGroup: "Menu", Method: "POST", Path: "/menu/getMenu", Description: "Get Menu Tree (required)"},
		{ApiGroup: "Menu", Method: "POST", Path: "/menu/deleteBaseMenu", Description: "Delete Menu"},
		{ApiGroup: "Menu", Method: "POST", Path: "/menu/updateBaseMenu", Description: "Update Menu"},
		{ApiGroup: "Menu", Method: "POST", Path: "/menu/getBaseMenuById", Description: "Get Menu by ID"},
		{ApiGroup: "Menu", Method: "POST", Path: "/menu/getMenuList", Description: "Get Basic Menu List"},
		{ApiGroup: "Menu", Method: "POST", Path: "/menu/getBaseMenuTree", Description: "Get User Dynamic Routes"},
		{ApiGroup: "Menu", Method: "POST", Path: "/menu/getMenuAuthority", Description: "Get Specified Role Menu"},
		{ApiGroup: "Menu", Method: "POST", Path: "/menu/addMenuAuthority", Description: "Add Menu and Role Association"},

		{ApiGroup: "Chunk Upload", Method: "GET", Path: "/fileUploadAndDownload/findFile", Description: "Find Target File (Fast Upload)"},
		{ApiGroup: "Chunk Upload", Method: "POST", Path: "/fileUploadAndDownload/breakpointContinue", Description: "Breakpoint Continue"},
		{ApiGroup: "Chunk Upload", Method: "POST", Path: "/fileUploadAndDownload/breakpointContinueFinish", Description: "Breakpoint Continue Finish"},
		{ApiGroup: "Chunk Upload", Method: "POST", Path: "/fileUploadAndDownload/removeChunk", Description: "Remove Uploaded File"},

		{ApiGroup: "File Upload and Download", Method: "POST", Path: "/fileUploadAndDownload/upload", Description: "File Upload (recommended)"},
		{ApiGroup: "File Upload and Download", Method: "POST", Path: "/fileUploadAndDownload/deleteFile", Description: "Delete File"},
		{ApiGroup: "File Upload and Download", Method: "POST", Path: "/fileUploadAndDownload/editFileName", Description: "Edit File Name or Remark"},
		{ApiGroup: "File Upload and Download", Method: "POST", Path: "/fileUploadAndDownload/getFileList", Description: "Get Uploaded File List"},

		{ApiGroup: "System Service", Method: "POST", Path: "/system/getServerInfo", Description: "Get Server Information"},
		{ApiGroup: "System Service", Method: "POST", Path: "/system/getSystemConfig", Description: "Get Configuration File Content"},
		{ApiGroup: "System Service", Method: "POST", Path: "/system/setSystemConfig", Description: "Set Configuration File Content"},

		{ApiGroup: "Customer", Method: "PUT", Path: "/customer/customer", Description: "Update Customer"},
		{ApiGroup: "Customer", Method: "POST", Path: "/customer/customer", Description: "Create Customer"},
		{ApiGroup: "Customer", Method: "DELETE", Path: "/customer/customer", Description: "Delete Customer"},
		{ApiGroup: "Customer", Method: "GET", Path: "/customer/customer", Description: "Get Single Customer"},
		{ApiGroup: "Customer", Method: "GET", Path: "/customer/customerList", Description: "Get Customer List"},

		{ApiGroup: "Code Generator", Method: "GET", Path: "/autoCode/getDB", Description: "Get All Databases"},
		{ApiGroup: "Code Generator", Method: "GET", Path: "/autoCode/getTables", Description: "Get Database Tables"},
		{ApiGroup: "Code Generator", Method: "POST", Path: "/autoCode/createTemp", Description: "Automated Code Generation"},
		{ApiGroup: "Code Generator", Method: "POST", Path: "/autoCode/preview", Description: "Preview Automated Code"},
		{ApiGroup: "Code Generator", Method: "GET", Path: "/autoCode/getColumn", Description: "Get All Fields of Selected Table"},
		{ApiGroup: "Code Generator", Method: "POST", Path: "/autoCode/installPlugin", Description: "Install Plugin"},
		{ApiGroup: "Code Generator", Method: "POST", Path: "/autoCode/pubPlug", Description: "Package Plugin"},

		{ApiGroup: "Template Configuration", Method: "POST", Path: "/autoCode/createPackage", Description: "Configure Template"},
		{ApiGroup: "Template Configuration", Method: "GET", Path: "/autoCode/getTemplates", Description: "Get Template Files"},
		{ApiGroup: "Template Configuration", Method: "POST", Path: "/autoCode/getPackage", Description: "Get All Templates"},
		{ApiGroup: "Template Configuration", Method: "POST", Path: "/autoCode/delPackage", Description: "Delete Template"},

		{ApiGroup: "Code Generator History", Method: "POST", Path: "/autoCode/getMeta", Description: "Get Meta Information"},
		{ApiGroup: "Code Generator History", Method: "POST", Path: "/autoCode/rollback", Description: "Rollback Auto-generated Code"},
		{ApiGroup: "Code Generator History", Method: "POST", Path: "/autoCode/getSysHistory", Description: "Query Rollback Records"},
		{ApiGroup: "Code Generator History", Method: "POST", Path: "/autoCode/delSysHistory", Description: "Delete Rollback Records"},
		{ApiGroup: "Code Generator History", Method: "POST", Path: "/autoCode/addFunc", Description: "Add Template Method"},

		{ApiGroup: "System Dictionary Detail", Method: "PUT", Path: "/sysDictionaryDetail/updateSysDictionaryDetail", Description: "Update Dictionary Content"},
		{ApiGroup: "System Dictionary Detail", Method: "POST", Path: "/sysDictionaryDetail/createSysDictionaryDetail", Description: "Create Dictionary Content"},
		{ApiGroup: "System Dictionary Detail", Method: "DELETE", Path: "/sysDictionaryDetail/deleteSysDictionaryDetail", Description: "Delete Dictionary Content"},
		{ApiGroup: "System Dictionary Detail", Method: "GET", Path: "/sysDictionaryDetail/findSysDictionaryDetail", Description: "Get Dictionary Content by ID"},
		{ApiGroup: "System Dictionary Detail", Method: "GET", Path: "/sysDictionaryDetail/getSysDictionaryDetailList", Description: "Get Dictionary Content List"},

		{ApiGroup: "System Dictionary", Method: "POST", Path: "/sysDictionary/createSysDictionary", Description: "Create Dictionary"},
		{ApiGroup: "System Dictionary", Method: "DELETE", Path: "/sysDictionary/deleteSysDictionary", Description: "Delete Dictionary"},
		{ApiGroup: "System Dictionary", Method: "PUT", Path: "/sysDictionary/updateSysDictionary", Description: "Update Dictionary"},
		{ApiGroup: "System Dictionary", Method: "GET", Path: "/sysDictionary/findSysDictionary", Description: "Get Dictionary by ID (recommended)"},
		{ApiGroup: "System Dictionary", Method: "GET", Path: "/sysDictionary/getSysDictionaryList", Description: "Get Dictionary List"},

		{ApiGroup: "Operation Record", Method: "POST", Path: "/sysOperationRecord/createSysOperationRecord", Description: "Create Operation Record"},
		{ApiGroup: "Operation Record", Method: "GET", Path: "/sysOperationRecord/findSysOperationRecord", Description: "Get Operation Record by ID"},
		{ApiGroup: "Operation Record", Method: "GET", Path: "/sysOperationRecord/getSysOperationRecordList", Description: "Get Operation Record List"},
		{ApiGroup: "Operation Record", Method: "DELETE", Path: "/sysOperationRecord/deleteSysOperationRecord", Description: "Delete Operation Record"},
		{ApiGroup: "Operation Record", Method: "DELETE", Path: "/sysOperationRecord/deleteSysOperationRecordByIds", Description: "Batch Delete Operation History"},

		{ApiGroup: "Resumable Upload (Plugin Version)", Method: "POST", Path: "/simpleUploader/upload", Description: "Plugin Version Chunk Upload"},
		{ApiGroup: "Resumable Upload (Plugin Version)", Method: "GET", Path: "/simpleUploader/checkFileMd5", Description: "File Integrity Verification"},
		{ApiGroup: "Resumable Upload (Plugin Version)", Method: "GET", Path: "/simpleUploader/mergeFileMd5", Description: "Upload Complete Merge File"},

		{ApiGroup: "Email", Method: "POST", Path: "/email/emailTest", Description: "Send Test Email"},
		{ApiGroup: "Email", Method: "POST", Path: "/email/sendEmail", Description: "Send Email"},

		{ApiGroup: "Button Permission", Method: "POST", Path: "/authorityBtn/setAuthorityBtn", Description: "Set Button Permission"},
		{ApiGroup: "Button Permission", Method: "POST", Path: "/authorityBtn/getAuthorityBtn", Description: "Get Existing Button Permissions"},
		{ApiGroup: "Button Permission", Method: "POST", Path: "/authorityBtn/canRemoveAuthorityBtn", Description: "Delete Button"},

		{ApiGroup: "Table Template", Method: "POST", Path: "/sysExportTemplate/createSysExportTemplate", Description: "Create Export Template"},
		{ApiGroup: "Table Template", Method: "DELETE", Path: "/sysExportTemplate/deleteSysExportTemplate", Description: "Delete Export Template"},
		{ApiGroup: "Table Template", Method: "DELETE", Path: "/sysExportTemplate/deleteSysExportTemplateByIds", Description: "Batch Delete Export Templates"},
		{ApiGroup: "Table Template", Method: "PUT", Path: "/sysExportTemplate/updateSysExportTemplate", Description: "Update Export Template"},
	}
	if err := db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, sysModel.SysApi{}.TableName()+"Table data initialization failed!")
	}
	next := context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}

func (i *initApi) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("path = ? AND method = ?", "/authorityBtn/canRemoveAuthorityBtn", "POST").
		First(&sysModel.SysApi{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
