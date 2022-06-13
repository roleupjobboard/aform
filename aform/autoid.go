package aform

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	defaultAutoID  = "id_%s"
	disabledAutoID = ""
)

// validAutoID checks if the provided auto ID is correct.
// Correct values are:
// 		- empty string to deactivate auto ID e.g. ""
// 		- string containing one verb '%s' to customize auto ID. e.g. "id_%s"
//
// A non-empty string without a verb '%s' is invalid.
func validAutoID(autoID string) error {
	if isDisabledAutoID(autoID) {
		return nil
	}
	switch strings.Count(autoID, "%s") {
	case 1:
		return nil
	default:
		return fmt.Errorf("autoID must contain one %%s verb (e.g. %s). To disable auto ID use DisableAutoID(). Given: %s", defaultAutoID, autoID)
	}
}

func hasID(fld fieldReader) bool {
	return !isDisabledAutoID(fld.AutoID())
}

// isDisabledAutoID returns true if autoID is the disabled Auto ID
func isDisabledAutoID(autoID string) bool {
	return autoID == disabledAutoID
}

func normalizedIDForField(fld fieldReader) string {
	if !hasID(fld) {
		return disabledAutoID
	}
	singleSpacePattern := regexp.MustCompile(`\s+`)
	return fmt.Sprintf(fld.AutoID(), strings.ToLower(singleSpacePattern.ReplaceAllString(fld.Name(), "_")))
}

func normalizedNameForField(fld fieldReader) string {
	return normalizedName(fld.Name())
}

func normalizedName(name string) string {
	singleSpacePattern := regexp.MustCompile(`\s+`)
	return strings.ToLower(singleSpacePattern.ReplaceAllString(name, "_"))
}

func normalizedDescribedByIDErrList(fld fieldReader, numberOfError int) []string {
	ids := make([]string, numberOfError)
	for i := 0; i < numberOfError; i++ {
		ids[i] = normalizedDescribedByIDForErr(fld, i)
	}
	return ids
}

// normalizedDescribedByIDForErr generates IDs for errors pointed by the aria-describedby tag.
func normalizedDescribedByIDForErr(fld fieldReader, index int) string {
	return fmt.Sprintf("err_%d_", index) + normalizedIDForField(fld)
}

func normalizedGroupID(normalizedID string, index uint) string {
	return normalizedID + fmt.Sprintf("_%d", index)
}

func normalizedSubGroupID(normalizedID string, index, subIndex uint) string {
	return normalizedID + fmt.Sprintf("_%d_%d", index, subIndex)
}

// normalizedDescribedByIDForHelpText generates the ID for the help text pointed by the aria-describedby tag.
func normalizedDescribedByIDForHelpText(fld fieldReader) string {
	return "helptext_" + normalizedIDForField(fld)
}
