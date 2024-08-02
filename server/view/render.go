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

func RenderPage(ctx fiber.Ctx, component templ.Component) error {
	hxRequest := ctx.GetReqHeaders()["Hx-Request"]

	if len(hxRequest) == 0 {

		return Render(
			ctx,
			layout.Html(
				layout.HtmlProps{
					MainId:      components.MainId,
					CurrentPath: ctx.Path(),
				},
				component,
			),
		)
	}

	return Render(ctx, component)
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

func Fragment(fragmentId string, handler func(fiber.Ctx) (templ.Component, error)) fiber.Handler {
	return html(
		skip.New(
			func(c fiber.Ctx) error {

				component, err := handler(c)

				if err != nil {
					return err
				}

				return RenderPage(c, component)

			},
			func(ctx fiber.Ctx) bool {
				hxTarget := ctx.GetReqHeaders()["Hx-Target"]

				return len(hxTarget) == 0 || hxTarget[0] != fragmentId
			},
		),
	)
}

func Page(handler func(c fiber.Ctx) (page templ.Component, title string, err error)) fiber.Handler {
	return html(
		skip.New(
			func(c fiber.Ctx) error {
				page, title, err := handler(c)

				if err != nil {
					return err
				}

				return RenderPage(c, layout.Layout(title, page))
			},
			func(ctx fiber.Ctx) bool {
				hxTarget := ctx.GetReqHeaders()["Hx-Target"]

				return len(hxTarget) != 0 && hxTarget[0] != components.MainId
			},
		),
	)
}

func accepts(contentType string) func(fiber.Ctx) bool {
	return func(ctx fiber.Ctx) bool {
		return ctx.Accepts("text/html", "application/json") != contentType
	}
}

func html(handler fiber.Handler) fiber.Handler {
	return skip.New(
		handler,
		accepts("text/html"),
	)
}
