package html

import "fmt"

type Attribute struct {
	Name  string
	Value SafeString
}

func (a *Attribute) Apply(n Node) error {
	e, ok := n.(*ElementNode)
	if !ok {
		return fmt.Errorf("attributes must be applied to elements, got node %v", n)
	}
	e.Attributes = append(e.Attributes, *a)
	return nil
}

// Lists the required trust level for the content of known HTML attributes. If
// an attribute is not on this list, then assume FullyTrusted is required.
var requiredTrustPerAttribute = map[string]StringTrust{
	"accept":          AttributeSafe,
	"accept-charset":  FullyTrusted,
	"action":          URLSafe,
	"alt":             AttributeSafe,
	"archive":         URLSafe,
	"async":           FullyTrusted,
	"autocomplete":    AttributeSafe,
	"autofocus":       AttributeSafe,
	"autoplay":        AttributeSafe,
	"background":      URLSafe,
	"border":          AttributeSafe,
	"checked":         AttributeSafe,
	"cite":            URLSafe,
	"challenge":       FullyTrusted,
	"charset":         FullyTrusted,
	"class":           AttributeSafe,
	"classid":         URLSafe,
	"codebase":        URLSafe,
	"cols":            AttributeSafe,
	"colspan":         AttributeSafe,
	"content":         FullyTrusted,
	"contenteditable": AttributeSafe,
	"contextmenu":     AttributeSafe,
	"controls":        AttributeSafe,
	"coords":          AttributeSafe,
	"crossorigin":     FullyTrusted,
	"data":            URLSafe,
	"datetime":        AttributeSafe,
	"default":         AttributeSafe,
	"defer":           FullyTrusted,
	"dir":             AttributeSafe,
	"dirname":         AttributeSafe,
	"disabled":        AttributeSafe,
	"draggable":       AttributeSafe,
	"dropzone":        AttributeSafe,
	"enctype":         FullyTrusted,
	"for":             AttributeSafe,
	"form":            FullyTrusted,
	"formaction":      URLSafe,
	"formenctype":     FullyTrusted,
	"formmethod":      FullyTrusted,
	"formnovalidate":  FullyTrusted,
	"formtarget":      AttributeSafe,
	"headers":         AttributeSafe,
	"height":          AttributeSafe,
	"hidden":          AttributeSafe,
	"high":            AttributeSafe,
	"href":            URLSafe,
	"hreflang":        AttributeSafe,
	"http-equiv":      FullyTrusted,
	"icon":            URLSafe,
	"id":              AttributeSafe,
	"ismap":           AttributeSafe,
	"keytype":         FullyTrusted,
	"kind":            AttributeSafe,
	"label":           AttributeSafe,
	"lang":            AttributeSafe,
	"language":        FullyTrusted,
	"list":            AttributeSafe,
	"longdesc":        URLSafe,
	"loop":            AttributeSafe,
	"low":             AttributeSafe,
	"manifest":        URLSafe,
	"max":             AttributeSafe,
	"maxlength":       AttributeSafe,
	"media":           AttributeSafe,
	"mediagroup":      AttributeSafe,
	"method":          FullyTrusted,
	"min":             AttributeSafe,
	"multiple":        AttributeSafe,
	"name":            AttributeSafe,
	"novalidate":      FullyTrusted,
	"open":            AttributeSafe,
	"optimum":         AttributeSafe,
	"pattern":         FullyTrusted,
	"placeholder":     AttributeSafe,
	"poster":          URLSafe,
	"profile":         URLSafe,
	"preload":         AttributeSafe,
	"pubdate":         AttributeSafe,
	"radiogroup":      AttributeSafe,
	"readonly":        AttributeSafe,
	"rel":             FullyTrusted,
	"required":        AttributeSafe,
	"reversed":        AttributeSafe,
	"rows":            AttributeSafe,
	"rowspan":         AttributeSafe,
	"sandbox":         FullyTrusted,
	"spellcheck":      AttributeSafe,
	"scope":           AttributeSafe,
	"scoped":          AttributeSafe,
	"seamless":        AttributeSafe,
	"selected":        AttributeSafe,
	"shape":           AttributeSafe,
	"size":            AttributeSafe,
	"sizes":           AttributeSafe,
	"span":            AttributeSafe,
	"src":             URLSafe,
	"srcdoc":          HTMLSafe,
	"srclang":         AttributeSafe,
	"srcset":          FullyTrusted,
	"start":           AttributeSafe,
	"step":            AttributeSafe,
	"style":           FullyTrusted,
	"tabindex":        AttributeSafe,
	"target":          AttributeSafe,
	"title":           AttributeSafe,
	"type":            FullyTrusted,
	"usemap":          URLSafe,
	"value":           FullyTrusted,
	"width":           AttributeSafe,
	"wrap":            AttributeSafe,
	"xmlns":           URLSafe,
}
