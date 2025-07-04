package lorem

import (
	"bytes"
	"strings"
)

const (
	// Title Lorem Ipsum
	Title = "Lorem Ipsum"
	// SubTitle Neque porro quisquam...
	SubTitle = "Neque porro quisquam est qui dolorem ipsum quia dolor sit amet, consectetur."
	// OneParagraph Lorem ipsum dolor sit amet, consectetur...
	OneParagraph = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nunc at aliquam justo. Donec dignissim nulla quis metus aliquet, eu sagittis diam elementum. Nullam interdum aliquet enim a maximus. Ut ornare mollis blandit. Phasellus ut nisi vel eros bibendum dictum. Praesent rhoncus viverra luctus. Pellentesque quam velit, commodo sed nunc non, faucibus luctus urna. Pellentesque aliquet, elit et tristique commodo, urna purus imperdiet libero, et tempor erat felis eget nibh. Integer pellentesque est id mi malesuada finibus. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas.`
	// OneParagraphHTML Lorem ipsum with <div><p>...
	OneParagraphHTML = `<div><p>` + OneParagraph + `</p></div>`
	// ThreeParagraph three long line
	ThreeParagraph = `Mauris eu quam ut elit ultricies pretium. Sed at ligula tempor, aliquam elit ac, porta leo. In hendrerit, sem ac tincidunt bibendum, ex purus viverra sapien, eu finibus ante nunc ut lacus. Sed ornare mi at nisi fringilla dapibus. In dapibus erat vitae odio venenatis, id dignissim lorem mollis. Mauris vel metus ullamcorper, varius ligula et, laoreet ligula. In vitae accumsan velit. Donec facilisis tortor at ex hendrerit, at commodo urna faucibus. Phasellus in eros ut ante consequat dignissim. Integer ac lobortis felis. Pellentesque suscipit, dui vitae aliquam tincidunt, diam metus placerat erat, quis dignissim dui neque pellentesque augue. Fusce vitae libero bibendum est auctor tristique eu eget dui.
Pellentesque porttitor condimentum velit sed mattis. In ultricies, urna ut pretium lobortis, augue urna imperdiet lorem, eu viverra massa tellus eget nibh. Nulla a nulla purus. Aenean fermentum rhoncus aliquam. Nullam sodales enim ante, nec aliquam risus imperdiet at. Vivamus rutrum efficitur tempor. Nullam eu libero ultrices metus interdum eleifend a eget odio. Nullam vel nisl enim. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Nam ac mi ac massa interdum eleifend.
Integer molestie ullamcorper sapien, ac semper odio vehicula imperdiet. Pellentesque at tellus libero. Cras consectetur mi sit amet ultricies pulvinar. Ut tempor quam in mi placerat, non viverra lacus faucibus. Cras maximus turpis facilisis mi consequat sollicitudin. Pellentesque at posuere turpis. Donec nec hendrerit ex. Integer commodo nisi eu consequat varius. Nunc quis dolor vitae ante tincidunt mollis. Duis efficitur metus sed purus efficitur, eget bibendum sapien porttitor.`
)

// ThreeParagraphHTML three long line with <p>. Wrapped with <div>
var ThreeParagraphHTML = func() string {
	buf := &bytes.Buffer{}
	buf.WriteString("<div>")
	lines := strings.Split(ThreeParagraph, "\n")
	for _, s := range lines {
		buf.WriteString("<p>")
		buf.WriteString(s)
		buf.WriteString("</p>")
	}
	buf.WriteString("</div>")
	return buf.String()
}()
