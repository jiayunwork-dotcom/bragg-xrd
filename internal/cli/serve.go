package cli

import "bragg-xrd/internal/server"

func runServe(args []string) int {
	fs := flagSet("serve")
	addr := fs.String("addr", ":8080", "HTTP listen address")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	srv := server.New(*addr)
	if err := srv.RunWithGracefulShutdown(); err != nil && !server.IsServerClosed(err) {
		return fail(err)
	}
	return 0
}
