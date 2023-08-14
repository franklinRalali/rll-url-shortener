package ucase

import (
	"net/http"

	"github.com/ralali/rll-url-shortener/internal/appctx"
	"github.com/ralali/rll-url-shortener/internal/consts"
	"github.com/ralali/rll-url-shortener/internal/presentations"
	"github.com/ralali/rll-url-shortener/internal/ucase/contract"
)

type allEndpoints struct {
}

func NewAllEndpoints() contract.UseCase {
	return &allEndpoints{}
}

func (a *allEndpoints) Serve(data *appctx.Data) (response appctx.Response) {
	// put all the available endpoints here
	endpoints := presentations.AllEndpoints{
		Endpoints: []presentations.Endpoint{
			{
				Path:        "/",
				Method:      http.MethodGet,
				Description: "List all available endpoints",
			},
			{
				Path:        "/health",
				Method:      http.MethodGet,
				Description: "Get system health status",
			},
			{
				Path:        "/{url}",
				Method:      http.MethodPost,
				Description: "Create new shorten URL",
			},
			{
				Path:        "/{short_code}",
				Method:      http.MethodGet,
				Description: "Get shorten URL with its origin URL",
			},
			{
				Path:        "/{short_code}",
				Method:      http.MethodPut,
				Description: "Update a short URL (currently can only update the origin URL)",
			},
			{
				Path:        "/{short_code}",
				Method:      http.MethodDelete,
				Description: "Delete a short URL",
			},
			{
				Path:        "/{short_code}/stats",
				Method:      http.MethodGet,
				Description: "Get short URL statistics (currently only visit count)",
			},
		},
	}

	response.SetName(consts.ResponseSuccess)
	response.SetData(endpoints)

	return response
}
