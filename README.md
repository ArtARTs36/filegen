# filegen

**filegen** - is console application for generate files of your templates

## Usage example

config.yaml
```yaml
vars:
  city: Moscow

files:
  - template_path: user.tmpl
    output_path: 'user_{{ vars.local.user.id }}.yaml'
    vars:
      user:
        id: 1
        name: "John"
```

user.tmpl
```yaml
user:
    id: {{ vars.local.user.id }}
    name: {{ vars.local.user.name }}
    city: {{ vars.global.user.city }}
```

Run `filegen config.yaml`
