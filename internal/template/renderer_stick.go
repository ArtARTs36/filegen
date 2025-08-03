package template

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/tyler-sommer/stick"
	"os"
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

	loadJSON := func(_ stick.Context, args ...stick.Value) stick.Value {
		if len(args) != 1 {
			panic("load_json expects one argument")
		}

		path, ok := args[0].(string)
		if !ok {
			panic("load_json expects string argument")
		}

		content, err := os.ReadFile(path)
		if err != nil {
			panic("unable to read " + path + ": " + err.Error())
		}

		var val interface{}

		if err = json.Unmarshal(content, &val); err != nil {
			panic("unable to unmarshal json " + path + ": " + err.Error())
		}

		return val
	}

	fileEngine.Functions["lower"] = lower
	fileEngine.Functions["upper"] = upper
	fileEngine.Functions["load_json"] = loadJSON

	stringEngine.Functions["lower"] = lower
	stringEngine.Functions["upper"] = upper
	stringEngine.Functions["load_json"] = loadJSON

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
	var res bytes.Buffer

	err := r.fileEngine.Execute(path, &res, r.remapVars(vars))
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
	return r.renderString(content, r.remapVars(vars))
}

func (r *StickRenderer) renderString(content string, vars map[string]stick.Value) ([]byte, error) {
	var res bytes.Buffer

	err := r.stringEngine.Execute(content, &res, vars)
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
