package util

import (
	"fmt"
	"strconv"
	"strings"

	projectapi "github.com/sh-miyoshi/hekate/pkg/apihandler/admin/v1/project"
	"github.com/stretchr/stew/slice"
)

// ParsePolicies ...
func ParsePolicies(policies []string) (projectapi.PasswordPolicy, error) {
	res := projectapi.PasswordPolicy{}
	for _, policy := range policies {
		val := strings.Split(policy, "=")
		switch val[0] {
		case "minLen":
			if len(val) != 2 {
				return res, fmt.Errorf("please set unsigned integer for minLen")
			}
			v, err := strconv.ParseUint(val[1], 10, 64)
			if err != nil {
				return res, fmt.Errorf("please set unsigned integer for minLen")
			}
			res.MinimumLength = uint(v)
		case "notUserName":
			if len(val) == 1 {
				res.NotUserName = true
				continue
			}
			v, err := strconv.ParseBool(val[1])
			if err != nil {
				return res, fmt.Errorf("plase set true or false for notUserName")
			}
			res.NotUserName = v
		case "useChar":
			if len(val) != 2 {
				return res, fmt.Errorf("please set lower or upper or both or either for useChar")
			}
			valid := []string{"lower", "upper", "both", "either"}
			if !slice.Contains(valid, val[1]) {
				return res, fmt.Errorf("please set lower or upper or both or either for useChar")
			}
			res.UseCharacter = val[1]
		case "useDigit":
			if len(val) == 1 {
				res.UseDigit = true
				continue
			}
			v, err := strconv.ParseBool(val[1])
			if err != nil {
				return res, fmt.Errorf("plase set true or false for useDigit")
			}
			res.UseDigit = v
		case "useSpecialChar":
			if len(val) == 1 {
				res.UseSpecialCharacter = true
				continue
			}
			v, err := strconv.ParseBool(val[1])
			if err != nil {
				return res, fmt.Errorf("plase set true or false for useSpecialChar")
			}
			res.UseSpecialCharacter = v
		case "blackLists":
			if len(val) != 2 {
				return res, fmt.Errorf("please set string separated by semicolon for blackLists")
			}
			res.BlackList = strings.Split(val[1], ";")
		default:
			return res, fmt.Errorf("Invalid policy type %s", policy)
		}
	}
	return res, nil
}
