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

//RESTLoadbalancerListenerUpdateRequest for update request for REST.
type RESTLoadbalancerListenerUpdateRequest struct {
	Data map[string]interface{} `json:"loadbalancer-listener"`
}

//RESTCreateLoadbalancerListener handle a Create REST service.
func (service *ContrailService) RESTCreateLoadbalancerListener(c echo.Context) error {
	requestData := &models.CreateLoadbalancerListenerRequest{
		LoadbalancerListener: models.MakeLoadbalancerListener(),
	}
	if err := c.Bind(requestData); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"resource": "loadbalancer_listener",
		}).Debug("bind failed on create")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format")
	}
	ctx := c.Request().Context()
	response, err := service.CreateLoadbalancerListener(ctx, requestData)
	if err != nil {
		return common.ToHTTPError(err)
	}
	return c.JSON(http.StatusCreated, response)
}

//CreateLoadbalancerListener handle a Create API
func (service *ContrailService) CreateLoadbalancerListener(
	ctx context.Context,
	request *models.CreateLoadbalancerListenerRequest) (*models.CreateLoadbalancerListenerResponse, error) {
	model := request.LoadbalancerListener
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
			return db.CreateLoadbalancerListener(ctx, tx, request)
		}); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"resource": "loadbalancer_listener",
		}).Debug("db create failed on create")
		return nil, common.ErrorInternal
	}
	return &models.CreateLoadbalancerListenerResponse{
		LoadbalancerListener: request.LoadbalancerListener,
	}, nil
}

//RESTUpdateLoadbalancerListener handles a REST Update request.
func (service *ContrailService) RESTUpdateLoadbalancerListener(c echo.Context) error {
	//id := c.Param("id")
	request := &models.UpdateLoadbalancerListenerRequest{}
	if err := c.Bind(request); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"resource": "loadbalancer_listener",
		}).Debug("bind failed on update")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format")
	}
	ctx := c.Request().Context()
	response, err := service.UpdateLoadbalancerListener(ctx, request)
	if err != nil {
		return common.ToHTTPError(err)
	}
	return c.JSON(http.StatusOK, response)
}

//UpdateLoadbalancerListener handles a Update request.
func (service *ContrailService) UpdateLoadbalancerListener(
	ctx context.Context,
	request *models.UpdateLoadbalancerListenerRequest) (*models.UpdateLoadbalancerListenerResponse, error) {
	model := request.LoadbalancerListener
	if model == nil {
		return nil, common.ErrorBadRequest("Update body is empty")
	}
	if err := common.DoInTransaction(
		service.DB,
		func(tx *sql.Tx) error {
			return db.UpdateLoadbalancerListener(ctx, tx, request)
		}); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"resource": "loadbalancer_listener",
		}).Debug("db update failed")
		return nil, common.ErrorInternal
	}
	return &models.UpdateLoadbalancerListenerResponse{
		LoadbalancerListener: model,
	}, nil
}

//RESTDeleteLoadbalancerListener delete a resource using REST service.
func (service *ContrailService) RESTDeleteLoadbalancerListener(c echo.Context) error {
	id := c.Param("id")
	request := &models.DeleteLoadbalancerListenerRequest{
		ID: id,
	}
	ctx := c.Request().Context()
	_, err := service.DeleteLoadbalancerListener(ctx, request)
	if err != nil {
		return common.ToHTTPError(err)
	}
	return c.JSON(http.StatusNoContent, nil)
}

//DeleteLoadbalancerListener delete a resource.
func (service *ContrailService) DeleteLoadbalancerListener(ctx context.Context, request *models.DeleteLoadbalancerListenerRequest) (*models.DeleteLoadbalancerListenerResponse, error) {
	if err := common.DoInTransaction(
		service.DB,
		func(tx *sql.Tx) error {
			return db.DeleteLoadbalancerListener(ctx, tx, request)
		}); err != nil {
		log.WithField("err", err).Debug("error deleting a resource")
		return nil, common.ErrorInternal
	}
	return &models.DeleteLoadbalancerListenerResponse{
		ID: request.ID,
	}, nil
}

//RESTGetLoadbalancerListener a REST Get request.
func (service *ContrailService) RESTGetLoadbalancerListener(c echo.Context) error {
	id := c.Param("id")
	request := &models.GetLoadbalancerListenerRequest{
		ID: id,
	}
	ctx := c.Request().Context()
	response, err := service.GetLoadbalancerListener(ctx, request)
	if err != nil {
		return common.ToHTTPError(err)
	}
	return c.JSON(http.StatusOK, response)
}

//GetLoadbalancerListener a Get request.
func (service *ContrailService) GetLoadbalancerListener(ctx context.Context, request *models.GetLoadbalancerListenerRequest) (response *models.GetLoadbalancerListenerResponse, err error) {
	spec := &models.ListSpec{
		Limit: 1,
		Filter: models.Filter{
			"uuid": []string{request.ID},
		},
	}
	listRequest := &models.ListLoadbalancerListenerRequest{
		Spec: spec,
	}
	var result *models.ListLoadbalancerListenerResponse
	if err := common.DoInTransaction(
		service.DB,
		func(tx *sql.Tx) error {
			result, err = db.ListLoadbalancerListener(ctx, tx, listRequest)
			return err
		}); err != nil {
		return nil, common.ErrorInternal
	}
	if len(result.LoadbalancerListeners) == 0 {
		return nil, common.ErrorNotFound
	}
	response = &models.GetLoadbalancerListenerResponse{
		LoadbalancerListener: result.LoadbalancerListeners[0],
	}
	return response, nil
}

//RESTListLoadbalancerListener handles a List REST service Request.
func (service *ContrailService) RESTListLoadbalancerListener(c echo.Context) error {
	var err error
	spec := common.GetListSpec(c)
	request := &models.ListLoadbalancerListenerRequest{
		Spec: spec,
	}
	ctx := c.Request().Context()
	response, err := service.ListLoadbalancerListener(ctx, request)
	if err != nil {
		return common.ToHTTPError(err)
	}
	return c.JSON(http.StatusOK, response)
}

//ListLoadbalancerListener handles a List service Request.
func (service *ContrailService) ListLoadbalancerListener(
	ctx context.Context,
	request *models.ListLoadbalancerListenerRequest) (response *models.ListLoadbalancerListenerResponse, err error) {
	if err := common.DoInTransaction(
		service.DB,
		func(tx *sql.Tx) error {
			response, err = db.ListLoadbalancerListener(ctx, tx, request)
			return err
		}); err != nil {
		return nil, common.ErrorInternal
	}
	return response, nil
}
