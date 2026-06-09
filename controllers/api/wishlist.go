package api

import (
	"TravelSphere/services"
	"TravelSphere/utils"
	"encoding/json"
	"net/http"

	"github.com/beego/beego/v2/core/validation"
)

type WishlistController struct {
	APIBaseController
}

type WishlistRequest struct {
	CountryName string `json:"country_name" valid:"Required"`
	Note string `json:"note"`
	Status string `json:"status" valid:"Required"`
}

func (c *WishlistController) CreateWishlist() {
	username := c.Ctx.Input.Header("Username")
	var req WishlistRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.SendError("Invalid payload structure", http.StatusBadRequest)
		return
	}
	validationEngine := validation.Validation{}
	ok, err := validationEngine.Valid(&req)
	if err != nil {
		c.SendError("Validation engine internal failure", http.StatusInternalServerError)
		return
	}
	if !utils.IsValidStatus(req.Status) {
		c.SendError("status: must be either Planned or Visited", http.StatusBadRequest)
		return
	}
	if !ok {
		for _, errItem := range validationEngine.Errors {
			c.SendError(errItem.Key+": "+errItem.Message, http.StatusBadRequest)
			return
		}
	}
	service := &services.WishlistService{}
	entry, err := service.AddToWishlist(username, req.CountryName, req.Note, req.Status)
	if err != nil {
		c.SendError(err.Error(), http.StatusBadRequest)
		return
	}
	c.SendSuccess("Destination successfully added to wishlist", entry, http.StatusCreated)
}

func (c *WishlistController) DeleteWishlist() {
	username := c.Ctx.Input.Header("Username")
	id := c.Ctx.Input.Param(":id")
	service := &services.WishlistService{}
	if err := service.DeleteWishlist(username, id); err != nil {
		c.SendError(err.Error(), http.StatusNotFound)
		return
	}
	c.SendSuccess("Wishlist entry deleted successfully", nil, http.StatusOK)
}
