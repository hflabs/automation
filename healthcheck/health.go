package healthcheck

import (
	"context"
	"fmt"
	"github.com/carlmjohnson/requests"
	"os"
	"time"
)

var url = os.Getenv("MONITORING_URL")
var serviceName = os.Getenv("SERVICE_NAME")

type Healthcheck struct {
	ServiceName string `json:"service_name"`
	Status      int    `json:"status"`
}

func HealthcheckSender() {
	for {
		_ = requests.
			URL(fmt.Sprintf("%s/%s", url, "healthcheck")).
			Post().
			BodyJSON(Healthcheck{
				ServiceName: serviceName,
				Status:      1,
			}).
			Fetch(context.Background())
		time.Sleep(time.Minute)
	}
}
