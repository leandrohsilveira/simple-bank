package view

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gofiber/fiber/v3/middleware/skip"
	"github.com/leandrohsilveira/simple-bank/server/components"
	"github.com/leandrohsilveira/simple-bank/server/layout"
)

func Handler(component templ.Component) fiber.Handler {
	return adaptor.HTTPHandler(templ.Handler(component))
}

func Render(ctx fiber.Ctx, component templ.Component) error {
	ctx.Set("Content-Type", "text/html")
	return component.Render(ctx.UserContext(), ctx.Response().BodyWriter())
}

func RenderPage(ctx fiber.Ctx, component templ.Component, title string) error {
	hxRequest := ctx.GetReqHeaders()["Hx-Request"]

	page := layout.Layout(
		title,
		component,
	)

	if len(hxRequest) == 0 {

		return Render(
			ctx,
			layout.Html(
				layout.HtmlProps{
					MainId:      components.MainId,
					CurrentPath: ctx.Path(),
				},
				page,
			),
		)
	}

	return Render(ctx, page)
}

func Json(handler fiber.Handler) fiber.Handler {
	return skip.New(
		handler,
		accepts("application/json"),
	)
}

func Htmx(handler fiber.Handler) fiber.Handler {
	return html(
		skip.New(
			handler,
			func(ctx fiber.Ctx) bool {
				hxRequest := ctx.GetReqHeaders()["Hx-Request"]
				return len(hxRequest) == 0
			},
		),
	)
}

func Fragment(fragmentId string, handler fiber.Handler) fiber.Handler {
	return html(
		skip.New(
			handler,
			func(ctx fiber.Ctx) bool {
				hxTarget := ctx.GetReqHeaders()["Hx-Target"]

				return len(hxTarget) == 0 || hxTarget[0] != fragmentId
			},
		),
	)
}

func Page(handler fiber.Handler) fiber.Handler {
	return html(
		skip.New(
			handler,
			func(ctx fiber.Ctx) bool {
				hxTarget := ctx.GetReqHeaders()["Hx-Target"]

				return len(hxTarget) != 0 && hxTarget[0] != components.MainId
			},
		),
	)
}

func Accept(ctx fiber.Ctx) string {
	return ctx.Accepts("text/html", "application/json")
}

func accepts(contentType string) func(fiber.Ctx) bool {
	return func(ctx fiber.Ctx) bool {
		return Accept(ctx) != contentType
	}
}

func html(handler fiber.Handler) fiber.Handler {
	return skip.New(
		handler,
		accepts("text/html"),
	)
}
