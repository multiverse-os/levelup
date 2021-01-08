package model

type Action func() error

type Hooks struct {
	BeforeGet []*Action
	AfterGet  []*Action

	BeforeSet []*Action
	AfterSet  []*Action

	BeforeDelete []*Action
	AfterDelete  []*Action

	AfterOpen   []*Action
	BeforeClose []*Action
}
