package html5

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/the80srobot/html5/bindings"
	"github.com/the80srobot/html5/safe"
)

type AttributeNode struct {
	Name          string
	Value         Value
	RequiredTrust safe.TrustLevel
}

func Attribute(name string, value Value) *AttributeNode {
	return &AttributeNode{Name: name, Value: value}
}

func DataAttribute(name string, value Value, trust safe.TrustLevel) *AttributeNode {
	return &AttributeNode{Name: "data-" + name, Value: value, RequiredTrust: trust}
}

// Apply will insert the attribute into the node, which must be ElementNode.
func (a *AttributeNode) Apply(n Node) error {
	e, ok := n.(*ElementNode)
	if !ok {
		return fmt.Errorf("attributes must be applied to elements, got node %v", n)
	}
	e.Attributes = append(e.Attributes, *a)
	return nil
}

func appendAttribute(tc *templateCompiler, a *AttributeNode) error {
	if _, err := fmt.Fprintf(tc, " %s=\"", a.Name); err != nil {
		return err
	}

	// Different attributes require different levels of trust (e.g. href
	// contains URLs).
	reqTrust, ok := requiredTrustPerAttribute[a.Name]
	if !ok {
		if strings.HasPrefix(a.Name, "data-") {
			reqTrust = safe.Default
		} else {
			reqTrust = safe.FullyTrusted
		}
	}
	reqTrust = safe.Max(reqTrust, a.RequiredTrust)

	switch v := a.Value.(type) {
	case safe.String:
		s, err := safe.Check(v, reqTrust)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(tc, "%s\"", s)
		return err
	case bindings.Var:
		tc.appendVar(v, reqTrust)
		_, err := fmt.Fprint(tc, "\"")
		return err
	default:
		return fmt.Errorf("value must be safe.String or *bindings.Var, %v (%v) is neither", v, reflect.TypeOf(v))
	}
}

// Lists the required trust level for the content of known HTML attributes. If
// an attribute is not on this list, then assume FullyTrusted is required.
//
// Current spec: https://html.spec.whatwg.org/multipage/indices.html#attributes-3
var requiredTrustPerAttribute = map[string]safe.TrustLevel{
	"accept":          safe.AttributeSafe,
	"accept-charset":  safe.FullyTrusted,
	"action":          safe.URLSafe,
	"alt":             safe.AttributeSafe,
	"archive":         safe.URLSafe,
	"async":           safe.FullyTrusted,
	"autocomplete":    safe.AttributeSafe,
	"autofocus":       safe.AttributeSafe,
	"autoplay":        safe.AttributeSafe,
	"background":      safe.URLSafe,
	"border":          safe.AttributeSafe,
	"checked":         safe.AttributeSafe,
	"cite":            safe.URLSafe,
	"challenge":       safe.FullyTrusted,
	"charset":         safe.FullyTrusted,
	"class":           safe.AttributeSafe,
	"classid":         safe.URLSafe,
	"codebase":        safe.URLSafe,
	"cols":            safe.AttributeSafe,
	"colspan":         safe.AttributeSafe,
	"content":         safe.FullyTrusted,
	"contenteditable": safe.AttributeSafe,
	"contextmenu":     safe.AttributeSafe,
	"controls":        safe.AttributeSafe,
	"coords":          safe.AttributeSafe,
	"crossorigin":     safe.FullyTrusted,
	"data":            safe.URLSafe,
	"datetime":        safe.AttributeSafe,
	"default":         safe.AttributeSafe,
	"defer":           safe.FullyTrusted,
	"dir":             safe.AttributeSafe,
	"dirname":         safe.AttributeSafe,
	"disabled":        safe.AttributeSafe,
	"draggable":       safe.AttributeSafe,
	"dropzone":        safe.AttributeSafe,
	"enctype":         safe.FullyTrusted,
	"for":             safe.AttributeSafe,
	"form":            safe.FullyTrusted,
	"formaction":      safe.URLSafe,
	"formenctype":     safe.FullyTrusted,
	"formmethod":      safe.FullyTrusted,
	"formnovalidate":  safe.FullyTrusted,
	"formtarget":      safe.AttributeSafe,
	"headers":         safe.AttributeSafe,
	"height":          safe.AttributeSafe,
	"hidden":          safe.AttributeSafe,
	"high":            safe.AttributeSafe,
	"href":            safe.URLSafe,
	"hreflang":        safe.AttributeSafe,
	"http-equiv":      safe.FullyTrusted,
	"icon":            safe.URLSafe,
	"id":              safe.AttributeSafe,
	"ismap":           safe.AttributeSafe,
	"keytype":         safe.FullyTrusted,
	"kind":            safe.AttributeSafe,
	"label":           safe.AttributeSafe,
	"lang":            safe.AttributeSafe,
	"language":        safe.FullyTrusted,
	"list":            safe.AttributeSafe,
	"longdesc":        safe.URLSafe,
	"loop":            safe.AttributeSafe,
	"low":             safe.AttributeSafe,
	"manifest":        safe.URLSafe,
	"max":             safe.AttributeSafe,
	"maxlength":       safe.AttributeSafe,
	"media":           safe.AttributeSafe,
	"mediagroup":      safe.AttributeSafe,
	"method":          safe.FullyTrusted,
	"min":             safe.AttributeSafe,
	"multiple":        safe.AttributeSafe,
	"name":            safe.AttributeSafe,
	"novalidate":      safe.FullyTrusted,
	"open":            safe.AttributeSafe,
	"optimum":         safe.AttributeSafe,
	"pattern":         safe.FullyTrusted,
	"placeholder":     safe.AttributeSafe,
	"poster":          safe.URLSafe,
	"profile":         safe.URLSafe,
	"preload":         safe.AttributeSafe,
	"pubdate":         safe.AttributeSafe,
	"radiogroup":      safe.AttributeSafe,
	"readonly":        safe.AttributeSafe,
	"rel":             safe.FullyTrusted,
	"required":        safe.AttributeSafe,
	"reversed":        safe.AttributeSafe,
	"rows":            safe.AttributeSafe,
	"rowspan":         safe.AttributeSafe,
	"sandbox":         safe.FullyTrusted,
	"spellcheck":      safe.AttributeSafe,
	"scope":           safe.AttributeSafe,
	"scoped":          safe.AttributeSafe,
	"seamless":        safe.AttributeSafe,
	"selected":        safe.AttributeSafe,
	"shape":           safe.AttributeSafe,
	"size":            safe.AttributeSafe,
	"sizes":           safe.AttributeSafe,
	"span":            safe.AttributeSafe,
	"src":             safe.URLSafe,
	"srcdoc":          safe.HTMLSafe,
	"srclang":         safe.AttributeSafe,
	"srcset":          safe.FullyTrusted,
	"start":           safe.AttributeSafe,
	"step":            safe.AttributeSafe,
	"style":           safe.FullyTrusted,
	"tabindex":        safe.AttributeSafe,
	"target":          safe.AttributeSafe,
	"title":           safe.AttributeSafe,
	"type":            safe.FullyTrusted,
	"usemap":          safe.URLSafe,
	"value":           safe.FullyTrusted,
	"width":           safe.AttributeSafe,
	"wrap":            safe.AttributeSafe,
	"xmlns":           safe.URLSafe,
}
