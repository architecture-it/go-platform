package web

import (
	"net/http"
	"os"
	"os/signal"
	"time"
	"github.com/gin-gonic/gin"
	"context"
	"github.com/zsais/go-gin-prometheus"
	"strings"
	"github.com/gin-contrib/cors"
	"github.com/eandreani/go-platform/log"
)
//Server un server http basado en gin-gonic
type Server struct {
	r *gin.Engine
	config Config
}

//GetRouter devuelve el router de gin
func (s *Server) GetRouter() *gin.Engine {
	return s.r
}

//NewServer crea un server nuevo con la config indicada
func NewServer(cfg Config) *Server {
	return &Server{gin.Default(),cfg}

}

// AddMetrics agrega un endpoint /metrics con las metricas de Prometheus para los requests 
func (s *Server) AddMetrics() *ginprometheus.Prometheus {
	p := ginprometheus.NewPrometheus("gin") 

	//esta funcion es para que se contabilicen agrupadas las metricas en cada endpoint mas alla de como cambie el ultimo elemento del path (el nombre de la cola o del topic)
	p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		url := c.Request.URL.String()
		for _, p := range c.Params {
			if p.Key == "name" {
				url = strings.Replace(url, p.Value, ":name", 1)
				break
			}
			if p.Key == "topic" {
				url = strings.Replace(url, p.Value, ":topic", 1)
				break
			}
		}
		return url
	}
	p.Use(s.r)

	return p
}


//AddHealth agrega un endpoint /health. Si alguno de los status_endpoints devuelve DOWN entonces /health va a devolver 404,
//si todos devuelve UP, 200-OK
//AddHealth(web.HealthAlwaysUp) siempre devuelve UP
//AddHealth(health.NewRedisHealthChecker(redisHealthChecker.Config{}),
//			health.NewMySqlHealthChecker(mySqlHealthChecker.Config{}),
//			...func())
func (s *Server) AddHealth(fs ...func()Status) {
	
	s.r.GET("/health", func (c *gin.Context) {
		result := make([]Status,len(fs))
		statusCode := http.StatusOK
		for i,f := range fs {
			check:=f()
			result[i]=check
			if check.Result != UP {
				statusCode = http.StatusNotFound
			}
		}
		c.JSON(statusCode,result)
	})
}
//ListenAndServe inicia el server http y bloquea hasta SIGINT
func (s *Server) ListenAndServe() {
	
	srv := &http.Server{
		Addr:    ":" + s.config.Port,
		Handler: s.r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal.Printf("listen: %s\n", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt) //(1) que nos notifique en el chanel quit @ SIGINT
	signal.Notify(quit, os.Kill)
	<-quit //esto se queda bloqueado aca hasta que (1) no sucede.
	log.Info.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal.Printf("Shutting down server: %s", err)
	}
	log.Info.Println("Farewell")

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