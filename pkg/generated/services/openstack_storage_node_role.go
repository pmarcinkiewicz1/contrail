package services

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/Juniper/contrail/pkg/common"
	"github.com/Juniper/contrail/pkg/generated/db"
	"github.com/Juniper/contrail/pkg/generated/models"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"

	log "github.com/sirupsen/logrus"
)

//RESTOpenstackStorageNodeRoleUpdateRequest for update request for REST.
type RESTOpenstackStorageNodeRoleUpdateRequest struct {
	Data map[string]interface{} `json:"openstack-storage-node-role"`
}

//RESTCreateOpenstackStorageNodeRole handle a Create REST service.
func (service *ContrailService) RESTCreateOpenstackStorageNodeRole(c echo.Context) error {
	requestData := &models.CreateOpenstackStorageNodeRoleRequest{
		OpenstackStorageNodeRole: models.MakeOpenstackStorageNodeRole(),
	}
	if err := c.Bind(requestData); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"resource": "openstack_storage_node_role",
		}).Debug("bind failed on create")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format")
	}
	ctx := c.Request().Context()
	response, err := service.CreateOpenstackStorageNodeRole(ctx, requestData)
	if err != nil {
		return common.ToHTTPError(err)
	}
	return c.JSON(http.StatusCreated, response)
}

//CreateOpenstackStorageNodeRole handle a Create API
func (service *ContrailService) CreateOpenstackStorageNodeRole(
	ctx context.Context,
	request *models.CreateOpenstackStorageNodeRoleRequest) (*models.CreateOpenstackStorageNodeRoleResponse, error) {
	model := request.OpenstackStorageNodeRole
	if model.UUID == "" {
		model.UUID = uuid.NewV4().String()
	}
	auth := common.GetAuthCTX(ctx)
	if auth == nil {
		return nil, common.ErrorUnauthenticated
	}

	if model.FQName == nil {
		if model.DisplayName == "" {
			return nil, common.ErrorBadRequest("Both of FQName and Display Name is empty")
		}
		model.FQName = []string{auth.DomainID(), auth.ProjectID(), model.DisplayName}
	}
	model.Perms2 = &models.PermType2{}
	model.Perms2.Owner = auth.ProjectID()
	if err := common.DoInTransaction(
		service.DB,
		func(tx *sql.Tx) error {
			return db.CreateOpenstackStorageNodeRole(ctx, tx, request)
		}); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"resource": "openstack_storage_node_role",
		}).Debug("db create failed on create")
		return nil, common.ErrorInternal
	}
	return &models.CreateOpenstackStorageNodeRoleResponse{
		OpenstackStorageNodeRole: request.OpenstackStorageNodeRole,
	}, nil
}

//RESTUpdateOpenstackStorageNodeRole handles a REST Update request.
func (service *ContrailService) RESTUpdateOpenstackStorageNodeRole(c echo.Context) error {
	//id := c.Param("id")
	request := &models.UpdateOpenstackStorageNodeRoleRequest{}
	if err := c.Bind(request); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"resource": "openstack_storage_node_role",
		}).Debug("bind failed on update")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format")
	}
	ctx := c.Request().Context()
	response, err := service.UpdateOpenstackStorageNodeRole(ctx, request)
	if err != nil {
		return common.ToHTTPError(err)
	}
	return c.JSON(http.StatusOK, response)
}

//UpdateOpenstackStorageNodeRole handles a Update request.
func (service *ContrailService) UpdateOpenstackStorageNodeRole(
	ctx context.Context,
	request *models.UpdateOpenstackStorageNodeRoleRequest) (*models.UpdateOpenstackStorageNodeRoleResponse, error) {
	model := request.OpenstackStorageNodeRole
	if model == nil {
		return nil, common.ErrorBadRequest("Update body is empty")
	}
	if err := common.DoInTransaction(
		service.DB,
		func(tx *sql.Tx) error {
			return db.UpdateOpenstackStorageNodeRole(ctx, tx, request)
		}); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"resource": "openstack_storage_node_role",
		}).Debug("db update failed")
		return nil, common.ErrorInternal
	}
	return &models.UpdateOpenstackStorageNodeRoleResponse{
		OpenstackStorageNodeRole: model,
	}, nil
}

//RESTDeleteOpenstackStorageNodeRole delete a resource using REST service.
func (service *ContrailService) RESTDeleteOpenstackStorageNodeRole(c echo.Context) error {
	id := c.Param("id")
	request := &models.DeleteOpenstackStorageNodeRoleRequest{
		ID: id,
	}
	ctx := c.Request().Context()
	_, err := service.DeleteOpenstackStorageNodeRole(ctx, request)
	if err != nil {
		return common.ToHTTPError(err)
	}
	return c.JSON(http.StatusNoContent, nil)
}

//DeleteOpenstackStorageNodeRole delete a resource.
func (service *ContrailService) DeleteOpenstackStorageNodeRole(ctx context.Context, request *models.DeleteOpenstackStorageNodeRoleRequest) (*models.DeleteOpenstackStorageNodeRoleResponse, error) {
	if err := common.DoInTransaction(
		service.DB,
		func(tx *sql.Tx) error {
			return db.DeleteOpenstackStorageNodeRole(ctx, tx, request)
		}); err != nil {
		log.WithField("err", err).Debug("error deleting a resource")
		return nil, common.ErrorInternal
	}
	return &models.DeleteOpenstackStorageNodeRoleResponse{
		ID: request.ID,
	}, nil
}

//RESTGetOpenstackStorageNodeRole a REST Get request.
func (service *ContrailService) RESTGetOpenstackStorageNodeRole(c echo.Context) error {
	id := c.Param("id")
	request := &models.GetOpenstackStorageNodeRoleRequest{
		ID: id,
	}
	ctx := c.Request().Context()
	response, err := service.GetOpenstackStorageNodeRole(ctx, request)
	if err != nil {
		return common.ToHTTPError(err)
	}
	return c.JSON(http.StatusOK, response)
}

//GetOpenstackStorageNodeRole a Get request.
func (service *ContrailService) GetOpenstackStorageNodeRole(ctx context.Context, request *models.GetOpenstackStorageNodeRoleRequest) (response *models.GetOpenstackStorageNodeRoleResponse, err error) {
	spec := &models.ListSpec{
		Limit: 1,
		Filter: models.Filter{
			"uuid": []string{request.ID},
		},
	}
	listRequest := &models.ListOpenstackStorageNodeRoleRequest{
		Spec: spec,
	}
	var result *models.ListOpenstackStorageNodeRoleResponse
	if err := common.DoInTransaction(
		service.DB,
		func(tx *sql.Tx) error {
			result, err = db.ListOpenstackStorageNodeRole(ctx, tx, listRequest)
			return err
		}); err != nil {
		return nil, common.ErrorInternal
	}
	if len(result.OpenstackStorageNodeRoles) == 0 {
		return nil, common.ErrorNotFound
	}
	response = &models.GetOpenstackStorageNodeRoleResponse{
		OpenstackStorageNodeRole: result.OpenstackStorageNodeRoles[0],
	}
	return response, nil
}

//RESTListOpenstackStorageNodeRole handles a List REST service Request.
func (service *ContrailService) RESTListOpenstackStorageNodeRole(c echo.Context) error {
	var err error
	spec := common.GetListSpec(c)
	request := &models.ListOpenstackStorageNodeRoleRequest{
		Spec: spec,
	}
	ctx := c.Request().Context()
	response, err := service.ListOpenstackStorageNodeRole(ctx, request)
	if err != nil {
		return common.ToHTTPError(err)
	}
	return c.JSON(http.StatusOK, response)
}

//ListOpenstackStorageNodeRole handles a List service Request.
func (service *ContrailService) ListOpenstackStorageNodeRole(
	ctx context.Context,
	request *models.ListOpenstackStorageNodeRoleRequest) (response *models.ListOpenstackStorageNodeRoleResponse, err error) {
	if err := common.DoInTransaction(
		service.DB,
		func(tx *sql.Tx) error {
			response, err = db.ListOpenstackStorageNodeRole(ctx, tx, request)
			return err
		}); err != nil {
		return nil, common.ErrorInternal
	}
	return response, nil
}
