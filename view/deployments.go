package view

import (
	"regexp"

	"github.com/hashicorp/nomad/api"
	"github.com/rivo/tview"

	"github.com/hcjulz/damon/component"
	"github.com/hcjulz/damon/models"
)

func (v *View) Deployments() {
	v.viewSwitch()

	v.components.Commands.Update(component.DeploymentCommands)

	v.Layout.Container.SetInputCapture(v.InputDeployments)

	v.state.Elements.TableMain = v.components.DeploymentTable.Table.Primitive().(*tview.Table)

	update := func() {
		v.components.DeploymentTable.Props.Data = v.filterDeployments(v.state.Deployments)
		v.components.DeploymentTable.Props.Namespace = v.state.SelectedNamespace
		v.components.DeploymentTable.Render()
		v.Draw()
	}

	v.components.Search.InputField.SetText("")
	v.components.Search.Props.ChangedFunc = func(text string) {
		v.state.Filter.Deployments = text
		update()
	}

	v.Watcher.Subscribe(update, api.TopicDeployment)

	update()

	v.components.Selections.Namespace.SetSelectedFunc(func(text string, index int) {
		v.state.SelectedNamespace = text
		v.Deployments()
	})

	v.addToHistory(v.state.SelectedNamespace, api.TopicDeployment, v.Deployments)
	v.Layout.Container.SetFocus(v.components.DeploymentTable.Table.Primitive())
}

func (v *View) filterDeployments(data []*models.Deployment) []*models.Deployment {
	filter := v.state.Filter.Deployments
	if filter != "" {
		rx, _ := regexp.Compile(filter)
		result := []*models.Deployment{}
		for _, dep := range v.state.Deployments {
			switch true {
			case rx.MatchString(dep.ID),
				rx.MatchString(dep.JobID),
				rx.MatchString(dep.Namespace),
				rx.MatchString(dep.Status),
				rx.MatchString(dep.StatusDescription):
				result = append(result, dep)
			}
		}

		return result
	}

	return data
}
