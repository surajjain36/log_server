package web

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

//WriteLog ...
func (s *Service) WriteLog(c *gin.Context) {
	var logData map[string]interface{}
	msg := "Something went wrong"
	statusCode := http.StatusBadRequest
	var currentTime = time.Now()

	if err := c.ShouldBindJSON(&logData); err != nil {
		log.Error("Error with i/p params: ", err.Error())
		s.responseWriter(c, gin.H{"message": msg}, statusCode)
		return
	}

	logData["CreatedAt"] = currentTime.UTC()
	logData["UpdatedAt"] = currentTime.UTC()

	if _, err := s.mdb.Insert("logs", &logData); err == nil {
		s.responseWriter(c, gin.H{"message": "Log is inserted successfully"}, http.StatusOK)
		return
	}

	s.responseWriter(c, gin.H{"message": msg}, statusCode)
	return
}

//ReadLog ...
func (s *Service) ReadLog(c *gin.Context) {
	log.Println("Hi")
	log.Println(c.DefaultQuery("source", ""))
	var logsData []map[string]interface{}
	statusCode := http.StatusBadRequest
	res := gin.H{"message": "Something went wrong", "data": nil}

	pn, err := strconv.Atoi(c.DefaultQuery("pn", "1"))
	if err != nil {
		log.WithError(err).Warn("Invalid page number")
		pn = 1
	}

	pp, err := strconv.Atoi(c.DefaultQuery("pp", "20"))
	if err != nil {
		log.WithError(err).Warn("Invalid per page value")
		pp = 20
	}

	if source := c.DefaultQuery("source", ""); source != "" {
		pipeline := bson.A{
			bson.M{"$match": bson.M{"source": source}},
			bson.M{"$skip": pp * (pn - 1)},
			bson.M{"$limit": pp},
		}
		log.Println(pipeline)
		if err := s.mdb.Aggregate("logs", pipeline, &logsData); err == nil {
			res["message"] = "Logs retrieved successfully"
			res["data"] = &logsData
			s.responseWriter(c, res, http.StatusOK)
			return
		}
	}

	s.responseWriter(c, res, statusCode)
	return
}
