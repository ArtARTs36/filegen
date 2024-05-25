package template

import "context"

type Renderer interface {
	Render(ctx context.Context, content []byte, vars map[string]interface{}) ([]byte, error)
}
