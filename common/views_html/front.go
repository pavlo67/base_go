package views_html

type Operator interface {
	// HTMLToEdit returns an html code to view the element with selected modelID
	HTMLToView(field Field, data ValuesString, formID string) string

	// HTMLToEdit returns an html code to edit the element with selected modelID
	HTMLToEdit(field Field, data ValuesString, formID string) string

	// HTMLToLoad returns an html code to load all required js/css partes
	HTMLToLoad() string
}

func New(htmlToView, htmlToEdit HTMLToShow, htmlFront string) Operator {
	return &frontComponent{
		htmlToView: htmlToView,
		htmlToEdit: htmlToEdit,
		htmlFront:  htmlFront,
	}
}

// -----------------------------------------------------------------------------------------------------

var _ Operator = &frontComponent{}

type HTMLToShow func(field Field, data ValuesString, formID string) string

type frontComponent struct {
	htmlToView HTMLToShow
	htmlToEdit HTMLToShow
	htmlFront  string
}

// HTMLToEdit returns an html code to join this front component to the element with selected modelID
func (fc frontComponent) HTMLToEdit(field Field, data ValuesString, formID string) string {
	return fc.htmlToEdit(field, data, formID)
}

// HTMLToView returns an html code to join this front component to the element with selected modelID
func (fc frontComponent) HTMLToView(field Field, data ValuesString, formID string) string {
	return fc.htmlToView(field, data, formID)
}

// HTMLToLoad returns an html code to load all required js/css partes
func (fc frontComponent) HTMLToLoad() string {
	return fc.htmlFront
}
