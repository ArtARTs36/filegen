package template

import (
	"bytes"
	"context"
	"github.com/tyler-sommer/stick"
	"strings"
)

type StickRenderer struct {
	fileEngine   *stick.Env
	stringEngine *stick.Env
}

func NewStickRenderer(loader stick.Loader) *StickRenderer {
	fileEngine := stick.New(loader)
	stringEngine := stick.New(nil)

	lower := func(_ stick.Context, args ...stick.Value) stick.Value {
		if len(args) != 1 {
			return nil
		}

		arg, ok := args[0].(string)
		if !ok {
			return nil
		}

		return strings.ToLower(arg)
	}

	upper := func(_ stick.Context, args ...stick.Value) stick.Value {
		if len(args) != 1 {
			return nil
		}

		arg, ok := args[0].(string)
		if !ok {
			return nil
		}

		return strings.ToUpper(arg)
	}

	fileEngine.Functions["lower"] = lower
	fileEngine.Functions["upper"] = upper

	stringEngine.Functions["lower"] = lower
	stringEngine.Functions["upper"] = upper

	return &StickRenderer{
		fileEngine:   fileEngine,
		stringEngine: stringEngine,
	}
}

func (r *StickRenderer) RenderFile(
	_ context.Context,
	path string,
	vars map[string]interface{},
) ([]byte, error) {
	stickVars := r.remapVars(vars)
	var res bytes.Buffer

	err := r.fileEngine.Execute(path, &res, stickVars)
	if err != nil {
		return nil, err
	}

	return res.Bytes(), nil
}

func (r *StickRenderer) RenderString(
	_ context.Context,
	content string,
	vars map[string]interface{},
) ([]byte, error) {
	stickVars := r.remapVars(vars)
	var res bytes.Buffer

	err := r.stringEngine.Execute(content, &res, stickVars)
	if err != nil {
		return nil, err
	}

	return res.Bytes(), nil
}

func (r *StickRenderer) remapVars(params map[string]interface{}) map[string]stick.Value {
	stickVars := map[string]stick.Value{}
	for key, val := range params {
		switch v := val.(type) {
		case map[string]interface{}:
			stickVars[key] = r.remapVars(v)
		default:
			stickVars[key] = v
		}
	}

	return stickVars
}
