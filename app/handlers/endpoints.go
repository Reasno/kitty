// Code generated by truss. DO NOT EDIT.
// Rerunning truss will overwrite this file.
// Version: d800079357
// Version Date: 2020-10-29T08:16:24Z

package handlers

import (
	// This Service
	"github.com/Reasno/kitty/app/svc"
	pb "github.com/Reasno/kitty/proto"
)

func NewEndpoints(service pb.AppServer) svc.Endpoints {

	// Endpoint domain.
	var (
		loginEndpoint      = svc.MakeLoginEndpoint(service)
		getcodeEndpoint    = svc.MakeGetCodeEndpoint(service)
		getinfoEndpoint    = svc.MakeGetInfoEndpoint(service)
		updateinfoEndpoint = svc.MakeUpdateInfoEndpoint(service)
		bindEndpoint       = svc.MakeBindEndpoint(service)
		unbindEndpoint     = svc.MakeUnbindEndpoint(service)
		refreshEndpoint    = svc.MakeRefreshEndpoint(service)
	)

	endpoints := svc.Endpoints{
		LoginEndpoint:      loginEndpoint,
		GetCodeEndpoint:    getcodeEndpoint,
		GetInfoEndpoint:    getinfoEndpoint,
		UpdateInfoEndpoint: updateinfoEndpoint,
		BindEndpoint:       bindEndpoint,
		UnbindEndpoint:     unbindEndpoint,
		RefreshEndpoint:    refreshEndpoint,
	}

	return endpoints
}
