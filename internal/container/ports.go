// internal/container/ports.go
package container

import (
    //"fmt"
    "net"
)

// FindFreePort finds an available port on the host.
func FindFreePort() (int, error) {
    addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
    if err != nil {
        return 0, err
    }

    listener, err := net.ListenTCP("tcp", addr)
    if err != nil {
        return 0, err
    }
    defer listener.Close()

    return listener.Addr().(*net.TCPAddr).Port, nil
}