package validator

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"social-network/services/media/internal/configs"
)

var (
	ErrInvalidImage     = errors.New("invalid image")
	ErrImageTooLarge    = errors.New("image exceeds size limit")
	ErrUnsupportedType  = errors.New("unsupported image type")
	ErrInvalidDimension = errors.New("invalid image dimensions")
)

type ImageValidator struct {
	Config configs.FileConstraints
}

func (v *ImageValidator) ValidateImage(ctx context.Context, r io.Reader) error {
	// 1️⃣ Enforce max upload size
	limited := io.LimitReader(r, v.Config.MaxImageUpload+1)
	buf, err := io.ReadAll(limited)
	if err != nil {
		return fmt.Errorf("read failed: %w", err)
	}

	if int64(len(buf)) > v.Config.MaxImageUpload {
		return ErrImageTooLarge
	}

	// 2️⃣ Detect MIME type by content (NOT filename)
	mime := http.DetectContentType(buf[:min(512, len(buf))])
	if !v.Config.AllowedMIMEs[mime] {
		return fmt.Errorf("%w: %s", ErrUnsupportedType, mime)
	}

	// 3️⃣ Decode config only (fast, safe)
	cfg, format, err := image.DecodeConfig(bytes.NewReader(buf))
	if err != nil {
		return ErrInvalidImage
	}

	// 4️⃣ Dimension validation (prevents decompression bombs)
	if cfg.Width <= 0 || cfg.Height <= 0 {
		return ErrInvalidDimension
	}

	if cfg.Width > v.Config.MaxWidth || cfg.Height > v.Config.MaxHeight {
		return fmt.Errorf(
			"%w: %dx%d",
			ErrInvalidDimension,
			cfg.Width,
			cfg.Height,
		)
	}

	// 5️⃣ Optional: restrict formats (jpeg/png/gif/etc)
	if !v.Config.AllowedMIMEs["image/"+format] {
		return fmt.Errorf("%w: %s", ErrUnsupportedType, format)
	}

	return nil
}
