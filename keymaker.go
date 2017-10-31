package main

import (
  "github.com/sheng/air"
  "github.com/astaxie/beego/orm"
  _ "github.com/mattn/go-sqlite3"
  "github.com/kamildrazkiewicz/go-flow"
  "math/rand"
  "time"
  "strconv"
  "fmt"
)

// Model Struct
type Key struct {
  Id   int    `orm:"auto"`
  Value string `orm:"size(256)"`
  Status string `orm:"size(10)"`
}

var o orm.Ormer

func init() {
  // register model
  orm.RegisterModel(new(Key))

  name := "default"

  maxIdle := 30
  maxConn := 30

  // set default database
  orm.RegisterDataBase(name, "sqlite3", "data.db", maxIdle, maxConn)

  force := true
  verbose := true

  orm.RunSyncdb(name, force, verbose)
}

func main() {
  a := air.New()
  a.GET("/v1/get", keyGetter)
  a.GET("/v1/validate/:id", keyValidate)

  o = orm.NewOrm()

  test_key := Key{Value: "test", Status: "test"}

  o.Insert(&test_key)

  a.Serve()
}


func keyGetter(c *air.Context) error {
  keyLength := 64

  value := RandStringBytesMaskImprSrc(keyLength)
  key := Key{Value: value, Status: "active"}
  go newKeyWorker(key)

  return c.String(value)
}

func keyValidate(c *air.Context) error {
  qs := o.QueryTable("key")
  var key Key
  qs.Filter("value", c.Param("id")).Filter("status", "active").One(&key)
  result := strconv.FormatBool(key.Value == c.Param("id"))
  return c.String(result)
}


// Hash function from http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang/31832326
var src = rand.NewSource(time.Now().UnixNano())
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_0123456789"
const (
    letterIdxBits = 6                    // 6 bits to represent a letter index
    letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
    letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytesMaskImprSrc(n int) string {
    b := make([]byte, n)
    // A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
    for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
        if remain == 0 {
            cache, remain = src.Int63(), letterIdxMax
        }
        if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
            b[i] = letterBytes[idx]
            i--
        }
        cache >>= letterIdxBits
        remain--
    }

    return string(b)
}

func newKeyWorker(k Key) {
  active := func(r map[string]interface{}) (interface{}, error) {
    time.Sleep(time.Millisecond * 500)
    fmt.Println("active started", "---", k.Value)
    o.Insert(&k)
    fmt.Println("active finished", "---", k.Value)
    return nil, nil
  }

  inactive := func(r map[string]interface{}) (interface{}, error) {
    time.Sleep(time.Millisecond * 30000)
    fmt.Println("inactive started", "---", k.Value)
    qs := o.QueryTable("key")
    var key Key
    qs.Filter("id", k.Id).One(&key)
    key.Status = "inactive"
    o.Update(&key, "Status")
    fmt.Println("inactive finished", "---", k.Value)
    return nil, nil
  }

  killed := func(r map[string]interface{}) (interface{}, error) {
    time.Sleep(time.Millisecond * 60000)
    fmt.Println("killed started", "---", k.Value)
    o.Delete(&Key{Id: k.Id})
    fmt.Println("killed finished", "---", k.Value)
    return nil, nil
  }

  goflow.New().
  Add("active", nil, active).
  Add("inactive", nil, inactive).
  Add("killed", nil, killed).
  Do()
}
