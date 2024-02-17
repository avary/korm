package korm

import (
	"context"
	"net/http"

	"github.com/kamalshkeir/ksmux"
)

func init() {
	ksmux.BeforeRenderHtml("korm-user", func(reqCtx context.Context, data *map[string]any) {
		const key ksmux.ContextKey = "user"
		user, ok := reqCtx.Value(key).(User)
		if ok {
			(*data)["IsAuthenticated"] = true
			(*data)["User"] = user
		} else {
			(*data)["IsAuthenticated"] = false
			(*data)["User"] = nil
		}
	})
}

func initAdminUrlPatterns(r *ksmux.Router) {
	media_root := http.FileServer(http.Dir("./" + MediaDir))
	r.Get(`/`+MediaDir+`/*path`, func(c *ksmux.Context) {
		http.StripPrefix("/"+MediaDir+"/", media_root).ServeHTTP(c.ResponseWriter, c.Request)
	})
	r.Get("/mon/ping", func(c *ksmux.Context) { c.Status(200).Text("pong") })
	r.Get("/offline", OfflineView)
	r.Get("/manifest.webmanifest", ManifestView)
	r.Get("/sw.js", ServiceWorkerView)
	r.Get("/robots.txt", RobotsTxtView)
	adminGroup := r.Group(adminPathNameGroup)
	adminGroup.Get("/", Admin(IndexView))
	adminGroup.Get("/login", Auth(LoginView))
	adminGroup.Post("/login", Auth(LoginPOSTView))
	adminGroup.Get("/logout", LogoutView)
	adminGroup.Get("/logs", Admin(LogsView))
	adminGroup.Post("/delete/row", Admin(DeleteRowPost))
	adminGroup.Post("/update/row", Admin(UpdateRowPost))
	adminGroup.Post("/create/row", Admin(CreateModelView))
	adminGroup.Post("/drop/table", Admin(DropTablePost))
	adminGroup.Get("/table/:model", Admin(AllModelsGet))
	adminGroup.Post("/table/:model/search", Admin(AllModelsSearch))
	adminGroup.Get("/get/:model/:id", Admin(SingleModelGet))
	adminGroup.Get("/export/:table", Admin(ExportView))
	adminGroup.Get("/export/:table/csv", Admin(ExportCSVView))
	adminGroup.Post("/import", Admin(ImportView))
}
