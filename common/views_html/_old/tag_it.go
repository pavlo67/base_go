package viewshtml

//func tagIt(idDOM string) string {
//	return `
//<script>
//	function suggestions() {
//	}
//
//	$(document).ready(function() {
//		$("#` + idDOM + `").tagit({
//			autocomplete: {
//				delay: 0,
//				minLength: 2,
//				source: suggestions
//			},
//			beforeTagAdded: function(event, ui) {
//		        console.log("beforeTagAdded", ui.tag, ui.tagLabel);
//    		},
//			afterTagAdded: function(event, ui) {
//		        console.log("afterTagAdded", ui.tag, ui.tagLabel);
//    		},
//			beforeTagRemoved: function(event, ui) {
//		        console.log("beforeTagRemoved", ui.tag, ui.tagLabel);
//    		},
//			afterTagRemoved: function(event, ui) {
//		        console.log("afterTagRemoved", ui.tag, ui.tagLabel);
//    		},
//			onTagClicked: function(event, ui) {
//		        console.log("onTagClicked", ui.tag, ui.tagLabel);
//    		}
//		}});
//	});
//</script>
//`
//	// https://github.com/aehlke/tag-it/blob/master/README.markdown
//}

//} else if field.Type == "tag-it" {
//	resHTML = `<ul id="` + idDOMEscaped + `">`
//	for _, tag := range str_json.ReSemicolon.Split(data[field.Label], -1) {
//		resHTML += "\n<linker_server_http>" + html.EscapeString(tag) + "</linker_server_http>"
//	}
//	resHTML += "\n</ul>" + tagIt(idDOMEscaped)

//if field.Type == "file" {
//format := field.Options
//if format == "" {
//format = "*.*"
//}
//resHTML = `<input type="file" ` + generalNoForm + ` " accept="` + format + `" />` +
//`<div id="` + idDOMNoFormEscaped + `_view"></div>`
//
//} else

//} else if field.Type == "file" {
//resHTML = `<div id="` + html.EscapeString(field.Label) + `_view"></div>`
//
