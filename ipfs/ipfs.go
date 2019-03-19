package ipfs

import (
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/plugin/loader"
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/repo/fsrepo"
	config "gx/ipfs/QmUAuYuiafnJRZxDDX7MuruMNsicYNuyub5vUeAcupUBNs/go-ipfs-config"
	"os"
	"path/filepath"
)

// LoadPlugins loads all the plugins to the IPFS
func LoadPlugins(basePath string) error {
	pluginpath := filepath.Join(basePath, "plugins")

	// check if repo is accessible before loading plugins
	plugins, err := loader.NewPluginLoader(pluginpath)
	if err != nil {
		return err
	}

	if err := plugins.Initialize(); err != nil {
		return err
	}

	if err := plugins.Inject(); err != nil {
		return err
	}

	return nil
}

// InitRepo initializes the repo if it is not yet
func InitRepo(path string) error {
	// Verifies if the repo is initialized
	isRepoInitialized := fsrepo.IsInitialized(path)

	if isRepoInitialized {
		return nil
	}

	// Loads all the plugins
	LoadPlugins(path)

	// Generates the initial config
	initialConfig, err := config.Init(os.Stdout, 2048)
	if err != nil {
		return err
	}

	initialConfig.Bootstrap = append(initialConfig.Bootstrap, "/ip4/206.189.200.98/tcp/4001/ipfs/QmQ7VQEj6asBAfUbW9XEC2PNMKUz1yWggSgKmtUsbYN6rt")

	// Initializes the repo
	if err := fsrepo.Init(path, initialConfig); err != nil {
		return err
	}

	return nil
}
