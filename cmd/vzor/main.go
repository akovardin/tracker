package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/cron"
	"github.com/pocketbase/pocketbase/tools/template"
	"go.uber.org/fx"

	"gohome.4gophers.ru/kovardin/vzor/app/handlers"
	"gohome.4gophers.ru/kovardin/vzor/app/tasks"
	"gohome.4gophers.ru/kovardin/vzor/app/tasks/uploader"
	"gohome.4gophers.ru/kovardin/vzor/static"
)

func main() {
	fx.New(
		handlers.Module,
		tasks.Module,

		fx.Provide(pocketbase.New),
		fx.Provide(template.NewRegistry),
		fx.Invoke(
			routing,
		),
		fx.Invoke(
			task,
		),
	).Run()
}

func routing(
	app *pocketbase.PocketBase,
	lc fx.Lifecycle,
	registry *template.Registry,
	home *handlers.Home,
	conversions *handlers.Conversions,
) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/", home.Home)

		e.Router.GET("/:name/fire", conversions.Fire)

		e.Router.GET("/static/*", func(c echo.Context) error {
			p := c.PathParam("*")

			path, err := url.PathUnescape(p)
			if err != nil {
				return fmt.Errorf("failed to unescape path variable: %w", err)
			}

			err = c.FileFS(path, static.FS)
			if err != nil && errors.Is(err, echo.ErrNotFound) {
				return c.FileFS("index.html", static.FS)
			}

			return err
		})

		return nil

	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go app.Start()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}

func task(app *pocketbase.PocketBase, uploader *uploader.Task) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		scheduler := cron.New()

		scheduler.MustAdd("hello", "0 */1 * * *", func() {
			uploader.Do()
		})

		scheduler.Start()

		return nil
	})
}
