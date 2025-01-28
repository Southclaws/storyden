package asker

import (
	"net/url"
	"strings"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/services/semdex"
)

func streamExtractor(iter func(yield func(string, error) bool)) semdex.AskResponseIterator {
	urlPeek := false
	peek := strings.Builder{}

	// Metadata is accumulated as new references are discovered in the stream.
	acc := semdex.AskResponseChunkMeta{
		Refs: []*datagraph.Ref{},
		URLs: []url.URL{},
	}

	return func(yield func(semdex.AskResponseChunk, error) bool) {
		yieldPeekedReference := func() bool {
			peekedURL := peek.String()
			peek.Reset()

			if !yield(&semdex.AskResponseChunkText{
				Chunk: peekedURL,
			}, nil) {
				return false
			}

			// Parse the URL to make sure it's valid.
			parsed, err := url.Parse(peekedURL)
			if err != nil {
				// do nothing, yield nothing, continue.
				return true
			}

			switch parsed.Scheme {
			case "http", "https":
				acc.URLs = append(acc.URLs, *parsed)

			case datagraph.RefScheme:
				ref, err := datagraph.NewRefFromSDR(*parsed)
				if err != nil {
					return true
				}
				acc.Refs = append(acc.Refs, ref)
			}

			return yield(&acc, nil)
		}

		tryYieldReference := func(chunk string) (string, bool) {
			// If we're peeking for a URL, take the next chunk and look for
			// a boundary that ends the URL. If there is one, we have a URL.
			// If not, continue peeking until we find one.
			urlEnd := findURLEnd(chunk)
			if urlEnd == -1 {
				peek.WriteString(chunk)
				return "", true
			}
			// We've found the end of the URL, reset the peek buffer.
			urlPeek = false

			peek.WriteString(chunk[:urlEnd])
			chunk = chunk[urlEnd:]

			// we have a full URL in `peek`, yield it and reset the peek buffer.
			if !yieldPeekedReference() {
				return "", false
			}

			// since yieldPeekedReference yields the URL part of the text chunk,
			// yield the remainder of the chunk here.
			return chunk, true
		}

		for chunk, err := range iter {
			if err != nil {
				yield(nil, err)
				return
			}

			if urlPeek {
				var y bool
				chunk, y = tryYieldReference(chunk)
				if !y {
					return
				}
			}

			urlStart := findURLStart(chunk)
			if urlStart != -1 {
				urlEnd := findURLEnd(chunk[urlStart:])
				if urlEnd == -1 {
					// End of the URL is in a future chunk, peek for it.
					urlPeek = true
					peek.WriteString(chunk[urlStart:])
					chunk = chunk[:urlStart]
				} else {
					// End of the URL is in this chunk, yield the URL now.
					peek.WriteString(chunk[urlStart:urlEnd])
					chunk = chunk[urlEnd:]

					yieldPeekedReference()
				}
			}

			if chunk == "" {
				continue
			}

			if !yield(&semdex.AskResponseChunkText{
				Chunk: chunk,
			}, nil) {
				return
			}

		}

		if urlPeek {
			yieldPeekedReference()
		}
	}
}

func findURLStart(chunk string) int {
	https := strings.Index(chunk, "https")
	if https != -1 {
		return https
	}

	http := strings.Index(chunk, "http")
	if http != -1 {
		return http
	}

	sdr := strings.Index(chunk, datagraph.RefScheme)
	if sdr != -1 {
		return sdr
	}

	return -1
}

func findURLEnd(chunk string) int {
	pos := strings.IndexAny(chunk, " \n")
	if pos != -1 {
		return pos
	}

	return -1
}
