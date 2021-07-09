package sealing

import (
	"bufio"
	"context"
	"io"
	"os"

	"github.com/filecoin-project/go-commp-utils/ffiwrapper"
	"github.com/filecoin-project/go-padreader"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-car"
	xerrors "golang.org/x/xerrors"
)

func (m *Sealing) AssignPieceIntoAnyRawSectors(ctx context.Context, carfile string) (*abi.PieceInfo, error) {

	m.inputLk.Lock()

	pcid, psize, err := m.generatePieceCid(ctx, carfile)
	if err != nil {
		return nil, xerrors.Errorf("generate piece cid %s", err)
	}

	if pcid == nil || psize == 0 {
		return nil, xerrors.Errorf("cannot determine piece info")
	}

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

	pieceReader, _ := padreader.New(rdr, uint64(stat.Size()))

	m.pendingPieces[*pcid] = &pendingPiece{
		size:     psize,
		data:     pieceReader,
		assigned: false,
		accepted: func(sn abi.SectorNumber, offset abi.UnpaddedPieceSize, err error) {},
	}

	sp, err := m.currentSealProof(ctx)
	if err != nil {
		return nil, xerrors.Errorf("getting current seal proof type: %w", err)
	}

	go func() {
		defer m.inputLk.Unlock()
		if err := m.updateInput(ctx, sp); err != nil {
			log.Errorf("%+v", err)
		}
	}()

	return &abi.PieceInfo{
		Size:     psize.Padded(),
		PieceCID: *pcid,
	}, nil
}

func (m *Sealing) generatePieceCid(ctx context.Context, path string) (*cid.Cid, abi.UnpaddedPieceSize, error) {
	rdr, err := os.Open(path)
	if err != nil {
		return nil, 0, err
	}
	defer rdr.Close() //nolint:errcheck

	stat, err := rdr.Stat()
	if err != nil {
		return nil, 0, err
	}

	// check that the data is a car file; if it's not, retrieval won't work
	_, _, err = car.ReadHeader(bufio.NewReader(rdr))
	if err != nil {
		return nil, 0, xerrors.Errorf("not a car file: %w", err)
	}

	if _, err := rdr.Seek(0, io.SeekStart); err != nil {
		return nil, 0, xerrors.Errorf("seek to start: %w", err)
	}

	pieceReader, pieceSize := padreader.New(rdr, uint64(stat.Size()))
	sp, err := m.currentSealProof(ctx)
	if err != nil {
		return nil, 0, xerrors.Errorf("getting current seal proof type: %w", err)
	}

	ssize, err := sp.SectorSize()
	if err != nil {
		return nil, 0, err
	}

	if pieceSize > abi.PaddedPieceSize(ssize).Unpadded() {
		return nil, 0, xerrors.Errorf("piece cannot fit into a sector")
	}

	commP, err := ffiwrapper.GeneratePieceCIDFromFile(sp, pieceReader, pieceSize)
	if err != nil {
		return nil, 0, err
	}
	return &commP, pieceSize, nil
}
