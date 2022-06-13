package aform

import (
	"golang.org/x/exp/slices"
)

func fieldGroupsToWidgetGroups(groups []choiceFieldOptionGroup, optionWidget Widget, baseID, htmlName string, selected []string) []map[string][]widgetOption {
	output := make([]map[string][]widgetOption, 0, len(groups))
	groupIndex := uint(0)
	for _, group := range groups {
		groupWidget := group.widget(optionWidget, baseID, htmlName, groupIndex, selected)
		groupIndex += indexesUsedByGroup(groupWidget)
		output = append(output, groupWidget)
	}
	return output
}

func indexesUsedByGroup(group map[string][]widgetOption) uint {
	output := uint(0)
	for name, options := range group {
		if len(name) > 0 {
			output += 1
		} else {
			output += uint(len(options))
		}
	}
	return output
}

type choiceFieldOptionGroup map[string][]ChoiceFieldOption

type ChoiceFieldOption struct {
	Value string
	Label string
}

func (g choiceFieldOptionGroup) values() []string {
	var values []string
	for _, options := range g {
		for _, option := range options {
			values = append(values, option.Value)
		}
	}
	return values
}

func (g choiceFieldOptionGroup) widget(inputWidget Widget, baseID, htmlName string, groupIndex uint, selected []string) map[string][]widgetOption {
	output := map[string][]widgetOption{}
	for groupName, groupOptions := range g {
		widgetOptions := make([]widgetOption, len(groupOptions))
		withID := len(baseID) > 0
		isSubGroup := len(groupName) > 0
		for i, option := range groupOptions {
			attrs := map[string]string{}
			if withID {
				if isSubGroup {
					attrs["id"] = normalizedSubGroupID(baseID, groupIndex, uint(i))
				} else {
					attrs["id"] = normalizedGroupID(baseID, groupIndex+uint(i))
				}
			}
			widgetOptions[i] = option.widget(inputWidget, htmlName, slices.Contains(selected, option.Value), attrs)
		}
		output[groupName] = widgetOptions
	}
	return output
}

func (o ChoiceFieldOption) widget(inputWidget Widget, htmlName string, selected bool, attrs map[string]string) widgetOption {
	if selectedAttr, ok := inputWidget.selectedAttr(selected); ok {
		attrs[selectedAttr.n] = selectedAttr.v
	}
	return widgetOption{
		Label:       o.Label,
		WrapLabel:   true,
		widgetInput: widgetInput{Type: inputWidget, Name: htmlName, Value: o.Value, Attrs: attrs},
	}
}
