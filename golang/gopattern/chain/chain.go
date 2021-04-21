package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
)

// PipelineConfig is the representation of a pipeline in the configuration.
type PipelineConfig struct {
	// Steps is the list of step in your pipeline.
	Steps map[string]PipelineStep `yaml:"steps"`

	// Root is the name of the first step in your pipeline, we will start by calling it, and it will call the next steps after it.
	Root string `yaml:"root"`
}

// PipelineStep is a step representation in the configuration.
type PipelineStep struct {
	// StepType is the type of the Handler to map for this step configuration, the list of
	// the available types is in the method getHandlerFromType
	StepType string `yaml:"type"`

	// Next is the next step we should call after this one. This param is not mandatory.
	Next string `yaml:"next"`
}

// NewPipeline will create a new Pipeline ready to be executed
func NewPipeline(pipelineConfigUrl string) (Pipeline, error) {
	resp, _ := http.Get(pipelineConfigUrl)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var pipelineConfig PipelineConfig
	_ = yaml.Unmarshal(body, &pipelineConfig)

	p := Pipeline{
		steps: pipelineConfig.Steps,
		root:  pipelineConfig.Root,
	}

	// Get handlers from  config
	p.handlers = make(map[string]Handler, len(p.steps))
	for name, step := range p.steps {
		handler, _ := p.getHandlerFromType(step.StepType)
		p.handlers[name] = handler
	}

	// Init all handlers
	for name, step := range p.steps {
		err := p.handlers[name].Init(name, step, p.handlers)
		if err != nil {
			return Pipeline{}, fmt.Errorf("impossible to init the step named '%s': %v", name, err)
		}
	}
	// Check that root step exists
	if _, ok := p.handlers[p.root]; !ok {
		return Pipeline{}, fmt.Errorf("impossible to start with step \"%s\" because it does not exists", p.root)
	}
	return p, nil
}

type Pipeline struct {
	root     string
	steps    map[string]PipelineStep
	handlers map[string]Handler
}

func (p *Pipeline) start() {
	context := make([]string, 0)
	err := p.Execute(&context)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Check what all steps have done
	fmt.Println(context)
}

// getHandlerFromType is mapping handler type name in your configuration to proper handlers.
func (p *Pipeline) getHandlerFromType(s string) (Handler, error) {
	// mapping list for the handlers
	handlers := map[string]Handler{
		"handlerImpl1": &HandlerImpl1{},
		"handlerImpl2": &HandlerImpl2{},
	}

	stepHandler, handlerExists := handlers[s]
	if !handlerExists {
		return nil, fmt.Errorf("impossible to find a matching step handler for %s", s)
	}
	return stepHandler, nil
}

// Execute the search Pipeline by taking the 1st Step and execute it.
func (p *Pipeline) Execute(context *[]string) error {
	return p.handlers[p.root].Execute(context)
}

// Handler is defining how a step looks like.
type Handler interface {
	// Init configure the step from the configuration file
	Init(name string, step PipelineStep, availableHandlers map[string]Handler) error

	// Execute apply the action of the step and move to the next step
	Execute(context *[]string) error
}

type HandlerImpl1 struct {
	next Handler
}

func (e *HandlerImpl1) Init(name string, step PipelineStep, availableHandlers map[string]Handler) error {
	// This is a simplified version of the init method, you can check that next step is not it-self
	// and that the handler is available.
	if step.Next != "" {
		e.next = availableHandlers[step.Next]
	}
	return nil
}

func (e *HandlerImpl1) Execute(context *[]string) error {
	// You can add logic before and after the next step is called.
	*context = append(*context, "HandlerImpl1: before the call")
	if e.next != nil {
		_ = e.next.Execute(context)
	}
	*context = append(*context, "HandlerImpl1: after the call")
	return nil
}

type HandlerImpl2 struct {
	next Handler
}

func (e *HandlerImpl2) Init(name string, step PipelineStep, availableHandlers map[string]Handler) error {
	if step.Next != "" {
		e.next = availableHandlers[step.Next]
	}
	return nil
}

func (e *HandlerImpl2) Execute(context *[]string) error {
	*context = append(*context, "HandlerImpl2 called")
	if e.next != nil {
		return e.next.Execute(context)
	}
	return nil
}
