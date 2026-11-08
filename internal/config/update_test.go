package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeConfigFile(t *testing.T, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), ".gitkeeper.json")
	t.Setenv("GITKEEPER_CONFIG", path)
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadMarksUnparseableConfigAsFailed(t *testing.T) {
	writeConfigFile(t, `{"forges": [`)

	cfg := Load()
	if !cfg.LoadFailed {
		t.Fatal("LoadFailed = false for a config that doesn't parse")
	}
	if len(cfg.Warnings) == 0 {
		t.Error("expected a warning explaining the failure")
	}
	if err := Save(cfg); err == nil {
		t.Error("Save wrote defaults over a config that failed to load")
	}
}

func TestUpdateRefusesToOverwriteUnparseableConfig(t *testing.T) {
	const broken = `{"forges": [`
	path := writeConfigFile(t, broken)

	_, err := Update(func(c *Config) { c.Tools.Shell = "zsh" })
	if err == nil {
		t.Fatal("Update succeeded over a config that doesn't parse")
	}
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(data) != broken {
		t.Errorf("file was rewritten to %q", data)
	}
}

func TestUpdateKeepsEditsMadeByAnotherWriter(t *testing.T) {
	writeConfigFile(t, `{}`)
	stale := Load()

	if _, err := Update(func(c *Config) {
		c.Repos["ns/repo"] = RepoLink{TagPrefix: "written elsewhere/v"}
	}); err != nil {
		t.Fatal(err)
	}

	got, err := Update(func(c *Config) { c.Tools.Shell = "zsh" })
	if err != nil {
		t.Fatal(err)
	}
	if got.Repos["ns/repo"].TagPrefix != "written elsewhere/v" {
		t.Errorf("tag prefix lost: %+v", got.Repos["ns/repo"])
	}
	if stale.Tools.Shell == "zsh" {
		t.Error("Update mutated a previously loaded config")
	}
}

func TestUpdateCreatesMissingConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", ".gitkeeper.json")
	t.Setenv("GITKEEPER_CONFIG", path)

	if _, err := Update(func(c *Config) { c.Tools.Shell = "zsh" }); err != nil {
		t.Fatal(err)
	}
	if Load().Tools.Shell != "zsh" {
		t.Error("change not persisted")
	}
}

func TestLoadNormalizesReposDir(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	writeConfigFile(t, `{"reposDir": "~/repos/../gitkeeper/"}`)

	if got, want := Load().ReposDir, filepath.Join(home, "gitkeeper"); got != want {
		t.Errorf("ReposDir = %q, want %q", got, want)
	}
}

func TestLoadRejectsUnknownDefaultProtocol(t *testing.T) {
	writeConfigFile(t, `{"behavior": {"defaultProtocol": "SSH"}}`)
	if got := Load().Behavior.DefaultProtocol; got != SchemeSSH {
		t.Errorf("DefaultProtocol = %q, want ssh", got)
	}

	writeConfigFile(t, `{"behavior": {"defaultProtocol": "ftp"}}`)
	cfg := Load()
	if cfg.Behavior.DefaultProtocol != SchemeHTTPS {
		t.Errorf("DefaultProtocol = %q, want https", cfg.Behavior.DefaultProtocol)
	}
	if len(cfg.Warnings) != 1 {
		t.Errorf("Warnings = %q, want one for the unknown protocol", cfg.Warnings)
	}
}

func TestConfigCloneSharesNoState(t *testing.T) {
	original := Default()
	original.Forges = []Forge{{Name: "github", Host: "github.com"}}
	original.Repos["ns/repo"] = RepoLink{Forge: "github"}

	clone := original.Clone()
	clone.Forges[0].Name = "changed"
	clone.Repos["ns/other"] = RepoLink{}

	if original.Forges[0].Name != "github" {
		t.Error("clone shares Forges")
	}
	if _, ok := original.Repos["ns/other"]; ok {
		t.Error("clone shares the Repos map")
	}
}

func TestRepoLinkWithForgeLinkKeepsTagFormat(t *testing.T) {
	stored := RepoLink{Forge: "a", TagPrefix: "v", TagSuffix: "-pkg"}

	got := stored.WithForgeLink(RepoLink{Forge: "b"})
	if got.Forge != "b" || got.TagPrefix != "v" || got.TagSuffix != "-pkg" {
		t.Errorf("got %+v, want the tag format kept", got)
	}
}
