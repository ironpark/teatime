package node

// Property Builder Pattern Usage Examples
//
// The property builder provides a fluent interface for creating NodeProperty instances.
//
// Basic usage:
//
//	prop := StringProp("key", "Label", Optional())
//	prop := IntProp("count", "Count", RequiredWithDefault(10))
//	prop := BoolProp("enabled", "Enabled", Toggle())
//
// With ranges:
//
//	volume := IntProp("volume", "Volume", RangeSlider(0, 100))
//	opacity := FloatProp("opacity", "Opacity", RangeSlider(0, 1))
//	percent := IntProp("percent", "Percentage", Percentage())
//
// With selections:
//
//	size := SelectProp("size", "Size", []string{"S", "M", "L"})
//	tags := MultiSelectProp("tags", "Tags", []string{"tag1", "tag2"})
//
// Complex properties:
//
//	email := StringProp("email", "Email",
//		WithDescription("User email address"),
//		WithPlaceholder("user@example.com"),
//		RequiredWithDefault(""),
//		WithMinLength(5),
//		WithMaxLength(100),
//	)
//
// Text areas:
//
//	desc := StringProp("description", "Description",
//		TextArea(10),  // 10 rows
//		Optional(),
//	)
//
// Output properties (read-only):
//
//	result := OutputProp(String, "result", "Result")
//