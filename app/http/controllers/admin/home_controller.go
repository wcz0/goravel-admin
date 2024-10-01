package admin

import (
	"goravel/app/tools"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/wcz0/gamis"
	"github.com/wcz0/gamis/renderers"
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
			h.frameworkInfo().Set("md", 5),
			gamis.Flex().Items([]any{
				h.pieChart(),
				h.cube(),
			}),
		}),
		gamis.Grid().Columns([]any{
			h.lineChart().Set("md", 8),
			gamis.Flex().ClassName("h-full").Items([]any{
				h.clock(),
				h.codeView(),
			}).Direction("column"),
		}),
	})
	return h.SuccessData(ctx, page)
}

func (h *HomeController) codeView() *renderers.Card {
	return gamis.Card().ClassName("h-full clear-card-mb rounded-md").Body([]any{
		gamis.Markdown().Options(map[string]any{"html": true, "breaks": true}).Value(`
### __The beginning of everything__

<br>
`),
	})
}

func (h *HomeController) clock() *renderers.Card {
	return gamis.Card().ClassName("h-full bg-blingbling mb-4").Header(map[string]string{"title": "Clock"}).Body([]any{
		gamis.Custom().Name("clock").Html(`
<div id="clock" class="text-4xl"></div><div id="clock-date" class="mt-5"></div>
`).OnMount(`
const clock = document.getElementById('clock');
const tick = () => {
	clock.innerHTML = (new Date()).toLocaleTimeString();
	requestAnimationFrame(tick);
};
tick();

const clockDate = document.getElementById('clock-date');
clockDate.innerHTML = (new Date()).toLocaleDateString();
`),
	})

}

func (h *HomeController) pieChart() *renderers.Card {
	return gamis.Card().ClassName("h-96").Body([]any{
		gamis.Chart().Height(350).Config(map[string]any{
			"backgroundColor": "",
			"tooltip": map[string]string{
				"trigger": "item",
			},
			"legend": map[string]string{
				"bottom": "0",
				"left":   "center",
			},
			"series": []map[string]any{
				{
					"name":              "Access From",
					"type":              "pie",
					"radius":            []string{"40%", "70%"},
					"avoidLabelOverlap": false,
					"itemStyle": map[string]any{
						"borderRadius": 10,
						"borderColor":  "#fff",
						"borderWidth":  2,
					},
					"label": map[string]any{"show": false, "position": "center"},
					"emphasis": map[string]any{
						"label": map[string]any{
							"show":       true,
							"fontSize":   "40",
							"fontWeight": "bold",
						},
					},
					"labelLine": map[string]any{"show": false},
					"data": []map[string]any{
						{"value": 1048, "name": "Search Engine"},
						{"value": 735, "name": "Direct"},
						{"value": 580, "name": "Email"},
						{"value": 484, "name": "Union Ads"},
						{"value": 300, "name": "Video Ads"},
					},
				},
			},
		}),
	})
}

func (h *HomeController) cube() *renderers.Card {
	return gamis.Card().ClassName("h-96 ml-4 w-8/12").Body([]any{
		gamis.Html().Html(`
		<style>
    .cube-box{ height: 300px; display: flex; align-items: center; justify-content: center; }
  .cube { width: 100px; height: 100px; position: relative; transform-style: preserve-3d; animation: rotate 10s linear infinite; }
  .cube:after {
    content: '';
    width: 100%;
    height: 100%;
    box-shadow: 0 0 50px rgba(0, 0, 0, 0.2);
    position: absolute;
    transform-origin: bottom;
    transform-style: preserve-3d;
    transform: rotateX(90deg) translateY(50px) translateZ(-50px);
    background-color: rgba(0, 0, 0, 0.1);
  }
  .cube div {
    background-color: rgba(64, 158, 255, 0.7);
    position: absolute;
    width: 100%;
    height: 100%;
    border: 1px solid rgb(27, 99, 170);
    box-shadow: 0 0 60px rgba(64, 158, 255, 0.7);
  }
  .cube div:nth-child(1) { transform: translateZ(-50px); animation: shade 10s -5s linear infinite; }
  .cube div:nth-child(2) { transform: translateZ(50px) rotateY(180deg); animation: shade 10s linear infinite; }
  .cube div:nth-child(3) { transform-origin: right; transform: translateZ(50px) rotateY(270deg); animation: shade 10s -2.5s linear infinite; }
  .cube div:nth-child(4) { transform-origin: left; transform: translateZ(50px) rotateY(90deg); animation: shade 10s -7.5s linear infinite; }
  .cube div:nth-child(5) { transform-origin: bottom; transform: translateZ(50px) rotateX(90deg); background-color: rgba(0, 0, 0, 0.7); }
  .cube div:nth-child(6) { transform-origin: top; transform: translateZ(50px) rotateX(270deg); }

  @keyframes rotate {
    0% { transform: rotateX(-15deg) rotateY(0deg); }
    100% { transform: rotateX(-15deg) rotateY(360deg); }
  }
  @keyframes shade { 50% { background-color: rgba(0, 0, 0, 0.7); } }
</style>
<div class="cube-box">
    <div class="cube">
        <div></div>
        <div></div>
        <div></div>
        <div></div>
        <div></div>
        <div></div>
    </div>
</div>
`),
	})
}

func (h *HomeController) lineChart() *renderers.Card {
	random1 := []int{50, 200, 100, 150, 180, 130, 160}
	random2 := []int{150, 100, 180, 130, 160, 50, 200}
	chart := gamis.Chart().Height(380).ClassName("h-96").Config(map[string]any{
		"backgroundColor": "",
		"title":           map[string]string{"text": "Users Behavior"},
		"tooltip":         map[string]string{"trigger": "axis"},
		"xAxis": []map[string]any{
			{"type": "category", "boundaryGap": false, "data": []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}},
		},
		"yAxis":  []map[string]string{{"type": "value"}},
		"grid":   map[string]string{"left": "7%", "right": "3%", "top": "60", "bottom": "30"},
		"legend": map[string][]string{"data": {"Visits", "Bounce Rate"}},
		"series": []map[string]any{
			{"name": "Visits", "data": random1, "type": "line", "areaStyle": []string{}, "smooth": true, "symbol": "none"},
			{"name": "Bounce Rate", "data": random2, "type": "line", "areaStyle": []string{}, "smooth": true, "symbol": "none"},
		},
	})
	return gamis.Card().ClassName("clear-card-mb").Body(chart)
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

func (h *HomeController) frameworkInfo() *renderers.Card {
	link := func(label any, link string) renderers.Action {
		return *gamis.Action().
			Level("link").
			ClassName("text-lg font-semibold").
			Label(label).
			Set("blank", true).
			ActionType("url").
			Link(link)
	}
	config := facades.Config()
	return gamis.Card().ClassName("h-96").Body(
		gamis.Wrapper().ClassName("h-full").Body(
			gamis.Flex().
				ClassName("h-full").
				Direction("column").
				Justify("center").
				AlignItems("center").
				Items([]any{
					gamis.Image().Src(tools.Url(config.GetString("app.logo"))),
					gamis.Wrapper().ClassName("text-3xl mt-9 font-bold").Body(config.GetString("admin.name")),
					gamis.Flex().ClassName("w-full mt-5").Justify("center").Items([]any{
						link("Github", "https://github.com/wcz0/goravel-admin"),
					}),
				}),
		),
	)
}

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
