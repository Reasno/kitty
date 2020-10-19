package cmd

import (
	"fmt"
	"github.com/Reasno/kitty/pkg/execprotoc"
	"github.com/Reasno/kitty/pkg/parsesvcname"
	"github.com/metaverse/truss/svcdef"
	"github.com/metaverse/truss/truss"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go/build"
	"golang.org/x/tools/go/packages"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	ggkconf "github.com/metaverse/truss/gengokit"
	gengokit "github.com/metaverse/truss/gengokit/generator"
)

func init() {
	generateCmd.Flags().StringVar(&svcout, "svc_out", "", "the service output dir")
	rootCmd.AddCommand(generateCmd)
}

var (
	svcout string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate kitty services",
	Long:  `Reflect changes from protobuf to code.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			args = append(args, "proto/app.proto")
		}

		cfg, err := parseInput(args)
		if err != nil {
			log.Fatal(errors.Wrap(err, "cannot parse input"))
		}

		// If there was no service found in parseInput, the rest can be omitted.
		if cfg == nil {
			return
		}

		sd, err := parseServiceDefinition(cfg)
		if err != nil {
			log.Fatal(errors.Wrap(err, "cannot parse input definition proto files"))
		}

		genFiles, err := generateCode(cfg, sd)
		if err != nil {
			log.Fatal(errors.Wrap(err, "cannot generate service"))
		}

		for path, file := range genFiles {
			err := writeGenFile(file, filepath.Join(cfg.ServicePath, path))
			if err != nil {
				log.Fatal(errors.Wrap(err, "cannot to write output"))
			}
		}

		cleanupOldFiles(cfg.ServicePath, strings.ToLower(sd.Service.Name))

		//if err = executeWire(); err != nil {
		//	log.Fatal(err)
		//}
	},
}

// parseInput constructs a *truss.Config with all values needed to parse
// service definition files.
func parseInput(args []string) (*truss.Config, error) {
	var cfg truss.Config

	// GOPATH
	cfg.GoPath = filepath.SplitList(os.Getenv("GOPATH"))
	if len(cfg.GoPath) == 0 {
		cfg.GoPath = filepath.SplitList(build.Default.GOPATH)
	}
	log.WithField("GOPATH", cfg.GoPath).Debug()

	// DefPaths
	var err error
	cfg.DefPaths, err = cleanProtofilePath(args)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse input arguments")
	}
	log.WithField("DefPaths", cfg.DefPaths).Debug()

	protoDir := filepath.Dir(cfg.DefPaths[0])
	p, err := packages.Load(nil, protoDir)
	if err != nil || len(p) == 0 {
		return nil, errors.Wrap(err, "proto files not found in importable go package")
	}

	cfg.PBPackage = p[0].PkgPath
	cfg.PBPath = protoDir
	log.WithField("PB Package", cfg.PBPackage).Debug()
	log.WithField("PB Path", cfg.PBPath).Debug()

	if err := execprotoc.GeneratePBDotGo(cfg.DefPaths, cfg.GoPath, cfg.PBPath); err != nil {
		return nil, errors.Wrap(err, "cannot create .pb.go files")
	}

	// Service Path
	svcName, err := parsesvcname.FromPaths(cfg.GoPath, cfg.DefPaths)
	if err != nil {
		log.Warnf("No valid service is defined; exiting now: %v", err)
		log.Info(".pb.go generation with protoc-gen-go was successful.")
		return nil, nil
	}

	svcName = strings.ToLower(svcName)

	svcDirName := svcName
	log.WithField("svcDirName", svcDirName).Debug()

	svcPath := filepath.Join(filepath.Dir(filepath.Dir(cfg.DefPaths[0])), svcDirName)

	if svcout != "" {
		svcOut := svcout
		log.WithField("svcPackageFlag", svcOut).Debug()

		// If the package flag ends in a seperator, file will be "".
		_, file := filepath.Split(svcOut)
		seperator := file == ""
		log.WithField("seperator", seperator)

		svcPath, err = parseSVCOut(svcOut, cfg.GoPath[0])
		if err != nil {
			return nil, errors.Wrapf(err, "cannot parse svcout: %s", svcOut)
		}

		// Join the svcDirName as a svcout ending with `/` should create it
		if seperator {
			svcPath = filepath.Join(svcPath, svcDirName)
		}
	}

	log.WithField("svcPath", svcPath).Debug()

	// Create svcPath for the case that it does not exist
	err = os.MkdirAll(svcPath, 0777)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create svcPath directory: %s", svcPath)
	}

	p, err = packages.Load(nil, svcPath)
	if err != nil || len(p) == 0 {
		return nil, errors.Wrap(err, "generated service not found in importable go package")
	}

	log.WithField("Service Packages", p).Debug()

	cfg.ServicePackage = p[0].PkgPath
	cfg.ServicePath = svcPath

	log.WithField("Service Package", cfg.ServicePackage).Debug()
	log.WithField("package name", p[0].Name).Debug()
	log.WithField("Service Path", cfg.ServicePath).Debug()

	// PrevGen
	cfg.PrevGen, err = readPreviousGeneration(cfg.ServicePath)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read previously generated files")
	}

	return &cfg, nil
}

// parseSVCOut handles the difference between relative paths and go package
// paths
func parseSVCOut(svcOut string, GOPATH string) (string, error) {
	if build.IsLocalImport(svcOut) {
		return filepath.Abs(svcOut)
	}
	return filepath.Join(GOPATH, "src", svcOut), nil
}

// parseServiceDefinition returns a svcdef which contains all necessary
// information for generating a truss service.
func parseServiceDefinition(cfg *truss.Config) (*svcdef.Svcdef, error) {
	protoDefPaths := cfg.DefPaths
	// Create the ServicePath so the .pb.go files may be place within it
	if cfg.PrevGen == nil {
		err := os.MkdirAll(cfg.ServicePath, 0777)
		if err != nil {
			return nil, errors.Wrap(err, "cannot create service directory")
		}
	}

	// Get path names of .pb.go files
	pbgoPaths := []string{}
	for _, p := range protoDefPaths {
		base := filepath.Base(p)
		barename := strings.TrimSuffix(base, filepath.Ext(p))
		pbgp := filepath.Join(cfg.PBPath, barename+".pb.go")
		pbgoPaths = append(pbgoPaths, pbgp)
	}
	pbgoFiles, err := openFiles(pbgoPaths)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open all .pb.go files")
	}

	pbFiles, err := openFiles(protoDefPaths)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open all .proto files")
	}

	// Create the svcdef
	sd, err := svcdef.New(pbgoFiles, pbFiles)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create service definition; did you pass ALL the protobuf files to truss?")
	}

	return sd, nil
}

// generateCode returns a map[string]io.Reader that represents a gokit
// service
func generateCode(cfg *truss.Config, sd *svcdef.Svcdef) (map[string]io.Reader, error) {
	conf := ggkconf.Config{
		PBPackage:     cfg.PBPackage,
		GoPackage:     cfg.ServicePackage,
		PreviousFiles: cfg.PrevGen,
	}

	genGokitFiles, err := gengokit.GenerateGokit(sd, conf)
	if err != nil {
		return nil, errors.Wrap(err, "cannot generate gokit service")
	}

	return genGokitFiles, nil
}

func openFiles(paths []string) (map[string]io.Reader, error) {
	rv := map[string]io.Reader{}
	for _, p := range paths {
		reader, err := os.Open(p)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot open file %q", p)
		}
		rv[p] = reader
	}
	return rv, nil
}

// writeGenFile writes a file at path to the filesystem
func writeGenFile(file io.Reader, path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return err
	}

	outFile, err := os.Create(path)
	if err != nil {
		return errors.Wrapf(err, "cannot create file %v", path)
	}

	str, err := ioutil.ReadAll(file)
	if err != nil {
		return errors.Wrapf(err, "cannot read file")
	}

	file = strings.NewReader(strings.Replace(string(str), "github.com/gogo/protobuf/jsonpb", "github.com/golang/protobuf/jsonpb", -1))

	_, err = io.Copy(outFile, file)
	if err != nil {
		return errors.Wrapf(err, "cannot write to %v", path)
	}
	return outFile.Close()
}

// cleanProtofilePath returns the absolute filepath of a group of files
// of the files, or an error if the files are not in the same directory
func cleanProtofilePath(rawPaths []string) ([]string, error) {
	var fullPaths []string

	// Parsed passed file paths
	for _, def := range rawPaths {
		log.WithField("rawDefPath", def).Debug()
		full, err := filepath.Abs(def)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get working directory of truss")
		}
		log.WithField("fullDefPath", full)

		fullPaths = append(fullPaths, full)

		if filepath.Dir(fullPaths[0]) != filepath.Dir(full) {
			return nil, errors.Errorf("passed .proto files in different directories")
		}
	}

	return fullPaths, nil
}

// readPreviousGeneration returns a map[string]io.Reader representing the files in serviceDir
func readPreviousGeneration(serviceDir string) (map[string]io.Reader, error) {
	if !fileExists(serviceDir) {
		return nil, nil
	}

	const handlersDirName = "handlers"
	files := make(map[string]io.Reader)

	addFileToFiles := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			switch info.Name() {
			// Only files within the handlers dir are used to
			// support regeneration.
			// See `gengokit/generator/gen.go:generateResponseFile`
			case filepath.Base(serviceDir), handlersDirName:
				return nil
			default:
				return filepath.SkipDir
			}
		}

		file, ioErr := os.Open(path)
		if ioErr != nil {
			return errors.Wrapf(ioErr, "cannot read file: %v", path)
		}

		// trim the prefix of the path to the proto files from the full path to the file
		relPath, err := filepath.Rel(serviceDir, path)
		if err != nil {
			return err
		}

		// ensure relPath is unix-style, so it matches what we look for later
		relPath = filepath.ToSlash(relPath)

		files[relPath] = file

		return nil
	}

	err := filepath.Walk(serviceDir, addFileToFiles)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot fully walk directory %v", serviceDir)
	}

	return files, nil
}

// fileExists checks if a file at the given path exists. Returns true if the
// file exists, and false if the file does not exist.
func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func cleanupOldFiles(servicePath, serviceName string) {
	serverCLI := filepath.Join(servicePath, "svc/server/cli")
	if _, err := os.Stat(serverCLI); err == nil {
		log.Warnf("Removing stale 'svc/server/cli' files")
		err := os.RemoveAll(serverCLI)
		if err != nil {
			log.Error(err)
		}
	}
	clientCLI := filepath.Join(servicePath, "svc/client/cli")
	if _, err := os.Stat(clientCLI); err == nil {
		log.Warnf("Removing stale 'svc/client/cli' files")
		err := os.RemoveAll(clientCLI)
		if err != nil {
			log.Error(err)
		}
	}

	oldServer := filepath.Join(servicePath, fmt.Sprintf("cmd/%s-server", serviceName))
	if _, err := os.Stat(oldServer); err == nil {
		log.Warnf(fmt.Sprintf("Removing stale 'cmd/%s-server' files, use cmd/%s going forward", serviceName, serviceName))
		err := os.RemoveAll(oldServer)
		if err != nil {
			log.Error(err)
		}
	}
}

func executeWire() error {
	wireExec := exec.Command(
		"wire",
		"./...",
	)

	outBytes, err := wireExec.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err,
			"wire exec failed.\nwire output:\n\n%v\nprotoc arguments:\n\n%v\n\n",
			string(outBytes), wireExec.Args)
	}
	return nil
}


