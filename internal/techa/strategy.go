package techa

import "errors"

var AVAILABLE_INDICATORS = []string{"SMA", "EMA", "MACD", "BollingerBands", "RSI", "StochRSI"}

type StrategyNodeVariable struct {
	result    float64
	indicator string
}

type StrategyNode struct {
	head       *StrategyNode `json:"head"`
	tail       *StrategyNode `json:"tail"`
	conditions []Condition   `json:"conditions"`
}

type Condition struct {
	operator      string               `json:"operator"`
	alphaVariable StrategyNodeVariable `json:"alpha_variable"`
	betaVariable  StrategyNodeVariable `json:"beta_variable"`
}

func (c *Condition) ValidateOperator(operator string) bool {
	return operator == ">" || operator == ">=" || operator == "=" || operator == "<=" || operator == "<"
}

func (c *Condition) SetOperator(operator string) error {
	if !c.ValidateOperator(operator) {
		return errors.New("incorrect operator format")
	}
	c.operator = operator
	return nil
}
func validateVariableClasses(class string) bool {
	return class == "constant" || class == "indicator"
}
func NewStrategyNodeVariable()

func NewCondition(operator string, alpha string)

type StrategyTree struct {
}
