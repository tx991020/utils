package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
)

const tp = "text/plain"

const MaxUint = ^uint64(0)
const MaxInt = int64(MaxUint >> 1)

func GRequestBody(c *gin.Context) {
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Data(http.StatusBadRequest, tp, []byte("Parse body failed."))
		c.Abort()
	} else {
		c.Set("requestBody", b)
	}
	return
}

func GRequestBodyMap(c *gin.Context) {
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Data(http.StatusBadRequest, tp, []byte("Parse body failed."))
		c.Abort()
	} else {
		m := map[string]interface{}{}
		if err = json.Unmarshal(b, &m); err != nil {
			c.Data(http.StatusBadRequest, tp, []byte(fmt.Sprintf("Json unmarshal failed: %s, %v, %v", string(b), m, err)))
		} else {
			c.Set("requestBody", m)
		}
	}
	return
}

// idl = {json, xml}
func GRequestBodyObject(t reflect.Type, idl string) func(*gin.Context) {
	return func(c *gin.Context) {
		instance := reflect.New(t).Interface()

		b, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.Data(http.StatusBadRequest, tp, []byte("Parse body failed."))
			c.Abort()
		} else {
			switch idl {
			case "json":
				if err = json.Unmarshal(b, instance); err != nil {
					c.Data(http.StatusBadRequest, tp, []byte(fmt.Sprintf("Json unmarshal failed: %s, %v, %v", string(b), instance, err)))
					c.Abort()
				} else {
					c.Set("requestBody", instance)
				}
			case "xml":
			default:
			}
		}
		return
	}
}

func GPathRequireInt(match string) func(*gin.Context) { return GPathInt(match, true, "") }
func GPathRequireIntAlias(match string, alias string) func(*gin.Context) {
	return GPathInt(match, true, alias)
}

func GPathOptionalInt(match string) func(*gin.Context) { return GPathInt(match, false, "") }
func GPathOptionalIntAlias(match string, alias string) func(*gin.Context) {
	return GPathInt(match, false, alias)
}

func GPathInt(match string, must bool, alias string) func(*gin.Context) {
	return func(c *gin.Context) {
		if n, err := strconv.ParseInt(c.Param(match), 10, 64); err != nil {
			if must {
				c.Data(http.StatusBadRequest, tp, []byte(fmt.Sprintf("Parse path %s failed.", match)))
				c.Abort()
			}
		} else {
			c.Set(match, n)
			if len(alias) > 0 {
				c.Set(alias, n)
			}
		}
	}
}

func GPathRequireString(match string) func(*gin.Context) {
	return GPathString(match, true, "")
}
func GPathRequireStringAlias(match string, must bool, alias string) func(*gin.Context) {
	return GPathString(match, true, alias)
}

func GPathOptionalString(match string) func(*gin.Context) {
	return GPathString(match, false, "")
}
func GPathOptionalStringAlias(match string, alias string) func(*gin.Context) {
	return GPathString(match, false, alias)
}

func GPathString(match string, must bool, alias string) func(*gin.Context) {
	return func(c *gin.Context) {
		if len(c.Param(match)) == 0 {
			if must {
				c.Data(http.StatusBadRequest, tp, []byte(fmt.Sprintf("Parse path param: %s failed.", match)))
				c.Abort()
			}
		} else {
			c.Set(match, c.Param(match))
			if len(alias) > 0 {
				c.Set(alias, c.Param(match))
			}
		}
	}
}

func GHeaderRequireInt(match string) func(*gin.Context) { return GHeaderInt(match, true, "") }
func GHeaderRequireIntAlias(match string, alias string) func(*gin.Context) {
	return GHeaderInt(match, true, alias)
}

func GHeaderOptionalInt(match string) func(*gin.Context) { return GHeaderInt(match, false, "") }
func GHeaderOptionalIntAlias(match string, alias string) func(*gin.Context) {
	return GHeaderInt(match, false, alias)
}

func GHeaderInt(match string, must bool, alias string) func(*gin.Context) {
	return func(c *gin.Context) {
		if len(c.Request.Header[match]) > 0 {
			if id, err := strconv.ParseInt(c.GetHeader(match), 10, 64); err != nil {
				if must {
					c.Data(http.StatusBadRequest, tp, []byte(fmt.Sprintf("Parse header: %s failed.", match)))
					c.Abort()
				}
			} else {
				c.Set(match, id)
				if len(alias) > 0 {
					c.Set(alias, id)
				}
			}
		} else if must {
			c.Data(http.StatusBadRequest, tp, []byte(fmt.Sprintf("Parse header: %s failed.", match)))
			c.Abort()
		}
	}
}

func GHeaderRequireString(match string) func(*gin.Context) {
	return GHeaderString(match, true, "")
}
func GHeaderRequireStringAlias(match string, alias string) func(*gin.Context) {
	return GHeaderString(match, true, alias)
}

func GHeaderOptionalString(match string) func(*gin.Context) {
	return GHeaderString(match, false, "")
}
func GHeaderOptionalStringAlias(match string, alias string) func(*gin.Context) {
	return GHeaderString(match, false, alias)
}

func GHeaderString(match string, must bool, alias string) func(*gin.Context) {
	return func(c *gin.Context) {
		if len(c.GetHeader(match)) == 0 {
			if must {
				c.Data(http.StatusBadRequest, tp, []byte(fmt.Sprintf("Parse header: %s failed.", match)))
				c.Abort()
			}
		} else {

			c.Set(match, c.GetHeader(match))
			if len(alias) > 0 {
				c.Set(alias, c.Request.Header[match])

			}
		}
	}
}

func GHeaderOptionalStringDefault(match string, defaultValue string) func(*gin.Context) {
	return GHeaderStringDefault(match, false, "", defaultValue)
}

func GHeaderStringDefault(match string, must bool, alias, defaultValue string) func(*gin.Context) {
	return func(c *gin.Context) {
		if len(c.Request.Header[match]) == 0 {
			if must {
				c.Data(http.StatusBadRequest, tp, []byte(fmt.Sprintf("Parse header: %s failed.", match)))
				c.Abort()
			}
			c.Set(match, defaultValue)
			if len(alias) > 0 {
				c.Set(alias, defaultValue)
			}
		} else {
			c.Set(match, c.Request.Header.Get(match))
			if len(alias) > 0 {
				c.Set(alias, c.Request.Header.Get(match))
			}
		}
	}
}

func GQueryRequirePositiveInt(match string) func(*gin.Context) {
	return GQueryPositiveInt(match, true, "")
}
func GQueryRequirePositiveIntAlias(match string, alias string) func(*gin.Context) {
	return GQueryPositiveInt(match, true, alias)
}

func GQueryOptionalPositiveInt(match string) func(*gin.Context) {
	return GQueryPositiveInt(match, false, "")
}
func GQueryOptionalPositiveIntAlias(match string, alias string) func(*gin.Context) {
	return GQueryPositiveInt(match, false, alias)
}

func GQueryPositiveInt(match string, must bool, alias string) func(*gin.Context) {
	return func(c *gin.Context) {
		q := c.Request.URL.Query()
		if q[match] == nil || len(q[match][0]) == 0 {
			if must {
				c.Data(http.StatusBadRequest, tp, []byte(fmt.Sprintf("Parse query param: %s not exist.", match)))
				c.Abort()
			}
		} else if n, err := strconv.ParseInt(q[match][0], 10, 64); err != nil {
			c.Data(http.StatusBadRequest, tp, []byte(fmt.Sprintf("Parse query param: %s failed.", match)))
			c.Abort()
		} else if n < 0 {
			c.Data(http.StatusBadRequest, tp, []byte(fmt.Sprintf("Parse query param: %s out of range.", match)))
			c.Abort()
		} else {
			c.Set(match, n)
			if len(alias) > 0 {
				c.Set(alias, n)
			}
		}
	}
}

func GQueryRequireInt(match string) func(*gin.Context) { return GQueryInt(match, true, "", MaxInt) }
func GQueryRequireIntAlias(match string, alias string) func(*gin.Context) {
	return GQueryInt(match, true, alias, MaxInt)
}

func GQueryOptionalInt(match string) func(*gin.Context) { return GQueryInt(match, false, "", MaxInt) }
func GQueryOptionalIntAlias(match string, alias string) func(*gin.Context) {
	return GQueryInt(match, false, alias, MaxInt)
}

func GQueryOptionalIntDefault(match string, defaultValue int64) func(*gin.Context) {
	return GQueryInt(match, false, "", defaultValue)
}

func GQueryInt(match string, must bool, alias string, defaultValue int64) func(*gin.Context) {
	return func(c *gin.Context) {
		q := c.Request.URL.Query()
		if q[match] == nil || len(q[match][0]) == 0 {
			if must {
				c.Data(http.StatusBadRequest, tp, []byte(fmt.Sprintf("Parse query param: %s not exist.", match)))
				c.Abort()
			}
			if defaultValue == MaxInt {
				return
			}
			c.Set(match, defaultValue)
			if len(alias) > 0 {
				c.Set(alias, defaultValue)
			}
		} else if n, err := strconv.ParseInt(q[match][0], 10, 64); err != nil {
			c.Data(http.StatusBadRequest, tp, []byte(fmt.Sprintf("Parse query param: %s failed.", match)))
			c.Abort()
		} else {
			c.Set(match, n)
			if len(alias) > 0 {
				c.Set(alias, n)
			}
		}
	}
}

func GQueryRequireString(match string) func(*gin.Context) { return GQueryString(match, true, "") }
func GQueryRequireStringAlias(match string, alias string) func(*gin.Context) {
	return GQueryString(match, true, alias)
}

func GQueryOptionalString(match string) func(*gin.Context) { return GQueryString(match, false, "") }
func GQueryOptionalStringAlias(match string, alias string) func(*gin.Context) {
	return GQueryString(match, false, alias)
}

func GQueryString(match string, must bool, alias string) func(*gin.Context) {
	return func(c *gin.Context) {
		q := c.Request.URL.Query()
		if q[match] == nil || len(q[match][0]) == 0 {
			if must {
				c.Data(http.StatusBadRequest, tp, []byte(fmt.Sprintf("Parse query param: %s failed.", match)))
				c.Abort()
			}
		} else {
			c.Set(match, q[match][0])
			if len(alias) > 0 {
				c.Set(alias, q[match][0])
			}
		}
	}
}

func GQueryOptionalStringDefault(match string, defaultValue string) func(*gin.Context) {
	return GQueryStringDefault(match, false, "", defaultValue)
}

func GQueryStringDefault(match string, must bool, alias, defaultValue string) func(*gin.Context) {
	return func(c *gin.Context) {
		q := c.Request.URL.Query()
		if q[match] == nil || len(q[match][0]) == 0 {
			if must {
				c.Data(http.StatusBadRequest, tp, []byte(fmt.Sprintf("Parse query param: %s failed.", match)))
				c.Abort()
			}
			c.Set(match, defaultValue)
			if len(alias) > 0 {
				c.Set(alias, defaultValue)
			}
		} else {
			c.Set(match, q[match][0])
			if len(alias) > 0 {
				c.Set(alias, q[match][0])
			}
		}
	}
}

func GJsonResponse(c *gin.Context) {
	c.Header("Content-Type", "application/json")
}

func GXmlResponse(c *gin.Context) {
	c.Header("Content-Type", "application/xml")
}
