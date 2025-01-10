package healthcheck

import (
	"cmp"
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
	PeriodInMin int    `json:"period"`
}

func HealthcheckSender(minutes int) {
	for {
		minutes = cmp.Or(minutes, 1)
		_ = requests.
			URL(fmt.Sprintf("%s/%s", url, "healthcheck")).
			Post().
			BodyJSON(Healthcheck{
				ServiceName: serviceName,
				Status:      1,
				PeriodInMin: minutes,
			}).
			Fetch(context.Background())
		time.Sleep(time.Duration(minutes) * time.Minute)
	}
}
