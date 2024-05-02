package main

import (
	"fmt"
	"os"

	"github.com/sakura-remote-desktop/godbus/v5"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}
	defer conn.Close()

	return Emit(conn, &Org_Freedesktop_DBus_Properties_PropertiesChangedSignal{
		Path: "/org/test",
		Body: &Org_Freedesktop_DBus_Properties_PropertiesChangedSignalBody{
			InterfaceName: "org.freedesktop.test",
		},
	})
}
