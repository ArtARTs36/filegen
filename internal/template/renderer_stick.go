package template

import (
	"bytes"
	"context"
	"github.com/tyler-sommer/stick"
	"strings"
)

type StickRenderer struct {
	engine *stick.Env
}

func NewStickRenderer() *StickRenderer {
	engine := stick.New(nil)

	engine.Functions["lower"] = func(_ stick.Context, args ...stick.Value) stick.Value {
		if len(args) != 1 {
			return nil
		}

		return strings.ToLower(args[0].(string))
	}

	engine.Functions["upper"] = func(_ stick.Context, args ...stick.Value) stick.Value {
		if len(args) != 1 {
			return nil
		}

		return strings.ToUpper(args[0].(string))
	}

	return &StickRenderer{
		engine: engine,
	}
}

func (r *StickRenderer) Render(
	_ context.Context,
	content []byte,
	vars map[string]interface{},
) ([]byte, error) {
	stickVars := r.remapVars(vars)
	var res bytes.Buffer

	err := r.engine.Execute(string(content), &res, stickVars)
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
