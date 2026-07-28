package core_test

import (
	"context"
	"testing"
	"time"

	"github.com/NordCoder/unmatched-web/internal/persistence"
	coreruntime "github.com/NordCoder/unmatched-web/internal/runtime"
)

func TestPostgresPersistenceCorrections(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	pool := openIsolatedPostgres(t, ctx)
	if err := persistence.ApplyMigrations(ctx, pool); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}

	registry := coreruntime.NewMemoryDefinitionRegistry()
	bundle := syntheticBundle()
	if err := registry.Register("synthetic@v1", bundle); err != nil {
		t.Fatalf("register definitions: %v", err)
	}

	t.Run("canonical identity survives reconstruction", func(t *testing.T) {
		assertCanonicalIdentitySurvivesReconstruction(t, ctx, pool, registry)
	})
	t.Run("saturated pool cannot deadlock command owner", func(t *testing.T) {
		assertSaturatedPoolProgress(t, ctx, pool, registry)
	})
	t.Run("authority mutation rolls back with late failure", func(t *testing.T) {
		assertAuthorityRollbackIsAtomic(t, ctx, pool, registry, bundle.Ref)
	})
}
