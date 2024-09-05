package handler

import (
	"fmt"
	"fundraising-web/campaign"
	"fundraising-web/helper"
	"fundraising-web/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
		return &campaignHandler{service}
}

//api/v1/campaigns
func(h *campaignHandler) GetCampaigns(c *gin.Context) {
		userID,_ := strconv.Atoi(c.Query("user_id"))

		campaigns, err := h.service.GetCampaigns(userID)
		if err != nil {
			response := helper.APIresponse("Error to get campaigns", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response := helper.APIresponse("List of campaigns", http.StatusBadRequest, "success", campaign.FormatCampaigns(campaigns))
		c.JSON(http.StatusOK, response)

}

func (h *campaignHandler) GetCampaign(c *gin.Context) {
	var input campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
			response := helper.APIresponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
	}

	campaignDetail, err := h.service.GetCampaignById(input)
	if err != nil {
		response := helper.APIresponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIresponse("Campaign detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}


func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput 

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIresponse("Failed to create campaign", http.StatusUnprocessableEntity, "success", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	
	input.User = currentUser

	newCampaign, err := h.service.CreateCampaign(input)
	if err != nil {
		response := helper.APIresponse("Failed to create campaign", http.StatusBadRequest, "success", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIresponse("Success to create campaign", http.StatusOK, "success", campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusBadRequest, response)
}

func (h *campaignHandler) UpdateCampaign(c *gin.Context) {
	var inputID campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&inputID)
	if err != nil {
			response := helper.APIresponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
	}

	var inputData campaign.CreateCampaignInput 

	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIresponse("Failed to update campaign", http.StatusUnprocessableEntity, "success", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	inputData.User = currentUser

	updatedCampaign, err := h.service.UpdateCampaign(inputID, inputData)
	if err != nil {
		response := helper.APIresponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIresponse("Success to update campaign", http.StatusOK, "success", campaign.FormatCampaign(updatedCampaign))
	c.JSON(http.StatusBadRequest, response)
}


func (h *campaignHandler) UploadImage (c *gin.Context) {
		var input campaign.CreateCampaignImageInput

		err := c.ShouldBind(&input)

		if err != nil {
				errors := helper.FormatValidationError(err)
				errorMessage := gin.H{"errors": errors}
		
				response := helper.APIresponse("Failed to upload campaign image", http.StatusUnprocessableEntity, "success", errorMessage)
				c.JSON(http.StatusUnprocessableEntity, response)
				return
		}

		currentUser := c.MustGet("currentUser").(user.User)
		input.User = currentUser
		userID := currentUser.ID

		file, err := c.FormFile("file")
		if err != nil {
			data := gin.H{"is_uploaded": false}
			response := helper.APIresponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
	
			c.JSON(http.StatusBadRequest, response)
			return 
		}
	
		path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

		err = c.SaveUploadedFile(file, path)
		if err != nil {
			data := gin.H{"is_uploaded": false}
			response := helper.APIresponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
	
			c.JSON(http.StatusBadRequest, response)
			return 
		}

		_,err = h.service.SaveCampaignImage(input, path)
		if err != nil {
			data := gin.H{"is_uploaded": false}
			response := helper.APIresponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
	
			c.JSON(http.StatusBadRequest, response)
			return 
		}
	
		data := gin.H{"is_uploaded": false}
		response := helper.APIresponse("Campaign image successfuly uploaded", http.StatusOK, "success", data)
	
		c.JSON(http.StatusOK, response)
		
}