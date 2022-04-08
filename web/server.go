package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v3"

	"github.com/architecture-it/go-platform/health"
	"github.com/architecture-it/go-platform/log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"

	"github.com/common-nighthawk/go-figure"
)

//Server un server http basado en gin-gonic
type Server struct {
	r      *gin.Engine
	config Config
}

//GetRouter devuelve el router de gin
func (s *Server) GetRouter() *gin.Engine {
	return s.r
}

//NewServer crea un server nuevo con la config indicada
func NewServer(cfg Config) *Server {
	return &Server{gin.Default(), cfg}
}

func serveJSONFromFile(c *gin.Context) {
	var d interface{}
	file, err := os.Open("docs/swagger.json")
	if err != nil {
		log.Logger.Error(err.Error())
	}
	defer file.Close()
	b, _ := ioutil.ReadAll(file)
	if err := json.Unmarshal(b, &d); err != nil {
		log.Logger.Error(err.Error())
	}
	c.JSON(200, &d)
}

func serveYAMLFromFile(c *gin.Context) {
	var d interface{}
	file, err := os.Open("docs/swagger.yaml")
	if err != nil {
		log.Logger.Error(err.Error())
	}
	defer file.Close()
	b, _ := ioutil.ReadAll(file)
	if err := yaml.Unmarshal(b, &d); err != nil {
		log.Logger.Error(err.Error())
	}
	c.YAML(200, &d)
}

func (s *Server) AddApiDocs() {
	s.r.GET("/openapi.json", func(c *gin.Context) {
		serveJSONFromFile(c)
	})

	s.r.GET("/openapi.yaml", func(c *gin.Context) {
		serveYAMLFromFile(c)
	})
}

// AddMetrics agrega un endpoint /metrics con las metricas de Prometheus para los requests
func (s *Server) AddMetrics(fs ...func() []string) *ginprometheus.Prometheus {
	p := ginprometheus.NewPrometheus("gin")

	//Esta funcion es para que se contabilicen agrupadas las metricas en cada endpoint mas alla de como cambie el ultimo elemento del path
	if fs != nil {
		p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
			url := c.Request.URL.String()
			for _, p := range c.Params {
				for _, parameter := range fs[0]() {
					if p.Key == parameter {
						url = strings.Replace(url, p.Value, ":"+parameter, 1)
						break
					}
				}
			}
			return url
		}
	}
	p.Use(s.r)
	return p
}

//AddHealth agrega un endpoint /health
//Si alguno de los status_endpoints devuelve DOWN entonces /health va a devolver 404, si todos devuelve UP, 200-OK
//	AddHealth(web.HealthAlwaysUp) siempre devuelve UP
//	AddHealth(health.NewRedisHealthChecker(redisHealthChecker.Config{}),
//			health.NewMySqlHealthChecker(mySqlHealthChecker.Config{}),
//			...func())

func (s *Server) AddHealth(fs ...func(k ...string) health.Checker) {
	s.r.GET("/health", func(c *gin.Context) {
		generalHealth := health.HealthAlwaysUp()
		result := make(map[string]interface{})
		statusCode := http.StatusOK
		for _, f := range fs {
			fmt.Println(fs)
			check := f()
			result[check.Name] = check.Health
			if check.Health.Status.Code != health.UP {
				generalHealth.Status.Code = health.DOWN
			}
		}
		generalHealth.Details = result
		c.JSON(statusCode, generalHealth)
	})
}

//ListenAndServe inicia el server http y bloquea hasta SIGINT
func (s *Server) ListenAndServe() {
	myFigure := figure.NewFigure("go-platform", "", true)
	myFigure.Print()

	srv := &http.Server{
		Addr:    ":" + s.config.Port,
		Handler: s.r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.SugarLogger.Errorf("listen: %s\n", err.Error())
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt) //(1) que nos notifique en el channel quit @SIGINT
	signal.Notify(quit, os.Kill)
	<-quit //Esto se queda bloqueado aca hasta que (1) no sucede
	log.Logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.SugarLogger.Errorf("Shutting down server: %s", err.Error())
	}
	log.Logger.Info("Farewell")

}

//AddCorsAllOrigins es autoexplicativo
func (s *Server) AddCorsAllOrigins() {
	s.r.Use(cors.Default())
	//see: https://github.com/gin-contrib/cors

	/*s.r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))*/
}
