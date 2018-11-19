package restful

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"io/ioutil"

	"net/http"
	"reflect"
	"strconv"
)

const tp = "text/plain"

const MaxUint = ^uint64(0)
const MaxInt = int64(MaxUint >> 1)

func GRequestBody(c iris.Context) {

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		c.JSON(iris.Map{"Code": http.StatusBadRequest, "Msg": fmt.Sprintf("Parse body failed"),"Data":nil})
		c.StopExecution()
	} else {
		c.Values().Set("requestBody", string(b))
		c.Next()
	}

}

func GRequestBodyMap(c iris.Context) {

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		c.JSON(iris.Map{"Code": http.StatusBadRequest, "Msg": fmt.Sprintf("Parse body failed"),"Data":nil})
		c.StopExecution()
	} else {
		m := map[string]interface{}{}
		if err = json.Unmarshal(b, &m); err != nil {
			c.JSON(iris.Map{"Code": http.StatusBadRequest, "Msg": fmt.Sprintf("Json unmarshal failed: %s, %v, %v", string(b), m, err),"Data":nil})
			c.StopExecution()
		} else {
			c.Values().Set("requestBody", m)
			c.Next()
		}
	}

}

func GRequestBodyObject(t reflect.Type) func(iris.Context) {
	return func(c iris.Context) {
		instance := reflect.New(t).Interface()
		b, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {

			c.JSON(iris.Map{"Code": http.StatusBadRequest, "Msg": fmt.Sprintf("Parse body failed"),"Data":nil})
			c.StopExecution()

		} else {
			if err = json.Unmarshal(b, instance); err != nil {
				c.JSON(iris.Map{"Code": http.StatusBadRequest, "Msg": fmt.Sprintf("Json unmarshal failed: %s, %v, %v", string(b), instance, err),"Data":nil})
				c.StopExecution()
			} else {
				c.Values().Set("requestBody", instance)
				c.Next()
			}
		}

	}
}

func GPathRequireInt(match string) func(iris.Context) { return GPathInt(match, true) }

func GPathOptionalInt(match string) func(iris.Context) { return GPathInt(match, false) }

func GPathInt(match string, must bool) func(iris.Context) {
	return func(c iris.Context) {
		if n, err := strconv.ParseInt(c.Params().Get(match), 10, 64); err != nil {
			if must {
				c.JSON(iris.Map{"Code": http.StatusBadRequest, "Msg": fmt.Sprintf("Parse path %s failed.", match),"Data":nil})
				c.StopExecution()
			}
		} else {
			c.Values().Set(match, n)
			c.Next()

		}
	}
}

func GPathRequireString(match string, must bool) func(iris.Context) {
	return GPathString(match, true)
}

func GPathOptionalString(match string) func(iris.Context) {
	return GPathString(match, false)
}

func GPathString(match string, must bool) func(iris.Context) {
	return func(c iris.Context) {
		if len((c.Params().Get(match))) == 0 {
			if must {
				c.JSON(iris.Map{"Code": http.StatusBadRequest, "Msg": fmt.Sprintf("Parse path %s failed.", match),"Data":nil})
				c.StopExecution()
			}
		} else {
			c.Values().Set(match, c.Params().Get(match))
			c.Next()

		}
	}
}

func GHeaderRequireInt(match string) func(iris.Context) { return GHeaderInt(match, true, "") }

func GHeaderOptionalInt(match string) func(iris.Context) { return GHeaderInt(match, false, "") }

func GHeaderInt(match string, must bool, alias string) func(iris.Context) {
	return func(c iris.Context) {
		if len(c.Request().Header[match]) > 0 {
			if id, err := strconv.ParseInt(c.Request().Header[match][0], 10, 64); err != nil {
				if must {
					c.JSON(iris.Map{"Code": http.StatusBadRequest, "Msg": fmt.Sprintf("Parse path %s failed.", match),"Data":nil})
					c.StopExecution()
				}
			} else {
				c.Values().Set(match, id)
				c.Next()

			}
		} else if must {
			c.JSON(iris.Map{"Code": http.StatusBadRequest, "Msg": fmt.Sprintf("Parse path %s failed.", match),"Data":nil})
			c.StopExecution()
		}
	}
}

func GHeaderRequireString(match string) func(iris.Context) {
	return GHeaderString(match, true)
}

func GHeaderOptionalString(match string) func(iris.Context) {
	return GHeaderString(match, false)
}

func GHeaderString(match string, must bool) func(iris.Context) {
	return func(c iris.Context) {
		if len(c.Request().Header[match]) == 0 {
			if must {
				c.JSON(iris.Map{"Code": http.StatusBadRequest, "Msg": fmt.Sprintf("Parse path %s failed.", match),"Data":nil})
				c.StopExecution()
			}
		} else {
			c.Values().Set(match, c.Request().Header[match])
			c.Next()

		}
	}
}

func GHeaderOptionalStringDefault(match string, defaultValue string) func(iris.Context) {
	return GHeaderStringDefault(match, false, defaultValue)
}

func GHeaderStringDefault(match string, must bool, defaultValue string) func(iris.Context) {
	return func(c iris.Context) {
		if len(c.Request().Header[match]) == 0 {
			if must {
				c.JSON(iris.Map{"Code": http.StatusBadRequest, "Msg": fmt.Sprintf("Parse path %s failed.", match),"Data":nil})
				c.StopExecution()
			}
			c.Values().Set(match, defaultValue)
			c.Next()

		} else {
			c.Values().Set(match, c.Request().Header.Get(match))
			c.Next()

		}
	}
}

func GQueryRequirePositiveInt(match string) func(iris.Context) {
	return GQueryPositiveInt(match, true)
}

func GQueryOptionalPositiveInt(match string) func(iris.Context) {
	return GQueryPositiveInt(match, false)
}

func GQueryPositiveInt(match string, must bool) func(iris.Context) {
	return func(c iris.Context) {
		q := c.Request().URL.Query()
		if q[match] == nil || len(q[match][0]) == 0 {
			if must {
				c.JSON(iris.Map{"Code": http.StatusBadRequest, "Msg": fmt.Sprintf("Parse path %s failed.", match),"Data":nil})
				c.StopExecution()
			}
		} else if n, err := strconv.ParseInt(q[match][0], 10, 64); err != nil {
			c.JSON(iris.Map{"Code": http.StatusBadRequest, "Msg": fmt.Sprintf("Parse path %s failed.", match),"Data":nil})
			c.StopExecution()
		} else if n < 0 {
			c.JSON(iris.Map{"Code": http.StatusBadRequest, "Msg": fmt.Sprintf("Parse path %s failed.", match),"Data":nil})
			c.StopExecution()
		} else {
			c.Values().Set(match, n)
			c.Next()

		}
	}
}

func GQueryRequireInt(match string) func(iris.Context) { return GQueryInt(match, true, MaxInt) }

func GQueryOptionalInt(match string) func(iris.Context) { return GQueryInt(match, false, MaxInt) }

func GQueryOptionalIntDefault(match string, defaultValue int64) func(iris.Context) {
	return GQueryInt(match, false, defaultValue)
}

func GQueryInt(match string, must bool, defaultValue int64) func(iris.Context) {
	return func(c iris.Context) {
		q := c.Request().URL.Query()
		if q[match] == nil || len(q[match][0]) == 0 {
			if must {
				c.JSON(iris.Map{"Code": http.StatusBadRequest, "Msg": fmt.Sprintf("Parse path %s failed.", match),"Data":nil})
				c.StopExecution()
			}
			if defaultValue == MaxInt {
				return
			}
			c.Values().Set(match, defaultValue)
			c.Next()

		} else if n, err := strconv.ParseInt(q[match][0], 10, 64); err != nil {
			c.JSON(iris.Map{"Code": http.StatusBadRequest, "Msg": fmt.Sprintf("Parse path %s failed.", match),"Data":nil})
			c.StopExecution()
		} else {
			c.Values().Set(match, n)
			c.Next()

		}
	}
}

func GQueryRequireString(match string) func(iris.Context) { return GQueryString(match, true) }

func GQueryOptionalString(match string) func(iris.Context) { return GQueryString(match, false) }

func GQueryString(match string, must bool) func(iris.Context) {
	return func(c iris.Context) {
		q := c.Request().URL.Query()
		if q[match] == nil || len(q[match][0]) == 0 {
			if must {
				c.JSON(iris.Map{"Code": http.StatusBadRequest, "Msg": fmt.Sprintf("Parse path %s failed.", match),"Data":nil})
				c.StopExecution()
			}
		} else {
			c.Values().Set(match, q[match][0])
			c.Next()

		}
	}
}

func GQueryOptionalStringDefault(match string, defaultValue string) func(iris.Context) {
	return GQueryStringDefault(match, false, defaultValue)
}

func GQueryStringDefault(match string, must bool, defaultValue string) func(iris.Context) {
	return func(c iris.Context) {
		q := c.Request().URL.Query()
		if q[match] == nil || len(q[match][0]) == 0 {
			if must {
				c.JSON(iris.Map{"Code": http.StatusBadRequest, "Msg": fmt.Sprintf("Parse path %s failed.", match),"Data":nil})
				c.StopExecution()
			}
			c.Values().Set(match, defaultValue)
			c.Next()

		} else {
			c.Values().Set(match, q[match][0])
			c.Next()

		}
	}
}

func GJsonResponse(c iris.Context) {
	c.Header("Content-Type", "application/json")
	c.Next()

}
