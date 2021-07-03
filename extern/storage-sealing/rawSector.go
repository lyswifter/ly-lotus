package sealing

import (
	"bufio"
	"context"
	"io"
	"os"

	"github.com/filecoin-project/go-commp-utils/ffiwrapper"
	"github.com/filecoin-project/go-padreader"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipld/go-car"
	xerrors "golang.org/x/xerrors"
)

func (m *Sealing) AssignPieceIntoAnyRawSectors(ctx context.Context, carfile string) (*abi.PieceInfo, error) {

	m.inputLk.Lock()

	rdr, err := os.Open(carfile)
	if err != nil {
		return nil, err
	}
	// defer rdr.Close() //nolint:errcheck

	stat, err := rdr.Stat()
	if err != nil {
		return nil, err
	}

	// check that the data is a car file; if it's not, retrieval won't work
	_, _, err = car.ReadHeader(bufio.NewReader(rdr))
	if err != nil {
		return nil, xerrors.Errorf("not a car file: %w", err)
	}

	if _, err := rdr.Seek(0, io.SeekStart); err != nil {
		return nil, xerrors.Errorf("seek to start: %w", err)
	}

	pieceReader, pieceSize := padreader.New(rdr, uint64(stat.Size()))

	if (padreader.PaddedSize(uint64(pieceSize))) != pieceSize {
		return nil, xerrors.Errorf("cannot allocate unpadded piece")
	}

	sp, err := m.currentSealProof(ctx)
	if err != nil {
		return nil, xerrors.Errorf("getting current seal proof type: %w", err)
	}

	ssize, err := sp.SectorSize()
	if err != nil {
		return nil, err
	}

	if pieceSize > abi.PaddedPieceSize(ssize).Unpadded() {
		return nil, xerrors.Errorf("piece cannot fit into a sector")
	}

	commP, err := ffiwrapper.GeneratePieceCIDFromFile(sp, pieceReader, pieceSize)
	if err != nil {
		return nil, err
	}

	log.Infof("Assigning piece for piecesize: %d pad-size: %d", pieceSize, pieceSize.Padded())

	m.pendingPieces[commP] = &pendingPiece{
		size:     pieceSize,
		data:     pieceReader,
		assigned: false,
		accepted: func(sn abi.SectorNumber, offset abi.UnpaddedPieceSize, err error) {},
	}

	go func() {
		defer m.inputLk.Unlock()
		if err := m.updateInput(ctx, sp); err != nil {
			log.Errorf("%+v", err)
		}
	}()

	return &abi.PieceInfo{
		Size:     pieceSize.Padded(),
		PieceCID: commP,
	}, nil
}
