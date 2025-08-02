package template

import "context"

type Renderer interface {
	RenderFile(ctx context.Context, path string, vars map[string]interface{}) ([]byte, error)
	RenderString(ctx context.Context, content string, vars map[string]interface{}) ([]byte, error)
}
