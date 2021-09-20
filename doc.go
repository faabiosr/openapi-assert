// Package assert provides methods that allow you to assert raw data using swagger files.
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
//      assert := assert.New(doc)
//
//      log.Println(
//          assert.RequestMediaType("text/html", "/pet", http.MethodPost),
//      )
//
//      log.Println(
//          assert.RequestMediaType("image/gif", "/v2/pet", http.MethodPost),
//      )
//  }
//
// Asserting http request object using the swagger schema file:
//  package main
//
//  import (
//      "fmt"
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
//      assert := assert.New(doc)
//
//      http.HandleFunc("/v2/pet", func(w http.ResponseWriter, r *http.Request) {
//          err := assert.Request(r)
//
//          fmt.Fprintf(w, err)
//      })
//
//      log.Fatal(
//          http.ListenAndServer("127.0.0.1:9000", nil),
//      )
//  }
//
// Asserting http response object using the swagger schema file:
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
//      assert := assert.New(doc
//
//      res, err := http.Get("https://petstore.swagger.io/v2/pet/111111422")
//
//      if err != nil {
//          log.Fatal(err)
//      }
//
//      log.Println(assert.Response(res))
//  }
package assert
