package domain

type SessionState int

const (
	StateSplash SessionState = iota
	StateMain
	StateConfirmDelete
	StateCloneRepo
	StateSettings
	StateForgeSelect   // choose the forge for repos being cloned
	StateForgeRegistry // manage the forge registry
	StateRepoSettings  // change a repo's forge and tag format
	StateKeybinds      // manage custom keybinds
	StateReposDir      // choose the repo directory before the repo list opens
	StateTagRelease    // tag and push the next version of a repo
)
