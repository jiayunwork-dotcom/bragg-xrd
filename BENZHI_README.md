# bragg-xrd — Go Bragg 衍射角 HTTP 后端与命令行工具

The server starts on `:8080` by default and exposes `POST /api/bragg` and
`POST /api/powder`; CLI subcommands `bragg` and `powder` provide the same
calculations in a terminal.

## Build / Run / Test

```text
go build ./...
go run . serve -addr :8080
go run . example -file example/fcc-cu-cu.json
go test ./...
```

## Evaluation Image

Evaluation-specific files (do not overwrite project Dockerfile/README):

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md` (this file)

Build and verify in container:

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
