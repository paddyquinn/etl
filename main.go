package main

import (
  "bufio"
  "encoding/json"
  "fmt"
  "net/http"
  "os"

  "github.com/paddyquinn/etl/models"
)

func main() {
  // Fetch the jsonl file.
  rsp, err := http.Get("https://s3-us-west-1.amazonaws.com/circleup-engr-interview-public/simple-etl.jsonl")
  if err != nil {
    fmt.Println(err.Error())
    os.Exit(1)
  }

  // Unmarshal each line into our object model and then transform each object and store it in our transformedObjects
  // slice.
  var transformedObjects []*models.TransformedObject
  scanner := bufio.NewScanner(rsp.Body)
  for scanner.Scan() {
    var o models.Object
    err = json.Unmarshal(scanner.Bytes(), &o)
    if err != nil {
      fmt.Println(err.Error())
      os.Exit(2)
    }
    transformedObjects = append(transformedObjects, o.Transform())
  }

  // Check if the scanner errored when exiting the loop.
  if err = scanner.Err(); err != nil {
    fmt.Println(err.Error())
    os.Exit(3)
  }

  // Find all the users whose first names start with J.
  var startsWithJList []*models.TransformedObject
  for _, transformedObject := range transformedObjects {
    if transformedObject.FullName[0] == 'J' {
      startsWithJList = append(startsWithJList, transformedObject)
    }
  }


  // Create the file to output solution 1.1 to. If the file exists already it will be overwritten.
  file, err := os.Create("1.1.json")
  if err != nil {
    fmt.Println(err.Error())
    os.Exit(4)
  }

  // Marshal our list into a JSON object.
  bytes, err := json.Marshal(startsWithJList)
  if err != nil {
    fmt.Println(err.Error())
    os.Exit(5)
  }

  // Write the marshaled bytes to our file.
  writer := bufio.NewWriter(file)
  writer.Write(bytes)
  writer.Flush()
}