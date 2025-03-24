package handler

import (
	"fmt"
	"net/http"

	"github.com/bignyap/kafka-go/pkg/utils"
)

func (app *Application) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(
		fmt.Sprintf("Server running on port %s", utils.GetEnvString("APPLICATION_PORT", "8080"))),
	)
}
