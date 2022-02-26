package nextcss

import (
	"bytes"
	"regexp"
)

const (
	BRACE_OPEN      = 123 // "{"
	BRACE_CLOSE     = 125 // "}"
	COLON           = 58  // ":"
	SEMI            = 59  // ";"
	COMMENT_SLASH   = 47  // "/"
	COMMENT_STAR    = 42  // "*"
	AT              = 64  // "@"
	DOUBLE_QUOTE    = 34  // "\""
	SINGLE_QUOTE    = 39  // "'"
	PAREN_LEFT      = 40  // "("
	PAREN_RIGHT     = 41  // ")"
	LINE_FEED       = 10  // "\n"
	CARRIAGE_RETURN = 13  // "\r"
	ESCAPE_SEQUENCE = 92  // "\"
)

type CSSParser struct {
	definitions *CSSDefinitionList
	defTree     *CSSDefinitionTree
	defRule     *CSSRule
	stack       []byte

	charPoint   int
	comment     bool
	quoting     bool
	singleQuote bool
	doubleQuote bool
	inSelector  bool
	skipping    bool
	isEscaping  bool
	inParen     bool
}

type CSSValue struct {
	Value     string `json:"data"`
	DefLine   int    `json:"line"`
	Point     int    `json:"column"`
	Semicolon bool   `json:"semicolon"`
	RawData   string `json:"raw"`
}

type CSSSelector struct {
	Selector        string `json:"selector"`
	ControlSelector bool   `json:"atrule"`
	RawData         string `json:"raw"`
	RawOffset       int    `json:"-"`
}

type CSSRule struct {
	Property  string    `json:"property"`
	Value     *CSSValue `json:"value"`
	DefLine   int       `json:"line"`
	Point     int       `json:"column"`
	RawData   string    `json:"raw"`
	RawPoint  int       `json:"-"`
	RawOffset int       `json:"-"`
}

type CSSParseResult struct {
	data []*CSSDefinition
}

type CSSDefinitionTree struct {
	definitions []*CSSDefinition
}

type CSSDefinition struct {
	Selector *CSSSelector     `json:"selector"`
	Rules    []*CSSRule       `json:"rules"`
	Controls []*CSSDefinition `json:"controls"`
	DefLine  int              `json:"line"`
	Point    int              `json:"column"`
	Parent   *CSSDefinition   `json:"-"`
}

type CSSDefinitionList struct {
	definitions []*CSSDefinition
}

func NewDefinitionList() *CSSDefinitionList {
	return &CSSDefinitionList{
		definitions: make([]*CSSDefinition, 0),
	}
}

func (l *CSSDefinitionList) Add(def *CSSDefinition) {
	l.definitions = append(l.definitions, def)
}

func (l *CSSDefinitionList) Merge(defs []*CSSDefinition) {
	l.definitions = append(l.definitions, defs...)
}
func (l *CSSDefinitionList) Get() []*CSSDefinition {
	return l.definitions
}

func NewDefinition(selector *CSSSelector, line, point int) *CSSDefinition {
	return &CSSDefinition{
		Selector: selector,
		DefLine:  line,
		Point:    point - selector.RawOffset,
		Rules:    make([]*CSSRule, 0),
		Controls: make([]*CSSDefinition, 0),
	}
}

func (d *CSSDefinition) AddRule(rule *CSSRule) {
	d.Rules = append(d.Rules, rule)
}

func (d *CSSDefinition) AddControl(control *CSSDefinition) {
	d.Controls = append(d.Controls, control)
}

func (d *CSSDefinition) AddChild(def *CSSDefinition) {
	def.Parent = d
}

func (d *CSSDefinition) GetParent() *CSSDefinition {
	return d.Parent
}

func (d *CSSDefinition) IsControl() bool {
	return d.Selector.IsControlSelector()
}

func NewDefinitionTree() *CSSDefinitionTree {
	return &CSSDefinitionTree{
		definitions: make([]*CSSDefinition, 0),
	}
}

func (l *CSSDefinitionTree) AddDefinitionToChild(def *CSSDefinition) {
	l.GetLastChild().AddChild(def)
}

func (l *CSSDefinitionTree) AddDefinition(def *CSSDefinition) {
	l.definitions = append(l.definitions, def)
}

func (l *CSSDefinitionTree) Remains() (remains bool) {
	if len(l.definitions) > 0 {
		remains = true
	}
	return
}

func (l *CSSDefinitionTree) HasParent() (has bool) {
	if len(l.definitions) > 1 {
		has = true
	}
	return
}

func (l *CSSDefinitionTree) GetLastChild() *CSSDefinition {
	return l.definitions[len(l.definitions)-1]
}

func (l *CSSDefinitionTree) Remove() {
	l.definitions = l.definitions[0 : len(l.definitions)-1]
}

func (c CSSParseResult) GetData() []*CSSDefinition {
	return c.data
}

func (c CSSParseResult) Walk(visitor func(*CSSDefinition)) {
	for _, v := range c.data {
		visitor(v)
	}

}

func (c CSSParseResult) Minisize() bytes.Buffer {
	var buffer bytes.Buffer
	for _, v := range c.data {
		s := v.Selector.Selector
		buffer.WriteString(s)
		buffer.WriteByte(BRACE_OPEN)
		for _, r := range v.Rules {
			p := r.Property
			v := r.Value.Value
			buffer.WriteString(p)
			buffer.WriteByte(COLON)
			buffer.WriteString(v)
			buffer.WriteByte(SEMI)
		}
		buffer.WriteByte(BRACE_CLOSE)

	}
	return buffer
}

func NewRule(property []byte, line, point int) *CSSRule {
	_, prop, _, offset := parseBytes(property)
	return &CSSRule{
		Property:  string(prop),
		DefLine:   line,
		Point:     point - offset,
		RawData:   string(property),
		RawPoint:  point,
		RawOffset: offset,
	}
}

func (r *CSSRule) IsSpecialProperty() (special bool) {
	if r.Property == "filter" {
		special = true
	}
	// todo
	return
}

func (r *CSSRule) SetValue(value []byte, index, point int, semicolon bool) {
	r.Value = NewValue(value, index, point, semicolon)
}

func NewSelector(selBytes []byte) *CSSSelector {
	_, selector, _, offset := parseBytes(selBytes)
	var isControl bool

	if len(selector) > 0 && selector[0] == AT {
		isControl = true
	} else {
		isControl = false
	}

	return &CSSSelector{
		Selector:        string(selector),
		ControlSelector: isControl,
		RawData:         string(selBytes),
		RawOffset:       offset,
	}
}

func (s *CSSSelector) String() string {
	return s.Selector
}

func (s *CSSSelector) IsControlSelector() bool {
	return s.ControlSelector
}

func NewValue(val []byte, line, point int, semicolon bool) *CSSValue {
	_, value, _, offset := parseBytes(val)
	return &CSSValue{
		Value:     string(value),
		DefLine:   line,
		Point:     point - offset,
		RawData:   string(val),
		Semicolon: semicolon,
	}
}

func NewParser() *CSSParser {
	return &CSSParser{
		definitions: NewDefinitionList(),
		stack:       []byte{},
		defTree:     NewDefinitionTree(),
		defRule:     nil,
		charPoint:   0,
	}
}

func (c *CSSParser) Parse(buffer []byte) CSSParseResult {
	c.execParse(buffer)

	return CSSParseResult{
		data: c.definitions.Get(),
	}
}

func (c *CSSParser) execParse(line []byte) {
	index := 1

	for point := 0; point < len(line); point++ {
		c.charPoint++

		if c.isCommentStart(line, point) {
			c.comment = true
			c.stack = append(c.stack, line[point])
			continue
		}

		if c.isCommentEnd(line, point) {
			c.comment = false
			c.stack = append(c.stack, COMMENT_STAR, COMMENT_SLASH)
			point++
			c.charPoint++
			continue
		}

		if c.comment {
			c.stack = append(c.stack, line[point])
			continue
		}

		switch line[point] {
		case ESCAPE_SEQUENCE:
			c.parseEscapeSequence()
			continue
		case PAREN_LEFT:
			if !c.quoting {
				c.inParen = true
			}
		case PAREN_RIGHT:
			if !c.quoting {
				c.inParen = false
			}
		case LINE_FEED:
			c.parseLineFeed(&index)
			index++
		case DOUBLE_QUOTE:
			if c.skipping || c.singleQuote {
				break
			}
			if c.doubleQuote {
				c.quoting = false
				c.doubleQuote = false
			} else {
				c.quoting = true
				c.doubleQuote = true
			}
		case SINGLE_QUOTE:
			if c.skipping || c.doubleQuote {
				break
			}
			if c.singleQuote {
				c.quoting = false
				c.singleQuote = false
			} else {
				c.quoting = true
				c.singleQuote = true
			}
		case BRACE_OPEN:
			if c.quoting || c.skipping || c.inParen {
				break
			}
			c.parseBraceOpen(&index)
			continue
		case COLON:
			if c.quoting || c.skipping || c.inParen {
				break
			}
			c.parseColon(&index)
			continue
		case SEMI:
			if c.quoting || c.skipping || c.inParen {
				break
			}
			c.parseSemi(&index)
			continue
		case BRACE_CLOSE:
			if c.quoting || c.skipping || c.inParen {
				break
			}
			c.parseBraceClose(&index)
			continue
		}

		if c.isEscaping {
			c.isEscaping = false
			c.skipping = false
		}

		c.stack = append(c.stack, line[point])

	}

	if c.defRule != nil {
		c.defRule.SetValue(c.stack, index, c.charPoint, false)
		c.defTree.GetLastChild().AddRule(c.defRule)
		c.defRule = nil
	}

	if c.defTree.Remains() {
		c.definitions.Add(c.defTree.GetLastChild())
	}

}

func (c *CSSParser) isCommentStart(line []byte, point int) (start bool) {
	if len(line) <= point+1 || c.quoting {
		return
	}

	if line[point] == COMMENT_SLASH && line[point+1] == COMMENT_STAR {
		start = true
	}

	return
}

func (c *CSSParser) isCommentEnd(line []byte, point int) (end bool) {
	if len(line) <= point+1 || c.quoting {
		return
	}

	if line[point] == COMMENT_STAR && line[point+1] == COMMENT_SLASH {
		end = true
	}

	return
}

func (c *CSSParser) parseEscapeSequence() {
	if !c.skipping {
		c.skipping = true
		c.isEscaping = true
	} else {
		c.skipping = false
		c.isEscaping = false
	}

	c.stack = append(c.stack, ESCAPE_SEQUENCE)
}

func (c *CSSParser) parseLineFeed(index *int) {
	val := bytes.Trim(c.stack, ";:\n\t ")

	if len(val) > 0 {
		if !c.inSelector && val[0] == AT {
			def := NewDefinition(
				NewSelector(c.stack),
				*index,
				c.charPoint,
			)
			c.definitions.Add(def)
			c.stack = []byte{}
		} else if c.defRule != nil {
			c.defRule.SetValue(
				c.stack,
				*index,
				c.charPoint,
				false,
			)
			c.defTree.GetLastChild().AddRule(c.defRule)
			c.defRule = nil
			c.stack = []byte{}
		}
	}
	c.charPoint = 0
}

func (c *CSSParser) parseBraceOpen(index *int) {
	def := NewDefinition(
		NewSelector(c.stack),
		*index,
		c.charPoint,
	)
	c.defTree.AddDefinition(def)
	c.inSelector = true
	c.stack = []byte{}
}

func (c *CSSParser) parseColon(index *int) {
	if c.defRule != nil && c.defRule.IsSpecialProperty() || !c.inSelector {
		c.stack = append(c.stack, AT)
		return
	}
	c.defRule = NewRule(c.stack, *index, c.charPoint)
	c.stack = []byte{}
}

func (c *CSSParser) parseSemi(index *int) {
	if !c.inSelector {
		def := NewDefinition(
			NewSelector(c.stack),
			*index,
			c.charPoint,
		)
		c.definitions.Add(def)
		c.stack = []byte{}
		return
	}
	if !isEmptyStack(c.stack) {
		c.defRule.SetValue(
			c.stack,
			*index,
			c.charPoint,
			true,
		)
		c.defTree.GetLastChild().AddRule(c.defRule)
		c.defRule = nil
	}
	c.stack = []byte{}
}

func (c *CSSParser) parseBraceClose(index *int) {
	cdef := c.defTree.GetLastChild()
	if c.defRule != nil {
		c.defRule.SetValue(
			c.stack,
			*index,
			c.charPoint,
			false,
		)
		cdef.AddRule(c.defRule)
		c.defRule = nil
	}
	c.defTree.Remove()

	if c.defTree.Remains() {
		c.defTree.GetLastChild().AddControl(cdef)
	} else {
		c.definitions.Add(cdef)
	}
	c.inSelector = false
	c.stack = []byte{}
}

var (
	leftRegex    = regexp.MustCompile("^([;\n\t ]*)")
	rightRegex   = regexp.MustCompile("[\n\t/\\* ]*$")
	spaceRegex   = regexp.MustCompile("[\n ]+")
	commentRegex = regexp.MustCompile("/\\*.*\\*/")
)

func parseBytes(data []byte) (before, value, after []byte, offset int) {
	size := len(data)
	data = commentRegex.ReplaceAll(data, []byte(""))
	left := leftRegex.FindSubmatchIndex(data)
	before = data[:left[1]]
	data = data[left[1]:]

	right := rightRegex.FindSubmatchIndex(data)
	value = data[:right[0]]
	after = data[right[0]:]

	offset = size - len(before)
	value = spaceRegex.ReplaceAll(value, []byte(" "))

	return
}

func isEmptyStack(stack []byte) (isEmpty bool) {
	if len(bytes.Trim(stack, "\r\n\t ")) == 0 {
		isEmpty = true
	}
	return
}
