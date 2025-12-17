package application

import (
	"context"
	"log"
	"time"

	ct "social-network/shared/go/customtypes"
)

// StartVariantWorker starts a background worker that periodically processes pending file variants
func (m *MediaService) StartVariantWorker(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := m.processPendingVariants(ctx); err != nil {
					log.Printf("Error processing pending variants: %v", err)
				}
			case <-ctx.Done():
				log.Println("Variant worker stopped")
				return
			}
		}
	}()
}

// processPendingVariants queries for file_variants with status 'pending' and calls GenerateVariant for each
func (m *MediaService) processPendingVariants(ctx context.Context) error {
	variants, err := m.Queries.GetPendingVariants(ctx)
	if err != nil {
		return err
	}

	for _, v := range variants {

		// Call GenerateVariant
		size, err := m.Clients.GenerateVariant(ctx, v.Bucket, v.ObjectKey, v.Variant)
		if err != nil {
			log.Printf("Failed to generate variant for file %d variant %s: %v", v.Id, v.Variant, err)
			// Update status to failed
			if updateErr := m.Queries.UpdateVariantStatusAndSize(ctx, v.Id, ct.Failed, size); updateErr != nil {
				log.Printf("Failed to update status to failed: %v", updateErr)
			}
		} else {
			log.Printf("Successfully generated variant for file %d variant %s", v.Id, v.Variant)
			// Update status to complete
			if updateErr := m.Queries.UpdateVariantStatusAndSize(ctx, v.Id, ct.Complete, size); updateErr != nil {
				log.Printf("Failed to update status to complete: %v", updateErr)
			}
		}
	}

	return nil
}
