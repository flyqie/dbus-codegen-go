package main

import (
	"context"
	"errors"
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

	sigc := make(chan *dbus.Signal, 1)
	conn.Signal(sigc)

	const name = "dbusgen.test"
	o := NewOrg_Freedesktop_DBus(conn.Object("org.freedesktop.DBus", "/org/freedesktop/DBus"))
	ret, err := o.RequestName(context.Background(), name, 3)
	if err != nil {
		return err
	}
	if ret != 1 {
		return errors.New("unexpected return code")
	}
Retry:
	s, err := LookupSignal(<-sigc)
	if err != nil {
		return err
	}
	sig := s.(*Org_Freedesktop_DBus_NameAcquiredSignal)
	if sig.Body.V0[0] == ':' {
		goto Retry
	}
	if sig.Sender() != "org.freedesktop.DBus" ||
		sig.Path != "/org/freedesktop/DBus" ||
		sig.Body.V0 != name {
		return fmt.Errorf("invalid signal = %v, want %v", sig, []interface{}{
			sig.Sender(), sig.Path, name,
		})
	}
	return nil
}
