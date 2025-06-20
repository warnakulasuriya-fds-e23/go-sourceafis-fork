package binarizer

import (
	"sync"

	"github.com/warnakulasuriya-fds-e23/go-sourceafis-fork/config"
	"github.com/warnakulasuriya-fds-e23/go-sourceafis-fork/extractor/filter/vote"
	"github.com/warnakulasuriya-fds-e23/go-sourceafis-fork/extractor/logger"
	"github.com/warnakulasuriya-fds-e23/go-sourceafis-fork/primitives"
)

type BinarizedImage struct {
	logger logger.TransparencyLogger
}

func New(logger logger.TransparencyLogger) *BinarizedImage {
	return &BinarizedImage{
		logger: logger,
	}
}

func (b *BinarizedImage) Binarize(input, baseline *primitives.Matrix, mask *primitives.BooleanMatrix, blocks *primitives.BlockMap) (*primitives.BooleanMatrix, error) {
	binarized := primitives.NewBooleanMatrixFromPoint(input.Size())
	var wg sync.WaitGroup
	numWorkers := config.Config.Workers
	workCh := make(chan primitives.IntPoint, numWorkers)
	blockIterator := blocks.Primary.Blocks.Iterator()
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for block := range workCh {
				rect := blocks.Primary.BlockPoint(block)
				for y := rect.Top(); y < rect.Bottom(); y++ {
					for x := rect.Left(); x < rect.Right(); x++ {
						if input.Get(x, y)-baseline.Get(x, y) > 0 {
							binarized.Set(x, y, true)
						}
					}
				}
			}
		}()
	}

	go func() {
		for blockIterator.HasNext() {
			block := blockIterator.Next()
			workCh <- block
		}
		close(workCh)
	}()

	wg.Wait()
	return binarized, b.logger.Log("binarized-image", binarized)
}

func (b *BinarizedImage) Cleanup(binary, mask *primitives.BooleanMatrix) error {

	var wg sync.WaitGroup
	numWorkers := config.Config.Workers

	size := binary.Size()
	inverted := primitives.NewBooleanMatrixFromBooleanMatrix(binary)
	inverted.Invert()
	islands := vote.Apply(inverted, mask, config.Config.BinarizedVoteRadius, config.Config.BinarizedVoteMajority, config.Config.BinarizedVoteBorderDistance)
	holes := vote.Apply(binary, mask, config.Config.BinarizedVoteRadius, config.Config.BinarizedVoteMajority, config.Config.BinarizedVoteBorderDistance)
	processSection := func(startY, endY int) {
		defer wg.Done()
		for y := startY; y < endY; y++ {
			for x := 0; x < size.X; x++ {
				binary.Set(x, y, binary.Get(x, y) && !islands.Get(x, y) || holes.Get(x, y))
			}
		}
	}
	sectionHeight := size.Y / numWorkers
	for i := 0; i < numWorkers; i++ {
		startY := i * sectionHeight
		endY := startY + sectionHeight
		if i == numWorkers-1 {
			endY = size.Y
		}
		wg.Add(1)
		go processSection(startY, endY)
	}

	wg.Wait()
	removeCrosses(binary)
	return b.logger.Log("filtered-binary-image", binary)
}

func removeCrosses(input *primitives.BooleanMatrix) {
	size := input.Size()
	any := true
	for any {
		any = false
		for y := 0; y < size.Y-1; y++ {
			for x := 0; x < size.X-1; x++ {
				if input.Get(x, y) && input.Get(x+1, y+1) && !input.Get(x, y+1) && !input.Get(x+1, y) || input.Get(x, y+1) && input.Get(x+1, y) && !input.Get(x, y) && !input.Get(x+1, y+1) {
					input.Set(x, y, false)
					input.Set(x, y+1, false)
					input.Set(x+1, y, false)
					input.Set(x+1, y+1, false)
					any = true
				}
			}
		}
	}
}

func (b *BinarizedImage) Invert(binary, mask *primitives.BooleanMatrix) *primitives.BooleanMatrix {
	size := binary.Size()
	inverted := primitives.NewBooleanMatrixFromPoint(size)
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			inverted.Set(x, y, !binary.Get(x, y) && mask.Get(x, y))
		}
	}
	return inverted
}
