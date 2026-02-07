package main

import (
	"context"
	"math"
	"os"

	"github.com/itchio/headway/state"
	"github.com/itchio/lake/pools/fspool"
	"github.com/itchio/savior/filesource"
	"github.com/itchio/wharf/pwr"
	"github.com/itchio/wharf/pwr/bowl"
	"github.com/itchio/wharf/pwr/patcher"

	_ "github.com/itchio/wharf/decompressors/cbrotli"
)

const TOTAL_PROG = 102;
const STEP_PATCH = 100;
const STEP_SIG = 101;
const STEP_GHOST = 102;

type progressSaveConsumer struct {
	shouldSave func() bool
	save       func(checkpoint *patcher.Checkpoint) (patcher.AfterSaveAction, error)
}

func (saveConsumer *progressSaveConsumer) ShouldSave() bool {
	return saveConsumer.shouldSave()
}

func (saveConsumer *progressSaveConsumer) Save(checkpoint *patcher.Checkpoint) (patcher.AfterSaveAction, error) {
	return saveConsumer.save(checkpoint)
}


func applyPatch(source string, destination string, patchFilename string, signatureFilename string, onProgress func(done int64, total int64)) error {

	patchReader, err := filesource.Open(patchFilename);
	if err != nil {
		return err;
	}
	defer patchReader.Close();

	os.MkdirAll(destination, 0755);

	consumer := &state.Consumer{
		OnProgress: func(w float64) {
			total := int64(TOTAL_PROG);
			done := int64(math.Ceil(w * float64(STEP_PATCH)));

			onProgress(done, total);
		},
	}

	p, err := patcher.New(patchReader, consumer);
	if err != nil {
		os.RemoveAll(destination);
		return err;
	}

	targetPool := fspool.New(p.GetTargetContainer(), source);

	b, err := bowl.NewFreshBowl(bowl.FreshBowlParams{
		SourceContainer: p.GetSourceContainer(),
				    TargetContainer: p.GetTargetContainer(),
				    TargetPool: targetPool,
				    OutputFolder: destination,
	});

	if err != nil {
		os.RemoveAll(destination);
		return err;
	}

	// dont actually care about checkpoints, just use them to track progress
	p.SetSaveConsumer(&progressSaveConsumer{
		shouldSave: func() bool {
			consumer.Progress(p.Progress());
			return false;
		},
		save: func(checkpoint *patcher.Checkpoint) (patcher.AfterSaveAction, error) {
			return patcher.AfterSaveContinue, nil;
		},
	})

	err = p.Resume(nil, targetPool, b);
	if err != nil {
		os.RemoveAll(destination);
		return err;
	}

	sigSource, err := filesource.Open(signatureFilename);
	if err != nil {
		os.RemoveAll(destination);
		return err;
	}

	onProgress(STEP_SIG, TOTAL_PROG);
	signatureInfo, err := pwr.ReadSignature(context.Background(), sigSource);
	err = pwr.AssertValid(destination, signatureInfo);
	if err != nil {
		os.RemoveAll(destination);
		return err;
	}

	onProgress(STEP_GHOST, TOTAL_PROG);
	err = pwr.AssertNoGhosts(destination, signatureInfo);
	if err != nil {
		os.RemoveAll(destination);
		return err;
	}

	return nil;
}
