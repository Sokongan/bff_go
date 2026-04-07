package permission_factory

import (
	"fmt"

	"sso-bff/modules/audit"

	permission_factory_modules "sso-bff/modules/permission/factory/modules"
	permission_handler "sso-bff/modules/permission/handler"
	permission_sdk "sso-bff/modules/permission/sdk"
)

type Module struct {
	Service *permission_factory_modules.PermissionService
	Handler *permission_handler.PermissionHandler
}

func NewPermissionModule(
	sdk *permission_sdk.PermissionSDK,
	auditWriter audit.AuditWriter,
) (*Module, error) {
	if sdk == nil {
		return nil, fmt.Errorf("permission sdk is nil")
	}

	gw, err := permission_factory_modules.NewPermissionGateway(sdk)
	if err != nil {
		return nil, err
	}

	svc := permission_factory_modules.NewPermissionService(gw)
	h := permission_handler.NewPermissionHandler(svc, auditWriter)

	return &Module{
		Service: svc,
		Handler: h,
	}, nil
}
