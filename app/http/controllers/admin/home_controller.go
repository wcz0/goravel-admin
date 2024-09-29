package admin

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/wcz0/gamis"
)

type HomeController struct {
	*Controller[string]
}

func NewHomeController() *HomeController {
	return &HomeController{
		Controller: NewAdminController[string](""),
	}
}

// GET /resource
func (h *HomeController) Index(ctx http.Context) http.Response {
	page := h.Controller.BasePage().Css(h.css(ctx)).Body([]any{
		gamis.Grid().ClassName("mb-1").Columns([]any{
			// h.frameworkInfo(ctx).Set("md", 5),
			// gamis.Flex().Items([]any{
			// )
		}),
	})
	return h.SuccessData(ctx, page)
}

// GET /resource/{id}
func (h *HomeController) Show(ctx http.Context) http.Response {
	return h.Success(ctx)
}

// POST /resource
func (h *HomeController) Store(ctx http.Context) http.Response {
	return h.Success(ctx)
}

// PUT /resource/{id}
func (h *HomeController) Update(ctx http.Context) http.Response {
	return h.Success(ctx)
}

// DELETE /resource/{id}
func (h *HomeController) Destroy(ctx http.Context) http.Response {
	return h.Success(ctx)
}

// func (h *HomeController) frameworkInfo(ctx http.Context) renderers.Card {
// 	link = func(label any, link string) renderers.Action {
// 		return *gamis.Action().
// 			Level("link").
// 			ClassName("text-lg font-semibold").
// 			Label(label).
// 			Set("blank", true).
// 			ActionType("url").
// 			Link(link)
// 	}
// 	config := facades.Config()
// 	return gamis.Card().ClassName("h-96").Body(
// 		gamis.Wrapper().ClassName("h-full").Body(
// 			gamis.Flex().ClassName("h-full").Direction("column").
// 				Justify("center").
// 				AlignItems("center").
// 				Items([]any{
// 					gamis.Image().Src(tools.Url(config.GetString("app.logo"))),
// 					gamis.Wrapper().ClassName("text-3xl mt-9 font-bold").Body(config.GetString("admin.name")),
// 					gamis.Flex().ClassName("w-full mt-5").Justify("center").Items([]any{
// 						link("Github", "https://github.com/wcz0/goravel-admin"),
// 						// link("Amis")
// 					}),
// }),
// 		),
// 	)
// }

func (h *HomeController) css(ctx http.Context) map[string]any {
	return map[string]any{
		".clear-card-mb": []string{"0 !important"},
		".cxd-Image": map[string]string{
			"border": "0",
		},
		".bg-blingbling": map[string]string{
			"color":             "#fff",
			"background":        "linear-gradient(to bottom right, #2C3E50, #FD746C, #FF8235, #ffff1c, #92fe9d, #00c9ff, #a044ff, #e73827 )",
			"background-repeat": "no-repeat",
			"background-size":   "100% 100%",
			"animation":         "gradient 60s ease infinite",
		},
		"@keyframes gradient": []string{
			"0%{background-position:0% 0%} 50%{background-position:100% 100%} 100%{background-position:0% 0%}",
		},
		".bg-blingbling .cxd-Card-title": map[string]string{
			"color": "#fff",
		},
	}
}
