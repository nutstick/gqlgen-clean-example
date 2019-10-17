package gqlapolloenginetracing

import (
	"regexp"

	"github.com/99designs/gqlgen/graphql"
)

func printWithReducedWhitespace(requestContext *graphql.RequestContext) string {
	r := regexp.MustCompile(`\s+`)
	r2 := regexp.MustCompile(`([^_a-zA-Z0-9]) `)
	r3 := regexp.MustCompile(` ([^_a-zA-Z0-9])`)
	r4 := regexp.MustCompile(`"([a-f0-9]+)"`)
	reducedWithSpace := r4.ReplaceAllString(
		r3.ReplaceAllString(
			r2.ReplaceAllString(
				r.ReplaceAllString(requestContext.RawQuery, " "),
				"$1",
			),
			"$1",
		),
		"$1",
	)
	return reducedWithSpace
}
