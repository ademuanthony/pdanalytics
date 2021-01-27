// Code generated by SQLBoiler 3.7.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/randomize"
	"github.com/volatiletech/sqlboiler/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testVoteReceiveTimeDeviations(t *testing.T) {
	t.Parallel()

	query := VoteReceiveTimeDeviations()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testVoteReceiveTimeDeviationsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VoteReceiveTimeDeviation{}
	if err = randomize.Struct(seed, o, voteReceiveTimeDeviationDBTypes, true, voteReceiveTimeDeviationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VoteReceiveTimeDeviation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := VoteReceiveTimeDeviations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testVoteReceiveTimeDeviationsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VoteReceiveTimeDeviation{}
	if err = randomize.Struct(seed, o, voteReceiveTimeDeviationDBTypes, true, voteReceiveTimeDeviationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VoteReceiveTimeDeviation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := VoteReceiveTimeDeviations().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := VoteReceiveTimeDeviations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testVoteReceiveTimeDeviationsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VoteReceiveTimeDeviation{}
	if err = randomize.Struct(seed, o, voteReceiveTimeDeviationDBTypes, true, voteReceiveTimeDeviationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VoteReceiveTimeDeviation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := VoteReceiveTimeDeviationSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := VoteReceiveTimeDeviations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testVoteReceiveTimeDeviationsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VoteReceiveTimeDeviation{}
	if err = randomize.Struct(seed, o, voteReceiveTimeDeviationDBTypes, true, voteReceiveTimeDeviationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VoteReceiveTimeDeviation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := VoteReceiveTimeDeviationExists(ctx, tx, o.BlockTime, o.Bin)
	if err != nil {
		t.Errorf("Unable to check if VoteReceiveTimeDeviation exists: %s", err)
	}
	if !e {
		t.Errorf("Expected VoteReceiveTimeDeviationExists to return true, but got false.")
	}
}

func testVoteReceiveTimeDeviationsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VoteReceiveTimeDeviation{}
	if err = randomize.Struct(seed, o, voteReceiveTimeDeviationDBTypes, true, voteReceiveTimeDeviationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VoteReceiveTimeDeviation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	voteReceiveTimeDeviationFound, err := FindVoteReceiveTimeDeviation(ctx, tx, o.BlockTime, o.Bin)
	if err != nil {
		t.Error(err)
	}

	if voteReceiveTimeDeviationFound == nil {
		t.Error("want a record, got nil")
	}
}

func testVoteReceiveTimeDeviationsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VoteReceiveTimeDeviation{}
	if err = randomize.Struct(seed, o, voteReceiveTimeDeviationDBTypes, true, voteReceiveTimeDeviationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VoteReceiveTimeDeviation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = VoteReceiveTimeDeviations().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testVoteReceiveTimeDeviationsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VoteReceiveTimeDeviation{}
	if err = randomize.Struct(seed, o, voteReceiveTimeDeviationDBTypes, true, voteReceiveTimeDeviationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VoteReceiveTimeDeviation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := VoteReceiveTimeDeviations().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testVoteReceiveTimeDeviationsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	voteReceiveTimeDeviationOne := &VoteReceiveTimeDeviation{}
	voteReceiveTimeDeviationTwo := &VoteReceiveTimeDeviation{}
	if err = randomize.Struct(seed, voteReceiveTimeDeviationOne, voteReceiveTimeDeviationDBTypes, false, voteReceiveTimeDeviationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VoteReceiveTimeDeviation struct: %s", err)
	}
	if err = randomize.Struct(seed, voteReceiveTimeDeviationTwo, voteReceiveTimeDeviationDBTypes, false, voteReceiveTimeDeviationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VoteReceiveTimeDeviation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = voteReceiveTimeDeviationOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = voteReceiveTimeDeviationTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := VoteReceiveTimeDeviations().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testVoteReceiveTimeDeviationsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	voteReceiveTimeDeviationOne := &VoteReceiveTimeDeviation{}
	voteReceiveTimeDeviationTwo := &VoteReceiveTimeDeviation{}
	if err = randomize.Struct(seed, voteReceiveTimeDeviationOne, voteReceiveTimeDeviationDBTypes, false, voteReceiveTimeDeviationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VoteReceiveTimeDeviation struct: %s", err)
	}
	if err = randomize.Struct(seed, voteReceiveTimeDeviationTwo, voteReceiveTimeDeviationDBTypes, false, voteReceiveTimeDeviationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VoteReceiveTimeDeviation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = voteReceiveTimeDeviationOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = voteReceiveTimeDeviationTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := VoteReceiveTimeDeviations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testVoteReceiveTimeDeviationsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VoteReceiveTimeDeviation{}
	if err = randomize.Struct(seed, o, voteReceiveTimeDeviationDBTypes, true, voteReceiveTimeDeviationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VoteReceiveTimeDeviation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := VoteReceiveTimeDeviations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testVoteReceiveTimeDeviationsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VoteReceiveTimeDeviation{}
	if err = randomize.Struct(seed, o, voteReceiveTimeDeviationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize VoteReceiveTimeDeviation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(voteReceiveTimeDeviationColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := VoteReceiveTimeDeviations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testVoteReceiveTimeDeviationsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VoteReceiveTimeDeviation{}
	if err = randomize.Struct(seed, o, voteReceiveTimeDeviationDBTypes, true, voteReceiveTimeDeviationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VoteReceiveTimeDeviation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = o.Reload(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testVoteReceiveTimeDeviationsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VoteReceiveTimeDeviation{}
	if err = randomize.Struct(seed, o, voteReceiveTimeDeviationDBTypes, true, voteReceiveTimeDeviationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VoteReceiveTimeDeviation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := VoteReceiveTimeDeviationSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testVoteReceiveTimeDeviationsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VoteReceiveTimeDeviation{}
	if err = randomize.Struct(seed, o, voteReceiveTimeDeviationDBTypes, true, voteReceiveTimeDeviationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VoteReceiveTimeDeviation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := VoteReceiveTimeDeviations().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	voteReceiveTimeDeviationDBTypes = map[string]string{`Bin`: `character varying`, `BlockHeight`: `bigint`, `BlockTime`: `bigint`, `ReceiveTimeDifference`: `double precision`}
	_                               = bytes.MinRead
)

func testVoteReceiveTimeDeviationsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(voteReceiveTimeDeviationPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(voteReceiveTimeDeviationAllColumns) == len(voteReceiveTimeDeviationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &VoteReceiveTimeDeviation{}
	if err = randomize.Struct(seed, o, voteReceiveTimeDeviationDBTypes, true, voteReceiveTimeDeviationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VoteReceiveTimeDeviation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := VoteReceiveTimeDeviations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, voteReceiveTimeDeviationDBTypes, true, voteReceiveTimeDeviationPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize VoteReceiveTimeDeviation struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testVoteReceiveTimeDeviationsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(voteReceiveTimeDeviationAllColumns) == len(voteReceiveTimeDeviationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &VoteReceiveTimeDeviation{}
	if err = randomize.Struct(seed, o, voteReceiveTimeDeviationDBTypes, true, voteReceiveTimeDeviationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VoteReceiveTimeDeviation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := VoteReceiveTimeDeviations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, voteReceiveTimeDeviationDBTypes, true, voteReceiveTimeDeviationPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize VoteReceiveTimeDeviation struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(voteReceiveTimeDeviationAllColumns, voteReceiveTimeDeviationPrimaryKeyColumns) {
		fields = voteReceiveTimeDeviationAllColumns
	} else {
		fields = strmangle.SetComplement(
			voteReceiveTimeDeviationAllColumns,
			voteReceiveTimeDeviationPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := VoteReceiveTimeDeviationSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testVoteReceiveTimeDeviationsUpsert(t *testing.T) {
	t.Parallel()

	if len(voteReceiveTimeDeviationAllColumns) == len(voteReceiveTimeDeviationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := VoteReceiveTimeDeviation{}
	if err = randomize.Struct(seed, &o, voteReceiveTimeDeviationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize VoteReceiveTimeDeviation struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert VoteReceiveTimeDeviation: %s", err)
	}

	count, err := VoteReceiveTimeDeviations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, voteReceiveTimeDeviationDBTypes, false, voteReceiveTimeDeviationPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize VoteReceiveTimeDeviation struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert VoteReceiveTimeDeviation: %s", err)
	}

	count, err = VoteReceiveTimeDeviations().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
