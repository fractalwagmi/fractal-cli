# Fractal CLI

CLI tool for interacting with Fractal's studio tools and APIs.

## Installation

You can install the CLI by building from source. You will need to have Go 1.17+
installed on your machine ([link to installer download](https://go.dev/dl/)).

```bash
go install github.com/fractalwagmi/fractal-cli/cmd/fractal@latest
```

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
