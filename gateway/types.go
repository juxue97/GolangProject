package main

import pb "github.com/juxue97/common/api"

type createOrderResponse struct {
	Order         *pb.Order `"json":order`
	RedirectToUrl string    `"json":redirectToUrl`
}
