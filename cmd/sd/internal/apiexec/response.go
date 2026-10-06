package apiexec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func renderResponse(out io.Writer, input Request, response *http.Response) error {
	if input.OutputFile != "" {
		file, err := os.OpenFile(input.OutputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
		_, copyErr := io.Copy(file, response.Body)
		closeErr := file.Close()
		if copyErr != nil {
			return copyErr
		}
		return closeErr
	}
	if input.Output == "raw" {
		_, err := io.Copy(out, response.Body)
		return err
	}
	if response.Request != nil && response.Request.Method == http.MethodHead {
		return writeJSON(out, map[string]any{"status": response.StatusCode, "headers": response.Header})
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if len(bytes.TrimSpace(body)) == 0 {
		return writeJSON(out, nil)
	}
	if !json.Valid(body) {
		return fmt.Errorf("response is not JSON; use --output raw or --output-file")
	}
	if input.Output == "json" {
		var pretty bytes.Buffer
		if err := json.Indent(&pretty, body, "", "  "); err != nil {
			return err
		}
		pretty.WriteByte('\n')
		_, err := out.Write(pretty.Bytes())
		return err
	}
	return nil
}

func writeJSON(out io.Writer, value any) error {
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}
