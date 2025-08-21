package node

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/samber/lo"
)

// ResolveInput resolves dynamic input properties for a node based on workflow state.
// It processes three types of property values:
//   - Bindings: @[bindingName] - replaced with values from workflow state
//   - Expressions: {{expression}} - evaluated using expr-lang with state context
//   - Plain values: returned as-is
//
// Returns a map of resolved property values ready for node execution.
func ResolveInput(properties []NodeProperty, states WorkflowState) (resolvedProperties map[string]any, err error) {
	re := regexp.MustCompile(`{{.*}}`)
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
				resolvedProperties[key], err = value.Cast(bindingValue)
				if err != nil {
					return nil, fmt.Errorf("failed to cast binding value for property %s:%s (%w)", key, binding, err)
				}
			} else if strings.HasPrefix(v, "{{") && strings.HasSuffix(v, "}}") {
				// it is a expression, evaluate it
				expression := v[2 : len(v)-2]
				resolvedProperties[key], err = Eval(expression, states)
				if err != nil {
					return nil, fmt.Errorf("failed to evaluate expression for property %s:%s (%w)", key, expression, err)
				}
			} else {
				// it is a value or embedded multiple expressions
				if re.MatchString(v) {
					expressions := re.FindAllString(v, -1)
					for _, expression := range expressions {
						evaluated, err := Eval(expression[2:len(expression)-2], states)
						if err != nil {
							return nil, fmt.Errorf("failed to evaluate expression for property %s:%s (%w)", key, expression, err)
						}
						// replace the expression with the evaluated value
						v = strings.Replace(v, expression, fmt.Sprintf("%v", evaluated), 1)
					}
					resolvedProperties[key] = v
				} else {
					resolvedProperties[key], err = value.Cast(v)
					if err != nil {
						return nil, fmt.Errorf("failed to cast value for property %s:%s (%w)", key, v, err)
					}
				}
			}
		default:
			resolvedProperties[key], err = value.Cast(v)
			if err != nil {
				return nil, fmt.Errorf("failed to cast value for property %s:%s (%w)", key, v, err)
			}
		}
	}
	return resolvedProperties, nil
}

// Eval evaluates a JavaScript-like expression using the expr-lang library.
// It creates an execution environment with workflow state data and built-in functions.
//
// States are structured as:
//   - {nodeId}.input.{propertyKey} = input values
//   - {nodeId}.output.{outputKey} = output values
//   - Direct key-value pairs for other workflow data
//
// Built-in functions include: len(), strContains(), toLowerCase(), toUpperCase(), toString()
func Eval(expression string, states WorkflowState) (any, error) {
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
