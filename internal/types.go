// 测试用例优先级
type Priority int

const (
	PriorityLow Priority = iota
	PriorityNormal
	PriorityHigh
)

// 测试用例结构体
type TestCase struct {
	Name           string                 `mapstructure:"name"`
	Url            string                 `mapstructure:"url"`
	Method         string                 `mapstructure:"method"`
	BodyTemplate   map[string]interface{} `mapstructure:"body_template"`
	ParamsTemplate map[string]string      `mapstructure:"params_template"`
	Parameters     []TestParameter        `mapstructure:"parameters"`
	DefaultStatus  int                    `mapstructure:"expected_status"`
	DefaultBody    string                 `mapstructure:"expected_body"`
	PreRequest     *PreRequestConfig      `mapstructure:"pre_request"`
	SourceFile     string                 `mapstructure:"-"`
	Project        string                 `mapstructure:"project"`
	Priority       Priority               `mapstructure:"priority"`
	Tags           []string              `mapstructure:"tags"` // 测试用例标签
}
