package arnlike

import (
	"fmt"
	"strings"

	"github.com/gobwas/glob"
)

const (
	arnDelimiter = ":"
	arnSectionsExpected  = 6
	arnPrefix    = "arn:"

	// zero-indexed
	sectionPartition = 1
	sectionService   = 2
	sectionRegion    = 3
	sectionAccountID = 4
	sectionResource  = 5

	// errors
	invalidPrefix   = "invalid prefix"
	invalidSections = "not enough sections"
)

// ArnLike takes an ARN returns true if it is matched by the pattern.
// Each component of the ARN is matched individually as per
// https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_condition_operators.html#Conditions_ARN
func ArnLike(arn, pattern string) (bool, error) {
	// "parse" the input arn into sections
	arnSections, err := parse(arn)
	if err != nil {
		return false, fmt.Errorf("Could not parse input arn: %v", err)
	}
	patternSections, err := parse(pattern)
	if err != nil {
		return false, fmt.Errorf("Could not parse ArnLike string: %v", err)
	}

	// Tidy glob special characters. Escape the ones not used in ArnLike.
	// Replace multiple ** with a single one
	tidyPatternSections(patternSections)
	fmt.Println(patternSections)

	for index := range arnSections {
		patternGlob, err := glob.Compile(patternSections[index])
		if err != nil {
			return false, fmt.Errorf("Could not parse %s: %v", patternSections[index], err)
		}

		if !patternGlob.Match(arnSections[index]) {
			return false, nil
		}
	}

	return true, nil
}

// parse is a copy of arn.Parse from the AWS SDK but represents the ARN as []string
func parse(input string) ([]string, error) {
	if !strings.HasPrefix(input, arnPrefix) {
		return nil, fmt.Errorf(invalidPrefix)
	}
	arnSections := strings.SplitN(input, arnDelimiter, arnSectionsExpected)
	if len(arnSections) != arnSectionsExpected {
		return nil, fmt.Errorf(invalidSections)
	}

	return arnSections, nil
}

// tidyPatternSections goes through each section of the arnLike slice and escapes any meta characters
// used by glob but not used by ArnLike
func tidyPatternSections(arnLikeSlice []string) {
	for index, section := range arnLikeSlice {
		b := make([]byte, 2*len(section))

		// a byte loop is correct because all meta characters are ASCII
		j := 0
		for i := 0; i < len(section); i++ {
			switch section[i] {
			case '{', '}', '[', ']', '\\':
				b[j] = '\\'
				j++
			}
			b[j] = section[i]
			j++
		}

		arnLikeSlice[index] = string(b)
	}
}
