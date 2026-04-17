package syncer

import (
	"context"
	"log/slog"
	"slices"
	"time"

	"github.com/isyuricunha/gitflect/internal/config"
	"github.com/isyuricunha/gitflect/internal/git"
	"github.com/isyuricunha/gitflect/internal/provider"
	"github.com/isyuricunha/gitflect/internal/provider/github"
	"github.com/isyuricunha/gitflect/internal/provider/gitlab"
)

type Syncer struct {
	cfg     *config.Config
	src     provider.Source
	dst     provider.Destination
	logger  *slog.Logger
	secrets []string
}

func New(cfg *config.Config, logger *slog.Logger) *Syncer {
	var src provider.Source
	if cfg.SourceProvider == "github" {
		src = github.New(cfg.SourceToken, cfg.SourceUser)
	}

	var dst provider.Destination
	if cfg.DestProvider == "gitlab" {
		dst = gitlab.New(cfg.DestToken, cfg.DestUser, cfg.DestURL)
	}

	// Used to sanitize logs
	secrets := []string{cfg.SourceToken, cfg.DestToken}

	return &Syncer{
		cfg:     cfg,
		src:     src,
		dst:     dst,
		logger:  logger,
		secrets: secrets,
	}
}

func (s *Syncer) Run(_ context.Context) {
	if s.src == nil || s.dst == nil {
		s.logger.Error("Provider is nil. Check configuration (SourceProvider / DestProvider)")
		return
	}

	s.logger.Info("Fetching repository list...")
	repos, err := s.src.ListRepos(s.cfg.Visibility)
	if err != nil {
		s.logger.Error("Failed to list source repos", "err", err)
		return
	}

	s.logger.Info("Found repositories", "count", len(repos))

	for _, repo := range repos {
		if !s.shouldSync(repo.Name) {
			s.logger.Info("Skipping excluded repository", "repo", repo.Name)
			continue
		}

		s.logger.Info("Syncing repository", "repo", repo.Name)
		pushURL, err := s.dst.EnsureRepo(repo.Name, repo.Private)
		if err != nil {
			s.logger.Error("Failed to ensure destination repo", "repo", repo.Name, "err", err)
			continue
		}

		start := time.Now()
		if err := git.Mirror(repo.CloneURL, pushURL, repo.Name, s.secrets); err != nil {
			s.logger.Error("Mirroring failed", "repo", repo.Name, "err", err)
			continue
		}

		s.logger.Info("Successfully synced", "repo", repo.Name, "duration", time.Since(start).String())
	}
}

func (s *Syncer) shouldSync(name string) bool {
	if len(s.cfg.Include) > 0 {
		return slices.Contains(s.cfg.Include, name)
	}
	if len(s.cfg.Exclude) > 0 {
		return !slices.Contains(s.cfg.Exclude, name)
	}
	return true
}
