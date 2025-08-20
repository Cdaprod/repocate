package main

import (
  "fmt"
  "os/exec"
)

func main() {
  fmt.Println("Bootstrapping your node...")
  cmd := exec.Command("bash", "-c", "scripts/bootstrap.sh")
  output, err := cmd.CombinedOutput()
  if err != nil {
    fmt.Printf("Error: %s\n", err)
  }
  fmt.Println(string(output))
}