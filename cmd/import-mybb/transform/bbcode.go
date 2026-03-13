package transform

import (
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/frustra/bbcode"
)

var bbcodeCompiler bbcode.Compiler

func init() {
	bbcodeCompiler = bbcode.NewCompiler(true, true)

	// Add custom MyBB BBCode tag mappings

	// [align=center], [align=left], [align=right], [align=justify]
	bbcodeCompiler.SetTag("align", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "div"

		alignment := node.GetOpeningTag().Value
		if alignment == "" {
			alignment = "left"
		}
		out.Attrs["style"] = "text-align: " + alignment + ";"

		return out, true
	})

	// [hr] - horizontal rule
	bbcodeCompiler.SetTag("hr", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "hr"
		return out, false
	})

	// [php] and [PHP] - code blocks for PHP
	phpHandler := func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "pre"

		code := bbcode.NewHTMLTag(node.GetOpeningTag().Raw)
		code.Name = "code"
		code.Attrs["class"] = "language-php"
		code.Value = bbcode.CompileText(node)

		out.AppendChild(code)
		return out, false
	}
	bbcodeCompiler.SetTag("php", phpHandler)
	bbcodeCompiler.SetTag("PHP", phpHandler)

	// [video=youtube] - YouTube embeds converted to links (iframes are sanitized out)
	bbcodeCompiler.SetTag("video", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "a"

		videoURL := bbcode.CompileText(node)
		videoType := node.GetOpeningTag().Value

		if videoType == "youtube" {
			// Extract video ID from various YouTube URL formats
			// Could be: https://www.youtube.com/watch?v=VIDEO_ID or just VIDEO_ID
			videoID := videoURL
			if len(videoURL) > 20 {
				// It's a full URL, try to extract the ID
				if idx := len(videoURL) - 11; idx >= 0 && len(videoURL) >= 11 {
					// YouTube video IDs are typically 11 characters
					videoID = videoURL[len(videoURL)-11:]
				}
			}
			out.Attrs["href"] = "https://www.youtube.com/watch?v=" + videoID
			out.Value = "YouTube Video: " + videoID
		} else {
			// Fallback for unknown video types
			out.Attrs["href"] = videoURL
			out.Value = "Video: " + videoURL
		}

		return out, false
	})

	// [font=...] - font family (we'll just ignore the font and keep the text)
	bbcodeCompiler.SetTag("font", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "span"
		// We don't set the font-family style to avoid security issues with user-controlled fonts
		// Just pass through the text content
		return out, true
	})

	// [list] and [list=1] - unordered and ordered lists
	bbcodeCompiler.SetTag("list", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")

		listType := node.GetOpeningTag().Value
		if listType == "1" || listType == "a" || listType == "A" {
			out.Name = "ol"
			if listType == "a" {
				out.Attrs["type"] = "a"
			} else if listType == "A" {
				out.Attrs["type"] = "A"
			}
		} else {
			out.Name = "ul"
		}

		return out, true
	})

	// [*] - list item (MyBB uses this for list items)
	bbcodeCompiler.SetTag("*", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "li"
		return out, true
	})
}

// convertBBCodeToHTML converts MyBB BBCode to HTML and returns a sanitized Content object
func convertBBCodeToHTML(input string) datagraph.Content {
	if input == "" {
		// Return empty content
		c, _ := datagraph.NewRichText("")
		return c
	}

	// Convert BBCode to HTML
	html := bbcodeCompiler.Compile(input)

	// Use NewRichText to sanitize and generate short summary
	content, err := datagraph.NewRichText(html)
	if err != nil {
		// If sanitization fails, return empty content
		// This shouldn't normally happen but handle it gracefully
		return datagraph.Content{}
	}

	return content
}
