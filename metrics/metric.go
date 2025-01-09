package metrics

import (
	"fmt"
	"github.com/carlmjohnson/requests"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

var url = os.Getenv("MONITORING_URL")
var serviceName = os.Getenv("SERVICE_NAME")

type Metric struct {
	ServiceName string  `json:"service_name"`
	Endpoint    string  `json:"endpoint"`
	Method      string  `json:"method"`
	Status      int     `json:"status"`
	Duration    float64 `json:"duration"`
}

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()
		_ = requests.
			URL(fmt.Sprintf("%s/%s", url, "metrics")).Post().
			BodyJSON(Metric{
				ServiceName: serviceName,
				Endpoint:    c.FullPath(),
				Method:      c.Request.Method,
				Status:      c.Writer.Status(),
				Duration:    duration,
			}).
			Fetch(c)
	}
}
