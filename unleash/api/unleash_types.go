package api

import "time"

type Response struct {
	Version int `json:"version"`
}

type FeatureResponse struct {
	Response
	Features []Feature `json:"features"`
}

type Feature struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Enabled     bool         `json:"enabled"`
	Strategies  []Strategy   `json:"strategies"`
	CreatedAt   time.Time    `json:"createdAt"`
	Strategy    string       `json:"strategy"`
	Parameters  ParameterMap `json:"parameters"`
}

type ParameterMap map[string]interface{}

func (fr FeatureResponse) FeatureMap() map[string]interface{} {
	features := map[string]interface{}{}
	for _, f := range fr.Features {
		features[f.Name] = f
	}
	return features
}

type StrategyResponse struct {
	Response
	Strategies []StrategyDescription `json:"strategies"`
}

type Strategy struct {
	Id          int          `json:"id"`
	Name        string       `json:"name"`
	Constraints []Constraint `json:"constraints"`
	Parameters  ParameterMap `json:"parameters"`
}

type ParameterDescription struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

type StrategyDescription struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  []ParameterDescription `json:"parameters"`
}

// Operator is a type representing a constraint operator
type Operator string

const (
	// OperatorIn indicates that the context values must be
	// contained within those specified in the constraint.
	OperatorIn Operator = "IN"

	// OperatorNotIn indicates that the context values must
	// NOT be contained within those specified in the constraint.
	OperatorNotIn Operator = "NOT_IN"
)

// Constraint represents a constraint on a particular context
// value.
type Constraint struct {
	ContextName string   `json:"contextName"`
	Operator    Operator `json:"operator"`
	Values      []string `json:"values"`
}
