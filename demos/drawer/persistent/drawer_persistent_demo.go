package main

import (
	"syscall/js"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
	"github.com/vecty-material/material/button"
	"github.com/vecty-material/material/demos/common"
	dcommon "github.com/vecty-material/material/demos/drawer/common"
	"github.com/vecty-material/material/drawer"
	"github.com/vecty-material/material/formfield"
	"github.com/vecty-material/material/radio"
)

type drawerDemoView struct {
	vecty.Core
	body   *vecty.HTML
	drawer *drawer.D
}

func main() {
	vecty.RenderBody(&drawerDemoView{})
}

func (c *drawerDemoView) Render() vecty.ComponentOrHTML {
	c.drawer = dcommon.NewDemoDrawer(drawer.Persistent)
	c.body = elem.Body(
		vecty.Markup(
			vecty.Class("mdc-typography"),
			vecty.Class("demo-body"),
		),
		c.drawer,
		elem.Div(
			vecty.Markup(
				vecty.Class("demo-content"),
			),
			&common.ToolbarHeader{
				Title:      "Persistent Drawer",
				Navigation: common.NavMenu,
				MenuHandler: func(e *vecty.Event) {
					c.drawer.Open = !c.drawer.Open
					vecty.Rerender(c.drawer)
				},
				NoFixed: true,
			},
			elem.Main(vecty.Markup(vecty.Class("demo-main")),
				elem.Heading1(
					vecty.Markup(vecty.Class("mdc-typography--display1")),
					vecty.Text("Persistent Drawer")),
				elem.Paragraph(
					vecty.Markup(vecty.Class("mdc-typography--body1")),
					vecty.Text("Click the menu icon above to open and "+
						"close the drawer."),
				),
				elem.Div(vecty.Markup(prop.ID("demo-radio-buttons")),
					&formfield.FF{
						Label: "Default",
						Input: vecty.List{
							&radio.R{
								Root: vecty.Markup(
									prop.ID("theme-radio-default"),
								),
								Name:    "theme",
								Checked: true,
								OnChange: func(thisR *radio.R,
									e *vecty.Event) {
									d := js.Global().Get("window").Get("document")
									dd := d.Call("querySelector", ".demo-drawer")
									dd.Get("classList").Call("remove",
										"demo-drawer--accessible")
									dd.Get("classList").Call("remove",
										"demo-drawer--custom")
								},
							},
						},
					},
					&formfield.FF{
						Label: "Custom Theme",
						Input: vecty.List{
							&radio.R{
								Root: vecty.Markup(
									prop.ID("theme-radio-custom"),
								),
								Name: "theme",
								OnChange: func(thisR *radio.R,
									e *vecty.Event) {
									d := js.Global().Get("window").Get("document")
									dd := d.Call("querySelector", ".demo-drawer")
									dd.Get("classList").Call("remove",
										"demo-drawer--accessible")
									dd.Get("classList").Call("add",
										"demo-drawer--custom")
								},
							},
						},
					},
					&formfield.FF{
						Label: "Accessible Theme",
						Input: vecty.List{
							&radio.R{
								Root: vecty.Markup(
									prop.ID("theme-radio-accessible"),
								),
								Name: "theme",
								OnChange: func(thisR *radio.R,
									e *vecty.Event) {
									d := js.Global().Get("window").Get("document")
									dd := d.Call("querySelector", ".demo-drawer")
									dd.Get("classList").Call("remove",
										"demo-drawer--custom")
									dd.Get("classList").Call("add",
										"demo-drawer--accessible")
								},
							},
						},
					},
				),
				elem.Div(vecty.Markup(vecty.Class("extra-content-wrapper")),
					&button.B{
						Root: vecty.Markup(vecty.Class(
							"demo-toolbar-example-heading__rtl-toggle-button"),
						),
						Label:    vecty.Text("Toggle RTL"),
						Outlined: true,
						Dense:    true,
						OnClick: func(thisB *button.B,
							e *vecty.Event) {
							b := c.body.Node()
							if b.Call("getAttribute",
								"dir").String() == "rtl" {
								b.Call("setAttribute", "dir", "ltr")
								return
							}
							b.Call("setAttribute", "dir", "rtl")
						},
					},
				),
			),
		),
	)
	return c.body
}
