package api

import (
	_ "net/url"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/feduxorg/letter-generator/assets"
	"github.com/feduxorg/letter-generator/converter"
	"github.com/feduxorg/letter-generator/letter"
	"github.com/feduxorg/letter-generator/letter_generator"
	"github.com/feduxorg/letter-generator/metadata"
	"github.com/feduxorg/letter-generator/recipients"
	"github.com/feduxorg/letter-generator/sender"
	log "github.com/sirupsen/logrus"
)

type LetterBuilder struct{}

func (lc *LetterBuilder) Build(config letter_generator.Config) error {
	metadata, err := readMetadata(config.MetadataFile)
	if err != nil {
		return err
	}

	sender, err := readSender(config.SenderFile)
	if err != nil {
		return err
	}

	recipientManager, err := readRecipients(config.RecipientsFile)
	if err != nil {
		return err
	}

	template, err := readLetterTemplate(config.TemplateFile)
	if err != nil {
		return err
	}

	outputDirectory, err := createOutputDirectory()
	if err != nil {
		return err
	}

	assets, err := findAssets(config.AssetsDirectory)
	if err != nil {
		return err
	}

	letters := generateLetterInstances(sender, metadata, recipientManager.Recipients)

	project := NewProject(letters, template, assets, outputDirectory)

	err = project.Build()
	if err != nil {
		return err
	}

	return nil
}

func readMetadata(srcFile string) (metadata.Metadata, error) {
	m := metadata.Metadata{}
	err := m.Read(srcFile)

	if err != nil {
		return metadata.Metadata{}, errors.Wrap(err, "read metadata")
	}

	log.WithField("file", srcFile).Debug("Reading metadata")

	return m, nil
}

func readSender(srcFile string) (sender.Sender, error) {
	s := sender.Sender{}
	err := s.Read(srcFile)

	if err != nil {
		return sender.Sender{}, errors.Wrap(err, "read sender information")
	}

	log.WithField("file", srcFile).Debug("Reading sender")

	return s, nil
}

func readRecipients(srcFile string) (recipients.RecipientManager, error) {
	recipientManager := recipients.RecipientManager{}

	err := recipientManager.Read(srcFile)
	if err != nil {
		return recipients.RecipientManager{}, errors.Wrap(err, "read recipients files")
	}

	logInfo := log.WithFields(log.Fields{
		"count(valid_recipients)": len(recipientManager.Recipients),
	})

	logInfo.Info("Reading recipients")
	logInfo.WithField("file", srcFile).Debug("Reading recipients")

	return recipientManager, nil
}

func generateLetterInstances(sender sender.Sender, metadata metadata.Metadata, recipients []recipients.Recipient) []letter.Letter {
	var letters []letter.Letter

	for _, r := range recipients {
		lt := letter.New(sender, r, metadata)
		letters = append(letters, lt)
	}

	log.WithFields(log.Fields{
		"count(letters)": len(letters),
	}).Debug("Creating letter instances")

	return letters
}

func readLetterTemplate(srcFile string) (converter.Template, error) {
	template := converter.Template{}

	err := template.Read(srcFile)
	if err != nil {
		return converter.Template{}, errors.Wrap(err, "read template")
	}

	log.WithFields(log.Fields{
		"file": srcFile,
	}).Debug("Reading letter template")

	return template, nil
}

func findAssets(srcDir string) ([]assets.Asset, error) {
	assetRepo := assets.NewRepository(srcDir)

	if err := assetRepo.Init(); err != nil {
		return []assets.Asset{}, errors.Wrap(err, "init asset repo")
	}

	assets := assetRepo.KnownAssets()

	log.WithFields(log.Fields{
		"len(assets)": len(assets),
	}).Debug("Initializing repository for assets")

	return assets, nil
}

func createOutputDirectory() (string, error) {
	cwd, err := os.Getwd()
	dir := filepath.Join(cwd, "letters")
	err = os.MkdirAll(dir, 0755)

	if err != nil {
		return "", errors.Wrap(err, "create output directory")
	}

	return dir, nil
}
