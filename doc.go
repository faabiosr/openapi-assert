// Package assert provides methods that allow you to write simple swagger validations.
//
// Example Usage
//
// See it in action:
//  package main
//
//  import (
//      assert "github.com/faabiosr/openapi-assert"
//      "log"
//      "net/http"
//  )
//
//  func main() {
//      doc, err := assert.LoadFromURI("http://petstore.swagger.io/v2/swagger.json")
//
//      if err != nil {
//          log.Fatal(err)
//      }
//
//      log.Println(
//          assert.RequestMediaType("text/html", doc, "/pet", http.MethodPost),
//      )
//  }
//
// If you want to assert many times, see below:
//  package main
//
//  import (
//      assert "github.com/faabiosr/openapi-assert"
//      "log"
//      "net/http"
//  )
//
//  func main() {
//      doc, err := assert.LoadFromURI("http://petstore.swagger.io/v2/swagger.json")
//      assert := assert.New(doc)
//
//      if err != nil {
//          log.Fatal(err)
//      }
//
//      log.Println(
//          assert.RequestMediaType("text/html", "/pet", http.MethodPost),
//      )
//
//      log.Println(
//          assert.RequestMediaType("image/gif", "/v2/pet", http.MethodPost),
//      )
//  }
package assert
