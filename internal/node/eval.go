package node

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/samber/lo"
)

func ResolveInput(node Node, states map[string]any) (resolvedProperties map[string]any, err error) {
	re := regexp.MustCompile(`{{.*}}`)
	properties := node.GetProperties(PropertyContext(states))
	propertiesMap := lo.Reduce(properties, func(acc map[string]NodeProperty, property NodeProperty, _ int) map[string]NodeProperty {
		acc[property.Key] = property
		return acc
	}, make(map[string]NodeProperty))
	resolvedProperties = make(map[string]any)
	for key, value := range propertiesMap {
		switch v := value.Value.(type) {
		// it can be a expression, binding, or a value
		// binding is @[bindingName]
		// expression is {{expression}}
		case string:
			if strings.HasPrefix(v, "@[") && strings.HasSuffix(v, "]") {
				// it is a binding
				binding := strings.TrimSpace(v[2 : len(v)-1])
				bindingValue, ok := states[binding]
				if !ok {
					return nil, fmt.Errorf("binding %s not found in states", binding)
				}
				resolvedProperties[key] = bindingValue
			} else if strings.HasPrefix(v, "{{") && strings.HasSuffix(v, "}}") {
				// it is a expression, evaluate it
				expression := v[2 : len(v)-2]
				resolvedProperties[key], err = Eval(expression, states, node.Ref())
				if err != nil {
					return nil, fmt.Errorf("failed to evaluate expression for property %s:%s (%w)", key, expression, err)
				}
			} else {
				// it is a value or embedded multiple expressions
				if re.MatchString(v) {
					expressions := re.FindAllString(v, -1)
					for _, expression := range expressions {
						evaluated, err := Eval(expression[2:len(expression)-2], states, node.Ref())
						if err != nil {
							return nil, fmt.Errorf("failed to evaluate expression for property %s:%s (%w)", key, expression, err)
						}
						// replace the expression with the evaluated value
						v = strings.Replace(v, expression, fmt.Sprintf("%v", evaluated), 1)
					}
					resolvedProperties[key] = v
				} else {
					resolvedProperties[key] = v
				}
			}
		default:
			resolvedProperties[key] = v
		}
	}
	return resolvedProperties, nil
}

// Eval evaluates an expression and returns the result.
// states is the current states of the workflow it have state of the nodes input,output datas
// - {node id}.input.{property key name} = value
// - {node id}.output.{output key name} = value
func Eval(expression string, states map[string]any, ref string) (any, error) {
	expression = strings.TrimSpace(expression)
	
	// Create evaluation environment with states and built-in variables
	env := make(map[string]any)
	
	// Parse states and create nested structure for better access
	for key, value := range states {
		parts := strings.Split(key, ".")
		if len(parts) == 3 {
			// Format: nodeId.type.property where type is "input" or "output"
			nodeId := parts[0]
			nodeType := parts[1]
			property := parts[2]
			
			// Create nested structure if not exists
			if _, ok := env[nodeId]; !ok {
				env[nodeId] = make(map[string]any)
			}
			if nodeMap, ok := env[nodeId].(map[string]any); ok {
				if _, ok := nodeMap[nodeType]; !ok {
					nodeMap[nodeType] = make(map[string]any)
				}
				if typeMap, ok := nodeMap[nodeType].(map[string]any); ok {
					typeMap[property] = value
				}
			}
		} else {
			// Direct key-value mapping for other states
			env[key] = value
		}
	}
	
	// Add special variables
	env["$ref"] = ref
	
	// Add built-in functions
	env["len"] = func(v any) int {
		switch val := v.(type) {
		case string:
			return len(val)
		case []any:
			return len(val)
		case map[string]any:
			return len(val)
		default:
			return 0
		}
	}
	
	// Add all functions to env - expr will recognize them
	// Use strContains to avoid conflict with expr's "contains" operator
	env["strContains"] = func(str, substr string) bool {
		return strings.Contains(str, substr)
	}
	env["toLowerCase"] = strings.ToLower
	env["toUpperCase"] = strings.ToUpper
	
	env["toString"] = func(v any) string {
		return fmt.Sprintf("%v", v)
	}
	
	// Compile and run the expression
	// Use AllowUndefinedVariables option to return nil for undefined variables
	program, err := expr.Compile(expression, 
		expr.Env(env),
		expr.AllowUndefinedVariables())
	if err != nil {
		return nil, fmt.Errorf("failed to compile expression: %w", err)
	}
	
	result, err := expr.Run(program, env)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate expression: %w", err)
	}
	
	return result, nil
}
