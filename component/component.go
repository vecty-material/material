package component // import "agamigo.io/material/component"

import (
	"errors"

	"agamigo.io/gojs"
	"github.com/gopherjs/gopherjs/js"
)

// Componenter is the base interface for every material component
// implementation.
type Componenter interface {
	// SetComponent should replace a component implementation's Component with
	// the provided component.
	SetComponent(c *Component)

	// GetComponent should return a pointer to the component implementation's
	// underlying Component. Implementors that embed a *Component directly
	// get this for free.
	GetComponent() (c *Component)
}

// AfterStarter is implemented by components that need further setup ran
// after their underlying MDC foundation has been initialized.
type AfterStarter interface {
	AfterStart() error
}

type ComponentTyper interface {
	ComponentType() ComponentType
}

// MDCClasser is an interface that allows component users to specify the MDC
// class object that will be used to create/initialize the component.
type MDCClasser interface {
	MDCClass() *js.Object
}

// ComponentStatus holds a component's lifecycle status.
type ComponentStatus int

const (
	// An Uninitialized component has not been associated with the MDC library
	// yet. This package does not provide a way to access an Uninitialized
	// component.
	Uninitialized ComponentStatus = iota

	// A Stopped component has been associated with a JS Object constructed from
	// a MDC class. New() returns a Stopped component, and Stop() will stop a
	// Running component.
	Stopped

	// A Running component has had its underlying MDC init() method called,
	// which attaches the component to a specific HTMLElement in the DOM. It is
	// ready to be used.
	Running
)

// Component is the base material component type. Types that embed Component and
// implement Componenter can use the component.Start and component.Stop
// functions.
type Component struct {
	mdc    *js.Object
	status ComponentStatus
}

// String returns the Component's StatusType as text.
func (c *Component) String() string {
	if c == nil || c.status == Uninitialized {
		return Uninitialized.String()
	}
	return c.Status().String()
}

// String returns the string version of a StatusType.
func (s ComponentStatus) String() string {
	switch s {
	case Stopped:
		return "stopped"
	case Running:
		return "running"
	}
	return "uninitialized"
}

// Status returns the component's StatusType. For the string version use
// Status().String().
func (c *Component) Status() ComponentStatus {
	return c.status
}

// Start takes a component implementation (c) and initializes it with an
// HTMLElement (rootElem). Upon success the component's status will be Running,
// and err will be nil.  If err is non-nil, it will contain any error thrown
// while calling the underlying MDC object's init() method, and the component's
// status will remain Stopped.
//
// Finding The MDC Library
//
// There are two ways Start knows of to find the MDC class needed to start a
// component. By default it uses values provided by the components in this
// project via the ComponentType method. This default works in the general case
// that the all-in-one MDC library is available under the global var "mdc".
//
// The second case, MDCClasser, is needed if the MDC code you need is elsewhere,
// for example if you are using individual MDC component "@material/checkbox"
// libraries instead of the all-in-one distribution. Implement the MDCClasser to
// provide Start with the exact object for the MDC component class.
//
// Implementing A Component
//
// If you are writing a component implementation the documentation for the
// Componenter interface provides useful information.
//
// If you need to perform additional work on the Component after initialization,
// read the AfterStarter interface documentation. If AfterStart returns a
// non-nill error then Stop will be called on the component.
//
// See: https://material.io/components/web/docs/framework-integration/
func Start(c Componenter, rootElem *js.Object) (err error) {
	defer gojs.CatchException(&err)

	switch {
	case rootElem == nil, rootElem == js.Undefined:
		return errors.New("rootElem is nil.")
	case c.GetComponent() == nil:
		c.SetComponent(&Component{})
	case c.GetComponent().status == Running:
		return errors.New("Component already started.")
	}

	// We create a new instance of the MDC component if c is Stopped or
	// Uninitialized.
	var newMDCClassObj *js.Object
	switch co := c.(type) {
	case MDCClasser:
		newMDCClassObj = co.MDCClass()
	case ComponentTyper:
		CCaseName := co.ComponentType().MDCCamelCaseName
		ClassName := co.ComponentType().MDCClassName
		if CCaseName == "" || ClassName == "" {
			return errors.New("Empty string in ComponentType")
		}
		mdcObject := js.Global.Get("mdc")
		newMDCClassObj = mdcObject.Get(CCaseName).Get(ClassName)
	default:
		return errors.New("The provided component does not implement " +
			"component.ComponentTyper or component.MDCClasser.")
	}

	// Create a new MDC component instance tied to rootElem
	newMDCObj := newMDCClassObj.New(rootElem)
	c.GetComponent().mdc = newMDCObj
	c.GetComponent().status = Running

	switch co := c.(type) {
	case AfterStarter:
		err = co.AfterStart()
		if err != nil {
			Stop(c)
			return err
		}
	}

	return err
}

// Stop stops a Running component, removing its association with an HTMLElement
// and cleaning up event listeners, etc. It changes the component's status to
// Stopped.
func Stop(c Componenter) (err error) {
	defer gojs.CatchException(&err)

	if c.GetComponent() == nil {
		return errors.New("GetComponent() returned nil.")
	}

	switch c.GetComponent().status {
	case Stopped:
		return errors.New("Component already stopped")
	case Uninitialized:
		return errors.New("Component is uninitialized")
	}
	c.GetComponent().mdc.Call("destroy")
	c.SetComponent(&Component{status: Stopped})
	return err
}

// GetComponent implements the Componenter interface. Component implementations
// can use this method as-is when embedding an exposed component.Component.
func (c *Component) GetComponent() *Component {
	return c
}

// ComponentType implements the Componenter interface. This should be shadowed
// by a component implementation.
func (c *Component) ComponentType() ComponentType {
	return ComponentType{}
}

// GetObject returns the component's MDC JavaScript object.
func (c *Component) GetObject() *js.Object {
	return c.mdc
}
