package smoothing

import (
	"github.com/warnakulasuriya-fds-e23/go-sourceafis-fork/config"
	"github.com/warnakulasuriya-fds-e23/go-sourceafis-fork/extractor/logger"
	"github.com/warnakulasuriya-fds-e23/go-sourceafis-fork/primitives"
)

type OrientedSmoothing struct {
	logger logger.TransparencyLogger
}

func New(logger logger.TransparencyLogger) *OrientedSmoothing {
	return &OrientedSmoothing{
		logger: logger,
	}
}

func (s *OrientedSmoothing) Parallel(input, orientation *primitives.Matrix, mask *primitives.BooleanMatrix, blocks *primitives.BlockMap) (*primitives.Matrix, error) {
	lines := lines(config.Config.ParallelSmoothingResolution, config.Config.ParallelSmoothingRadius, config.Config.ParallelSmoothingStep)
	smoothed := smooth(input, orientation, mask, blocks, 0, lines)

	return smoothed, s.logger.Log("parallel-smoothing", smoothed)
}

func (s *OrientedSmoothing) Orthogonal(input, orientation *primitives.Matrix, mask *primitives.BooleanMatrix, blocks *primitives.BlockMap) (*primitives.Matrix, error) {
	lines := lines(config.Config.OrthogonalSmoothingResolution, config.Config.OrthogonalSmoothingRadius, config.Config.OrthogonalSmoothingStep)
	smoothed := smooth(input, orientation, mask, blocks, primitives.Pi, lines)

	return smoothed, s.logger.Log("orthogonal-smoothing", smoothed)
}

func lines(resolution, radius int, step float64) [][]primitives.IntPoint {
	result := make([][]primitives.IntPoint, resolution)
	for orientationIndex := 0; orientationIndex < resolution; orientationIndex++ {
		line := []primitives.IntPoint{primitives.ZeroIntPoint()}
		direction := primitives.BucketCenter(orientationIndex, resolution).FromOrientation().ToVector()
		for r := float64(radius); r >= 0.5; r /= step {
			sample := direction.Multiply(r).Round()
			var isFound bool
			for _, samp := range line {
				if samp.Equals(sample) {
					isFound = true
				}
			}
			if !isFound {
				line = append(line, sample)
				line = append(line, sample.Negate())
			}
		}
		result[orientationIndex] = line
	}
	return result
}

func smooth(input, orientation *primitives.Matrix, mask *primitives.BooleanMatrix, blocks *primitives.BlockMap, angle float64, lines [][]primitives.IntPoint) *primitives.Matrix {
	output := primitives.NewMatrixFromPoint(input.Size())
	it := blocks.Primary.Blocks.Iterator()
	for it.HasNext() {
		block := it.Next()
		if mask.GetPoint(block) {
			line := lines[primitives.AngleAdd(orientation.GetPoint(block), angle).Quantize(len(lines))]
			for _, linePoint := range line {
				target := blocks.Primary.BlockPoint(block)
				source := target.Move(linePoint).Intersect(primitives.IntRect{
					X:      0,
					Y:      0,
					Width:  blocks.Pixels.X,
					Height: blocks.Pixels.Y,
				})
				target = source.Move(linePoint.Negate())
				for y := target.Top(); y < target.Bottom(); y++ {
					for x := target.Left(); x < target.Right(); x++ {
						output.Add(x, y, input.Get(x+linePoint.X, y+linePoint.Y))
					}
				}
			}
			blockArea := blocks.Primary.BlockPoint(block)
			for y := blockArea.Top(); y < blockArea.Bottom(); y++ {
				for x := blockArea.Left(); x < blockArea.Right(); x++ {
					output.Multiply(x, y, 1.0/float64(len(line)))
				}
			}
		}
	}

	return output
}
