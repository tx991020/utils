package xiris

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"

	"github.com/kataras/iris/v12"
)

const MaxUint = ^uint64(0)
const MaxInt = int64(MaxUint >> 1)

func GRequestBody(c iris.Context) {
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		c.Text("Parse body failed.")
		c.StopExecution()
	} else {

		c.Values().Set("requestBody", b)
		c.Next()
	}
	return
}

func GRequestBodyMap(c iris.Context) {
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		c.Text("Parse body failed.")
		c.StopExecution()

	} else {
		m := map[string]interface{}{}
		if err = json.Unmarshal(b, &m); err != nil {
			c.Text(fmt.Sprintf("Json unmarshal failed: %s, %v, %v", string(b), m, err))
			c.StopExecution()
		} else {
			c.Values().Set("requestBody", m)
			c.Next()
		}
	}
	return
}

// idl = {json, xml}
func GRequestBodyObject(t reflect.Type, idl string) func(iris.Context) {
	return func(c iris.Context) {
		instance := reflect.New(t).Interface()

		b, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			c.Text("Parse body failed.")
			c.StopExecution()
		} else {
			switch idl {
			case "json":
				if err = json.Unmarshal(b, instance); err != nil {
					c.Text(fmt.Sprintf("Json unmarshal failed: %s, %v, %v", string(b), instance, err))
					c.StopExecution()
				} else {
					c.Values().Set("requestBody", instance)
					c.Next()
				}
			case "xml":
			default:
			}
		}
		return
	}
}

func GPathRequireInt(match string) func(iris.Context) { return GPathInt(match, true, "") }
func GPathRequireIntAlias(match string, alias string) func(iris.Context) {
	return GPathInt(match, true, alias)
}

func GPathOptionalInt(match string) func(iris.Context) { return GPathInt(match, false, "") }
func GPathOptionalIntAlias(match string, alias string) func(iris.Context) {
	return GPathInt(match, false, alias)
}

func GPathInt(match string, must bool, alias string) func(iris.Context) {
	return func(c iris.Context) {
		c.Params()
		if n, err := strconv.ParseInt(c.Params().Get(match), 10, 64); err != nil {
			if must {
				c.Text(fmt.Sprintf("Parse path %s failed.", match))
				c.StopExecution()
			}
		} else {
			c.Values().Set(match, n)
			c.Next()
			if len(alias) > 0 {
				c.Values().Set(alias, n)
				c.Next()
			}
		}
	}
}

func GPathRequireString(match string) func(iris.Context) {
	return GPathString(match, true, "")
}
func GPathRequireStringAlias(match string, must bool, alias string) func(iris.Context) {
	return GPathString(match, true, alias)
}

func GPathOptionalString(match string) func(iris.Context) {
	return GPathString(match, false, "")
}
func GPathOptionalStringAlias(match string, alias string) func(iris.Context) {
	return GPathString(match, false, alias)
}

func GPathString(match string, must bool, alias string) func(iris.Context) {
	return func(c iris.Context) {
		if len(c.Params().Get(match)) == 0 {
			if must {

				c.Text(fmt.Sprintf("Parse path param: %s failed.", match))
				c.StopExecution()
			}
		} else {
			c.Values().Set(match, c.Params().Get(match))
			c.Next()
			if len(alias) > 0 {
				c.Values().Set(alias, c.Params().Get(match))
				c.Next()
			}
		}
	}
}

func GHeaderRequireInt(match string) func(iris.Context) { return GHeaderInt(match, true, "") }
func GHeaderRequireIntAlias(match string, alias string) func(iris.Context) {
	return GHeaderInt(match, true, alias)
}

func GHeaderOptionalInt(match string) func(iris.Context) { return GHeaderInt(match, false, "") }
func GHeaderOptionalIntAlias(match string, alias string) func(iris.Context) {
	return GHeaderInt(match, false, alias)
}

func GHeaderInt(match string, must bool, alias string) func(iris.Context) {
	return func(c iris.Context) {
		if len(c.GetHeader(match)) > 0 {
			if id, err := strconv.ParseInt(c.GetHeader(match), 10, 64); err != nil {
				if must {
					c.Text(fmt.Sprintf("Parse header: %s failed.", match))
					c.StopExecution()
				}
			} else {
				c.Values().Set(match, id)
				c.Next()
				if len(alias) > 0 {
					c.Values().Set(alias, id)
					c.Next()
				}
			}
		} else if must {
			c.Text(fmt.Sprintf("Parse header: %s failed.", match))
			c.StopExecution()
		}
	}
}

func GHeaderRequireString(match string) func(iris.Context) {
	return GHeaderString(match, true, "")
}
func GHeaderRequireStringAlias(match string, alias string) func(iris.Context) {
	return GHeaderString(match, true, alias)
}

func GHeaderOptionalString(match string) func(iris.Context) {
	return GHeaderString(match, false, "")
}
func GHeaderOptionalStringAlias(match string, alias string) func(iris.Context) {
	return GHeaderString(match, false, alias)
}

func GHeaderString(match string, must bool, alias string) func(iris.Context) {
	return func(c iris.Context) {
		if len(c.GetHeader(match)) == 0 {
			if must {
				c.Text(fmt.Sprintf("Parse header: %s failed.", match))
				c.StopExecution()
			}
		} else {

			c.Values().Set(match, c.GetHeader(match))
			c.Next()
			if len(alias) > 0 {
				c.Values().Set(alias, c.GetHeader(match))
				c.Next()

			}
		}
	}
}

func GHeaderOptionalStringDefault(match string, defaultValue string) func(iris.Context) {
	return GHeaderStringDefault(match, false, "", defaultValue)
}

func GHeaderStringDefault(match string, must bool, alias, defaultValue string) func(iris.Context) {
	return func(c iris.Context) {
		if len(c.GetHeader(match)) == 0 {
			if must {
				c.Text(fmt.Sprintf("Parse header: %s failed.", match))
				c.StopExecution()
			}
			c.Values().Set(match, defaultValue)
			c.Next()
			if len(alias) > 0 {
				c.Values().Set(alias, defaultValue)
				c.Next()
			}
		} else {
			c.Values().Set(match, c.GetHeader(match))
			c.Next()
			if len(alias) > 0 {
				c.Values().Set(alias, c.GetHeader(match))
				c.Next()
			}
		}
	}
}

func GQueryRequirePositiveInt(match string) func(iris.Context) {
	return GQueryPositiveInt(match, true, "")
}
func GQueryRequirePositiveIntAlias(match string, alias string) func(iris.Context) {
	return GQueryPositiveInt(match, true, alias)
}

func GQueryOptionalPositiveInt(match string) func(iris.Context) {
	return GQueryPositiveInt(match, false, "")
}
func GQueryOptionalPositiveIntAlias(match string, alias string) func(iris.Context) {
	return GQueryPositiveInt(match, false, alias)
}

func GQueryPositiveInt(match string, must bool, alias string) func(iris.Context) {
	return func(c iris.Context) {
		q := c.Request().URL.Query()
		if q[match] == nil || len(q[match][0]) == 0 {
			if must {
				c.Text(fmt.Sprintf("Parse query param: %s failed.", match))
				c.StopExecution()
			}
		} else if n, err := strconv.ParseInt(q[match][0], 10, 64); err != nil {
			c.Text(fmt.Sprintf("Parse query param: %s failed.", match))
			c.StopExecution()
		} else if n < 0 {
			c.Text(fmt.Sprintf("Parse query param: %s out of range.", match))
			c.StopExecution()
		} else {
			c.Values().Set(match, n)
			c.Next()
			if len(alias) > 0 {
				c.Values().Set(alias, n)
				c.Next()
			}
		}
	}
}

func GQueryRequireInt(match string) func(iris.Context) { return GQueryInt(match, true, "", MaxInt) }
func GQueryRequireIntAlias(match string, alias string) func(iris.Context) {
	return GQueryInt(match, true, alias, MaxInt)
}

func GQueryOptionalInt(match string) func(iris.Context) { return GQueryInt(match, false, "", MaxInt) }
func GQueryOptionalIntAlias(match string, alias string) func(iris.Context) {
	return GQueryInt(match, false, alias, MaxInt)
}

func GQueryOptionalIntDefault(match string, defaultValue int64) func(iris.Context) {
	return GQueryInt(match, false, "", defaultValue)
}

func GQueryInt(match string, must bool, alias string, defaultValue int64) func(iris.Context) {
	return func(c iris.Context) {
		q := c.Request().URL.Query()
		if q[match] == nil || len(q[match][0]) == 0 {
			if must {
				c.Text(fmt.Sprintf("Parse query param: %s failed.", match))
				c.StopExecution()
			}
			if defaultValue == MaxInt {
				return
			}
			c.Values().Set(match, defaultValue)
			c.Next()
			if len(alias) > 0 {
				c.Values().Set(alias, defaultValue)
				c.Next()
			}
		} else if n, err := strconv.ParseInt(q[match][0], 10, 64); err != nil {
			c.Text(fmt.Sprintf("Parse query param: %s failed.", match))
			c.StopExecution()
		} else {
			c.Values().Set(match, n)
			c.Next()
			if len(alias) > 0 {
				c.Values().Set(alias, n)
				c.Next()
			}
		}
	}
}

func GQueryRequireString(match string) func(iris.Context) { return GQueryString(match, true, "") }
func GQueryRequireStringAlias(match string, alias string) func(iris.Context) {
	return GQueryString(match, true, alias)
}

func GQueryOptionalString(match string) func(iris.Context) { return GQueryString(match, false, "") }
func GQueryOptionalStringAlias(match string, alias string) func(iris.Context) {
	return GQueryString(match, false, alias)
}

func GQueryString(match string, must bool, alias string) func(iris.Context) {
	return func(c iris.Context) {
		q := c.Request().URL.Query()
		if q[match] == nil || len(q[match][0]) == 0 {
			if must {
				c.Text(fmt.Sprintf("Parse query param: %s failed.", match))
				c.StopExecution()
			}
		} else {
			c.Values().Set(match, q[match][0])
			c.Next()
			if len(alias) > 0 {
				c.Values().Set(alias, q[match][0])
				c.Next()
			}
		}
	}
}

func GQueryOptionalStringDefault(match string, defaultValue string) func(iris.Context) {
	return GQueryStringDefault(match, false, "", defaultValue)
}

func GQueryStringDefault(match string, must bool, alias, defaultValue string) func(iris.Context) {
	return func(c iris.Context) {
		q := c.Request().URL.Query()
		if q[match] == nil || len(q[match][0]) == 0 {
			if must {
				c.Text(fmt.Sprintf("Parse query param: %s failed.", match))
				c.StopExecution()
			}
			c.Values().Set(match, defaultValue)
			if len(alias) > 0 {
				c.Values().Set(alias, defaultValue)
			}
		} else {
			c.Values().Set(match, q[match][0])
			c.Next()

			if len(alias) > 0 {
				c.Values().Set(alias, q[match][0])
				c.Next()
			}
		}
	}
}

func GJsonResponse(c iris.Context) {
	c.Header("Content-Type", "application/json")
}
