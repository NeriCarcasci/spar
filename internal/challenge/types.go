package challenge

type Difficulty string

const (
	Easy   Difficulty = "easy"
	Medium Difficulty = "medium"
	Hard   Difficulty = "hard"
)

type Challenge struct {
	ID               string     `yaml:"id"`
	Title            string     `yaml:"title"`
	Collection       string     `yaml:"collection"`
	Difficulty       Difficulty `yaml:"difficulty"`
	Category         string     `yaml:"category"`
	Tags             []string   `yaml:"tags"`
	Languages        []string   `yaml:"languages"`
	InputType        string     `yaml:"input_type"`
	OutputType       string     `yaml:"output_type"`
	Description      string     `yaml:"description"`
	Examples         []Example  `yaml:"examples"`
	Constraints      []string   `yaml:"constraints"`
	TimeLimitSeconds int        `yaml:"time_limit_seconds"`
	Hints            []string   `yaml:"hints"`
	Path             string     `yaml:"-"`
}

type Example struct {
	Input       string `yaml:"input"`
	Output      string `yaml:"output"`
	Explanation string `yaml:"explanation,omitempty"`
}

type TestSuite struct {
	Tests []TestCase `yaml:"tests"`
}

type TestCase struct {
	Input        map[string]interface{} `yaml:"input,omitempty"`
	InputFile    string                 `yaml:"input_file,omitempty"`
	Expected     interface{}            `yaml:"expected,omitempty"`
	ExpectedFile string                 `yaml:"expected_file,omitempty"`
	Visible      bool                   `yaml:"visible"`
}

type IndexEntry struct {
	ID         string     `yaml:"id"`
	Title      string     `yaml:"title"`
	Difficulty Difficulty `yaml:"difficulty"`
	Category   string     `yaml:"category"`
	Tags       []string   `yaml:"tags"`
	Languages  []string   `yaml:"languages"`
	Path       string     `yaml:"path"`
}

type Index struct {
	Challenges []IndexEntry `yaml:"challenges"`
}
