package ctx

import (
	"github.com/gin-gonic/gin"
	"github.com/philipbo/gin_project_tmpl/db"
)

const (
	CTX_DS = "CTX_DS"
)

func MustGetDs(c *gin.Context) *db.Ds {
	return c.MustGet(CTX_DS).(*db.Ds)
}
