package sealing

import (
	"context"

	"golang.org/x/xerrors"

	abi "github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/specs-storage/storage"
)

func (m *Sealing) PledgeSector(ctx context.Context, data storage.Data, pieceInfo *abi.PieceInfo) (storage.SectorRef, error) {
	m.startupWait.Wait()

	m.inputLk.Lock()
	defer m.inputLk.Unlock()

	cfg, err := m.getConfig()
	if err != nil {
		return storage.SectorRef{}, xerrors.Errorf("getting config: %w", err)
	}

	if cfg.MaxSealingSectors > 0 {
		if m.stats.curSealing() >= cfg.MaxSealingSectors {
			return storage.SectorRef{}, xerrors.Errorf("too many sectors sealing (curSealing: %d, max: %d)", m.stats.curSealing(), cfg.MaxSealingSectors)
		}
	}

	spt, err := m.currentSealProof(ctx)
	if err != nil {
		return storage.SectorRef{}, xerrors.Errorf("getting seal proof type: %w", err)
	}

	sid, err := m.createSector(ctx, cfg, spt)
	if err != nil {
		return storage.SectorRef{}, err
	}

	if pieceInfo != nil && data != nil {
		sp, err := m.currentSealProof(ctx)
		if err != nil {
			return storage.SectorRef{}, xerrors.Errorf("getting current seal proof type: %w", err)
		}

		ssize, err := sp.SectorSize()
		if err != nil {
			return storage.SectorRef{}, err
		}

		size := pieceInfo.Size.Unpadded()

		if size > abi.PaddedPieceSize(ssize).Unpadded() {
			return storage.SectorRef{}, xerrors.Errorf("piece cannot fit into a sector")
		}

		log.Infof("Creating CC sector: %d with piece data: %v inside", sid, pieceInfo.PieceCID)
		m.pendingPieces[pieceInfo.PieceCID] = &pendingPiece{
			size:     size,
			data:     data,
			assigned: false,
			accepted: func(sn abi.SectorNumber, offset abi.UnpaddedPieceSize, err error) {},
		}

		return m.minerSector(spt, sid), m.sectors.Send(uint64(sid), SectorStart{
			ID:         sid,
			SectorType: spt,
		})
	}

	log.Infof("Creating CC sector %d", sid)
	return m.minerSector(spt, sid), m.sectors.Send(uint64(sid), SectorStartCC{
		ID:         sid,
		SectorType: spt,
	})
}
