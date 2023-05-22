package installer

import (
	"github.com/NubeIO/lib-systemctl-go/systemctl"
	"path"
)

const fileMode = 0755
const defaultTimeout = 30

type Installer struct {
	RootDir         string // /data
	StoreDir        string // <root_dir>/store
	TmpDir          string // /data/tmp
	BackupDir       string // <root_dir>/backup
	FileMode        int    // 0755
	DefaultTimeout  int    // 30
	AppsDownloadDir string // <root_dir>/rubix-service/apps/download
	AppsInstallDir  string // <root_dir>/rubix-service/apps/install
	SystemCtl       *systemctl.SystemCtl
}

func New(app *Installer) *Installer {
	if app == nil {
		app = &Installer{}
	}
	if app.RootDir == "" {
		app.RootDir = "/data"
	}
	if app.FileMode == 0 {
		app.FileMode = fileMode
	}
	if app.DefaultTimeout == 0 {
		app.DefaultTimeout = defaultTimeout
	}
	if app.StoreDir == "" {
		app.StoreDir = path.Join(app.RootDir, "store")
	}
	if app.TmpDir == "" {
		app.TmpDir = path.Join(app.RootDir, "tmp")
	}
	if app.BackupDir == "" {
		app.BackupDir = path.Join(app.RootDir, "backup")
	}
	if app.AppsDownloadDir == "" {
		app.AppsDownloadDir = path.Join(app.RootDir, "rubix-service/apps/download")
	}
	if app.AppsInstallDir == "" {
		app.AppsInstallDir = path.Join(app.RootDir, "rubix-service/apps/install")
	}
	app.SystemCtl = systemctl.New(false, app.DefaultTimeout)
	return app
}
