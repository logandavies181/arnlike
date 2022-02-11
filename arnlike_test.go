package arnlike

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type arnLikeInput struct {
	arn, pattern string
}

type quoteMetaInput struct {
	input, expected string
}

func TestArnLikePostiveMatches(t *testing.T) {
	inputs := []arnLikeInput{
		{
			arn:     `arn:aws:iam::000000000000:role/some-role`,
			pattern: `arn:aws:iam::000000000000:role/some-role`,
		},
		{
			arn:     `arn:aws:iam::000000000000:role/some-role`,
			pattern: `arn:aws:iam::000000000000:*`,
		},
		{
			arn:     `arn:aws:iam::000000000000:role/some-role`,
			pattern: `arn:*:*:*:*:*`,
		},
		{
			arn:     `arn:aws:iam::000000000000:role/some-role`,
			pattern: `arn:aws:iam::000000000000:**`,
		},
		{
			arn:     `arn:aws:iam::000000000000:role/some-role`,
			pattern: `arn:aws:iam::000000000000:*role*`,
		},
		{
			arn:     `arn:aws:iam::000000000000:role/some-role`,
			pattern: `arn:aws:iam::000000000000:ro*`,
		},
		{
			arn:     `arn:aws:iam::000000000000:role/some-role`,
			pattern: `arn:aws:iam::000000000000:??????????????`,
		},
		{
			arn:     `arn:aws:testservice::000000000000:some/wacky-new-[resource]{with}\metacharacters`,
			pattern: `arn:aws:testservice::000000000000:some/wacky-new-[resource]{with}\metacharacters`,
		},
		{
			arn:     `arn:aws:testservice::000000000000:some/wacky-new-[resource]{with}\metacharacters`,
			pattern: `arn:aws:testservice::000000000000:some/wacky-new-[reso*`,
		},
	}

	for _, v := range inputs {
		output, err := ArnLike(v.arn, v.pattern)
		assert.Nil(t, err, "Expected no error for input arn: %s pattern: %s", v.arn, v.pattern)

		assert.True(t, output, "Expected true for input arn: %s pattern: %s", v.arn, v.pattern)
	}
}

func TestArnLikeNetagiveMatches(t *testing.T) {
	inputs := []arnLikeInput{
		{
			arn:     `arn:aws:iam::111111111111:role/some-role`,
			pattern: `arn:aws:iam::000000000000:role/some-role`,
		},
		{
			arn:     `arn:aws:testservice::000000000000:some/wacky:resource:with:colon:delims`,
			pattern: `arn:aws:testservice::**:delims`,
		},
	}

	for _, v := range inputs {
		output, err := ArnLike(v.arn, v.pattern)
		assert.Nil(t, err, "Expected no error for input arn: %s pattern: %s", v.arn, v.pattern)

		assert.False(t, output, "Expected false for input arn: %s pattern: %s", v.arn, v.pattern)
	}
}

func TestArnLikeInvalidArns(t *testing.T) {
	invalidPrefixArn := `nar:aws:iam::000000000000:role/some-role`
	invalidSectionsArn := `arn:aws:iam:000000000000:role/some-role`
	validArn := `arn:aws:iam::000000000000:role/some-role`

	// invalid prefix
	output, err := ArnLike(invalidPrefixArn, validArn)
	assert.False(t, output,
		"Expected false result on error for input arn: %s, pattern: %s", invalidPrefixArn, validArn)

	assert.Equal(t, "Could not parse input arn: invalid prefix", err.Error())

	// invalid sections
	output, err = ArnLike(invalidSectionsArn, validArn)
	assert.False(t, output,
		"Expected false result on error for input arn: %s, pattern: %s", invalidSectionsArn, validArn)

	assert.Equal(t, "Could not parse input arn: not enough sections", err.Error())
}

func TestQuoteMeta(t *testing.T) {
	inputs := []quoteMetaInput{
		{
			input:    `**`,
			expected: `.*.*`,
		},
		{
			input:    `??`,
			expected: `.?.?`,
		},
		{
			input:    `abdcEFG`,
			expected: `abdcEFG`,
		},
		{
			input:    `abd.EFG`,
			expected: `abd\.EFG`,
		},
		{
			input:    `\.+()|[]{}^$`,
			expected: `\\\.\+\(\)\|\[\]\{\}\^\$`,
		},
		{
			input:    `\.+()|[]{}^$*?`,
			expected: `\\\.\+\(\)\|\[\]\{\}\^\$.*.?`,
		},
	}

	for _, v := range inputs {
		output := quoteMeta(v.input)
		assert.Equal(t, v.expected, output)
	}
}
