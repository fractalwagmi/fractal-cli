# Fractal CLI

CLI tool for interacting with Fractal's studio tools and APIs.

## Installation

You can install the CLI by building from source. You will need to have Go 1.17+
installed on your machine ([link to installer download](https://go.dev/dl/)).

```bash
go install github.com/fractalwagmi/fractal-cli/cmd/fractal@latest
```

For convenience, you can add the Go bin directory to your PATH so that you can
invoke the CLI anywhere using `fractal` instead of the full path to the binary.

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Alternatively, you can download and run a pre-built binary from the project's
[releases page](https://github.com/fractalwagmi/fractal-cli/releases).

## Usage

The only command currently supported is `upload` which can be used to upload and
configure a new build of your game. It should be called once for each platform
you wish to upload.

Deployments are still manual and should be triggered in FStudio after the build
is uploaded.

```bash
fractal upload \
  -zip=<path to a .zip archive containing your game files> \
  -clientId=<your client ID obtained from FStudio> \
  -clientSecret=<your client secret obtained from FStudio> \
  -version=<globally (to your project) unique version name>
  -platform=<windows|mac|universal> \
  -exeFile=<only windows or universal, path to .exe game executable> \
  -macAppDirectory=<only mac or universal, path to .app directory> \
  -macInnerExecutable=<only mac or universal, path to mac inner executable inside .app/Contents/MacOS> \
```

## Support

Please file a github issue here if you have any questions or want to report a
bug.
