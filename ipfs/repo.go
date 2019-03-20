package ipfs

import (
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/plugin/loader"
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/repo"
	"gx/ipfs/QmPDEJTb3WBHmvubsLXCaqRPC8dRgvFz7A4p96dxZbJuWL/go-ipfs/repo/fsrepo"
	config "gx/ipfs/QmUAuYuiafnJRZxDDX7MuruMNsicYNuyub5vUeAcupUBNs/go-ipfs-config"
	"os"
	"path/filepath"
)

// loadPlugins loads all the plugins to the IPFS
func loadPlugins(basePath string) error {
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

// initRepo initializes the repo if it is not yet
func initRepo(path string) error {
	// Verifies if the repo is initialized
	if isRepoInitialized := isRepoInitialized(path); isRepoInitialized {
		return nil
	}

	// Loads all the plugins
	loadPlugins(path)

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

// openRepo Opens the repo
func openRepo(path string) (repo.Repo, error) {
	// Gets the repo
	nodeRepo, err := fsrepo.Open(path)
	if err != nil {
		return nil, err
	}

	return nodeRepo, nil
}

// isRepoInitialized return is the repo is initialized
func isRepoInitialized(path string) bool {
	return !fsrepo.IsInitialized(path)
}
