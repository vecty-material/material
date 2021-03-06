package ul

import (
	"syscall/js"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
	"github.com/vecty-material/material/base"
)

type nativeInputer interface {
	NativeInput() (*vecty.HTML, string)
}

// L is a vecty-material list component.
type L struct {
	*base.MDC
	vecty.Core
	Root           vecty.MarkupOrChild
	Items          []vecty.ComponentOrHTML
	Dense          bool
	Avatar         bool
	NonInteractive bool
	OnClick        func(thisL *L, thisI *Item, e *vecty.Event)
	GroupSubheader string
	twoLine        bool
}

// Item is a vecty-material list-item component.
type Item struct {
	*base.MDC
	vecty.Core
	Root      vecty.MarkupOrChild
	Primary   vecty.ComponentOrHTML
	Secondary vecty.ComponentOrHTML
	Graphic   vecty.ComponentOrHTML
	Meta      vecty.ComponentOrHTML
	Selected  bool
	Activated bool
	OnClick   func(i *Item, e *vecty.Event)
	Href      string
}

// Group is a vecty-material list-group component.
type Group struct {
	*base.MDC
	vecty.Core
	Root  vecty.MarkupOrChild
	Lists []vecty.ComponentOrHTML
}

type divider struct {
	vecty.Core
}

// Render implements the vecty.Component interface.
func (c *L) Render() vecty.ComponentOrHTML {
	rootMarkup := base.MarkupOnly(c.Root)
	if c.Root != nil && rootMarkup == nil {
		// User supplied root element.
		return elem.UnorderedList(c.Root)
	}

	items := make([]vecty.MarkupOrChild, len(c.Items))
	c.twoLine = false
	for i, li := range c.Items {
		switch t := li.(type) {
		case *Item:
			if t.Secondary != nil {
				c.twoLine = true
			}
			items[i] = t
		case *vecty.HTML:
			items[i] = base.RenderStoredChild(t)
		default:
			items[i] = li
		}
	}

	root := elem.UnorderedList(items...)
	vecty.Markup(
		c,
		vecty.MarkupIf(rootMarkup != nil, *rootMarkup),
	).Apply(root)
	return root
}

func (c *L) Apply(h *vecty.HTML) {
	switch {
	case c.MDC == nil:
		c.MDC = &base.MDC{}
	}

	vecty.Markup(
		vecty.Class("mdc-list"),
		vecty.MarkupIf(c.twoLine,
			vecty.Class("mdc-list--two-line")),
		vecty.MarkupIf(c.Dense,
			vecty.Class("mdc-list--dense")),
		vecty.MarkupIf(c.Avatar,
			vecty.Class("mdc-list--avatar-list")),
		vecty.MarkupIf(c.NonInteractive,
			vecty.Class("mdc-list--non-interactive")),
	).Apply(h)
	c.MDC.RootElement = h
}

// Render implements the vecty.Component interface.
func (c *Item) Render() vecty.ComponentOrHTML {
	tag := "li"
	if c.Href != "" {
		tag = "a"
	}

	rootMarkup := base.MarkupOnly(c.Root)
	if c.Root != nil && rootMarkup == nil {
		// User supplied root element.
		return vecty.Tag(tag, c.Root)
	}

	graphic := setupGraphicOrMeta(c.Graphic)
	if graphic != nil {
		if g, ok := graphic.(*vecty.HTML); ok {
			vecty.Class("mdc-list-item__graphic").Apply(g)
			vecty.Attribute("role", "presentation").Apply(g)
		}
	}

	meta := setupGraphicOrMeta(c.Meta)
	if meta != nil {
		if g, ok := meta.(*vecty.HTML); ok {
			vecty.Class("mdc-list-item__meta").Apply(g)
			vecty.Attribute("role", "presentation").Apply(g)
		}
	}

	var text vecty.ComponentOrHTML
	switch {
	case c.Secondary != nil:
		text = elem.Span(vecty.Markup(vecty.Class("mdc-list-item__text")),
			c.Primary,
			elem.Span(vecty.Markup(
				vecty.Class("mdc-list-item__secondary-text")),
				c.Secondary,
			))
	default:
		text = c.Primary
	}

	return vecty.Tag(tag,
		vecty.Markup(
			c,
			vecty.MarkupIf(rootMarkup != nil, *rootMarkup),
		),
		graphic,
		base.RenderStoredChild(text),
		meta,
	)
}

func (c *Item) Apply(h *vecty.HTML) {
	switch {
	case c.MDC == nil:
		c.MDC = &base.MDC{}
	}

	vecty.Markup(
		vecty.Class("mdc-list-item"),
		vecty.MarkupIf(c.Selected,
			vecty.Class("mdc-list-item--selected")),
		vecty.MarkupIf(c.Activated,
			vecty.Class("mdc-list-item--activated")),
		vecty.MarkupIf(c.OnClick != nil,
			event.Click(c.wrapOnClick()),
		),
		vecty.MarkupIf(c.Href != "", prop.Href(c.Href)),
	).Apply(h)
	c.MDC.RootElement = h
}

// Render implements the vecty.Component interface.
func (c *Group) Render() vecty.ComponentOrHTML {
	rootMarkup := base.MarkupOnly(c.Root)
	if c.Root != nil && rootMarkup == nil {
		// User supplied root element.
		return elem.Div(c.Root)
	}

	return elem.Div(
		vecty.Markup(
			c,
			vecty.MarkupIf(rootMarkup != nil, *rootMarkup),
		),
		c.listList(),
	)
}

func (c *Group) Apply(h *vecty.HTML) {
	switch {
	case c.MDC == nil:
		c.MDC = &base.MDC{}
	}

	vecty.Markup(
		vecty.Class("mdc-list-group"),
	).Apply(h)
	c.MDC.RootElement = h
}

func ListDivider() vecty.ComponentOrHTML {
	d := elem.HorizontalRule(
		vecty.Markup(
			vecty.Class("mdc-list-divider"),
		),
	)
	return base.RenderStoredChild(d)
}

func ListDividerInset() vecty.ComponentOrHTML {
	d := elem.HorizontalRule(
		vecty.Markup(
			vecty.Class("mdc-list-divider"),
			vecty.Class("mdc-list-divider--inset"),
		),
	)
	return base.RenderStoredChild(d)
}

func ItemDivider() vecty.ComponentOrHTML {
	d := elem.ListItem(
		vecty.Markup(
			vecty.Class("mdc-list-divider"),
			vecty.Attribute("role", "separator"),
		),
	)
	return base.RenderStoredChild(d)
}

func ItemDividerInset() vecty.ComponentOrHTML {
	d := elem.ListItem(
		vecty.Markup(
			vecty.Class("mdc-list-divider"),
			vecty.Class("mdc-list-divider--inset"),
			vecty.Attribute("role", "separator"),
		),
	)
	return base.RenderStoredChild(d)
}

func (c *Group) listList() vecty.List {
	lists := make(vecty.List, len(c.Lists)*2)
	for _, cList := range c.Lists {
		if list, ok := cList.(*L); ok {
			if list.GroupSubheader != "" {
				lists = append(lists,
					elem.Heading3(vecty.Markup(
						vecty.Class("mdc-list-group__subheader")),
						vecty.Text(list.GroupSubheader)),
				)
			}
			lists = append(lists, list)
		}
	}
	return lists
}

func (c *L) wrapOnClick() func(i *Item, e *vecty.Event) {
	return func(i *Item, e *vecty.Event) {
		c.OnClick(c, i, e)
	}
}

func (c *Item) wrapOnClick() func(e *vecty.Event) {
	return func(e *vecty.Event) {
		c.OnClick(c, e)
	}
}

func setupGraphicOrMeta(c vecty.ComponentOrHTML) vecty.ComponentOrHTML {
	var graphic vecty.ComponentOrHTML
	if c != nil {
		graphic = c
		if js.InternalObject(graphic).Get("tag").String() != "img" {
			graphic = elem.Span(graphic)
		}
	}
	return graphic
}
