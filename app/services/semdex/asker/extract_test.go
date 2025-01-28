package asker

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Southclaws/storyden/app/services/semdex"
)

func Test_streamExtractor(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		a := assert.New(t)

		stream := []string{
			"For", " finance data", ", several", " datasets can", " be useful", " depending on", " the specific", " needs of", " your project", ". Here", " are some", " recommendations:\n\n", "1.", " **Credit", " Analytics Dataset", " by S", "\u0026P Global", "**: This", " dataset provides", " a comprehensive", " range of", " credit risk", " indicators,", " including market", "-based and", " fundamental-based", " scores.", " It covers", " over ", "881,", "000 public", " and private", " companies across", " 160", " industries,", " 183", " countries,", " and ", "241 S", "\u0026P Dow", " Jones Indices", ". It", " includes pre", "-calculated", " probability of", " default and", " credit scores", ", early", " warning indicators", " for credit", " deterioration,", " and systematic", " risk factors", " like country", " and industry", " risk scores", ".\n\n2", ". **", "Prosper", " Loan Dataset", " on Kag", "gle**:", " This dataset", " contains information", " on loan", " origination", " dates,", " original amounts", ", payments", ", returns", ", and", " loss rates", ". It", " is useful", " for analyzing", " loan performance", " and can", " be used", " to create", " a portfolio", " project in", " finance analytics", ". The", " dataset includes", " metrics viable", " from ", "2009", ", and", " it can", " be filtered", " to exclude", " missing values", " and focus", " on specific", " time periods", ".\n\n3", ". **", "Finance and", " Risk Analytics", " Dataset on", " GitHub**:", " This dataset", " contains ", "6 years", " of weekly", " stock information", " for ", "10 different", " Indian stocks", ". It", " can be", " used for", " market risk", " analysis and", " includes data", " on stock", " prices,", " returns,", " and share", " insights.", " This dataset", " is particularly", " useful for", " those interested", " in stock", " market data", " and risk", " analysis.\n\n", "4.", " **Data", "hub.io", "**: This", " platform offers", " a variety", " of business", " and finance", " datasets,", " including stock", " market data", ", property", " prices,", " inflation,", " and logistics", ". The", " data is", " often updated", " monthly or", " daily,", " providing a", " constant flow", " of real", "-time information", ".\n\nThese", " datasets can", " be used", " to create", " a portfolio", " project in", " finance analytics", ", covering", " various aspects", " of finance", " such as", " credit risk", ", loan", " performance,", " and market", " risk.\n\n", "### Sources", "\n-", " https://", "www.market", "place.s", "pglobal", ".com/en", "/datasets", "/credit", "-analytics-(", "146)\n", "- https", "://www", ".youtube.com", "/watch?v", "=rS", "0NU", "ngQ", "cbU", "\n-", " https://", "github.com", "/Honey", "28Git", "/Finance", "-and-R", "isk-An", "alytics\n", "- https", "://career", "foundry", ".com/en", "/blog/data", "-analytics/", "where-to", "-find-free", "-datasets", "/",
		}

		iter := func(yield func(string, error) bool) {
			for _, v := range stream {
				if !yield(v, nil) {
					return
				}
			}
		}

		acc := []semdex.AskResponseChunk{}
		buf := ""
		meta := semdex.AskResponseChunkMeta{}
		for i := range streamExtractor(iter) {
			switch v := i.(type) {
			case *semdex.AskResponseChunkText:
				buf += v.Chunk
			case *semdex.AskResponseChunkMeta:
				meta = *v
			}

			acc = append(acc, i)
		}

		a.Equal(strings.Join(stream, ""), buf)
		a.Len(meta.URLs, 4)
	})

	t.Run("basic_sdrs", func(t *testing.T) {
		a := assert.New(t)

		stream := []string{
			"For", " finance data", ", several", " datasets can", " be useful", " depending on", " the specific", " needs of", " your project", ". Here", " are some", " recommendations:\n\n", "1.", " **Credit", " Analytics Dataset", " by S", "\u0026P Global", "**: This", " dataset provides", " a comprehensive", " range of", " credit risk", " indicators,", " including market", "-based and", " fundamental-based", " scores.", " It covers", " over ", "881,", "000 public", " and private", " companies across", " 160", " industries,", " 183", " countries,", " and ", "241 S", "\u0026P Dow", " Jones Indices", ". It", " includes pre", "-calculated", " probability of", " default and", " credit scores", ", early", " warning indicators", " for credit", " deterioration,", " and systematic", " risk factors", " like country", " and industry", " risk scores", ".\n\n2", ". **", "Prosper", " Loan Dataset", " on Kag", "gle**:", " This dataset", " contains information", " on loan", " origination", " dates,", " original amounts", ", payments", ", returns", ", and", " loss rates", ". It", " is useful", " for analyzing", " loan performance", " and can", " be used", " to create", " a portfolio", " project in", " finance analytics", ". The", " dataset includes", " metrics viable", " from ", "2009", ", and", " it can", " be filtered", " to exclude", " missing values", " and focus", " on specific", " time periods", ".\n\n3", ". **", "Finance and", " Risk Analytics", " Dataset on", " GitHub**:", " This dataset", " contains ", "6 years", " of weekly", " stock information", " for ", "10 different", " Indian stocks", ". It", " can be", " used for", " market risk", " analysis and", " includes data", " on stock", " prices,", " returns,", " and share", " insights.", " This dataset", " is particularly", " useful for", " those interested", " in stock", " market data", " and risk", " analysis.\n\n", "4.", " **Data", "hub.io", "**: This", " platform offers", " a variety", " of business", " and finance", " datasets,", " including stock", " market data", ", property", " prices,", " inflation,", " and logistics", ". The", " data is", " often updated", " monthly or", " daily,", " providing a", " constant flow", " of real", "-time information", ".\n\nThese", " datasets can", " be used", " to create", " a portfolio", " project in", " finance analytics", ", covering", " various aspects", " of finance", " such as", " credit risk", ", loan", " performance,", " and market", " risk.\n\n", "### Sources", "\n-", " sdr:", "thread", "/cpvf89ifunp0qr2aqp2g\n", "- sdr", ":node/", "cpvf89ifunp0qr2aqp8g", "\n-", " sdr:profile", "/cpvf89ifunp0qr2aqp8g",
		}

		iter := func(yield func(string, error) bool) {
			for _, v := range stream {
				if !yield(v, nil) {
					return
				}
			}
		}

		acc := []semdex.AskResponseChunk{}
		buf := ""
		meta := semdex.AskResponseChunkMeta{}
		for i := range streamExtractor(iter) {
			switch v := i.(type) {
			case *semdex.AskResponseChunkText:
				buf += v.Chunk
			case *semdex.AskResponseChunkMeta:
				meta = *v
			}

			acc = append(acc, i)
		}

		a.Equal(strings.Join(stream, ""), buf)
		a.Len(meta.URLs, 0)
		a.Len(meta.Refs, 3)
	})

	t.Run("edge_case_full_url_in_single_chunk", func(t *testing.T) {
		a := assert.New(t)

		stream := []string{
			"this", "is", "a", "chunk", "with", "a", "full", "url", "https://www.example.com", "in", "it",
		}

		iter := func(yield func(string, error) bool) {
			for _, v := range stream {
				if !yield(v, nil) {
					return
				}
			}
		}

		acc := []semdex.AskResponseChunk{}
		buf := ""
		meta := semdex.AskResponseChunkMeta{}
		for i := range streamExtractor(iter) {
			switch v := i.(type) {
			case *semdex.AskResponseChunkText:
				buf += v.Chunk
			case *semdex.AskResponseChunkMeta:
				meta = *v
			}

			acc = append(acc, i)
		}

		a.Equal(strings.Join(stream, ""), buf)
		a.Len(meta.URLs, 1)
		a.Len(meta.Refs, 0)
	})

	t.Run("edge_case_full_url_plus_more_in_single_chunk", func(t *testing.T) {
		a := assert.New(t)

		stream := []string{
			"this ", "is", " a", " chunk", " with", " a", " full", " url ", "https://www.example.com\nin", " it",
		}

		iter := func(yield func(string, error) bool) {
			for _, v := range stream {
				if !yield(v, nil) {
					return
				}
			}
		}

		acc := []semdex.AskResponseChunk{}
		buf := ""
		meta := semdex.AskResponseChunkMeta{}
		for i := range streamExtractor(iter) {
			switch v := i.(type) {
			case *semdex.AskResponseChunkText:
				buf += v.Chunk
			case *semdex.AskResponseChunkMeta:
				meta = *v
			}

			acc = append(acc, i)
		}

		a.Equal(strings.Join(stream, ""), buf)
		a.Len(meta.URLs, 1)
		a.Len(meta.Refs, 0)
	})

	t.Run("supports_break", func(t *testing.T) {
		a := assert.New(t)

		stream := []string{
			"For", " finance data", ", several", " datasets can", " be useful", " depending on", " the specific", " needs of", " your project", ". Here", " are some", " recommendations:\n\n", "1.", " **Credit", " Analytics Dataset", " by S", "\u0026P Global", "**: This", " dataset provides", " a comprehensive", " range of", " credit risk", " indicators,", " including market", "-based and", " fundamental-based", " scores.", " It covers", " over ", "881,", "000 public", " and private", " companies across", " 160", " industries,", " 183", " countries,", " and ", "241 S", "\u0026P Dow", " Jones Indices", ". It", " includes pre", "-calculated", " probability of", " default and", " credit scores", ", early", " warning indicators", " for credit", " deterioration,", " and systematic", " risk factors", " like country", " and industry", " risk scores", ".\n\n2", ". **", "Prosper", " Loan Dataset", " on Kag", "gle**:", " This dataset", " contains information", " on loan", " origination", " dates,", " original amounts", ", payments", ", returns", ", and", " loss rates", ". It", " is useful", " for analyzing", " loan performance", " and can", " be used", " to create", " a portfolio", " project in", " finance analytics", ". The", " dataset includes", " metrics viable", " from ", "2009", ", and", " it can", " be filtered", " to exclude", " missing values", " and focus", " on specific", " time periods", ".\n\n3", ". **", "Finance and", " Risk Analytics", " Dataset on", " GitHub**:", " This dataset", " contains ", "6 years", " of weekly", " stock information", " for ", "10 different", " Indian stocks", ". It", " can be", " used for", " market risk", " analysis and", " includes data", " on stock", " prices,", " returns,", " and share", " insights.", " This dataset", " is particularly", " useful for", " those interested", " in stock", " market data", " and risk", " analysis.\n\n", "4.", " **Data", "hub.io", "**: This", " platform offers", " a variety", " of business", " and finance", " datasets,", " including stock", " market data", ", property", " prices,", " inflation,", " and logistics", ". The", " data is", " often updated", " monthly or", " daily,", " providing a", " constant flow", " of real", "-time information", ".\n\nThese", " datasets can", " be used", " to create", " a portfolio", " project in", " finance analytics", ", covering", " various aspects", " of finance", " such as", " credit risk", ", loan", " performance,", " and market", " risk.\n\n", "### Sources", "\n-", " https://", "www.market", "place.s", "pglobal", ".com/en", "/datasets", "/credit", "-analytics-(", "146)\n", "- https", "://www", ".youtube.com", "/watch?v", "=rS", "0NU", "ngQ", "cbU", "\n-", " https://", "github.com", "/Honey", "28Git", "/Finance", "-and-R", "isk-An", "alytics\n", "- https", "://career", "foundry", ".com/en", "/blog/data", "-analytics/", "where-to", "-find-free", "-datasets", "/",
		}

		iter := func(yield func(string, error) bool) {
			for _, v := range stream {
				if !yield(v, nil) {
					return
				}
			}
		}

		acc := []semdex.AskResponseChunk{}
		buf := ""
		meta := semdex.AskResponseChunkMeta{}
		count := 0
		for i := range streamExtractor(iter) {
			switch v := i.(type) {
			case *semdex.AskResponseChunkText:
				buf += v.Chunk
			case *semdex.AskResponseChunkMeta:
				meta = *v
			}

			acc = append(acc, i)

			// break at the 173rd token. This token is the beginning of a URL.
			// This means the source tokens are advanced until the end of the
			// URL and the actual consumed tokens are 179 because the range from
			// 173 to 179 is the full URL:
			// "https://www.marketplace.spglobal.com/en/datasets/credit-analytics-(146)"
			count++
			if count == 173 {
				break
			}
		}

		// 179 is the last index of the stream before the break
		// remove the trailing newline because while this char is actually part
		// of the 179th token, it's not returned by yieldPeekedReference because
		// it doesn't yield actual tokens, it yields the entire parsed URL.
		wantBuf := strings.TrimRight(strings.Join(stream[:179], ""), "\n")

		a.Equal(wantBuf, buf)
		a.Len(meta.URLs, 1)
	})

	t.Run("supports_break_on_input_stream", func(t *testing.T) {
		a := assert.New(t)

		stream := []string{
			"For", " finance data", ", several", " datasets can", " be useful", " depending on", " the specific", " needs of", " your project", ". Here", " are some", " recommendations:\n\n", "1.", " **Credit", " Analytics Dataset", " by S", "\u0026P Global", "**: This", " dataset provides", " a comprehensive", " range of", " credit risk", " indicators,", " including market", "-based and", " fundamental-based", " scores.", " It covers", " over ", "881,", "000 public", " and private", " companies across", " 160", " industries,", " 183", " countries,", " and ", "241 S", "\u0026P Dow", " Jones Indices", ". It", " includes pre", "-calculated", " probability of", " default and", " credit scores", ", early", " warning indicators", " for credit", " deterioration,", " and systematic", " risk factors", " like country", " and industry", " risk scores", ".\n\n2", ". **", "Prosper", " Loan Dataset", " on Kag", "gle**:", " This dataset", " contains information", " on loan", " origination", " dates,", " original amounts", ", payments", ", returns", ", and", " loss rates", ". It", " is useful", " for analyzing", " loan performance", " and can", " be used", " to create", " a portfolio", " project in", " finance analytics", ". The", " dataset includes", " metrics viable", " from ", "2009", ", and", " it can", " be filtered", " to exclude", " missing values", " and focus", " on specific", " time periods", ".\n\n3", ". **", "Finance and", " Risk Analytics", " Dataset on", " GitHub**:", " This dataset", " contains ", "6 years", " of weekly", " stock information", " for ", "10 different", " Indian stocks", ". It", " can be", " used for", " market risk", " analysis and", " includes data", " on stock", " prices,", " returns,", " and share", " insights.", " This dataset", " is particularly", " useful for", " those interested", " in stock", " market data", " and risk", " analysis.\n\n", "4.", " **Data", "hub.io", "**: This", " platform offers", " a variety", " of business", " and finance", " datasets,", " including stock", " market data", ", property", " prices,", " inflation,", " and logistics", ". The", " data is", " often updated", " monthly or", " daily,", " providing a", " constant flow", " of real", "-time information", ".\n\nThese", " datasets can", " be used", " to create", " a portfolio", " project in", " finance analytics", ", covering", " various aspects", " of finance", " such as", " credit risk", ", loan", " performance,", " and market", " risk.\n\n", "### Sources", "\n-", " https://", "www.market", "place.s", "pglobal", ".com/en", "/datasets", "/credit", "-analytics-(", "146)\n", "- https", "://www", ".youtube.com", "/watch?v", "=rS", "0NU", "ngQ", "cbU", "\n-", " https://", "github.com", "/Honey", "28Git", "/Finance", "-and-R", "isk-An", "alytics\n", "- https", "://career", "foundry", ".com/en", "/blog/data", "-analytics/", "where-to", "-find-free", "-datasets", "/",
		}

		iter := func(yield func(string, error) bool) {
			for i, v := range stream {
				// Break the input stream at the 173rd token. This tests that
				// the streamExtractor break handling functionality works.
				if i == 173 {
					break
				}

				if !yield(v, nil) {
					return
				}
			}
		}

		acc := []semdex.AskResponseChunk{}
		buf := ""
		meta := semdex.AskResponseChunkMeta{}
		count := 0
		for i := range streamExtractor(iter) {
			switch v := i.(type) {
			case *semdex.AskResponseChunkText:
				buf += v.Chunk
			case *semdex.AskResponseChunkMeta:
				meta = *v
			}

			acc = append(acc, i)

			// break at the 173rd token. This token is the beginning of a URL.
			// This means the source tokens are advanced until the end of the
			// URL and the actual consumed tokens are 179 because the range from
			// 173 to 179 is the full URL:
			// "https://www.marketplace.spglobal.com/en/datasets/credit-analytics-(146)"
			count++
			if count == 173 {
				break
			}
		}

		// 179 is the last index of the stream before the break
		// remove the trailing newline because while this char is actually part
		// of the 179th token, it's not returned by yieldPeekedReference because
		// it doesn't yield actual tokens, it yields the entire parsed URL.
		wantBuf := strings.TrimRight(strings.Join(stream[:173], ""), "\n")

		a.Equal(wantBuf, buf)
		a.Len(meta.URLs, 1)
	})
}
