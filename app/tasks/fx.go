package tasks

import (
	"go.uber.org/fx"

	"gohome.4gophers.ru/kovardin/vzor/app/tasks/uploader"
)

var Module = fx.Module(
	"tasks",
	fx.Provide(uploader.New),
)
