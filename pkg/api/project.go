package api

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/fedux-org/letter-generator-go/assets"
	"github.com/fedux-org/letter-generator-go/converter"
	"github.com/fedux-org/letter-generator-go/latex"
	"github.com/fedux-org/letter-generator-go/letter"
	lgos "github.com/fedux-org/letter-generator-go/os"
	log "github.com/sirupsen/logrus"
)

type Project struct {
	letters  []letter.Letter
	template converter.Template
	assets   []assets.Asset
	outDir   string
	workDir  string
}

func NewProject(letters []letter.Letter, template converter.Template, assets []assets.Asset, outDir string) Project {
	p := Project{letters: letters, template: template, assets: assets, outDir: outDir}
	p.SetupWorkDir()

	return p
}

func (p *Project) SetupWorkDir() error {
	dir, err := ioutil.TempDir("", "lg")
	if err != nil {
		return errors.Wrap(err, "create temporary work dir")
	}

	p.workDir = dir

	return nil
}

func (p *Project) Build() error {
	var movableAssets []MovableFile = make([]MovableFile, len(p.assets))
	for i, d := range p.assets {
		movableAssets[i] = &d
	}

	texFiles, err := generateTexFiles(p.template, p.letters)

	defer func() {
		log.Debug("Invoke clean up function")
		for _, f := range texFiles {
			f.Destroy()
		}
	}()

	if err != nil {
		return errors.Wrap(err, "generate tex file")
	}

	for _, f := range texFiles {
		err := moveFilesToDir(movableAssets, f.Dir)
		if err != nil {
			return errors.Wrap(err, "move files")
		}
	}

	pdfFiles, err := compileTexFilesIntoPdf(texFiles)
	if err != nil {
		return errors.Wrap(err, "compile tex into pdf")
	}

	var movablePdfFiles []MovableFile = make([]MovableFile, len(pdfFiles))
	for i, d := range pdfFiles {
		movablePdfFiles[i] = &d
	}

	err = moveFilesToDir(movablePdfFiles, p.outDir)
	if err != nil {
		return errors.Wrap(err, "move files")
	}

	err = os.RemoveAll(p.workDir)
	if err != nil {
		return errors.Wrap(err, "remove work dir")
	}

	for _, f := range texFiles {
		log.WithField("working_directory", f.Dir).Debug("Remove working directory")
		os.RemoveAll(f.Dir)
	}

	log.WithField("working_directory", p.workDir).Debug("Remove working directory")

	files, err := filepath.Glob(filepath.Join(p.outDir, "*.pdf"))

	log.WithFields(log.Fields{"count(letters)": len(files), "files": strings.Join(files, ",")}).Info("Generate letters")

	return nil
}

func generateTexFiles(template converter.Template, letters []letter.Letter) ([]converter.TexFile, error) {
	var texFiles []converter.TexFile

	for _, l := range letters {
		texFile, err := renderTemplate(l, template)
		texFiles = append(texFiles, texFile)

		if err != nil {
			return texFiles, err
		}
	}

	log.WithFields(log.Fields{"count(tex files)": len(texFiles)}).Info("Generated tex files")

	return texFiles, nil
}

func compileTexFilesIntoPdf(texFiles []converter.TexFile) ([]converter.PdfFile, error) {
	compiler := latex.NewCompiler()
	var pdfFiles []converter.PdfFile

	for _, f := range texFiles {
		pdfFile, err := compiler.Compile(f)

		if err != nil {
			return []converter.PdfFile{}, err
		}

		log.WithFields(log.Fields{
			"input_file":  f.Path,
			"output_file": pdfFile.Path,
		}).Debug("Render letter as PDF")

		pdfFiles = append(pdfFiles, pdfFile)
	}

	return pdfFiles, nil
}

func moveFilesToDir(files []MovableFile, dir string) error {
	for _, f := range files {
		filename := filepath.Base(f.GetPath())
		newPath := filepath.Join(dir, filename)

		err := lgos.Copy(f.GetPath(), newPath)
		if err != nil {
			return err
		}

		log.WithFields(log.Fields{
			"source":      f.GetPath(),
			"destination": newPath,
		}).Debug("Moving file to dir")
	}

	return nil
}

func renderTemplate(l letter.Letter, t converter.Template) (converter.TexFile, error) {
	templateConverter := converter.NewConverter()
	texFile, err := templateConverter.Transform(l, t)

	if err != nil {
		return converter.TexFile{}, errors.Wrap(err, "render template into tex file")
	}

	log.WithFields(log.Fields{
		"path(tex_file)": texFile.Path,
		"path(template)": t.Path,
	}).Debug("Creating tex file from template")

	return texFile, nil

}
