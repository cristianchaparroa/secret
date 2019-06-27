package api

import (
	"fmt"
	"os"

	"github.com/cristianchaparroa/secret/initializer"
	"github.com/cristianchaparroa/secret/metrics"
	mw "github.com/cristianchaparroa/secret/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

// SecretServer is the server of secrets
type SecretServer struct {
	Engine  *gin.Engine
	DB      *gorm.DB
	Monitor *metrics.Monitor
}

// NewSecretServer returns a pointer to SecretServer
func NewSecretServer() *SecretServer {
	gin.ForceConsoleColor()
	r := gin.Default()
	r.Use(gin.Recovery())
	return &SecretServer{Engine: r}
}

// Setup is in charge to up all the required configurations in the server
func (ss *SecretServer) Setup() {
	ss.setupDB()
	ss.setupMetrics()
	ss.setupEndpoints()

}

// setupEndpoints is in charge to confugure the API endpoints
func (ss *SecretServer) setupEndpoints() {

	sh := NewSecretHandler(ss.DB)
	ss.Engine.GET("/v1/secret/:hash", sh.FindSecret)
	ss.Engine.POST("/v1/secret", sh.CreateSecret)
}

func (ss *SecretServer) setupMetrics() {

	prefix := "codersrank"
	ss.Monitor = metrics.NewMonitor(ss.Engine, prefix)

	ss.Engine.Use(mw.ResponseTime(ss.Monitor.RespTime))
	ss.Engine.Use(mw.CounterRequest(ss.Monitor.RequestTotal))
	ss.Engine.Use(mw.PercentileResponseTime(ss.Monitor.PercentilRespTime))

}

// SetupDB is charge to initialize the database connection
func (ss *SecretServer) setupDB() {

	user := os.Getenv("USER_DB")
	pass := os.Getenv("PASSWORD_DB")
	dbName := os.Getenv("NAME_DB")
	host := os.Getenv("HOST_DB")

	datasource := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", user, pass, host, dbName)

	db, err := gorm.Open("postgres", datasource)
	//db.LogMode(true)

	if err != nil {
		panic(err)
	}

	ss.DB = db
	im := initializer.NewInitialzerManager()
	im.Run(ss.DB)

}

// Run start the server
func (ss *SecretServer) Run() {
	ss.Engine.Run()
}

// Close all process in the server
func (ss *SecretServer) Close() {
	if ss.DB != nil {
		ss.DB.Close()
	}
}
