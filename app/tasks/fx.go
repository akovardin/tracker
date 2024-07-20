package tasks

import (
	"go.uber.org/fx"

	"gohome.4gophers.ru/kovardin/tracker/app/tasks/uploader"
)

var Module = fx.Module(
	"tasks",
	fx.Provide(uploader.New),
)
