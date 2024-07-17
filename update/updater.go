package update

import (
	"context"
	// "derhasi/awesome-simracing-client/backend"
	// "derhasi/awesome-simracing-client/backend/checks"
	// "derhasi/awesome-simracing-client/backend/settings"
	// "derhasi/awesome-simracing-client/backend/util"
	"os"

	"github.com/minio/selfupdate"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Updater struct {
	info     *backend.BuildInfo
	settings *settings.Settings
	http     *util.HttpHelper
	exe      util.Executable
	ctx      context.Context
}

func NewUpdater(build *backend.BuildInfo, settings *settings.Settings, http *util.HttpHelper) *Updater {
	return &Updater{info: build, settings: settings, http: http}
}

func (u *Updater) WailsOnStartup(ctx context.Context) {
	u.ctx = ctx
}
func (u *Updater) WailsBind() *Bindings {
	return &Bindings{updater: u}
}

func (u *Updater) update() error {

	executable, err := util.GetExecutable()
	if err != nil {
		return err
	}

	runtime.LogInfof(u.ctx, "Executable found %v", executable)

	// Loads info about the latest release.
	release := checks.GetRelease(u.settings.PreRelease, u.http)

	// Check if there is a binary file at all.
	_, err = release.GetBinary()
	if err != nil {
		return err
	}

	// Download all the assets.
	for _, file := range release.Files {
		err = u.processFile(file, executable)
		if err != nil {
			runtime.LogErrorf(u.ctx, "Could not process %s: %s", file.Name, err.Error())
			// Only quit on binary file errors.
			if file.IsBinary() {
				return err
			}
		}
	}

	return nil
}

func prepareBackupDir(info *backend.BuildInfo, exec util.Executable) (string, error) {

	dir := exec.Dir + "/old_versions/" + info.Version

	err := os.MkdirAll(dir, os.ModeDir|os.ModePerm)
	if err != nil {
		return "", err
	}

	return dir, nil
}

func (u *Updater) processFile(file checks.ReleaseFile, executable util.Executable) error {

	target := executable.Dir + "/" + file.Name
	if file.IsBinary() {
		target = executable.Exec
	}

	exists, err := util.LocalFileExists(target)
	if err != nil {
		return err
	}

	if exists {
		return u.updateFile(target, file, executable)
	}
	return util.DownloadFile(target, file.Url)
}

func (u *Updater) updateFile(target string, file checks.ReleaseFile, executable util.Executable) error {
	backupDir, err := prepareBackupDir(u.info, executable)
	if err != nil {
		return err
	}

	source, err := u.http.HttpGet(file.Url, false)
	if err != nil {
		return err
	}
	defer source.Body.Close()

	opt := selfupdate.Options{
		TargetPath:  target,
		OldSavePath: backupDir + "/" + file.Name,
	}

	return selfupdate.Apply(source.Body, opt)
}

// see https://github.com/wailsapp/wails/discussions/2223
func (u *Updater) restart() error {
	err := util.RestartSelf()
	if err != nil {
		runtime.LogErrorf(u.ctx, "Could not restart: %s", err.Error())
	}
	return err
}
