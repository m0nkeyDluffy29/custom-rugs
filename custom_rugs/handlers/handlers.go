package handlers

import (
	"net/http"
	"strconv"
	"time"

	"custom_rugs/db"
	"custom_rugs/models"

	"github.com/gin-gonic/gin"
)

func SubmitRugRequest(c *gin.Context) {
	var request models.RugRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err := db.DB.Prepare("INSERT INTO CUSTOM_RUGS (NAME, EMAIL, DETAILS) VALUES (?, ?, ?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(request.Name, request.EMAIL, request.Details)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert request"})
		return
	}
	request.STATUS = "PENDING"
	id, err := res.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get last insert ID"})
		return
	}
	request.ID = int(id)
	c.JSON(http.StatusCreated, gin.H{"message": "Rug request submitted successfully!", "request": request})
}

func GetAllRugRequests(c *gin.Context) {
	rows, err := db.DB.Query("SELECT ID, NAME, EMAIL, DETAILS, STATUS, CREATED_AT FROM CUSTOM_RUGS ORDER BY CREATED_AT DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve requests"})
		return
	}
	defer rows.Close()

	var requests []models.RugRequest
	for rows.Next() {
		var req models.RugRequest
		var createdAtStr string
		if err := rows.Scan(&req.ID, &req.Name, &req.EMAIL, &req.Details, &req.STATUS, &createdAtStr); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan request"})
			return
		}
		req.CREATED_AT, err = time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			req.CREATED_AT = time.Time{}
		}
		requests = append(requests, req)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over requests"})
		return
	}
	c.JSON(http.StatusOK, requests)

}

func UpdateRugRequestStatus(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	var updatereq struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&updatereq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validStatuses := map[string]bool{"pending": true, "approved": true, "rejected": true, "in_progress": true, "completed": true}
	if _, ok := validStatuses[updatereq.Status]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
		return
	}
	res, err := db.DB.Exec("UPDATE CUSTOM_RUGS SET status = ? WHERE id = ?", updatereq.Status, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update request status"})
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get rows affected"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Request status updated successfully"})
}
